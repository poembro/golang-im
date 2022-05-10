package gerrors

import (
	"context"
	"fmt"
	"golang-im/pkg/logger"
	"golang-im/pkg/util"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"runtime"
	"strings"

	"github.com/golang/protobuf/ptypes/any"
	spb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const name = "im"

const TypeUrlStack = "type_url_stack"

func WrapError(err error) error {
	if err == nil {
		return nil
	}

	s := &spb.Status{
		Code:    int32(codes.Unknown),
		Message: err.Error(),
		Details: []*any.Any{
			{
				TypeUrl: TypeUrlStack,
				Value:   util.Str2bytes(stack()),
			},
		},
	}
	return status.FromProto(s).Err()
}

func WrapRPCError(err error) error {
	if err == nil {
		return nil
	}
	e, _ := status.FromError(err)
	s := &spb.Status{
		Code:    int32(e.Code()),
		Message: e.Message(),
		Details: []*any.Any{
			{
				TypeUrl: TypeUrlStack,
				Value:   util.Str2bytes(GetErrorStack(e) + " --grpc调用栈-- \n" + stack()),
			},
		},
	}
	return status.FromProto(s).Err()
}

func GetErrorStack(s *status.Status) string {
	pbs := s.Proto()
	for i := range pbs.Details {
		if pbs.Details[i].TypeUrl == TypeUrlStack {
			return util.B2S(pbs.Details[i].Value)
		}
	}
	return ""
}

// Stack 获取堆栈信息
func stack() string {
	var pc = make([]uintptr, 20)
	n := runtime.Callers(3, pc)

	var build strings.Builder
	for i := 0; i < n; i++ {
		f := runtime.FuncForPC(pc[i] - 1)
		file, line := f.FileLine(pc[i] - 1)
		n := strings.Index(file, name)
		if n != -1 {
			s := fmt.Sprintf(" %s:%d \n", file[n:], line)
			build.WriteString(s)
		}
	}
	return build.String()
}

// LogPanic 用于GRPC过滤器 函数出错时打印日志
func LogPanic(serverName string, ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, err *error) {
	p := recover()
	if p != nil {
		logger.Logger.Error(serverName+" panic", zap.Any("info", info), zap.Any("ctx", ctx), zap.Any("req", req),
			zap.Any("panic", p), zap.String("stack", util.GetStackInfo()))
		*err = ErrUnknown
	}
}
