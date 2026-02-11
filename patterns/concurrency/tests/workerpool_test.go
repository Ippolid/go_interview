package tests

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Ippolid/go_interview/patterns/concurrency"
)

// TestWorkerPoolBasic - базовый тест обработки задач
func TestWorkerPoolBasic(t *testing.T) {
	fmt.Println("Тест : Базовая обработка задач")

	pool := concurrency.NewWorkerPool(3)
	pool.Start()

	var completed int32
	const totalTasks = 10

	// Запускаем чтение результатов ДО отправки задач
	done := make(chan struct{})
	go func() {
		for range pool.Results() {
		}
		close(done)
	}()

	// Отправляем задачи
	for i := 1; i <= totalTasks; i++ {
		taskID := i
		pool.Submit(func() error {
			fmt.Printf("Выполняется задача %d\n", taskID)
			time.Sleep(50 * time.Millisecond)
			atomic.AddInt32(&completed, 1)
			return nil
		})
	}

	pool.Stop()
	<-done // Ждем завершения чтения

	total := atomic.LoadInt32(&completed)
	if total != totalTasks {
		t.Errorf("Ожидалось %d задач, выполнено %d", totalTasks, total)
	}

	fmt.Printf("Тест пройден: выполнено %d/%d задач\n", total, totalTasks)
}

// TestWorkerPoolErrors - тест обработки ошибок
func TestWorkerPoolErrors(t *testing.T) {
	fmt.Println("Тест: Обработка ошибок")

	pool := concurrency.NewWorkerPool(2)
	pool.Start()

	var successCount, errorCount int32

	// Читаем результаты ДО отправки задач
	var resultsRead int32
	done := make(chan struct{})
	go func() {
		for result := range pool.Results() {
			atomic.AddInt32(&resultsRead, 1)
			if result.Error != nil {
				fmt.Printf("Получена ошибка: %v\n", result.Error)
			}
		}
		close(done)
	}()

	// Отправляем задачи с ошибками
	for i := 1; i <= 10; i++ {
		taskID := i
		pool.Submit(func() error {
			if taskID%3 == 0 {
				atomic.AddInt32(&errorCount, 1)
				return fmt.Errorf("ошибка в задаче %d", taskID)
			}
			atomic.AddInt32(&successCount, 1)
			return nil
		})
	}

	pool.Stop()
	<-done

	success := atomic.LoadInt32(&successCount)
	errors := atomic.LoadInt32(&errorCount)
	results := atomic.LoadInt32(&resultsRead)

	fmt.Printf("Успешных задач: %d\n", success)
	fmt.Printf("Задач с ошибками: %d\n", errors)
	fmt.Printf("Результатов прочитано: %d\n", results)

	if success+errors != 10 {
		t.Errorf("Не все задачи выполнены: %d успешных + %d с ошибками", success, errors)
	}

	if results != 10 {
		t.Errorf("Не все результаты получены: %d из 10", results)
	}

	fmt.Println("Тест пройден: ошибки обработаны корректно")
}

// TestWorkerPoolHighLoad - тест высокой нагрузки
func TestWorkerPoolHighLoad(t *testing.T) {
	fmt.Println("Тест: Высокая нагрузка")

	const numWorkers = 5
	const totalTasks = 1000

	pool := concurrency.NewWorkerPool(numWorkers)
	pool.Start()

	var completed int32
	start := time.Now()

	// Читаем результаты ДО отправки задач
	done := make(chan struct{})
	go func() {
		for range pool.Results() {
		}
		close(done)
	}()

	// Отправляем много задач
	for i := 1; i <= totalTasks; i++ {
		pool.Submit(func() error {
			time.Sleep(time.Millisecond)
			atomic.AddInt32(&completed, 1)
			return nil
		})
	}

	pool.Stop()
	<-done
	elapsed := time.Since(start)

	total := atomic.LoadInt32(&completed)
	fmt.Printf("Выполнено задач: %d/%d\n", total, totalTasks)
	fmt.Printf("Время выполнения: %v\n", elapsed)

	if total != totalTasks {
		t.Errorf("Не все задачи выполнены: %d из %d", total, totalTasks)
	}

	fmt.Println("✅ Тест пройден: высокая нагрузка обработана")
}
