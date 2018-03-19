---
title: "Let's Build Something"
date: 2018-02-05T21:15:53-07:00
weight: 20
draft: true
---

Now that Jiffy is installed, lets build a simple service to test it out!

Jiffy's source tree comes with a number of sample model files that you can use to get the hang of things.  We are going to use a simple model file that contains an entity named 'Person'.  The model file can be found in the Jiffy source tree, or pulled directly from the Jiffy github repository.

* $GOPATH/src/github.com/1414C/jiffy/support/testing_models/simpleSingleEntityModel.json
* [jiffy repository](https://github.com/1414C/jiffy/blob/master/support/testing_models/simpleSingleEntityModel.json)

For now we are not going to worry about the content of the model file, but we will look at the structure of our new 'Person' entity briefly. 

```golang

    // Person structure
    type Person struct {
	    ID           uint64   `json:"id" sqac:"primary_key:inc;start:10000000"`
	    Href         string   `json:"href" sqac:"-"`
	    Name         *string  `json:"name,omitempty" sqac:"nullable:true;index:non-unique"`
	    Age          *uint    `json:"age,omitempty" sqac:"nullable:true"`
	    Weight       *float64 `json:"weight,omitempty" sqac:"nullable:true"`
	    ValidLicense *bool    `json:"valid_license,omitempty" sqac:"nullable:true;index:non-unique"`
    }

```

We can see that jiffy will create a Person model with a small set of fields, each with a number of attributes.  For the moment, we will concern ourselves only with the field names and types, as we will need to use this information to construct some tests for the new service.  

{{% notice info %}}
'sqac' tags are used to pass information to jiffy's ORM layer.
{{% /notice %}}

<br/>

Create a new target directory for your project under $GOPATH/src.

```bash

    $ cd $GOPATH/src
    $ mkdir jiffy_tests

```
<br/>

Execute the jiffy binary, specifying the model file to use, as well as the target directory/project name.  

```bash

    $ jiffy -m $GOPATH/src/github.com/1414C/jiffy/support/testing_models/simpleSingleEntityModel.json -p /jiffy_tests/first_app

```

-m tells jiffy which model file should be used to construct the new service.

-p tells jiffy where to write the generated application code. 

{{% notice tip %}}
Note that the location of the folder specified by the -p flag is deemed to be relative to $GOPATH/src.
{{% /notice %}}
 