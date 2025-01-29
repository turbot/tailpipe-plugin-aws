package config

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/rs/dnscache"
	"golang.org/x/sync/semaphore"
)

const PluginName = "aws"

type AwsConnection struct {
	Profile               *string `hcl:"profile"`
	AccessKey             *string `hcl:"access_key"`
	SecretKey             *string `hcl:"secret_key"`
	SessionToken          *string `hcl:"session_token"`
	MaxErrorRetryAttempts *int    `hcl:"max_error_retry_attempts"`
	MinErrorRetryDelay    *int    `hcl:"min_error_retry_delay"`
	EndpointUrl           *string `hcl:"endpoint_url"`
	S3ForcePathStyle      *bool   `hcl:"s3_force_path_style"`
}

func (c *AwsConnection) Validate() error {
	if c.AccessKey != nil && c.SecretKey == nil {
		return fmt.Errorf("access_key set without secret_key")
	}

	if c.AccessKey == nil && c.SecretKey != nil {
		return fmt.Errorf("secret_key set without access_key")
	}

	if c.MinErrorRetryDelay != nil && *c.MinErrorRetryDelay < 1 {
		return fmt.Errorf("min_error_retry_delay must be greater than or equal to 10")
	}

	if c.MaxErrorRetryAttempts != nil && *c.MaxErrorRetryAttempts < 1 {
		return fmt.Errorf("max_error_retry_attempts must be greater than or equal to 1")
	}

	return nil
}

func (c *AwsConnection) Identifier() string {
	return PluginName
}

func (c *AwsConnection) GetClientConfiguration(ctx context.Context, overrideRegion *string) (*aws.Config, error) {
	var configOptions []func(*config.LoadOptions) error

	// profile
	if c.Profile != nil {
		profile := aws.ToString(c.Profile)
		configOptions = append(configOptions, config.WithSharedConfigProfile(profile))
	}

	// access keys
	if c.AccessKey != nil && c.SecretKey != nil {
		accessKey := aws.ToString(c.AccessKey)
		secretKey := aws.ToString(c.SecretKey)
		sessionToken := ""
		if c.SessionToken != nil {
			sessionToken = aws.ToString(c.SessionToken)
		}
		provider := credentials.NewStaticCredentialsProvider(accessKey, secretKey, sessionToken)
		configOptions = append(configOptions, config.WithCredentialsProvider(provider))
	}

	// shared http client
	configOptions = append(configOptions, config.WithHTTPClient(sharedHTTPClient))

	// load base config
	cfg, err := config.LoadDefaultConfig(ctx, configOptions...)
	if err != nil {
		return nil, fmt.Errorf("error loading AWS config: %w", err)
	}

	// if no region from base config, apply default region
	if overrideRegion != nil {
		cfg.Region = *overrideRegion
	} else if cfg.Region == "" {
		cfg.Region = "us-east-1"
	}

	// retry handling
	maxRetries := getConfigOrEnvInt(c.MaxErrorRetryAttempts, "AWS_MAX_ATTEMPTS", 9)
	var minRetryDelay = 25 * time.Millisecond
	if c.MinErrorRetryDelay != nil {
		minRetryDelay = time.Duration(*c.MinErrorRetryDelay) * time.Millisecond
	}

	retryer := retry.NewStandard(func(o *retry.StandardOptions) {
		// resetting state of rand to generate different random values
		rand.New(rand.NewSource(time.Now().UnixNano()))
		o.MaxAttempts = maxRetries
		o.MaxBackoff = 5 * time.Minute
		o.RateLimiter = NoOpRateLimit{} // With no rate limiter
		o.Backoff = NewExponentialJitterBackoff(minRetryDelay, maxRetries)
	})
	cfg.Retryer = func() aws.Retryer {
		// UnknownError is the code returned for a 408 from the aws go sdk, these can be frequent on large accounts especially around SNS Topics, etc.
		additionalErrors := []string{"UnknownError"}
		return retry.AddWithErrorCodes(retryer, additionalErrors...)
	}

	// custom endpoint
	endpointUrl := getConfigOrEnv(c.EndpointUrl, "AWS_ENDPOINT_URL")
	if endpointUrl != "" {
		customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{
				PartitionID:   "aws",
				URL:           endpointUrl,
				SigningRegion: region,
			}, nil
		})
		newCfg, err := config.LoadDefaultConfig(ctx, config.WithEndpointResolverWithOptions(customResolver))
		if err != nil {
			return nil, fmt.Errorf("error loading AWS config with custom endpoint resolver: %w", err)
		}
		newCfg.Retryer = cfg.Retryer
		newCfg.Region = cfg.Region
		cfg = newCfg
	}

	return &cfg, nil
}

// Helper function to get value from Config or environment variable
func getConfigOrEnv(configValue *string, env string) string {
	if configValue != nil {
		return *configValue
	}

	return os.Getenv(env)
}

func getConfigOrEnvInt(configValue *int, env string, defaultValue int) int {
	if configValue != nil {
		return *configValue
	}

	return readEnvVarToInt(env, defaultValue)
}

