---
title: "Get Some Persons"
date: 2018-02-06T22:34:10-07:00
weight: 130
draft: true
---


So far we have created two Person entities.  We have observed that upon successful creation of an entity, a JSON representation of that entity is passed back to us via the response-writer.  Let's now look at how we can get a list of all of our Person entities.

Create a new tab in Postman and specify a target URL of `http://127.0.0.1:3000/persons` making sure to select the http GET method.  Next, add the following key-value pair to the http header:

* `Authorization` : `Bearer *paste-your-JWT-here*`

When you have finished, the test session should look as follows and it is time to read some Person entities from the database.  Click 'Send' to issue our read request to the application.

![persons-a](../images/l2-get-persons-a.jpg)
<br/>

We just read the complete list of Person entities!  Adding an 's' to the entity name and issuing the request with a GET http verb tells jiffy to read all of the person entities.  In some cases this looks odd, but it makes it quite easy to consume the services.  Notice that the 'href' field of each Person entity provides a direct link to the entity that it is a part of.

![person-b](../images/l2-get-persons-b.jpg)
<br/>


