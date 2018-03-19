---
title: "Login"
date: 2018-02-06T22:34:10-07:00
weight: 100
draft: true
---


Launch Postman and specify a target URL of: `http://127.0.0.1:3000/usr/login` making sure to select the http POST method.  Maintain the request body to provide a user-id and password as shown in the following JSON snippet.  Typically the user-id for a jiffy application is an email address, but we make an exception for the default administration user.

```bash

    {
	    "email": "admin",
	    "password": "initpass"
    }

```
<br/>

When you have finshed and your Postman (or other test utility) looks like the following image, click the 'Send' button to post your login request to the running application. 
![Enter login credentials](../images/login-a.jpg)

<br/>
If all goes well, you will get a http response code of 200 (status ok), and a block of JSON with a single 'token' tag containing a jumble of letters and numbers.  This is the JWT that will be used to validate our authorization to access the Person entity's service end-points.  If you want to read more about JWT's [jwt.io](https://jwt.io) is a good place to start.
![Login](../images/login-b.jpg)
