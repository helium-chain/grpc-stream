package main

import (
	"context"
	"fmt"
	"log"

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

	// 连接server端，使用ssl加密通信
	conn, err := grpc.NewClient("127.0.0.1:9090", opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()

	// 建立连接
	client := pb.NewSayHelloClient(conn)
	// 执行rpc调用(这个方法在服务器端来实现并返回结构)
	resp, err := client.SayHello(context.Background(), &pb.HelloRequest{RequestName: "gh"})

	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	fmt.Println(resp.GetResponseMsg())
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
