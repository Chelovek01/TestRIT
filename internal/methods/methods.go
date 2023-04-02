package methods

import (
	"fmt"
	"os"
	"time"

	"TestRIT/internal/logger"

	"github.com/djherbis/times"
)

// Create creates or truncates the named file.
func CreateFile(name string) (err error) {

	file, err := os.Create(name)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return
	}
	file.Close()
	return

}

// Renames the selected file
// File must be in project root
func RenameFile(oldPath string, newPath string) (err error) {
	err = os.Rename(oldPath, newPath)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return
	}
	return
}

// Remove a file from project root
func RemoveFile(name string) (err error) {
	err = os.Remove(name)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return
	}
	return
}

// Returns the creation date of a file
func bTime(name string) (t times.Timespec, err error) {
	t, err = times.Stat(name)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return
	}
	return
}

// Write strings at file
// File must be in project root
func WriteString(name string, data string) (err error) {

	file, err := os.OpenFile(name, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return
	}

	defer file.Close()

	_, err = file.WriteString(data)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return
	}
	return
}

// Parses the selected JSON file
// JSON must be valid
func СheckJson(data interface{}, useParams bool, params ...interface{}) (myMap map[string]interface{}) {

	myMap = data.(map[string]interface{})

	if len(params) > 0 && useParams {
		myMap["params"] = params[0]
		useParams = true
	}

	_, ok := myMap["next_params"]
	if ok {
		useParams = true
	} else {
		useParams = false
	}

	switch myMap["action"] {

	case "create":

		name := myMap["params"].([]interface{})[0].(string)
		err := CreateFile(name)
		if err != nil {
			logger.ErrorLogger.Println(err)
			myMap["result"] = err.Error()
			break
		}

		myMap["result"] = fmt.Sprintf("file %s created", name)

	case "condition":

		conditionTime := time.Now()

		cond := myMap["params"].([]interface{})[0].(string)
		timej, err := time.Parse("2006-01-02 15:04:05", cond)
		if err != nil {
			logger.ErrorLogger.Println(err)
			myMap["result"] = "wrong time"
		}
		if conditionTime.After(timej) {

			myMap["result"] = "condition was true"
			if useParams {
				СheckJson(myMap["next_true"], useParams, myMap["next_params"])
			} else {
				СheckJson(myMap["next_true"], useParams)
			}

		} else {
			myMap["result"] = "condition was false"
			if useParams {
				СheckJson(myMap["next_false"], useParams, myMap["next_params"])
			} else {
				СheckJson(myMap["next_false"], useParams)
			}
		}

	case "write_string":

		file_name := myMap["params"].([]interface{})[0].(string)

		dataWrite := myMap["params"].([]interface{})[1].(string)

		err := WriteString(file_name, dataWrite)
		if err != nil {
			logger.ErrorLogger.Println(err)
			myMap["result"] = err.Error()
			break
		}
		myMap["result"] = fmt.Sprintf("string '%s' wrote in %s", dataWrite, file_name)

	case "change_name":

		oldPath := myMap["params"].([]interface{})[0].(string)

		newPath := myMap["params"].([]interface{})[1].(string)

		err := RenameFile(oldPath, newPath)
		if err != nil {
			logger.ErrorLogger.Println(err)
			myMap["result"] = err.Error()
			break
		}
		myMap["result"] = fmt.Sprintf("the file  %s was renamed to %s", oldPath, newPath)

	case "get_btime":
		fileName := myMap["params"].([]interface{})[0].(string)
		time, err := bTime(fileName)
		if err != nil {
			logger.ErrorLogger.Println(err)
			myMap["result"] = err.Error()
			break
		}
		myMap["result"] = time

	case "remove_file":

		fileName := myMap["params"].([]interface{})[0].(string)
		err := RemoveFile(fileName)
		if err != nil {
			logger.ErrorLogger.Println(err)
			myMap["result"] = err.Error()
			break
		}
		myMap["result"] = fmt.Sprintf("file  %s removed", fileName)
	}

	_, ok = myMap["next"]

	if ok {
		if useParams {
			СheckJson(myMap["next"], useParams, myMap["next_params"])
		} else {
			СheckJson(myMap["next"], useParams)
		}
	}

	return myMap
}
