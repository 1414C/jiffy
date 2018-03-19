---
title: "Update a Person"
date: 2018-02-06T22:34:10-07:00
weight: 150
draft: true
---

If you have been following along, we have created Person entities, read Person entities both in bulk and by 'id' key.  We are now going to take a look at updating an existing entity.  

Read Person entity 10000001 from the database as shown in the [Get a Person](../l2-content-e) example.

Once you have successfully read Person entity 10000001 into your Postman session, copy the content of the response body to your clipboard.  We are going to reuse our GET Person Postman tab to perform an update, so change the http verb from GET to PUT.  Changing the http verb to PUT will allow us to maintain the request-body in our Postman tab.  Next, add the following key-value pair to the http header:

* `Content-Type`   : `application\json`

Now paste the GET Person response-body into the request-body of our new PUT request, then edit it so that Opus's full name is given, as well as his correct weight.  Strictly speaking, you do not need to include the 'id' or 'href' fields in an update, but it does not hurt anything to do so.

```json

    {
        "id": 10000001,
        "href": "http://127.0.0.1:8080/person/10000001",
        "name": "Opus the Penguin",
        "age": 8,
        "weight": 385,
        "valid_license": false
    }

```
<br/>

When you have finished, you should have something that looks as follows.  Click 'Send' to issue the PUT request to the application.

![l2-update-person-a](../images/l2-update-person-a.jpg)
<br/>

The PUT request should update the entity using the 'id' key as its update criteria, and then return a JSON representation of the updated entity.  This is a little different than other approaches, where the result of an update is measured simply by the http response code.  If you don't like the way this works, it is very easy to update the generated source code to omit the response body following a PUT.

![l2-update-person-b](../images/l2-update-person-b.jpg)


