package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"job-detail-api-go/internal"
	pb "job-detail-api-go/pkg/grpc"
)

func main() {
	// AWS設定の読み込み
	// LocalStack対応: AWS_ENDPOINT_URLがある場合はダミー認証情報を設定
	var opts []func(*config.LoadOptions) error
	if os.Getenv("AWS_ENDPOINT_URL") != "" {
		opts = append(opts, config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("dummy", "dummy", "")))
		if os.Getenv("AWS_REGION") == "" {
			opts = append(opts, config.WithRegion("ap-northeast-1"))
		}
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(), opts...)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	// DynamoDBクライアントの作成
	// LocalStack対応: AWS_ENDPOINT_URL環境変数があればBaseEndpointとして設定
	svc := dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		if endpoint := os.Getenv("AWS_ENDPOINT_URL"); endpoint != "" {
			o.BaseEndpoint = aws.String(endpoint)
		}
	})

	tableName := os.Getenv("TABLE_NAME")
	if tableName == "" {
		tableName = "Jobs"
	}

	h := &internal.Handler{
		Client:    svc,
		TableName: tableName,
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterJobServiceServer(s, h)
	reflection.Register(s)

	log.Printf("gRPC server starting on port %s...", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
