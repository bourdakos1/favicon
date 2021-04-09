package pack

import (
	"log"
	"os"
	"path/filepath"

	"github.com/bourdakos1/favicon/favicon"
	"github.com/spf13/cobra"
)

func Run(cmd *cobra.Command, args []string) {
	folder := args[0]
	output := args[1]

	var files []string
	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	favicon.Pack(files, output)
}
