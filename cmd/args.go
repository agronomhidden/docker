package cmd

import (
	"flag"
	"fmt"
	"os"
)

type Args struct {
	DockerImage      string
	BashCommand      string
	CloudWatchGroup  string
	CloudWatchStream string
	AwsAccessKey     string
	AwsSecretKey     string
	AwsRegion        string
}

func (e *Args) Parse() {
	flag.StringVar(&e.DockerImage, "docker-image", "", "docker image")
	flag.StringVar(&e.BashCommand, "bash-command", "", "bash command")
	flag.StringVar(&e.CloudWatchGroup, "cloudwatch-group", "", "cloudwatch group")
	flag.StringVar(&e.CloudWatchStream, "cloudwatch-stream", "", "cloudwatch stream")
	flag.StringVar(&e.AwsAccessKey, "aws-access-key-id", "", "aws-access-key-id")
	flag.StringVar(&e.AwsSecretKey, "aws-secret-access-key", "", "aws-secret-access-key")
	flag.StringVar(&e.AwsRegion, "aws-region", "", "aws region")
	flag.Parse()
}
func (e *Args) Validate() {
	if e.DockerImage == "" || e.BashCommand == "" || e.AwsAccessKey == "" || e.AwsSecretKey == "" || e.AwsRegion == "" || e.CloudWatchGroup == "" || e.CloudWatchStream == "" {
		fmt.Println("some of required params are empty")
		flag.PrintDefaults()
		os.Exit(1)
	}
}
