package concurrency

import (
	"fmt"
	"sync"
	"time"
)

// воркер, результаты, которых мы должны собрать
func worker(id int, jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		fmt.Printf("Worker %d начал работу над задачей %d\n", id, job)
		time.Sleep(time.Second) // эмуляция работы
		results <- job * 2
		fmt.Printf("Worker %d завершил задачу %d\n", id, job)
	}
}

func FanOut(numWorkers int, jobs <-chan int) <-chan int {
	results := make(chan int, numWorkers)
	var wg sync.WaitGroup

	// Запускаем воркеры
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(w, jobs, results, &wg)
	}

	// Горутина для закрытия канала результатов
	go func() {
		wg.Wait()
		close(results)
	}()

	return results
}

func main() {
	jobs := make(chan int, 10)

	// Запускаем fan-out с 3 воркерами
	results := FanOut(4, jobs)

	// Отправляем задачи
	for j := 1; j <= 9; j++ {
		jobs <- j
	}
	close(jobs)

	// Читаем результаты
	for result := range results {
		fmt.Println("Результат:", result)
	}
}
