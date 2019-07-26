package db

import (
	"log"

	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq"
)

var engine *xorm.Engine

func newPostgresEngine() (*xorm.Engine, error) {
	var err error
	engine, err = xorm.NewEngine("postgres", postgresURI)
	if err == nil {
		engine.SetMaxIdleConns(2)
		engine.ShowSQL(true)
	} else {
		log.Fatal(err)
	}
	return engine, err
}

func GetEngine() *xorm.Engine {
	if engine != nil {
		return engine
	}

	engine, _ = newPostgresEngine()
	return engine
}
