package database

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/robertvitoriano/go-encoder-microservice/domain"
)

type Dabase struct {
	Db            *gorm.DB
	Dsn           string
	DsnTest       string
	DbType        string
	DbTypeTest    string
	Debug         bool
	AutoMigrateDb bool
	Env           string
}

func NewDb() *Dabase {
	return &Dabase{}
}

func NewDbTest() *gorm.DB {
	dbInstance := NewDb()
	dbInstance.Env = "Test"
	dbInstance.DbTypeTest = "sqlite3"
	dbInstance.DsnTest = ":memory"
	dbInstance.AutoMigrateDb = true
	dbInstance.Debug = true

	connection, err := dbInstance.Connect()

	if err != nil {
		log.Fatalf("Test db error: %v", err)
	}
	return connection
}

func (database *Dabase) Connect() (*gorm.DB, error) {
	var err error

	if database.Env != "Test" {
		database.Db, err = gorm.Open(database.DbType, database.Dsn)

		if err != nil {
			return nil, err
		}

		return database.Db, nil
	}

	database.Db, err = gorm.Open(database.DbTypeTest, database.DsnTest)

	if err != nil {
		return nil, err
	}

	if database.Debug {
		database.Db.LogMode(true)
	}

	if database.AutoMigrateDb {
		database.Db.AutoMigrate(&domain.Video{}, &domain.Job{})
		database.Db.Model(domain.Job{}).AddForeignKey("video_id", "videos (id)", "CASCADE", "CASCADE")
	}

	return database.Db, nil
}
