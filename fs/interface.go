package fs

import "os"

// FileSystem 文件系统
type FileSystem interface {
	Open(f string) (*os.File, error)
}
