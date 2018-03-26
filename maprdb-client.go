package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

type User struct {
	Id         string `json:"_id"`
	Age        int    `json:"age"`
	First_name string `json:"first_name"`
	Last_name  string `json:"last_name"`
}

type Result struct {
	DocumentStream []User `json:"DocumentStream"`
}

type JWToken struct {
	Token string `json:"token"`
}

func main() {

	flag.Parse()

	// retrieve command line parameters
	maprRestServer := *fServer
	username := *fUsername
	password := *fPassword
	table := *fTable
	condition := *fCondition
	drop := *fDrop
	create := *fCreate

	fmt.Println(" \n\n====  Start  Application ====\n Server : ", maprRestServer)

	var token = authenticateUser(maprRestServer, username, password)

	if drop || create {

		action := "create"
		if drop {
			action = "drop"
		}

		// Create or Drop table
		tableOperation(maprRestServer, token, table, action)
		os.Exit(0)
	}

	insertOrReplaceSampleUsers(maprRestServer, token, table)

	querySimpleUser(maprRestServer, token, table, "user003")

	getMultipleUsers(maprRestServer, token, table, condition)

	//Create new User and Insert or Replace it
	newUser := User{Id: "user999", First_name: "Peter", Last_name: "Parker", Age: 23}

	insertOrReplaceUser(maprRestServer, token, table, newUser)

	// Get new user
	querySimpleUser(maprRestServer, token, table, "user999")

	// Update Age to 44
	updateUserAge(maprRestServer, token, table, "user999", 44)

	// Get new user
	querySimpleUser(maprRestServer, token, table, "user999")

	// Delete user
	deleteUser(maprRestServer, token, table, "user999")

	// Print all users
	getMultipleUsers(maprRestServer, token, table, condition)

	fmt.Println(" \n\n====  End of Application ====")

}

func authenticateUser(maprServer string, username string, password string) (token JWToken) {

	var buffer bytes.Buffer
	buffer.WriteString(maprServer)
	buffer.WriteString("/auth/v2/token")

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, buffer.String(), nil)
	req.SetBasicAuth(username, password)

	res, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	bodyContent, err2 := ioutil.ReadAll(res.Body)
	if err2 != nil {
		log.Fatal(err2)
	}

	if res.StatusCode < 200 || res.StatusCode >= 299 {
		fmt.Println("Error ", res.StatusCode, string(bodyContent))
		os.Exit(res.StatusCode)
	}

	json.Unmarshal([]byte(bodyContent), &token)

	return

}

func tableOperation(maprServer string, token JWToken, table string, action string) {

	fmt.Println("\n\n===================================")
	fmt.Println("===      tableOperation()        ==")

	var buffer bytes.Buffer
	buffer.WriteString(maprServer)
	buffer.WriteString("/api/v2/table/")
	buffer.WriteString(url.QueryEscape(table))

	client := &http.Client{}

	httpVerb := http.MethodPut

	if action == "drop" {
		httpVerb = http.MethodDelete
	}

	req, err := http.NewRequest(httpVerb, buffer.String(), nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.Token))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		log.Fatal(err)
	}

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	bodyContent, err2 := ioutil.ReadAll(res.Body)
	if err2 != nil {
		log.Fatal(err2)
	}

	if res.StatusCode < 200 || res.StatusCode >= 299 {
		fmt.Println("Error ", res.StatusCode, string(bodyContent))
		os.Exit(res.StatusCode)
	}

	fmt.Printf("\t Table %s\n===================================\n", action)

}

func insertOrReplaceSampleUsers(maprServer string, token JWToken, table string) {

	fmt.Println("\n\n===================================")
	fmt.Println("===  insertOrReplaceSampleUsers() ==")

	// create an array of users

	userList := [4]User{
		{Id: "user001", First_name: "John", Last_name: "Doe", Age: 28},
		{Id: "user002", First_name: "Jane", Last_name: "Doe", Age: 30},
		{Id: "user003", First_name: "Simon", Last_name: "Davis", Age: 43},
		{Id: "user004", First_name: "Paul", Last_name: "Duran", Age: 37}}

	jsonObject, errorJSON := json.Marshal(userList)
	if errorJSON != nil {
		log.Fatal(errorJSON)
	}

	var buffer bytes.Buffer
	buffer.WriteString(maprServer)
	buffer.WriteString("/api/v2/table/")
	buffer.WriteString(url.QueryEscape(table))

	client := &http.Client{}

	req, err := http.NewRequest(http.MethodPost, buffer.String(), bytes.NewBuffer(jsonObject))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.Token))
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	bodyContent, err2 := ioutil.ReadAll(res.Body)
	if err2 != nil {
		log.Fatal(err2)
	}

	if res.StatusCode < 200 || res.StatusCode >= 299 {
		fmt.Println("Error ", res.StatusCode, string(bodyContent))
		os.Exit(res.StatusCode)
	}

	fmt.Println("==    Users inserted/updated    ==")
	fmt.Println("===================================")

}

func querySimpleUser(maprServer string, token JWToken, table string, id string) {

	fmt.Println("\n\n===================================")
	fmt.Println("===       querySimpleUser()      ==")

	var buffer bytes.Buffer
	buffer.WriteString(maprServer)
	buffer.WriteString("/api/v2/table/")
	buffer.WriteString(url.QueryEscape(table))
	buffer.WriteString("/document/")
	buffer.WriteString(id)

	client := &http.Client{}

	req, err := http.NewRequest(http.MethodGet, buffer.String(), nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.Token))

	res, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	bodyContent, err2 := ioutil.ReadAll(res.Body)
	if err2 != nil {
		log.Fatal(err2)
	}

	if res.StatusCode < 200 || res.StatusCode >= 299 {
		fmt.Println("Error ", res.StatusCode, string(bodyContent))
		os.Exit(res.StatusCode)
	}

	var user001 User

	json.Unmarshal([]byte(bodyContent), &user001)

	fmt.Printf("Id : %s \n", user001.Id)
	fmt.Printf("First Name : %s \n", user001.First_name)
	fmt.Printf("Age : %d \n", user001.Age)

	fmt.Println("===================================")

}

