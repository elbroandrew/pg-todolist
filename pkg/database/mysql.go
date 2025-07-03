package database

import (
	"context"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type MySQLDB struct {
	db *gorm.DB
	cfg *Config
}

func InitMySQL(cfg *Config) Database {

	return &MySQLDB{cfg: cfg}
}

func (m *MySQLDB) Connect(ctx context.Context) error {


	db, err := gorm.Open(mysql.Open(m.cfg.MySQLDSN()), &gorm.Config{
		Logger:  logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time { return time.Now().UTC() },
		PrepareStmt: true,
	})
	if err != nil {
		log.Fatalf("ERROR CONNECTION to MYSQL: %v", err)
	}

	log.Printf("SUCCESS CONNECTION! MYSQL")

	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(m.cfg.MaxConns)

	m.db = db
	return nil

}

func (m *MySQLDB) Migrate(models ...interface{}) error {
	return m.db.AutoMigrate(models...)
}

func (m *MySQLDB) GetDB() *gorm.DB {
	return m.db
}

func (m *MySQLDB) Close() error {
	sqlDB, err := m.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (m *MySQLDB) HealthCheck(ctx context.Context) error {
	sqlDB, err := m.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.PingContext(ctx)
}
