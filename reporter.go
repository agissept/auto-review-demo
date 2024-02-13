package main

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
)

type report struct {
	ChecklistsCompleted []string `json:"checklists_completed"`
	Message             string   `json:"message"`
}

func generateReport(c checklists) {
	var messages []string
	var checklistCompleted []string
	fields := reflect.VisibleFields(reflect.TypeOf(c))
	r := reflect.ValueOf(c)

	for _, field := range fields {
		f := reflect.Indirect(r).FieldByName(field.Name)
		message := f.FieldByName("comment").String()
		if message != "" {
			list := "<li>" + message + "</li>"
			fmt.Println(list)
			messages = append(messages, "<li>"+message+"</li>")
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

	reportJson, err := json.Marshal(report)
	if err != nil {
		unhandledException(err)
	}

	err = os.WriteFile("report.json", reportJson, os.ModePerm)
	if err != nil {
		unhandledException(err)
	}
}
