package main

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
)

/*

## Задание: Реализация декоратора для преобразования метрик в реальном времени

**Цель задания**:
Создать гибкий декоратор для каналов, который будет автоматически преобразовывать метрики
серверов из байтов в мегабайты перед отправкой в API. Используя паттерн `TRANSFORMER`
---

### Описание задачи

В системе мониторинга серверов:
1. **Источник данных**: Канал `metrics <-chan ServerMetric` получает метрики в формате:
   ```go
   type ServerMetric struct {
       Name  string  // Название метрики (например, "memory_usage")
       Value float64 // Значение в байтах
   }
   ```

# Задание: Реализация конвейерной обработки данных (Pipeline паттерн)

**Цель задания**:
Создать конвейер из трех этапов для обработки строковых данных:
1. **Парсинг** — добавление метки "parsed" к данным.
2. **Разделение** — распределение данных между N каналами (round-robin).
3. **Отправка** — параллельная обработка данных в N горутинах с добавлением метки "sent".

---

## Описание задачи

Ваша задача — реализовать систему, которая:
- Обрабатывает данные в строгом порядке: **Parse → Split → Send**.
- Корректно закрывает все каналы после завершения работы.
- Гарантирует потокобезопасность и отсутствие утечек горутин.

### Этапы конвейера

1. **Parse**:
   - Принимает канал сырых данных (`<-chan string`).
   - Добавляет к каждой строке префикс "parsed - ".
   - Возвращает канал обработанных данных.

2. **Split**:
   - Принимает канал данных и число `N` (количество выходных каналов).
   - Распределяет данные между `N` каналами в порядке round-robin.
   - Возвращает слайс каналов (`[]<-chan string`).

3. **Send**:
   - Принимает слайс каналов и запускает `N` горутин.
   - Каждая горутина добавляет к данным префикс "sent - ".
   - Возвращает объединенный канал результатов.

Но если что мы с тобой и так пройдем эти темы. А если хочешь прям догнать,то вот ресурсы
и пиши по вопросам. Можем отдельно встречу организовать по вопросам:

https://www.youtube.com/watch?v=luQlkud-jKE&t=5s

https://habr.com/ru/companies/pt/articles/764850/
*/

type ServerMetric struct {
	Name  string  // Название метрики (например, "memory_usage")
	Value float64 // Значение в байтах
}

func parseDecorator(metrics <-chan ServerMetric) chan ServerMetric {
	reworkedData := make(chan ServerMetric)
	go func() {
		defer close(reworkedData)
		for m := range metrics {
			m.Name = "parsed - " + m.Name
			reworkedData <- m
		}
	}()
	return reworkedData
}

func splitDecorator(reworkedData <-chan ServerMetric, n int) []chan ServerMetric {
	if n <= 0 {
		return []chan ServerMetric{}
	}
	splitChans := make([]chan ServerMetric, n)
	datafillChans := make([]chan ServerMetric, n)

	for i := 0; i < n; i++ {
		splitChans[i] = make(chan ServerMetric)
		datafillChans[i] = make(chan ServerMetric)
	}

	go func() {
		defer func() {
			for _, ch := range splitChans {
				ch := ch
				close(ch)
			}
		}()
		i := 0
		for m := range reworkedData {
			splitChans[i] <- m
			i++
			if i == n {
				i = 0
			}
		}
	}()
	return splitChans
}

func sendDecorator(splitChans []chan ServerMetric) chan ServerMetric {
	out := make(chan ServerMetric)
	wg := sync.WaitGroup{}
	wg.Add(len(splitChans))
	for _, splitChan := range splitChans {
		go func() {
			defer wg.Done()
			for m := range splitChan {
				m.Name = strings.Replace(m.Name, "parsed - ", "sent - ", 1)
				out <- m
			}
		}()
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func main() {
	metrics := make(chan ServerMetric)

	go func() {
		defer close(metrics)
		for i := 0; i < 10; i++ {
			metrics <- ServerMetric{
				Name:  "memory_usage",
				Value: float64(rand.Intn(1000)),
			}
		}
	}()

	parse := parseDecorator(metrics)

	split := splitDecorator(parse, 5)

	out := sendDecorator(split)

	for m := range out {
		fmt.Println(m.Name, m.Value)
	}

}
