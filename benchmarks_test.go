package rivers_test

import (
	"errors"
	"github.com/drborges/rivers"
	"github.com/drborges/rivers/stream"
	. "github.com/smartystreets/goconvey/convey"
	"math/rand"
	"runtime"
	"testing"
	"time"
)

func slowProcess(data stream.T) {
	time.Sleep(time.Second)
}

func BenchmarkParallel(b *testing.B) {
	b.ResetTimer()

	runtime.GOMAXPROCS(runtime.NumCPU())
	numbers := rivers.FromRange(1, 10000).Parallel()

	for i := 0; i < b.N; i++ {
		numbers = numbers.Each(slowProcess)
	}

	numbers.Drain()
}

func TestErrorInjection(t *testing.T) {
	err := errors.New("Error Injected")

	injectError := func(data stream.T) stream.T {
		panicFactor := rand.Intn(10)
		if panicFactor%2 == 0 {
			time.Sleep(500 * time.Millisecond)
			panic(err)
		} else {
			time.Sleep(100 * time.Millisecond)
		}
		return data
	}

	Convey("From Range -> Group By", t, func() {
		runtime.GOMAXPROCS(runtime.NumCPU())
		items, err := rivers.FromRange(1, 1000).Parallel().Map(injectError).Collect()

		So(items, ShouldNotBeEmpty)
		So(err, ShouldEqual, err)
	})
}
