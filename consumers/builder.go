package consumers

import "github.com/drborges/rivers/rx"

type Builder struct {
	context rx.Context
}

func New(c rx.Context) *Builder {
	return &Builder{c}
}

func (builder *Builder) Drainer() rx.Consumer {
	return &drainer{builder.context}
}

func (builder *Builder) DataCollector(data *[]rx.T) rx.Consumer {
	return &dataCollector{
		context: builder.context,
		data: data,
	}
}