package database

import (
	"log"

	"github.com/Codeflix-FullCycle/encoder/domain"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/lib/pq"
)

type Database struct {
	Db            *gorm.DB
	Dns           string
	DnsTest       string
	DbType        string
	DbTypeTest    string
	Debbug        bool
	AutoMigrateDb bool
	Env           string
}

func NewDb() *Database {
	return &Database{}
}

func NewDbTest() *gorm.DB {
	dbInstance := NewDb()
	dbInstance.Env = "test"
	dbInstance.DbTypeTest = "sqlite3"
	dbInstance.DnsTest = ":memory:"
	dbInstance.AutoMigrateDb = true
	dbInstance.Debbug = true

	connection, err := dbInstance.Connect()

	if err != nil {
		log.Fatalf("Test db error: %v", err)
	}

	return connection
}

func (d *Database) Connect() (*gorm.DB, error) {
	var err error

	if d.Env != "test" {
		d.Db, err = gorm.Open(d.DbType, d.Dns)
	} else {
		d.Db, err = gorm.Open(d.DbTypeTest, d.DnsTest)
	}

	if err != nil {
		return nil, err
	}

	if d.Debbug {
		d.Db.LogMode(true)
	}

	if d.AutoMigrateDb {
		d.Db.AutoMigrate(&domain.Video{}, &domain.Job{})
		d.Db.Model(&domain.Job{}).AddForeignKey("video_id", "video(id)", "CASCADE", "CASACADE")
	}

	return d.Db, nil
}