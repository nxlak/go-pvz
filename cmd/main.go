package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"

	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/nxlak/go-pvz/internal/config"

	//openapi "github.com/nxlak/go-pvz/internal/handler/http/server"
	"github.com/nxlak/go-pvz/internal/handler/grpc/server"
	//customMiddleware "github.com/nxlak/go-pvz/internal/middleware"
	orderPostgres "github.com/nxlak/go-pvz/internal/repository/storage/postgres"
	"github.com/nxlak/go-pvz/pkg/client/postgres"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	//order_v1 "github.com/nxlak/go-pvz/pkg/openapi/order/v1"
	orderV1 "github.com/nxlak/go-pvz/pkg/proto/order/v1"
)

const (
	httpPort          = 8080
	grpcPort          = 8081
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

func main() {
	// gRPC VERSION
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		return
	}
	defer func() {
		if cerr := lis.Close(); cerr != nil {
			log.Printf("failed to close listener: %v\n", cerr)
		}
	}()

	cfg := config.GetConfig()

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	client, err := postgres.NewClient(ctx, cfg.Storage)
	if err != nil {
		log.Fatalf("err %v", err)
	}

	orderRepo := orderPostgres.NewRepositoty(client)

	// Создаем gRPC сервер
	s := grpc.NewServer()

	// Регистрируем наш сервис
	service := server.NewOrderService(orderRepo)
	orderV1.RegisterOrderServiceServer(s, service)

	// Включаем рефлексию для отладки
	reflection.Register(s)

	go func() {
		log.Printf("🚀 gRPC server listening on %d\n", grpcPort)
		err = s.Serve(lis)
		if err != nil {
			log.Printf("failed to serve: %v\n", err)
			return
		}
	}()

	// Запускаем HTTP сервер с gRPC Gateway и Swagger UI
	var gwServer *http.Server
	go func() {
		// Создаем контекст с отменой
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// Создаем мультиплексор для HTTP запросов
		mux := runtime.NewServeMux()

		// Настраиваем опции для соединения с gRPC сервером
		opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

		// Регистрируем gRPC-gateway хендлеры
		err = orderV1.RegisterOrderServiceHandlerFromEndpoint(
			ctx,
			mux,
			fmt.Sprintf("localhost:%d", grpcPort),
			opts,
		)
		if err != nil {
			log.Printf("Failed to register gateway: %v\n", err)
			return
		}

		// Создаем файловый сервер для swagger-ui
		fileServer := http.FileServer(http.Dir("../api/swagger"))

		// Создаем HTTP маршрутизатор
		httpMux := http.NewServeMux()

		// Регистрируем API эндпоинты
		httpMux.Handle("/api/", mux)

		// Swagger UI эндпоинты
		httpMux.Handle("/swagger-ui.html", fileServer)
		httpMux.Handle("/swagger.json", fileServer)

		// Редирект с корня на Swagger UI
		httpMux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/" {
				http.Redirect(w, r, "/swagger-ui.html", http.StatusMovedPermanently)
				return
			}
			fileServer.ServeHTTP(w, r)
		}))

		// Создаем HTTP сервер
		gwServer = &http.Server{
			Addr:              fmt.Sprintf(":%d", httpPort),
			Handler:           httpMux,
			ReadHeaderTimeout: readHeaderTimeout,
		}

		// Запускаем HTTP сервер
		log.Printf("🌐 HTTP server with gRPC-Gateway and Swagger UI listening on %d\n", httpPort)
		err = gwServer.ListenAndServe()
		if err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Printf("Failed to serve HTTP: %v\n", err)
			return
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down servers...")

	// Сначала аккуратно останавливаем HTTP сервер
	if gwServer != nil {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := gwServer.Shutdown(shutdownCtx); err != nil {
			log.Printf("HTTP server shutdown error: %v", err)
		}
		log.Println("✅ HTTP server stopped")
	}

	// В конце останавливаем gRPC сервер
	s.GracefulStop()
	log.Println("✅ gRPC server stopped")

	// HTTP VERSION
	// orderHandler := openapi.NewOrderHandler(orderRepo)

	// orderServer, err := order_v1.NewServer(orderHandler)
	// if err != nil {
	// 	log.Fatalf("ошибка создания сервера OpenAPI: %v", err)
	// }

	// r := chi.NewRouter()
	// r.Use(middleware.Logger)
	// r.Use(middleware.Recoverer)
	// r.Use(middleware.Timeout(10 * time.Second))
	// r.Use(customMiddleware.RequestLogger)

	// r.Mount("/", orderServer)

	// server := &http.Server{
	// 	Addr:              net.JoinHostPort("localhost", httpPort),
	// 	Handler:           r,
	// 	ReadHeaderTimeout: readHeaderTimeout, // Защита от Slowloris атак - тип DDoS-атаки, при которой
	// 	// атакующий умышленно медленно отправляет HTTP-заголовки, удерживая соединения открытыми и истощая
	// 	// пул доступных соединений на сервере. ReadHeaderTimeout принудительно закрывает соединение,
	// 	// если клиент не успел отправить все заголовки за отведенное время.
	// }

	// go func() {
	// 	log.Printf("🚀 HTTP-сервер запущен на порту %s\n", httpPort)
	// 	err = server.ListenAndServe()
	// 	if err != nil && !errors.Is(err, http.ErrServerClosed) {
	// 		log.Printf("❌ Ошибка запуска сервера: %v\n", err)
	// 	}
	// }()

	// // Graceful shutdown
	// quit := make(chan os.Signal, 1)
	// signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// <-quit

	// log.Println("🛑 Завершение работы сервера...")

	// err = server.Shutdown(ctx)
	// if err != nil {
	// 	log.Printf("❌ Ошибка при остановке сервера: %v\n", err)
	// }

	// log.Println("✅ Сервер остановлен")
}
