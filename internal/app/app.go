package app

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"strings"

	"TestRIT/internal/logger"
	"TestRIT/internal/methods"
)

func RunApp() {
	//Inti logger
	logger.Init()
	logger.InfoLogger.Println("Starting the application...")

	//Read console's args
	readArgs := os.Args[1:]
	pathToFile := strings.Join(readArgs, "")

	//Open file
	f, err := os.Open(pathToFile)
	if err != nil {
		logger.ErrorLogger.Println(err)
	}
	defer f.Close()

	//Read file with reader
	wr := bytes.Buffer{}
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		wr.WriteString(sc.Text())
	}

	//Unmarshal JSON
	var jsonFile interface{}
	err = json.Unmarshal(wr.Bytes(), &jsonFile)
	if err != nil {
		logger.ErrorLogger.Println(err)
	}

	//Writing result JSON
	result := methods.СheckJson(jsonFile, false)

	resulJson, resErr := json.Marshal(result)
	if resErr != nil {
		logger.ErrorLogger.Println(resErr)
	}

	pathToSave := filepath.Dir(pathToFile)
	savePath := fmt.Sprintf("%s\\resultJSON.json", pathToSave)

	err = os.WriteFile(savePath, resulJson, 0644)
	if err != nil {
		logger.ErrorLogger.Println(err)
	}
}
