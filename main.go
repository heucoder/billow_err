package main

import (
	"fmt"

	"github.com/heucoder/billow_err/billow_err"
	"github.com/heucoder/billow_err/f"
)

//go:generate go get github.com/heucoder/billow_err/gen
//go:generate go run github.com/heucoder/billow_err/gen -path ./
func main() {
	err := f.FuncL1(10001)
	if err != nil {
		fmt.Printf("err:%+v\n", err)
	}
	billow_err.NewBaseError(12, "test 1")
	billow_err.NewBaseError(13, "test 2")
}
