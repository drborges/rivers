# Rivers ![Basic Stream](https://raw.githubusercontent.com/drborges/rivers/master/https://raw.githubusercontent.com/drborges/rivers/master/docs/rivers-logo.png)

[![Build Status](https://travis-ci.org/drborges/rivers.svg?branch=master)](https://travis-ci.org/drborges/rivers)

Data Stream Processing API for GO

# Overview

Rivers provides a simple though powerful API for processing streams of data built on top of `goroutines`, `channels` and the [pipeline pattern](https://blog.golang.org/pipelines).

```go
err := rivers.From(NewGithubRepositoryProducer(httpClient)).
	Filter(hasForks).
	Filter(hasRecentActivity).
	Drop(ifStarsCountIsLessThan(50)).
	Map(extractAuthorInformation).
	Batch(200).
	Each(saveBatch).
	Drain()
```

With a few basic building blocks based on the `Producer-Consumer` model, you can compose and create complex data processing pipelines for solving a variety of problems.

A pipeline often assumes the following format: `Producer`, one or more `Transformers` and an optional `Consumer`

![Basic Stream](https://raw.githubusercontent.com/drborges/rivers/master/docs/stream-transformation.png)

For more complex formats check the `Combiners` and `Dispatchers` sections.

# Building Blocks

A particular stream pipeline may be built composing building blocks such as `producers`,  `consumers`,  `transformers`, `combiners` and `dispatchers`.

### Stream ![Basic Stream](https://raw.githubusercontent.com/drborges/rivers/master/docs/stream.png)

Streams are simply readable or writable channels where data flows through `asynchronously`. They are usually created by `producers` providing data from a particular data source, for example `files`, `network` (socket data, API responses), or even as simple as regular `slice` of data to be processed.

Rivers provides a `stream` package with a constructor function for creating streams as follows:

```go
capacity := 100
readable, writable := stream.New(capacity)
```

Streams are buffered and the `capacity` parameter dictates how many items can be produced into the stream without being consumed until the producer is blocked. This blocking mechanism is natively implemented by Go channels, and is a form of `back-pressuring` the pipeline.

### Producers ![Basic Stream](https://raw.githubusercontent.com/drborges/rivers/master/docs/producer.png)

Asynchronously emits data into a stream. Any type implementing the `stream.Producer` interface can be used as a producer in rivers.

```go
type Producer interface {
	Produce() stream.Readable
}
```

Producers implement the [pipeline pattern](https://blog.golang.org/pipelines) in order to asynchronously produce items that will be eventually consumed by a further stage in the pipeline. Rivers provides a few implementations of producers such as:

- `rivers.FromRange(0, 1000)`
- `rivers.FromSlice(slice)`
- `rivers.FromData(1, 2, "a", "b", Person{Name:"Diego"})`
- `rivers.FromFile(aFile).ByLine()`
- `rivers.FromSocket("tcp", ":8484")`

A good producer implementation takes care of at least 3 important aspects:

1. Checks if rivers context is still opened before emitting any item
2. Defers the recover function from rivers context as part of the goroutine execution (For more see `Cancellation` topic)
3. Closes the writable stream at the end of the go routine. By closing the channel further stages of the pipeline know when their work is done.

Lets see how one would go about converting a slice of numbers into a stream with a simple Producer implementation:

```go
type NumbersProducer struct {
	context stream.Context
	numbers []int
}

func (producer *NumbersProducer) Produce() stream.Readable {
	readable, writable := stream.New(len(producer.numbers))

	go func() {
		defer producer.context.Recover()
		defer close(writable)

		for _, n := range producer.numbers {
			select {
			case <-producer.context.Closed():
				return
			default:
				writable <- n
			}
		}
	}()

	return readable
}
```

The code above is a complaint `rivers.Producer` implementation and it gives the developer full control of the process. Rivers also provides an `Observable` type that implements `stream.Producer` covering the basic 3 aspects mentioned above that you can use for most cases: `producers.Observable`.

Our producer implementation in terms of an observable would then look like:

```go
func NewNumbersProducer(numbers []int) stream.Producer {
	return &Observable{
		Capacity: len(numbers),
		Emit: func(emitter stream.Emitter) {
			for _, n := range numbers {
				emitter.Emit(n)
			}
		},
	}
}
```

### Consumers ![Basic Stream](https://raw.githubusercontent.com/drborges/rivers/master/docs/consumer.png)

Consumes data from a particular stream. Consumers block the process until there is no more data to be consumed out of the stream.

You can use consumers to collect the items reaching the end of the pipeline, or any errors that might have happened during the execution.

It is very likely you will most often need a final consumer in your pipeline for waiting for the pipeline result before moving on.

Rivers has a few built-in consumers, among them you will find:

1. `Drainers` which block draining the stream until there is no more data flowing through and returning any possible errors.

2. `Collectors` collect all items that reached the end of the pipeline and any possible error.

Say we have a stream where instances of `Person` are flowing through, then you can collect items off the stream like so:

```go
type Person struct {
	Name string
}

diego := Person{"Diego"}
borges := Person{"Borges"}

items, err := rivers.FromData(diego, borges).Collect()
// items == []stream.T{Person{"Diego"}, Person{"Borges"}}

item, err := rivers.FromData(diego, borges).CollectFirst()
// item == Person{"Diego"}

item, err := rivers.FromData(diego, borges).CollectLast()
// item == Person{"Borges"}

var people []Person
err := rivers.FromData(diego, borges).CollectAs(&people)
// people == []Person{{"Diego"}, {"Borges"}}

var diego Person
err := rivers.FromData(diego, borges).CollectFirstAs(&diego)

var borges Person
err := rivers.FromData(diego, borges).CollectLastAs(&diego)
```

### Transformers ![Dispatching To Streams](https://raw.githubusercontent.com/drborges/rivers/master/docs/transformer.png)

Reads data from a particular stream applying a transformation function to it, optionally forwarding the result to an output channel. Transformers implement the interface `stream.Transformer`

```go
type Transformer interface {
	Transform(in stream.Readable) (out stream.Readable)
}
```

There are a variety of transform operations built-in in rivers, to name a few: `Map`, `Filter`, `Each`, `Flatten`, `Drop`, `Take`, etc...

Basic Stream Transformation Pipeline: `Producer -> Transformer -> Consumer`

![Basic Stream](https://raw.githubusercontent.com/drborges/rivers/master/docs/stream-transformation.png)

Aiming extensibility, rivers allow you to implement your own version of `stream.Transformer`. The following code implements a `filter` in terms of `stream.Transformer`:

```go
type Filter struct {
	context stream.Context
	fn      stream.PredicateFn
}

func (filter *Filter) Transform(in stream.Readable) stream.Readable {
	readable, writable := stream.New(in.Capacity())

	go func() {
		defer filter.context.Recover()
		defer close(writable)

		for {
			select {
			case <-filter.context.Closed():
				return
			default:
				data, more := <-in
				if !more {
					return
				}
				if filter.fn(data) {
					writable <- data
				}
			}
		}
	}()

	return readable
}
```

Note that the transformer above also takes care of those `3 aspects` mentioned in the `producer` implementation. You could use this transformer like so:

```go
stream := rivers.FromRange(1, 10)

evensOnly := func(data stream.T) bool {
	return data.(int) % 2 == 0
}

filter := &Filter{stream.Context, evensOnly}

evens, err := stream.Apply(filter).Collect()
```

In order to reduce some of the boilerplate, rivers provides a generic implementation of `stream.Transformer` that you can use to implement many use cases: `transformers.Observer`. The filter above can be rewritten as:

```go
func NewFilter(fn stream.PredicateFn) stream.Transformer {
	return &Observer{
		OnNext: func(data stream.T, emitter stream.Emitter) error {
			if fn(data) {
				emitter.Emit(data)
			}
			return nil
		},
	}
}

evens, err := rivers.FromRange(1, 10).
	Apply(NewFilter(evensOnly)).
	Collect()
```

The observer `OnNext` function may return `stream.Done` in order to finish its work explicitly stopping the pipeline. This is useful for implementing short-circuit operations such as `find`, `any`, etc...

### Combiners ![Dispatching To Streams](https://raw.githubusercontent.com/drborges/rivers/master/docs/combiner.png)

Combining streams is often a useful operation and rivers makes it easy with its pre-baked combiner implementations `FIFO`, `Zip` and `ZipBy`. A combiner implements `stream.Combiner` interface:

```go
type Combiner interface {
	Combine(in ...Readable) (out Readable)
}
```

It essentially takes one or more readable streams and gives you back a readable stream which is the result of combining all the inputs.

Combining Streams Pipeline: `Producers -> Combiner -> Transformer -> Consumer`

![Combining Streams](https://raw.githubusercontent.com/drborges/rivers/master/docs/stream-combiner.png)

The following example combines data from 3 different streams into a single stream:

```go
facebookMentions := rivers.From(FacebookPostsProducer(httpClient)).
	Map(ExtractMentionsTo("diego.rborges"))

twitterMentions := rivers.From(TwitterFeedProducer(httpClient)).
	Map(ExtractMentionsTo("dr_borges")).Stream

githubMentions := rivers.From(GithubRiversCommitsProducer(httpClient)).
	Map(ExtractMentionsTo("drborges")).Stream

err := facebookMentions.Merge(twitterMentions, githubMentions).
	Take(mentionsFrom3DaysAgo).
	Each(prettyPrint).
	Drain()
```

### Dispatchers ![Dispatching To Streams](https://raw.githubusercontent.com/drborges/rivers/master/docs/dispatcher.png)

Dispatchers forward data from a readable stream to one or more writable streams returning a new readable stream from where the  non dispatched data can still be processed.

Dispatchers may conditionally dispatch data based on a `stream.PreficateFn`. The `Partition` operation is an example of conditional dispatcher.

Rivers implement 3 types of dispatchers: `Split`, `SplitN` and `Partition`. Dispatchers implement `stream.Dispatcher` interface:

```go
type Dispatcher interface {
	Dispatch(from Readable, to ...Writable) (out Readable)
}
```

Dispatching data to multiple targets in a pipeline would look like: `Producer -> Dispatcher -> Transformers -> Consumers`

![Dispatching To Streams](https://raw.githubusercontent.com/drborges/rivers/master/docs/stream-dispatcher.png)

The following code show a use case for the `Partition` operation mentioned above:

```go
inactiveUsers, activeUsers := rivers.From(UserSessionsAPI(httpClient)).
	Partition(inactiveForOver7Days)
```

The example above forks the original stream into two others based on the given `predicate`. Data is then dynamically dispatched to either stream based on the predicate result.

Under the hood the partition operation makes use of a conditional dispatcher to redirect data to the resulting streams. For more detailed information on how to leverage the built-in dispatchers, take a look at the `dispatchers` package.  

# Examples

```go
evensOnly := func(data stream.T) bool { return data.(int) % 2 == 0 }
addOne := func(data stream.T) stream.T { return data.(int) + 1 }

data, err := rivers.FromRange(1, 10).
	Filter(evensOnly).
	Map(addOne).
	Collect()

fmt.Println("data:", data)
fmt.Println("err:", err)

// Output:
// data: []stream.T{3, 5, 7, 9, 11}
// err: nil
```

# Built-in Filters and Mappers

TODO

# The Cancellation Problem

Imagine you have a program that fires off a great deal of concurrent execution threads, in Go that translates to `goroutines`. Making sure each one of these goroutines come to an end by normally finishing their execution or more importantly canceling themselves upon any fatal failure in the pipeline is an interesting challenge. That is called the `cancellation` problem. 

Keeping track of all potential `goroutines` running in a system might not scale very well specially in scenarios where the number of concurrent goroutines might get to the hundreds of thousands. The concurrency model based on `goroutines` implemented by Go along with its primitives for working with `channels` and panic recovering allow developers to handle problems like this in a very elegant manner.

Rivers solves the `cancellation` problem by applying one of the properties of `closed channels`: [A closed channel is always ready to receive](http://golang.org/ref/spec#Receive_operator).

Every rivers operation before producing, consuming or transforming incoming data, first checks whether or not the context is `Done` within a [select channel block](https://gobyexample.com/non-blocking-channel-operations), take the `filter` example previously mentioned:

```go
go func() {
	defer filter.context.Recover()
	defer close(writable)

	select {
	case <-filter.context.Done():
		return
	default:
		data, more := <-in
		if !more {
			return
		}
		if filter.fn(data) {
			writable <- data
		}
	}
}()
```

Note that before it filters the incoming data -- default block -- it tries to read from the context `Done` channel. If it is able to read from it then the context Done channel is already closed therefore always read to receive. That causes the select block to peek the context done case, returning and finishing the goroutine right away.

Rivers will close the done channel whenever it recovers from a panic anywhere in the pipeline -- as long as the panicing goroutine defers `context.Recover()`.

Now imagine every single goroutine fired by rivers applying this pattern, the whole system continuation/cancellation logic can be controlled by simply closing or not a single channel. No need to keep track nor synchronize hundreds of thousands of execution threads, just a single close statement for the rescue. That is how powerful Go concurrency model is <3

# Going Parallel

Rivers parallelization support is an experimental feature and it is under active development. Currently the following operations suport parallelization: `Each`, `Map`, `Filter` and `OnData`.

When parallelizing an opearation, lets say `Each`, rivers will create a new parallel transformer (`Each` transformer in this case) for each slot in the pipeline capacity (a.k.a the producer's capacity), each of these transformers will be consuming data out of the same input readable stream having their output streams merged into a single pipeline where you can attach new pipeline operations to it.  

Enabling parallelization is feirly simple, take the following code for instance:

```go
facebookMentions := rivers.From(FacebookPostsProducer(httpClient)).
	Map(ExtractMentionsTo("diego.rborges"))

err := facebookMentions.Parallel().
	Drop(ifAlreadyInDatabase).
	Batch(500).
	Each(saveBatchInDatabase).
	Each(indexBatchInElasticsearch).
	Drain()
```

In the example above the `Drop`, and `Each` operations will be `parallelized`. Assuming the `FacebookPostsProducer` capacity is `1000` items then 1000 parallel transformers will be created for the `Drop` stage and a `1000` more parallel transformers for the `Each` stage.

# Troubleshooting

TODO `rivers.DebugEnabled`

# Future Work / Improvements

- [ ] Dynamic back-pressure?
- [ ] Monitoring?

# Contributing

TODO

# License

TODO