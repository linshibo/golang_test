package test6

import "testing"

type BigStruct struct {
	c [10]int
}

func Invoke1(a *BigStruct) uint64 {
    return 0
}

func Invoke2(a BigStruct) uint64 {
    return 0
}

func Invoke3(a interface{}) uint64 {
    return 0
}

func Benchmark_Invoke1(b *testing.B) {
    var a = new(BigStruct)

    for i := 0; i < b.N; i++ {
        Invoke1(a)
    }
}

func Benchmark_Invoke2(b *testing.B) {
    var a = BigStruct{}

    for i := 0; i < b.N; i++ {
        Invoke2(a)
    }
}

func Benchmark_Invoke3(b *testing.B) {
    var a = BigStruct{}

    for i := 0; i < b.N; i++ {
        Invoke3(a)
    }
}

