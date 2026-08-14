package main

import (
	"context"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	userv1 "github.com/example/grpc-user-service/gen/user/v1"
	"github.com/example/grpc-user-service/internal/gateway"
	"github.com/example/grpc-user-service/internal/service"
	"github.com/example/grpc-user-service/internal/store"
	"github.com/example/grpc-user-service/web"
)

func main() {
	grpcAddr := flag.String("addr", ":50051", "gRPC listen address")
	httpAddr := flag.String("http", ":8080", "HTTP console / JSON gateway address")
	flag.Parse()

	users := service.NewUserServer(store.NewMemoryStore())

	grpcLis, err := net.Listen("tcp", *grpcAddr)
	if err != nil {
		log.Fatalf("grpc listen: %v", err)
	}

	grpcSrv := grpc.NewServer()
	userv1.RegisterUserServiceServer(grpcSrv, users)
	reflection.Register(grpcSrv)

	go func() {
		log.Printf("gRPC listening on %s", *grpcAddr)
		if err := grpcSrv.Serve(grpcLis); err != nil {
			log.Fatalf("grpc serve: %v", err)
		}
	}()

	mux := http.NewServeMux()
	gateway.NewHandler(users).Mount(mux)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(web.FS))))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		data, err := fs.ReadFile(web.FS, "index.html")
		if err != nil {
			http.Error(w, "index missing", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write(data)
	})

	httpSrv := &http.Server{
		Addr:              *httpAddr,
		Handler:           withCORS(mux),
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("web console listening on http://127.0.0.1%s", *httpAddr)
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http serve: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = httpSrv.Shutdown(ctx)
	grpcSrv.GracefulStop()
	fmt.Println("bye")
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
