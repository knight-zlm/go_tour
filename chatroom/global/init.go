package global

import (
	"os"
	"path/filepath"
	"sync"
)

func init() {
	Init()
}

var (
	RootDir string
	once    = new(sync.Once)
)

func Init() {
	once.Do(func() {
		inferRootDir()
	})
}

func inferRootDir() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	var infer func(string) string
	infer = func(d string) string {
		if exists(d + "/template") {
			return d
		}

		return infer(filepath.Dir(d))
	}

	RootDir = infer(cwd)
}

func exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}
