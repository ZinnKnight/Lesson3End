package main

import (
	"fmt"
	"sync"
)

/*

# SYNC/WAIT

# Задание: Синхронизация горутин с использованием `sync.WaitGroup`

**Цель задания**:
Исправить код, чтобы все горутины корректно выводили значения от 0 до 99,
и обеспечить завершение всех горутин перед выходом из программы.
```go
package main

import "fmt"

func main() {
	cnt := 100
	for i := 0; i < cnt; i++ {
		go func() {
			fmt.Println(i)
		}()
	}
}
```

*/

func main() {
	wg := sync.WaitGroup{}
	cnt := 100
	for i := 0; i < cnt; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println(i)
		}()
	}
	wg.Wait()
}
