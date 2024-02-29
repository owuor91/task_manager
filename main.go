package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
	"task_manager_api/handlers"
	"task_manager_api/model"
)

var db *bun.DB

func main() {

	db, err := connectToDatabase()
	if err != nil {
		fmt.Println("Failed to connect to the db:", err)
		return
	}

	defer db.Close()

	err = createTasksTable(db)
	if err != nil {
		fmt.Println("Failed to create table:", err)
		return
	}

	db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(true),
		bundebug.FromEnv("BUNDEBUG"),
	))

	err = db.Ping()
	if err != nil {
		fmt.Println("Failed to connect to the db")
		return
	}

	fmt.Println("Connected to the database")

	handlers.Db = db

	router := gin.Default()
	router.GET("/ping", ping)

	router.GET("/tasks", handlers.GetTasks)

	router.GET("/tasks/:id", handlers.GetTask)

	router.PUT("/tasks/:id", handlers.EditTask)

	router.POST("/tasks", handlers.PostTask)

	router.Run()
}

func ping(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "pong",
	})
}

func connectToDatabase() (*bun.DB, error) {
	dsn := "postgres://owuor:password@localhost:5432/task_manager?sslmode=disable"
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db = bun.NewDB(sqldb, pgdialect.New())
	return db, nil
}

func createTasksTable(db *bun.DB) error {
	ctx := context.Background()
	_, err := db.NewCreateTable().Model((*model.Task)(nil)).IfNotExists().Exec(ctx)
	return err
}
