package ml_slice

import (
	"golang.org/x/exp/constraints"
	"math/rand"
	"sort"
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
	out := make([]T, 0, len(slice))
	for i := 0; i < len(slice); i++ {
		if i == index {
			continue
		}
		out = append(out, slice[i])
	}
	return out
	// return append(slice[:index], slice[index+1:]...)
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
	if len(slice) == 0 {
		return *new(T)
	}
	return slice[r.Intn(len(slice))]
}

func PickRandomIndex[T any](slice []T) int {
	if len(slice) == 0 {
		return -1
	}
	return r.Intn(len(slice))
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

func Prepend[T any](slice []T, value []T) []T {
	return append(value, slice...)
}

func NotNil[T any](slice []T) []T {
	if slice == nil {
		return make([]T, 0)
	}
	return slice
}

func SortAZ[T constraints.Ordered](slice []T) {
	fn := func(i, j int) (T, T) { return slice[i], slice[j] }

	sort.SliceStable(slice, func(i, j int) bool {
		a, b := fn(i, j)
		return a < b
	})
}

func SortZA[T constraints.Ordered](slice []T) {
	fn := func(i, j int) (T, T) { return slice[i], slice[j] }

	sort.SliceStable(slice, func(i, j int) bool {
		a, b := fn(i, j)
		return a > b
	})
}

func SortAZBy[T any, N constraints.Ordered](slice []T, fn func(i int, j int) (N, N)) {
	sort.SliceStable(slice, func(i, j int) bool {
		a, b := fn(i, j)
		return a < b
	})
}

func SortZABy[T any, N constraints.Ordered](slice []T, fn func(i int, j int) (N, N)) {
	sort.SliceStable(slice, func(i, j int) bool {
		a, b := fn(i, j)
		return a > b
	})
}

func Reverse[T any](slice []T) []T {
	out := make([]T, len(slice))
	copy(out, slice)

	for i, j := 0, len(out)-1; i < j; i, j = i+1, j-1 {
		out[i], out[j] = out[j], out[i]
	}

	return out
}

// Permute the values at index i to len(a)-1.
func perm[T any](slice []T, f func([]T), i int) {
	if i > len(slice) {
		f(slice)
		return
	}
	perm(slice, f, i+1)
	for j := i + 1; j < len(slice); j++ {
		slice[i], slice[j] = slice[j], slice[i]
		perm(slice, f, i+1)
		slice[i], slice[j] = slice[j], slice[i]
	}
}

func Permutation[T any](slice []T) [][]T {
	out := make([][]T, 0, len(slice))
	perm(slice, func(a []T) {
		n := make([]T, len(a))
		copy(n, a)
		out = append(out, n)
	}, 0)
	return out
}

// PullFirst elements and remove from slice
func PullFirst[T any](slice *[]T, n int) []T {
	if n >= len(*slice) {
		result := append([]T{}, *slice...)
		*slice = (*slice)[:0] // очищаем исходный срез
		return result
	}

	result := (*slice)[:n]
	*slice = (*slice)[n:]
	return result
}
