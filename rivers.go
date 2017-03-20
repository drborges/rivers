package rivers

import (
	"github.com/drborges/rivers/ctxtree"
	"github.com/drborges/rivers/stream"
)

// Producer represents the source of a rivers rivers. Bound to the given context,
// the producer produces data until there is no more data to be produced or until
// the context is closed by any of its downstreams, due to an error or becase no
// further data is required.
type Producer func(ctx ctxtree.Context) stream.Reader

// Transformer represents intermidiary stages of the pipeline, responsible for
// applying a transformation function to every item flowing through the stream,
// forwarding the result to the downstream.
type Transformer func(stream.Reader) stream.Reader

// Consumer the final stage of the pipeline, which consumes all stream items that
// make it to this final stage.
type Consumer func(stream.Reader)

// Splitter represents a stage in the pipeline where the stream is splitted into
// two new ones.
type Splitter func(stream.Reader) (stream.Reader, stream.Reader)

// Aggregator implements an aggregation between two streams returning a new
// stream.
type Aggregator func(stream.Reader, stream.Reader) stream.Reader

// Predicate is a function that given an input it returns true if the input matches
// the predicate, false otherwise.
type Predicate func(stream.T) bool

type Runner func() error

func (fn Runner) Run() error {
	return fn()
}

type Pipeline interface {
	Apply(Transformer) Pipeline
	Then(Consumer) Runner
}

func From(p Producer) Pipeline {
	reader, runner := NewRunner(p)
	return &pipeline{
		upstream: reader,
		runner:   runner,
	}
}

type pipeline struct {
	upstream stream.Reader
	runner   Runner
}

func (pipe *pipeline) Apply(t Transformer) Pipeline {
	pipe.upstream = t(pipe.upstream)
	return pipe
}

func (pipe *pipeline) Then(c Consumer) Runner {
	c(pipe.upstream)
	return pipe.runner
}

func NewRunner(producer Producer) (stream.Reader, Runner) {
	ctx := ctxtree.New()
	run := make(chan struct{})
	runner := func() error {
		close(run)
		<-ctx.Done()
		return ctx.Err()
	}

	reader := func(ctx ctxtree.Context) stream.Reader {
		reader, writer := stream.New(ctx.NewChild())
		upstream := producer(ctx)

		go func() {
			defer writer.Close(nil)

			<-run
			for data := range upstream.Read() {
				if err := writer.Write(data); err != nil {
					return
				}
			}
		}()

		return reader
	}(ctx)

	return reader, runner
}
