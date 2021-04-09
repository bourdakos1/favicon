package unpack

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bourdakos1/favicon/favicon"
	"github.com/spf13/cobra"
)

func Run(cmd *cobra.Command, args []string) {
	path := args[0]

	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer file.Close()

	icon, err := favicon.New(file)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	extension := filepath.Ext(path)
	baseName := path[0 : len(path)-len(extension)]

	err = icon.SaveAsPNGs(baseName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
