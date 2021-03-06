package cmd

import (
	"clean/infra/config"
	"clean/infra/conn"
	"clean/infra/logger"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	RootCmd = &cobra.Command{
		Use:   "clean code",
		Short: "implementing clean architecture in golang",
	}
)

func init() {
	RootCmd.AddCommand(serveCmd)
	RootCmd.AddCommand(seedCmd)
}

// Execute executes the root command
func Execute() {
	config.LoadConfig()
	//conn.ConnectDbMysql()
	conn.ConnectDbSqlite()
	//conn.ConnectRedis()

	logger.Info("about to start the application")

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
