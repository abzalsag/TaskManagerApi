package main

import (
	"TaskManagerApi/kafka"
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"TaskManagerApi/database"
	"TaskManagerApi/grpcserver"
	"TaskManagerApi/handler"
	"TaskManagerApi/logger"
	"TaskManagerApi/middleware"
	"TaskManagerApi/proto"
	"TaskManagerApi/repository"
	"TaskManagerApi/routes"
	"TaskManagerApi/service"

	"google.golang.org/grpc"
)

func main() {
	appLogger := logger.New()

	db, err := database.NewPostgres()
	if err != nil {
		appLogger.Error("database connection failed", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	appLogger.Info("database connected")

	taskRepository := repository.NewTaskRepository(db)
	userRepository := repository.NewUserRepository(db)

	taskProducer := kafka.NewProducer()
	defer taskProducer.Close()
	consumerCtx, consumerCancel := context.WithCancel(context.Background())
	defer consumerCancel()
	go kafka.StartConsumer(consumerCtx)

	taskService := service.NewTaskService(taskRepository, taskProducer)
	userService := service.NewUserService(userRepository)

	taskHandler := handler.NewTaskHandler(taskService)
	userHandler := handler.NewUserHandler(userService)

	mux := http.NewServeMux()

	routes.RegisterRoutes(
		mux,
		taskHandler,
		userHandler,
	)

	var httpHandler http.Handler = mux

	httpHandler = middleware.Logger(appLogger)(httpHandler)
	httpHandler = middleware.Recovery(appLogger)(httpHandler)

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: httpHandler,
	}

	grpcServer := grpc.NewServer()

	grpcTaskServer := grpcserver.NewTaskGRPCServer(taskService)

	proto.RegisterTaskServiceServer(
		grpcServer,
		grpcTaskServer,
	)

	grpcListener, err := net.Listen("tcp", ":50051")
	if err != nil {
		appLogger.Error("gRPC listener failed", "error", err)
		os.Exit(1)
	}

	go func() {
		appLogger.Info("gRPC server started", "port", 50051)

		if err := grpcServer.Serve(grpcListener); err != nil {
			appLogger.Error("gRPC server stopped", "error", err)
		}
	}()

	go func() {
		appLogger.Info("HTTP server started", "port", 8080)

		if err := httpServer.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {
			appLogger.Error("HTTP server error", "error", err)
		}
	}()

	stop := make(chan os.Signal, 1)

	signal.Notify(
		stop,
		os.Interrupt,
		syscall.SIGTERM,
	)

	<-stop

	appLogger.Info("shutdown signal received")

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		appLogger.Error("HTTP shutdown error", "error", err)
	}

	grpcServer.GracefulStop()

	consumerCancel()
	appLogger.Info("servers stopped gracefully")
}
