package lock

import (
	"sync"
	"sync/atomic"
	"time"
)

const (
	mutexUnLocked int32 = iota
	mutexLocked

	tryLockTime = 1 * time.Millisecond
)

type Logger interface {
	Infof(format string, v ...interface{})
}

type mutex struct {
	mux *sync.Mutex

	// lock status
	state int32

	// holder information if you set
	holdInfo interface{}

	// record how long the lock was held
	lockTime time.Time
	log      Logger
}

func NewMutex(log Logger) *mutex {
	return &mutex{
		mux:      new(sync.Mutex),
		state:    0,
		holdInfo: nil,
		log:      log,
	}
}

// Lock blocking mode, wait until lock
func (m *mutex) Lock(holdInfo interface{}) {
	m.mux.Lock()
	if atomic.CompareAndSwapInt32(&m.state, mutexUnLocked, mutexLocked) {
		m.holdInfo = holdInfo
		m.lockTime = time.Now()
		return
	}

	panic("lock state err: the original state should be mutexUnLocked")
}

// Unlock unlocks m.
// It is a run-time error if m is not locked on entry to Unlock.
func (m *mutex) Unlock() {
	if atomic.CompareAndSwapInt32(&m.state, mutexLocked, mutexUnLocked) {
		go func(holdInfo interface{}, oldTime time.Time) {
			m.log.Infof("%s the time the lock is held is %v Milliseconds", holdInfo, time.Since(oldTime).Milliseconds())
		}(m.holdInfo, m.lockTime)
		
		m.holdInfo = nil
		m.mux.Unlock()
		return
	}

	panic("lock state err: the original state should be mutexLocked")
}

// GetHoldInfo This information is only available when the lock is held
func (m mutex) GetHoldInfo() interface{} {
	return m.holdInfo
}

// LockWithTimeOut acquire lock with timeout
func (m *mutex) LockWithTimeOut(holdInfo interface{}, timeOut time.Duration) bool {
	if timeOut <= 0 {
		m.Lock(holdInfo)
		return true
	}

	ticker := time.NewTimer(timeOut)
	defer ticker.Stop()

	lockChan, quitChan := make(chan struct{}, 1), make(chan struct{}, 1)
	go m.loopAcquireLock(holdInfo, lockChan, quitChan)

	for {
		select {
		case <-ticker.C:
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

func (m *mutex) loopAcquireLock(holdInfo interface{}, lockChan chan struct{}, quitChan chan struct{}) {
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
func (m *mutex) TryLock(holdInfo interface{}) bool {
	// fail fast
	if m.state == mutexLocked {
		return false
	}

	if !m.mux.TryLock() {
		return false
	}

	if atomic.CompareAndSwapInt32(&m.state, mutexUnLocked, mutexLocked) {
		m.holdInfo = holdInfo
		m.lockTime = time.Now()
		return true
	}

	panic("lock state err: the original state should be mutexUnLocked")
}
