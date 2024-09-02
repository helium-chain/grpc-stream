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

	// 连接server端，使用ssl加密通信
	conn, err := grpc.NewClient("127.0.0.1:9090", grpc.WithTransportCredentials(creds))
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
