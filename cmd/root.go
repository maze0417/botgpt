package cmd

import (
	"botgpt/cmd/botgpt"
	"botgpt/internal/config"
	"fmt"
	"github.com/gin-gonic/gin"
	"os"

	"github.com/spf13/cobra"
)

var (
	env      string
	service  string
	version  string
	serverID string
)

var rootCmd = &cobra.Command{
	Use:   "server",
	Short: "A server with different modes and services",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("cli argu :: environment %s, service %s version %s serverID %s \n",
			env, service, version, serverID)

		gin.SetMode(gin.ReleaseMode)

		switch service {
		case "botgpt":
			botgpt.RegisterFactory()
			config.Init(env, service, version, serverID)
			botgpt.Run()
		default:
			panic("no match service : " + service)
		}
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&env, "env", "e", "local", "Environment mode (production, uat, development, local)")
	rootCmd.PersistentFlags().StringVarP(&service, "service", "s", "botgpt", "Service type (api, job)")
	rootCmd.PersistentFlags().StringVarP(&version, "version", "v", "1.0.0", "Version")
	rootCmd.PersistentFlags().StringVarP(&serverID, "serverID", "i", "1", "Server ID")
}

// Execute is the entry point for the Cobra CLI.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
