package f

import (
	"fmt"

	"github.com/heucoder/billow_err/billow_err"
)

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
		return billow_err.NewBaseError(14, "a>10000")
	}
	if a < 10 {
		return billow_err.NewBaseError(15, "a<10")
	}
	fmt.Printf("FuncCommon req:%v\n", a)
	return nil
}
