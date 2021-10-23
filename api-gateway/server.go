package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime/pprof"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"

	grpcClient "api-gateway/internal/grpc"
	httpHandler "api-gateway/internal/http"
	"api-gateway/internal/middleware"
	"api-gateway/pkg/proto"
)

const (
	defaultHTTPPort      = ":8000"
	defaultPProfHTTPPort = ":6060"
)

func init() {
	_ = godotenv.Load()
}

func main() {
	pprofServer := new(http.Server)
	enableProf, _ := strconv.ParseBool(os.Getenv("ENABLE_PPROF"))
	if enableProf {
		if cpuProfile := os.Getenv("CPU_PPROF_FILE_PATH"); cpuProfile != "" {
			f, err := os.Create(cpuProfile)
			if err != nil {
				log.Println(err)
			} else {
				defer func() {
					if err := f.Close(); err != nil {
						log.Printf("Error closing cpu profile file: %v", err)
					}
				}()

				_ = pprof.StartCPUProfile(f)
			}
		}
		if memProfile := os.Getenv("MEMORY_PPROF_FILE_PATH"); memProfile != "" {
			f, err := os.Create(memProfile)
			if err != nil {
				log.Println(err)
			} else {
				defer func() {
					if err := f.Close(); err != nil {
						log.Printf("Error closing memory profile file: %v", err)
					}
				}()

				_ = pprof.WriteHeapProfile(f)
			}
		}

		pprofHTPPPort := os.Getenv("PPROF_HTTP_PORT")
		if pprofHTPPPort == "" {
			pprofHTPPPort = defaultPProfHTTPPort
		}

		pprofServer = &http.Server{
			Addr:    pprofHTPPPort,
			Handler: nil,
		}

		go func() {
			if err := pprofServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Printf("Error starting pprof server: %v", err)
			}
		}()
	}

	userGRPCClientConn, err := grpc.Dial(
		fmt.Sprintf("%s%s", os.Getenv("USER_SERVICE_HOST"), os.Getenv("USER_SERVICE_PORT")),
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatal(err)
	}

	bookGRPCClientConn, err := grpc.Dial(
		fmt.Sprintf("%s%s", os.Getenv("BOOK_SERVICE_HOST"), os.Getenv("BOOK_SERVICE_PORT")),
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatal(err)
	}

	lendingGRPCClientConn, err := grpc.Dial(
		fmt.Sprintf("%s%s", os.Getenv("LENDING_SERVICE_HOST"), os.Getenv("LENDING_SERVICE_PORT")),
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatal(err)
	}

	userServiceClient := proto.NewUserServiceClient(userGRPCClientConn)
	bookServiceClient := proto.NewBookServiceClient(bookGRPCClientConn)
	lendingServiceClient := proto.NewLendingServiceClient(lendingGRPCClientConn)

	userGRPCService := grpcClient.NewUserGRPCService(userServiceClient)
	bookGRPCService := grpcClient.NewBookGRPCService(bookServiceClient)
	lendingGRPCService := grpcClient.NewLendingGRPCService(lendingServiceClient)

	server := gin.Default()
	server.GET("/", httpHandler.GraphPlaygroundHandler())
	server.POST("/query", middleware.GinJWT(), httpHandler.GraphQLHandler(
		userGRPCService,
		bookGRPCService,
		lendingGRPCService,
	))

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = defaultHTTPPort
	}

	httpServer := &http.Server{
		Addr:    httpPort,
		Handler: server,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	wg := new(sync.WaitGroup)
	wg.Add(1)

	go func() {
		defer wg.Done()
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := httpServer.Shutdown(ctx); err != nil {
			log.Println(err)
		}
		pprof.StopCPUProfile()
		if err := pprofServer.Shutdown(ctx); err != nil {
			log.Println(err)
		}
	}()

	if err = httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
	wg.Wait()
	log.Println("service is gracefully shutdown")
}
