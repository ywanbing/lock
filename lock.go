package lock

import "time"

const (
	lockStatus_UnLock = 0
	lockStatus_Lock   = 1

	tryLockTime = 1 * time.Millisecond
)

type Locker interface {
	Lock(holdInfo string)
	Unlock()
	GetHoldInfo() string
	LockWithTimeOut(holdInfo string, timeOut time.Duration) bool
	TryLock(holdInfo string) bool
	SetLogger(log Logger)
}

// NewDefMutex create default mutexGo18
func NewDefMutex() Locker {
	return NewMutex(NewLog())
}
