package rivers_test

import (
	"github.com/drborges/rivers"
	"github.com/drborges/rivers/scanners"
	"github.com/drborges/rivers/stream"
	. "github.com/smartystreets/goconvey/convey"
	"net"
	"strings"
	"testing"
)

func TestRiversAPI(t *testing.T) {
	toString := func(data stream.T) stream.T { return string(data.([]byte)) }
	nonEmptyLines := func(data stream.T) bool { return data.(string) != "" }
	splitWord := func(data stream.T) stream.T { return strings.Split(data.(string), " ") }
	evensOnly := func(data stream.T) bool { return data.(int)%2 == 0 }
	sum := func(a, b stream.T) stream.T { return a.(int) + b.(int) }
	add := func(n int) stream.MapFn {
		return func(data stream.T) stream.T { return data.(int) + n }
	}

	append := func(c string) stream.MapFn {
		return func(data stream.T) stream.T { return data.(string) + c }
	}

	addOrAppend := func(n int, c string) stream.MapFn {
		return func(data stream.T) stream.T {
			if num, ok := data.(int); ok {
				return num + n
			}
			if letter, ok := data.(string); ok {
				return letter + "_"
			}
			return data
		}
	}

	alphabeticOrder := func(a, b stream.T) bool {
		return a.(string) < b.(string)
	}

	listen := func() (net.Listener, string) {
		port := ":8282"
		ln, err := net.Listen("tcp", port)
		if err != nil {
			port = ":8383"
			ln, _ = net.Listen("tcp", port)
		}
		return ln, port
	}

	Convey("rivers API", t, func() {

		Convey("From Range -> Filter -> Map -> Reduce -> Each -> Sink", func() {
			data, _ := rivers.FromRange(1, 5).
				Filter(evensOnly).
				Map(add(1)).
				Reduce(0, sum).
				Collect()

			So(data, ShouldResemble, []stream.T{8})
		})

		Convey("From Data -> Flatten -> Map -> Sort By -> Batch -> Sink", func() {
			data, _ := rivers.FromData([]stream.T{"a", "c"}, "b", []stream.T{"d", "e"}).
				Flatten().
				Map(append("_")).
				SortBy(alphabeticOrder).
				Batch(2).
				Collect()

			So(data, ShouldResemble, []stream.T{
				[]stream.T{"a_", "b_"},
				[]stream.T{"c_", "d_"},
				[]stream.T{"e_"},
			})
		})

		Convey("From Slice -> Dispatch If -> Map -> Sink", func() {
			in, out := stream.New(2)

			notDispatched, _ := rivers.FromSlice([]stream.T{1, 2, 3, 4, 5}).
				DispatchIf(evensOnly, out).
				Map(add(2)).
				Collect()

			So(in.Read(), ShouldResemble, []stream.T{2, 4})
			So(notDispatched, ShouldResemble, []stream.T{3, 5, 7})
		})

		Convey("Zip -> Map -> Sink", func() {
			numbers := rivers.FromData(1, 2, 3, 4)
			letters := rivers.FromData("a", "b", "c")

			combined, _ := rivers.Zip(numbers, letters).
				Map(addOrAppend(1, "_")).Collect()

			So(combined, ShouldResemble, []stream.T{2, "a_", 3, "b_", 4, "c_", 5})
		})

		Convey("Zip By -> Map -> Sink", func() {
			numbers := rivers.FromData(1, 2, 3, 4)
			moreNumbers := rivers.FromData(4, 4, 1)

			combined, _ := rivers.ZipBy(sum, numbers, moreNumbers).Filter(evensOnly).Collect()

			So(combined, ShouldResemble, []stream.T{6, 4, 4})
		})

		Convey("Merge -> Map -> Sink", func() {
			numbers := rivers.FromData(1, 2, 3, 4)
			moreNumbers := rivers.FromData(4, 4, 1)

			combined, _ := rivers.Merge(numbers, moreNumbers).Collect()

			So(combined, ShouldResemble, []stream.T{1, 2, 3, 4, 4, 4, 1})
		})

		Convey("From Data -> Drain", func() {
			numbers := rivers.FromData(1, 2, 3, 4)
			numbers.Drain()

			data, opened := <-numbers.Sink()
			So(data, ShouldBeNil)
			So(opened, ShouldBeFalse)
		})

		Convey("From Socket -> Map -> Filter -> Map -> Flatten -> Sink", func() {
			ln, port := listen()

			go func() {
				conn, _ := ln.Accept()
				defer conn.Close()
				conn.Write([]byte("Hello there\n"))
				conn.Write([]byte("\n"))
				conn.Write([]byte("rivers!\n"))
				conn.Write([]byte("super cool!\n"))
			}()

			words := rivers.FromSocketWithScanner("tcp", port, scanners.NewLineScanner()).
				Map(toString).
				Filter(nonEmptyLines).
				Map(splitWord).
				Flatten().
				Sink()

			So(words.Read(), ShouldResemble, []stream.T{"Hello", "there", "rivers!", "super", "cool!"})
		})

		Convey("From Range -> Partition -> Sink", func() {
			evens, odds := rivers.FromRange(1, 10).Partition(evensOnly)

			So(evens.Sink().Read(), ShouldResemble, []stream.T{2, 4, 6, 8, 10})
			So(odds.Sink().Read(), ShouldResemble, []stream.T{1, 3, 5, 7, 9})
		})

		Convey("From Range -> Slipt -> Sink", func() {
			lhs, rhs := rivers.FromRange(1, 4).Split()

			So(lhs.Sink().Read(), ShouldResemble, []stream.T{1, 2, 3, 4})
			So(rhs.Sink().Read(), ShouldResemble, []stream.T{1, 2, 3, 4})
		})

		Convey("From Range -> Slipt N -> Sink", func() {
			streams := rivers.FromRange(1, 4).SplitN(3)

			So(streams[0].Sink().Read(), ShouldResemble, []stream.T{1, 2, 3, 4})
			So(streams[1].Sink().Read(), ShouldResemble, []stream.T{1, 2, 3, 4})
			So(streams[2].Sink().Read(), ShouldResemble, []stream.T{1, 2, 3, 4})
		})

		Convey("From Range -> OnData -> Sink", func() {
			processor := rivers.FromRange(1, 4).OnData(func(data stream.T, out stream.Writable) {
				if data.(int)%2 == 0 {
					out <- data
				}
			})

			So(processor.Sink().Read(), ShouldResemble, []stream.T{2, 4})
		})

		Convey("From Range -> Take -> Sink", func() {
			taken, _ := rivers.FromRange(1, 4).TakeFirst(2).Collect()

			So(taken, ShouldResemble, []stream.T{1, 2})
		})

		Convey("From Range -> TakeIf -> Sink", func() {
			processor := rivers.FromRange(1, 4).Take(evensOnly)

			So(processor.Sink().Read(), ShouldResemble, []stream.T{2, 4})
		})

		Convey("From Range -> DropIf -> Sink", func() {
			processor := rivers.FromRange(1, 4).Drop(evensOnly)

			So(processor.Sink().Read(), ShouldResemble, []stream.T{1, 3})
		})

		Convey("From Range -> Collect -> Sink", func() {
			data, err := rivers.FromRange(1, 4).Collect()

			So(err, ShouldBeNil)
			So(data, ShouldResemble, []stream.T{1, 2, 3, 4})
		})

		Convey("From Range -> CollectFirst -> Sink", func() {
			data, err := rivers.FromRange(1, 4).CollectFirst()

			So(err, ShouldBeNil)
			So(data, ShouldEqual, 1)
		})

		Convey("From Range -> CollectFirstAs -> Sink", func() {
			var data int
			err := rivers.FromRange(1, 4).CollectFirstAs(&data)

			So(err, ShouldBeNil)
			So(data, ShouldEqual, 1)
		})

		Convey("From Range -> CollectLast -> Sink", func() {
			data, err := rivers.FromRange(1, 4).CollectLast()

			So(err, ShouldBeNil)
			So(data, ShouldEqual, 4)
		})

		Convey("From Range -> CollectLastAs -> Sink", func() {
			var data int
			err := rivers.FromRange(1, 4).CollectLastAs(&data)

			So(err, ShouldBeNil)
			So(data, ShouldEqual, 4)
		})

		Convey("From Range -> CollectAs -> Sink", func() {
			var numbers []int
			err := rivers.FromRange(1, 4).CollectAs(&numbers)

			So(err, ShouldBeNil)
			So(numbers, ShouldResemble, []int{1, 2, 3, 4})
		})
	})
}
