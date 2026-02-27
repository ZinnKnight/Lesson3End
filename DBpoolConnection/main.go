package main

import (
	"fmt"
	"sync"
	"time"
)

/*

## Пул подключений к БД с использованием `sync.Cond`

### Описание задачи
Реализовать пул подключений к базе данных с ограничением на максимальное количество активных подключений.
Если все подключения заняты, новые запросы должны блокироваться до освобождения ресурсов.
Использовать `sync.Cond` для синхронизации.

---

### Требования
1. Реализовать методы:
    - `Get() *Connection` — возвращает свободное подключение или блокирует горутину.
    - `Release(*Connection)` — освобождает подключение и уведомляет ожидающих.
2. Ограничить максимальное количество подключений (например, 3).
3. Гарантировать потокобезопасность.
4. Смоделировать работу с задержками (имитация запросов к БД).

---

```go
func main() {
    pool := NewConnectionPool(3) // Пул на 3 подключения

    for i := 0; i < 10; i++ {
        go func(id int) {
            conn := pool.Get()
            defer pool.Release(conn)

            fmt.Printf("Горутина %d: подключение %d получено\n", id, conn.ID)
            time.Sleep(2 * time.Second) // Имитация работы
        }(i)
    }

    time.Sleep(10 * time.Second)
}

```
Но если что мы с тобой и так пройдем эти темы. А если хочешь прям догнать,то вот дополнительные ресурсы. Можем отдельно встречу организовать по вопросам::
https://ubiklab.net/posts/go-sync-cond/
https://dev.to/func25/go-synccond-the-most-overlooked-sync-mechanism-1fgd
https://wcademy.ru/go-multithreading-sync-cond/

*/

type Connection struct {
	ID int
}

type ConnectionPool struct {
	mutex            sync.Mutex
	cond             *sync.Cond
	maxConnectToDB   int
	totalConnectToDB int
	freeConnections  []*Connection
}

func NewConnectionPool(maxConnectToDB int) *ConnectionPool {
	pool := &ConnectionPool{
		maxConnectToDB: maxConnectToDB,
	}
	pool.cond = sync.NewCond(&pool.mutex)
	return pool
}

func (pool *ConnectionPool) Get() *Connection {
	pool.mutex.Lock()
	defer pool.mutex.Unlock()

	for {
		if len(pool.freeConnections) > 0 {
			connection := pool.freeConnections[len(pool.freeConnections)-1]
			pool.freeConnections = pool.freeConnections[:len(pool.freeConnections)-1]
			return connection
		}

		if pool.totalConnectToDB < pool.maxConnectToDB {
			pool.totalConnectToDB++
			return &Connection{
				ID: pool.totalConnectToDB,
			}
		}

		pool.cond.Wait()
	}
}

func (pool *ConnectionPool) Release(con *Connection) {
	pool.mutex.Lock()
	defer pool.mutex.Unlock()

	pool.freeConnections = append(pool.freeConnections, con)
	pool.cond.Signal()
}

func main() {
	pool := NewConnectionPool(3) // Пул на 3 подключения

	for i := 0; i < 10; i++ {
		go func(id int) {
			conn := pool.Get()
			defer pool.Release(conn)

			fmt.Printf("Горутина %d: подключение %d получено\n", id, conn.ID)
			time.Sleep(2 * time.Second) // Имитация работы
		}(i)
	}
	time.Sleep(10 * time.Second)
}
