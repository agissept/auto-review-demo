package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func findFile(submissionPath string, fileName string) (*string, error) {
	var file *string

	err := filepath.Walk(submissionPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && (strings.Contains(path, "node_modules") || strings.Contains(path, ".git")) {
			return filepath.SkipDir
		}

		if !info.IsDir() && filepath.Base(path) == fileName {
			file = &path
			return io.EOF
		}

		return nil
	})

	if err != nil && err != io.EOF {
		return nil, err
	}

	return file, nil
}

func getParams() (string, string) {
	submissionsPath := flag.String("submission", "", "Specify submissions path")
	reportPath := flag.String("report", "", "Specify report path")
	flag.Parse()

	if *submissionsPath == "" {
		fmt.Println("submission path (-submission) tidak boleh kosong")
		os.Exit(1)
	}

	if *reportPath == "" {
		fmt.Println("report path (-report) tidak boleh kosong")
		os.Exit(1)
	}

	return *submissionsPath, *reportPath
}

func runNpmInstall(projectPath string) {
	cmd := exec.Command("bash", "-c", "npm install")
	cmd.Dir = projectPath
	err := cmd.Start()
	if err != nil {

		unhandledException(err)
	}
}

func runMainJs(fileJsPath string) {
	cmd := exec.Command("bash", "-c", "node "+fileJsPath)
	if err := cmd.Start(); err != nil {
		unhandledException(err)
	}
}

func unhandledException(err error) {
	stopServer()

	pc := make([]uintptr, 10)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	file, line := f.FileLine(pc[0])
	fmt.Printf("%s:%d %s\n", file, line, f.Name())
	fmt.Printf(err.Error())

	os.Exit(1)
}

func stopServer() {
	err := exec.Command("bash", "-c", "kill -9 $(lsof -t -i:5000)").Start()
	if err != nil {
		unhandledException(err)
	}
	fmt.Println("Server stopped")
}

type autoReviewConfig struct {
	SubmitterId int    `json:"submitter_id"`
	Username    string `json:"submitter_name"`
}

func getAutoReviewConfig(submissionPath string) autoReviewConfig {
	configPath := filepath.Join(submissionPath, "auto-review-config.json")
	config, err := os.ReadFile(configPath)
	if err != nil {
		unhandledException(err)
	}

	var autoReviewConfig autoReviewConfig

	err = json.Unmarshal(config, &autoReviewConfig)
	if err != nil {
		unhandledException(err)
	}

	return autoReviewConfig
}
