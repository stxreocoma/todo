package models

import "fmt"

var (
	ErrRepeat       = fmt.Errorf("wrong repeat format")
	ErrTaskNotFound = fmt.Errorf("задача не найдена")
	ErrSearch       = fmt.Errorf("wrong search format")
	ErrID           = fmt.Errorf("не указан идентификатор")
	ErrTitle        = fmt.Errorf("no title")
)
