# rgen

## Overview
A simple code generation utility to create RESTful services with a Postgres DB backend.
<br/>

## Work-In-Progress
1.  Consider simply enforcing the use of ID & Href in the model definition as explicitly named fields.
2.  Fully implement nullable / pointer support
3.  Ensure that rune and byte types are fully accounted for

## Features
* login / session management via jwt
* built-in support for the creation of signing-keys for jwt
* JSON configuration (model) file for Entity, Index and Relationship definitions
* automatically creates backend database artifacts based on the model file (tables, indices)
* supports single and composite index declarations via the model file
* built-in support for https
* baked in normalization and validation
* generates a working set of CRUD-type RESTful services based on the model file
* generates working query end-points based on the model fie 
* generates a comprehensive set of working tests (go test)
* generated code is easily extended
<br/>

## Installation and Execution
In order to run the application generator, ensure the following:

1.  From $GOPATH/src, use go get to install the following:
    * go get -u github.com/gorilla/mux
    * go get -u github.com/dgrijalva/jwt-go
    * go get -u github.com/golang.org/x/crypto/bcrypt
    * go get -u github.com/1414C/sqac
    * go get -u github.com/lib/pq
    * go get -u github.com/SAP/go-hdb/driver
    * go get -u github.com/go-sql-driver/mysql
    * go get -u github.com/mattn/go-sqlite3
    * go get -u github.com/MSSQL

2.  Install the application into your local $GOPATH/src directory:
    * go get -u github.com/1414C/rgen

3.  You will need access to a Postgres, MySQL, MSSQL or SAP Hana database, either locally or over the network.  It is also possible to run tests with SQLite3.
    
4.  The application can be started in two ways:
    * From $GOPATH/src/github.com/1414C/rgen you may execute the application by typing:
        * go run main.go     
    * A binary can also be build from $GOPATH/src/github.com/1414C/rgen by typing the following:
        * go build .
        * The application can then be started from the same directory by typing:
            * ./rgen
<br/>

## Flags
Flags are generally not used, as the configuration files (models.json) are easier to deal with.  There are however, a few flags that can be appended to the execution command:
* go run *.go -p
	* The -p switch is used to specify the target directory for generated application source-code relative to $GOPATH/src.

```bash

    $ go run main.go -p "github.com/footle.com/myrestfulsvc"

```

* go run main.go -m "./my_model.json"
    * By default, the application will attempt to use ./models.json as the model source, but inclusion of the -m flag permits the use of an alternate model file.
    * The path of model file in the application base directory must be prefaced with ./ .  If the model file is not located in the base directory of the application, the full path must be specified when using the -m flag.

---
<br/>

## Model Creation
Create a model file containing the Entities, Indexes and Relations that you wish to generate services for.  Entity model defintion consists of an array of JSON objects, with each object being limited to a flat hierarchy and basic go-data-types.  By default, the generator expects a *models.json* file in the execution directory, but a correctly formatted JSON file can be loaded from any location by executing with the *-m* flag.  

 A sample models.json file is installed with the application and can be found in the root application folder, as shown below:

models.json
```JSON
[
    {
        "type": "entity",
        "typeName": "Person",
        "properties": {
            "name": {
                "type": "string",
                "format": "", 
                "required": false,
                "unique": false,
                "index": "nonUnique",
                "selectable": "eq,like"
            },
            "age": {
                "type": "uint",
                "format": "", 
                "required": false,
                "unique": false,
                "index": "",
                "selectable": "eq,lt,gt"
            },
            "weight": {
                "type": "float64",
                "format": "", 
                "required": false,
                "unique": false,
                "index": "",
                "selectable": "eq,lt,le,gt,ge"
            },
            "validLicense": {
                "type": "bool",
                "format": "", 
                "required": false,
                "unique": false,
                "index": "nonUnique",
                "selectable": "eq,ne"
            },
            "homeAddressID": {
                "type": "uint",
                "format": "", 
                "required": false,
                "unique": false,
                "index": "nonUnique",
                "selectable": "eq",
                "relation": "Address",
                "relationFld": "addressID",
                "relationCrd": "1:1"
            }
        },
        "compositeIndexes": [ 
            {"index": "name, age"}, 
            {"index": "validlicense, age"} 
        ]
    },
    {
        "type": "entity",
        "typeName": "Address",
        "properties": {
            "city": {
                "type": "string",
                "format": "", 
                "required": true,
                "index": "nonUnique",
                "selectable": "eq,like"
            },
            "street": {
                "type": "string",
                "dbType": "varchar(100)",
                "format": "", 
                "required": false,
                "index": "nonUnique",
                "selectable": "eq,like"
            },
            "streetNumber": {
                "type": "string",
                "format": "", 
                "required": false,
                "index": "",
                "selectable": ""
            },
            "postCode": {
                "type": "string",
                "format": "", 
                "required": false,
                "index": "nonUnique",
                "selectable": "eq,like"
            }
        }
    }
]
```
<br/>

