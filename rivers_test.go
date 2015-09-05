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

		Convey("From Range -> Filter -> Map -> Reduce -> Each", func() {
			data, _ := rivers.FromRange(1, 5).
				Filter(evensOnly).
				Map(add(1)).
				Reduce(0, sum).
				Collect()

			So(data, ShouldResemble, []stream.T{8})
		})

		Convey("From Data -> Flatten -> Map -> Sort By -> Batch", func() {
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

		Convey("From Slice -> Dispatch If -> Map", func() {
			in, out := stream.New(2)

			notDispatched, _ := rivers.FromSlice([]stream.T{1, 2, 3, 4, 5}).
				DispatchIf(evensOnly, out).
				Map(add(2)).
				Collect()

			data := in.ReadAll()
			So(data, ShouldContain, 2)
			So(data, ShouldContain, 4)

			So(notDispatched, ShouldContain, 3)
			So(notDispatched, ShouldContain, 5)
			So(notDispatched, ShouldContain, 7)
		})

		Convey("Zip -> Map", func() {
			numbers := rivers.FromData(1, 2, 3, 4)
			letters := rivers.FromData("a", "b", "c")

			combined, _ := numbers.Zip(letters.Stream).
				Map(addOrAppend(1, "_")).Collect()

			So(combined, ShouldResemble, []stream.T{2, "a_", 3, "b_", 4, "c_", 5})
		})

		Convey("Zip By -> Filter -> Collect", func() {
			numbers := rivers.FromData(1, 2, 3, 4)
			moreNumbers := rivers.FromData(4, 4, 1)

			combined, err := numbers.ZipBy(sum, moreNumbers.Stream).Filter(evensOnly).Collect()

			So(err, ShouldBeNil)
			So(combined, ShouldResemble, []stream.T{6, 4, 4})
		})

		Convey("Merge -> Map", func() {
			numbers := rivers.FromData(1, 2)
			moreNumbers := rivers.FromData(3, 4)

			combined, _ := numbers.Merge(moreNumbers.Stream).Collect()

			So(len(combined), ShouldEqual, 4)
			So(combined, ShouldContain, 1)
			So(combined, ShouldContain, 2)
			So(combined, ShouldContain, 3)
			So(combined, ShouldContain, 4)
		})

		Convey("From Data -> Drain", func() {
			numbers := rivers.FromData(1, 2, 3, 4)
			numbers.Drain()

			data, opened := <-numbers.Stream
			So(data, ShouldBeNil)
			So(opened, ShouldBeFalse)
		})

		Convey("From Socket -> Map -> Filter -> Map -> Flatten", func() {
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
				Stream

			So(words.ReadAll(), ShouldResemble, []stream.T{"Hello", "there", "rivers!", "super", "cool!"})
		})

		Convey("From Range -> Partition", func() {
			evensStage, oddsStage := rivers.FromRange(1, 4).Partition(evensOnly)
			evens, _ := evensStage.Collect()
			odds, _ := oddsStage.Collect()

			So(evens, ShouldContain, 2)
			So(evens, ShouldContain, 4)

			So(odds, ShouldContain, 1)
			So(odds, ShouldContain, 3)
		})

		Convey("From Range -> Slipt", func() {
			lhs, rhs := rivers.FromRange(1, 2).Split()

			lhsData := lhs.Stream.ReadAll()
			rhsData := rhs.Stream.ReadAll()
			So(lhsData, ShouldContain, 1)
			So(lhsData, ShouldContain, 2)

			So(rhsData, ShouldContain, 1)
			So(rhsData, ShouldContain, 2)
		})

		Convey("From Range -> Slipt N", func() {
			pipelines := rivers.FromRange(1, 2).SplitN(3)

			data0 := pipelines[0].Stream.ReadAll()
			data1 := pipelines[1].Stream.ReadAll()
			data2 := pipelines[2].Stream.ReadAll()

			So(data0, ShouldContain, 1)
			So(data0, ShouldContain, 2)

			So(data1, ShouldContain, 1)
			So(data1, ShouldContain, 2)

			So(data2, ShouldContain, 1)
			So(data2, ShouldContain, 2)
		})

		Convey("From Range -> OnData", func() {
			pipeline := rivers.FromRange(1, 4).OnData(func(data stream.T, emitter stream.Emitter) {
				if data.(int)%2 == 0 {
					emitter.Emit(data)
				}
			})

			So(pipeline.Stream.ReadAll(), ShouldResemble, []stream.T{2, 4})
		})

		Convey("From Range -> TakeFirst N -> Collect", func() {
			taken, _ := rivers.FromRange(1, 4).TakeFirst(2).Collect()

			So(taken, ShouldResemble, []stream.T{1, 2})
		})

		Convey("From Range -> Take", func() {
			pipeline := rivers.FromRange(1, 4).Take(evensOnly)

			So(pipeline.Stream.ReadAll(), ShouldResemble, []stream.T{2, 4})
		})

		Convey("From Range -> Drop", func() {
			pipeline := rivers.FromRange(1, 4).Drop(evensOnly)

			So(pipeline.Stream.ReadAll(), ShouldResemble, []stream.T{1, 3})
		})

		Convey("From Range -> Drop First 2", func() {
			pipeline := rivers.FromRange(1, 5).DropFirst(2)

			So(pipeline.Stream.ReadAll(), ShouldResemble, []stream.T{3, 4, 5})
		})

		Convey("From Range -> Collect", func() {
			data, err := rivers.FromRange(1, 4).Collect()

			So(err, ShouldBeNil)
			So(data, ShouldResemble, []stream.T{1, 2, 3, 4})
		})

		Convey("From Range -> CollectFirst", func() {
			data, err := rivers.FromRange(1, 4).CollectFirst()

			So(err, ShouldBeNil)
			So(data, ShouldEqual, 1)
		})

		Convey("From Range -> CollectFirstAs", func() {
			var data int
			err := rivers.FromRange(1, 4).CollectFirstAs(&data)

			So(err, ShouldBeNil)
			So(data, ShouldEqual, 1)
		})

		Convey("From Range -> CollectLast", func() {
			data, err := rivers.FromRange(1, 4).CollectLast()

			So(err, ShouldBeNil)
			So(data, ShouldEqual, 4)
		})

		Convey("From Range -> CollectLastAs", func() {
			var data int
			err := rivers.FromRange(1, 4).CollectLastAs(&data)

			So(err, ShouldBeNil)
			So(data, ShouldEqual, 4)
		})

		Convey("From Range -> CollectAs", func() {
			var numbers []int
			err := rivers.FromRange(1, 4).CollectAs(&numbers)

			So(err, ShouldBeNil)
			So(numbers, ShouldResemble, []int{1, 2, 3, 4})
		})

		Convey("From Data -> Map From Struct To JSON", func() {
			type Account struct{ Name string }

			items := rivers.FromData(Account{"Diego"}).Map(from.StructToJSON).Stream.ReadAll()

			So(items, ShouldResemble, []stream.T{[]byte(`{"Name":"Diego"}`)})
		})

		Convey("From Data -> Map From JSON To Struct", func() {
			type Account struct{ Name string }

			items := rivers.FromData([]byte(`{"Name":"Diego"}`)).Map(from.JSONToStruct(Account{})).Stream.ReadAll()

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
