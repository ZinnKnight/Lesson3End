package main

import (
	"errors"
	"sync"
)

/*

# SYNC/COND
## Реализация очереди с ограниченной емкостью на sync.Cond

### Описание задачи
В распределенных системах часто требуется синхронизировать работу продюсеров (добавляющих задачи)
и консьюмеров (обрабатывающих задачи). Очередь с фиксированной емкостью (`BoundedQueue`)
решает следующие проблемы:
- **Блокировка продюсеров** при заполнении очереди.
- **Блокировка консьюмеров** при опустошении очереди.
- **Потокобезопасность** в многогоруточной среде. // много-горутинной
- **Корректное завершение** работы через `Shutdown()`.

**Цель:**
Реализовать очередь, использующую `sync.Cond` для эффективной синхронизации горутин.

---

### Требования
1. Реализация методов:
    - `Put(task interface{})` — блокируется, если очередь заполнена.
    - `Get() interface{}` — блокируется, если очередь пуста.
    - `Shutdown()` — завершает работу очереди.
2. Использование `sync.Cond` и `sync.Mutex` для синхронизации.
3. Гарантия отсутствия гонок и утечек.

*/

// Не совсем осознал почему мы не можем использовать например context, но вот мои почеркушки

var ErrorClosedQueue = errors.New("очередь закрыта")

type BoundedQueue struct {
	capacity        int
	queueData       []interface{}
	mutex           sync.Mutex
	notFullCheckup  *sync.Cond
	notEmptyCheckup *sync.Cond
	closed          bool
}

func NewBoundedQueue(capacity int) *BoundedQueue {
	q := &BoundedQueue{
		capacity:  capacity,
		queueData: make([]interface{}, 0, capacity),
	}
	q.notFullCheckup = sync.NewCond(&q.mutex)
	q.notEmptyCheckup = sync.NewCond(&q.mutex)
	return q
}

func (q *BoundedQueue) Put(task interface{}) error {
	q.mutex.Lock()
	// не уверен что это хорошая практика, но человек в примере что я нашёл так делал и это вроде как ок
	// в данной задаче
	defer q.mutex.Unlock()

	for len(q.queueData) == q.capacity {
		if q.closed {
			return ErrorClosedQueue
		}
		q.notFullCheckup.Wait()
	}
	q.queueData = append(q.queueData, task)
	q.notEmptyCheckup.Signal()
	return nil
}

func (q *BoundedQueue) Get() (interface{}, error) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	for len(q.queueData) == q.capacity && !q.closed {
		q.notFullCheckup.Wait()
	}
	if len(q.queueData) == 0 && q.closed {
		return nil, ErrorClosedQueue
	}
	// если я правильно понимаю как это должно работать, то тут мы вроде держим ссылку на наш объект
	// а значит её нужно будет снимать
	task := q.queueData[0]
	q.queueData = nil
	q.queueData = q.queueData[1:]

	q.notFullCheckup.Signal()
	return task, nil
}

func (q *BoundedQueue) Shutdown() error {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	q.closed = true
	q.notFullCheckup.Broadcast()
	q.notEmptyCheckup.Broadcast()
	return nil
}

// половину кода в наглую скомуниздил из разных примеров как из документации, так и других людей, пока сильно
// туплю с пакетом sync помимо условного мьютекса и waitgroup/errorgroup
