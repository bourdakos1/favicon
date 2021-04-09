package cmd

import (
	"fmt"
	"os"

	"github.com/bourdakos1/favicon/cmd/pack"
	"github.com/bourdakos1/favicon/cmd/unpack"
	"github.com/bourdakos1/favicon/version"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "bloop",
	Short: "A command line tool for creating favicons",
	Long: fmt.Sprintf(`bloop (%s)

A command line tool for creating favicons.

For more info visit: https://example.com`, version.BuildVersion()),
	Version: version.BuildVersion(),
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "pack <folder> <favicon.ico>",
		Short: "Pack PNGs into a favicon",
		Long:  `Pack a folder of PNGs into a favicon.ico file.`,
		Run:   pack.Run,
	})

	rootCmd.AddCommand(&cobra.Command{
		Use:   "unpack <favicon.ico>",
		Short: "Unpack the PNGs inside a favicon",
		Long:  `Unpack all the PNG files stored inside the favicon`,
		Run:   unpack.Run,
	})
}
