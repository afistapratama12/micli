package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "micli",
	Short: "Market information-cli",
	Long: `Market information-cli is a CLI application for retrieving market information.

The application retrieves market information from a variety of sources and displays it to the user.	
	`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	// configFile = ".cobra-cli-samples.yml"
	// viper.SetConfigType("yaml")
	// viper.SetConfigFile(configFile)

	// viper.AutomaticEnv()
	// viper.SetEnvPrefix("COBRACLISAMPLES")
	// helper.HandleError(viper.BindEnv("API_KEY"))
	// helper.HandleError(viper.BindEnv("API_SECRET"))
	// helper.HandleError(viper.BindEnv("USERNAME"))
	// helper.HandleError(viper.BindEnv("PASSWORD"))

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using configuration file: ", viper.ConfigFileUsed())
	}
}
