package rpc

import (
	"context"
	"golang-im/pkg/gerrors"

	"google.golang.org/grpc"
)

//拦截器在作用于每一个 RPC 调用，通常用来做日志，认证，metric 等等
func interceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	//1.预处理 可以通过参数获取 context, method 名称,发送的请求, CallOption
	//2.调用(invoker)RPC 方法
	err := invoker(ctx, method, req, reply, cc, opts...)
	//3.调用后处理
	return gerrors.WrapRPCError(err)
}