// ********************************************************************* //
// The below code is taken from Steampipe AWS SDK client initialization. //
// ********************************************************************* //

// Initialize a single HTTP client that is optimized for Steampipe and shared
// across all AWS SDK clients. We have hundreds of AWS SDK clients (one per
// account region) that are all sharing this same HTTP client - creating shared
// caching and controls over parallelism.
//
// The AWS SDK defaults are good, but not great for our highly parallel use in
// Steampipe. Specific problems this client aims to solve:
// 1. DNS floods - performing thousands of simultaneous API calls creates a DNS
// lookup for each one (even if the same domain). This can overwhelm the DNS
// server and cause "no such host" errors.
// 2. HTTP connection floods - the AWS SDK defaults to no limit on the number of
// HTTP connections per host. Thousands of connections created simultaneously to
// the same host is hard on both the client and the target server.
// 3. DNS caching - Golang does not cache DNS lookups by default. We end up
// looking up the same host thousands of times both within a query and across
// queries.
func initializeHTTPClient() aws.HTTPClient {

	// DNS lookup floods are a real problem with highly parallel AWS SDK calls. Every
	// API request leads to a DNS lookup by default (since Go doesn't cache them). We
	// employ a DNS lookup cache, but we also need to limit the number of parallel DNS
	// requests to avoid overwhelming the underlying DNS server. For example, listing
	// S3 buckets will create 2 DNS lookup requests per bucket which is a lot of
	// pressure on the DNS layer of your network.
	// This setting will limit the number of parallel DNS lookups. An appropriate setting
	// depends on the capabilities of your DNS server. The default is 25, which is low
	// enough for a Macbook M1 to work without "no such host" errors when using the cgo
	// network stack. It's high enough to work great in most cases, except maybe massive
	// S3 bucket listing (which is rare). Notably on the same Macbook M1, when the plugin
	// is compiled using netgo (our default on Mac) DNS lookups will succeed with virtually
	// no upper limit on this setting. So, bottom line, 25 is a guess to try and ensure
	// it works reliably and optimally enough.
	dnsLookupMaxParallel := readEnvVarToInt("TAILPIPE_AWS_DNS_LOOKUP_MAX_PARALLEL", 25)

	// The DNS cache will be refreshed at this interval. A refresh means that
	// any unused entries are removed and any entries that were used since the
	// last refresh will be re-looked up to ensure they are current.
	// This setting should be large enough to get the benefit of caching and short
	// enough to prevent stale entries from being used for too long.
	// Set to 0 to disable the refresh completely (not a good idea).
	// Set to -1 to disable the DNS cache completely (the AWS default).
	dnsCacheRefreshIntervalSecs := readEnvVarToInt("TAILPIPE_AWS_DNS_CACHE_REFRESH_INTERVAL_SECS", 300)

	// This is the maximum number of HTTPS API connections used for each host
	// (e.g.  iam.amazonaws.com). We want a number that is high enough to do a
	// lot of parallel work, but not so high that we have an excess number of
	// sockets open.
	// There is a trade off here. Tables like S3 have a lot of hosts - i.e. two
	// per bucket (one for the central region to get the creation time and one
	// for the actual bucket region), while services like IAM use a single host
	// for all queries.
	// Set to 0 to remove the limit (which is the AWS SDK default).
	httpTransportMaxConnsPerHost := readEnvVarToInt("TAILPIPE_AWS_HTTP_TRANSPORT_MAX_CONNS_PER_HOST", 5000)

	// Our DNS resolver should automatically refresh itself on this schedule.
	var resolver = &dnscache.Resolver{}
	if dnsCacheRefreshIntervalSecs > 0 {
		go func() {
			t := time.NewTicker(time.Duration(dnsCacheRefreshIntervalSecs) * time.Second)
			defer t.Stop()
			for range t.C {
				resolver.Refresh(true)
			}
		}()
	}

	// The AWS SDK has a special "buildable" HTTP client so it can be combined
	// with specific options such as custom certificate bundles. It matches the
	// interface of a HTTPClient, but has specific approaches for setting
	// transport options etc. Our goal is to use the default AWS settings (e.g.
	// timeouts, etc) as much as possible and just override the specific
	// behavior of parallelism for DNS lookups and HTTP requests.
	client := awshttp.NewBuildableClient()

	// Limit the max connections per host, but only if set. The AWS SDK default
	// is no limit.
	if httpTransportMaxConnsPerHost > 0 {
		client = client.WithTransportOptions(func(tr *http.Transport) {
			tr.MaxConnsPerHost = httpTransportMaxConnsPerHost
		})
	}

	// Use a DNS cache if it's set, otherwise we just avoid changing the dialer behavior
	// of the AWS HTTP client.
	if dnsCacheRefreshIntervalSecs >= 0 {

		// A semaphore is used to control the number of parallel DNS lookups.
		sem := semaphore.NewWeighted(int64(dnsLookupMaxParallel))

		// A dialer for testing connections
		dialer := client.GetDialer()

		client = client.WithTransportOptions(func(tr *http.Transport) {
			tr.DialContext = func(ctx context.Context, network string, addr string) (conn net.Conn, err error) {

				host, port, err := net.SplitHostPort(addr)
				if err != nil {
					return nil, err
				}

				// Acquire a semaphore slot, blocking until one is available.
				if err := sem.Acquire(ctx, 1); err != nil {
					return nil, err
				}

				// Actually resolve the host, using a cached result if possible.
				// Returns an array of IPs for the host.
				ips, err := resolver.LookupHost(ctx, host)

				// Release the semaphore, even if there was an error.
				sem.Release(1)

				// If there was an error during lookup, we give up immediately.
				if err != nil {
					return nil, err
				}

				// Now, look through the IP addresses until we manage to create a good connection.
				// This is less optimal than the parallelized native golang approach, but good
				// enough and much simpler. Comparison - https://cs.opensource.google/go/go/+/refs/tags/go1.21.5:src/net/dial.go;l=454-507
				for _, ip := range ips {
					conn, err = dialer.DialContext(ctx, network, net.JoinHostPort(ip, port))
					if err == nil {
						break
					}
				}

				return
			}
		})
	}

	return client
}

