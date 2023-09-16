package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/mmlt/jo/cmd"
	"github.com/spf13/viper"
)

func main() {
	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()
	err := cmd.NewRootCmd().Execute()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
