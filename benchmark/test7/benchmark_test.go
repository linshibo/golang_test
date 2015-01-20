package test7

import (
	"testing"
	"sync"
	"fmt"
)


func Benchmark_InvokeFunc(b *testing.B) {
    for i := 0; i < b.N; i++ {
		func()int{
			return 0
		}()
    }
}

func Benchmark_InvokeGoroutine(b *testing.B) {
    for i := 0; i < b.N; i++ {
		go func()int{
			return 0
		}()
    }
}

func Benchmark_Mutex(b *testing.B) {
	var lock sync.Mutex
    for i := 0; i < b.N; i++ {
		lock.Lock()
		lock.Unlock()
    }
}


func Benchmark_RWMutexWrite(b *testing.B) {
	var lock sync.RWMutex
    for i := 0; i < b.N; i++ {
		lock.Lock()
		lock.Unlock()
    }
}
func Benchmark_RWMutexRead(b *testing.B) {
	var lock sync.RWMutex
    for i := 0; i < b.N; i++ {
		lock.RLock()
		lock.RUnlock()
    }
}

func Benchmark_MutexFunc(b *testing.B) {
	var lock sync.Mutex
    for i := 0; i < b.N; i++ {
		func(){
			lock.Lock()
			lock.Unlock()
		}()
    }
}

func Benchmark_DeferMutexFunc(b *testing.B) {
	var lock sync.Mutex
    for i := 0; i < b.N; i++ {
		func(){
			lock.Lock()
			defer lock.Unlock()
		}()
    }
}
func Benchmark_Sprintf(b *testing.B) {
	s := "hello"
	c := ""
    for i := 0; i < b.N; i++ {
		c=fmt.Sprintf("%s,%s",s,s)
    }
	s=c
}
func Benchmark_stringplus(b *testing.B) {
	s := "hello"
	c := "" 
    for i := 0; i < b.N; i++ {
		c=s+s
    }
	s=c
}
