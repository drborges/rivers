package rivers_test

import (
	"github.com/drborges/rivers"
	"github.com/drborges/rivers/rx"
	"github.com/drborges/rivers/scanners"
	. "github.com/smartystreets/goconvey/convey"
	"net"
	"strings"
	"testing"
)

func TestRiversAPI(t *testing.T) {
	toString := func(data rx.T) rx.T { return string(data.([]byte)) }
	nonEmptyLines := func(data rx.T) bool { return data.(string) != "" }
	splitWord := func(data rx.T) rx.T { return strings.Split(data.(string), " ") }
	evensOnly := func(data rx.T) bool { return data.(int)%2 == 0 }
	sum := func(a, b rx.T) rx.T { return a.(int) + b.(int) }
	add := func(n int) rx.MapFn {
		return func(data rx.T) rx.T { return data.(int) + n }
	}

	append := func(c string) rx.MapFn {
		return func(data rx.T) rx.T { return data.(string) + c }
	}

	addOrAppend := func(n int, c string) rx.MapFn {
		return func(data rx.T) rx.T {
			if num, ok := data.(int); ok {
				return num + n
			}
			if letter, ok := data.(string); ok {
				return letter + "_"
			}
			return data
		}
	}

	alphabeticOrder := func(a, b rx.T) bool {
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
			data, _ := rivers.New().FromRange(1, 5).
				Filter(evensOnly).
				Map(add(1)).
				Reduce(0, sum).
				Collect()

			So(data, ShouldResemble, []rx.T{8})
		})

		Convey("From Data -> Flatten -> Map -> Sort By -> Batch -> Sink", func() {
			data, _ := rivers.New().FromData([]rx.T{"a", "c"}, "b", []rx.T{"d", "e"}).
				Flatten().
				Map(append("_")).
				SortBy(alphabeticOrder).
				Batch(2).
				Collect()

			So(data, ShouldResemble, []rx.T{
				[]rx.T{"a_", "b_"},
				[]rx.T{"c_", "d_"},
				[]rx.T{"e_"},
			})
		})

		Convey("From Slice -> Dispatch If -> Map -> Sink", func() {
			in, out := rx.NewStream(2)

			notDispatched, _ := rivers.New().FromSlice([]rx.T{1, 2, 3, 4, 5}).
				DispatchIf(evensOnly, out).
				Map(add(2)).
				Collect()

			So(in.Read(), ShouldResemble, []rx.T{2, 4})
			So(notDispatched, ShouldResemble, []rx.T{3, 5, 7})
		})

		Convey("Combine Zipping -> Map -> Sink", func() {
			streams := rivers.New()
			numbers := streams.FromData(1, 2, 3, 4)
			letters := streams.FromData("a", "b", "c")

			combined, _ := streams.CombineZipping(numbers.Sink(), letters.Sink()).
				Map(addOrAppend(1, "_")).Collect()

			So(combined, ShouldResemble, []rx.T{2, "a_", 3, "b_", 4, "c_", 5})
		})

		Convey("Combine Zipping By -> Map -> Sink", func() {
			streams := rivers.New()
			numbers := streams.FromData(1, 2, 3, 4)
			moreNumbers := streams.FromData(4, 4, 1)

			combined, _ := streams.CombineZippingBy(sum, numbers.Sink(), moreNumbers.Sink()).Filter(evensOnly).Collect()

			So(combined, ShouldResemble, []rx.T{6, 4, 4})
		})

		Convey("From Data -> Drain", func() {
			numbers := rivers.New().FromData(1, 2, 3, 4)
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

			words := rivers.New().FromSocketWithScanner("tcp", port, scanners.NewLineScanner()).
				Map(toString).
				Filter(nonEmptyLines).
				Map(splitWord).
				Flatten().
				Sink()

			So(words.Read(), ShouldResemble, []rx.T{"Hello", "there", "rivers!", "super", "cool!"})
		})

		Convey("From Range -> Partition -> Sink", func() {
			evens, odds := rivers.New().FromRange(1, 10).Partition(evensOnly)

			So(evens.Sink().Read(), ShouldResemble, []rx.T{2, 4, 6, 8, 10})
			So(odds.Sink().Read(), ShouldResemble, []rx.T{1, 3, 5, 7, 9})
		})

		Convey("From Range -> Slipt -> Sink", func() {
			lhs, rhs := rivers.New().FromRange(1, 4).Split()

			So(lhs.Sink().Read(), ShouldResemble, []rx.T{1, 2, 3, 4})
			So(rhs.Sink().Read(), ShouldResemble, []rx.T{1, 2, 3, 4})
		})

		Convey("From Range -> Slipt N -> Sink", func() {
			streams := rivers.New().FromRange(1, 4).SplitN(3)

			So(streams[0].Sink().Read(), ShouldResemble, []rx.T{1, 2, 3, 4})
			So(streams[1].Sink().Read(), ShouldResemble, []rx.T{1, 2, 3, 4})
			So(streams[2].Sink().Read(), ShouldResemble, []rx.T{1, 2, 3, 4})
		})

		Convey("From Range -> ProcessWith -> Sink", func() {
			processor := rivers.New().FromRange(1, 4).ProcessWith(func (data rx.T, out rx.Writable) {
				if data.(int) % 2 == 0 {
					out <- data
				}
			})

			So(processor.Sink().Read(), ShouldResemble, []rx.T{2, 4})
		})

		Convey("From Range -> Take -> Sink", func() {
			taken, _ := rivers.New().FromRange(1, 4).Take(2).Collect()

			So(taken, ShouldResemble, []rx.T{1, 2})
		})

		Convey("From Range -> TakeIf -> Sink", func() {
			processor := rivers.New().FromRange(1, 4).TakeIf(evensOnly)

			So(processor.Sink().Read(), ShouldResemble, []rx.T{2, 4})
		})

		Convey("From Range -> DropIf -> Sink", func() {
			processor := rivers.New().FromRange(1, 4).DropIf(evensOnly)

			So(processor.Sink().Read(), ShouldResemble, []rx.T{1, 3})
		})

		Convey("From Range -> Collect -> Sink", func() {
			data, err := rivers.New().FromRange(1, 4).Collect()

			So(err, ShouldBeNil)
			So(data, ShouldResemble, []rx.T{1, 2, 3, 4})
		})

		Convey("From Range -> CollectFirst -> Sink", func() {
			data, err := rivers.New().FromRange(1, 4).CollectFirst()

			So(err, ShouldBeNil)
			So(data, ShouldEqual, 1)
		})

		Convey("From Range -> CollectFirstAs -> Sink", func() {
			var data int
			err := rivers.New().FromRange(1, 4).CollectFirstAs(&data)

			So(err, ShouldBeNil)
			So(data, ShouldEqual, 1)
		})

		Convey("From Range -> CollectLast -> Sink", func() {
			data, err := rivers.New().FromRange(1, 4).CollectLast()

			So(err, ShouldBeNil)
			So(data, ShouldEqual, 4)
		})

		Convey("From Range -> CollectLastAs -> Sink", func() {
			var data int
			err := rivers.New().FromRange(1, 4).CollectLastAs(&data)

			So(err, ShouldBeNil)
			So(data, ShouldEqual, 4)
		})

		Convey("From Range -> CollectAs -> Sink", func() {
			var numbers []int
			err := rivers.New().FromRange(1, 4).CollectAs(&numbers)

			So(err, ShouldBeNil)
			So(numbers, ShouldResemble, []int{1, 2, 3, 4})
		})
	})
}
