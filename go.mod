module github.com/turbot/tailpipe-plugin-aws

go 1.24

toolchain go1.24.1

//replace github.com/turbot/tailpipe-plugin-sdk => ../tailpipe-plugin-sdk

require (
	github.com/aws/aws-sdk-go v1.55.7
	github.com/aws/aws-sdk-go-v2 v1.36.5
	github.com/aws/aws-sdk-go-v2/config v1.29.16
	github.com/aws/aws-sdk-go-v2/credentials v1.17.69
	github.com/aws/aws-sdk-go-v2/feature/s3/manager v1.17.75
	github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs v1.50.3
	github.com/aws/aws-sdk-go-v2/service/guardduty v1.54.5
	github.com/aws/aws-sdk-go-v2/service/s3 v1.79.3
	github.com/aws/aws-sdk-go-v2/service/securityhub v1.57.4
	github.com/elastic/go-grok v0.3.1
	github.com/hashicorp/hcl/v2 v2.23.0
	github.com/rs/dnscache v0.0.0-20230804202142-fc85eb664529
	github.com/rs/xid v1.6.0
	github.com/stoewer/go-strcase v1.3.0
	github.com/turbot/go-kit v1.3.0
	github.com/turbot/pipe-fittings/v2 v2.6.0
	github.com/turbot/tailpipe-plugin-sdk v0.9.2
	golang.org/x/sync v0.15.0
)

