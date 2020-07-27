package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	sdk "agones.dev/agones/sdks/go"
)

// go build -ldflags "-X main.version=..." inserts version
var version = "v0"

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	agones, err := setupAgones(ctx)
	if err != nil {
		log.Fatalf("failed to setup agones SDK: %+v", err)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)
	sig := <-sigCh
	log.Printf("Signal received: %+v", sig)
	if err := agones.Shutdown(); err != nil {
		log.Fatalf("failed to send shutdown to agones SDK: %+v", err)
	}
}

func setupAgones(ctx context.Context) (*sdk.SDK, error) {
	agones, err := sdk.NewSDK()
	if err != nil {
		return nil, err
	}
	if err := agones.Ready(); err != nil {
		return nil, err
	}
	gs, err := agones.GameServer()
	if err != nil {
		return nil, err
	}
	log.Printf("gameserver is ready! (version: %s, name: %s)", version, gs.ObjectMeta.Name)

	go func() {
		ticker := time.NewTicker(2 * time.Second)
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := agones.Health(); err != nil {
					log.Printf("failed to health to agones SDK: %+v", err)
				}
			}
		}
	}()
	return agones, err
}
