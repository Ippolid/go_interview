package tests

import (
	"fmt"
	"testing"

	"github.com/Ippolid/go_interview/patterns/concurrency"
)

// TestFanOutBasic - базовый тест распределения задач
func TestFanOutBasic(t *testing.T) {
	fmt.Println("Тест 1: Базовое распределение задач")

	jobs := make(chan int, 5)

	// Отправляем 5 задач
	for i := 1; i <= 5; i++ {
		jobs <- i
		fmt.Printf("Отправлена задача: %d\n", i)
	}
	close(jobs)

	// Запускаем 2 воркера
	results := concurrency.FanOut(2, jobs)

	// Собираем результаты
	resultMap := make(map[int]bool)
	for result := range results {
		fmt.Printf("Получен результат: %d\n", result)
		resultMap[result] = true
	}

	// Проверяем, что все задачи обработаны (job * 2)
	expected := []int{2, 4, 6, 8, 10}
	if len(resultMap) != len(expected) {
		t.Errorf("Ожидалось %d результатов, получено %d", len(expected), len(resultMap))
	}

	for _, exp := range expected {
		if !resultMap[exp] {
			t.Errorf("Отсутствует результат: %d", exp)
		}
	}

	fmt.Println("Тест пройден: все задачи обработаны")
}

// TestFanOutMultipleWorkers - тест с разным количеством воркеров
func TestFanOutMultipleWorkers(t *testing.T) {
	fmt.Println("Тест 2: Различное количество воркеров")

	testCases := []struct {
		workers int
		jobs    int
	}{
		{1, 3},  // 1 воркер, 3 задачи
		{3, 6},  // 3 воркера, 6 задач
		{5, 10}, // 5 воркеров, 10 задач
	}

	for _, tc := range testCases {
		fmt.Printf("\n Тест: %d воркеров, %d задач\n", tc.workers, tc.jobs)

		jobsChan := make(chan int, tc.jobs)
		for i := 1; i <= tc.jobs; i++ {
			jobsChan <- i
		}
		close(jobsChan)

		results := concurrency.FanOut(tc.workers, jobsChan)

		count := 0
		for result := range results {
			count++
			fmt.Printf("Результат: %d\n", result)
		}

		if count != tc.jobs {
			t.Errorf("Воркеров: %d, ожидалось %d результатов, получено %d",
				tc.workers, tc.jobs, count)
		}

		fmt.Printf("Обработано %d задач %d воркерами\n", count, tc.workers)
	}
}
