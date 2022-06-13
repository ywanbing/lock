//go:build (go1.14 || go1.15 || go1.16 || go1.17) && !go1.18
// +build go1.14 go1.15 go1.16 go1.17
// +build !go1.18

package lock

import (
	"log"
	"os"
	"sync/atomic"
	"time"
)

type mutex struct {
	state int32

	// holder information if you set
	holdInfo string

	// record how long the lock was held
	lockTime time.Time
	log      Logger
}

// NewMutex create mutex with logger
func NewMutex(log Logger) Locker {
	return &mutex{
		log: log,
	}
}

func NewLog() Logger {
	return &Log{
		Logger: log.New(os.Stderr, "", log.LstdFlags),
	}
}

// Lock blocking mode, wait until lock
func (m *mutex) Lock(holdInfo string) {
	lockChan, quitChan := make(chan struct{}, 1), make(chan struct{}, 1)
	go m.loopAcquireLock(holdInfo, lockChan, quitChan)
	<-lockChan
}

// Unlock unlocks m.
// It is a run-time error if m is not locked on entry to Unlock.
func (m *mutex) Unlock() {
	if m.state == lockStatus_UnLock {
		panic("cannot unlock an already unlocked lock")
	}

	go func(holdInfo interface{}, oldTime time.Time) {
		m.log.PrintLockUsageTime("%s the time the lock is held is %d Milliseconds", holdInfo, time.Since(oldTime).Milliseconds())
	}(m.holdInfo, m.lockTime)

	if atomic.CompareAndSwapInt32(&m.state, lockStatus_Lock, lockStatus_UnLock) {
		m.holdInfo = ""
	}
}

// GetHoldInfo This information is only available when the lock is held
func (m mutex) GetHoldInfo() string {
	return m.holdInfo
}

// LockWithTimeOut acquire lock with timeout
func (m *mutex) LockWithTimeOut(holdInfo string, timeOut time.Duration) bool {
	if timeOut <= 0 {
		m.Lock(holdInfo)
		return true
	}

	timer := time.NewTimer(timeOut)
	defer timer.Stop()

	lockChan, quitChan := make(chan struct{}, 1), make(chan struct{}, 1)
	go m.loopAcquireLock(holdInfo, lockChan, quitChan)

	for {
		select {
		case <-timer.C:
			quitChan <- struct{}{}

			// trying to read once
			if _, ok := <-lockChan; ok {
				return true
			}
			return false
		case <-lockChan:
			return true
		}
	}
}

func (m *mutex) loopAcquireLock(holdInfo string, lockChan chan struct{}, quitChan chan struct{}) {
	ticker := time.NewTicker(tryLockTime)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if m.TryLock(holdInfo) {
				lockChan <- struct{}{}
				close(lockChan)
				return
			}
		case <-quitChan:
			close(lockChan)
			return
		}
	}
}

// TryLock tries to lock m and reports whether it succeeded.
func (m *mutex) TryLock(holdInfo string) bool {
	// quick judgment
	if m.state == lockStatus_Lock {
		return false
	}

	if atomic.CompareAndSwapInt32(&m.state, lockStatus_UnLock, lockStatus_Lock) {
		m.holdInfo = holdInfo
		m.lockTime = time.Now()
		return true
	}

	return false
}

// SetLogger set up a log component
func (m *mutex) SetLogger(log Logger) {
	m.log = log
}