var sharedHTTPClient = initializeHTTPClient()

// Helper function for integer based environment variables.
func readEnvVarToInt(name string, defaultVal int) int {
	val := defaultVal
	envValue := os.Getenv(name)
	if envValue != "" {
		i, err := strconv.Atoi(envValue)
		if err == nil {
			val = i
		}
	}
	return val
}

// Given a region (including wildcards), guess at the best last resort region
// based on the partition. Examples:
//
//	us-gov-* -> us-gov-west-1
//	cn* -> cn-northwest-1
//	us-west-2 -> us-east-1
//	* -> us-east-1
//	crap -> ""
func awsLastResortRegionFromRegionWildcard(regionWildcard string) string {

	// Check prefixes for obscure partitions
	if strings.HasPrefix(regionWildcard, "us-gov") {
		return "us-gov-west-1"
	} else if strings.HasPrefix(regionWildcard, "cn") {
		return "cn-northwest-1"
	} else if strings.HasPrefix(regionWildcard, "us-isob") {
		return "us-isob-east-1"
	} else if strings.HasPrefix(regionWildcard, "us-iso") {
		return "us-iso-east-1"
	}

	// Check if the prefix is for a commercial region.
	// Must be done after obscure partitions, because they have the same
	// prefixes but longer.
	for _, prefix := range awsCommercialRegionPrefixes() {
		if strings.HasPrefix(regionWildcard, prefix) {
			return "us-east-1"
		}
	}

	// Unknown partition
	return ""
}

//
// AWS STANDARD REGIONS
//
// Source: https://docs.aws.amazon.com/general/latest/gr/rande.html#regional-endpoints
//
// Maintain a hard coded list of regions to use when API calls to the region
// list endpoint are not possible. This list must be updated manually as new
// regions are announced.
//

func awsCommercialRegionPrefixes() []string {
	return []string{
		"af",
		"ap",
		"ca",
		"eu",
		"me",
		"sa",
		"us",
	}
}

// NoOpRateLimit https://github.com/aws/aws-sdk-go-v2/issues/543
type NoOpRateLimit struct{}

func (NoOpRateLimit) AddTokens(uint) error { return nil }
func (NoOpRateLimit) GetToken(context.Context, uint) (func() error, error) {
	return noOpToken, nil
}
func noOpToken() error { return nil }

// ExponentialJitterBackoff provides backoff delays with jitter based on the
// number of attempts.
type ExponentialJitterBackoff struct {
	minDelay           time.Duration
	maxBackoffAttempts int
}

// NewExponentialJitterBackoff returns an ExponentialJitterBackoff configured
// for the max backoff.
func NewExponentialJitterBackoff(minDelay time.Duration, maxAttempts int) *ExponentialJitterBackoff {
	return &ExponentialJitterBackoff{minDelay, maxAttempts}
}

// BackoffDelay returns the duration to wait before the next attempt should be
// made. Returns an error if unable get a duration.
func (j *ExponentialJitterBackoff) BackoffDelay(attempt int, err error) (time.Duration, error) {
	minDelay := j.minDelay

	// The calculated jitter will be between [0.8, 1.2)
	var jitter = float64(rand.Intn(120-80)+80) / 100

	retryTime := time.Duration(int(float64(int(minDelay.Nanoseconds())*int(math.Pow(3, float64(attempt)))) * jitter))

	// Cap retry time at 5 minutes to avoid too long a wait
	if retryTime > (5 * time.Minute) {
		retryTime = time.Duration(5 * time.Minute)
	}

	// Low level method to log retries since we don't have context etc here.
	// Logging is helpful for visibility into retries and choke points in using
	// the API.
	slog.Info("BackoffDelay:", "attempt", attempt, "retry_time", retryTime.String(), "error", err)

	return retryTime, nil
}
