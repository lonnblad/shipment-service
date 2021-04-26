package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/lonnblad/shipment-service-backend/boundaries/rest"
	"github.com/lonnblad/shipment-service-backend/businesslogic"
	"github.com/lonnblad/shipment-service-backend/config"
	"github.com/lonnblad/shipment-service-backend/storage/go-memdb"
	"github.com/lonnblad/shipment-service-backend/trace"
)

func main() {
	ctx := context.Background()

	if err := trace.Start(); err != nil {
		log.Println(err)
		return
	}

	defer trace.Stop(ctx)

	shipmentStorage, err := memdb.NewShipmentStorage()
	if err != nil {
		log.Println(err)
		return
	}

	logic := businesslogic.New(shipmentStorage)

	restAPI, err := rest.New(config.GetRestURL(), logic)
	if err != nil {
		log.Println(err)
		return
	}

	restAPI.ListenAndServe(config.GetRestPort())

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), config.GetShutdownTimeout())
	defer cancel()

	restAPI.Shutdown(ctx)
}
