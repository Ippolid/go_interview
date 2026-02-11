package tests

import (
	"fmt"
	"sort"
	"testing"

	"github.com/Ippolid/go_interview/patterns/concurrency"
)

// TestFanInBasic - базовый тест на слияние каналов
func TestFanInBasic(t *testing.T) {
	fmt.Println("Тест 1: Базовое слияние каналов")

	// Создаем 3 канала
	ch1 := make(chan int, 2)
	ch2 := make(chan int, 2)
	ch3 := make(chan int, 2)

	// Заполняем каналы
	ch1 <- 1
	ch1 <- 2
	close(ch1)

	ch2 <- 3
	ch2 <- 4
	close(ch2)

	ch3 <- 5
	ch3 <- 6
	close(ch3)

	// Объединяем
	merged := concurrency.FanInOrMergeChannels(ch1, ch2, ch3)

	// Собираем результаты
	var results []int
	for val := range merged {
		fmt.Printf("Получено значение: %d\n", val)
		results = append(results, val)
	}

	// Проверяем количество
	if len(results) != 6 {
		t.Errorf("Ожидалось 6 значений, получено %d", len(results))
	}

	// Проверяем, что все значения есть
	sort.Ints(results)
	expected := []int{1, 2, 3, 4, 5, 6}
	for i, v := range expected {
		if results[i] != v {
			t.Errorf("Ожидалось %d, получено %d", v, results[i])
		}
	}

	fmt.Println("Тест пройден: все значения получены")
}

// TestFanInWithStrings - тест с другим типом данных
func TestFanInWithStrings(t *testing.T) {
	fmt.Println("Тест 2: Работа с string")

	ch1 := make(chan string, 2)
	ch2 := make(chan string, 2)

	ch1 <- "Hello"
	ch1 <- "World"
	close(ch1)

	ch2 <- "Go"
	ch2 <- "Patterns"
	close(ch2)

	merged := concurrency.FanInOrMergeChannels(ch1, ch2)

	var results []string
	for str := range merged {
		fmt.Printf("Получена строка: %s\n", str)
		results = append(results, str)
	}

	if len(results) != 4 {
		t.Errorf("Ожидалось 4 строки, получено %d", len(results))
	}

	fmt.Println(" Тест пройден: generics работают корректно")
}
