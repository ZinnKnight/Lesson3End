package main

import (
	"fmt"
	"math/rand"
	"sync"
)

/*

# CHANNELS
## Задание: Анализ и исправление кода с гонками данных

### Описание задачи
1. Внимательно изучить код.
2. Найти все ошибки, описать их в комментариях прямо в коде.
3. Исправить код, обеспечив корректную работу.

```golang
package main

import (
	"fmt"
	"math/rand"
	"sync"
)

func main() {
	alreadyStored := make(map[int]struct{})
	capacity := 1000
	doubles := make([]int, 0, capacity)
	for i := 0; i < capacity; i++ { // нет переопределения типа i := i, я полагаю это может вызвать проблему
		doubles = append(doubles, rand.Intn(10))
	}
	uniqueIDs := make(chan int, capacity) // буфер на 1к элементов не выглядит привлекательно
	wg := sync.WaitGroup{}
	for i := 0; i < capacity; i++ {
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, ok := alreadyStored[doubles[i]]; !ok { // не увидел мьютекса, ну или синхромапы - мы просто поломаем
// данные при чтнеии и попытке записи
				alreadyStored[doubles[i]] = struct{}{}
				uniqueIDs <- doubles[i]
			}
		}()
	}
	wg.Wait()
	for val := range uniqueIDs { // из канала читаем, но сам канал при этом не закрывается нигде
// по идее можно и wgwait и закрытие канала вынести в отдельную горутинку
		fmt.Println(val)
	}
	fmt.Println(uniqueIDs) // не совсем понял зачем нам печатать сам канал, если мы уже напечатали содержимое из него
}
```
*/

func main() {
	alreadyStored := make(map[int]struct{})
	capacity := 1000
	doubles := make([]int, 0, capacity)
	for i := 0; i < capacity; i++ {
		doubles = append(doubles, rand.Intn(10))
	}
	uniqueIDs := make(chan int, 10) // если правильно понял то псевдорандомизированных значений у нас всего 10
	// (выводится тож 10 при запуске), значит буфер излишен - подтёр
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	for i := 0; i < capacity; i++ {
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			if _, ok := alreadyStored[doubles[i]]; !ok {
				alreadyStored[doubles[i]] = struct{}{}
				uniqueIDs <- doubles[i]
			}
			mu.Unlock()
		}()
	}

	go func() {
		wg.Wait()
		close(uniqueIDs)
	}()

	for val := range uniqueIDs {
		fmt.Println(val)
	}
	fmt.Println(uniqueIDs)
}

// рандом вида rand.Intn - в целом +- достаточно рандомен для такой штуки, но как буд-то лучше через seed использовать,
// хотя Goшка говорит что лучше так не делать
// я не уверен в альтернативе - буду рад пояснениям
