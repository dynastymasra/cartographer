package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/dynastymasra/cartographer/country"
	"github.com/dynastymasra/cartographer/infrastructure/web"
	"github.com/dynastymasra/cartographer/region"
	"gopkg.in/tylerb/graceful.v1"

	"github.com/dynastymasra/cartographer/config"
	"github.com/dynastymasra/cartographer/console"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func init() {
	config.Load()
	config.Logger().Setup()

}

func main() {
	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	log := logrus.WithFields(logrus.Fields{
		"service_name": config.ServiceName,
		"version":      config.Version,
	})

	log.Infoln("Prepare start service")

	driver, err := config.Neo4J().Driver()
	if err != nil {
		log.WithError(err).Fatalln("Failed create neo4j driver")
	}

	migration, err := console.Migration(driver)
	if err != nil {
		log.WithError(err).Fatalln("Failed run migration")
	}

	regionRepo := region.NewRepository(driver)
	countryRepo := country.NewRepository(driver)

	clientApp := cli.NewApp()
	clientApp.Name = config.ServiceName
	clientApp.Version = config.Version

	clientApp.Action = func(c *cli.Context) error {
		webServer := &graceful.Server{
			Timeout: 0,
		}

		router := web.NewRouter(config.ServiceName, driver, regionRepo, countryRepo)

		go web.Run(webServer, config.ServerAddress(), router)

		select {
		case sig := <-stop:
			<-webServer.StopChan()

			log.Warnln(fmt.Sprintf("Service shutdown because %+v", sig))
			os.Exit(0)
		}

		return nil
	}

	clientApp.Commands = []*cli.Command{
		{
			Name:        "migrate:run",
			Description: "Running database migration",
			Action: func(c *cli.Context) error {
				logrus.Infoln("Start database migration")

				if err := console.RunMigration(migration); err != nil {
					logrus.WithError(err).Errorln("Failed run database migration")
					os.Exit(1)
				}

				logrus.Infoln("Success run database migration to latest")

				return nil
			},
		}, {
			Name:        "migrate:rollback",
			Description: "Rollback database migration",
			Action: func(c *cli.Context) error {
				logrus.Infoln("Rollback database migration to previous version")

				if err := console.RollbackMigration(migration); err != nil {
					logrus.WithError(err).Errorln("Failed rollback database migration")
					os.Exit(1)
				}

				logrus.Infoln("Success rollback database migration")

				return nil
			},
		}, {
			Name:        "migrate:create",
			Description: "Create up and down migration files with timestamp",
			Action: func(c *cli.Context) error {
				return console.CreateMigrationFiles(c.Args().Get(0))
			},
		},
	}

	if err := clientApp.Run(os.Args); err != nil {
		panic(err)
	}
}
