/*
Copyright © 2019 Sebastian Green-Husted <geoffcake@gmail.com>

*/
package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"janus/model"
)

// clientCmd represents the client command
var createClientCmd = &cobra.Command{
	Use:   "client",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		client := &model.Client{}

		dname, err := cmd.Flags().GetString("display-name")
		if err != nil {
			log.Fatal(err)
		}

		redir, err := cmd.Flags().GetString("redirect")
		if err != nil {
			log.Fatal(err)
		}

		client.Context = args[0]
		client.DisplayName = dname
		client.BaseUri = redir

		client.SetSecret(args[1])

		if err := model.CreateClient(client); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	createCmd.AddCommand(createClientCmd)

	createClientCmd.Flags().StringP("display-name", "d", "", "display name for user")
	createClientCmd.Flags().StringP("redirect", "r", "", "redirect url")
}
