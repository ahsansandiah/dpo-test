package manager

import (
	"database/sql"

	"github.com/DATA-DOG/go-sqlmock"
	httpClient "github.com/ahsansandiah/dpo-test/packages/client"
	"github.com/ahsansandiah/dpo-test/packages/config"
	"github.com/ahsansandiah/dpo-test/packages/json"
	logger "github.com/ahsansandiah/dpo-test/packages/log"
	"github.com/ahsansandiah/dpo-test/packages/server"
	"github.com/golang/mock/gomock"
)

type FakeManager interface {
	Manager
}

type fakeManager struct {
	config     *config.Config
	server     *server.Server
	database   *sql.DB
	logger     logger.Log
	httpClient httpClient.Http
	json       json.Json
}

func NewFakeInit(ctrl *gomock.Controller) (FakeManager, error) {
	lg := logger.NewMockLog(ctrl)

	cfg := &config.Config{}

	srv := server.NewServer(cfg)

	dbMysql, _, err := sqlmock.New()
	if err != nil {
		return nil, err
	}
	defer dbMysql.Close()

	json := json.NewMockJson(ctrl)
	clHttp := httpClient.NewMockHttp(ctrl)

	return &fakeManager{
		config:     cfg,
		server:     srv,
		database:   dbMysql,
		logger:     lg,
		json:       json,
		httpClient: clHttp,
	}, nil
}

func (fm *fakeManager) GetConfig() *config.Config {
	return fm.config
}

func (fm *fakeManager) GetServer() *server.Server {
	return fm.server
}

func (fm *fakeManager) GetDB() *sql.DB {
	return fm.database
}

func (fm *fakeManager) GetLog() logger.Log {
	return fm.logger
}

func (fm *fakeManager) GetJson() json.Json {
	return fm.json
}

func (fm *fakeManager) GetHttp() httpClient.Http {
	return fm.httpClient
}
