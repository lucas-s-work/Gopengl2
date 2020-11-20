package util

import (
	"os"
	"path"
)

func RelativePath(relpath string) string {
	return path.Join(os.Getenv("root_file_path"), relpath)
}
