package grpclib

import (
	"context"
	"golang-im/pkg/gerrors"
	"golang-im/pkg/logger"

	"google.golang.org/grpc/metadata"
)

//本文件 为ctx.go 的副本 主要差异在返回值类型不一样

// ContextWithRequstIdStr 向http请求头header中加个 请求id, 原理就是 用1个context携带一个map并发送出去
func ContextWithRequstIdStr(ctx context.Context, requestId string) context.Context {
	return metadata.NewOutgoingContext(ctx, metadata.Pairs(CtxRequestId, requestId))
}

// GetCtxRequestId 获取ctx的app_id
func GetCtxRequestIdStr(ctx context.Context) string {
	var dst string
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return dst
	}

	requstIds, ok := md[CtxRequestId]
	if !ok && len(requstIds) == 0 {
		return dst
	}
	return requstIds[0]
}

// GetCtxDataStr 获取ctx的用户数据，依次返回user_id,device_id
func GetCtxDataStr(ctx context.Context) (string, string, error) {
	var (
		userId   string
		deviceId string
	)

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		logger.Sugar.Infow("---->处理鉴权逻辑失败 ctx的用户数据 1")
		return userId, deviceId, gerrors.ErrUnauthorized
	}

	userIdStrs, ok := md[CtxUserId]
	if !ok && len(userIdStrs) == 0 {
		logger.Sugar.Infow("---->处理鉴权逻辑失败 ctx的用户数据 2")
		return userId, deviceId, gerrors.ErrUnauthorized
	}
	userId = userIdStrs[0]

	deviceIdStrs, ok := md[CtxDeviceId]
	if !ok && len(deviceIdStrs) == 0 {
		logger.Sugar.Infow("---->处理鉴权逻辑失败 ctx的用户数据 4")
		return userId, deviceId, gerrors.ErrUnauthorized
	}
	deviceId = deviceIdStrs[0]

	return userId, deviceId, nil
}
