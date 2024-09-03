package utils

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

var MonthsMap = map[string]string{
	"1":  "January",
	"2":  "February",
	"3":  "March",
	"4":  "April",
	"5":  "May",
	"6":  "June",
	"7":  "July",
	"8":  "August",
	"9":  "September",
	"10": "October",
	"11": "November",
	"12": "December",
}

const (
	DateFormat = "20060102"
)

func lenMonth(date time.Time) int {
	switch {
	case date.Month().String() == MonthsMap["1"] || date.Month().String() == MonthsMap["3"] || date.Month().String() == MonthsMap["5"] || date.Month().String() == MonthsMap["7"] || date.Month().String() == MonthsMap["8"] || date.Month().String() == MonthsMap["10"] || date.Month().String() == MonthsMap["12"]:
		return 31
	case date.Month().String() == MonthsMap["2"]:
		if (date.Year()%4 == 0 && date.Year()%100 != 0) || date.Year()%400 == 0 {
			return 29
		} else {
			return 28
		}
	case date.Month().String() == MonthsMap["4"] || date.Month().String() == MonthsMap["6"] || date.Month().String() == MonthsMap["9"] || date.Month().String() == MonthsMap["11"]:
		return 30
	}
	return 0
}

func NextDate(now time.Time, date string, repeat string) (string, error) {
	d, err := time.Parse(DateFormat, date)
	if err != nil {
		log.Println(err)
		return "", err
	}

	params := strings.Split(repeat, " ")
	if len(params) == 1 && params[0] != "y" {
		return "", fmt.Errorf("wrong parameters")
	}

	switch params[0] {
	case "y":
		for d.Unix() <= now.Unix() {
			d = d.AddDate(1, 0, 0)
		}

		if d.Format(DateFormat) == date {
			d = d.AddDate(1, 0, 0)
		}

		return d.Format(DateFormat), nil

	case "d":
		param, err := strconv.Atoi(params[1])
		if err != nil {
			return "", err
		} else if param > 400 {
			return "", fmt.Errorf("wrong number: %d\nmax number: 400", param)
		}
		for d.Unix() <= now.Unix() {
			d = d.AddDate(0, 0, param)
		}

		if d.Format(DateFormat) == date {
			d = d.AddDate(0, 0, param)
		}

		return d.Format(DateFormat), nil

	case "w":
		days := strings.Split(params[1], ",")

		for {
			d = d.AddDate(0, 0, 1)
			if d.Unix() > now.Unix() {
				for _, v := range days {
					day, err := strconv.Atoi(v)
					if err != nil {
						return "", err
					} else if day > 7 {
						return "", fmt.Errorf("wrong number: %d\nmax number: 7", day)
					}
					weekDay := int(d.Weekday())
					if weekDay == 0 {
						weekDay += 7
					}
					log.Println((int(d.Weekday())), day)
					if weekDay == day {
						return d.Format(DateFormat), nil
					}
				}
			}
		}

	case "m":
		days := strings.Split(params[1], ",")

		if len(params) == 2 {
			for {
				d = d.AddDate(0, 0, 1)
				if d.Unix() > now.Unix() {
					for _, v := range days {
						day, err := strconv.Atoi(v)
						if err != nil {
							return "", err
						} else if day > 31 || day < -2 {
							return "", fmt.Errorf("wrong number: %d\nmax number: 31\nmin number: -2", day)
						}

						if d.Day() == day || d.Day() == lenMonth(d)+day+1 {
							return d.Format(DateFormat), nil
						}
					}
				}
			}
		} else {
			months := strings.Split(params[2], ",")

			for {
				d = d.AddDate(0, 0, 1)
				if d.Unix() > now.Unix() {
					for _, v1 := range months {
						month, err := strconv.Atoi(v1)
						if err != nil {
							return "", nil
						}

						if month > 12 || month < 1 {
							return "", fmt.Errorf("wrong number: %d\nmax number: 12", month)
						}

						for _, v2 := range days {
							day, err := strconv.Atoi(v2)
							if err != nil {
								return "", err
							} else if day > 31 || day < -2 {
								return "", fmt.Errorf("wrong number: %d\nmax number: 31\nmin number: -2", day)
							}

							if (d.Day() == day || d.Day() == lenMonth(d)+day+1) && int(d.Month()) == month {
								return d.Format(DateFormat), nil
							}
						}
					}
				}
			}
		}

	default:
		return "", fmt.Errorf("wrong repeat format")
	}
}
