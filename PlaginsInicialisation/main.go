package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

/*



## Инициализация плагинов с `sync.Once`

**Цель задания**
Реализовать систему безопасной инициализации плагинов, где:
- Каждый плагин инициализируется **только один раз**
- Инициализация потокобезопасна
- Ошибки при инициализации корректно обрабатываются
- Плагины доступны для использования из разных компонентов
---

**Требования**
1. **Структура `PluginManager`**:
    - Хранит загруженные плагины
    - Использует `sync.Once` для каждого плагина
    - Поддерживает конкурентный доступ

2. **Методы**:
    - `GetPlugin(name string) (Plugin, error)` – возвращает инициализированный плагин
    - `RegisterPlugin()` – регистрирует плагины (симуляция)
```golang
package main

import (
	"fmt"
	"log"
	"sync"
)

// Интерфейс для всех плагинов
type Plugin interface {
	Execute() string
}

// Управляет инициализацией и доступом к плагинам
type PluginManager struct {
	plugins map[string]*pluginEntry
	mu      sync.RWMutex
}

type pluginEntry struct {
	//Добавить необходимые поля для однократной инициализации
	initFn func() (Plugin, error)
}

// NewPluginManager создает новый менеджер плагинов
func NewPluginManager() *PluginManager {
	return &PluginManager{
		plugins: make(map[string]*pluginEntry),
	}
}

// RegisterPlugin регистрирует новый плагин
func (pm *PluginManager) RegisterPlugin(name string, initFn func() (Plugin, error)) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.plugins[name] = &pluginEntry{
		initFn: initFn,
	}
}

// GetPlugin возвращает инициализированный плагин
func (pm *PluginManager) GetPlugin(name string) (Plugin, error) {
	// Реализовать:
	// 1. Проверку существования плагина
	// 2. Потокобезопасную однократную инициализацию
	// 3. Обработку и кэширование ошибок
	// 4. Возврат кэшированного результата
	return nil, fmt.Errorf("not implemented")
}

// DemoPlugin реализация плагина
type DemoPlugin struct{}

func (p *DemoPlugin) Execute() string {
	return "DemoPlugin executed successfully!"
}

func initDemo() (Plugin, error) {
	// Имитация длительной инициализации
	// time.Sleep(500 * time.Millisecond)
	return &DemoPlugin{}, nil
}

func main() {
	pm := NewPluginManager()

	pm.RegisterPlugin("demo", initDemo)
	pm.RegisterPlugin("broken", func() (Plugin, error) {
		return nil, fmt.Errorf("simulated error")
	})

	var wg sync.WaitGroup

	// Тестирование рабочего плагина
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			p, err := pm.GetPlugin("demo")
			if err != nil {
				log.Printf("Goroutine %d error: %v", id, err)
				return
			}
			log.Printf("Goroutine %d: %s", id, p.Execute())
		}(i)
	}

	// Тестирование плагина с ошибкой
	for i := 5; i < 7; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			_, err := pm.GetPlugin("broken")
			if err != nil {
				log.Printf("Goroutine %d error: %v", id, err)
			}
		}(i)
	}

	wg.Wait()
}
```
Но если что мы с тобой и так пройдем эти темы. А если хочешь прям догнать,то вот дополнительные ресурсы. Можем отдельно встречу организовать по вопросам::
https://victoriametrics.com/blog/go-sync-once/
https://dev.to/jones_charles_ad50858dbc0/a-developers-guide-to-synconce-your-go-concurrency-lifesaver-3kf2
https://backendinterview.ru/goLang/concurrency/sync.html


*/

type Plugin interface {
	Execute() string
}

type pluginEntry struct {
	plugOnce sync.Once
	initFn   func() (Plugin, error)
	plugin   Plugin
	err      error
}

type PluginManager struct {
	plugins map[string]*pluginEntry
	mu      sync.RWMutex
}

func NewPluginManager() *PluginManager {
	return &PluginManager{
		plugins: make(map[string]*pluginEntry),
	}
}

func (pm *PluginManager) RegisterPlugin(name string, initFn func() (Plugin, error)) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.plugins[name] = &pluginEntry{
		initFn: initFn,
	}
}

func (pm *PluginManager) GetPlugin(name string) (Plugin, error) {
	pm.mu.RLock()
	plg, ok := pm.plugins[name]
	pm.mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("plugin %s not found", name)
	}
	plg.plugOnce.Do(func() {
		plg.plugin, plg.err = plg.initFn()
	})
	return plg.plugin, plg.err
}

type DemoPlugin struct{}

func (p *DemoPlugin) Execute() string {
	return "DemoPlugin executed successfully!"
}

func initDemo() (Plugin, error) {
	// Имитация длительной инициализации // эта часть была закоменченна - раскомил что б посмотреть
	time.Sleep(500 * time.Millisecond)
	return &DemoPlugin{}, nil
}

func main() {
	pm := NewPluginManager()

	pm.RegisterPlugin("demo", initDemo)
	pm.RegisterPlugin("broken", func() (Plugin, error) {
		return nil, fmt.Errorf("simulated error")
	})

	var wg sync.WaitGroup

	// Тестирование рабочего плагина
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			p, err := pm.GetPlugin("demo")
			if err != nil {
				log.Printf("Goroutine %d error: %v", id, err)
				return
			}
			log.Printf("Goroutine %d: %s", id, p.Execute())
		}(i)
	}

	// Тестирование плагина с ошибкой
	for i := 5; i < 7; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			_, err := pm.GetPlugin("broken")
			if err != nil {
				log.Printf("Goroutine %d error: %v", id, err)
			}
		}(i)
	}

	wg.Wait()
}
