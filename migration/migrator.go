package migration

import (
	"botgpt/internal/config"
	"botgpt/internal/repository"

	"fmt"
	"github.com/golang-migrate/migrate/v4"

	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"log"
)

type MigrateLogger struct {
}

func (m *MigrateLogger) Printf(format string, v ...interface{}) {
	log.Printf(format, v)
}
func (m *MigrateLogger) Verbose() bool {
	return true
}

func Migrate() {

	repository.CreateDbIfNotExist()

	c := config.GetConfig()
	migrationFile := c.GetString("mysql.migration")

	migrateDsn := fmt.Sprintf("mysql://%s", repository.GetMysqlDsn())

	log.Printf("migrate used %v %v", migrationFile, migrateDsn)
	m, err := migrate.New(
		migrationFile,
		migrateDsn,
	)
	if err != nil {
		log.Fatalf("migrate open mysql fatal %v", err)
	}
	m.Log = &MigrateLogger{}

	if err := m.Up(); err == migrate.ErrNoChange {
		log.Printf("migrate %v", err)
	}
	if err != nil {
		log.Fatalf("migrate up fatal %v", err)
	}

	fmt.Println("migration completed")
}
