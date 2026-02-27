package main

import (
	"fmt"
	"sync"
)

/*

## Задание: Объединение каналов в Go

**Цель задания**:
Написать функцию `mergeChannels`, которая объединяет данные из нескольких каналов в один общий канал,
используя паттерн `FAN-IN`.

---

### Описание задачи

Дано:
- `n` каналов типа `<-chan int`.
- Функция должна вернуть канал `<-chan int`, в который попадают все значения из исходных каналов.

Требования:
1. Все значения из входных каналов должны быть отправлены в выходной канал.
2. Выходной канал должен быть закрыт после завершения всех входных каналов.
3. Решение должно быть потокобезопасным и эффективным.

*/

func mergeChannels(n ...<-chan int) <-chan int {
	wg := sync.WaitGroup{}
	out := make(chan int)

	go func() {
		for chanNum := range n {
			wg.Add(1)
			go func() {
				defer wg.Done()

				for val := range chanNum {
					out <- val
				}
			}()
		}
		wg.Wait()
		close(out)
	}()
	return out
}

/*
Я если честно не увидел ниже начало функции, если брать представленный формат то по факту код поменяется не сильно
уйдёт "основная" горутинка - на её место встанет main и каналы нужно будет прежде собрать в источник данных, ну или же по одному
проходить и данные кидать

Ниже написал вариант с переделанным форматом

package main

func mergeChannels(channels ...<-chan int) <-chan int {

}

func main() {
	a := make(chan int)
	b := make(chan int)
	c := make(chan int)
}

*/

func mergeChannels(n ...<-chan int) <-chan int {
	wg := sync.WaitGroup{}
	wg.Add(len(n))

	out := make(chan int)

	for _, ch := range n {
		ch := ch

		go func(c <-chan int) {
			defer wg.Done()

			for val := range c {
				out <- val
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func main() {
	a := make(chan int)
	b := make(chan int)
	c := make(chan int)

	go func() {
		defer func() {
			close(a)
			close(b)
			close(c)
		}()

		// по факту тут просто наполнение данными для наглядности, но это можно спокойно скипнуть
		for i := 0; i < 10; i++ {
			a <- i
			b <- i + i
			c <- i * i
		}
	}()
	for val := range mergeChannels(a, b, c) {
		fmt.Println(val)
	}
}
