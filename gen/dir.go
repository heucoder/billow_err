package main

import (
	"io/ioutil"
	"os"
	"strings"
)

func getDirAllFilePaths(dirname string) ([]string, error) {

	dirname = strings.TrimSuffix(dirname, string(os.PathSeparator))

	infos, err := ioutil.ReadDir(dirname)
	if err != nil {
		return nil, err
	}

	paths := make([]string, 0, len(infos))
	for _, info := range infos {
		path := dirname + string(os.PathSeparator) + info.Name()
		if info.IsDir() {
			tmp, err := getDirAllFilePaths(path)
			if err != nil {
				return nil, err
			}
			paths = append(paths, tmp...)
			continue
		}
		pathNames := strings.Split(path, ".")
		if pathNames[len(pathNames)-1] == "go" {
			paths = append(paths, path)
		}
	}
	return paths, nil
}
