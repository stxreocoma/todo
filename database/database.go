package database

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite"
)

type DBinstance struct {
	Db *sql.DB
}

var DB DBinstance

type GormInstance struct {
	Db *gorm.DB
}

var Gorm GormInstance

func ConnectDB() {
	appPath, err := os.Executable()
	log.Println(appPath)
	if err != nil {
		log.Fatal(err)
	}
	dbFile := filepath.Join(filepath.Dir(appPath), "scheduler.db")
	_, err = os.Stat(dbFile)
	log.Println(err.Error())

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

	DB.Db, err = sql.Open("sqlite", "scheduler.db")
	if err != nil {
		log.Fatal(err, "connect")
		return
	}

	if install {
		_, err = DB.Db.Exec(`CREATE TABLE scheduler (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date CHAR(8),
			title TEXT,
			comment TEXT,
			repeat VARCHAR(128)
		);
		CREATE INDEX scheduler_date ON scheduler (date);`)
		log.Println(err.Error(), "con")
	}
}

func ConnectGorm() {
	appPath, err := os.Executable()
	log.Println(appPath)
	if err != nil {
		log.Fatal(err)
	}
	dbFile := filepath.Join(filepath.Dir(appPath), "scheduler.db")
	_, err = os.Stat(dbFile)
	log.Println(err.Error())

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
		Gorm.Db.Exec(`CREATE TABLE scheduler (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date CHAR(8),
			title TEXT,
			comment TEXT,
			repeat VARCHAR(128)
		);
		CREATE INDEX scheduler_date ON scheduler (date);`)
	}
}
