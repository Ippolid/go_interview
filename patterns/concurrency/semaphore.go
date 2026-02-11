package concurrency

//Семафор — это паттерн синхронизации, который ограничивает количество одновременно выполняющихся горутин

// Семафор на основе буфф канала
type Semaphore struct {
	sem chan struct{}
}

// Создаем семафор нужной емкости
func NewSemaphore(maxItems int) *Semaphore {
	return &Semaphore{
		sem: make(chan struct{}, maxItems),
	}
}

// Захватываем семафор ( если переполнен, то блокируется при записи)
func (s *Semaphore) Acquire() {
	s.sem <- struct{}{}
}

// Освобождаем место
func (s *Semaphore) Release() {
	<-s.sem
}
