package main

import (
	"context"
	"fmt"
	"github.com/agronomhidden/docker/awsClient"
	"github.com/agronomhidden/docker/cmd"
	"github.com/agronomhidden/docker/containerManager"
	"os"
	"os/signal"
	"syscall"
)
var args = cmd.Args{}

func init() {
	args.Parse()
	args.Validate()
}

func main() {
	defer recoverPanic()

	client, err := awsClient.New(args.AwsRegion, args.AwsAccessKey, args.AwsSecretKey, args.CloudWatchGroup, args.CloudWatchStream)
	if err != nil {
		fmt.Printf("unable to init awsClient: %s", err.Error())
		return
	}

	manager, err := containerManager.New(context.Background(), args.DockerImage, args.BashCommand)
	if err != nil {
		fmt.Printf("unable to create container: %s", err.Error())
		return

	}
	ctx, cancel := context.WithCancel(context.Background())
	streamCh := manager.Run(ctx)

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGALRM)
		<-c
		fmt.Println("stopping service...")
		//graceful shutdown:
		cancel()
	}()

	for {
		message, ok := <-streamCh
		if !ok {
			fmt.Println("container output was complete")
			break
		}
		if err := client.Push(message); err != nil {
			fmt.Println("error pushing message to aws")
		}
	}
	client.List()
	fmt.Println("done")
}

func recoverPanic() {
	if r := recover(); r != nil {
		switch x := r.(type) {
		case string:
			fmt.Println(x)
		case error:
			fmt.Println(x.Error())
		default:
			fmt.Println("unknown panic")
		}
	}
}
