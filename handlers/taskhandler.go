package handlers

import (
	
	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
	"net/http"
	"task_manager_api/model"
)

var Db *bun.DB

func HomePage(ctx *gin.Context){
	ctx.JSON(http.StatusOK, gin.H{"message": "Welcome to the taskmanager API"})
}

func GetTasks(ctx *gin.Context) {
	var tasks []model.Task
	err := Db.NewSelect().Model(&tasks).Scan(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

func GetTask(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	task := &model.Task{}
	err := Db.NewSelect().Model(task).Where("id = ?", id).Scan(ctx.Request.Context())

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	if task.ID == "" {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	ctx.JSON(http.StatusOK, task)
}

func EditTask(ctx *gin.Context) {
	id := ctx.Param("id")
	updatedTask := &model.Task{}

	if err := ctx.ShouldBindJSON(&updatedTask); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := Db.NewUpdate().Model(updatedTask).
		Set("title = ?", updatedTask.Title).
		Set("description = ?", updatedTask.Description).
		Where("id = ?", id).
		Exec(ctx.Request.Context())

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Task updated"})
}

func PostTask(ctx *gin.Context) {
	newTask := &model.Task{}
	if err := ctx.ShouldBindJSON(&newTask); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := Db.NewInsert().Model(newTask).Exec(ctx.Request.Context())

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"message": "Task created"})

}