### Model Entity Components

A subset of the sample Entity model is shown along with a description of each field.
```code
[
    {
        "type": "entity",
        Field 'type' refers to the category of object the block of JSON represents.  
        Valid values include {"entity", "relation" or "index"}.
        This is a mandatory model element.

        "typeName": "Person",
        Field 'typeName' refers to the formal and common name of the object type specified in the 
        current block's 'type' field.  An Entity given a typeName of "Person" will result in an 
        internal model object of type Person.
        This is a mandatory model element.
        
        "properties": {
        Field 'properties' refers to the content of the 'type'-'typeName' combination.
        In the context of the "Person" "entity" 'type', 'properties' refer the data fields of 
        generated "Person" model structure.  'properties' are a collection of free-form name-tags,
        each with a child-block containing the 'property' attributes.
        This is a mandatory model element.

        The "name" property block is described below:

            "name": {
                "type": "string",
                Field 'type' in a 'properties'->'name' JSON-clock refers to the go data-type associated
                with the current 'property'.
                'type' is a mandatory field in an "entity" 'property' block.

                "dbtype": "varchar(100)",
                Field 'dbtype' can be used to specify a native db-field-type for the current 'property'.
                This is an optional field, and the cast to the DB-type is handled in the ORM layer.

                "format": "", 
                Field 'format' is not currently used, but is intended to deal with field conversion from
                strings / floats to timestamp formats etc.
                This is an optional field.

                "required": false,
                Field 'required' can be used to instruct the database that the current 'property' is a 
                required field in its related db table column.
                Allowed values include {true, false}.
                This is a mandatory field.

                "unique": false,
                Field 'unique' can be used to instruct the database not to accept duplicate values in the
                database column related to the current 'property'.
                Allowed values include {true, false}.
                This is a mandatory field.

                "index": "nonUnique",
                Field 'index' can be used to instruct the database to create an index on the db table-column
                related to the current 'property'.
                See the 'indexes' element in the type definition for the creation of compound indices.
                Allowed values include {"unique", "nonUnique"}.
                This is an optional field.

                "selectable": "eq,like"
                Field 'selectable' can be used to instruct the code-generator to create simple REST query 
                accessor routes for the current 'property'.  The generator creates routes to permit GET
                operations that can be called based on the entity 'typeName' and 'property' values.
                Allowed values include {"EQ", "eq", "LT", "lt", "GT", "gt", "GE", "ge", "LIKE", "like", "NE", "ne"}
                Additional restrictions are imposed based on the 'type' field value.  For example, a bool
                type need not support LT or GT operators.
                Sample routes for Person->Name selection with "eq,like" are shown:
                    
                    https://localhost:<port>/persons/name(EQ '<sel_string>')
                    https://localhost:<port>/persons/name(LIKE '<sel_string>')

            },
            "age": {
                "type": "uint",
                "format": "", 
                "required": false,
                "unique": false,
                "index": "",
                "selectable": "eq,lt,le,gt,ge,ne"
            },
            ...
            ...
        },
        Composite indices for an Entity can be declared in the 'compositeIndexes' JSON array element.  Specify the 
        column names in the order that they should appear in the composite index, making sure to match the case etc.
        "compositeIndexes": [ 
            {"index": "name, age"}, 
            {"index": "validlicense, age"} 
        ]
    },    
    {
        ... next entity definition
    }
]
```

# Using the Generated Code

1.  Edit the generated .prd.config.json file to define your production configuration.

2.  Edit the generated .dev.config.json file to define your development / testing configuration.  

3.  When using SSL to test locally, SSL certs will be needed.  See the SSL setup section below for 
    instructions regarding the generation of certificates suitable for *local testing* via go test.
___
<br/>

## Execution
The generated server runs based on a generated JSON configuration file as shown below.  

