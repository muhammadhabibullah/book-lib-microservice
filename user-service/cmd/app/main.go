package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"runtime/pprof"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"user-service/internal/service"
	"user-service/pkg/mongodb"
	"user-service/pkg/proto"
)

const (
	defaultGRPCPort      = ":8000"
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

	mongodb.GetDatabase()

	userService := service.NewUserGRPCService()
	server := grpc.NewServer()
	proto.RegisterUserServiceServer(server, userService)

	reflection.Register(server)

	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = defaultGRPCPort
	}

	listener, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalln(err)
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

		server.GracefulStop()
		pprof.StopCPUProfile()
		if err := pprofServer.Shutdown(ctx); err != nil {
			log.Println(err)
		}
	}()

	log.Println("starting to serve")
	if err = server.Serve(listener); err != nil {
		log.Println(err)
	}
	wg.Wait()
	log.Println("service is gracefully shutdown")
}
