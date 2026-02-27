package main

import (
	"fmt"
	"sync"
)

/*

## Конфигуратор приложения с `sync.Once`

**Описание**
Этот проект реализует **потокобезопасный** менеджер конфигурации, который загружает настройки
**только один раз** при первом запросе.
Используется `sync.Once`, чтобы избежать повторной загрузки при одновременном
доступе из нескольких горутин.
---

**Возможности**
1. Ленивая инициализация – загрузка конфигурации только при первом вызове.
2. Потокобезопасность – отсутствие гонок данных при многопоточной работе.
3. Гибкость – возможность загружать конфигурацию из файла, переменных окружения или базы данных.
---

**Реализованные методы**
- `LoadConfig()` – загружает конфигурацию **один раз** и сохраняет в памяти.
- `Get(key string) string` – возвращает значение конфигурации по ключу.
- `PrintConfig()` – выводит загруженные параметры.
---

```go
// Имитация загрузки конфигурации
cm.config = map[string]string{
"app_name":  "MyApp",
"port":      "8080",
"log_level": "debug",
}

func main() {
    keys := []string{"app_name", "port", "log_level"}
    configManager.PrintConfig()
}
```
*/

type AppConfiguration struct {
	config map[string]string
	conn   sync.Once
}

var (
	conf *AppConfiguration
	once sync.Once
)

func NewAppConfiguration() *AppConfiguration {
	once.Do(func() {
		conf = &AppConfiguration{}
	})
	return conf
}

func (appCf *AppConfiguration) LoadConfig() {
	appCf.conn.Do(func() {
		appCf.config = make(map[string]string)
	})
}

func (appCf *AppConfiguration) Get(key string) string {
	appCf.LoadConfig()
	return appCf.config[key]
}

func (appCf *AppConfiguration) PrintConfig() {
	appCf.LoadConfig()
	for key, value := range appCf.config {
		fmt.Println(key, value)
	}
}

func main() {
	wg := sync.WaitGroup{}
	confManager := NewAppConfiguration()
	keys := []string{"app_name", "port", "log_level"}

	for _, key := range keys {
		wg.Add(1)
		go func() {
			defer wg.Done()
			confManager.Get(key)
		}()
	}
	wg.Wait()
	confManager.PrintConfig()
}
