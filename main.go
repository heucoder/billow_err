package main

import (
	"fmt"

	"github.com/heucoder/billow_err/billow_err"
)

//go:generate go get github.com/heucoder/billow_err/gen
//go:generate go run github.com/heucoder/billow_err/gen -path ./
func main() {
	err := FuncL1(10001)
	if err != nil {
		fmt.Printf("err:%+v\n", err)
	}
}

func FuncL1(a int64) error {
	if a > 100 {
		return FuncL2(a)
	} else {
		return FuncL3(a)
	}
}

func FuncL2(a int64) error {
	err := FuncCommon(a)
	if err != nil {
		return err
	}
	return nil
}

func FuncL3(a int64) error {
	err := FuncCommon(a)
	if err != nil {
		return err
	}
	return nil
}

func FuncCommon(a int64) error {
	if a > 10000 {
		return billow_err.NewBaseError(1, "a>10000")
	}
	fmt.Printf("FuncCommon req:%v\n", a)
	return nil
}
