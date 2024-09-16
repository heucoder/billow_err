package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
)

type filehandler func(string)

func main() {
	var path string
	flag.StringVar(&path, "path", "", "main文件所在的文件夹")
	flag.Parse()
	if path == "" {
		log.Fatalln("路径为空")
		return
	}
	log.Println("begin path:", path)
	log.Println("--------newBaseError统计开始--------------")
	walk(path, generateFileError)
	log.Println("--------newBaseError统计完成,开始生成新的code--------------")
	walk(path, countFileNewBaseError)
	log.Println("----------------code生成完成，开始输出json---------------------------")
	toJson(path)
	log.Println("--------------------json输出完成，结束----------------------------")
}

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

func walk(dirName string, f filehandler) {
	paths, _ := getDirAllFilePaths(dirName)
	ch := make(chan struct{}, 10)
	wg := sync.WaitGroup{}
	for _, path := range paths {
		ch <- struct{}{}
		wg.Add(1)
		go func(path string) {
			defer func() {
				wg.Done()
				<-ch
			}()
			fmt.Println(path)
			f(path)
		}(path)
	}
	wg.Wait()
}

func parseFile(fileName string) (*ast.File, *token.FileSet, error) {
	content, err := os.ReadFile(fileName)
	if err != nil {
		panic(err.Error())
	}
	srcCode := string(content)

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", srcCode, parser.ParseComments)
	return file, fset, err
}

func toJson(dirName string) {
	walk(dirName, getFileErrors)
	jsonData, err := json.Marshal(baseErrors)
	if err != nil {
		log.Fatalf("JSON marshaling failed: %s", err)
	}

	// 创建文件
	file, err := os.Create("errors.json")
	if err != nil {
		log.Fatalf("File creation failed: %s", err)
	}
	defer file.Close()

	// 将JSON数据写入文件
	_, err = file.Write(jsonData)
	if err != nil {
		log.Fatalf("File write failed: %s", err)
	}

	log.Println("")
}
