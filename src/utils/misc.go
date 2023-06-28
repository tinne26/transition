package utils

import "os"

// Linear, slow, whatever.
func OsArgReceived(arg string) bool {
	for _, osArg := range os.Args {
		if osArg == arg { return true }
	}
	return false
}
