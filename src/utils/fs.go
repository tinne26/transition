package utils

import "io/fs"
import "image"
import "image/png"

import "github.com/hajimehoshi/ebiten/v2"

// The path must use forward slashes /, as the fs.FS interface
// operates like a unix filesystem.
func LoadFsEbiImage(filesys fs.FS, path string) (*ebiten.Image, error) {
	img, err := LoadFsImage(filesys, path)
	if err != nil { return nil, err }
	return ebiten.NewImageFromImage(img), nil
}

func LoadFsImage(filesys fs.FS, path string) (image.Image, error) {
	file, err := filesys.Open(path)
	if err != nil { return nil, err }
	defer file.Close()
	return png.Decode(file)
}