```code

{
    "port": 3000,
    'port' is used to instruct the generated server which tcp port to publish the service end-points on.

    "env": "def",     
    'env' is used to inform the generated server which mode to run in.  The material difference
    between "dev", "def" and "prod" is slight; the "dev" and "def" modes run the ORM in debugging
    mode, thereby causing the generated SQL statements to be written as a log to stdout.

    "pepper": "secret-pepper-key",
    'pepper' is used as a pepper seed to the bcrypt password hash.  The generated server handles
    user login authentication via bcrypt hashing of the password the user entered, then comparing
    the resulting hash to the stored bcrypt password hash that was created when the user set their
    initial password.  Passwords are not kept anywhere in the system.

    "hmac_Key": "secret-hmac-key",
    'hmac_key' is a legacy configuration option left over from the old CSRF implementation.  
    This field has been deprecated.

    "database": {
        "db_dialect": "postgres".
        "host":       "localhost",
		"port":       5432,
		"user":       "godev",
		"password":   "gogogo123",
		"name":       "glrestgen"
    },
    'db_dialect' refers to the backend database type that will be used by the generated application.  Currently,
    the following db_dialects are supported:
    
    |  Database               | JSON Value for db_dialect field    |
    |------------------------:|:----------------------------------:|
    | Postgres                | "db_dialect": "postgres"           |   
    | MSSQL (2008+)           | "db_dialect": "mssql"              |
    | SAP Hana                | "db_dialect": "hdb"                |
    | SQLite3                 | "db_dialect": "sqlite3"            |
    | MySQL / MariaDB         | "db_dialect": "mysql"              |
    |                         |                                    |
    |                         |                                    |  

    'database' is a JSON block holding the access information for the database.  At the moment, 
    only Postgres is implemented and it is assumed that the application will use the pg default
    scehema 'Public'.  Sqlite3 will be implemented as a proof-of-concept for multi-db-support,
    and the target database schema will be added to the JSON block at that time.

    "cert_file": "",
    'cert_file' should point to the location of a self-signed or purchased certificate file and
    is used to support https.  Maintaining a 'cert_file' and 'key_file' in the configuration 
    informs the generated server to publish via https.

    "key_file": "",
    'key_file' should point to the location of the key-file for the self-signed or purchased 
    certificate file referenced in the 'key_file' configuration key.  Maintaining a 'cert_file'
    and 'key_file' in the configuration informs the generated server to publish via https.

    "jwt_priv_key_file": "jwtkeys/private.pem",
    "jwt_pub_key_file": "jwtkeys/public.pem"
    Application access is handled via claims embedded in JWT tokens.  JWT content is encrypted
    via ECDSDA-384, thereby requiring a set of valid key-files to support the initial encryption 
    for the login reponse, as well as the subsequent decryption of the 'Authorization' http 
    header field for each incoming request.  The JWT key-files are automatically generated when 
    the server codebase is created.  Leaving the default values in these fields is recommended; 
    they have been included in the configuration in order to support their storage in an alternate 
    location.
}

```
<br/>

### Default Config
    $ go run main.go 

    This will run the program using a set of default configuration that has been compiled into the binary.  
    Default configuration may be edited in the generated appobj/appconf.go file to suit local requirements.  
    The default application settings are shown in the server configuration file format.  
    The default configuration publishes the end-points on port 3000 over http due to the absence of the 
    'cert_file' and 'key_file' values.

```JSON

{
    "port": 3000,    
    "env": "def",     
    "pepper": "secret-pepper-key",  
    "hmac_Key": "secret-hmac-key",
    "database": {
        "db_dialect": "postgres",
        "host":       "localhost",
		"port":       5432,
		"user":       "godev",
		"password":   "gogogo123",
		"name":       "glrestgen"
    },
    "cert_file": "",
    "key_file": "",
    "jwt_priv_key_file": "jwtkeys/private.pem",
    "jwt_pub_key_file": "jwtkeys/public.pem"
}

```
<br/>

### Development Config
    $ go run main.go -dev

    The program will be executed using the configuration specified in the content of .dev.config.json.  
    The generated sample dev configuration file should be edited to match the local environment.

```JSON

{
    "port": 3000,    
    "env": "def",     
    "pepper": "secret-pepper-key",  
    "hmac_Key": "secret-hmac-key",
    "database": {
        "db_dialect": "postgres",
        "host":       "localhost",
		"port":       5432,
		"user":       "godev",
		"password":   "gogogo123",
		"name":       "glrestgen"
    },
    "cert_file": "",
    "key_file": "",
    "jwt_priv_key_file": "jwtkeys/private.pem",
    "jwt_pub_key_file": "jwtkeys/public.pem"
}

```
<br/>

