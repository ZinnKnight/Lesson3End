package main

import (
	"fmt"
	"math/rand"
	"sync"
)

/*

Задание: Конвейер чисел на Go

**Цель задания**:
Реализовать конвейер для обработки чисел с использованием горутин и каналов.
Числа из первого канала должны читаться по мере поступления, обрабатываться
(например, возводиться в квадрат) и записываться во второй канал.

---

### Описание задачи

Даны два канала:
- `naturals` (для передачи исходных чисел),
- `squares` (для передачи обработанных чисел).

Необходимо:
1. **Генерировать** числа и отправлять их в канал `naturals`.
2. **Читать** числа из `naturals`, обрабатывать их (возводить в квадрат) и отправлять результат в `squares`.
3. **Выводить** результаты из `squares` в консоль.

```go
package main

func main() {
	naturals := make(chan int)
	squares := make(chan int)
}
```

*/

func main() {
	wg := sync.WaitGroup{}
	naturals := make(chan int)
	squares := make(chan int)

	num := rand.Intn(10)

	wg.Add(num)
	go func() {
		defer wg.Done()
		for n := range num {
			naturals <- n
		}
		close(naturals)
	}()

	go func() {
		defer wg.Done()
		for sqn := range naturals {
			squares <- sqn * sqn
		}
		close(squares)
	}()

	go func() {
		wg.Wait()
	}()

	for res := range squares {
		fmt.Println(res)
	}

}
