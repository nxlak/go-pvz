package main

import (
	"context"
	"log"
	"time"

	"github.com/nxlak/go-pvz/internal/config"
	"github.com/nxlak/go-pvz/internal/domain/model"
	orderPostgres "github.com/nxlak/go-pvz/internal/repository/storage/postgres"
	"github.com/nxlak/go-pvz/pkg/client/postgres"
)

func main() {
	cfg := config.GetConfig()

	client, err := postgres.NewClient(context.TODO(), cfg.Storage)
	if err != nil {
		log.Fatalf("err %v", err)
	}

	orderRepo := orderPostgres.NewRepositoty(client)
	testOrder := model.Order{
		Id:        "1",
		UserId:    "123",
		Status:    model.StatusAccepted,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(time.Hour * 24),
	}
	if err := orderRepo.Create(context.TODO(), &testOrder); err != nil {
		log.Fatalf("err %v", err)
	}

}
