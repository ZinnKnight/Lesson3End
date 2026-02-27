package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

/*

## Задание: Параллельный подсчет слов в файлах с использованием паттерна Fan-Out

### Цель задания
Реализовать параллельную обработку текстовых файлов с использованием паттерна **Fan-Out**,
чтобы ускорить подсчет слов в каждом файле.

### Описание задачи
Есть директория с текстовыми файлами. Нужно:
1. Прочитать все файлы.
2. Распределить их обработку между несколькими горутинами.
3. Подсчитать количество слов в каждом файле.
4. Вывести общую статистику. // не совсем понял про статистику, иммется в виду просто вывод данных, или
// буквально статистика сравнения? Если именно статистика то это тесты и сравнение метрик

### Требования
- Использовать паттерн **Fan-Out** для распределения задач.
- Обработка каждого файла должна выполняться в отдельной горутине.
- Результаты должны агрегироваться в основном потоке.

*/

// Кароч я всё равно не понял как по человечески сделать, так что будет как осилил - коряво

// Стуктурка что представляет собой "хранилище" под дату

type ResultStruct struct {
	FileName   string
	WordAmount int
	Errors     error
}

// worker - по сути сущность, что и будт делать всю работу

func worker(jobChan chan string, resultChan chan ResultStruct, wg *sync.WaitGroup) {
	defer wg.Done()

	for file := range jobChan {
		wc, err := wordsCounter(file)
		resultChan <- ResultStruct{
			FileName:   file,
			WordAmount: wc,
			Errors:     err,
		}
	}
}

// Как я понял читать будем реальные файлы, так что взял пакет os
// Если нужно читать просто из каналов - код поменяется офк, но не на столько что бы всё приложение переписывать

func wordsCounter(file string) (int, error) {
	fileData, err := os.Open(file)
	if err != nil {
		return 0, err
	}
	defer fileData.Close()

	scan := bufio.NewScanner(fileData)
	scan.Split(bufio.ScanWords)

	counter := 0

	for scan.Scan() {
		word := strings.TrimSpace(scan.Text())
		word = strings.ToLower(word)
		if word != "" {
			counter++
		}
	}
	if err := scan.Err(); err != nil {
		return 0, err
	}
	return counter, nil
}

func main() {
	fileDir := "..." // типо директория к файлам

	jobChan := make(chan string)
	resultChan := make(chan ResultStruct)

	wg := sync.WaitGroup{}
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go worker(jobChan, resultChan, &wg)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	go func() {
		filepath.Walk(fileDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				jobChan <- path
			}
			return nil
		})
		close(jobChan)
	}()
	for res := range resultChan {
		if res.Errors != nil {
			fmt.Printf("Ошибка при обработке файлов: %s\n", res.Errors)
		}
		fmt.Printf("Файл: %s, Слов: %d", res.FileName, res.WordAmount)
	}
}
