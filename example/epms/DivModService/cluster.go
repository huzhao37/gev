package DivModService

import (
	"context"
	"fmt"
	"github.com/huzhao37/gev/example/epms/DivModService/protocols/gen-go/cluster"
	"net/http"
)

type ClusterHandler struct {
	log map[int]*cluster.Result_
}

func NewClusterHandler() *ClusterHandler {
	return &ClusterHandler{log: make(map[int]*cluster.Result_)}

}

func (p *ClusterHandler) DoDivMod(ctx context.Context, arg1, arg2 int64) (*cluster.Result_, error) {
	//parentContext, _, _ := thrift.ServerInterceptor(ctx, "DoDivMod")
	//fmt.Println(ctx)
	fmt.Print("DoDivMod(", arg1, arg2, ")\n")
	divRes := int64(arg1 / arg2)
	modRes := int64(arg1 % arg2)
	// 生成的用于生成自定义数据对象的函数
	res := cluster.NewResult_()
	res.Div = divRes
	res.Mod = modRes

	return res, nil
}

func (p *ClusterHandler) DoDivMod2(ctx context.Context, arg1, arg2 int64) (*cluster.Result_, error) {
	//parentContext, _, _ := thrift.ServerInterceptor(ctx, "DoDivMod")
	//fmt.Println(ctx)
	fmt.Print("DoDivMod2(", arg1, arg2, ")\n")
	divRes := int64(arg1 % arg2)
	modRes := int64(arg1 / arg2)
	// 生成的用于生成自定义数据对象的函数
	res := cluster.NewResult_()
	res.Div = divRes
	res.Mod = modRes

	return res, nil
}

func (p *ClusterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("cluster divmod \n")
	p.DoDivMod(r.Context(), 100, 1)
}
