package interceptor

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func UnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		// method: /SayHello/SayHello
		// req: 接口装载的 *server.HelloRequest
		// reply: 接口返回的 *server.HelloReply
		// cc: grpc.ClientConn
		// invoker: grpc.UnaryInvoker
		// opts: []grpc.CallOption

		fmt.Printf("req: %v\n", req)

		// 获取 gRPC 方法名称
		// method: /SayHello/SayHello
		fmt.Printf("opts: %s\n", opts)

		// 获取 gRPC 方法的请求参数
		// req: 接口装载的 *server.HelloRequest
		fmt.Printf("req: %v\n", req)
		fmt.Printf("type: %T\n", req)

		// 获取 gRPC 方法的返回参数
		// reply: 接口返回的 *server.HelloReply
		fmt.Printf("reply: %v\n", reply)
		fmt.Printf("type: %T\n", reply)

		// 获取 gRPC 方法的 gRPC 客户端连接
		// cc: grpc.Client

		// 1) 预处理阶段
		start := time.Now()

		cos := runtime.GOOS // 获取操作系统
		// 将操作系统信息附加到元数据传出请求
		ctx = metadata.AppendToOutgoingContext(ctx, "client-os", cos)

		// 2) 调用 gRPC 方法
		err := invoker(ctx, method, req, reply, cc, opts...)

		// 3) 后处理阶段
		end := time.Now()
		//  RPC: /SayHello/SayHello, client-os: linux, req: requestName:"gh", start-time: 2024-09-03 21:25:51.715107181 +0800 CST m=+0.000534541, end-time: 2024-09-03 21:25:51.717004906 +0800 CST m=+0.002432245, err: <nil>
		log.Printf("RPC: %s, client-os: %s, req: %v, start-time: %s, end-time: %s, err: %v\n", method, cos, req, start, end, err)
		// *server.HelloRequest
		fmt.Printf("req: %v\n", req)
		fmt.Printf("type: %T\n", req)

		return err
	}
}

func StreamClientInterceptor() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		// 1) 预处理
		log.Printf("method: %v\n", method) // /SayHello/Channel

		// 因为SendMsg和RecvMsg方法是ClientStream接口内方法，我们需要先调用streamer函数获取到ClientStream
		// 再对它进行封装，实现自己的 SendMsg 和 RecvMsg 方法。
		stream, err := streamer(ctx, desc, cc, method, opts...)
		return newStreamClient(stream), err
	}
}

type streamClient struct {
	grpc.ClientStream
}

func newStreamClient(c grpc.ClientStream) grpc.ClientStream {
	return &streamClient{c}
}

func (s *streamClient) SendMsg(m interface{}) error {
	// 2) 发送前，我们可以再这里对发送的消息处理
	fmt.Printf("SendMsg: %v\n m: %T\n", m, m) // SendMsg: value:"张三"\n m: *service.Request
	return s.ClientStream.SendMsg(m)
}

func (s *streamClient) RecvMsg(m interface{}) error {
	// 3) 在这里，我们可以对接收到的消息进行处理(发送前)
	fmt.Printf("RecvMsg: %v\n m:%T\n", m, m) // RecvMsg: m:*service.Response
	return s.ClientStream.RecvMsg(m)
	// 发送后
}
