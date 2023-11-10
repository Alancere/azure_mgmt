package common

import (
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
)

func AutorestCommand(armPath string, goExtension, goTestExtension string, azcoreVersion string) error {

	if azcoreVersion == "" {
		azcoreVersion = DefaultAzcore
	}

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
		"--testmodeler.generate-fake-test=true",
		fmt.Sprintf("--azcore-version=%s", azcoreVersion),
		"--debug",
		// "--gotest.debugger",
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

func Gomodtidy(dir string) error {
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	log.Printf("Result of `go mod tidy` execution: \n%s", string(output))
	if err != nil {
		return fmt.Errorf("failed to execute `go mod tidy` '%s': %+v", string(output), err)
	}
	return nil
}

func GoVet(dir string) error {
	cmd := exec.Command("go", "vet", ".")
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	log.Printf("Result of `go vet` execution: \n%s", string(output))
	if err != nil {
		return fmt.Errorf("failed to execute `go vet` '%s': %+v", string(output), err)
	}
	return nil
}

func RunFakeTest(dir string) error {
	cmd := exec.Command("go", "test", "-v", "fake_test.go")
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	log.Printf("Result of `run fake test` execution: \n%s", string(output))
	if err != nil {
		return fmt.Errorf("failed to execute `run fake test` '%s': %+v \n%s", string(output), err, dir)
	}
	return nil
}
