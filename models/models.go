package models

import "gorm.io/gorm"

type Task struct {
	gorm.Model `json:"-"`
	ID         int    `json:"id" gorm:"<-:create"`
	Date       string `json:"date"`
	Title      string `json:"title"`
	Comment    string `json:"comment"`
	Repeat     string `json:"repeat"`
}

type CorrectResponse struct {
	ID string `json:"id"`
}

type CorrectGetTasks struct {
	Tasks []Task `json:"tasks"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type Tabler interface {
	TableName() string
}

func (Task) TableName() string {
	return "scheduler"
}
