package util

import (
	"os"
	"path"
)

func RelativePath(relpath string) string {
	return path.Join(os.Getenv("root_file_path"), relpath)
}

// thread safe wrappers for use when sharing data
type SafeMap struct {
}

type SafeSlice struct {
}
