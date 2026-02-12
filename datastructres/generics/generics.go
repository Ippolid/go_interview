package generics

import (
	"golang.org/x/exp/constraints"
)

// Числовые типы
type Number interface {
	constraints.Float | constraints.Integer | constraints.Complex
}

// умножение любого числа на 2
func Double[T Number](value T) T {
	return value * 2
}

// Сумма проивзедений 2 числовых слайсов
func DotProduct[T Number](s1, s2 []T) T {
	if len(s1) != len(s2) {
		panic("DotProduct: slices of unequal length")
	}
	var r T
	for i := range s1 {
		r += s1[i] * s2[i]
	}
	return r
}

// Сумма значений мапы
func Sum[K comparable, V Number](m map[K]V) V {
	var s V
	for _, v := range m {
		s += v
	}
	return s
}

// Делаем из слайса маппу
func MapSlice[T any, R any](items []T, mapper func(T) R) []R {
	result := make([]R, len(items))
	for i, item := range items {
		result[i] = mapper(item)
	}
	return result
}

// Фильтр любого слайса по зданной функции
func FilterSlice[T any](items []T, keep func(T) bool) []T {
	result := make([]T, 0, len(items))
	for _, item := range items {
		if keep(item) {
			result = append(result, item)
		}
	}
	return result
}

// Возвращаем переданное значение
func Identity[T any](value T) T {
	return value
}

// Поиск значения в слайсе
func Contains[T comparable](items []T, target T) bool {
	for _, item := range items {
		if item == target {
			return true
		}
	}
	return false
}
