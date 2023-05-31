package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Jeffail/gabs/v2"
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

func fetchUserId(wpurl string, username string, password string, line string) string {
	client := http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest("GET", wpurl + "?search=" + line, strings.NewReader(""))
	if err != nil {
		exitWithError(errors.New("error building http request"))
	}

	req.SetBasicAuth(username, password)
	
	resp, err := client.Do(req)
	if err != nil {
		exitWithError(fmt.Errorf("http get request error: %s", err.Error()))
	}

	defer resp.Body.Close()

	fmt.Printf("HTTP Status: %d - %s\n", resp.StatusCode, resp.Status)

	responseData, err := io.ReadAll(resp.Body)
	checkErr(err)
	
	jsonParsed, err := gabs.ParseJSON(responseData)
	checkErr(err)

	gObj, err := jsonParsed.JSONPointer("/0/id")
	checkErr(err)

	readuserid, _ := gObj.Data().(float64)
	userid := int(readuserid)

	fmt.Printf("Found id to delete: %d for username %s\n", userid, line)

	return strconv.Itoa(userid)
}

func deleteUserById(wpurl string, username string, password string, userid string) {
	client := http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest("DELETE", wpurl + "/" + userid + "?reassign=0&force=true", strings.NewReader(""))
	if err != nil {
		exitWithError(errors.New("error building http request"))
	}

	req.SetBasicAuth(username, password)
	
	resp, err := client.Do(req)
	if err != nil {
		exitWithError(fmt.Errorf("http delete request error: %s", err.Error()))
	}

	defer resp.Body.Close()

	fmt.Printf("DELETE : HTTP Status: %d - %s\n", resp.StatusCode, resp.Status)

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

	var line []string
	reader := csv.NewReader(file)

	for {
		line, err = reader.Read()
		if err == io.EOF {
			break
		}
		checkErr(err)

		fmt.Printf("processing: %s\n", line)

		userid := fetchUserId(wpurl, username, password, line[0])

		deleteUserById(wpurl, username, password, userid)
	}

}