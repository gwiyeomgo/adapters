package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
)

type Task map[string]string

func (t *Task) Set(key string, val string) {
	(*t)[key] = val
}

func (t *Task) Get(key string) string {
	return (*t)[key]
}

func (t Task) New() string {
	requestBody := map[string]interface{}{
		"users": map[string]interface{}{
			"to": []interface{}{
				map[string]interface{}{
					"type": "member",
					"member": map[string]interface{}{
						"organizationMemberId": t.Get("organizationMemberId"),
					},
				},
			},
		},
		"subject": "title",
		"body": map[string]interface{}{
			"mimeType": t.Get("mimeType"), //"text/html",
			"content":  t.Get("content"),
		},
		"dueDateFlag": true,
		"priority":    "none",
	}

	b, err := json.Marshal(requestBody)
	if err != nil {
		panic(err)
	}
	return string(b)
}

type DoorayAdapter struct {
	task      *Task
	doorayUrl string
	apiKey    string
}

func (d DoorayAdapter) Send() error {
	reqBody := bytes.NewBufferString(d.task.New())
	//request
	req, _ := http.NewRequest("POST", d.doorayUrl, reqBody)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", d.apiKey)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func ChangeContentType(content map[string]interface{}) string {
	var str string
	for key, item := range content {
		t := reflect.ValueOf(item).Type()
		if t.Kind() == reflect.Map {
			var line string
			myMap := item.(map[string]interface{})
			for k, i := range myMap {
				line += fmt.Sprintf("* %s \n```\n %s \n```\n", k, i)
			}
			str += line
		} else if t.Kind() == reflect.String {
			str += fmt.Sprintf("* %s \n```\n %s \n```\n", key, item)
		}
	}
	return str
}
