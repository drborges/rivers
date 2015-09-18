package rivers_test

import (
	"github.com/drborges/rivers"
	"github.com/drborges/rivers/stream"
	"testing"
	"runtime"
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
