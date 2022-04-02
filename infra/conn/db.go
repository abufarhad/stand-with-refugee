package conn

import (
	"clean/app/domain"
	"clean/infra/config"
	"clean/infra/logger"
	"encoding/json"
	"fmt"
	"gorm.io/driver/mysql"

	"gorm.io/driver/sqlite" // Sqlite driver based on GGO
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"io/ioutil"
	"os"
	"time"
)

var db *gorm.DB

func ConnectDbMysql() {
	conf := config.Db()

	logger.Info("connecting to mysql at " + conf.Host + ":" + conf.Port + "...")

	logMode := gormlogger.Silent
	if conf.Debug {
		logMode = gormlogger.Info
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", conf.User, conf.Pass, conf.Host, conf.Port, conf.Schema)

	dB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		PrepareStmt: true,
		Logger:      gormlogger.Default.LogMode(logMode),
	})

	if err != nil {
		panic(err)
	}

	sqlDb, err := dB.DB()
	if err != nil {
		panic(err)
	}

	if conf.MaxIdleConn != 0 {
		sqlDb.SetMaxIdleConns(conf.MaxIdleConn)
	}
	if conf.MaxOpenConn != 0 {
		sqlDb.SetMaxOpenConns(conf.MaxOpenConn)
	}
	if conf.MaxConnLifetime != 0 {
		sqlDb.SetConnMaxLifetime(conf.MaxConnLifetime * time.Second)
	}

	db = dB

	db.AutoMigrate(
		&domain.User{},
		&domain.Role{},
		&domain.Permission{},
		&domain.RolePermission{},
		&domain.Specialization{},
		&domain.Commitments{},
		&domain.Place{},
		&domain.Symptom{},
		&domain.Help{},
	)
	logger.Info("mysql connection successful...")
}

func ConnectDbSqlite() {
	conf := config.Db()

	logger.Info("connecting to sqlite ...")

	dB, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	sqlDb, err := dB.DB()
	if err != nil {
		panic(err)
	}

	if conf.MaxIdleConn != 0 {
		sqlDb.SetMaxIdleConns(conf.MaxIdleConn)
	}
	if conf.MaxOpenConn != 0 {
		sqlDb.SetMaxOpenConns(conf.MaxOpenConn)
	}
	if conf.MaxConnLifetime != 0 {
		sqlDb.SetConnMaxLifetime(conf.MaxConnLifetime * time.Second)
	}

	db = dB

	db.AutoMigrate(
		&domain.User{},
		&domain.Role{},
		&domain.Permission{},
		&domain.RolePermission{},
		&domain.Specialization{},
		&domain.Commitments{},
		&domain.Place{},
		&domain.Symptom{},
		&domain.Help{},
	)

	logger.Info("mysql connection successful...")
}

func Db() *gorm.DB {
	return db
}

type Seed struct {
	Name string
	Run  func(db *gorm.DB, truncate bool) error
}

func SeedAll() []Seed {
	return []Seed{
		{
			Name: "Place",
			Run: func(db *gorm.DB, truncate bool) error {
				if err := seedPlaces(db, "/infra/seed/place.json", truncate); err != nil {
					return err
				}
				return nil
			},
		},
		{
			Name: "Specialization",
			Run: func(db *gorm.DB, truncate bool) error {
				if err := seedSpecialization(db, "/infra/seed/specialization.json", truncate); err != nil {
					return err
				}
				return nil
			},
		},
		{
			Name: "Symptoms",
			Run: func(db *gorm.DB, truncate bool) error {
				if err := seedSymptoms(db, "/infra/seed/symptoms.json", truncate); err != nil {
					return err
				}
				return nil
			},
		},
	}
}

func seedPlaces(db *gorm.DB, jsonfilPath string, truncate bool) error {
	file, _ := readSeedFile(jsonfilPath)
	places := []domain.Place{}

	_ = json.Unmarshal([]byte(file), &places)

	if truncate {
		db.Exec("DELETE TABLE refugee.places;")
		db.Exec("DELETE TABLE refugee.specializations;")
		db.Exec("DELETE TABLE refugee.symptoms;")
	}

	var count int64

	db.Model(&domain.Place{}).Count(&count)
	if count == 0 {
		db.Create(&places)
	}

	return nil
}

func seedSpecialization(db *gorm.DB, jsonfilPath string, truncate bool) error {
	file, _ := readSeedFile(jsonfilPath)
	spcial := []domain.Specialization{}

	_ = json.Unmarshal([]byte(file), &spcial)

	var count int64

	db.Model(&domain.Specialization{}).Count(&count)
	if count == 0 {
		db.Create(&spcial)
	}

	return nil
}

func seedSymptoms(db *gorm.DB, jsonfilPath string, truncate bool) error {
	file, _ := readSeedFile(jsonfilPath)
	symp := []domain.Symptom{}

	_ = json.Unmarshal([]byte(file), &symp)

	var count int64

	db.Model(&domain.Symptom{}).Count(&count)
	if count == 0 {
		db.Create(&symp)
	}

	return nil
}

func readSeedFile(jsonfilPath string) ([]byte, error) {
	BaseDir, _ := os.Getwd()
	seedFile := BaseDir + jsonfilPath
	if BaseDir == "/" {
		seedFile = jsonfilPath
	}
	fmt.Println("seed folder: ", seedFile)

	return ioutil.ReadFile(seedFile)
}
