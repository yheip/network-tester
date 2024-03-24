package curl

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
)

func Run(ctx context.Context, url string) (<-chan string, <-chan string, error) {
	cmd := exec.CommandContext(ctx, "curl", "-vs", url)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, nil, err
	}

	outChan := make(chan string)
	errChan := make(chan string)

	go readPipeAndSend(stdout, outChan)
	go readPipeAndSend(stderr, errChan)

	go func() {
		<-ctx.Done()
		cmd.Process.Kill()
	}()

	return outChan, errChan, nil
}

func readPipeAndSend(pipe io.ReadCloser, outChan chan<- string) {
	defer close(outChan)
	reader := bufio.NewReader(pipe)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				fmt.Println("Error reading from pipe:", err)
			}
			return
		}
		outChan <- line
	}
}
