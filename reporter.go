package main

import (
	"bytes"
	"encoding/json"
	"os"
	"reflect"
	"strings"
)

type report struct {
	ChecklistsCompleted []string `json:"checklists_completed"`
	Message             string   `json:"message"`
}

func generateReport(checklist checklists, reportPath string, username string) {
	report := createReport(checklist)
	report.Message = generateTemplatedMessage(report, username)
	save(report, reportPath)
}

func generateTemplatedMessage(report report, username string) string {
	if isSubmissionApproved(report) {
		return "Selamat <b>" + username + "!!</b> kamu telah lolos submission ini."
	}

	return "Hallo " + username + " masih terdapat beberapa kesalahan, berikut adalah kesalahan yang terjadi <ul>" + report.Message + "</ul>. Silakan diperbaiki yaa."
}

func isSubmissionApproved(report report) bool {
	return len(report.ChecklistsCompleted) == 6
}

func createReport(checklists checklists) report {
	var messages []string
	var checklistCompleted []string
	fields := reflect.VisibleFields(reflect.TypeOf(checklists))
	r := reflect.ValueOf(checklists)

	for _, field := range fields {
		f := reflect.Indirect(r).FieldByName(field.Name)
		message := f.FieldByName("comment").String()
		if message != "" {
			list := "<li>" + message + "</li>"
			messages = append(messages, list)
		}

		status := f.FieldByName("status").Bool()
		if status {
			checklistCompleted = append(checklistCompleted, field.Tag.Get("json"))
		}

	}

	messageString := strings.Join(messages, "")

	report := report{
		ChecklistsCompleted: checklistCompleted,
		Message:             messageString,
	}

	return report
}

func save(report report, reportPath string) {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetIndent("", "    ")
	enc.SetEscapeHTML(false)
	err := enc.Encode(report)

	if err != nil {
		unhandledException(err)
	}

	err = os.MkdirAll(reportPath, os.ModePerm)
	if err != nil {
		unhandledException(err)
	}

	err = os.WriteFile(reportPath+"/report.json", buf.Bytes(), os.ModePerm)
	if err != nil {
		unhandledException(err)
	}
}
