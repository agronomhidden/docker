package awsClient

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"time"
)

type Client struct {
	svc           *cloudwatchlogs.CloudWatchLogs
	group         string
	stream        string
	sequenceToken *string
}

func New(awsRegion, awsKeyId, awsKeySecret, logGroupName, logStreamName string) (*Client, error) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(awsRegion),
		Credentials: credentials.NewStaticCredentials(awsKeyId, awsKeySecret, ""),
	}))
	instance := &Client{
		svc:    cloudwatchlogs.New(sess),
		group:  logGroupName,
		stream: logStreamName,
	}
	if err := instance.createLogGroupIfNotExist(); err != nil {
		return nil, err
	}
	if err := instance.createLogStreamIfNotExist(); err != nil {
		return nil, err
	}

	return instance, nil
}

func (e *Client) createLogGroupIfNotExist() error {
	listGroupResp, err := e.svc.DescribeLogGroups(&cloudwatchlogs.DescribeLogGroupsInput{})
	if err != nil {
		return err
	}
	for _, group := range listGroupResp.LogGroups {
		if *group.LogGroupName == e.group {
			return nil
		}
	}
	_, err = e.svc.CreateLogGroup(&cloudwatchlogs.CreateLogGroupInput{LogGroupName: &e.group})

	return err
}

func (e *Client) createLogStreamIfNotExist() error {
	res, err := e.svc.DescribeLogStreams(&cloudwatchlogs.DescribeLogStreamsInput{LogGroupName: &e.group})
	if err != nil {
		return err
	}
	for _, item := range res.LogStreams {
		if *item.LogStreamName == e.stream {
			e.sequenceToken = item.UploadSequenceToken
			return nil
		}
	}
	_, err = e.svc.CreateLogStream(&cloudwatchlogs.CreateLogStreamInput{LogGroupName: &e.group, LogStreamName: &e.stream})

	return err
}

func (e *Client) Push(message string) error {
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	resp, err := e.svc.PutLogEvents(&cloudwatchlogs.PutLogEventsInput{LogEvents: []*cloudwatchlogs.InputLogEvent{
		{Timestamp: &timestamp, Message: &message},
	},
		LogGroupName:  &e.group,
		LogStreamName: &e.stream,
		SequenceToken: e.sequenceToken,
	})
	if err != nil {
		return err
	}
	e.sequenceToken = resp.NextSequenceToken
	return nil
}

// @todo for debug :)
func (e *Client) List() error {
	var limit int64 = 100
	list, err := e.svc.GetLogEvents(&cloudwatchlogs.GetLogEventsInput{

		Limit:         &limit,
		LogGroupName:  &e.group,
		LogStreamName: &e.stream,
	})

	if err != nil {
		return err
	}

	for _, item := range list.Events {
		fmt.Println(*item.Timestamp, *item.Message)
	}
	return nil
}
