package service

import "context"

func Run(ctx context.Context, input string) (string, error) {
	// Implement your business logic here
	return "Processed: " + input, nil
}
