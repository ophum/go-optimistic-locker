package optimisticlocker

type PreconditionCheckOption interface {
	apply(*PreconditionCheckOptions)
}

type PreconditionCheckOptions struct {
	EtagGenerator EtagGenerator
}

type etagGeneratorOption struct {
	etagGenerator EtagGenerator
}

func (o *etagGeneratorOption) apply(opts *PreconditionCheckOptions) {
	opts.EtagGenerator = o.etagGenerator
}

func WithEtagGenerator(generator EtagGenerator) PreconditionCheckOption {
	return &etagGeneratorOption{generator}
}