func getMultipleUsers(maprServer string, token JWToken, table string, condition string) {

	fmt.Println("\n\n===================================")
	fmt.Println("===       getMultipleUser()      ==")

	var buffer bytes.Buffer
	buffer.WriteString(maprServer)
	buffer.WriteString("/api/v2/table/")
	buffer.WriteString(url.QueryEscape(table))

	if condition != "" {
		buffer.WriteString("?condition=")
		buffer.WriteString(url.QueryEscape(condition))
		fmt.Println("\t condition ", condition)
	}

	fmt.Println(buffer.String())

	client := &http.Client{}

	req, err := http.NewRequest(http.MethodGet, buffer.String(), nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.Token))

	res, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	var result Result

	bodyContent, err2 := ioutil.ReadAll(res.Body)
	if err2 != nil {
		log.Fatal(err2)
	}

	if res.StatusCode < 200 || res.StatusCode >= 299 {
		fmt.Println("Error ", res.StatusCode, string(bodyContent))
		os.Exit(res.StatusCode)
	}

	json.Unmarshal([]byte(bodyContent), &result)

	fmt.Println("result %s ", result)

	for _, userDoc := range result.DocumentStream {
		// element is the element from someSlice for where we are

		fmt.Printf("User %s : [ Name :  %s %s , Age : %d  ]\n", userDoc.Id, userDoc.First_name, userDoc.Last_name, userDoc.Age)
	}

	fmt.Println("===================================")

}

func insertOrReplaceUser(maprServer string, token JWToken, table string, newUser User) {

	fmt.Println("\n\n===================================")
	fmt.Println("===       insertOrReplaceUser()      ==")

	var buffer bytes.Buffer
	buffer.WriteString(maprServer)
	buffer.WriteString("/api/v2/table/")
	buffer.WriteString(url.QueryEscape(table))

	jsonObject, errorJSON := json.Marshal(newUser)
	if errorJSON != nil {
		log.Fatal(errorJSON)
	}

	client := &http.Client{}

	req, err := http.NewRequest(http.MethodPost, buffer.String(), bytes.NewBuffer(jsonObject))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.Token))
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	bodyContent, err2 := ioutil.ReadAll(res.Body)
	if err2 != nil {
		log.Fatal(err2)
	}

	if res.StatusCode < 200 || res.StatusCode >= 299 {
		fmt.Println("Error ", res.StatusCode, string(bodyContent))
		os.Exit(res.StatusCode)
	}

	fmt.Println("===================================")

}

func updateUserAge(maprServer string, token JWToken, table string, id string, age int) {

	fmt.Println("\n\n===================================")
	fmt.Println("===       updateUserAge()      ==")

	var buffer bytes.Buffer
	buffer.WriteString(maprServer)
	buffer.WriteString("/api/v2/table/")
	buffer.WriteString(url.QueryEscape(table))
	buffer.WriteString("/document/")
	buffer.WriteString(id)

	var bufferMutation bytes.Buffer
	bufferMutation.WriteString("{\"$set\":{\"age\":")
	bufferMutation.WriteString(fmt.Sprintf("%v", age))
	bufferMutation.WriteString("}}")

	fmt.Println(buffer.String())
	fmt.Println(bufferMutation.String())

	client := &http.Client{}

	req, err := http.NewRequest(http.MethodPost, buffer.String(), bytes.NewReader(bufferMutation.Bytes()))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.Token))
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	bodyContent, err2 := ioutil.ReadAll(res.Body)
	if err2 != nil {
		log.Fatal(err2)
	}

	if res.StatusCode < 200 || res.StatusCode >= 299 {
		fmt.Println("Error ", res.StatusCode, string(bodyContent))
		os.Exit(res.StatusCode)
	}

	fmt.Println("===================================")

}

func deleteUser(maprServer string, token JWToken, table string, id string) {

	fmt.Println("\n\n===================================")
	fmt.Println("===       updateUserAge()      ==")

	var buffer bytes.Buffer
	buffer.WriteString(maprServer)
	buffer.WriteString("/api/v2/table/")
	buffer.WriteString(url.QueryEscape(table))
	buffer.WriteString("/document/")
	buffer.WriteString(id)

	client := &http.Client{}

	req, err := http.NewRequest(http.MethodDelete, buffer.String(), nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.Token))
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	bodyContent, err2 := ioutil.ReadAll(res.Body)
	if err2 != nil {
		log.Fatal(err2)
	}

	if res.StatusCode < 200 || res.StatusCode >= 299 {
		fmt.Println("Error ", res.StatusCode, string(bodyContent))
		os.Exit(res.StatusCode)
	}

	fmt.Println("===================================")

}

var fServer = flag.String("server", "http://localhost:8085", "MapR-DB REST API Server")
var fUsername = flag.String("user", "mapr", "Username")
var fPassword = flag.String("password", "", "Password")
var fTable = flag.String("table", "/apps/employee", "Table path")
var fCondition = flag.String("condition", "", "OJAI JSON Condition used by  getMultipleUsers function")
var fCreate = flag.Bool("create", false, "Create Table")
var fDrop = flag.Bool("drop", false, "Drop Table")
