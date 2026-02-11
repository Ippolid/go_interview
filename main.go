package main

import (
	"fmt"
	"time"

	"github.com/Ippolid/go_interview/patterns/concurrency"
)

func main() {
	// Создаем генераторы данных
	ch1 := generator("Источник-1", 3)
	ch2 := generator("Источник-2", 4)
	ch3 := generator("Источник-3", 5)

	// Объединяем каналы
	merged := concurrency.FanInOrMergeChannels(ch1, ch2, ch3)

	// Читаем из объединенного канала
	for msg := range merged {
		fmt.Println(msg)
	}

	fmt.Println("Все каналы обработаны")
}

// Вспомогательная функция для генерации данных

func generator(name string, count int) <-chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)
		for i := 1; i <= count; i++ {
			ch <- i
			fmt.Printf("%s отправил: %d\n", name, i)
			time.Sleep(100 * time.Millisecond)
		}
	}()
	return ch
}
