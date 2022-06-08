# lock
内部记录了是谁在持有锁，在解锁时会打印使用锁的时间。

Internally records who is holding the lock, and prints the time when the lock is used when unlocking.
## Feature  特性
使用`sync.mutex`封装，带有获取锁超时的额外功能。

Wrapped with `sync.mutex`, with additional functionality for acquiring lock timeouts.

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