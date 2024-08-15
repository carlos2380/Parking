package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"parking/internal/errors"
	"parking/internal/handlers"
	"parking/internal/register/redis"
	"parking/internal/server"
	"strconv"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	port := flag.String("port", "8000", "Port on which the server will be listening for incoming requests.")
	portRedis := flag.String("port_redis", "6379", "Port on which the Redis server will be serving connections.")
	ipRedis := flag.String("ip_redis", "0.0.0.0", "IP address on which the Redis server will be listening.")
	flag.Parse()

	portRedisInt, err := strconv.Atoi(*portRedis)
	if err != nil {
		log.Fatal(errors.Wrap(err, *errors.ErrIntConvert))
	}
	registerRedis := &redis.RegisterRedis{}
	if err := registerRedis.InitRedis(*ipRedis, portRedisInt, ""); err != nil {
		log.Fatal("Error redis:", err)
	}

	pHandler := &handlers.ParckingHandler{
		Register: registerRedis,
	}

	router := server.NewRouter(pHandler).(*mux.Router)

	srv := &http.Server{
		Addr:    ":" + *port,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	log.Print("Server Started")
	log.Printf("Listening on 0.0.0.0:%s", *port)
	<-done

	log.Print("Server Stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Print("Server Exited Properly")
}
