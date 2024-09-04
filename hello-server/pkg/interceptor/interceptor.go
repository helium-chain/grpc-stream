package interceptor

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// req：请求参数类型 *service.HelloRequest
		// info:

		fmt.Printf("req: %T\n", req)

		// 1) 预处理

		start := time.Now()

		// 从传入上下文获取元数据
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, fmt.Errorf("上下文中没有元数据")
		}

		// 检索客户端操作系统，如果它不存在说明此值为空
		os := md.Get("client-os")
		// 获取客户端IP地址
		ip, err := getClientIp(ctx)
		if err != nil {
			return nil, err
		}

		// 2) 调用 RPC 方法
		m, err := handler(ctx, req)

		// 3) 后处理
		end := time.Now()

		fmt.Printf("RPC: %s, client-os: '%v' and  IP: '%v', req: %v, start-time: %s, end-time: %s, err: %v",
			info.FullMethod, os, ip, req, start.Format(time.RFC3339), end.Format(time.RFC3339), err)

		return m, err
	}
}

func getClientIp(ctx context.Context) (string, error) {
	p, ok := peer.FromContext(ctx)

	if !ok {
		return "", fmt.Errorf("没有获取到IP")
	}

	return p.Addr.String(), nil
}
