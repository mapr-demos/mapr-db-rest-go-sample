---
{
    "authorDisplayName": "Tugdual Grall",
    "category": "mapr-platform",
    "collection": "blog",
    "disqus": {
        "developer": 1,
        "domain": "maprblog",
        "identifier": "getting-started-with-mapr-db-json-rest-api",
        "status": true,
        "title": "Getting started with MapR-DB JSON REST API",
        "url": "https://mapr.com/blog/getting-started-with-mapr-db-json-rest-api/"
    },
    "featuredImage": "assets/mapr.png",
    "layout": "blog-item.html",
    "locale": "en",
    "publish": "2018-04-09T12:00:00.000",
    "title": "Getting started with MapR-DB JSON REST API"
}
---

## Introduction

In this project you will learn how to use the MapR-DB JSON REST API to:

* Create and Delete tables
* Create, Read, Update and Delete documents (CRUD)


MapR Extension Package 5.0 (MEP) introduced the MapR-DB JSON REST API that allow application to use REST to interact with MapR-DB JSON.

You can find information about the MapR-DB JSON REST API in the documentation: [Using the MapR-DB JSON REST API](https://maprdocs.mapr.com/home/MapR-DB/JSON_DB/UsingMapRDBJSONRESTAPI.html)


### Prerequisites

You system should have the following components:

* A running  MapR 6.0.1 & MEP 5.0 cluster with the MapR-DB REST API service installed
* `curl` or equivalent tool


## Discover the MapR-DB JSON REST API

The easiest way to discover it, is to use `curl` command (or equivalent).

**1 - Create a table**

```
curl -X PUT \
  'https://mapr-node:8243/api/v2/table/%2Fapps%2Femp' \
  -u root:mapr \
  -k 

```

In this command:

* the MapR-DB REST Service (MapR Data Access Gateway) is running on the `mapr-node` host with the default port `8243` using HTTPS
* the HTTP verb `PUT` on `/api/v2/table/` endoint creates a new table
* the protocol is HTTP since HTTPS is not enabled on this cluster
* the new table will be created wit the path  `/apps/emp` that is encoded to `%2Fapps%2Femp`
* the user `root` with the password `mapr` is used  for authentication, using basic authentication
* the `-k' parameter is used to indicate to turn off curl's verification of the certificate. 

In this example, you use the basic authentication, it is also possible to use [JSON Web Token](https://jwt.io/introduction/). You will learn more about this when you will write an application in Go.

**2 - Insert Documents**

Insert one document

```
curl -X POST \
  'https://mapr-node:8243/api/v2/table/%2Fapps%2Femp' \
  -u root:mapr \
  -H 'Content-Type: application/json' \
  -d '{"_id":"user001","first_name":"John","last_name":"Doe", "age" : 28}' \
  -k
```

In this command:

* we specify the Content-Type header to be `application/json`
* the HTTP verb `POST` is used to insert or replace a document to the `/apps/emp` table

It is also possible to send multiple document using a JSON array, as you can see in the following command:

Insert multiple documents

```
curl -X POST \
  'https://mapr-node:8243/api/v2/table/%2Fapps%2Femp' \
  -u root:mapr \
  -H 'Content-Type: application/json' \
  -d '[ 
     {"_id":"user002","first_name":"Jane","last_name":"Doe", "age" : 30}, 
     {"_id":"user003","first_name":"Simon","last_name":"Davis", "age" : 43}, 
     {"_id":"user004","first_name":"Paul","last_name":"Duran", "age" : 37}
   ]' \
  -k
```

**3 - Query Documents**

Query document by `_id`

```
curl -X GET \
  'https://mapr-node:8243/api/v2/table/%2Fapps%2Femp%2F/document/user001' \
  -u root:mapr \
  -k
```

In this command:

* the URI `/api/v2/table/{path}/document/{id}` with the HTTP verb `GET` is used to retrieve a single document

If you are trying to retrieve a document that does not exist, the REST API will return a 404 HTTP Code (not found).

Query documents using a condition

You can also query document using condition using the following command:


```
curl -g -X GET \
  'https://mapr-node:8243/api/v2/table/%2Fapps%2Femp?condition={"$eq":{"last_name":"Doe"}}' \
  -u root:mapr \
  -k
```

In this command:

* the `/api/v2/table/{path}` with the verb `GET` is used with a `condition` query parameter
* the OJAI JSON syntax is used to express the condition: `{"$eq":{"last_name":"Doe"}}`


**3 - Update a document**

The following example will increment the age by 1 and update the last name.

```
curl -X POST \
  'https://mapr-node:8243/api/v2/table/%2Fapps%2Femp%2F/document/user001' \
  -u root:mapr \
  -H 'Content-Type: application/json' \
  -d '{"$set" : {"last_name" : "New Doe"}, "$increment" : {"age":1}}' \
  -k 
```

In this comamnd:

* the URL points to the document `_id` to update
* the HTTP verb `POST` is used to modify the resource
* the request body `-d` is the OJAI JSON Mutation that update the last name and increment the age.


You can check that the document has been updated using the following command:

```
curl -X GET \
  'https://mapr-node:8243/api/v2/table/%2Fapps%2Femp%2F/document/user001' \
  -u root:mapr \
  -k
```


**4 -  Delete a document**

Delete the document with the `_id` user001.

```
curl -X DELETE \
  'https://mapr-node:8243/api/v2/table/%2Fapps%2Femp%2F/document/user001' \
  -u root:mapr \
  -k
```


In this command:

* the URI `/api/v2/table/{path}/document/{id}` with the HTTP verb `DELETE` is used to delete the document

**5 - Delete the MapR-DB JSON table**

The last step of this tutorial is to delete the table using the following command:

```
curl -X DELETE \
  'https://mapr-node:8243/api/v2/table/%2Fapps%2Femp' \
  -u root:mapr \
  -k
```


### MapR-DB REST Intro

In this tutorial you have learned how to use the MapR-DB JSON REST API to:

* Create a table
* Insert and query documents
* Update and delete documents
* Drop table


You can now use the API to create MapR-DB JSON Application using your favorite language.





