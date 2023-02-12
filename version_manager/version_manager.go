package versionmanager

import "context"

type VersionManager interface {
	Get(ctx context.Context, key string) (string, error)
	Create(ctx context.Context, key string, opts ...Option) (string, error)
	Update(ctx context.Context, key string, opts ...Option) (string, error)
	Delete(ctx context.Context, key string) error
}

type Options struct {
	params map[string]any
}

func (o *Options) GetString(key string) string {
	v, ok := o.params[key]
	if !ok {
		return ""
	}

	vv, ok := v.(string)
	if !ok {
		return ""
	}
	return vv
}

func (o *Options) Apply(opts ...Option) {
	for _, opt := range opts {
		opt.Apply(o)
	}
}

type Option interface {
	Apply(*Options)
}

type withParamOption struct {
	key   string
	value any
}

func (o *withParamOption) Apply(opts *Options) {
	opts.params[o.key] = o.value
}

func WithParam(key string, value any) Option {
	return &withParamOption{key, value}
}
