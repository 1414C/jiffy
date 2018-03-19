---
title: "Let's Run Something"
date: 2018-02-05T21:40:32-07:00
weight: 30
draft: true
---

If everything has gone according to plan, a new application has been generated and should be ready to run. Let's take a look at what jiffy generated for us!

```bash

    $ cd $GOPATH/src/jiffy_tests/first_app
    $ ls -l

```

This should result in a list of the files and folders comprising our new application.

```bash

    drwxr-xr-x  12 stevem  staff    384  5 Feb 21:44 .
    drwxr-xr-x  10 stevem  staff    320  5 Feb 21:44 ..
    -rwxr-xr-x   1 stevem  staff    588  5 Feb 21:44 .dev.config.json
    -rwxr-xr-x   1 stevem  staff    611  5 Feb 21:44 .prd.config.json
    drwxr-xr-x   4 stevem  staff    128  5 Feb 21:44 appobj
    drwxr-xr-x  10 stevem  staff    320  5 Feb 21:44 controllers
    drwxr-xr-x   4 stevem  staff    128  5 Feb 21:44 jwtkeys
    -rwxr-xr-x   1 stevem  staff    839  5 Feb 21:44 main.go
    -rwxr-xr-x   1 stevem  staff  18421  5 Feb 21:44 main_test.go
    drwxr-xr-x   3 stevem  staff     96  5 Feb 21:44 middleware
    drwxr-xr-x  12 stevem  staff    384  5 Feb 21:44 models
    drwxr-xr-x   3 stevem  staff     96  5 Feb 21:44 util
    $

```

Jiffy generates two sample configuration files each time it is executed.  We are going to run our application with development environment settings, so lets take a quick look at .dev.config.json to make sure there are no horrible surprises.

```json

    {
        "port": 8080,    
        "env": "dev",     
        "pepper": "secret-pepper-key",  
        "database": {
            "db_dialect": "sqlite",
            "host":       "127.0.0.1",
		    "port":       0,
		    "usr":        "",
		    "password":   "",
		    "name":       "testdb.sqlite"
        },
        "cert_file": "",
        "key_file": "",
        "jwt_priv_key_file": "jwtkeys/private.pem",
        "jwt_pub_key_file": "jwtkeys/public.pem",
        "service_activations": [
            {
                "service_name":   "Person",
                "service_active": true
            }
            ]
    }

```

Jiffy decides by default to run against a sqlite database and generates what *should* be a suitable configuration file for most systems.  Again, we are not going to worry too much about what is in the configuration file at this point, but know that database file 'testdb.sqlite' will be created in the generated application's root folder.

<br/>

Let's execute our first application!

```bash

    $ cd $GOPATH/src/jiffy_tests/first_app
    $ go run main.go -dev -rs

```

Exectuing with the -dev and -rs flags instructs our new application to initialize itself using the development settings file, and forces a rebuild of the 'Super' authorization-group.  Consequently, you will see some warning and info messages scroll up the screen which is perfectly normal.  

```bash

    2018/02/05 22:29:26 package sqac init is running
    2018/02/05 22:29:26 successfully loaded the config file...
    2018/02/05 22:29:26 JWTPrivKeyFile: jwtkeys/private.pem
    2018/02/05 22:29:26 JWTPubKeyFile: jwtkeys/public.pem
    2018/02/05 22:29:26 warning: auth usr.GET_SET not found in the db Auth master data
    2018/02/05 22:29:26 warning: auth usr.CREATE not found in the db Auth master data
    2018/02/05 22:29:26 warning: auth usr.GET_ID not found in the db Auth master data
    ...
    ...
    2018/02/05 22:29:26 info: creating auth usr.GET_SET in the db Auth master data
    2018/02/05 22:29:26 warning: new auth usr.GET_SET must be added to at least one group
    2018/02/05 22:29:26 info: creating auth usr.CREATE in the db Auth master data
    2018/02/05 22:29:26 warning: new auth usr.CREATE must be added to at least one group
    ...
    ...
    2018/02/05 22:29:26 warning: new auth person.STATICFLTR_ByValidLicense must be added to at least one group
    2018/02/05 22:29:26 info: creating auth person.STATICFLTR_CMD_ByValidLicense in the db Auth master data
    2018/02/05 22:29:26 warning: new auth person.STATICFLTR_CMD_ByValidLicense must be added to at least one group
    2018/02/05 22:29:26 The Super UsrGroup has been initialized with 42 Auth objects.
    2018/02/05 22:29:26 re-initializing local middleware to accomodate Super group changes.
    2018/02/05 22:29:27 admin user created with ID: 1 and initial password of initpass
    Development settings selected...
    Starting http server on port...  8080

```

You should see a message indicating that the application is running.  During the startup, the application executed a number of steps.

1. Loaded the development configuration file.
1. Initialized a handle to the underlying ORM.
1. Checked for and loaded the public and private keys for JWT support.
1. Checked for and created the user, auth, and usergroup tables in the database.
1. Checked for and created the person table in the database based on the Person model.
1. Checked for and created authorizations in the database for each service end-point.
1. Checked for and created the Super user-group in the database.
1. Assigned all authorizations to the Super user-group.
1. Created the 'admin' user and assigned it to the 'Super' user-group.
1. Initialized the authorization cache in the router.
1. Started the router.

Congratulations!  Your first application is now open for business at http://127.0.0.1:3000 


 