package main

import (
	"context"
	"crud/cmd/crud/app"
	"crud/pkg/crud/services"
	"crud/pkg/crud/services/burgers"
	"flag"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"net"
	"net/http"
	"path/filepath"
)

const (
	EHost = "HOST"
	EPort = "PORT"
	EDsn  = "DATABASE_URL"
)

var (
	hostF = flag.String("host", "", "Server host")
	portF = flag.String("port", "", "Server port")
	dsnF  = flag.String("dsn", "", "Postgres DSN")
)

func main() {
	flag.Parse()
	host, ok := services.FlagOrEnv(*hostF, EHost)
	if !ok {
		log.Panic("can't get port")
	}

	port, ok := services.FlagOrEnv(*portF, EPort)
	if !ok {
		log.Panic("can't get port")
	}

	log.Println("set address to connect")
	addr := net.JoinHostPort(host, port)
	log.Printf("address to connect: %s", addr)

	dsn, ok := services.FlagOrEnv(*dsnF, EDsn)
	if !ok {
		log.Panic("can't get dsn")
	}
	start(addr, dsn)
}

func start(addr string, dsn string) {
	router := app.NewExactMux()
	pool, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		panic(err)
	}

	burgersSvc := burgers.NewBurgersSvc(pool)
	err = burgersSvc.InitDB()
	server := app.NewServer(
		router,
		pool,
		burgersSvc, // DI + Containers
		filepath.Join("web", "templates"),
		filepath.Join("web", "assets"),
	)
	server.InitRoutes()

	panic(http.ListenAndServe(addr, server))
}
