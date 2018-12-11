package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

type CreateDomainOptions struct{}

var cdo = &CreateDomainOptions{}

var createDomainCmd = &cobra.Command{
	Use:   "domain",
	Short: "Create external domain",
	Long:  `Create secured external domain to access your application with one command`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := RunCreateDomain(cdo); err != nil {
			fmt.Printf("RunCreateDomain error: %s", err.Error())
		}
	},
}

func RunCreateDomain(cdo *CreateDomainOptions) error {
	fmt.Println("RunCreateDomain not yet implemented")
	return nil
}

func init() {
	createCmd.AddCommand(createDomainCmd)
}
