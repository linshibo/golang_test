package test2

import(
 "testing"
 "math/rand"
)

type BigStruct struct {
	C01 int
	C02 int
	C03 int
	C04 int
	C05 int
	C06 int
	C07 int
	C08 int
	C09 int
	C10 int
	C11 int
	C12 int
	C13 int
	C14 int
	C15 int
	C16 int
	C17 int
	C18 int
	C19 int
	C20 int
	C21 int
	C22 int
	C23 int
	C24 int
	C25 int
	C26 int
	C27 int
	C28 int
	C29 int
	C30 int
}
func Loop1(a []BigStruct, key int) int {
	for i := 0; i < len(a); i++ {
		if a[i].C30 == key {
			return i
		}
	}

	return -1
}

func Loop2(a map[int]BigStruct, key int) int {
	return  a[key].C30
}

func Loop3(a []BigStruct, key int) int {
	for i, _ := range a {
		if a[i].C30 == key {
			return a[i].C30 
		}
	}
	return -1
}

func Loop4(a []BigStruct, key int) int {
	for _, x := range a {
		if x.C30 == key {
			return x.C30 
		}
	}
	return -1
}

func Loop5(a map[int]*BigStruct, key int) int {
	return  a[key].C30
}


const(
	size = 100
)

func Benchmark_Loop1(b *testing.B) {
	var a = make([]BigStruct, size)
	for i := range a{
		a[i].C30=rand.Intn(1000000)
	}
	for i := 0; i < b.N; i++ {
		index := rand.Intn(size)
		Loop1(a, a[index].C30)
	}
}

func Benchmark_Loop2(b *testing.B) {
	var a = make(map[int]BigStruct, size)
	for i:=0;i<size;i++{
		a[i]=BigStruct{C30: 1}
	}
	for i := 0; i < b.N; i++ {
		index := rand.Intn(size)
		Loop2(a, index)
	}
}

func Benchmark_Loop3(b *testing.B) {
	var a = make([]BigStruct, size)

	for i := range a{
		a[i].C30=rand.Intn(1000000)
	}
	for i := 0; i < b.N; i++ {
		index := rand.Intn(size)
		Loop3(a, a[index].C30)
	}
}
func Benchmark_Loop4(b *testing.B) {
	var a = make([]BigStruct, size)

	for i := range a{
		a[i].C30=rand.Intn(1000000)
	}
	for i := 0; i < b.N; i++ {
		index := rand.Intn(size)
		Loop4(a, a[index].C30)
	}
}

func Benchmark_Loop5(b *testing.B) {
	var a = make(map[int]*BigStruct, size)
	for i:=0;i<size;i++{
		a[i]=&BigStruct{C30: 1}
	}
	for i := 0; i < b.N; i++ {
		index := rand.Intn(size)
		Loop5(a, index)
	}
}


