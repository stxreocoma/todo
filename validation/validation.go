package validation

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/stxreocoma/utils"
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
				return fmt.Errorf("wrong repeat format")
			} else if day < 1 || day > 400 {
				return fmt.Errorf("wrong repeat format")
			}
		}
	case "w":
		days := strings.Split(parts[1], ",")
		for _, v := range days {
			day, err := strconv.Atoi(v)
			if err != nil {
				return fmt.Errorf("wrong repeat format")
			} else if day < 1 || day > 7 {
				return fmt.Errorf("wrong repeat format")
			}
		}
	case "m":
		days := strings.Split(parts[1], ",")
		for _, v := range days {
			day, err := strconv.Atoi(v)
			if err != nil {
				return fmt.Errorf("wrong repeat format")
			} else if day < -2 || day > 31 {
				return fmt.Errorf("wrong repeat format")
			}
		}
		if len(repeat) == 2 {
			months := strings.Split(parts[2], ",")
			for _, month := range months {
				if _, ok := utils.MonthsMap[month]; !ok {
					return fmt.Errorf("wrong repeat format")
				}
			}
		}
	default:
		return fmt.Errorf("wrong repeat format")
	}
}
