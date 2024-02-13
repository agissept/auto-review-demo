package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func main() {
	c := newChecklists()
	studentId := 1234567

	submissionPath, _ := getParams()
	projectPath, packageJsonExists := getProjectPath(submissionPath)
	c.packageJsonExists = packageJsonExists

	fileJsPath, mainJsExists := getMainJs(submissionPath)
	c.mainJsExists = mainJsExists

	if projectPath != nil {
		runNpmInstall(*projectPath)

		if fileJsPath != nil {
			runMainJs(*fileJsPath)

			isServerUp := waitUntilServerUp()
			c.serveInPort5000 = isServerUp

			if isServerUp.status == true {
				html, rootShowingHtml := rootIsServingHtml()
				c.rootShowingHtml = rootShowingHtml

				if html != nil {
					c.htmlContainH1ElementWithStudentId = h1ElementIsCorrect(studentId, *html)
				}
			}

		}
	}

	if fileJsPath != nil {
		c.mainJsHaveStudentIdComment = checkCommentInMainJs(studentId, *fileJsPath)
	}

	stopServer()
	generateReport(c)
}

func getProjectPath(submissionPath string) (*string, checklist) {
	packageJsonPath, err := findFile(submissionPath, "package.json")
	if err != nil {
		unhandledException(err)
	}

	if packageJsonPath == nil {
		return nil, checklist{comment: "Kami tidak bisa menemukan file package.json pada submission yang kamu kirimkan", status: false}
	}

	projectPath := strings.Replace(*packageJsonPath, "/package.json", "", 1)

	return &projectPath, checklist{status: true}
}

func getMainJs(submissionPath string) (*string, checklist) {
	mainJsPath, err := findFile(submissionPath, "main.js")
	if err != nil {
		unhandledException(err)
	}

	if mainJsPath == nil {
		return nil, checklist{comment: "Kami tidak bisa menemukan file main.js pada submission yang kamu kirimkan", status: false}
	}

	return mainJsPath, checklist{status: true}
}

func checkCommentInMainJs(studentId int, fileJsPath string) checklist {
	file, err := os.ReadFile(fileJsPath)

	code := string(file)
	matchString, err := regexp.MatchString("//.*"+strconv.Itoa(studentId)+"|/\\*.*"+strconv.Itoa(studentId), code)
	if err != nil {
		unhandledException(err)
	}

	if matchString {
		return checklist{
			status:  true,
			comment: "",
		}
	}
	return checklist{
		status:  false,
		comment: "Kami tidak bisa menemukan user id " + strconv.Itoa(studentId) + " pada file main.js",
	}
}

func rootIsServingHtml() (*string, checklist) {
	response, err := http.Get("http://localhost:5000")
	if err != nil {
		unhandledException(err)
	}

	contentType := response.Header.Get("Content-Type")
	if !strings.Contains(contentType, "html") {
		return nil, checklist{status: true, comment: "Content yang berada pada root bukanlah html, melainkan " + contentType}
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		unhandledException(err)
	}

	html := string(responseData)
	return &html, checklist{status: true, comment: ""}
}

func h1ElementIsCorrect(studentId int, html string) checklist {
	compile := regexp.MustCompile("<h1>" + strconv.Itoa(studentId) + "</h1>")
	elementIsCorrect := compile.MatchString(html)

	if elementIsCorrect {
		return checklist{
			status:  true,
			comment: "",
		}
	}

	return checklist{
		status:  false,
		comment: "Kami tidak bisa menemukan user id " + strconv.Itoa(studentId) + " di url root pada aplikasi yang kamu buat",
	}
}

func waitUntilServerUp() checklist {
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
				return checklist{status: true}

			}
			if time.Since(start) > timeout {
				fmt.Println("Port 5000 is not running")
				return checklist{status: false, comment: "Kami tidak bisa mendeteksi port 5000 setelah aplikasi dijalankan, mohon periksa kembali apakah port yang kamu gunakan adalah 5000"}
			}
		}
		i++
	}
}
