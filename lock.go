package lock

import (
	"sync"
	"time"
)

const (
	tryLockTime = 1 * time.Millisecond
)

type mutex struct {
	mux *sync.Mutex

	// holder information if you set
	holdInfo string

	// record how long the lock was held
	lockTime time.Time
	log      Logger
}

// NewDefMutex create default mutex
func NewDefMutex() *mutex {
	return NewMutex(NewLog())
}

// NewMutex create mutex with logger
func NewMutex(log Logger) *mutex {
	return &mutex{
		mux:      new(sync.Mutex),
		holdInfo: "",
		log:      log,
	}
}

// Lock blocking mode, wait until lock
func (m *mutex) Lock(holdInfo string) {
	m.mux.Lock()
	m.holdInfo = holdInfo
	m.lockTime = time.Now()
}

// Unlock unlocks m.
// It is a run-time error if m is not locked on entry to Unlock.
func (m *mutex) Unlock() {
	go func(holdInfo interface{}, oldTime time.Time) {
		m.log.PrintLockUsageTime("%s the time the lock is held is %d Milliseconds", holdInfo, time.Since(oldTime).Milliseconds())
	}(m.holdInfo, m.lockTime)

	m.holdInfo = ""
	m.mux.Unlock()
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
	if !m.mux.TryLock() {
		return false
	}

	m.holdInfo = holdInfo
	m.lockTime = time.Now()
	return true
}

// SetLogger set up a log component
func (m *mutex) SetLogger(log Logger) {
	m.log = log
}
