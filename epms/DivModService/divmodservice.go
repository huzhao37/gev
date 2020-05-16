/**
 * @Author: hiram
 * @Date: 2020/5/9 14:58
 */
package DivModService

import (
	"context"
	"fmt"
)

type DivModResult struct {
	Div int64
	Mod int64
}
type DivModHandler struct {
	log map[int]*DivModResult
}

func NewDivModHandler() *DivModHandler {
	return &DivModHandler{log: make(map[int]*DivModResult)}

}

func (p *DivModHandler) DoDivMod(ctx context.Context, arg1, arg2 int64) (*DivModResult, error) {
	fmt.Print("DoDivMod(", arg1, arg2, ")\n")
	divRes := int64(arg1 / arg2)
	modRes := int64(arg1 % arg2)
	// 生成的用于生成自定义数据对象的函数
	res := &DivModResult{}
	res.Div = divRes
	res.Mod = modRes

	return res, nil
}

func (p *DivModHandler) DoDivMod2(ctx context.Context, arg1, arg2 int64) (*DivModResult, error) {
	fmt.Print("DoDivMod2(", arg1, arg2, ")\n")
	divRes := int64(arg1 % arg2)
	modRes := int64(arg1 / arg2)
	// 生成的用于生成自定义数据对象的函数
	res := &DivModResult{}
	res.Div = divRes
	res.Mod = modRes

	return res, nil
}
