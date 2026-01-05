package utils

// Contains 检查切片是否包含指定元素
func Contains[T comparable](slice []T, item T) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

// Unique 去重切片
func Unique[T comparable](slice []T) []T {
	seen := make(map[T]struct{})
	result := make([]T, 0, len(slice))

	for _, item := range slice {
		if _, exists := seen[item]; !exists {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}

	return result
}

// Filter 过滤切片
func Filter[T any](slice []T, predicate func(T) bool) []T {
	result := make([]T, 0)
	for _, item := range slice {
		if predicate(item) {
			result = append(result, item)
		}
	}
	return result
}

// Map 映射切片
func Map[T any, R any](slice []T, mapper func(T) R) []R {
	result := make([]R, len(slice))
	for i, item := range slice {
		result[i] = mapper(item)
	}
	return result
}

// Reduce 归约切片
func Reduce[T any, R any](slice []T, initial R, reducer func(R, T) R) R {
	result := initial
	for _, item := range slice {
		result = reducer(result, item)
	}
	return result
}

// Find 查找第一个满足条件的元素
func Find[T any](slice []T, predicate func(T) bool) (T, bool) {
	for _, item := range slice {
		if predicate(item) {
			return item, true
		}
	}
	var zero T
	return zero, false
}

// FindIndex 查找第一个满足条件的元素索引
func FindIndex[T any](slice []T, predicate func(T) bool) int {
	for i, item := range slice {
		if predicate(item) {
			return i
		}
	}
	return -1
}

// Any 检查是否有任意元素满足条件
func Any[T any](slice []T, predicate func(T) bool) bool {
	for _, item := range slice {
		if predicate(item) {
			return true
		}
	}
	return false
}

// All 检查是否所有元素都满足条件
func All[T any](slice []T, predicate func(T) bool) bool {
	for _, item := range slice {
		if !predicate(item) {
			return false
		}
	}
	return true
}

// Chunk 将切片分块
func Chunk[T any](slice []T, size int) [][]T {
	if size <= 0 {
		return nil
	}

	var chunks [][]T
	for i := 0; i < len(slice); i += size {
		end := i + size
		if end > len(slice) {
			end = len(slice)
		}
		chunks = append(chunks, slice[i:end])
	}

	return chunks
}

// Reverse 反转切片
func Reverse[T any](slice []T) []T {
	result := make([]T, len(slice))
	for i, item := range slice {
		result[len(slice)-1-i] = item
	}
	return result
}

// Intersection 求两个切片的交集
func Intersection[T comparable](slice1, slice2 []T) []T {
	set := make(map[T]struct{})
	for _, item := range slice1 {
		set[item] = struct{}{}
	}

	result := make([]T, 0)
	seen := make(map[T]struct{})
	for _, item := range slice2 {
		if _, exists := set[item]; exists {
			if _, added := seen[item]; !added {
				result = append(result, item)
				seen[item] = struct{}{}
			}
		}
	}

	return result
}

// Difference 求两个切片的差集（在slice1但不在slice2中的元素）
func Difference[T comparable](slice1, slice2 []T) []T {
	set := make(map[T]struct{})
	for _, item := range slice2 {
		set[item] = struct{}{}
	}

	result := make([]T, 0)
	for _, item := range slice1 {
		if _, exists := set[item]; !exists {
			result = append(result, item)
		}
	}

	return result
}

// Union 求两个切片的并集
func Union[T comparable](slice1, slice2 []T) []T {
	set := make(map[T]struct{})
	result := make([]T, 0)

	for _, item := range slice1 {
		if _, exists := set[item]; !exists {
			set[item] = struct{}{}
			result = append(result, item)
		}
	}

	for _, item := range slice2 {
		if _, exists := set[item]; !exists {
			set[item] = struct{}{}
			result = append(result, item)
		}
	}

	return result
}

// Flatten 展平二维切片
func Flatten[T any](slice [][]T) []T {
	result := make([]T, 0)
	for _, subSlice := range slice {
		result = append(result, subSlice...)
	}
	return result
}

// GroupBy 根据键函数分组
func GroupBy[T any, K comparable](slice []T, keyFunc func(T) K) map[K][]T {
	result := make(map[K][]T)
	for _, item := range slice {
		key := keyFunc(item)
		result[key] = append(result[key], item)
	}
	return result
}

// Partition 将切片分为两组（满足条件和不满足条件）
func Partition[T any](slice []T, predicate func(T) bool) ([]T, []T) {
	trueSlice := make([]T, 0)
	falseSlice := make([]T, 0)

	for _, item := range slice {
		if predicate(item) {
			trueSlice = append(trueSlice, item)
		} else {
			falseSlice = append(falseSlice, item)
		}
	}

	return trueSlice, falseSlice
}
