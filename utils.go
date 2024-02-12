package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"time"
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
	waitUntilServerUp()
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

func waitUntilServerUp() {
	host := "localhost"
	port := "5000"
	timeout := time.Second * 3

	var i int
	for start := time.Now(); ; {
		if i%10 == 0 {
			conn, _ := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout)
			if conn != nil {
				conn.Close()
				fmt.Println("Port 5000 is running")
				break
			}
			if time.Since(start) > timeout {
				fmt.Println("Port 5000 is not running")
				break
			}
		}
		i++
	}
}

func stopServer() {
	err := exec.Command("bash", "-c", "kill -9 $(lsof -t -i:5000)").Start()
	if err != nil {
		unhandledException(err)
	}
	fmt.Println("Server stopped")
}

type report struct {
	checklistsCompleted []string
	message             string
}

func generateReport(c checklists) {
	var messages []string
	fields := reflect.VisibleFields(reflect.TypeOf(c))
	r := reflect.ValueOf(c)

	for _, field := range fields {
		f := reflect.Indirect(r).FieldByName(field.Name)
		messages = append(messages, f.FieldByName("comment").String())
	}

	fmt.Println(messages)

	//fmt.Println(v)

	//var report report
	//if checklists.packageJsonExists.status {
	//	report.checklistsCompleted = append(report.checklistsCompleted, "package_json_exist")
	//}
	//
	//if checklists.mainJsExists.status {
	//	report.checklistsCompleted = append(report.checklistsCompleted, "main_js_exist")
	//}
	//
	//if checklists.serveInPort5000.status {
	//	report.checklistsCompleted = append(report.checklistsCompleted, "using_port_5000")
	//}
	//
	//if checklists.rootShowingHtml.status {
	//	report.checklistsCompleted = append(report.checklistsCompleted, "root_show_html")
	//}
	//
	//if checklists.htmlContainH1ElementWithStudentId.status {
	//	report.checklistsCompleted = append(report.checklistsCompleted, "html_contain_h1_element")
	//}
	//
	//if checklists.mainJsHaveStudentIdComment.status {
	//	report.checklistsCompleted = append(report.checklistsCompleted, "main_js_contain_student_id")
	//}

}
