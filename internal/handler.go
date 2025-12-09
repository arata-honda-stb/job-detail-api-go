package internal

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type JobRequest struct {
	JobIDs []string `json:"jobIds"`
}

type JobResponse struct {
	Jobs []Job `json:"job"`
}

type Handler struct {
	Client    *dynamodb.Client
	TableName string
}

func (h *Handler) GetJobs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req JobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if len(req.JobIDs) == 0 {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JobResponse{Jobs: []Job{}})
		return
	}

	// DynamoDB BatchGetItem Keys作成
	// 重複を除去するなどしたほうが良いが、簡易実装とする
	keys := make([]map[string]types.AttributeValue, 0, len(req.JobIDs))
	seen := make(map[string]bool)

	for _, id := range req.JobIDs {
		if _, exists := seen[id]; !exists {
			seen[id] = true
			keys = append(keys, map[string]types.AttributeValue{
				"jobId": &types.AttributeValueMemberS{Value: id},
			})
		}
	}

	// BatchGetItemの制限（最大100件）があるため、実際の運用では分割が必要だが
	// 簡易実装のため、ここではそのままリクエストする
	input := &dynamodb.BatchGetItemInput{
		RequestItems: map[string]types.KeysAndAttributes{
			h.TableName: {
				Keys: keys,
			},
		},
	}

	result, err := h.Client.BatchGetItem(context.TODO(), input)
	if err != nil {
		log.Printf("failed to get items: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var jobs []Job
	if items, ok := result.Responses[h.TableName]; ok {
		if err := attributevalue.UnmarshalListOfMaps(items, &jobs); err != nil {
			log.Printf("failed to unmarshal items: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

	if jobs == nil {
		jobs = []Job{}
	}

	resp := JobResponse{
		Jobs: jobs,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}
