package main

import (
	"fmt"
	"sync"
	"time"
)

/*
## Задание: Реализация паттерна "Tee" для записи в несколько реплик БД

**Цель задания**:
Реализовать паттерн "Разветвитель", при котором данные из одного источника параллельно записываются в
несколько реплик базы данных (имитированных каналами).

---

### Описание задачи

Есть сервис, который записывает данные в кластер БД, состоящий из нескольких реплик. Требуется:
1. Принимать данные из входного канала.
2. Параллельно отправлять их во все реплики (каналы).
3. Гарантировать, что данные записаны во все реплики.
4. Корректно закрыть реплики после завершения работы.

```go
package main

import (

	"fmt"
	"time"

)

// Реплика БД (имитация)

	func dbReplica(name string, in <-chan int) {
		for data := range in {
			fmt.Printf("Запись в %s: %d\n", name, data)
			time.Sleep(100 * time.Millisecond) // Имитация задержки записи
		}
		fmt.Printf("Реплика %s закрыта\n", name)
	}

	func main() {
		input := make(chan int) // Канал для входящих данных
		replicas := []chan int{ // Реплики БД (каналы)
			make(chan int),
			make(chan int),
			make(chan int),
		}
	}
*/
func dbReplica(name string, in <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for data := range in {
		fmt.Printf("Запись в %s: %d\n", name, data)
		time.Sleep(100 * time.Millisecond) // Имитация задержки записи
	}
	fmt.Printf("Реплика %s закрыта\n", name)
}

func tee(input <-chan int, replicas []chan int) {
	// Я не думаю что такое решение самое оптимальное, но в голову пока ничего другого не приходит
	for value := range input {
		wg := sync.WaitGroup{}

		wg.Add(len(replicas))

		for _, replica := range replicas {
			go func() {
				defer wg.Done()
				replica <- value
			}()
		}
		wg.Wait()
	}
	for _, replica := range replicas {
		close(replica)
	}
}

func main() {
	input := make(chan int) // Канал для входящих данных
	replicas := []chan int{ // Реплики БД (каналы)
		make(chan int),
		make(chan int),
		make(chan int),
	}

	replicWG := sync.WaitGroup{}
	replicWG.Add(len(replicas))

	for i, replica := range replicas {
		go dbReplica(fmt.Sprintf("Реплика :%d", i+1), replica, &replicWG)
	}

	go tee(input, replicas)

	// проверяем отправку данных - можно вырезать без потерь
	for i := 0; i < 10; i++ {
		input <- i
	}
	close(input)

	replicWG.Wait()
}
