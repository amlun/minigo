package service

import (
	"context"
	"testing"
)

func TestRun(t *testing.T) {

	t.Run("Test Run function", func(t *testing.T) {
		ctx := context.Background()
		input := "test input"
		expectedOutput := "Processed: test input"

		output, err := Run(ctx, input)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if output != expectedOutput {
			t.Fatalf("Expected output %q, got %q", expectedOutput, output)
		}
	})
}