require (
	cloud.google.com/go v0.115.0 // indirect
	cloud.google.com/go/auth v0.7.2 // indirect
	cloud.google.com/go/auth/oauth2adapt v0.2.3 // indirect
	cloud.google.com/go/compute/metadata v0.5.2 // indirect
	cloud.google.com/go/iam v1.1.10 // indirect
	cloud.google.com/go/storage v1.42.0 // indirect
	github.com/AdaLogics/go-fuzz-headers v0.0.0-20230811130428-ced1acdcaa24 // indirect
	github.com/Masterminds/semver/v3 v3.2.1 // indirect
	github.com/acarl005/stripansi v0.0.0-20180116102854-5a71ef0e047d // indirect
	github.com/agext/levenshtein v1.2.3 // indirect
	github.com/apache/arrow-go/v18 v18.1.0 // indirect
	github.com/apparentlymart/go-cidr v1.1.0 // indirect
	github.com/apparentlymart/go-textseg/v15 v15.0.0 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.6.11 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.16.31 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.36 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.36 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.3 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.3.35 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.12.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.7.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.12.16 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.18.16 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.25.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.30.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.33.21 // indirect
	github.com/aws/smithy-go v1.22.4 // indirect
	github.com/bgentry/go-netrc v0.0.0-20140422174119-9fd32a8b3d3d // indirect
	github.com/bgentry/speakeasy v0.1.0 // indirect
	github.com/bmatcuk/doublestar v1.3.4 // indirect
	github.com/briandowns/spinner v1.23.0 // indirect
	github.com/btubbs/datetime v0.1.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/containerd/containerd v1.7.27 // indirect
	github.com/containerd/errdefs v0.3.0 // indirect
	github.com/containerd/log v0.1.0 // indirect
	github.com/containerd/platforms v0.2.1 // indirect
	github.com/danwakefield/fnmatch v0.0.0-20160403171240-cbb64ac3d964 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/dgraph-io/ristretto v0.2.0 // indirect
	github.com/duckdb/duckdb-go-bindings v0.1.13 // indirect
	github.com/duckdb/duckdb-go-bindings/darwin-amd64 v0.1.8 // indirect
	github.com/duckdb/duckdb-go-bindings/darwin-arm64 v0.1.8 // indirect
	github.com/duckdb/duckdb-go-bindings/linux-amd64 v0.1.8 // indirect
	github.com/duckdb/duckdb-go-bindings/linux-arm64 v0.1.8 // indirect
	github.com/duckdb/duckdb-go-bindings/windows-amd64 v0.1.8 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/fatih/color v1.17.0 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/fsnotify/fsnotify v1.8.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.3 // indirect
	github.com/gertd/go-pluralize v0.2.1 // indirect
	github.com/go-git/go-git/v5 v5.13.0 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.20.0 // indirect
	github.com/go-viper/mapstructure/v2 v2.3.0 // indirect
	github.com/goccy/go-json v0.10.5 // indirect
	github.com/goccy/go-yaml v1.11.2 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/flatbuffers v25.1.24+incompatible // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/google/s2a-go v0.1.7 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.3.2 // indirect
	github.com/googleapis/gax-go/v2 v2.13.0 // indirect
	github.com/gopherjs/gopherjs v1.17.2 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-getter v1.7.5 // indirect
	github.com/hashicorp/go-hclog v1.6.3 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-plugin v1.6.1 // indirect
	github.com/hashicorp/go-safetemp v1.0.0 // indirect
	github.com/hashicorp/go-uuid v1.0.3 // indirect
	github.com/hashicorp/go-version v1.7.0 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/hashicorp/terraform-registry-address v0.2.1 // indirect
	github.com/hashicorp/terraform-svchost v0.1.1 // indirect
	github.com/hashicorp/yamux v0.1.1 // indirect
	github.com/hokaccha/go-prettyjson v0.0.0-20211117102719-0474bc63780f // indirect
	github.com/iancoleman/strcase v0.3.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.14.3 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.3.3 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.6.0 // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	github.com/jedib0t/go-pretty/v6 v6.5.9 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/jtolds/gls v4.20.0+incompatible // indirect
	github.com/karrick/gows v0.3.0 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/klauspost/cpuid/v2 v2.2.9 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/logrusorgru/aurora v2.0.3+incompatible // indirect
	github.com/magefile/mage v1.15.0 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/marcboeker/go-duckdb/arrowmapping v0.0.6 // indirect
	github.com/marcboeker/go-duckdb/mapping v0.0.6 // indirect
	github.com/marcboeker/go-duckdb/v2 v2.1.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.15 // indirect
	github.com/mitchellh/colorstring v0.0.0-20190213212951-d06e56a500db // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/go-testing-interface v1.14.1 // indirect
	github.com/mitchellh/go-wordwrap v1.0.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/moby/locker v1.0.1 // indirect
	github.com/oklog/run v1.0.0 // indirect
	github.com/olekukonko/tablewriter v0.0.5 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.1.0 // indirect
	github.com/pelletier/go-toml/v2 v2.2.2 // indirect
	github.com/pierrec/lz4/v4 v4.1.22 // indirect
	github.com/pjbgf/sha1cd v0.3.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/sagikazarmark/locafero v0.4.0 // indirect
	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
	github.com/satyrius/gonx v1.4.0 // indirect
	github.com/sethvargo/go-retry v0.3.0 // indirect
	github.com/shiena/ansicolor v0.0.0-20230509054315-a9deabde6e02 // indirect
	github.com/shirou/gopsutil v3.21.11+incompatible // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/smarty/assertions v1.16.0 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/afero v1.11.0 // indirect
	github.com/spf13/cast v1.6.0 // indirect
	github.com/spf13/cobra v1.8.1 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/spf13/viper v1.19.0 // indirect
	github.com/stevenle/topsort v0.2.0 // indirect
	github.com/stretchr/testify v1.10.0 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	github.com/tklauser/go-sysconf v0.3.9 // indirect
	github.com/tklauser/numcpus v0.3.0 // indirect
	github.com/tkrajina/go-reflector v0.5.8 // indirect
	github.com/turbot/pipes-sdk-go v0.12.0 // indirect
	github.com/turbot/terraform-components v0.0.0-20231213122222-1f3526cab7a7 // indirect
	github.com/ulikunitz/xz v0.5.10 // indirect
	github.com/xlab/treeprint v1.2.0 // indirect
	github.com/yusufpapurcu/wmi v1.2.2 // indirect
	github.com/zclconf/go-cty v1.14.4 // indirect
	github.com/zclconf/go-cty-yaml v1.0.3 // indirect
	github.com/zeebo/xxh3 v1.0.2 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.49.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.49.0 // indirect
	go.opentelemetry.io/otel v1.31.0 // indirect
	go.opentelemetry.io/otel/metric v1.31.0 // indirect
	go.opentelemetry.io/otel/trace v1.31.0 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.9.0 // indirect
	golang.org/x/crypto v0.36.0 // indirect
	golang.org/x/exp v0.0.0-20250128182459-e0ece0dbea4c // indirect
	golang.org/x/mod v0.22.0 // indirect
	golang.org/x/net v0.38.0 // indirect
	golang.org/x/oauth2 v0.27.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/term v0.30.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	golang.org/x/time v0.5.0 // indirect
	golang.org/x/tools v0.29.0 // indirect
	golang.org/x/xerrors v0.0.0-20240903120638-7835f813f4da // indirect
	google.golang.org/api v0.189.0 // indirect
	google.golang.org/genproto v0.0.0-20240722135656-d784300faade // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20241015192408-796eee8c2d53 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241104194629-dd2ea8efbc28 // indirect
	google.golang.org/grpc v1.69.2 // indirect
	google.golang.org/protobuf v1.36.1 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	oras.land/oras-go/v2 v2.5.0 // indirect
	sigs.k8s.io/yaml v1.4.0 // indirect
)
