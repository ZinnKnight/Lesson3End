package main

import "sync"

/*

## Оптимизация HTTP-обработчика с sync.Pool

### Описание задачи
В высоконагруженных сервисах, обрабатывающих тысячи HTTP-запросов в секунду,
частая аллокация объектов для декодирования JSON становится узким местом. Каждый вызов `json.NewDecoder`
создает новый экземпляр `RequestData`, что приводит к:
- Высокой нагрузке на GC (сборщик мусора).
- Увеличению времени обработки запросов.
- Нестабильной работе при пиковых нагрузках.

**Цель:**
Использовать `sync.Pool` для переиспользования объектов `RequestData`,
сократив аллокации и улучшив производительность.

---

### Требования
1. **Реализация пула объектов**
    - Создать пул для структур `RequestData` с предварительной инициализацией вложенных полей
(например, `map` или `slice`).
    - Гарантировать потокобезопасность.

2. **Метод `Reset()`**
    - Очистить все поля объекта перед возвратом в пул.
    - Для слайсов: сохранить базовый массив (`items = items[:0]`).
    - Для мап: явно удалить все ключи.

3. **Отсутствие утечек данных**
    - Убедиться, что объекты из пула не сохраняют данные предыдущих запросов.

---
```go
func main() {
    http.HandleFunc("/", handleRequest)
    fmt.Println("Server started at :8080")
    http.ListenAndServe(":8080", nil)
}
```

*/

// я просто собрал самый простой пример куда запихнул слайс и мапу, по факту это можно расширить

type RequestData struct {
	jsonMapData   map[string]string
	jsonSliceData []string
}

func (rd *RequestData) Reset() {
	rd.jsonSliceData = rd.jsonSliceData[:0]
	for key := range rd.jsonMapData {
		delete(rd.jsonMapData, key)
	}
}

var reqPool = sync.Pool{
	New: func() interface{} {
		return &RequestData{
			jsonMapData: make(map[string]string, 10), // я не совсем понимаю какой оринтир или же hint ставить
			// поэтому поставил от балды
			jsonSliceData: make([]string, 10),
		}
	},
}
