package validation

import (
	"strconv"
	"strings"
	"time"

	"github.com/stxreocoma/todo/models"
	"github.com/stxreocoma/todo/utils"
)

func Repeat(repeat string) error {
	if len(repeat) == 0 || repeat == "y" {
		return nil
	}

	parts := strings.Split(repeat, " ")

	switch parts[0] {
	case "d":
		days := strings.Split(parts[1], ",")
		for _, v := range days {
			day, err := strconv.Atoi(v)
			if err != nil {
				return models.ErrRepeat
			} else if day < 1 || day > 400 {
				return models.ErrRepeat
			}
		}
	case "w":
		days := strings.Split(parts[1], ",")
		for _, v := range days {
			day, err := strconv.Atoi(v)
			if err != nil {
				return models.ErrRepeat
			} else if day < 1 || day > 7 {
				return models.ErrRepeat
			}
		}
	case "m":
		days := strings.Split(parts[1], ",")
		for _, v := range days {
			day, err := strconv.Atoi(v)
			if err != nil {
				return models.ErrRepeat
			} else if day < -2 || day > 31 {
				return models.ErrRepeat
			}
		}
		if len(repeat) == 2 {
			months := strings.Split(parts[2], ",")
			for _, month := range months {
				if _, ok := utils.MonthsMap[month]; !ok {
					return models.ErrRepeat
				}
			}
		}
	default:
		return models.ErrRepeat
	}

	return nil
}

func Date(task models.Task) (string, error) {
	date, err := time.Parse(utils.DateFormat, task.Date)
	if err != nil {
		if task.Date == "" {
			date = time.Now()
		} else {
			return "", err
		}
	} else if date.Unix() < time.Now().Unix() && (time.Now().Unix()-date.Unix() >= 86400 || time.Now().Day() != date.Day()) {
		if task.Repeat == "" {
			date = time.Now()
		} else {
			dateString, err := utils.NextDate(time.Now(), date.Format(utils.DateFormat), task.Repeat)
			if err != nil {
				return "", err
			}
			date, err = time.Parse(utils.DateFormat, dateString)
			if err != nil {
				return "", err
			}
		}
	}

	return date.Format(utils.DateFormat), nil
}
