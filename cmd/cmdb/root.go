// Copyright 2023 NJWS Inc.

package main

import (
	"git.fg-tech.ru/listware/cmdb/internal/server"
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
	Use:     "cmdb",
	Short:   "CMDB agent",
	Long:    `CMDB agent`,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Println("version", version, "release", release)
		return server.Run(cmd.Context())
	},
}

func init() {
	rootCmd.SetVersionTemplate(versionTemplate)

	// Add commands
	rootCmd.AddCommand(autoShellCmd)
}

var autoShellCmd = &cobra.Command{
	Use:    "autoshell",
	Short:  "Generate bash completion script",
	Long:   "Generate bash completion script",
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Root().GenBashCompletionFile("/etc/bash_completion.d/cmdb.sh")
	},
}
