package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/dapr/go-sdk/service/common"
	daprd "github.com/dapr/go-sdk/service/grpc"
	"github.com/google/uuid"
)

var (
	appPort = os.Getenv("APP_PORT")
)

func main() {
	if appPort == "" {
		appPort = "5050"
	}

	dapr, err := daprd.NewService(fmt.Sprintf(":%s", appPort))
	if err != nil {
		log.Fatalf("genid: service creation: %v", err)
	}

	// Define service endpoint /genid
	if err := dapr.AddServiceInvocationHandler("/genid", generateId); err != nil {
		log.Fatalf("genid: invocation handler setup: %v", err)
	}

	// start the service
	if err := dapr.Start(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("genid start: %v", err)
	}
}

// generateId service handler
func generateId(ctx context.Context, in *common.InvocationEvent) (*common.Content, error) {
	id := uuid.New()
	out := &common.Content{
		Data:        []byte(id.String()),
		ContentType: in.ContentType,
		DataTypeURL: in.DataTypeURL,
	}

	return out, nil
}
