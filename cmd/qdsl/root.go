// Copyright 2023 NJWS Inc.
// Copyright 2022 Listware

package main

import (
	"errors"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var log = logrus.New()

var (
	version         = "v1.0.0"
	release         = "dev"
	versionTemplate = `{{printf "%s" .Short}}
{{printf "Version: %s" .Version}}
Release: ` + release + `
`
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Version: version,
	Use:     "qdsl QUERY",
	Short:   "CMDB query language",
	Long:    `CMDB query language for getting information about nodes`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires at least one arg")
		}

		return nil
	},
	Run: qdslQuery,
}

func init() {
	rootCmd.SetVersionTemplate(versionTemplate)

	// Add commands
	rootCmd.AddCommand(autoShellCmd)

	rootCmd.Flags().BoolVarP(&options.Key, "key", "k", false, "add key to result")
	rootCmd.Flags().BoolVarP(&options.Id, "id", "i", false, "add id to result")
	rootCmd.Flags().BoolVarP(&options.Type, "type", "t", false, `add type to result`)
	rootCmd.Flags().BoolVarP(&options.Object, "object", "o", false, "add object to result")
	rootCmd.Flags().BoolVarP(&options.Link, "link", "l", false, `add link to result`)
	rootCmd.Flags().BoolVarP(&options.LinkId, "linkid", "I", false, "add link id to result")
	rootCmd.Flags().BoolVarP(&options.Name, "name", "n", false, "add name in particular topology to result")
	rootCmd.Flags().BoolVarP(&options.Path, "path", "p", false, `add path to result`)

	rootCmd.Flags().BoolVarP(&options.Remove, "remove", "r", false, "remove result")
	rootCmd.Flags().BoolVarP(&confirm, "confirm", "y", false, "confirm remove")
}

var autoShellCmd = &cobra.Command{
	Use:    "autoshell",
	Short:  "Generate bash completion script",
	Long:   "Generate bash completion script",
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Root().GenBashCompletionFile("/etc/bash_completion.d/qdsl.sh")
	},
}
