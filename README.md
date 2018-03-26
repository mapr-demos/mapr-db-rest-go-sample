# Create a MapR-DB application in Go using REST API


## Introduction

In this project you will learn how to use the MapR-DB JSON REST API using Go

MapR Extension Package 5.0 (MEP) introduced the MapR-DB JSON REST API that allow application to use REST to interact with MapR-DB JSON.

You can find information about the MapR-DB JSON REST API in the documentation: [Using the MapR-DB JSON REST API](https://maprdocs.mapr.com/home/MapR-DB/JSON_DB/UsingMapRDBJSONRESTAPI.html)


### Prerequisites

You system should have the following components:

* A running  MapR 6.0.1 & MEP 5.0 cluster with the MapR-DB REST API service installed
* [Go](http://golang.org)
* Git


## Run your first Go/MapR-DB Application

**1- Get the source and Build the Application**

Cline and build the repository using the following commands:

```bash
$ git clone https://github.com/mapr-demos/mapr-db-rest-go-sample.git

$ cd mapr-db-rest-go-sample

$ go build -o mapr-db-go

```

You have now a Go application named `mapr-db-go`.


**2- Run the Application**

The `mapr-db-go` is a simple command line application.

2.1 List of parameters

You can look at the options using:

```
mapr-db-go -help

Usage of mapr-db-go:
  -condition string
        OJAI JSON Condition used by  getMultipleUsers function
  -create
        Create Table
  -drop
        Drop Table
  -password string
        Password
  -server string
        MapR-DB REST API Server (default "http://localhost:8085")
  -table string
        Table path (default "/apps/employee")
  -user string
        Username (default "mapr")

```

2.1 Create a table

The following command will create a table if does not already exist, and return an error if the table already exists.

```
$ mapr-db-go --server http://mapr-node:8085 --user mapr --password mapr --table /apps/emp --create
```

2.2 Run the full application

```
$ mapr-db-go --server http://mapr-node:8085 --user mapr --password mapr --table /apps/emp 
```

This command will run all operations contains in the application:

* Insert or replace User documents
* Query a single user by its `_id`
* Get multiple users without condition
* Insert a new user
* Update a user
* Delete a user

2.3 Query with a specific OJAI Condition

```
$ mapr-db-go --server http://mapr-node:8085 --user mapr --password mapr --table /apps/emp --condition '{"$eq":{"last_name":"Doe"}}'
```

This command will run the same commands as before, but the it will return the multiple users based on the condition.


2.4 Drop the table

You can drop the table usinf the following command:

```
$ mapr-db-go --server http://mapr-node:8085 --user mapr --password mapr --table /apps/emp --drop
```


## A quick look to the Go application code

The `main()` capture the command line parameters using the [`flag`](https://golang.org/pkg/flag/) library and call various functions.

Before doing any operation on MapR-DB tables, the application must create an authentication token, this is done in the [authencateUser](https://github.com/mapr-demos/mapr-db-rest-go-sample/blob/master/maprdb-client.go#L89) function.

### Authenticate user function

```golang
func authenticateUser(maprServer string, username string, password string) (token JWToken) {
    ...
    client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, buffer.String(), nil)
    req.SetBasicAuth(username, password)
    ...
  	bodyContent, err2 := ioutil.ReadAll(res.Body)

	if res.StatusCode < 200 || res.StatusCode >= 299 {
		fmt.Println("Error ", res.StatusCode, string(bodyContent))
		os.Exit(res.StatusCode)
	}

	json.Unmarshal([]byte(bodyContent), &token)

    return
}
```

The authenticateUser function has the following signature:

* `maprServer`: the URL of the MapR-RB REST API
* `username`: the MapR user
* `password`: the MapR user password
* `token`: the return value, a JSON Web Token that will be used in sub sequent HTTP requests.

The `buffer.String()` contains the REST API URL for authentication, that looks like: ` http://mapr-node:8085/auth/v2/token`


### Insert documents in MapR-DB Table


The signature of the `insertOrReplaceSampleUsers` function is:

```golang
func insertOrReplaceSampleUsers(maprServer string, token JWToken, table string) {
...
}
```

where:

* `maprServer`: the URL of the MapR-RB REST API
* `token`: is the JSON Web Token use to authenticate the API Call
* `table`: the path of the table

These 3 parameters are present in all functions of the application; depending of the needs of the function you may have additional parameter, for example a user `_id` or a User object.

Go language has a native library to marshall and unmarshall JSON ojbect into Go structure, so the first thing to do is to create a structure that match the JSON Documents that represents the User. This structure looks like:

```golang
type User struct {
	Id         string `json:"_id"`
	Age        int    `json:"age"`
	First_name string `json:"first_name"`
	Last_name  string `json:"last_name"`
}
```

The application instianciates the object using the following code present in the `insertOrReplaceSampleUsers` function:

```golang

userList := [4]User{
		{Id: "user001", First_name: "John", Last_name: "Doe", Age: 28},
		{Id: "user002", First_name: "Jane", Last_name: "Doe", Age: 30},
		{Id: "user003", First_name: "Simon", Last_name: "Davis", Age: 43},
		{Id: "user004", First_name: "Paul", Last_name: "Duran", Age: 37}}
```

This create an Array of 4 Users.

You can now marshall this list of objects into a JSON object using the following code:

```golang
jsonObject, errorJSON := json.Marshal(userList)
if errorJSON != nil {
    log.Fatal(errorJSON)
}
```

The `jsonObject` can now be used in the MapR-DB JSON HTTP request 

```golang
var buffer bytes.Buffer
buffer.WriteString(maprServer)
buffer.WriteString("/api/v2/table/")
buffer.WriteString(url.QueryEscape(table))

client := &http.Client{}

req, err := http.NewRequest(http.MethodPost, buffer.String(), bytes.NewBuffer(jsonObject))
req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.Token))
req.Header.Set("Content-Type", "application/json")

res, err := client.Do(req)

```

Let's look at the code above in detail:

* The first 4 lines are used to create the URL to use to post new document into the table. This will look like `http://mapr-node:8085/api/v2/table/%2Fapps%2Femp`
* Then the HTTP request is prepared using `http.NewRequest(http.MethodPost, buffer.String(), bytes.NewBuffer(jsonObject))`
    * The first parameter is set do an HTTP POST
    * The second parameter is the URL of the operation
    * the third parameter is the body of the request, in this case the JSON that contains an array of User
*  The request is used to set HTTP Headers:
    * Use the JWT Token in the Authorization header: `req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.Token))`
    * Set the content type of the request to `application/json'


The application use the same business logic to insert multiple document (`insertOrReplaceSampleUsers()`), or to insert a single document (`insertOrReplaceUser()`).



### Find Documents

You know now how to authenticate, and how to insert documents. Let's see how you can query documents in Go using the REST API.

The `getMultipleUsers()` and `querySimpleUser()` are similar, si let's focus on the method that uses OJAI Conditions. 


```golang
func getMultipleUsers(maprServer string, token JWToken, table string, condition string) {
...
}
```


where:

* `maprServer`: the URL of the MapR-RB REST API
* `token`: is the JSON Web Token use to authenticate the API Call
* `table`: the path of the table
* `condition`: is a OJAI Condition, that if present will be used in the URL call

The HTTP call to get the users from the JSON Table is done using the following code:

```golang
...
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
...
```

The first 4 lines are used to create the URL to get documents from the table, something that will look like `http://mapr-node:8085/api/v2/table/%2Fapps%2Femp`

The if a condition is present, the query parameter is added to build something like:

```
http://mapr-node:8085/api/v2/table/%2Fapps%2Femp?condition={"$eq":{"last_name":"Doe"}}
```

Then you just have to use the Go HTTP library to do the call.


Once the call is done you can parse the result using the following code:

1- Create a Go structure that will be used to represent the result of a MapR-DB REST call:

```golang
type Result struct {
	DocumentStream []User `json:"DocumentStream"`
}
```

2- Create an object from the HTTP result:

```golang

var result Result
...
bodyContent, err2 := ioutil.ReadAll(res.Body)
...
json.Unmarshal([]byte(bodyContent), &result)

for _, userDoc := range result.DocumentStream {
	// element is the element from someSlice for where we are
	fmt.Printf("User %s : [ Name :  %s %s , Age : %d  ]\n", userDoc.Id, userDoc.First_name, userDoc.Last_name, userDoc.Age)
}

```

* Once you have read the Body of the response you can umarshall it to create a Go object using `json.Unmarshal([]byte(bodyContent), &result)
`
* Then you can simply iterate over the list of user using `for _, userDoc := range result.DocumentStream`

### Other operations

The application use other function to:

* Create and Drop Table
* Update a document
* Delete a document.

The Go code is very similar to what you have seen on the operations detailed previously. The differences are mostly, the URI and the HTTP verb used to call the REST API from your code.



## Conclusion

In this article you have learned how to call the MapR-DB REST API from Go language, and manipulates documents and tables.

You can now create a richer application in Go, or any other programming language that provide good HTTP and JSON libraries.


