package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/structpb"

	dapr "github.com/dapr/go-sdk/client"
)

// createUserMessageInput is a helper method to create user messages in expected proto format
func createUserMessageInput(msg string) *dapr.ConversationInputAlpha2 {
	return &dapr.ConversationInputAlpha2{
		Messages: []*dapr.ConversationMessageAlpha2{
			{
				ConversationMessageOfUser: &dapr.ConversationMessageOfUserAlpha2{
					Content: []*dapr.ConversationMessageContentAlpha2{
						{
							Text: &msg,
						},
					},
				},
			},
		},
	}
}

func main() {
	client, err := dapr.NewClient()
	if err != nil {
		panic(err)
	}

	inputMsg := "What is dapr?"
	conversationComponent := "anthropic"

	// Optional: structured outputs and prompt cache retention
	responseFormat, err := structpb.NewStruct(map[string]any{
		"type": "object",
		"properties": map[string]any{
			"answer": map[string]any{"type": "string"},
		},
		"required": []any{"answer"},
	})
	if err != nil {
		log.Fatalf("failed to build response_format: %v", err)
	}

	request := dapr.ConversationRequestAlpha2{
		Name:                 conversationComponent,
		Inputs:               []*dapr.ConversationInputAlpha2{createUserMessageInput(inputMsg)},
		ResponseFormat:       responseFormat,
		PromptCacheRetention: durationpb.New(24 * time.Hour),
	}

	fmt.Println("Input sent:", inputMsg)

	resp, err := client.ConverseAlpha2(context.Background(), request)
	if err != nil {
		log.Fatalf("err: %v", err)
	}

	firstOut := resp.Outputs[0]
	if firstOut.Model != nil && *firstOut.Model != "" {
		fmt.Println("Model:", *firstOut.Model)
	}
	if firstOut.Usage != nil {
		fmt.Printf("Usage: prompt_tokens=%d completion_tokens=%d total_tokens=%d\n",
			firstOut.Usage.PromptTokens, firstOut.Usage.CompletionTokens, firstOut.Usage.TotalTokens)
	}
	fmt.Println("Output response:", firstOut.Choices[0].Message.Content)

	select {} // Keep app running
}
