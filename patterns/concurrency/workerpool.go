package concurrency

import "sync"

type Task func() error

type Result struct {
	Error error
}

type WorkerPool struct {
	workerCount int
	tasks       chan Task
	results     chan Result
	wg          sync.WaitGroup
	once        sync.Once // Защита от повторного Start
	stopOnce    sync.Once // Защита от повторного Stop
}

func NewWorkerPool(workerCount int) *WorkerPool {
	return &WorkerPool{
		workerCount: workerCount,
		tasks:       make(chan Task),   // Небуферизованный , чтобы не переполнялось при загрузки тасками
		results:     make(chan Result), // Небуферизованный - требует обязательного чтения
	}
}

func (wp *WorkerPool) Start() {
	wp.once.Do(func() {
		// Запускаем воркеры (единожды)
		for i := 0; i < wp.workerCount; i++ {
			wp.wg.Add(1)
			go func() {
				//воркеры читают таски и записывают резалт
				defer wp.wg.Done()
				for task := range wp.tasks {
					err := task()
					wp.results <- Result{Error: err}
				}
			}()
		}

		// Автоматически закрываем results после завершения всех воркеров
		go func() {
			wp.wg.Wait()
			close(wp.results)
		}()
	})
}

func (wp *WorkerPool) Submit(task Task) {
	wp.tasks <- task
}

func (wp *WorkerPool) Stop() {
	//единожды останавливаем
	wp.stopOnce.Do(func() {
		close(wp.tasks)
	})
}

func (wp *WorkerPool) Results() <-chan Result {
	return wp.results
}
