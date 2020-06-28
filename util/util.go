package util

import (
	"os"
	"path"
)

/*
:)
*/
func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

func RelativePath(relpath string) string {
	return path.Join(os.Getenv("root_file_path"), relpath)
}
