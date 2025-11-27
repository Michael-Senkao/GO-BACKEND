package data

import (
	"errors"
	"time"

	"task_manager/models"
)

var tasks = []models.Task{
	{ID: "1", Title: "Task 1", Description: "First task", DueDate: time.Now(), Status: "Pending"},
	{ID: "2", Title: "Task 2", Description: "Second task", DueDate: time.Now().AddDate(0, 0, 1), Status: "In Progress"},
	{ID: "3", Title: "Task 3", Description: "Third task", DueDate: time.Now().AddDate(0, 0, 2), Status: "Completed"},
}

func GetAllTasks() []models.Task {
	return tasks
}

func GetTaskByID(id string) (*models.Task, error) {
	for _, t := range tasks {
		if t.ID == id {
			return &t, nil
		}
	}
	return nil, errors.New("task not found")
}

func AddTask(task models.Task) {
	tasks = append(tasks, task)
}

func UpdateTask(id string, updated models.Task) error {
	for i, t := range tasks {
		if t.ID == id {
			if updated.Title != "" {
				tasks[i].Title = updated.Title
			}
			if updated.Description != "" {
				tasks[i].Description = updated.Description
			}
			if updated.Status != "" {
				tasks[i].Status = updated.Status
			}
			return nil
		}
	}
	return errors.New("task not found")
}

func DeleteTask(id string) error {
	for i, t := range tasks {
		if t.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			return nil
		}
	}
	return errors.New("task not found")
}
