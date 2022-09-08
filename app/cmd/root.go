package cmd

import (
	"database/sql"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	domain "github.com/LieAlbertTriAdrian/clean-arch-golang/domain"
	"github.com/LieAlbertTriAdrian/clean-arch-golang/app/config"
	"github.com/LieAlbertTriAdrian/clean-arch-golang/ebus"
	eBusRabbit "github.com/LieAlbertTriAdrian/clean-arch-golang/ebus/rabbitmq"
	eBusSNS "github.com/LieAlbertTriAdrian/clean-arch-golang/ebus/sns"

	_todoUsecase "github.com/LieAlbertTriAdrian/clean-arch-golang/todo/usecase"
)

var (
	EnvFilePath string
	rootCmd     = &cobra.Command{
		Use:   "template",
		Short: "Template is template management application",
	}
)

var (
	rootConfig  config.Root
	database    *sql.DB
	eBus        *ebus.Bus
	todoUsecase domain.ITodoUsecase
)

// Execute will call the root command execute
func Execute() {
	rootCmd.PersistentFlags().StringVarP(&EnvFilePath, "env", "e", ".env", ".env file to read from")
	if err := rootCmd.Execute(); err != nil {
		logrus.Fatal(err)
	}
}

func init() {
	cobra.OnInitialize(initConfigReader, initLogLevel, initDB, initApp)
}

func initLogLevel() {
	if rootConfig.App.Env != "production" {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Warn("Template is Running in Debug Mode")
		return
	}
	logrus.SetLevel(logrus.InfoLevel)
	logrus.Warn("Template is Running in Production Mode")
}

func initConfigReader() {
	rootConfig = config.Load(EnvFilePath)
}

func initDB() {
	var err error
	database, err = config.OpenDatabaseConnection(rootConfig.Postgres)
	if err != nil {
		logrus.WithError(err).Fatal("Could not establish connection to database")
	}
}

func initApp() {
	/*****************
	 * Ebus Handler
	 *****************/
	eBus = ebus.NewEbus()
	// Ebus handler for RabbitMQ
	rabbitEbus := eBusRabbit.NewEbusSubscriber(nil) //TODO(LieAlbertTriAdrian): fill the RabbitMQ Client
	eBus.Subscribe(rabbitEbus)
	// Ebus handler for SNS
	SNSEbus := eBusSNS.NewEbusSubscriber(nil) //TODO(LieAlbertTriAdrian): fill the SNS Client
	eBus.Subscribe(SNSEbus)

	/*****************
	 * Todo Service
	 *****************/

	todoUsecase = _todoUsecase.NewTxService(database)
}