### Production Config
    $ go run main.go -prod

    The program will be executed using the configuration specified in the content of .prd.config.json.
    The generated sample prd configuration file should be edited to match the local environment.  A 
    sample .prd.config.json file is shown below.  This file will instruct the server to publish the 
    end-points on port 8080 using https.

```JSON

{
    "port": 8080,    
    "env": "prod",     
    "pepper": "secret-pepper-key",  
    "hmac_Key": "secret-hmac-key",
    "database": {
        "db_dialect": "postgres",
        "host":       "localhost",
		"port":       5432,
		"user":       "godev",
		"password":   "gogogo123",
		"name":       "glrestgen"
    },
    "cert_file": "srvcert.cer",
    "key_file": "srvcert.key",
    "jwt_priv_key_file": "jwtkeys/private.pem",
    "jwt_pub_key_file": "jwtkeys/public.pem"
}

```
___    

<br/>

## Generate Self-Signed Certs for https Testing
If you wish to perform local https-based testing, it is possible to do so through the use of self-signed
certificates.  Self-signed certificates can be easily created through the use of the openssl tool on 
*nix systems.  
<br/>

### Verify the OpenSSL Installation

Open a terminal session and verify that openssl is available:
```code
$ which -a openssl
/usr/bin/openssl
$

```
If openssl is not shown in the 'which' command output, check your path to ensure you have access to /usr/bin 
or /usr/local/bin.  If you have access to the ./bin directories, but still cannot find the openssl tool, 
it can be downloaded from https://www.openssl.org/source/ .  Follow the directions on the site to correctly 
download and install the tool.
<br/>

### Generate a Private Certificate Authority (CA) Certificate

Open a terminal session and execute the openssl command as shown:
```code
$ openssl genrsa -out "myCA.key" "2048"
Generating RSA private key, 2048 bit long modulus
...................................+++
..........................................................................................+++
e is 65537 (0x10001)
$

```
Verify that a file called "myCA.key" has been created.
<br/>

### Generate a Private Certificate Authority (CA) Certificate

Open a terminal session and execute the openssl command as shown:
```code
$ openssl req -x509 -new -days 365 -key "myCA.key" -out "myCA.cer" -subj "/CN=\""MyCompanyName"\""

```
There is no ouput to this command, so verify that a file called "myCA.cer" has been created.
<br/>

### Generate a Private Server Key

Open a terminal session and execute the openssl command as shown:
```code
$ openssl genrsa -out "srvcert.key" "2048"
Generating RSA private key, 2048 bit long modulus
..............................................................................................+++
.....+++
e is 65537 (0x10001)
$

```
Verify that a file called "srvcert.key" has been created.
<br/>

### Create a Private Server Certificate Signing Request

This generates an intermediate certificate signing request file (.csr) based on the Private Server 
Key created in the previous step.  The creation of the CSR is an interogative process, but for 
self-signed testing, most of the inputs can safely be ignored.  Follow the prompts as per the 
example shown below:
```code
$ openssl req -new -key srvcert.key -out srvcert.csr
You are about to be asked to enter information that will be incorporated
into your certificate request.
What you are about to enter is what is called a Distinguished Name or a DN.
There are quite a few fields but you can leave some blank
For some fields there will be a default value,
If you enter '.', the field will be left blank.
-----
Country Name (2 letter code) [AU]:CA
State or Province Name (full name) [Some-State]:AB
Locality Name (eg, city) []:
Organization Name (eg, company) [Internet Widgits Pty Ltd]:MyCompany
Organizational Unit Name (eg, section) []:
Common Name (e.g. server FQDN or YOUR name) []:
Email Address []:

Please enter the following 'extra' attributes
to be sent with your certificate request
A challenge password []:
An optional company name []:
$

```
Verify that a file called "srvcert.crt" has been created.
<br/>

### Create a Private Server Certificate

This is the final step in getting the required certicate and key files to support local https
testing.  In this step, the CA certificate and private key files will be used in conjunction 
with the private server key and private server signing-request to generate a private server 
certificate.
```code
$ openssl x509 -req -in srvcert.csr -out srvcert.cer -CAkey myCA.key -CA myCA.cer -days 365 -CAcreateserial -CAserial 123456
Signature ok
subject=/C=CA/ST=AB/O=MyCompany
Getting CA Private Key

```
Verify that a file called "srvcert.cer" has been created.
<br/>

