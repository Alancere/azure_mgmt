package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	// autorestGo     = "@autorest/go@4.0.0-preview.55"
	// autorestGoTest = "@autorest/go@4.6.2"

	// specify github.com/Azure/azure-sdk-for-go/sdk/azcore version
	azcoreVersion = "1.8.0-beta.2"
	generatefake  = fmt.Sprintf(`azcore-version: %s
generate-fakes: true
inject-spans: true
`, azcoreVersion)
)

// autorest generate fake
func main() {
	count := 0

	// replace with your env
	mgmtRepo := "D:/Go/src/github.com/Azure/dev/azure-sdk-for-go/sdk/resourcemanager"
	goExtension := "D:/Go/src/github.com/Azure/autorest.go/packages/autorest.go"
	goTestExtension := "D:/Go/src/github.com/Azure/autorest.go/packages/autorest.gotest"
	_ = goTestExtension

	// If specified as empty, all services will be run.
	specifyArmService := "armcdn"

	if err := filepath.WalkDir(mgmtRepo, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() || !strings.Contains(d.Name(), "arm") {
			return nil
		}

		if specifyArmService != "" {
			if specifyArmService != d.Name() {
				return nil
			}
		}

		count++
		fmt.Println(path)

		err = generateFake(path)
		if err != nil {
			return err
		}

		// run autorest command
		if err := autorestCommand(path, goExtension, goTestExtension, azcoreVersion); err != nil {
			// uninterrupted
			fmt.Println(err)
		}

		// go mod tidy
		if err := gomodtidy(path); err != nil {
			fmt.Println(err)
		}

		// go vet
		if err := goVet(path); err != nil {
			fmt.Println(err)
		}

		return nil
	}); err != nil {
		log.Fatal(err)
	}

	fmt.Println("run autorest command counts:", count)
}

func generateFake(armPath string) error {
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
		fail := append([]byte(generatefake), autorestData[index:]...)
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

func autorestCommand(armPath string, goExtension, goTestExtension string, azcorVersion string) error {

	args := []string{}
	if goExtension != "" {
		args = append(args, fmt.Sprintf("--use=%s", goExtension))
	}

	if goTestExtension != "" {
		args = append(args, fmt.Sprintf("--use=%s", goTestExtension))
	}

	args = append(args,
		"--go",
		"--track2",
		fmt.Sprintf("--output-folder=%s", armPath),
		"--go.clear-output-folder=false",
		"--generate-sdk=true",
		"--testmodeler.generate-sdk-example=true",
		fmt.Sprintf("--azcore-version=%s", azcoreVersion),
		"--debug",
		filepath.Join(armPath, "autorest.md"))
	cmd := exec.Command("autorest", args...)

	fmt.Printf("autorest command:::%s\n", cmd.Args)
	output, err := cmd.CombinedOutput()
	log.Printf("Result of `generate fake test` execution: \n%s", string(output))
	if err != nil {
		return fmt.Errorf("failed to execute `generate fake test` '%s': %+v", string(output), err)
	}
	return nil
}

func gomodtidy(dir string) error {
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	log.Printf("Result of `go mod tidy` execution: \n%s", string(output))
	if err != nil {
		return fmt.Errorf("failed to execute `go mod tidy` '%s': %+v", string(output), err)
	}
	return nil
}

func goVet(dir string) error {
	cmd := exec.Command("go", "vet", ".")
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	log.Printf("Result of `go vet` execution: \n%s", string(output))
	if err != nil {
		return fmt.Errorf("failed to execute `go vet` '%s': %+v", string(output), err)
	}
	return nil
}
