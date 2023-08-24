/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// dryRunCmd represents the dryRun command
var dryRunCmd = &cobra.Command{
	Use:   "dryRun",
	Short: "Outputs the migrations that will be performed without actually executing them.",
	Long:  `Outputs the migrations that will be performed without actually executing them.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("dryRun called")
	},
}

func init() {
	rootCmd.AddCommand(dryRunCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// dryRunCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// dryRunCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
