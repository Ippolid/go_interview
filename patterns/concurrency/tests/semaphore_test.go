package tests

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Ippolid/go_interview/patterns/concurrency"
)

// TestSemaphoreBasic - базовый тест ограничения параллелизма
func TestSemaphoreBasic(t *testing.T) {
	fmt.Println("Тест 1: Базовое ограничение параллелизма")

	sem := concurrency.NewSemaphore(2) // максимум 2 одновременных задачи
	var active int32
	var maxActive int32

	var wg sync.WaitGroup

	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			sem.Acquire()
			defer sem.Release()

			// Увеличиваем счетчик активных горутин
			current := atomic.AddInt32(&active, 1)
			fmt.Printf("Задача %d начата (активных: %d)\n", id, current)

			// Обновляем максимум
			for {
				max := atomic.LoadInt32(&maxActive)
				if current <= max || atomic.CompareAndSwapInt32(&maxActive, max, current) {
					break
				}
			}

			time.Sleep(100 * time.Millisecond)

			atomic.AddInt32(&active, -1)
			fmt.Printf("Задача %d завершена\n", id)
		}(i)
	}

	wg.Wait()

	max := atomic.LoadInt32(&maxActive)
	fmt.Printf("Максимум одновременных задач: %d\n", max)

	if max > 2 {
		t.Errorf("Семафор не ограничивает параллелизм: максимум был %d вместо 2", max)
	}

	fmt.Println("Тест пройден: параллелизм ограничен корректно")
}
