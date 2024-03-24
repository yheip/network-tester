package curl

import (
	"context"
	"fmt"
	"testing"
)

func TestRun(t *testing.T) {
	ctx := context.Background()
	url := "https://www.google.com"

	outChan, errChan, err := Run(ctx, url)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Print stdout
	go func() {
		for output := range outChan {
			fmt.Print(output)
		}
	}()

	// Print stderr
	go func() {
		for output := range errChan {
			fmt.Print(output)
		}
	}()

	// Waiting for both channels to close
	<-outChan
	<-errChan
}
