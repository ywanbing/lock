# lock
内部记录了是谁在持有锁，在解锁时会打印使用锁的时间。

Internally records who is holding the lock, and prints the time when the lock is used when unlocking.
## Feature  特性
go.1.18 使用`sync.mutex`封装，带有获取锁超时的额外功能。

go.1.18 Wrapped with `sync.mutex`, with additional functionality for acquiring lock timeouts.

小于go1.18 的版本采用原子操作进行封装，对于锁的使用应该更快，当前是在通常情况下。

Versions smaller than go1.18 are encapsulated with atomic operations, and the use of locks should be faster, which is currently the usual case.

## HowToUse 使用方式

example: use mutex 
```go
package main

import "github.com/ywanbing/lock"

func main() {
    mu := lock.NewDefMutex()
    mu.Lock("lock1")
    defer mu.Unlock()
    
    // do something ...
}
```

example: use mutex with timeOut
```go
package main

import (
	"time"

	"github.com/ywanbing/lock"
)

func main() {
    mu := lock.NewDefMutex()
    
    if !mu.LockWithTimeOut("lock1", 1*time.Second) {
        return
    }
    	
    defer mu.Unlock()
    // do something ...
}
```

example: use self log component
```go
package main

import (
	"github.com/ywanbing/lock"
)

func main() {
    // Implement lock.Logger to replace it with your own log
    log := lock.NewLog()
    mu := lock.NewMutex(log)
    
    mu.Lock("lock1")
    defer mu.Unlock()
    // do something  	
}
```