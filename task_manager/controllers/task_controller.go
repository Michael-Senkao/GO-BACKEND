package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"task_manager/data"
	"task_manager/models"
)

func GetTasks(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, data.GetAllTasks())
}

func GetTaskByID(ctx *gin.Context) {
	id := ctx.Param("id")
	task, err := data.GetTaskByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, task)
}

func AddTask(ctx *gin.Context) {
	var task models.Task
	if err := ctx.ShouldBindJSON(&task); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	data.AddTask(task)
	ctx.JSON(http.StatusCreated, gin.H{"message": "task created"})
}

func UpdateTask(ctx *gin.Context) {
	id := ctx.Param("id")
	var task models.Task
	if err := ctx.ShouldBindJSON(&task); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := data.UpdateTask(id, task); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "task updated"})
}

func DeleteTask(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := data.DeleteTask(id); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "task removed"})
}
