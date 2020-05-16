package DivModService

import (
	"context"
	"fmt"
	"github.com/huzhao37/gev/epms/DivModService/protocols/gen-go/cluster"
	"net/http"
)

type ClusterHandler struct {
	log map[int]*cluster.Result_
}

func NewClusterHandler() *ClusterHandler {
	return &ClusterHandler{log: make(map[int]*cluster.Result_)}

}

func (p *ClusterHandler) DoDivMod(ctx context.Context, args *cluster.DivModDoDivModArgs, res *cluster.DivModDoDivModResult) error {
	//parentContext, _, _ := thrift.ServerInterceptor(ctx, "DoDivMod")
	//fmt.Println(ctx)
	fmt.Print("DoDivMod(", args.Arg1, args.Arg2, ")\n")
	divRes := int64(args.Arg1 / args.Arg2)
	modRes := int64(args.Arg1 % args.Arg2)
	// 生成的用于生成自定义数据对象的函数
	//res = &cluster.DivModDoDivModResult{Success:cluster.NewResult_()}
	res.Success.Div = divRes
	res.Success.Mod = modRes

	return nil
}

func (p *ClusterHandler) DoDivMod2(ctx context.Context, args *cluster.DivMod2DoDivMod2Args, res *cluster.DivMod2DoDivMod2Result) error {
	//parentContext, _, _ := thrift.ServerInterceptor(ctx, "DoDivMod")
	//fmt.Println(ctx)
	fmt.Print("DoDivMod2(", args.Arg1, args.Arg2, ")\n")
	divRes := int64(args.Arg1 % args.Arg2)
	modRes := int64(args.Arg1 / args.Arg2)
	// 生成的用于生成自定义数据对象的函数
	res = &cluster.DivMod2DoDivMod2Result{Success: cluster.NewResult_()}
	res.Success.Div = divRes
	res.Success.Mod = modRes

	return nil
}

func (p *ClusterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("cluster divmod \n")
	//p.DoDivMod(r.Context(), 100, 1)
}
