package utils

import (
	"math"
)

// SafeSumInt64 sums a slice of int64, returning math.MaxInt64 on overflow.
// Note: Go's default behavior is to wrap around on integer overflow.
// This function explicitly checks for overflow to mimic .NET's checked sum behavior
// if that level of fidelity is required. For simple counts, direct sum might be fine.
func SafeSumInt64(source []int64) int64 {
	var sum int64
	for _, val := range source {
		if val > 0 && sum > math.MaxInt64-val {
			return math.MaxInt64 // Overflow
		}
		if val < 0 && sum < math.MinInt64-val {
			// return math.MinInt64 // Underflow (though C# SafeSum returns MaxValue)
			// To strictly mimic C# (always MaxValue on overflow/underflow):
			if (val > 0 && sum > math.MaxInt64-val) || (val < 0 && sum < math.MinInt64-val) {
				return math.MaxInt64
			}
		}
		sum += val
	}
	return sum
}

// SafeSumInt for int type.
func SafeSumInt(source []int) int {
	var sum int
	for _, val := range source {
		if val > 0 && sum > math.MaxInt-val {
			return math.MaxInt
		}
		if val < 0 && sum < math.MinInt-val {
			// return math.MinInt // or math.MaxInt to mimic C#
			return math.MaxInt
		}
		sum += val
	}
	return sum
}

// TakeLast returns the last 'count' elements from a slice.
// If count is larger than the slice length, returns the whole slice.
// If count is non-positive, returns an empty slice.
func TakeLast[T any](source []T, count int) []T {
	if count <= 0 {
		return []T{}
	}
	if count >= len(source) {
		return append([]T{}, source...) // Return a copy
	}
	return append([]T{}, source[len(source)-count:]...) // Return a copy
}

// ToSet converts a slice of comparable elements into a set (map[T]struct{}).
func ToSet[T comparable](slice []T) map[T]struct{} {
	set := make(map[T]struct{}, len(slice))
	for _, item := range slice {
		set[item] = struct{}{}
	}
	return set
}

// DistinctBy returns a new slice containing only the distinct elements from the source slice,
// based on a key selected by the keySelector function.
// The order of the returned elements is preserved from their first appearance.
func DistinctBy[T any, K comparable](source []T, keySelector func(T) K) []T {
	if source == nil {
		return nil
	}
	seen := make(map[K]struct{})
	result := make([]T, 0, len(source))
	for _, item := range source {
		key := keySelector(item)
		if _, ok := seen[key]; !ok {
			seen[key] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}
