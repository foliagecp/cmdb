// Copyright 2023 NJWS Inc.
// Copyright 2022 Listware

package main

import (
	"fmt"
	"os"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
