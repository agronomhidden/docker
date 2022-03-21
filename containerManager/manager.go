package containerManager

import (
	"bufio"
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"io"
	"os"
	"time"
)

type manager struct {
	doneCh chan struct{}
	cli    *client.Client
	cnt    container.ContainerCreateCreatedBody
}

func New(ctx context.Context, image, cmd string) (_ *manager, err error) {
	instance := &manager{}
	if instance.cli, err = client.NewClientWithOpts(); err != nil {
		return nil, err
	}

	reader, err := instance.cli.ImagePull(ctx, image, types.ImagePullOptions{})
	if err != nil {
		return nil, err
	}
	io.Copy(os.Stdout, reader)

	cntConf := &container.Config{
		Image: image,
		Cmd:   []string{"sh", "-c", cmd},
		Tty:   true,
	}
	hostConf := &container.HostConfig{AutoRemove: true}

	if instance.cnt, err = instance.cli.ContainerCreate(ctx, cntConf, hostConf, nil, nil, ""); err != nil {
		return nil, err
	}
	return instance, nil
}

func (e *manager) Run(ctx context.Context) <-chan string {
	resCh := make(chan string, 10)
	stop := make(chan struct{})

	if err := e.cli.ContainerStart(context.Background(), e.cnt.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}
	go func() {
		reader, err := e.cli.ContainerLogs(ctx, e.cnt.ID, types.ContainerLogsOptions{
			ShowStdout: true,
			ShowStderr: true,
			Follow:     true,
			Timestamps: false,
		})
		if err != nil {
			panic(err)
		}
		defer reader.Close()

		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			resCh <- scanner.Text()
		}
		<-stop
		close(resCh)
	}()
	go func() {
		<-ctx.Done()
		d := time.Second * 10
		if err := e.cli.ContainerStop(context.Background(), e.cnt.ID, &d); err != nil {
			panic(err)
		}
	}()
	go func() {
		stateCh, errCh := e.cli.ContainerWait(context.Background(), e.cnt.ID, container.WaitConditionNotRunning)
		select {
		case <-stateCh:
			stop <- struct{}{}
		case <-errCh:
			stop <- struct{}{}
		}
	}()
	return resCh
}
