---
title: "Delete a Person"
date: 2018-02-06T22:34:10-07:00
weight: 170
draft: true
---


Time to delete one of our Person entities!  Create a new tab in Postman and specify a target URL of `http://127.0.0.1:3000/person/10000000` with the http DELETE method.  Next, add the following key-value pairs to the http header:

* `Content-Type`   : `application\json`
* `Authorization` : `Bearer *paste-your-JWT-here*`

When you have something that looks as follows, click the 'Send' button to issue the delete request to the application.

![delete-person-a](../images/delete-person-a.jpg)

If the delete request was successful, you will see a http response-code of 202 (Accepted).  Try to read the entity by converting your delete request into a get request and verify that Person entity 10000000 has truly been deleted.

![delete-person-b](../images/delete-person-b.jpg)