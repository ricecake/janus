/*
Copyright © 2019 Sebastian Green-Husted <geoffcake@gmail.com>
*/
package cmd

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"janus/model"
)

// userCmd represents the user command
var createUserCmd = &cobra.Command{
	Use:   "user",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		user := &model.Identity{}

		fmt.Printf("create user [%s] password [%s]\n", args[0], args[1])

		active, err := cmd.Flags().GetBool("active")
		if err != nil {
			log.Fatal(err)
		}

		pname, err := cmd.Flags().GetString("preferred-name")
		if err != nil {
			log.Fatal(err)
		}

		user.Email = args[0]
		user.Active = active
		user.PreferredName = pname

		if err := model.CreateIdentity(user); err != nil {
			log.Fatal(err)
		}

		user.SetPassword(args[1])
	},
}

func init() {
	createCmd.AddCommand(createUserCmd)

	createUserCmd.Flags().BoolP("active", "a", true, "User is active")
	createUserCmd.Flags().StringP("preferred-name", "p", "", "Preferred name for user")
}
