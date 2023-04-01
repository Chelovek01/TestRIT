package app

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	// "strings"

	"github.com/djherbis/times"
)

type File struct {
	Data interface{}
}

type Action struct {
	Action      string
	Result      string
	Params      []string
	Condition   string
	Next_action string
}

func RunApp() {

	f, err := os.Open("test.json")
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	// Чтение файла с ридером
	wr := bytes.Buffer{}
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		wr.WriteString(sc.Text())
	}

	var jsonFile interface{}
	err = json.Unmarshal(wr.Bytes(), &jsonFile)
	if err != nil {
		fmt.Println(err)
	}

	СheckJson(jsonFile)

}

func CreateFile(name string) (err error) {

	file, err := os.Create(name)
	if err != nil {
		fmt.Println(err)
	}
	file.Close()
	return

}

func RenameFile(oldPath string, newPath string) {
	err := os.Rename(oldPath, newPath)
	if err != nil {
		fmt.Println(err)
	}
}

func RemoveFile(name string) {
	err := os.Remove(name)
	if err != nil {
		fmt.Println(err)
	}
}

func bTime(name string) {
	t, err := times.Stat(name)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(t.AccessTime())
}

func writeString(name string, data string) {

	file, err := os.OpenFile(name, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		fmt.Println(err)
	}

	defer file.Close()

	file.WriteString(data)
}

func СheckJson(data interface{}) {

	myMap := data.(map[string]interface{})

	// nextLevlJson := make(map[string]interface{})

	switch myMap["action"] {

	case "create":

		name := myMap["params"].([]interface{})[0].(string)
		CreateFile(name)

		fmt.Println("created", name)

		result := fmt.Sprintf("file %s created", name)

		myMap["result"] = result

	case "condition":

		conditionTime := "2006-01-02 15:04:01"

		cond := myMap["params"].([]interface{})[0].(string)

		if cond > conditionTime {

			result := "condition was true"
			myMap["result"] = result

			fmt.Println(result)

			СheckJson(myMap["next_true"])
		}

		if cond < conditionTime {
			result := "condition was false"
			myMap["result"] = result

			fmt.Println(result)

			СheckJson(myMap["next_false"])
		}

	case "write_string":

		file_name := myMap["params"].([]interface{})[0].(string)

		dataWrite := myMap["params"].([]interface{})[1].(string)

		result := fmt.Sprintf("random string wrote in %s", file_name)
		myMap["result"] = result

		fmt.Println(result)

		writeString(file_name, dataWrite)

	case "change_name":

		oldPath := myMap["params"].([]interface{})[0].(string)

		newPath := myMap["params"].([]interface{})[1].(string)

		result := fmt.Sprintf("the file  %s was renamed to %s", oldPath, newPath)
		myMap["result"] = result

		fmt.Println(result)

		RenameFile(oldPath, newPath)

	case "get_btime":
		fileName := myMap["params"].([]interface{})[0].(string)

		result := "was born"
		myMap["result"] = result

		fmt.Println(result)

		bTime(fileName)

	case "remove_file":

		fileName := myMap["params"].([]interface{})[0].(string)

		result := fmt.Sprintf("file  %s removed", fileName)
		myMap["result"] = result

		fmt.Println("file removed", fileName)

		RemoveFile(fileName)

	}

	_, ok := myMap["next"]

	if ok {
		СheckJson(myMap["next"])
	} else {

		fmt.Println(myMap)

		return
	}

}
