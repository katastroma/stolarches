//revive:disable:package-comments
package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	healthpb "google.golang.org/grpc/health/grpc_health_v1"

	foundationclient "git.sonicoriginal.software/grpc-foundation/client"
	foundationotel "git.sonicoriginal.software/grpc-foundation/otel"
	"git.sonicoriginal.software/grpc-foundation/server"

	pb "github.com/katastroma/diataxis"

	"github.com/katastroma/stolarches/internal/config"
	grpc_health "github.com/katastroma/stolarches/internal/health/grpc"
	http_health "github.com/katastroma/stolarches/internal/health/http"
	"github.com/katastroma/stolarches/internal/order"
	"github.com/katastroma/stolarches/internal/order/cliutils"
	"github.com/katastroma/stolarches/internal/order/helm"
	"github.com/katastroma/stolarches/internal/provisioner"
	"github.com/katastroma/stolarches/internal/serve"
)

const shutdownTimeout = 10 * time.Second

func main() {
	mainCtx := context.Background()

	log, shutdown, err := foundationotel.Init(mainCtx, "stolarches")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize otel: %v\n", err)
		os.Exit(1)
	}
	defer shutdown(mainCtx)

	provisionerAddr, err := config.RequireEnv("PROVISIONER_ADDR")
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	provisionerConn, err := foundationclient.New(provisionerAddr, nil, nil)
	if err != nil {
		log.Error("provisioner connection failed", "error", err)
		os.Exit(1)
	}
	defer provisionerConn.Close()

	var router order.Router
	router.Register(pb.OrdererType_ORDERER_TYPE_HELM, helm.New())
	router.Register(pb.OrdererType_ORDERER_TYPE_CLI_UTILS, cliutils.New())

	streamFn := provisioner.NewStreamFunc(log, provisionerConn)
	service := serve.New(log, &router, streamFn)

	mux := http.NewServeMux()
	mux.Handle("GET /healthz", http_health.New(log))

	httpPort := config.StringEnv("PORT", "8080")
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", httpPort),
		Handler: mux,
	}

	grpcServer := server.New(log)
	healthpb.RegisterHealthServer(grpcServer, grpc_health.New())
	pb.RegisterOrdererServiceServer(grpcServer, service)

	sigNotifyContext, stop := context.WithCancel(mainCtx)
	defer stop()

	var wg sync.WaitGroup

	wg.Go(func() {
		log.Info("starting http server", "port", httpPort)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("http server error", "error", err)
			stop()
		}
	})

	wg.Go(func() {
		addr := server.Address()
		lis, err := net.Listen("tcp", addr)
		if err != nil {
			log.Error("grpc listen error", "error", err)
			stop()
			return
		}
		log.Info("starting grpc server", "address", addr)
		if err := grpcServer.Serve(lis); err != nil {
			log.Error("grpc server error", "error", err)
			stop()
		}
	})

	server.HandleGracefulShutdown(sigNotifyContext, stop, log, grpcServer, shutdownTimeout)

	shutdownCtx, shutdownCancel := context.WithTimeout(mainCtx, shutdownTimeout)
	defer shutdownCancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Error("http shutdown error", "error", err)
	}

	wg.Wait()
}
