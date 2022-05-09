package interceptor

import (
	"context"
	"golang-im/pkg/gerrors"
	"golang-im/pkg/grpclib"
	"golang-im/pkg/logger"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// NewInterceptor 生成GRPC过滤器
func NewInterceptor(name string, urlWhitelist map[string]int) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer gerrors.LogPanic(name, ctx, req, info, &err)

		md, _ := metadata.FromIncomingContext(ctx)
		resp, err = handleWithAuth(ctx, req, info, handler, urlWhitelist)

		logger.Logger.Debug(name, zap.String("desc", "grpc gateway"), zap.Any("method", info.FullMethod), zap.Any("md", md), zap.Any("req", req),
			zap.Any("resp", resp), zap.Error(err))

		s, _ := status.FromError(err)
		if s.Code() != 0 && s.Code() < 1000 {
			md, _ := metadata.FromIncomingContext(ctx)
			logger.Logger.Error(name, zap.String("method", info.FullMethod), zap.Any("md", md), zap.Any("req", req),
				zap.Any("resp", resp), zap.Error(err), zap.String("stack", gerrors.GetErrorStack(s)))
		}
		return
	}
}

// handleWithAuth 处理鉴权逻辑
func handleWithAuth(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler, urlWhitelist map[string]int) (interface{}, error) {
	if _, ok := urlWhitelist[info.FullMethod]; !ok {
		requestId := grpclib.GetCtxRequestIdStr(ctx)

		userId, deviceId, err := grpclib.GetCtxDataStr(ctx)
		if err != nil {
			logger.Sugar.Infow("---->处理鉴权逻辑失败 grpc头信息没有user_id device_id")
			return nil, err
		}
		token, err := grpclib.GetCtxToken(ctx)
		if err != nil {
			logger.Sugar.Infow("---->处理鉴权逻辑失败 grpc头信息没有token")
			return nil, err
		}

		// TODO 去处理业务接口验证token
		//_, err = rpc.BusinessIntClient.Auth(ctx, &pb.AuthReq{
		//    UserId:   userId,
		//    DeviceId: deviceId,
		//    Token:    token,
		//})

		logger.Logger.Debug("handleWithAuth",
			zap.String("desc", "grpc gateway"),
			zap.String("userId", userId),
			zap.String("deviceId", deviceId),
			zap.String("token", token),
			zap.String("requestId", requestId))
	}
	return handler(ctx, req)
}
