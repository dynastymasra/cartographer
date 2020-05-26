package console

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	j "github.com/neo4j/neo4j-go-driver/neo4j"

	"github.com/sirupsen/logrus"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/neo4j"

	_ "github.com/golang-migrate/migrate/v4/database/neo4j"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	migrationSourcePath = "file://migration"
	migrationFilePath   = "./migration"
)

func CreateMigrationFiles(filename string) error {
	if len(filename) == 0 {
		return errors.New("migration filename is not provided")
	}

	timeStamp := time.Now().Unix()
	upMigrationFilePath := fmt.Sprintf("%s/%d_%s.up.cypher", migrationFilePath, timeStamp, filename)
	downMigrationFilePath := fmt.Sprintf("%s/%d_%s.down.cypher", migrationFilePath, timeStamp, filename)

	if err := createFile(upMigrationFilePath); err != nil {
		return err
	}
	log.Println("created", upMigrationFilePath)

	if err := createFile(downMigrationFilePath); err != nil {
		os.Remove(upMigrationFilePath)
		return err
	}

	log.Println("created", downMigrationFilePath)

	return nil
}

func Migration(client j.Driver) (*migrate.Migrate, error) {
	config := &neo4j.Config{
		MigrationsLabel: neo4j.DefaultMigrationsLabel,
		MultiStatement:  true,
	}

	driver, err := neo4j.WithInstance(client, config)
	if err != nil {
		logrus.WithError(err).Errorln("Failed open instance")

		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance(migrationSourcePath, "neo4j", driver)
	if err != nil {
		logrus.WithError(err).Errorln("Failed migration data")

		return nil, err
	}

	return m, nil
}

func RunMigration(migration *migrate.Migrate) error {
	if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
		logrus.WithError(err).Errorln("Failed run database migration")
		return err
	}
	return nil
}

func RollbackMigration(migration *migrate.Migrate) error {
	if err := migration.Steps(-1); err != nil {
		logrus.WithError(err).Errorln("Failed rollback database migration")
		return err
	}
	return nil
}

func createFile(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}

	return nil
}
