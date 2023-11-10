/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// releaseCmd represents the release command
var releaseCmd = &cobra.Command{
	Use:   "release",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("release called")
		if err := release(); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(releaseCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// releaseCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// releaseCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func release() error {
	f := Flags{}
	err := viper.Unmarshal(&f)
	if err != nil {
		return err
	}

	armServicesData := viper.GetStringMap("ArmServices")

	armServices := ParseArmServices(armServicesData)



	return nil
}

type Flags struct {
	ReleaseV2Flags struct {
		ForceStableVersion  string
		GoVersion           string
		PackageConfig       string
		PackageTitle        string
		ReleaseDate         string
		SkipCreateBranch    bool
		SkipGenerateExample bool
		UpdateSpecVersion   bool
		SpecRepo            string
		SpecRPName          string
		Token               string
		VersionNumber       string
	}
}

// --force-stable-version true --force-stable-version=
const (
	ForceStableVersion  = "--force-stable-version"
	GoVersion           = "--go-version"
	PackageConfig       = "--package-config"
	PackageTitle        = "--package-title"
	ReleaseDate         = "--release-date"
	SkipCreateBranch    = "--skip-create-branch"
	SkipGenerateExample = "--skip-generate-example"
	UpdateSpecVersion   = "--update-spec-version"
	SpecRepo            = "--spec-repo"
	SpecRPName          = "--spec-rp-name"
	Token               = "--token"
	VersionNumber       = "--version-number"
)

func (f Flags) String() string {

	flagString := ""
	if f.ReleaseV2Flags.ForceStableVersion != "" {
		flagString = fmt.Sprintf("%s --force-stable-version %s", flagString, f.ReleaseV2Flags.ForceStableVersion)
	}

	if f.ReleaseV2Flags.GoVersion != "" {
		flagString = fmt.Sprintf("%s --go-version %s", flagString, f.ReleaseV2Flags.GoVersion)
	}

	if f.ReleaseV2Flags.PackageConfig != "" {
		flagString = fmt.Sprintf("%s --package-config %s", flagString, f.ReleaseV2Flags.PackageConfig)
	}

	if f.ReleaseV2Flags.PackageTitle != "" {
		flagString = fmt.Sprintf("%s --package-title %s", flagString, f.ReleaseV2Flags.PackageTitle)
	}

	if f.ReleaseV2Flags.ReleaseDate != "" {
		flagString = fmt.Sprintf("%s --release-date %s", flagString, f.ReleaseV2Flags.ReleaseDate)
	}

	if f.ReleaseV2Flags.SkipCreateBranch {
		flagString = fmt.Sprintf("%s --skip-create-branch %t", flagString, f.ReleaseV2Flags.SkipCreateBranch)
	}

	if f.ReleaseV2Flags.SkipGenerateExample {
		flagString = fmt.Sprintf("%s --skip-generate-example %t", flagString, f.ReleaseV2Flags.SkipGenerateExample)
	}

	if !f.ReleaseV2Flags.UpdateSpecVersion {
		flagString = fmt.Sprintf("%s --update-spec-version=%t", flagString, f.ReleaseV2Flags.UpdateSpecVersion)
	}

	if f.ReleaseV2Flags.SpecRepo != "" {
		flagString = fmt.Sprintf("%s --spec-repo %s", flagString, f.ReleaseV2Flags.SpecRepo)
	}

	if f.ReleaseV2Flags.SpecRPName != "" {
		flagString = fmt.Sprintf("%s --spec-rp-name %s", flagString, f.ReleaseV2Flags.SpecRPName)
	}

	if f.ReleaseV2Flags.Token != "" {
		flagString = fmt.Sprintf("%s --token %s", flagString, f.ReleaseV2Flags.Token)
	}

	if f.ReleaseV2Flags.VersionNumber != "" {
		flagString = fmt.Sprintf("%s --version-number %s", flagString, f.ReleaseV2Flags.VersionNumber)
	}

	return flagString
}

func (f Flags)Copy() Flags {
	newF := Flags{}
	
}

// --\S* \S*
func (f *Flags) ReplaceConfig(config string) {
	flags := regexp.MustCompile(`--\S* \S*`).FindAllString(config, -1)
	for _, v := range flags {
		before, after, _ := strings.Cut(v, " ")
		switch strings.TrimSpace(before) {
		case ForceStableVersion:
			f.ReleaseV2Flags.ForceStableVersion = after
		case GoVersion:
			f.ReleaseV2Flags.GoVersion = after
		case PackageConfig:
			f.ReleaseV2Flags.PackageConfig = after
		case PackageTitle:
			f.ReleaseV2Flags.PackageTitle = after
		case ReleaseDate:
			f.ReleaseV2Flags.ReleaseDate = after
		case SkipCreateBranch:
			temp, _ := strconv.ParseBool(after)
			f.ReleaseV2Flags.SkipCreateBranch = temp
		case SkipGenerateExample:
			temp, _ := strconv.ParseBool(after)
			f.ReleaseV2Flags.SkipGenerateExample = temp
		case UpdateSpecVersion:
			temp, _ := strconv.ParseBool(after)
			f.ReleaseV2Flags.UpdateSpecVersion = temp
		case SpecRepo:
			f.ReleaseV2Flags.SpecRepo = after
		case SpecRPName:
			f.ReleaseV2Flags.SpecRPName = after
		case Token:
			f.ReleaseV2Flags.Token = after
		case VersionNumber:
			f.ReleaseV2Flags.VersionNumber = after
		}
	}

}

type ArmService struct {
	ServiceName    string
	ArmServiceName string
	Config         string
}

func ParseArmServices(data map[string]any) []ArmService {
	armServices := make([]ArmService, 0, len(data))

	for k,v := range data {
		if v == nil {
			arm := ArmService{
				ServiceName: k,
				// ArmServiceName: fmt.Sprintf("arm%s", k),
			}
			armServices = append(armServices, arm)
		}else {
			 switch v.(type) {
			 case string:
				arm := ArmService{
					ServiceName: k,
					Config: v.(string),
				}
				armServices = append(armServices, arm)
			 case map[string]interface{}:
				for armServiceName, c := range v.(map[string]any) {
					arm := ArmService{
						ServiceName: k,
						ArmServiceName: armServiceName,
						Config: c.(string),
					}
					armServices = append(armServices, arm)
				}
			 }
		}
	}

	return armServices
}
