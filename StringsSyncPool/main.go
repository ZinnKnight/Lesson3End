package main

import (
	"fmt"
	"sync"
)

/*

# SYNC/POOL
## Оптимизация обработки строк с sync.Pool

### Описание задачи
В высоконагруженном сервисе частые аллокации буферов для преобразования строк создают нагрузку на GC.
Цель — реализовать оптимизированную функцию `ProcessString` с использованием `sync.Pool`,
чтобы переиспользовать буферы `[]byte`.

### Требования
1. Функция `ProcessString(s string) string` преобразует строку в верхний регистр.
2. Использование `sync.Pool` для буферов `[]byte`.
3. Потокобезопасность, отсутствие утечек памяти.
```go
func main() {
	examples := []string{
		"hello, world!",
		"gopher",
		"lorem ipsum dolor sit amet",
	}

	for _, s := range examples {
		processed := ProcessString(s)
		fmt.Printf("Original: %q\nProcessed: %q\n\n", s, processed)
	}
}
```

*/

var SincPool = sync.Pool{
	New: func() interface{} {
		b := make([]byte, 0, 64)
		return &b
	},
}

func ProcessString(s string) string {
	bufPool := SincPool.Get().(*[]byte)
	bufP := *bufPool
	// в процессе увидел проблему (а так же подглядел в коде у другого человека)
	// что гонять без присваивания к переменной не получится, т.к не смогу вернуть потом под следющие операции
	bufP = bufP[:0]

	for i := 0; i < len(s); i++ {
		// еслия я правильно понял, то синтаксического сахара, или чего-то схожего со strings
		//сюда запихнуть не выйдет, без корявых костылей, так что я взял переписанный вариант что нашёл в инете
		check := s[i]
		// тут английский алфавит подразумеваются, но я полагаю что увеличив объём выделенной памяти можно
		// и любой другой засунуть
		if check >= 'a' && check <= 'z' {
			check = check - 'a' + 'A'
		}
		bufP = append(bufP, check)
	}
	res := string(bufP)

	*bufPool = bufP
	SincPool.Put(bufPool)
	return res
}
func main() {
	examples := []string{
		// ради интереса немного поменял изначальный тест, просто подвердить то что хотел
		"Привет, Мир!",
		"gopher",
		"lorem ipsum dolor sit amet",
	}

	for _, s := range examples {
		processed := ProcessString(s)
		fmt.Printf("Original: %q\nProcessed: %q\n\n", s, processed)
	}
}
