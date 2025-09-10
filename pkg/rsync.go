package rsync

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

func Wrapper(ctx context.Context, path, homeDir string) error {
	configFilePath := filepath.Join(path, WRAPPER_CONFIG)
	globalConfigFilePath := filepath.Join(homeDir, WRAPPER_CONFIG)

	projectConfigFile, err := os.Stat(configFilePath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	_, err = os.Stat(globalConfigFilePath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	var configFileContent []byte
	if projectConfigFile != nil {
		logrus.Debug("using project debug file")
		configFileContent, err = os.ReadFile(configFilePath)
	} else {
		logrus.Debug("using global config file")
		configFileContent, err = os.ReadFile(globalConfigFilePath)
	}

	if err != nil {
		return err
	}

	logrus.Debugf("config file contents \n%v\n", string(configFileContent))
	configObj := &RsyncWrapper{}
	if err := yaml.Unmarshal(configFileContent, configObj); err != nil {
		return err
	}

	cmd, args := generateCommand(path, *configObj)

	return executeCommand(ctx, cmd, args)
}

func generateCommand(sourcePath string, config RsyncWrapper) (string, []string) {
	command := "rsync"
	commandArgs := []string{"-zarP"}
	if config.RemoveDestinationFiles {
		commandArgs = append(commandArgs, "--delete")
	}

	for _, v := range config.ExcludeDirs {

		commandArgs = append(commandArgs, fmt.Sprintf("--exclude=%s", strings.TrimSpace(v)))
	}

	if config.DryRun {
		commandArgs = append(commandArgs, "--dry-run")
	}
	commandArgs = append(commandArgs, sourcePath)
	dest := fmt.Sprintf("%s:%s", config.DestinationAddress, config.DestinationPath)
	commandArgs = append(commandArgs, dest)
	return command, commandArgs
}

func executeCommand(ctx context.Context, cmd string, args []string) error {
	var wg sync.WaitGroup
	cmdPath, err := exec.LookPath(cmd)
	if err != nil {
		return err
	}

	cmdObj := exec.CommandContext(ctx, cmdPath, args...)
	stdOut, err := cmdObj.StdoutPipe()
	if err != nil {
		return err
	}

	stdErr, err := cmdObj.StderrPipe()
	if err != nil {
		return err
	}

	wg.Add(2)
	go func() {
		defer wg.Done()
		print(stdOut)
	}()

	go func() {
		defer wg.Done()
		print(stdErr)
	}()
	if err := cmdObj.Start(); err != nil {
		return err
	}

	wg.Wait()
	return cmdObj.Wait()
}

func print(reader io.Reader) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		fmt.Print(scanner.Text())
	}
}
