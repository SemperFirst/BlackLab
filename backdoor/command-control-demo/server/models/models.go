package models

import (
	"fmt"
	"os"
	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
)

var (
	Engine *xorm.Engine
)

func init() {
	Engine, err := NewDbEngine()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	err = Engine.Sync2(new(Agent))
	err = Engine.Sync2(new(Command))
	fmt.Println(err)
}

func NewDbEngine() (*xorm.Engine, error) {
	engine, err := xorm.NewEngine("sqlite3", "c_c.db")
	return engine, err
}
