package plugin

import (
	"context"
	"errors"
	"sync"

	"github.com/turbot/tailpipe-plugin-aws/collection"
	"github.com/turbot/tailpipe-plugin-aws/source"

	sdkcollection "github.com/turbot/tailpipe-plugin-sdk/collection"
	"github.com/turbot/tailpipe-plugin-sdk/plugin"
	sdksource "github.com/turbot/tailpipe-plugin-sdk/source"
)

type AwsPlugin struct {
	ctx context.Context

	// observers is a list of observers that will be notified of events.
	observers      []plugin.PluginObserver
	observersMutex sync.RWMutex
}

func (p *AwsPlugin) Identifier() string {
	return "aws"
}

func (p *AwsPlugin) Init(ctx context.Context) error {
	p.ctx = ctx
	return nil
}

func (p *AwsPlugin) Context() context.Context {
	return p.ctx
}

func (p *AwsPlugin) Validate() error {
	return errors.ErrUnsupported
}

func (p *AwsPlugin) AddObserver(observer plugin.PluginObserver) {
	p.observersMutex.Lock()
	defer p.observersMutex.Unlock()
	p.observers = append(p.observers, observer)
}

func (p *AwsPlugin) RemoveObserver(observer plugin.PluginObserver) {
	p.observersMutex.Lock()
	defer p.observersMutex.Unlock()
	for i, o := range p.observers {
		if o == observer {
			p.observers = append(p.observers[:i], p.observers[i+1:]...)
			break
		}
	}
}

func (p *AwsPlugin) Sources() map[string]sdksource.Plugin {
	return map[string]sdksource.Plugin{
		"aws_s3_bucket": &source.AwsS3BucketSource{},
	}
}

func (p *AwsPlugin) Collections() map[string]sdkcollection.Plugin {
	return map[string]sdkcollection.Plugin{
		"aws_cloudtrail_log": &collection.AwsCloudTrailLogCollection{},
	}
}
