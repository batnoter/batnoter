package cmd

import (
	"os"

	"github.com/iamolegga/enviper"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vivekweb2013/gitnoter/internal/config"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gitnoter",
	Short: "A simple note taking app",
	Long:  `A simple web application to store and retrieve notes`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

var conf config.Config

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	initLogger()
	initConfig()

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initLogger() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}

func initConfig() {
	// The enviper is a wrapper over viper that loads the config from file & overrides
	// them with env variables if available. If the config file is missing it simply ignores it
	e := enviper.New(viper.New())

	var cfgFile string
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is .gitnoter.yaml)")
	if cfgFile != "" {
		e.SetConfigFile(cfgFile)
	} else {
		e.AddConfigPath(".")
		e.SetConfigName(".gitnoter")
	}

	if err := e.Unmarshal(&conf); err == nil {
		logrus.Infof("using the config file: %s", e.ConfigFileUsed())
	}
}
