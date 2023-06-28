package utils

import "io/fs"
import "image/png"

import "github.com/hajimehoshi/ebiten/v2"

// The path must use forward slashes /, as the fs.FS interface
// operates like a unix filesystem.
func LoadFsImage(filesys fs.FS, path string) (*ebiten.Image, error) {
	file, err := filesys.Open(path)
	if err != nil { return nil, err }
	defer file.Close()
	img, err := png.Decode(file)
	if err != nil { return nil, err }
	return ebiten.NewImageFromImage(img), nil
}
