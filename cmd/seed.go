package cmd

import (
	server "clean/app/http"
	"clean/infra/conn"
	"clean/infra/logger"
	"fmt"

	"github.com/spf13/cobra"
)

var truncate bool

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Prepopulate data",
	Long:  `prepopulate roles, permissions, role_premissions data into mysql`,
	Run:   seed,
}

func seed(cmd *cobra.Command, args []string) {
	// seed roles, permissions, role_permissions
	//conn.ConnectDbMysql()
	conn.ConnectDbSqlite()

	db := conn.Db()
	truncate, _ = cmd.Flags().GetBool("truncate")
	fmt.Println("truncate=", truncate)

	for _, seed := range conn.SeedAll() {
		if err := seed.Run(db, truncate); err != nil {
			logger.Error(fmt.Sprintf("Running seed '%s', failed with error:", seed.Name), err)
		}
	}

	server.Start()
}

func init() {
	seedCmd.PersistentFlags().BoolVarP(&truncate, "truncate", "t", false, "will truncate tables")
}
