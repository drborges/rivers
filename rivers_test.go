package rivers_test

import (
	"github.com/drborges/rivers"
	"github.com/drborges/rivers/producers"
	"github.com/drborges/rivers/scanners"
	"github.com/drborges/rivers/stream"
	"github.com/drborges/rivers/transformers/from"
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

	concat := func(c string) stream.MapFn {
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
				Map(concat("_")).
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

			combined, _ := numbers.Zip(letters.Sink()).
				Map(addOrAppend(1, "_")).Collect()

			So(combined, ShouldResemble, []stream.T{2, "a_", 3, "b_", 4, "c_", 5})
		})

		Convey("Zip By -> Filter -> Collect", func() {
			numbers := rivers.FromData(1, 2, 3, 4)
			moreNumbers := rivers.FromData(4, 4, 1)

			combined, err := numbers.ZipBy(sum, moreNumbers.Sink()).Filter(evensOnly).Collect()

			So(err, ShouldBeNil)
			So(combined, ShouldResemble, []stream.T{6, 4, 4})
		})

		Convey("Merge -> Map -> Sink", func() {
			numbers := rivers.FromData(1, 2, 3, 4)
			moreNumbers := rivers.FromData(4, 4, 1)

			combined, _ := numbers.Merge(moreNumbers.Sink()).Collect()

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

			words := rivers.From(producers.FromSocket("tcp", port, scanners.NewLineScanner())).
				Map(toString).
				Filter(nonEmptyLines).
				Map(splitWord).
				Flatten().
				Sink()

			So(words.Read(), ShouldResemble, []stream.T{"Hello", "there", "rivers!", "super", "cool!"})
		})

		Convey("From Range -> Partition -> Sink", func() {
			evensStream, oddsStream := rivers.FromRange(1, 4).Partition(evensOnly)
			evens, _ := evensStream.Collect()
			odds, _ := oddsStream.Collect()

			So(evens, ShouldContain, 2)
			So(evens, ShouldContain, 4)
			So(odds, ShouldContain, 1)
			So(odds, ShouldContain, 3)
		})

		Convey("From Range -> Slipt -> Sink", func() {
			lhs, rhs := rivers.FromRange(1, 4).Split()

			lhsData := lhs.Sink().Read()
			rhsData := rhs.Sink().Read()
			So(lhsData, ShouldContain, 1)
			So(lhsData, ShouldContain, 2)
			So(lhsData, ShouldContain, 3)
			So(lhsData, ShouldContain, 4)

			So(rhsData, ShouldContain, 1)
			So(rhsData, ShouldContain, 2)
			So(rhsData, ShouldContain, 3)
			So(rhsData, ShouldContain, 4)
		})

		Convey("From Range -> Slipt N -> Sink", func() {
			streams := rivers.FromRange(1, 2).SplitN(3)

			data0 := streams[0].Sink().Read()
			data1 := streams[1].Sink().Read()
			data2 := streams[2].Sink().Read()

			So(data0, ShouldContain, 1)
			So(data0, ShouldContain, 2)

			So(data1, ShouldContain, 1)
			So(data1, ShouldContain, 2)

			So(data2, ShouldContain, 1)
			So(data2, ShouldContain, 2)
		})

		Convey("From Range -> OnData -> Sink", func() {
			pipe := rivers.FromRange(1, 4).OnData(func(data stream.T, emitter stream.Emitter) {
				if data.(int)%2 == 0 {
					emitter.Emit(data)
				}
			})

			So(pipe.Sink().Read(), ShouldResemble, []stream.T{2, 4})
		})

		Convey("From Range -> TakeFirst N -> Sink", func() {
			taken, _ := rivers.FromRange(1, 4).TakeFirst(2).Collect()

			So(taken, ShouldResemble, []stream.T{1, 2})
		})

		Convey("From Range -> Take -> Sink", func() {
			pipe := rivers.FromRange(1, 4).Take(evensOnly)

			So(pipe.Sink().Read(), ShouldResemble, []stream.T{2, 4})
		})

		Convey("From Range -> Drop -> Sink", func() {
			pipe := rivers.FromRange(1, 4).Drop(evensOnly)

			So(pipe.Sink().Read(), ShouldResemble, []stream.T{1, 3})
		})

		Convey("From Range -> Drop First 2 -> Sink", func() {
			pipe := rivers.FromRange(1, 5).DropFirst(2)

			So(pipe.Sink().Read(), ShouldResemble, []stream.T{3, 4, 5})
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

		Convey("From Data -> Map From Struct To JSON -> Sink", func() {
			type Account struct{ Name string }

			items := rivers.FromData(Account{"Diego"}).Map(from.StructToJSON).Sink().Read()

			So(items, ShouldResemble, []stream.T{[]byte(`{"Name":"Diego"}`)})
		})

		Convey("From Data -> Map From JSON To Struct -> Sink", func() {
			type Account struct{ Name string }

			items := rivers.FromData([]byte(`{"Name":"Diego"}`)).Map(from.JSONToStruct(Account{})).Sink().Read()

			So(items, ShouldResemble, []stream.T{Account{"Diego"}})
		})

		Convey("From Range -> Collect By", func() {
			items := []stream.T{}
			err := rivers.FromRange(1, 5).CollectBy(func(data stream.T) {
				items = append(items, data)
			})

			So(err, ShouldBeNil)
			So(items, ShouldResemble, []stream.T{1, 2, 3, 4, 5})
		})
	})
}
