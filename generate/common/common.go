package common

const (
	DefaultAzcore = "1.9.0"
)

// slice convert to map
func SliceConvertMap[T comparable](s []T) map[T]bool {
	m := make(map[T]bool)
	for _, v := range s {
		m[v] = true
	}

	return m
}
