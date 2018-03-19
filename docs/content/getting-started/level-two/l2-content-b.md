---
title: "Create a Person"
date: 2018-02-06T22:34:10-07:00
weight: 110
draft: true
---


Now that we have successfully logged into the application and received our first JWT, it is time to create a new Person entity.  Start by copying the content of the 'token' tag from the login response body to the clipboard.  This JWT must henceforth be included in the http header of every subsequent request.

Create a new tab in Postman and specify a target URL of `http://127.0.0.1:3000/person` with the http POST method.  Next, add the following key-value pairs to the http header:

* `Content-Type`   : `application\json`
* `Authorization` : `Bearer *paste-your-JWT-here*`

![person-a](../images/person-a.jpg)

When you have finished maintaining the http header-values, click on 'Body' and maintain it using the 'raw' setting.  This will allow you to paste the following JSON code snippet into the request's body:

```JSON

    {
	    "name": "Steve Dallas",
	    "age": 46,
	    "weight": 185,
	    "valid_license": true,
	    "license_class": "A"
    }

```
<br/>

When you have finished, the test session should look as follows and it is time to create our first Person entity.  Click 'Send' to post the new entity to the application!

![person-b](../images/person-b.jpg)
<br/>

Congratulations!  You have created your first Person entity!

![person-c](../images/person-c.jpg)

The router matched the request URL to a route (service end-point), the middleware layer in the matched route examined the JWT, verified it was okay to proceed and then passed the raw JSON from the request body to the Person entity's controller.  The controller deserialized the JSON into a Person struct and then passed the result to the Create method in the model/validation layer.  Validation of the Person struct's content occured, and then a call was made to the underlying ORM to create the entity on the database!  

The ORM-layer returned the new entity to the application's model-layer, where it was checked and passed back to the controller layer, whereupon it was serialized (struct content becomes JSON) and written to the the response-writer.

This is a high-level view of what transpired, but the *general* flow of things is accurate.  

Notice that the entity passed back to us seems to have a couple of extra fields?  All entities created via a jiffy model file are injected with a primary-key of 'id' as well as a non-persistent 'href' field.  In this example, our entity's 'id' field was specified to be auto-incrementing with a starting value of 10000000.  See the sqac-tag section in the documentation for details regarding key options.

Href is included in each entity's GET responses, and acts as a self-reference providing the entity's direct access URI to the consumer.
