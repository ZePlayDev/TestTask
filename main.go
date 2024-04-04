package main

import (
	"fmt"
	"sync"
	"time"
)

type Ttype struct {
	id         int
	cT         time.Time // Использую тип Time вместо строки
	fT         time.Time
	taskRESULT []byte
}

func main() {
	superChan := make(chan Ttype, 10)
	doneTasks := make(chan Ttype)
	undoneTasks := make(chan error)

	var wg sync.WaitGroup

	// Создание задач
	wg.Add(1)
	go taskCreator(superChan, &wg)

	// Обработка задач
	for i := 0; i < 5; i++ { // Использую несколько воркеров для обработки задач
		wg.Add(1)
		go taskWorker(superChan, doneTasks, undoneTasks, &wg)
	}

	// Обработка результатов в отдельных горутинах
	go processResults(doneTasks, undoneTasks)

	wg.Wait()
	close(doneTasks)
	close(undoneTasks)

	// Добавлено время для демонстрации результатов обработки
	time.Sleep(1 * time.Second)
}

func taskCreator(tasks chan<- Ttype, wg *sync.WaitGroup) {
	defer wg.Done()
	ticker := time.NewTicker(500 * time.Millisecond) // Генерация задач с интервалом
	defer ticker.Stop()

	for i := 0; i < 20; i++ { // Создаю ограниченное количество задач для демонстрации
		<-ticker.C
		tasks <- Ttype{id: i, cT: time.Now()}
	}
	close(tasks)
}

func taskWorker(tasks <-chan Ttype, doneTasks chan<- Ttype, undoneTasks chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()
	for t := range tasks {
		time.Sleep(150 * time.Millisecond) // Имитация работы
		if time.Now().Unix()%2 == 0 {      // Условие ошибки
			t.taskRESULT = []byte("task has been successed")
			doneTasks <- t
		} else {
			t.taskRESULT = []byte("something went wrong")
			undoneTasks <- fmt.Errorf("Task id %d, error %s", t.id, t.taskRESULT)
		}
	}
}

func processResults(doneTasks <-chan Ttype, undoneTasks <-chan error) {
	for {
		select {
		case task, ok := <-doneTasks:
			if !ok {
				return
			}
			fmt.Printf("Task %d completed successfully\n", task.id)
		case err, ok := <-undoneTasks:
			if !ok {
				return
			}
			fmt.Printf("Error: %s\n", err)
		}
	}
}

/*В этом коде я сделал следующие изменения:
Использовал time.Time для времени создания и завершения задачи вместо строк.
Ввел taskCreator для создания задач с использованием time.Ticker, что позволяет контролировать частоту создания задач.
taskWorker теперь обрабатывает задачи параллельно с помощью нескольких горутин, используя каналы для передачи завершенных и неудачных задач.
Ввел функция processResults для асинхронной обработки результатов задач.
Добавил sync.WaitGroup для контроля за завершением всех горутин перед закрытием каналов и завершением программы.*/
