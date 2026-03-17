package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	dapr "github.com/dapr/go-sdk/client"
	"github.com/dapr/go-sdk/service/common"
	daprd "github.com/dapr/go-sdk/service/http"
)

var (
	daprClient dapr.Client
	appPort    = os.Getenv("APP_PORT")
)

func main() {
	if appPort == "" {
		appPort = "6010"
	}

	dc, err := dapr.NewClient()
	if err != nil {
		log.Fatalf("order proc: dapr client: %s", err)
	}
	daprClient = dc
	defer daprClient.Close()

	// 1. Create the Dapr Service (to listen for Cron)
	s := daprd.NewService(fmt.Sprintf(":%s", appPort))

	// 2. Register a handler for the Cron binding
	// The route name MUST match the metadata.name in your cron-binding.yaml
	s.AddBindingInvocationHandler("heartbeat-cron", cronHandler)

	log.Println("Starting server..")
	if err := s.Start(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("error listening: %v", err)
	}
}

func cronHandler(ctx context.Context, in *common.BindingEvent) (out []byte, err error) {
	log.Println("Cron binding triggered! Inserting record to Postgres...")

	// 4. Invoke Postgres Output Binding
	// This executes a SQL command defined in the 'sql' metadata field
	req := &dapr.InvokeBindingRequest{
		Name:      "postgres-binding",
		Operation: "exec",
		Data:      []byte(""), // No data needed if the query is static
		Metadata: map[string]string{
			"sql": "INSERT INTO logs (message, created_at) VALUES ('Dapr was here', now());",
		},
	}

	err = daprClient.InvokeOutputBinding(ctx, req)
	if err != nil {
		log.Printf("Failed to write to postgres: %v", err)
		return nil, err
	}

	return []byte("Success"), nil
}
