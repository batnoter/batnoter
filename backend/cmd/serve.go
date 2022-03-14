package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vivekweb2013/gitnoter/internal/applicationconfig"
	"github.com/vivekweb2013/gitnoter/internal/db"
	"github.com/vivekweb2013/gitnoter/internal/httpservice"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts http server",
	Long:  `Starts the http server with configured options`,
	RunE: func(cmd *cobra.Command, args []string) error {
		db, err := db.Connect(conf.Database)
		if err != nil {
			return err
		}
		applicationconfig := applicationconfig.NewApplicationConfig(conf, db)
		return httpservice.Run(applicationconfig)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
