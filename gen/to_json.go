package main

import (
	"go/ast"
	"strconv"
	"sync"

	"github.com/heucoder/billow_err/billow_err"
)

var baseErrors = []*billow_err.BaseErr{}
var getMu sync.Mutex

func getFileErrors(fileName string) {
	file, _, err := parseFile(fileName)
	if err != nil {
		panic(err.Error())
	}

	ast.Inspect(file, func(node ast.Node) bool {
		if callExpr, ok := (node).(*ast.CallExpr); ok {
			if selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok && selectorExpr.Sel.Name == "NewBaseError" {
				if len(callExpr.Args) != 2 {
					return false
				}
				var errorNum int32
				var errMsg string
				if val, ok := callExpr.Args[0].(*ast.BasicLit); ok {
					n, _ := strconv.ParseInt(val.Value, 10, 32)
					errorNum = int32(n)
				}
				if val, ok := callExpr.Args[1].(*ast.BasicLit); ok {
					errMsg = val.Value
				}
				if errorNum != 0 {
					getMu.Lock()
					baseErrors = append(baseErrors, &billow_err.BaseErr{
						Code: errorNum,
						Msg:  errMsg,
					})
					getMu.Unlock()
				}
			}
		}
		return true
	})
}
