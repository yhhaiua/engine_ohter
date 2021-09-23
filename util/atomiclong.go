package util

import "sync/atomic"

type AtomicLog struct {
	value int64
}

func (a *AtomicLog)Set(newValue int64)  {
	a.value = newValue
}

func (a *AtomicLog)IncrementAndGet()int64  {
	return  atomic.AddInt64(&a.value,1)
}

func (a *AtomicLog)DecrementAndGet()int64  {
	return  atomic.AddInt64(&a.value,-1)
}

