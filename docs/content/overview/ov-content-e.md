---
title: "Jiffy Development Steps"
date: 2018-02-14T22:11:10-07:00
weight: 50
draft: true
---

### Jiffy pre-generation workflow

Jiffy is intended to generate a clean, straight-forward and secure application services platform just like the one you would write by hand.  Jiffy originally started as a few go templates I used to generate the boiler-plate code that I do not like to type.  Generation is a lot nicer than cut-and-paste.  There are places in the generated code where things could be more elegant, but the code is intended to be easy to work on even if one is not familliar with it.

#### Design

1. You are building a back-end of some sort where you need to reliably read and write data.  Before you start messing around with a real Jiffy model of your application, try some of demo models to get a feel for what Jiffy is and what it isn't.  Its a lot easier to design if you know what the platform provides, and what it does not.

1. Start messing around with entity ideas.  What is the best way to represent your data?  Transfering your ideas into a few rough diagrams often changes your understanding of what you are modelling.

1. Consider the relationships you would like to have between your entities.  Does your model still hold up?

1. Add some detail to your proposed entities in terms of fields.  Remember that Jiffy will insert a primary-key (id), as well as a self-referencing entity href into the model for you.  This is important when you are considering relationships, as Jiffy applications have some expectations regarding the name of the referencing field.  See the Relations section of this documentation set for details.

1. Once you are comfortable with the entities and the fields that they will contain it is time to create the model in the Jiffy model file format.

1. Jiffy model files are simple JSON and at the moment need to be coded up by hand.  Taking a copy of one of the sample model-files is the best way to start.  Short of creating a graphical model file generator, maintaining the models via direct file maintenance is the most direct and transparent way to edit them.  Remember that you do not have to put your entire model in one file.  Jiffy is quite happy to accept any number of model files for your project with the entities defined in any order.  The generator will sort it all out.

1. Check your model files as best as possible before feeding them into Jiffy.  If you get an error, don't worry; we get them all the time too.  Look at the error message, then check the model-file for problems.  If you can't see anything wrong, log an issue and we will take a look.  Usual suspects for errors are missing double-quotes on field-names or string-values, missing closing parentheses on a list, missing brace on struct or sub-struct, putting quotes on a bool value, putting quotes on an int or dec value etc.

1. Again, don't worry too much about formal inter-entity rules in the beginning.  The idea is to get something up and running fast so you can try out the services.  You can generate your application many times over without touching anything but the model-file / model-files.  Generation is fast - less than 5 seconds for most models.

1. Generate a version of your application and start it up.  Try it out with a tool like Postman.  See the Quickstart for an overview of application generation through to application testing, or go to the formal section dealing with each of the appliation creation steps.

1. It is not (typically) neccessary to do early development testing via https.  If this is not the case for you, see the testing with https section in this document set.  Step-by-step instructions are given to create and install self-signed certificates that will permit you to test locally with https.

1. As you test your application, you will probably see some fields missing, some fields that you don't care for, and maybe the need for an additional entity.  Update your model files and generate the application again.  Try it out.  Repeat.

1. As your model becomes more refined, consider formally adding relationships.  Generate, test, repeat.

1. Consider foreign-keys.  Generate, test, repeat.

1. Consider start-values for id.  Not everybody is okay to start at 1, particularly if you are planning on migrating existing data into your application.  Generate and test.

1. Consider static-queries and add the selection options to each field that can be queried.  Generate, test, repeat.  See the static query section for details.

1. Consider indices and add indexes to the model definitions where needed.  Generate, test, repeat.

1. Access restrictions can be assigned at the end-point-level, and this is built into the generated application.  Don't worry about testing authorization and authentication for now.

#### Development

1. When you are satisfied with your model, generate a version of the application and place it under source-code control.

1. Implement extension-points in the controller and model go source as per your requirements.  Implemented extension-points will not be over-written should you need to regenerate the application.  Test often.

1. Implement normalization and validation at the field-level in the generated entity model go source file.  This part is slightly contentious, as field normalization and validiation is performed directly in the generated code, rather than off to the side like extension-points.  We suggest that if you are worried about over-writing your normalization and validation code with an accidental regeneration, you create a *model_normalization* package for each model and implement field-level normalizations and checks there.  Of course, you are free to ignore the provided field-level normalization and validation methods and perform all checks in one of the model extension-points.  Its up to you, and it will make more sense once you take a look through one of the generated files.

1. When you are happy with the way things are working, think about user access.  Sketch out user-groups and assign end-points to them until you have something that you like.  Use Postman to create UserGroups and assign Auths to them, then create new users and allocate the relevant UserGroups to their ID.  That is all you need to do; the Jiffy middleware will take care of the rest.  See the Authentication and Authorization section for a detailed discussion of how users are authenticated and authorized.

#### Deployment

Jiffy generated applications can be deployed in any number of ways.  For example:

* On your laptop
* On a server under your desk
* On a blade running the os of your choice
* On a VM (xen etc.)
* Docker image / Droplet / Cloud Foundry etc.

At the software level, the Jiffy generated application can be deployed as:

* A single application server talking to a database on your SAN
* Multiple application servers talking to a database on your SAN
* Application server and Database server on the same box/image etc.


#### Considerations for deployment

* SSL certificates
* Do you need to support JWT's from other SSO providers?  If so, you will need their public-keys.
* JWT expiration policy
* You have tested your user-access revocation?  Jiffy makes provisions for the revocation of user-access in spite of what some say about JWT.
* Data migration; have you run real test migrations with the key relationships etc. in place?
* If you disabled certain database features for migration, have you turned them back on again?  Looking at you foreign-keys...
* User access;  users have been defined/create/update/delete end-points are locked down in accordance with your use-cases?
* Do you need a reverse-proxy / is there a need to route traffic based on expected load?
* Service activations; is there a need to route traffic based on expected load?
* If you need to scale horizontally or vertically do you have a plan?








Jiffy approaches the API from a services perspective.  Each entity has a corresponding service that can be started when the application initializes.  The Usr, UsrGroup, Auth and GroupAuth services are always generated by default when creating a Jiffy application.  Additional services are generated based on the content of your project's model files.

Generated application services can be broken down into five high-level areas:
<br/>

{{<mermaid align="center">}}
graph TD;
    subgraph 
    A(End-Points)-->B(Middleware)
    B-->C(Controllers)
    C-->D(Models)
    D-->E(Database)
    end
{{< /mermaid >}}
<br/>

* **End-Points** expose the service API's to the consumer such as a web-app or another server.  End-points may be customized by way of the application model files.
* **Middleware** provides user authentication / authorization services and is tightly-coupled to the end-point definitions.  The middleware offers comprehensive services such as authorization via JWT claims inspection, as well as some caching of user and group authorization details.  This is an area of active development.
* **Controllers** are the entry point into the application proper, and are called after a request has been granted access to the end-point by the middleware.  It is here that the body of the request is deserialized and mapped into the correct go model structure.  Extension-points conforming to standard Jiffy interfaces are provided in the controllers for post-generation enhancements.
* **Models** are where the entity data from the request is checked, normalized and prepared for submission to the database.  Extension-points conforming to standard Jiffy interfaces are provided in the models for post-generation enhancements.
* **Database** refers to the backend DBMS that is being used to house the entity data.  Jiffy generated applications can connect to PostgreSQL, MariaDB/MySQL, MSSQL, SAP HanaDB or SQLite.  It is easy to extend the database support to other relational platforms provided that there is an existing go sql driver for the database in question.  It is possible to override the generated Jiffy call to the database and 'roll-your-own' should the need arise.
