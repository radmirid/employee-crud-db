package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/radmirid/employee-crud-db/internal/config"
	psql "github.com/radmirid/employee-crud-db/internal/repository/psql"
	"github.com/radmirid/employee-crud-db/internal/service"
	grpc "github.com/radmirid/employee-crud-db/internal/transport/grpc"
	rest "github.com/radmirid/employee-crud-db/internal/transport/rest"
	db "github.com/radmirid/employee-crud-db/pkg/db"
	"github.com/radmirid/employee-crud-db/pkg/hash"

	_ "github.com/lib/pq"

	log "github.com/sirupsen/logrus"
)

const (
	CONFIG_DIR  = "configs"
	CONFIG_FILE = "main"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func main() {
	cfg, err := config.New(CONFIG_DIR, CONFIG_FILE)
	if err != nil {
		log.Fatal(err)
	}

	db, err := db.NewPostgresConnection(db.ConnectionInfo{
		Host:     cfg.DB.Host,
		Port:     cfg.DB.Port,
		Username: cfg.DB.Username,
		DBName:   cfg.DB.Name,
		SSLMode:  cfg.DB.SSLMode,
		Password: cfg.DB.Password,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	hasher := hash.NewSHA1Hasher("salt")

	employeesRepo := psql.NewEmployees(db)
	employeesService := service.NewEmployees(employeesRepo)

	usersRepo := psql.NewUsers(db)
	tokensRepo := psql.NewTokens(db)

	grpcLogger, err := grpc.NewClient(9000)
	if err != nil {
		log.Fatal(err)
	}

	usersService := service.NewUsers(usersRepo, tokensRepo, grpcLogger, hasher, []byte("secret"))

	handler := rest.NewHandler(employeesService, usersService)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: handler.InitRouter(),
	}

	log.Info("Running")

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
