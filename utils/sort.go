package utils

import (
	"time"

	"github.com/stxreocoma/todo/models"
)

func Sort(tasks []models.Task) {
	for i := 0; i < len(tasks)-1; i++ {
		for j := i + 1; j < len(tasks); j++ {
			aDate, _ := time.Parse(DateFormat, tasks[i].Date)
			bDate, _ := time.Parse(DateFormat, tasks[j].Date)
			if aDate.Unix() > bDate.Unix() {
				tasks[i], tasks[j] = tasks[j], tasks[i]
			}
		}
	}
}
