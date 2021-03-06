package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"net/http"
	"os"
	"time"
)

func checkEndpoint(endpoint string) int {
	resp, err := http.Get(endpoint)
	if err != nil {
		return 503
	}

	return resp.StatusCode
}

func healthcheck() {
	endpoint := os.Getenv("HTTP_HEALTHCHECK_ENDPOINT")
	metricName := "http-healthcheck-" + os.Getenv("APP_NAME") + "-" + os.Getenv("AWS_REGION")
	timestamp := time.Now()
	metricValue := float64(checkEndpoint(endpoint))

	data := cloudwatch.MetricDatum{
		MetricName: &metricName,
		Timestamp:  &timestamp,
		Value:      &metricValue,
	}

	metricDataList := []*cloudwatch.MetricDatum{&data}
	metricDataInput := cloudwatch.PutMetricDataInput{
		MetricData: metricDataList,
		Namespace:  &metricName,
	}

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	}))
	svc := cloudwatch.New(sess)
	svc.PutMetricData(&metricDataInput)
}

func main() {
	lambda.Start(healthcheck)
}
