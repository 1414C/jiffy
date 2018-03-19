---
title: "Get a Person"
date: 2018-02-06T22:34:10-07:00
weight: 140
draft: true
---

What if we need to read a single Person, or isolate a Person entity from a list of entities?  Let's try reading a Person entity using its 'id' key.

Create a new tab in Postman and specify a target URL of `http://127.0.0.1:3000/person/10000001` making sure to select the http GET method.  Next, add the following key-value pair to the http header:

* `Authorization` : `Bearer *paste-your-JWT-here*`

When you have finished, the test session should look as follows and it is time to read Person 10000001 from the database.  Click 'Send' to issue our read request to the application.

![l2-get-person-a](../images/l2-get-person-a.jpg)
<br/>

We just read the Person entity with 'id' key 10000001!  While this is not a very human-friendly way to search for a Person, it is a simple way to isolate and reference an entity for reading, updating or deletion. 

![l2-get-person-b](../images/l2-get-person-b.jpg)
<br/>