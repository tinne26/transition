//go:build windows

package utils

import "os"

func OnWindowsPreferOpenGL() error {
	// allow directX if passed as program flag
	for _, arg := range os.Args {
		if arg == "--directX" || arg == "--directx" { return nil }
	}

	// set openGL as the graphics backend otherwise
	return os.Setenv("EBITENGINE_GRAPHICS_LIBRARY", "opengl")
}
