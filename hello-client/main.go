package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"example.com/learn-grpc/hello-client/pkg/interceptor"
	pb "example.com/learn-grpc/hello-client/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	// 配置ssl，"*.heliu.site"在实际开发中从浏览器中取获取，证书路径使用绝对路径
	creds, _ := credentials.NewClientTLSFromFile(
		"/root/workspace/learn-grpc/key/test.pem",
		"*.heliu.site",
	)

	var opts []grpc.DialOption
	// 不带TLS这里是grpc.WithTransportCredentials(insecure.NewCredentials())
	opts = append(opts, grpc.WithTransportCredentials(creds))
	opts = append(opts, grpc.WithPerRPCCredentials(&ClientTokenAuth{}))
	// 添加客户端拦截器
	opts = append(opts, grpc.WithUnaryInterceptor(interceptor.UnaryClientInterceptor()))
	// 添加流拦截器
	opts = append(opts, grpc.WithStreamInterceptor(interceptor.StreamClientInterceptor()))

	// 连接server端，使用ssl加密通信
	conn, err := grpc.NewClient("127.0.0.1:9090", opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()

	// 建立连接
	client := pb.NewSayHelloClient(conn)
	// 执行rpc调用(这个方法在服务器端来实现并返回结构)
	resp, err := client.SayHello(context.Background(), &pb.HelloRequest{RequestName: "gh", Age: 12})

	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	fmt.Println(resp.GetResponseMsg())

	// 客户端调用Channel方法，获取返回的流对象
	stream, err := client.Channel(context.Background())
	if err != nil {
		log.Fatalf("error creating stream: %v", err)
	}

	// 在客户端将发送和接收放到两个独立的 goroutine

	// 向服务器发送数据：
	go func() {
		for {
			req := &pb.Request{
				Value: "张三",
			}

			if err := stream.Send(req); err != nil {
				log.Fatalf("error sending message: %v", err)
				return
			}
			time.Sleep(time.Second)
		}
	}()

	// 然后再循环中接收服务端返回的数据
	for {
		reply, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("error receiving message: %v", err)
			return
		}
		fmt.Printf("Received: %s\n", reply.GetValue())
	}
}

// ClientTokenAuth 自定义Token认证
type ClientTokenAuth struct{}

func (c ClientTokenAuth) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	// uri: 是请求地址
	fmt.Printf("uri: %v\n", uri) // uri: [https://*.heliu.site/SayHello]

	fmt.Printf("%v\n", ctx)

	return map[string]string{
		"Authorization": "Bearer YOUR_ACCESS_TOKEN",
	}, nil
}

func (c ClientTokenAuth) RequireTransportSecurity() bool {
	return true
}
