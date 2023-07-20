package utils

import "os"

// Linear, slow, whatever.
func OsArgReceived(arg string) bool {
	for _, osArg := range os.Args {
		if osArg == arg { return true }
	}
	return false
}

func FastFill[T any](buffer []T, value T) {
	if len(buffer) <= 24 { // no-copy case
		for i, _ := range buffer {
			buffer[i] = value
		}
	} else { // copy case
		for i, _ := range buffer[:16] {
			buffer[i] = value
		}
		for i := 16; i < len(buffer); i *= 2 {
			copy(buffer[i:], buffer[:i])
		}
	}
}

func IterDelete[T any](slice []T, deleteFn func(T) bool) []T {
	deleteCount := 0
	i := 0
	for i < len(slice) - deleteCount {
		if deleteFn(slice[i]) {
			last := len(slice) - deleteCount - 1
			slice[i], slice[last] = slice[last], slice[i]
			deleteCount += 1
		} else {
			i += 1
		}
	}
	return slice[0 : len(slice) - deleteCount]
}