### Ensure myCA.cer is Trusted Locally
Ensure that myCA.cer is fully-trusted in your local certificate store.  The process to do this will differ per operating system, so look online for instructions regarding 'trusting a self-signed CA certificate'.

### Add Certificates to the Configuration File
In order to publish the generated services over https, add the "srvcert.cer" and "svrcert.key" files to the 'cert_file' and 'key_file' keys respectively in the appropriate configuration file.  Additionally, the myCA.key file must be placed in the same directory as the "srvcert.*" files in order for go's https (TLS) server to operate correctly.

___
<br/>

## Automated Testing
Automated testing can be performed using the standard go test tooling.  Tests can be run using http
or https, and run against any port that the application is presently serving on.  Remember, the 
application must be running prior to executing the test.

The generated CRUD tests check the availability of the end-points, and attempt to perform CRUD
activities using representitive data for the field-types.  If customization has occurred in the
model normalization and validation enhancement points, the field values used in the generated
main_test.go file should be updated accordingly.  The generated CRUD tests are provided as a 
starting point for your own testing. 

The generated simple selector tests check the availabilty of the end-points, and attempt to peform
a GET for each of the selection operators specified in the Enitity->selectable field in the models.json
file.  It is not neccessary to have values populated in the dataabase in order for the simple selector
 tests to run.


### Run go test With https
```code
    $ go test -v -https -port "8080"

```

### Run go test Without https
```code
    $go test -v -port "8080"

```

___

<br/>

## Discrete CRUD Testing
TODO: add steps describing how to test the generated server from a tool like Postman


![alt text](/md_images/postman_login.png "Postman Login")

___

## Pending Changes
  - [ ] add service activation to the config
  - [x] add support for additional db platforms via the dialect
    - [ ] write a dialect for db2 community edition
    - [x] write a dialect for hana as a relational-db
    - [ ] hana hybrid model(...)
  - [ ] use of claims as scopes in the middleware to dicate access to routes / actions
  - [ ] add server-side user creation / disallow open user creation route
    - [ ] web-based interface for API documentation?
  - [ ] add method-chaining to new 'cust' package to allow for code-regen
  - [ ] add opportunistic locking via etag concept
    - [ ] look at fast hash algorithms (murmur-2??)
  - [x] add Href to entities as a common self-referential field
  - [ ] add code to support the links via child-href
  - [ ] add code to support expansion of child-href
    - [ ] 	Href string  `rgen:"-" json:"Href,omitempty"`
	- [ ]   Test string  `rgen:"-" json:"Test,omitempty"`
  - [ ] add code to support filtering of expansions
  - [ ] add scopes to config
    - [ ] use scopes in JWT to allow / disallow access
  - [x] update main_test.go to create and delete the test user
  - [x] replace custom model interpretation code with https://golang.org/pkg/encoding/json/#Unmarshal
  - [ ] enhance model
    - [x] support single-field index creation via model attribute
    - [x] support not-nullable directive via model attribute
    - [x] support native dbType column directive via model attribute
    - [x] support selectable directive via model attribute
      - [x] create a test handler for User{}ByID in the router
      - [x] create template for single-field lookup based on User{}ByID()
      - [x] call template following the CRUD method creations (controller.gotmpl & model.gotmpl)
      - [x] add handlers following the CRUD handler processing (appobj.gotmpl)
    - [x] support compound index directive via model attribute
  - [x] disallow snake case in the ddlconfig element names
  - [x] convert camelCase model field name to snake_case using the gorm conversion routine
  - [x] add a flag for model file i.e.   $ go run main.go -m "/Users/tomthedog/config/mymodel.json
  - [x] look at how gorilla.mux handles routes like  ../product?Attr1='foo'&&Attr2
    - *see https://stackoverflow.com/questions/45378566/gorilla-mux-optional-query-values*
  - [x] add support for a dev config.json file
  - [ ] add support for LetsEncrypt
  - [x] add capability of generating keys for JWT via ecdsa256
  - [x] add automated default tests 
  - [x] run go fmt on each file immediately following generation?
  - [x] extend model to support JSON-type tags for GORM etc.
  - [x] remove the gorilla csrf dependency; the use of JWT's in a stateless application obviates the need for CSRF protection. 
  - [x] run goimports on generated code  
  - [ ] add the capability of automatically running go get (look at go dep) for missing packages in the dependency list
  - [ ] add capability to generate self-signed certs for local ssl testing
  - [ ] create github repo for gnerated code via https://godoc.org/github.com/google/go-github/github#RepositoriesService
