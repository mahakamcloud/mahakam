package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

type LoginOptions struct{}

var lo = &LoginOptions{}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login with your gate credentials",
	Long:  `Login with your gate credentials`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := RunLogin(lo); err != nil {
			fmt.Printf("RunLogin error: %s", err.Error())
		}
	},
}

func RunLogin(lo *LoginOptions) error {
	fmt.Println("RunLogin not yet implemented")
	return nil
}

func init() {
	RootCmd.AddCommand(loginCmd)
}
