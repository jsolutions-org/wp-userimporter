package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func exitWithError(err error) {
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	os.Exit(1)
}

func checkErr(e error) {
	if e != nil {
		exitWithError(e)
	}
}

func createRequestBody(headerline []string, userline []string) (string, error) {
	if len(userline) != len(headerline) {
		return "", errors.New("line must match headers format")
	}

	var sb strings.Builder

	// create json body
	sb.WriteString("{")
	for i, key := range headerline {
		sb.WriteString("\"" + key + "\"")
		sb.WriteString(":")
		sb.WriteString("\"" + userline[i] + "\"")
		if i < len(headerline)-1 {
			sb.WriteString(",")
		}
	}
	sb.WriteString("}")

	reqBody := sb.String()

	return reqBody, nil
}

func callWpApi(wpurl string, username string, password string, requestbody string) {
	client := http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest("POST", wpurl, strings.NewReader(requestbody))
	if err != nil {
		exitWithError(errors.New("error building http request"))
	}

	req.SetBasicAuth(username, password)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		exitWithError(fmt.Errorf("http post request error: %s", err.Error()))
	}

	fmt.Printf("HTTP Status: %d - %s\n", resp.StatusCode, resp.Status)

	defer resp.Body.Close()
}

func main() {
	if len(os.Args) != 5 {
		fmt.Printf("Usage: %s <csv-file> <wp-url> <username> <password>\n", os.Args[0])
		fmt.Printf("For more information see https://github.com/jsolutions-org/wp-userimporter\n")
		err := errors.New("wrong argument count")
		exitWithError(err)
	}

	// open csv file
	file, err := os.Open(os.Args[1])
	checkErr(err)
	defer file.Close()

	// create request header
	wpurl := os.Args[2] + "/wp-json/wp/v2/users"
	username := os.Args[3]
	password := os.Args[4]

	var headerline, userline []string

	reader := csv.NewReader(file)
	headerline, err = reader.Read()
	checkErr(err)

	for {
		userline, err = reader.Read()

		//TEST
		fmt.Printf("processing: %s\n",userline)

		if err == io.EOF {
			break
		}
		checkErr(err)

		// form request body
		requestbody, err := createRequestBody(headerline, userline)
		if err != nil {
			fmt.Printf("Error in Line: %sError: %s\n", userline, err)
			fmt.Printf("Line will not be processed. Continiue with next line.\n")
			continue
		}

		// call WP REST API
		callWpApi(wpurl, username, password, requestbody)
	}

}
