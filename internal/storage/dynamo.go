package storage

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"int-status/internal"
	"sort"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamoDBStorage struct {
	client *dynamodb.Client
	table  string
}

func NewDynamoDBStorage(region string, table string) (*DynamoDBStorage, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return nil, err
	}

	client := dynamodb.NewFromConfig(cfg)
	return &DynamoDBStorage{
		client: client,
		table:  table,
	}, nil
}

func (s *DynamoDBStorage) UpdateHistory(statuses []internal.Status) error {
	const maxBatchSize = 25 // DynamoDB BatchWriteItem의 최대 크기
	var writeRequests []types.WriteRequest

	for _, status := range statuses {
		data, err := s.toDynamoDBData(status)
		if err != nil {
			return fmt.Errorf("failed to marshal status data: %v", err)
		}

		writeRequests = append(writeRequests, types.WriteRequest{
			PutRequest: &types.PutRequest{
				Item: data,
			},
		})

		if len(writeRequests) == maxBatchSize {
			if err := s.batchWrite(writeRequests); err != nil {
				return fmt.Errorf("failed to execute batch write: %v", err)
			}
			writeRequests = nil
		}
	}

	if len(writeRequests) > 0 {
		if err := s.batchWrite(writeRequests); err != nil {
			return fmt.Errorf("failed to execute final batch write: %v", err)
		}
	}

	return nil
}

func (s *DynamoDBStorage) GetDailyHistory(service string) ([]internal.Status, error) {
	today := time.Now().Format("2006-01-02")

	result, err := s.client.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:              aws.String(s.table),
		KeyConditionExpression: aws.String("service = :service AND #timestamp BETWEEN :start AND :end"),
		ExpressionAttributeNames: map[string]string{
			"#timestamp": "timestamp",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":service": &types.AttributeValueMemberS{Value: service},
			":start":   &types.AttributeValueMemberS{Value: today + "T00:00:00+09:00"},
			":end":     &types.AttributeValueMemberS{Value: today + "T23:59:59+09:00"},
		},
		ScanIndexForward: aws.Bool(false),
		Limit:            aws.Int32(10),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to query DynamoDB: %v", err)
	}

	var statuses []internal.Status
	for _, item := range result.Items {
		var status internal.Status
		err := attributevalue.UnmarshalMap(item, &status)
		if err != nil {
			return nil, err
		}
		statuses = append(statuses, status)
	}

	if len(statuses) == 0 {
		return nil, fmt.Errorf("no status data found. date = %s", today)
	}

	for i, j := 0, len(statuses)-1; i < j; i, j = i+1, j-1 {
		statuses[i], statuses[j] = statuses[j], statuses[i]
	}

	return statuses, nil
}

func (s *DynamoDBStorage) GetDailyIncidents(service string) ([]internal.Incident, error) {
	today := time.Now().Format("2006-01-02")

	result, err := s.client.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:              aws.String(s.table),
		KeyConditionExpression: aws.String("service = :service AND #timestamp BETWEEN :start AND :end"),
		FilterExpression:       aws.String("#status = :status"),
		ExpressionAttributeNames: map[string]string{
			"#timestamp": "timestamp",
			"#status":    "status",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":service": &types.AttributeValueMemberS{Value: service},
			":status":  &types.AttributeValueMemberS{Value: "DOWN"},
			":start":   &types.AttributeValueMemberS{Value: today + "T00:00:00+09:00"},
			":end":     &types.AttributeValueMemberS{Value: today + "T23:59:59+09:00"},
		},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to query service %s: %v", service, err)
	}

	var statuses []internal.Status
	for _, item := range result.Items {
		var status internal.Status
		if err := attributevalue.UnmarshalMap(item, &status); err != nil {
			return nil, fmt.Errorf("failed to unmarshal status: %v", err)
		}
		statuses = append(statuses, status)
	}

	if len(statuses) == 0 {
		return nil, nil
	}

	sort.Slice(statuses, func(i, j int) bool {
		return statuses[i].Timestamp.Before(statuses[j].Timestamp)
	})

	var incidents []internal.Incident
	var currentIncident *internal.Incident

	for i, status := range statuses {
		if currentIncident == nil {
			currentIncident = &internal.Incident{
				Service:   service,
				StartTime: status.Timestamp,
			}
		}

		if i < len(statuses)-1 {
			timeDiff := statuses[i+1].Timestamp.Sub(status.Timestamp)
			if timeDiff.Minutes() >= 2 {
				currentIncident.EndTime = status.Timestamp
				incidents = append(incidents, *currentIncident)
				currentIncident = nil
			}
		} else {
			currentIncident.EndTime = status.Timestamp
			incidents = append(incidents, *currentIncident)
		}
	}

	return incidents, nil
}

// internals
func (s *DynamoDBStorage) toDynamoDBData(status internal.Status) (map[string]types.AttributeValue, error) {
	av, err := attributevalue.MarshalMap(map[string]interface{}{
		"service":   status.Service,
		"timestamp": status.Timestamp.Format(time.RFC3339),
		"status":    status.Status,
		"latency":   status.Latency,
	})
	if err != nil {
		return nil, err
	}
	return av, nil
}

func (s *DynamoDBStorage) batchWrite(writeRequests []types.WriteRequest) error {
	_, err := s.client.BatchWriteItem(context.TODO(), &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]types.WriteRequest{
			s.table: writeRequests,
		},
	})
	return err
}
