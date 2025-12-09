package internal

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	pb "job-detail-api-go/pkg/grpc"
)

type Handler struct {
	pb.UnimplementedJobServiceServer
	Client    *dynamodb.Client
	TableName string
}

func (h *Handler) GetJobs(ctx context.Context, req *pb.GetJobsRequest) (*pb.GetJobsResponse, error) {
	if len(req.JobIds) == 0 {
		return &pb.GetJobsResponse{}, nil
	}

	// DynamoDB BatchGetItem Keys作成
	keys := make([]map[string]types.AttributeValue, 0, len(req.JobIds))
	seen := make(map[string]bool)

	for _, id := range req.JobIds {
		if _, exists := seen[id]; !exists {
			seen[id] = true
			keys = append(keys, map[string]types.AttributeValue{
				"jobId": &types.AttributeValueMemberS{Value: id},
			})
		}
	}

	input := &dynamodb.BatchGetItemInput{
		RequestItems: map[string]types.KeysAndAttributes{
			h.TableName: {
				Keys: keys,
			},
		},
	}

	result, err := h.Client.BatchGetItem(ctx, input)
	if err != nil {
		log.Printf("failed to get items: %v", err)
		return nil, err
	}

	var jobs []Job
	if items, ok := result.Responses[h.TableName]; ok {
		if err := attributevalue.UnmarshalListOfMaps(items, &jobs); err != nil {
			log.Printf("failed to unmarshal items: %v", err)
			return nil, err
		}
	}

	// マッピング処理 (DB Model -> API Model)
	pbJobs := make([]*pb.Job, 0, len(jobs))
	for _, j := range jobs {
		pbJobs = append(pbJobs, &pb.Job{
			JobId:      j.JobID,
			JobTitle:   j.JobTitle,
			JobContent: j.JobContent,
		})
	}

	return &pb.GetJobsResponse{
		Jobs: pbJobs,
	}, nil
}
