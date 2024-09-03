package models

import "gorm.io/gorm"

type Task struct {
	gorm.Model `json:"-"`
	ID         string `json:"id" gorm:"<-:create"`
	Date       string `json:"date"`
	Title      string `json:"title"`
	Comment    string `json:"comment"`
	Repeat     string `json:"repeat"`
}

type Tabler interface {
	TableName() string
}

func (Task) TableName() string {
	return "scheduler"
}
