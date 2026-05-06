package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prachi-satbhai0741/Taskboard/db"
	"github.com/prachi-satbhai0741/Taskboard/models"
)

func GetTasks(c *gin.Context) {
	var tasks []models.Task
	db.DB.Find(&tasks)
	c.JSON(http.StatusOK, tasks)
}

func CreateTask(c *gin.Context) {
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db.DB.Create(&task)
	c.JSON(http.StatusCreated, task)
}

func UpdateTask(c *gin.Context) {
	id := c.Param("id")
	var task models.Task

	if err := db.DB.First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	var input struct {
		Status models.Status `json:"status"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.DB.Model(&task).Update("status", input.Status)
	c.JSON(http.StatusOK, task)
}

func DeleteTask(c *gin.Context) {
	id := c.Param("id")
	var task models.Task

	if err := db.DB.First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	db.DB.Delete(&task)
	c.JSON(http.StatusOK, gin.H{"message": "Task deleted"})
}
