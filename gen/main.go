package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	var path string
	flag.StringVar(&path, "path", "", "main文件所在的文件夹")
	flag.Parse()
	if path == "" {
		log.Fatalln("路径为空")
		return
	}
	fmt.Println("path:", path)
	generateBillowsError(path)
	log.Println("成功")
}

var errorN int

func generateError3(fileName, srcCode string) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", srcCode, parser.ParseComments)
	if err != nil {
		panic(err.Error())
	}

	ast.Inspect(file, func(node ast.Node) bool {
		if callExpr, ok := (node).(*ast.CallExpr); ok {
			if selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok && selectorExpr.Sel.Name == "TodoBaseErr" {
				if errorMsg, ok := callExpr.Args[0].(*ast.BasicLit); ok {
					newCallExpr := CreateNewBaseErrorExpr(errorMsg.Value)
					callExpr.Fun = newCallExpr.Fun
					callExpr.Args = newCallExpr.Args
				}
			}
		}
		return true
	})

	if err := saveASTToFile(fileName, fset, file); err != nil {
		log.Panicln(err.Error())
	}
}

func CreateNewBaseErrorExpr(msg string) *ast.CallExpr {

	printFunc := &ast.SelectorExpr{
		X:   ast.NewIdent("billow_err"),
		Sel: ast.NewIdent("NewBaseError"),
	}

	errorN += 1
	codeLit := &ast.BasicLit{
		Kind:  token.INT,
		Value: fmt.Sprintf("%v", errorN),
	}
	msgLit := &ast.BasicLit{
		Kind:  token.STRING,
		Value: msg,
	}

	callExpr := &ast.CallExpr{
		Fun:  printFunc,
		Args: []ast.Expr{codeLit, msgLit},
	}
	return callExpr
}

func saveASTToFile(fileName string, fset *token.FileSet, node *ast.File) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	var buf bytes.Buffer
	err = format.Node(&buf, fset, node)
	if err != nil {
		return err
	}

	_, err = file.WriteString(buf.String())
	return err
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

func generateBillowsError(dirname string) {
	paths, _ := getDirAllFilePaths(dirname)
	for _, path := range paths {
		fmt.Println(path)
		content, err := os.ReadFile(path)
		if err != nil {
			log.Fatal(err)
		}

		srcCode := string(content)
		generateError3(path, srcCode)
	}
}
