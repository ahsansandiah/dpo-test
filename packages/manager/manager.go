package manager

import (
	"context"
	"database/sql"

	httpClient "github.com/ahsansandiah/dpo-test/packages/client"
	"github.com/ahsansandiah/dpo-test/packages/config"
	"github.com/ahsansandiah/dpo-test/packages/json"
	logger "github.com/ahsansandiah/dpo-test/packages/log"
	"github.com/ahsansandiah/dpo-test/packages/server"
	database "github.com/ahsansandiah/dpo-test/packages/storage/mysql"
)

type Manager interface {
	GetConfig() *config.Config
	GetServer() *server.Server
	GetDB() *sql.DB
	GetLog() logger.Log
	GetJson() json.Json
	GetHttp() httpClient.Http
}

type manager struct {
	config     *config.Config
	server     *server.Server
	db         *sql.DB
	logger     logger.Log
	json       json.Json
	httpClient httpClient.Http
}

func NewInit() (Manager, error) {
	lg := logger.NewLog()
	ctx := context.Background()

	cfg, err := config.NewConfig()
	if err != nil {
		lg.ErrorLog(ctx, err)
		return nil, err
	}

	srv := server.NewServer(cfg)
	database, err := database.NewMySQL(cfg).Connect()
	if err != nil {
		lg.ErrorLog(ctx, err)
		return nil, err
	}

	clHttp := httpClient.NewHttp(cfg, lg)
	clHttp.Connect()

	json := json.NewJson(lg)

	return &manager{
		config:     cfg,
		server:     srv,
		db:         database,
		logger:     lg,
		httpClient: clHttp,
		json:       json,
	}, nil
}

func (sm *manager) GetConfig() *config.Config {
	return sm.config
}

func (sm *manager) GetServer() *server.Server {
	return sm.server
}

func (sm *manager) GetDB() *sql.DB {
	return sm.db
}

func (sm *manager) GetLog() logger.Log {
	return sm.logger
}

func (sm *manager) GetJson() json.Json {
	return sm.json
}

func (sm *manager) GetHttp() httpClient.Http {
	return sm.httpClient
}
