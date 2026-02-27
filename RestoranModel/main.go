package main

import "sync"

/*

## Моделирование работы ресторана с использованием `sync.Cond`

### Описание задачи
Реализовать систему управления столиками в ресторане, где:
- Количество столиков фиксировано (например, 5).
- Посетители (горутины) занимают столики, если они свободны.
- Если все столики заняты, посетители ожидают в очереди.
- При освобождении столика его получает первый ожидающий посетитель.


**Цель:**
Научиться синхронизировать горутины с помощью `sync.Cond`,
моделируя реальный сценарий с ограниченными ресурсами.

---

### Требования
1. Реализовать структуру `Restaurant` с методами:
    - `OccupyTable()` — блокируется, если нет свободных столиков.
    - `ReleaseTable()` — освобождает столик и уведомляет ожидающих.
2. Использовать `sync.Cond` для управления очередью ожидания.

*/

type Restaurant struct {
	tables    int
	available int
	queue     []int
	mutex     sync.Mutex
	cond      *sync.Cond
}

// как с прошлой задачкой определил структуру и что-то типо билдера для неё

func NewRestoran(tables int) *Restaurant {
	r := &Restaurant{
		tables:    tables,
		available: tables,
		queue:     make([]int, tables),
	}
	r.cond = sync.NewCond(&r.mutex)
	return r
}

func (r *Restaurant) OccupyTable(id int) {
	r.mutex.Lock()

	r.queue = append(r.queue, id)
	// я не совсем понял к чему привязать for, что бы не билдить полу-бесконечный цикл, но пока ничего не придумал
	// так что будет пока что так, в дальнейшем вернусь поумневшим и сделаю лучше
	for {
		if r.available >= id && r.queue[0] == id {
			break
		}
		r.cond.Wait()
	}
	r.queue = r.queue[1:]
	r.available--

	r.mutex.Unlock()
}

func (r *Restaurant) ReleaseTable() {
	r.mutex.Lock()
	r.available++

	r.cond.Broadcast()
	r.mutex.Unlock()
}

// кароч пакет sync на 3х функциях не заканчивается, кто бы мог подумать, буду теперь долбить активно cond и once
