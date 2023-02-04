package ml_slice

import (
	"math/rand"
	"time"
)

func IndexOf[T comparable](slice []T, v T) int {
	for i := 0; i < len(slice); i++ {
		if slice[i] == v {
			return i
		}
	}

	return -1
}

func IndexBy[T any](slice []T, fn func(*T) bool) int {
	for i := 0; i < len(slice); i++ {
		if fn(&slice[i]) {
			return i
		}
	}

	return -1
}

func RemoveAt[T any](slice []T, index int) []T {
	return append(slice[:index], slice[index+1:]...)
}

func Includes[T comparable](slice []T, v T) bool {
	for i := 0; i < len(slice); i++ {
		if slice[i] == v {
			return true
		}
	}

	return false
}

func Unique[T comparable](slice []T) []T {
	keys := make(map[T]bool)
	list := make([]T, 0)

	for i := 0; i < len(slice); i++ {
		if _, value := keys[slice[i]]; !value {
			keys[slice[i]] = true
			list = append(list, slice[i])
		}
	}

	return list
}

func UniqueBy[T any](slice []T, fn func(*T) any) []T {
	keys := make(map[any]T)
	list := make([]T, 0)

	for i := 0; i < len(slice); i++ {
		keys[fn(&slice[i])] = slice[i]
	}

	for _, v := range keys {
		list = append(list, v)
	}

	return list
}

func FilterBy[T any](slice []T, filter func(*T) bool) []T {
	filtered := make([]T, 0)

	for i := 0; i < len(slice); i++ {
		if filter(&slice[i]) {
			filtered = append(filtered, slice[i])
		}
	}
	return filtered
}

func ToAny[T any](slice []T) []any {
	filtered := make([]any, len(slice))

	for i := 0; i < len(slice); i++ {
		filtered[i] = any(slice[i])
	}

	return filtered
}

func Find[T any](slice []T, filter func(*T) bool) (T, bool) {
	for i := 0; i < len(slice); i++ {
		if filter(&slice[i]) {
			return slice[i], true
		}
	}
	return *new(T), false
}

func GetKeys[K comparable, V comparable](mp map[K]V) []K {
	l := make([]K, 0)
	for k, _ := range mp {
		l = append(l, k)
	}
	return l
}

func Paginate[T any](slice []T, offset int, limit int) []T {
	if offset < 0 {
		offset = 0
	}
	if limit < 0 {
		limit = 0
	}
	if offset > len(slice) {
		offset = len(slice)
	}

	end := offset + limit
	if end > len(slice) {
		end = len(slice)
	}
	return slice[offset:end]
}

func Map[T any, R any](slice []T, m func(T) R) []R {
	mapped := make([]R, 0)
	for _, v := range slice {
		mapped = append(mapped, m(v))
	}
	return mapped
}

var r = rand.New(rand.NewSource(time.Now().Unix()))

func PickRandom[T any](slice []T) T {
	return slice[r.Intn(len(slice))]
}

func Combine[T any](slices ...[]T) []T {
	finalSlice := make([]T, 0)

	for i := 0; i < len(slices); i++ {
		for j := 0; j < len(slices[i]); j++ {
			finalSlice = append(finalSlice, slices[i][j])
		}
	}

	return finalSlice
}
