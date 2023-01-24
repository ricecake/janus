/*
Copyright © 2019 Sebastian Green-Husted <geoffcake@gmail.com>
*/
package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"janus/model"
)

// ContextCmd represents the Context command
var createContextCmd = &cobra.Command{
	Use:   "context",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		context := &model.Context{}

		context.Name = args[0]

		if err := model.CreateContext(context); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	createCmd.AddCommand(createContextCmd)
}
