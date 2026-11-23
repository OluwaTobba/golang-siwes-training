package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "devtool",
	Short: "DevOps utility CLI",
	Long:  "devtool - a collection of network and HTTP inspection utilities.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().String("config", "", "config file (default: $HOME/.devtool.yaml)")
	rootCmd.PersistentFlags().Bool("verbose", false, "enable verbose output")
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
}

func initConfig() {
	viper.SetConfigName(".devtool")
	viper.SetConfigType("yaml")
	if home, err := os.UserHomeDir(); err == nil {
		viper.AddConfigPath(home)
	}
	viper.SetEnvPrefix("devtool") // DEVTOOL_VERBOSE=true -> viper.Get("verbose")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()
	viper.ReadInConfig()
}
