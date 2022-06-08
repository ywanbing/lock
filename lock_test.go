package lock

import (
	"fmt"
	"log"
	"sync"
	"testing"
	"time"
)

type Log struct {
	*log.Logger
}

func (l *Log) Infof(format string, args ...interface{}) {
	l.Printf(format, args...)
}

func NewLog() Logger {
	return &Log{
		Logger: log.Default(),
	}
}

func TestMutex_Lock(t *testing.T) {
	mu := NewMutex(NewLog())

	wg := sync.WaitGroup{}
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			mu.Lock(fmt.Sprintf("func.Lock%d", idx))
			defer mu.Unlock()
			fmt.Println("lock success ------>  ", idx)
			time.Sleep(time.Duration(idx) * time.Second)
		}(i)
	}
	wg.Wait()
}

func TestMutex_LockWithTimeOut(t *testing.T) {
	mu := NewMutex(NewLog())

	wg := &sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			if !mu.LockWithTimeOut(fmt.Sprintf("func.Lock%d", idx), 10*time.Millisecond) {
				fmt.Printf("lock filed idx:%d ------>  holdInfo:%v \n", idx, mu.GetHoldInfo())
				return
			}

			defer func() {
				mu.Unlock()
				fmt.Println("unlock success ------>  ", idx)
			}()

			fmt.Println("lock success ------>  ", idx)
			time.Sleep(2 * time.Millisecond)
		}(i)
	}
	wg.Wait()
}

func TestMutex_TryLock(t *testing.T) {
	mu := NewMutex(NewLog())
	wg := sync.WaitGroup{}
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			if !mu.TryLock(fmt.Sprintf("func.Lock%d", idx)) {
				fmt.Printf("lock filed idx:%d ------>  holdInfo:%v \n", idx, mu.GetHoldInfo())
				return
			}
			defer func() {
				mu.Unlock()
				fmt.Println("unlock success ------>  ", idx)
			}()
			fmt.Println("lock success ------>  ", idx)
		}(i)
	}
	wg.Wait()
}
