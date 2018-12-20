package otaws

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/mocktracer"
	"testing"
	"time"
)

// Test requires running local instance of DynamoDB
func TestAWS(t *testing.T) {
	tracer := mocktracer.New()
	opentracing.InitGlobalTracer(tracer)

	sess, err := session.NewSession(&aws.Config{
		Region:   aws.String("us-west-2"),
		Endpoint: aws.String("http://localhost:8000")})
	if err != nil {
		return
	}
	dbClient := dynamodb.New(sess)

	AddOTHandlers(dbClient.Client)

	fmt.Println("ListTables")
	result, err := dbClient.ListTables(&dynamodb.ListTablesInput{})
	if err != nil {
		fmt.Println("Error")
	}

	for _, table := range result.TableNames {
		fmt.Println(table)
	}

	spans := tracer.FinishedSpans()

	count := 20

	for len(spans) != 1 && count > 0 {
		spans = tracer.FinishedSpans()
		count -= 1
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println(spans)

	if len(spans) != 1 {
		t.Errorf("Expected 1 finished span. Found: %d", len(spans))
	}
}
