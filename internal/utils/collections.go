package utils

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
