package main

import (
	"bytes"
	"encoding/json"
	"sync"
	"time"
)

/*
## JSON-кэш с `sync.Pool` и `map + RWMutex`

### **Описание**
Этот проект демонстрирует **потокобезопасный JSON-кэш** с поддержкой TTL и оптимизированной сериализацией.
Используется `sync.Pool` для **эффективной работы с JSON**, а также `map + RWMutex`
для **более быстрого доступа к данным**.

---

### **Основные возможности**
- **Хранение объектов в `map` (с TTL)**
- **Автоматическое удаление устаревших объектов**
- **Быстрая сериализация JSON с `sync.Pool`**
- **Использование `sync.RWMutex` для конкурентного доступа**

---

### **Методы**
#### **Базовые операции**
- `Set(key string, value interface{})` – **добавить объект в кэш**
- `Get(key string) (interface{}, bool)` – **получить объект по ключу**
- `Delete(key string)` – **удалить объект**
- `ToJSON() ([]byte, error)` – **сериализовать кэш в JSON**

---

### **Как это работает?**
- Все объекты хранятся в **`map[string]item`** (ключ → объект с TTL).
- `sync.Pool` позволяет **переиспользовать JSON-буферы**, снижая нагрузку на GC.
- Очистка устаревших данных выполняется **в отдельной горутине**.

```go

	func main() {
	    cache := NewObjectCache(5 * time.Second)

	    // Добавляем данные в кэш
	    cache.Set("user:1", map[string]string{"name": "Alice", "role": "admin"})
	    cache.Set("user:2", map[string]string{"name": "Bob", "role": "user"})

	    // Получаем объект
	    if user, found := cache.Get("user:1"); found {
	        fmt.Println("Найден:", user)
	    }

	    // Выводим JSON
	    jsonData, _ := cache.ToJSON()
	    fmt.Println("Кэш в JSON:", string(jsonData))

	    // Ждём истечения TTL и проверяем снова
	    time.Sleep(6 * time.Second)
	    _, found := cache.Get("user:1")
	    fmt.Println("После TTL, user:1 найден?", found)
	}

```
Но если что мы с тобой и так пройдем эти темы. А если хочешь прям догнать,то вот дополнительные ресурсы. Можем отдельно встречу организовать по вопросам::
https://ubiklab.net/posts/go-pool-and-mechanics-behind-it/
https://reliasoftware.com/blog/golang-sync-pool
https://dev.to/func25/go-syncpool-and-the-mechanics-behind-it-52c1
https://engineer.yadro.com/article/three-ways-to-optimize-memory-performance-on-go-with-memory-pools/
https://leapcell.io/blog/boost-go-performance-sync-pool
https://www.sobyte.net/post/2022-06/go-sync-pool/
https://goperf.dev/01-common-patterns/object-pooling/
*/
type dataPice struct {
	data interface{}
	ttl  int64
}

type Cashe struct {
	mu          sync.RWMutex
	pool        sync.Pool
	dataPices   map[string]dataPice
	expiringTTL time.Duration
	signalChan  chan struct{}
}

// на подобии примера из (вроде как 2го урока), разбил ещё на парочку базовых структур

func NewCashe(expiringTTL time.Duration) *Cashe {
	cashe := &Cashe{
		dataPices:   make(map[string]dataPice),
		signalChan:  make(chan struct{}),
		expiringTTL: expiringTTL,
		pool: sync.Pool{
			New: func() interface{} {
				return new(bytes.Buffer)
			},
		},
	}
	go cashe.CleanUp()
	return cashe
}

func (c *Cashe) CleanUp() {
	ticker := time.NewTicker(c.expiringTTL)
	defer ticker.Stop()

	for {
		select {
		case <-c.signalChan:
			return
		case <-ticker.C:
			now := time.Now().Unix()
			c.mu.Lock()
			for k, obj := range c.dataPices {
				if now > obj.ttl {
					delete(c.dataPices, k)
				}
			}
			c.mu.Unlock()
		}
	}
}

func (c *Cashe) Stop() {
	close(c.signalChan)
}

func (c *Cashe) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.dataPices[key] = dataPice{
		data: value,
		ttl:  time.Now().Add(c.expiringTTL).Unix(),
	}
}

func (c *Cashe) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	obj, ok := c.dataPices[key]
	c.mu.RUnlock()
	if !ok {
		return nil, false
	}
	if time.Now().Unix() > obj.ttl {
		return nil, false
	}
	return obj.data, true
}

func (c *Cashe) Delete(key string) {
	c.mu.Lock()
	delete(c.dataPices, key)
	c.mu.Unlock()
}

func (c *Cashe) ToJSON() ([]byte, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	buf := c.pool.Get().(*bytes.Buffer)
	buf.Reset()
	defer c.pool.Put(buf)

	coder := json.NewEncoder(buf)
	err := coder.Encode(c.dataPices)
	if err != nil {
		return nil, err
	}
	res := make([]byte, len(buf.Bytes()))
	copy(res, buf.Bytes())
	return res, nil
}
