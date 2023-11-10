/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/Alancere/azure-mgmt/generate/common"
)

/*
	运行autorest cmd
*/
// autorestCmd represents the autorest command
var autorestCmd = &cobra.Command{
	Use:   "autorest",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("autorest called")
		if err := autorest(); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(autorestCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// autorestCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// autorestCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

const (
	generatefakeConfig = `azcore-version: %s
generate-fakes: true
inject-spans: true
`
)

func autorest() error {
	var count int
	mgmtRepo := viper.GetString("mgmtRepo")

	armServices := common.SliceConvertMap(viper.GetStringSlice("ArmServices"))
	skipArmServices := common.SliceConvertMap(viper.GetStringSlice("SkipArmServices"))

	if err := filepath.WalkDir(mgmtRepo, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() || !strings.Contains(d.Name(), "arm") {
			return nil
		}

		if _, ok := skipArmServices[d.Name()]; ok {
			return nil
		}

		if _, ok := armServices[d.Name()]; !ok {
			return nil
		}

		count++
		fmt.Println(path)

		err = generateFake(path, viper.GetString("azcoreVersion"))
		if err != nil {
			return err
		}

		// run autorest command
		if err := common.AutorestCommand(path, viper.GetString("CodeGenVersion.go"), viper.GetString("CodeGenVersion.gotest"), viper.GetString("azcoreVersion")); err != nil {
			// uninterrupted
			fmt.Println(err)
		}

		// go mod tidy
		if err := common.Gomodtidy(path); err != nil {
			fmt.Println(err)
		}

		// go vet
		if err := common.GoVet(path); err != nil {
			fmt.Println(err)
		}

		return nil
	}); err != nil {
		return err
	}

	fmt.Println("run autorest command counts:", count)
	return nil
}

func generateFake(armPath string, azcoreVersion string) error {
	// azcore-version
	if azcoreVersion == "" {
		azcoreVersion = common.DefaultAzcore
	}

	//aad generate-fakes: true in autorest.md
	// 1. read autprest.md
	var newAutorestData []byte
	autorestPath := filepath.Join(armPath, "autorest.md")
	autorestData, err := os.ReadFile(autorestPath)
	newAutorestData = autorestData
	if err != nil {
		return err
	}
	if !strings.Contains(string(autorestData), "generate-fakes: true") {
		index := bytes.LastIndex(autorestData, []byte("```"))
		if index <= 0 {
			return nil
		}
		fail := append([]byte(fmt.Sprintf(generatefakeConfig, azcoreVersion)), autorestData[index:]...)
		newAutorestData = append(autorestData[:index], fail...)
	}

	// 2. Delete the module set in autorest.md, autorest.go@4.0.0-preview.55 fixed
	if strings.Contains(string(autorestData), "generate-fakes: true") {
		if strings.Contains(string(autorestData), "module:") {
			// remove module line
			index := bytes.Index(autorestData, []byte("module:"))
			if index != 0 {
				head := autorestData[:index]
				tail := autorestData[index:]
				split := strings.Split(string(tail), "\n")
				split = split[1:]
				newAutorestData = append(head, []byte(strings.Join(split, "\n"))...)
			}
		}
	}

	if len(newAutorestData) != 0 && len(autorestData) != len(newAutorestData) {
		if err := os.WriteFile(autorestPath, newAutorestData, 0666); err != nil {
			return err
		}
	}

	return nil
}
