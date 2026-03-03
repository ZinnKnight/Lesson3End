package main

import (
	"fmt"
	"sync"
	"time"
)

/*

# Параллельная загрузка данных из нескольких источников
---
## Описание задачи
Реализовать систему параллельной загрузки данных из независимых источников:
1. Асинхронная загрузка комментариев из БД
2. Параллельная загрузка данных пользователей на основе полученных комментариев
3. Загрузка данных сессии и условная загрузка вложений

// Я ваще не понял последний пункт, тип 1,2 ещё ладно, но 3й звучит странно

**Цель:**
Освоить работу с горутинами, `sync.Once` и синхронизацией через `sync.WaitGroup`.

---
## Требования
1. Загрузка комментариев и данных сессии должна выполняться параллельно
2. Загрузка данных пользователей должна стартовать только после получения комментариев
3. Загрузка вложений должна выполняться только при наличии session-id
4. Использовать минимум 3 горутины для разных этапов
5. Синхронизировать все операции перед завершением
Но если что мы с тобой и так пройдем эти темы. А если хочешь прям догнать,то вот дополнительные ресурсы.
Можем отдельно встречу организовать по вопросам::
https://victoriametrics.com/blog/go-sync-once/

*/

// кароч я или балбес или обделенный, тз звучит примерно как собери мини версию вк по рофлу с затычками
// гиперполизированно офк, но я знатно потупил
// Ниже реализовал как смог

type Comment struct {
	ID      int
	UseID   string
	Message string
}

type User struct {
	ID       int
	UserName string
}

type Session struct {
	ID               int
	SessionTimeStomp time.Time
}

// я хз чё ваще имелось под "вложениями" поэтому написал вот такую затычку

type Zakladka struct {
	Name string
}

func main() {
	mu := sync.Mutex{}
	wg := sync.WaitGroup{}
	syncOnes := sync.Once{}

	commChan := make(chan []Comment)
	sessionChan := make(chan Session)

	var sessionID = 1
	var sessionCounter = 1

	wg.Add(1)
	go func() {
		defer wg.Done()
		// возможно стоило бы сделать так что бы коммы по одному подгружались и собирались потом в один слайс, или что-то
		// такое. Но если у нас отправка сессионная, как я понял - то и отправлять можно сразу бахнув слайс
		comments := []Comment{
			{ID: 1, UseID: "Pervonax", Message: "kys"},
			{ID: 2, UseID: "Hitler", Message: "ломай шмотки и в окно"},
			{ID: 3, UseID: "Поздняков.Подписатся", Message: "черти токсичные"},
		}
		commChan <- comments
	}()

	for i := 0; i < len(commChan); i++ {
		wg.Add(3)

		go func() {
			defer wg.Done()
			mu.Lock()
			session := Session{
				ID:               sessionCounter,
				SessionTimeStomp: time.Now(),
			}
			mu.Unlock()

			sessionCounter++
			sessionChan <- session
		}()

		go func() {
			defer wg.Done()

			comments := <-commChan
			var users []User
			for _, comment := range comments {
				users = append(users, User{
					ID:       comment.ID,
					UserName: comment.UseID,
				})
			}
		}()

		go func() {
			defer wg.Done()
			session := <-sessionChan
			sessionID = session.ID

			if sessionID != 0 {
				syncOnes.Do(func() {
					zakladka := Zakladka{
						Name: "Zakladka",
					}
					fmt.Println("Вложение полученно:", zakladka)
				})
			}
			mu.Lock()
			sessionID++
			mu.Unlock()
		}()
	}
	go func() {
		wg.Wait()
		close(commChan)
		close(sessionChan)
	}()

	fmt.Println(<-commChan)
	fmt.Println(<-sessionChan)
}
