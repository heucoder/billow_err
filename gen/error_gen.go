package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"os"
	"strconv"
	"sync"
)

var mu sync.Mutex

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

func generateError(fileName string) {
	file, fset, err := parseFile(fileName)
	if err != nil {
		panic(err.Error())
	}
	ast.Inspect(file, func(node ast.Node) bool {
		if callExpr, ok := (node).(*ast.CallExpr); ok {
			if selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok && selectorExpr.Sel.Name == "TodoBaseErr" {
				if errorMsg, ok := callExpr.Args[0].(*ast.BasicLit); ok {
					newCallExpr := createNewBaseErrorExpr(errorMsg.Value)
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

func createNewBaseErrorExpr(msg string) *ast.CallExpr {

	BaseErrorFunc := &ast.SelectorExpr{
		X:   ast.NewIdent("billow_err"),
		Sel: ast.NewIdent("NewBaseError"),
	}

	errorN := errorNumGenerate.Add(1)
	codeLit := &ast.BasicLit{
		Kind:  token.INT,
		Value: fmt.Sprintf("%v", errorN),
	}
	msgLit := &ast.BasicLit{
		Kind:  token.STRING,
		Value: msg,
	}

	callExpr := &ast.CallExpr{
		Fun:  BaseErrorFunc,
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

func countFileNewBaseError(fileName string) {
	file, _, err := parseFile(fileName)
	if err != nil {
		panic(err.Error())
	}

	ast.Inspect(file, func(node ast.Node) bool {
		if callExpr, ok := (node).(*ast.CallExpr); ok {
			if selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok && selectorExpr.Sel.Name == "NewBaseError" {
				if val, ok := callExpr.Args[0].(*ast.BasicLit); ok {
					errorNum, _ := strconv.ParseInt(val.Value, 10, 32)

					mu.Lock()
					if errorNumGenerate.Load() < int32(errorNum) {
						errorNumGenerate.Store(int32(errorNum))
					}
					mu.Unlock()
				}
			}
		}
		return true
	})
}
