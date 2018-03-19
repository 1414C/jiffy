---
title: "Let's Test Something"
date: 2018-02-06T22:40:21-07:00
weight: 50
draft: true
---

We have a running application, but what can we do with it?  If you clicked on the http://127.0.0.1:3000 link at the end of the preceding section, you didn't get a very good looking response...

Jiffy services are best tested using a RESTful test utility.  If you have a tool that works for you, use that to follow along.  If you don't have a test utility, Google's [Postman is a great choice](https://www.getpostman.com/) and that is what we are going to use for the rest of the quick-start.

Let's make a quick list of things that we are going to do in order to test our new Person service.

1. Login
1. Create a new Person entity
1. Create another new Person entity
1. Read each Person by their key
1. Read a list of Person entities
1. Update a Person entity
1. Create yet another new Person entity
1. See what options we can add to an entity request
1. Delete an entity
