package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"strconv"

	"example.com/learn-grpc/hello-server/pkg/interceptor"
	pb "example.com/learn-grpc/hello-server/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

func main() {
	// 绝对地址
	creds, err1 := credentials.NewServerTLSFromFile(
		"/root/workspace/learn-grpc/key/test.pem",
		"/root/workspace/learn-grpc/key/test.key",
	)

	if err1 != nil {
		fmt.Printf("证书错误：%v", err1)
		return
	}

	// 开启端口
	listen, _ := net.Listen("tcp", ":9090")
	// 创建grpc服务
	grpcServer := grpc.NewServer(grpc.Creds(creds),
		grpc.UnaryInterceptor(interceptor.UnaryServerInterceptor()),
		grpc.StreamInterceptor(interceptor.StreamServerInterceptor()))
	// 在grpc服务端中注册我们自己编写的服务
	pb.RegisterSayHelloServer(grpcServer, &server{})

	// 启动服务
	err := grpcServer.Serve(listen)
	if err != nil {
		fmt.Println(err)
		return
	}
}

var _ pb.SayHelloServer = (*server)(nil)

type server struct {
	pb.UnimplementedSayHelloServer
}

func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	// 获取元数据信息，这里的逻辑可以搬去拦截器中
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("metadata not found")
	}
	// metadata: map[:authority:[*.heliu.site] authorization:[Bearer YOUR_ACCESS_TOKEN] client-os:[linux] content-type:[application/grpc] user-agent:[grpc-go/1.66.0]]
	fmt.Println("metadata:", md)

	userId, ok := md["authorization"]
	if !ok {
		return nil, fmt.Errorf("metadata not found 2")
	}
	fmt.Printf("Authorization: %s\n", userId[0])

	if userId[0] != "Bearer YOUR_ACCESS_TOKEN" {
		return nil, fmt.Errorf("metadata not found 3")
	}

	// 正常的业务处理
	fmt.Printf("Received: %s\n", req.RequestName)
	return &pb.HelloResponse{ResponseMsg: "hello," + req.RequestName + " age:" + strconv.FormatInt(req.Age, 10)}, nil
}

func (s *server) Channel(stream pb.SayHello_ChannelServer) error {
	// 服务端在循环中接收客户端发来的数据
	for {
		req, err := stream.Recv()
		if err != nil {
			// 如果遇到io.EOF标识客户端流关闭
			if err == io.EOF {
				return nil
			}
			// 其他的io.Error 类型都返回
			return err
		}
		fmt.Printf("Received: %s\n", req.GetValue())
		// 假设需要返回处理结果
		err = stream.Send(&pb.Response{Value: "hello," + req.GetValue()})
		if err != nil {
			// 服务器发送异常，函数退出，服务端流关闭
			return err
		}
	}
}
