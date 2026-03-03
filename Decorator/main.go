package main

import (
	"fmt"
	"strconv"
	"sync"
)

/*

## Задание: Реализация декоратора для преобразования метрик в реальном времени

**Цель задания**:
Создать гибкий декоратор для каналов, который будет автоматически преобразовывать метрики серверов из байтов в
мегабайты перед отправкой в API. Используя паттерн `TRANSFORMER`
---

### Описание задачи

В системе мониторинга серверов:
1. **Источник данных**: Канал `metrics <-chan ServerMetric` получает метрики в формате:
   ```go
   type ServerMetric struct {
       Name  string  // Название метрики (например, "memory_usage")
       Value float64 // Значение в байтах
   }

*/

type ServerMetric struct {
	Name  string  // Название метрики (например, "memory_usage")
	Value float64 // Значение в байтах
}

func decorator(metrics chan ServerMetric) chan ServerMetric {
	outChan := make(chan ServerMetric)
	// я не уверен что так правильно, но работает и работу делает
	// хотя по идее условные побитовые сдвиги или операции с байтами напрямую были бы более удобные наверное
	const bytesInMegabites = 1024 * 1024

	go func() {
		defer close(outChan)
		for mx := range metrics {
			mx.Value = mx.Value / float64(bytesInMegabites)
			outChan <- mx
		}
	}()
	return outChan
}

func main() {
	metrics := make(chan ServerMetric)
	
	wg := sync.WaitGroup{}

	for i := 0; i < 10; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()
			metrics <- ServerMetric{
				Name:  "server" + strconv.Itoa(i),
				Value: 2.0 * float64(i),
			}
		}(i)
	}

	go func() {
		wg.Wait()
		close(metrics)
	}()

	for m := range decorator(metrics) {
		fmt.Println(m)
	}
}
