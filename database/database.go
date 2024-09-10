package database

import (
	_ "embed"
	"log"
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite"
)

type GormInstance struct {
	Db *gorm.DB
}

var Gorm GormInstance

//go:embed scheduler.sql
var scheduler string

func ConnectGorm() {
	appPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	dbFile := filepath.Join(filepath.Dir(appPath), "scheduler.db")
	_, err = os.Stat(dbFile)
	if err != nil {
		log.Println(err)
	}

	var install bool
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("База данных будет создана")
			install = true
		} else {
			log.Println("Не получилось проверить файл")
			log.Fatal(err)
		}
	}

	Gorm.Db, err = gorm.Open(sqlite.Open("scheduler.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err, "connect")
		return
	}

	if install {
		Gorm.Db.Exec(scheduler)
	}
}
