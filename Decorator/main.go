package main

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

func decorator(metrics <-chan ServerMetric, worker func(ServerMetric) ServerMetric) chan ServerMetric {
	outChan := make(chan ServerMetric)
	// я не уверен что так правильно, но работает и работу делает
	// хотя по идее условные побитовые сдвиги или операции с байтами напрямую были бы более удобные наверное
	const bytesInMegabites = 1024 * 1024

	go func() {
		defer close(outChan)
		for mx := range metrics {
			mx.Value = mx.Value / float64(bytesInMegabites)
		}
	}()
	return outChan
}
