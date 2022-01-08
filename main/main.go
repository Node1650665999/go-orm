package main

import (
	"geeorm"
	"geeorm/log"
	"reflect"
)

type User struct {
	Name string `geeorm:"PRIMARY KEY"`
	Age  int
}

func main()  {
	Migrate()
}

//Migrate 实现数据库迁移
func Migrate() {
	engine := OpenDB()
	defer engine.Close()
	s := engine.NewSession()
	_, _ = s.Raw("DROP TABLE IF EXISTS User;").Exec()
	_, _ = s.Raw("CREATE TABLE User(Name text PRIMARY KEY, XXX integer);").Exec()
	_, _ = s.Raw("INSERT INTO User(`Name`) values (?), (?)", "Tom", "Sam").Exec()
	engine.Migrate(&User{})

	rows, _ := s.Raw("SELECT * FROM User").QueryRows()
	columns, _ := rows.Columns()
	if !reflect.DeepEqual(columns, []string{"Name", "Age"}) {
		log.Error("Failed to migrate table User, got columns", columns)
	}
}

//OpenDB 实例化 Engine
func OpenDB() *geeorm.Engine {
	engine, err := geeorm.NewEngine("sqlite3", "gee.db")
	if err != nil {
		log.Error("failed to connect", err)
	}
	return engine
}