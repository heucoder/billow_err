package main

import (
	"flag"
	"fmt"
	"log"
	"sync"
)

func main() {
	var path string
	flag.StringVar(&path, "path", "", "main文件所在的文件夹")
	flag.Parse()
	if path == "" {
		log.Fatalln("路径为空")
		return
	}
	fmt.Println("begin path:", path)
	countBillowsError(path)
	generateBillowsError(path)
	log.Println("--------------------成功----------------------------")
}

func generateBillowsError(dirname string) {
	paths, _ := getDirAllFilePaths(dirname)
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
			generateError(path)
		}(path)
	}
	wg.Wait()
}

func countBillowsError(dirname string) {
	paths, _ := getDirAllFilePaths(dirname)
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
			countFileNewBaseError(path)
		}(path)
	}
	wg.Wait()
	fmt.Println(errorNumGenerate.Load())
}
