package main

import (
	"fmt"
	"sync"
)

/*

# CONCURRENCY
## Реализация потокобезопасного кеша

### Описание задачи
Ваша задача — реализовать потокобезопасный кеш для хранения данных в формате ключ-значение.
Кеш должен безопасно обрабатывать одновременные операции записи и чтения из множества горутин.

### Требования
1. Реализовать структуру `SafeCache` с методами:
   - `Set(key string, value string)` — добавляет значение в кеш.
   - `Get(key string) (string, bool)` — возвращает значение по ключу.
2. Гарантировать отсутствие data race при параллельном доступе.

*/

type SafeCache struct {
	casheMap map[string]string
	mutex    sync.RWMutex
	key      string
	value    string
}

func (sc *SafeCache) Set(key string, value string) {
	if key == "" || value == "" {
		return
	}
	sc.mutex.Lock()
	sc.casheMap[key] = value
	sc.mutex.Unlock()

}

func (sc *SafeCache) Get(key string) (string, bool) {
	if key == "" {
		return "", false
	}
	sc.mutex.RLock()
	value, ok := sc.casheMap[key]
	sc.mutex.RUnlock()
	return value, ok
}

func main() {
	test := &SafeCache{
		casheMap: make(map[string]string),
		key:      "",
		value:    "",
	}

	test.Set("test", "Go")
	fmt.Println(test.Get("test2"))

}
