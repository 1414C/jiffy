# rgen

## Overview and Features
Rgen is a model-based application services generator written in go.  It was developed as an experiment to offer an alternative avenue when developing cloud native applications for SAP Hana.  The rgen application allows a developer to treat the data persistence layer as an abstraction, thereby removing the need to make use of CDS and the SAP XS libraries.  While this is not for everybody, it does reduce the mental cost of entry and allows deployment of a web-based application to SAP Hana with virtually no prior Hana knowledge.

#### Why write in Go?
* Go has a very strong standard library, thereby keeping dependencies on public packages to a minimum
* Go offers true concurrency via lightweight threads known as goroutines 
    - no blocking in the i/o layer during compute intensive tasks
    - no 'lost' callbacks or 'broken' promises
    - goroutines will use all available cores to handle incoming requests
* Go it is a small language that offers type-safety
* Go projects compile to a static single binary which simplifies deployments
* Go cross-compiles to virtually any platform and architecture; write and test on a chromebook - deploy to z/OS
* Go is making inroads into areas that have been dominated by other languages and packages

#### What does the Rgen application provide?
* generated apps can be connected to Postgres, MSSQL, SAP Hana, SQLite or MariaDB
* no database specific code is compiled into the binary; an app can be pointed from SQLite to SAP Hana with no code changes
* login / session management via jwt
* built-in support for the creation of signing-keys for jwt
* bcrypt salt/pepper based authentication scheme where passwords are never stored in the db
* JSON configuration (model) file for Entity, Index and Relationship definitions
* models support persistent and non-persistent fields
* automatically creates backend database artifacts based on the model file (tables, indices)
* supports single and composite index declarations via the model file
* built-in support for https
* baked in normalization and validation in the model-layer
* generates a working set of CRUD-type RESTful services for each entity in the model file
* supports and generates working end-points for hasOne, hasMany and belongsTo entity relationships 
* generates working query end-points based on the model fie 
* end-points are secured by way of scope inspection (jwt claims) in the route handler middleware
* generates a comprehensive set of working tests (go test)
* generated code is easily extended
<br/>

#### What does an application look like?
The generated application can be pointed at the DBMS of your choice without the need to recompile the binary (architecture differences not withstanding).  This means that a developer can build a model, fully test it locally using SQLite and then redirect the appplication to a formal testing environment running SAP Hana, or any of the other supported database systems.  This is achievable due to the ORM layer that the Rgen application is built upon.  The ORM is easily extendable to accomodate other databases if required (oracle, db2, SAP ASE are candidates here).

Applications are generated based on model files which are encoded as simple JSON.  The concepts of entity and resource-id form the cornerstones upon which the model, application and RESTful end-points are built upon.

Entities can be thought of anything that needs to be modelled; Order, Customer, Invoice, Truck, Oven, ..., ... Each entity is mandated to have an ID field, which is analagous to a primary-key or row-id in the backend database.  ID is used as the primary resource identifier for an entity, and is generally be setup as an auto-incrementing column in the database.  

Accessing an entity via the generated CRUD interface is very simple.  For example, a customer could be defined in the model and then accessed via the application as follows:

1.  Create a customer entity:
    - https://servername:port/customer  + {JSON body}

2.  Update a customer entity:
    - https://servername:port/customer/:id  + {JSON body}

3.  Read a customer entity:
    - https://servername:port/customer/:id

4.  Delete a customer entity:
    - https://servername:port/customer/:id

5.  Read all customer entities:
    - https://servername:port/customers


Additional routes can also be generated based on the model file, including custom filters for GET operations, static end-points for common GET operations, HasOne, HasMany and BelongsTo relationships:

1.  Use a filter to Get customers where the last name is 'Smith':
    - https://servername:port/customers/?last_name=Smith

2.  Use a generated static end-point to Get customers where credit score is less than 4:
    - https://servername:port/customers/credit_score(LT 4)

3.  Use a generated relationship to retrieve all orders for a customer:
    - https://servername:port/customer/10023/orders

4.  Use a generated relationship to retrieve a specific order for a customer:
    - https://servername:port/customer/10023/order/99000022

5.  Use a generated relationship to retrieve the customer for a specific order:
    - https://servername:port/order/99000022/customer


This is just a sample of what the model files have to offer.  More details regarding application modlelling are contained in later sections of this file.

#### Access Control
Access to resources (entities) is controlled in three ways:

1. Configuration based service activation
2. Secure user authentication
3. JWT tokens with claim inspection in middleware applied to the protected routes (end-points) 

An internal service is created for each of the modelled entities in the application.  Services can be marked as active or inactive in the service configuration, thereby allowing a single application to be generated, but also allowing selective service deployment.  For example, there may be cases where it is desirable to route certain services to a particular server and another set of services to the rest of the pool.  In such a case, NGix could be configured to route the end-points appropriately, and the deployed service configurations would be adjusted accordingly.

User authentication is conducted using bcrypt in such a manner that passwords are never stored in the application database.  When a user is created, their user-id is stored in the backend database along with the salt/peppered bcrypt hash of their password.  This ensures that in the event of a breach no plain-text passwords can be obtained.  

The bcrypt hashes are not very useful to would-be attackers for the following reasons:
* bcrypt hashes are salt/peppered making rainbow tables useless
* bcrypt is slow by design, making brute force reversal a time-consuming and expensive proposition
* as increased computing power becomes available, the bcrypt cost parameter can be increased (current = 14)
* the hash itself is not used for authentication; it is the by-product of successful authentication

When a user logs into the application the following steps occur:
* the user-name and stored bcrypt hash is looked up in the back-end db
* the provided password is hashed in memory using the standard lib Go bcrypt functions and the protected salt/pepper values
* the computed bcrypt hash is compared to the stored hash for the user
* if the hash values match, a JWT (token) is created using ECDSA-384
* the JWT is passed back to the caller and must henceforth be included in the http header of all requests in the Authorization field
* in addtion to fullfilling the authorization requirements, the JWT is also used as a CSRF equivalent
* see the Authorization section for more details regarding the content and use of the JWT content/claims

<br/>

## Work-In-Progress
1.  [ ]Ensure that rune and byte types are fully accounted for
2.  [ ]Extend claims support in the route middleware
3.  [ ]Add support for 'sqac:"default:xxxyyyzzz"' directives
    - default value
    - default (sqac) function (datetime defaults for example)
4.  [ ]Add option for Foreign Key defintion / enforcement in relations
5.  [ ]Droplet deployment
6.  [ ]NGinx
7. [-]Complete service activations
 - config model files and file generation updated
<br/>

## Installation and Execution
In order to run the application generator, ensure the following:

1.  Make sure go has been installed in the test environment.  See http://www.golang.org for installation files and instructions.

2.  From $GOPATH/src, use go get to install the following:
    * go get -u github.com/gorilla/mux
    * go get -u github.com/dgrijalva/jwt-go
    * go get -u github.com/golang.org/x/crypto/bcrypt
    * go get -u github.com/1414C/sqac
    * go get -u github.com/lib/pq
    * go get -u github.com/SAP/go-hdb/driver
    * go get -u github.com/go-sql-driver/mysql
    * go get -u github.com/mattn/go-sqlite3
    * go get -u github.com/MSSQL

    ** godep will be incorporated in order to eliminate this installation step

3.  Install the application into your local $GOPATH/src directory:
    * go get -u github.com/1414C/rgen

4.  You will need access to a Postgres, MySQL, MSSQL or SAP Hana database, either locally or over the network.  It is also possible to run tests with SQLite3.
    
5.  The application can be started in two ways:
    * From $GOPATH/src/github.com/1414C/rgen you may execute the application by typing:
        * go run main.go     
    * A binary can also be built from $GOPATH/src/github.com/1414C/rgen by typing the following:
        * go build .
        * The application can then be started from the same directory by typing:
            * ./rgen
<br/>

## Flags
Flags are generally not used, as the configuration files (models.json) are easier to deal with.  There are however, a few flags that can e appended to the execution command:

* go run *.go -p <target_dir>
	* The -p switch is used to specify the target directory for generated application source-code relative to $GOPATH/src.

```bash

    $ go run main.go -p "github.com/footle.com/myrestfulsvc"

```

* go run main.go -m <model_file>.json
    * By default, the application will attempt to use ./models.json as the model source, but inclusion of the -m flag permits the use of an alternate model file.
    * The path of model file in the application base directory must be prefaced with ./ .  If the model file is not located in the base directory of the application, the full path must be specified when using the -m flag.

```bash

    go run main.go -m "./my_model.json"

```

---
<br/>

## Model Creation
Create a model file containing the Entities, Indexes and Relations that you wish to generate services for.  Entity model defintion consists of an array of JSON objects, with each object being limited to a flat hierarchy and basic go-data-types, although this is easily extended.  By default, the generator expects a *models.json* file in the execution directory, but a correctly formatted JSON file can be loaded from any location by executing with the *-m* flag. 

Sample <models>.json files are installed with the application and can be found in the testing_models folder.  The sample models are used as the basis for the following sections.

### Simple Single Entity Model

The following JSON illustrates the defintion of a simple single-entity model file.  In this case, a model entity called 'Person' will be created in the generated application, along with corresponding database table 'person'.  Table 'person' will be created (if it does not already exist) when the application is started for the first time.  See the application startup sequence section of this document for details regarding database artifact creation and updates.

```JSON

{
    "entities":  [
        {
            "typeName": "Person",
            "properties": {
                "name": {
                    "type": "string",
                    "db_type": "",
                    "no_db": false,
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
                }
            }
        }
    ]
}

```
Sample model ![simpleSingleEntityModel.json](/testing_models/simpleSingleEntityModel.json "simpleSingleEntityModel.json")

The simpleSingleEntityModel.json file structure and content is explained below:

```code
{
    "entities": [
    The 'entities' block contains an array of entities belonging to the application model.  Each entity relates
    directly to a database table (or view).  Entities contain information that the application generator uses 
    to create and update database artifacts such as tables, indexes, sequences and foreign-keys, as well as
    information informing the application runtime of the member field properties.
    This is a mandatory model element.
    {
        "typeName": "Person"
        Field 'typeName' refers to the name of an entity.  It should be capitalized and written in CamelCase.
        An Entity given a typeName of "Person" will result in an internal model object of type Person and a 
        database table called 'person'.
        This is a mandatory model element.

        "properties": {
        The 'properties' block contains 1:n entity member field definitions.  Member fields should be defined
        in camelCase and can start with a lower or upper-case character.  In the context of the "entity" with
        a 'typeName' of 'Person', 'properties' refer to the data fields of the generated "Person" model 
        structure.  'properties' are a collection of free-form name-tags, each with a child-block containing 
        the 'property' attributes.
        This is a mandatory model element.
        
        The "name" property block is described below:

            "name": {
                "type": "string",
                Field 'type' in a 'properties'->'name' JSON-block refers to the go data-type associated
                with the current 'property'.
                'type' is a mandatory field in an "entity" 'property' block.

                "dbtype": "varchar(100)",
                Field 'dbtype' can be used to specify a native db-field-type for the property.  This feature
                can be useful if for example, the developer is confident that a string will never exceed 100
                characters in length.  Care should be taken to ensure that the specified DB-Type is consistent
                with the go-type that will be generated in the model.<Entity> defintion in the application.
                Consider also that making use of this field to some extent limits the backend portability of 
                the generated code.  For example, not all database systems have a TINYINT data-type, so
                specifying a 'db_type' of TINYINT could be problematic if multiple database systems are 
                being used for testing.
                This is an optional field. 

                "no_db":
                Field 'no_db' can be used to instruct the generator to create the field as a member in the 
                enitity struture, but to prevent the field from being persisted to the backend database.
                Data like passwords for example should never be persisted to the database, but it handy to
                have in the user entitiy definition to help with the login process.  Non-persisted fields
                are not created in the database table schemas, and values passed into the application 
                are wiped from their respective internal structures following use.

                "format": "", 
                Field 'format' is not currently used, but is intended to deal with field conversion from
                strings / floats to timestamp formats etc.
                This is an optional field.

                "required": false,
                Field 'required' is used to instruct the generator to set a 'NOT NULL' database constraint
                on the column related to the property.  
                Allowed values include {true, false}.
                This is a mandatory field.

                "unique": false,
                Field 'unique' is used to instruct the database not to accept duplicate values in the
                database column related to the property.  Setting this field to true will cause a 'UNIQUE' 
                constraint to be applied to the related database column.
                Allowed values include {true, false}.
                This is a mandatory field.

                "index": "nonUnique",
                Field 'index' is used to instruct the database to create an index on the db table-column
                related to the property. 
                See the 'indexes' element in the type definition for the creation of compound indices.
                Allowed values include {"unique", "nonUnique", ""}.
                This is an optional field.

                "selectable": "eq,like"
                Field 'selectable' can be used to instruct the code-generator to create simple REST query 
                accessor routes for the current 'property'.  The generator creates routes to permit GET
                operations that can be called based on the entity 'typeName' and 'property' values.
                Allowed values include {"EQ", "eq", "LT", "lt", "GT", "gt", "GE", "ge", "LIKE", "like", "NE", "ne", ""}
                Additional restrictions are imposed based on the 'type' field value.  For example, a bool
                type need not support LT or GT operators.
                Sample routes for Person->Name selection with "eq,like" are shown:
                    
                    https://localhost:<port>/persons/name(EQ '<sel_string>')
                    https://localhost:<port>/persons/name(LIKE '<sel_string>')

                Note that this is not the same thing as filtering insofar as setting the selectable options 
                results in the creation of parameterized static routes in the application mux.

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
    }
    },    
    {
        ... next entity definition
    }
    ]
}
```

### Entity ID
The ID field is visibly absent from the preceding entity declarations.  The original intent was to support any name for the primary key / resource identifier of an entity.  While it is possile to do this, it seems that ID is the universal 'non-standard' way of representing object identifiers in RESTful-type services, so we went with it.  As a result, ID is injected into the model defintion of every entity as a uint64 field and is marked as the primary-key in the database backend.  By default, the ID is created as an auto-incrementing column in the DBMS, but this functionality can be suppressed (future).  The ability to allow a specific starting point for the ID key range exists in the ORM, and is in the process of being added to the model file. 

If the ID field really needs to be known as CustomerNumber for example, the generated code can be edited in a few locations to support the change.  It is worth mentioning that the number of edits required to rename 'ID' increases in direct relation to the number and complexity of entity relations (both to and from).

As an alternative to renaming ID, it is also conceivable that it can be ignored.  Ignoring the ID means that the generated CRUD controller/model/routes are not as useful as they could be, but they offer a great starting point for your own coding.  Entities can be defined with column constraints that mimic those of DBMS primary / complex keys, then the generated CRUD artifacts based on ID can be ignored, copied then ignored, or modified to accmodate the modelled entities.

It is also possible to go completely custom and write your own models and controllers from scratch using a generated model as a reference template.  In addition to exposing a generic internal CRUD interface to the backend, the more interesting go/sql calls are exposed internally along with some lightly wrapped and super useful calls from jmoirons widely used sqlx package.

https://github.com/jmoiron/sqlx

http://jmoiron.github.io/sqlx/

Although rgen eschews non-standard lib packages wherever possible, sqlx is worth making an exception for.


### Simple Two Entity Model

The following JSON illustrates the defintion of a simple two-entity model file.  In this case, model entities 'Person' and 'Country' will be created in the generated application, along with corresponding database tables 'person' and 'country'.  No relationships have been defined between the two entities; this example simply illustrates how to add multiple entity definitions to the model file.

models.json
```JSON

{
    "entities":  [
        {
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
                }
            }
        },
        {
            "typeName": "Country",
            "properties": {
                "name": {
                    "type": "string",
                    "format": "", 
                    "required": false,
                    "unique": false,
                    "index": "nonUnique",
                    "selectable": "eq,like"
                },
                "isoCode": {
                    "type": "uint",
                    "format": "", 
                    "required": false,
                    "unique": false,
                    "index": "unique",
                    "selectable": "eq,lt,gt"
                }
            }
        }
    ]
}

```
Sample model ![simpleTwoEntityModel.json](/testing_models/simpleTwoEntityModel.json "simpleTwoEntityModel.json")


### Two Entity Model With Composite Index

The following JSON illustrates the addition of a composite-index to an entity definition.  An index composed of the 'name' and 'province' fields has been declared in the 'Owner' entity.  This declaration will result in the creation of a non-unique b-tree index for columns 'name' and 'province' in the database.  Any number of composite indices may be declared for an entity.  No relationships have been defined between the two entities; this example simply illustrates how to declare a composite-index for an entity.  

```JSON

{
    "entities":  [
        {
            "typeName": "Owner",
            "properties": {
                "Name": {
                    "type": "string",
                    "format": "", 
                    "required": false,
                    "unique": false,
                    "index": "nonUnique",
                    "selectable": "eq,like"
                },
                "LicenseNumber": {
                    "type": "uint",
                    "format": "", 
                    "required": false,
                    "unique": true,
                    "index": "",
                    "selectable": "eq,lt,gt"
                },
                "Province": {
                    "type": "string",
                    "format": "", 
                    "required": false,
                    "unique": false,
                    "index": "",
                    "selectable": "eq,lt,gt"
                }
            },
            "compositeIndexes": [ 
                {"index": "name, province"}
            ]
    },
    {
        "typeName": "Car",
        "properties": {
            "Model": {
                "type": "string",
                "format": "", 
                "required": true,
                "index": "nonUnique",
                "selectable": "eq,like"
            },
            "Make": {
                "type": "string",
                "format": "", 
                "required": true,
                "index": "nonUnique",
                "selectable": "eq,like"
            }
        }
    }
    ]
}

```
Sample model ![twoEntityWithCompositeIndex.json](/testing_models/twoEntityWithCompositeIndex.json "twoEntityWithCompositeIndex.json")


## Entity Relationships

Relationships between entities can be declared in the application model file  via the addition of a 'relations' block inside an entity's declaration.  Relationships are based on resource id's by default, although it is possible to specify non-default key fields in the configuration, or implement complex joins directly by maintaining the entity's controller and model.  'relations' blocks look as follows:

```JSON

            "relations": [
                    { 
                    "relName": "ToOwner",
                        "properties": {
                            "refKey": "",
                            "relType": "hasOne",
                            "toEntity": "Owner",
                            "foreignPK": "" 
                        }
                    }
                ]

```
The sample relations block illustrates the declaration of a HasOne relationship between Car and Owner making use of default-keys.


### HasOne Relationship

HasOne relationships establish a one-to-one relationship between two model entities.  As an example, let's posit that a Car can have one Owner.  If the Car and Owner were modelled as entities, we could declare a HasOne (1:1) relationship between them.  The relation would be added in the 'relations' block inside the Car entity definition (as shown above).

A break-down of the relations block fields is as follows:
```code
{
    "relations": [
    The 'entities' block contains an array of relations belonging to the containing entity definition.
    Each relation is defined from the perspective of the containing entity having a relationship of 
    the specified type (in this case hasOne), with the entity referenced in the declaration.  A Car
    has one Owner - in our example at least.
    {
        "relName": "Owner"
        Field 'relName' refers to the name the relationship will be known by inside the application 
        and in the end-point definition in the mux routes.  It must be capitalized and written in 
        CamelCase.  Any name may be chosen for this field, but keep in mind this name will be exposed
        to the service consumer via the URI, so something respecting the relationship enities and 
        cardinaliy is best.  For the example, we have chosen a relName of 'ToOwner' to demonstrate 
        the difference between the toEntity and relName fields.
        relName is a mandatory field in a relations declaration.
            
            "properties": {
            The 'properties' block contains the details of the relationship.

            "refKey":
            Field 'refKey' can be used to specify an non-default reference key belonging to the 
            containing (from) entity.  By leaving this field empty, the default field of 'ID' will 
            be used, which is what most relationships will use most of the time.  For those times 
            where the default 'from' key cannot be 'ID', you may specify your own as long as the 
            chosen field is an existing member in the containing (from) entity and is of go-type 
            uint64 or *uint64.  The refKey will be matched in the selection of the toEntity when 
            the relationship is accessed.
            This is an optional field.

            "relType":
            Field 'relType' is used to indicate what sort of relationship is being declared between
            the containing (from) entity and the toEntity.  Valid values are {HasOne, HasMany and 
            BelongsTo}.
            This is a mandatory field.

            "toEntity":
            Field 'toEntity' is used to specify the target entity in the relationship. The toEnity 
            must be capitalized and provided in CamelCase that matches that used in the toEntity's 
            declaration.  The toEntity need not appear prior to the containing entity in the model 
            file or files.
            This is a mandatory field.

            "foreignPK":
            Field 'foreignPK' can be used to specify the field in the toEntity to which the containing
            entity will match the 'refKey'.  As such, both fields must be of the same go-type 
            (uint64/*uint64).  By leaving this field empty, the application will attempt to use 
            <ContainingEntityName>ID as the column to which the containing (from) entity will attempt
            to match its refKey to.  In the given example of Car -> Owner, the application will attempt
            to find the Car's Owner as shown in the following pseudo-code:

            SELECT * FROM owner WHERE owner.CarID = car.ID LIMIT 1;

            }
    }
    ]      
}
```

### HasMany Relationship

HasMany relationships establish a one-to-many relationship between two model entities.  As an example, let's posit that a Libary can have many Books.  If Library and Book were modelled as entities, we could declare a HasMany (1:N) relationship between them.  The relation would be added in the 'relations' block inside the Library entity definition: 

```JSON

{
    "entities":  [
        {
            "typeName": "Library",
            "properties": {
                "Name": {
                    "type": "string",
                    "format": "", 
                    "required": true,
                    "unique": false,
                    "index": "nonUnique",
                    "selectable": "eq,like"
                },
                "City": {
                    "type": "string",
                    "format": "", 
                    "required": true,
                    "unique": false,
                    "index": "",
                    "selectable": "eq,lt,gt"
                }
            },
            "compositeIndexes": [ 
                {"index": "name, city"}
            ],
            "relations": [
                    { 
                    "relName": "ToBooks",
                        "properties": {
                            "refKey": "",
                            "relType": "hasMany",
                            "toEntity": "Book",
                            "foreignPK": ""
                        }
                    }
                ]
    },
    {
        "typeName": "Book",
        "properties": {
            "Title": {
                "type": "string",
                "format": "", 
                "required": true,
                "index": "nonUnique",
                "selectable": "eq,like"
            },
            "Hardcover": {
                "type": "bool",
                "format": "", 
                "required": true,
                "index": "",
                "selectable": "eq,ne"
            },
            "Copies": {
                "type": "uint64",
                "format": "", 
                "required": true,
                "index": "",
                "selectable": "eq,lt,gt"
            },
            "LibraryID": {
                "type": "uint64",
                "format": "", 
                "required": true,
                "index": "nonUnique",
                "selectable": "eq,like"
            }
        }
    }
    ]
}

```
Sample model ![hasManyDefaultKeys.json](/testing_models/hasManyDefaultKeys.json "hasManyDefaultKeys.json")


### BelongsTo Relationship

BelongsTo relationships are used to form the inverse of the HasOne and HasMany relations.  Consider the Library HasMany Books example; A Library has many books, but we can also posit that 'a Book belongs to a Library'; this is an example of a BelongsTo relationship.  The JSON below extends the Library -> Book example by adding the BelongsTo relationship to the Book entity definition:

```JSON

{
    "entities":  [
        {
            "typeName": "Library",
            "properties": {
                "Name": {
                    "type": "string",
                    "format": "", 
                    "required": true,
                    "unique": false,
                    "index": "nonUnique",
                    "selectable": "eq,like"
                },
                "City": {
                    "type": "string",
                    "format": "", 
                    "required": true,
                    "unique": false,
                    "index": "",
                    "selectable": "eq,lt,gt"
                }
            },
            "compositeIndexes": [ 
                {"index": "name, city"}
            ],
            "relations": [
                    { 
                    "relName": "ToBooks",
                        "properties": {
                            "refKey": "",
                            "relType": "hasMany",
                            "toEntity": "Book",
                            "foreignPK": ""
                        }
                    }
                ]
    },
    {
        "typeName": "Book",
        "properties": {
            "Title": {
                "type": "string",
                "format": "", 
                "required": true,
                "index": "nonUnique",
                "selectable": "eq,like"
            },
            "Hardcover": {
                "type": "bool",
                "format": "", 
                "required": true,
                "index": "",
                "selectable": "eq,ne"
            },
            "Copies": {
                "type": "uint64",
                "format": "", 
                "required": true,
                "index": "",
                "selectable": "eq,lt,gt"
            },
            "LibraryID": {
                "type": "uint64",
                "format": "", 
                "required": true,
                "index": "nonUnique",
                "selectable": "eq,like"
            }
        },
        "relations": [
                    { 
                    "relName": "ToLibrary",
                        "properties": {
                            "refKey": "",
                            "relType": "belongsTo",
                            "toEntity": "Library",
                            "foreignPK": ""
                        }
                    }
                ]
    }
    ]
}

```

By relying on the default key determinations for the BelongsTo relationship, the generator determines that the Book.LibraryID field should be matched against field Library.ID.  If alternate keys are desired, they can be specified in the 'refKey' and 'foreignKey' property fields in the BelongsTo relation declaration.


### What if more complex relationships are required?

At the moment the generator only supports HasOne, HasMany and BelongsTo relations, as in practice these tend to be the most widely used.  The generated code can be extended to accomodate additional relationships and joins if need be.  There is a tentative plan to support more complex relations in the generator in the future.


<br/>

# What gets generated?

Running the rgen generator creates a set of files that comprise a basic working application.  Incoming requests are handled by a mux, which validates the request, and then matches it to a route.  The selected route passes the request to a controller specific to the entity-type, where the incoming information is mapped into a go struct matching the entity declaration.  The controller then calls the appropriate model function for the http operation and entity-type combination, passing it the entity structure.  The model handler passes the entity struct through a member-field validation layer, and then to the model's interface to the underlying sqac ORM.  The database request is handled by the ORM, and then the response is passed from the model back to the controller where it is packaged as a JSON payload and sent back to the caller in the response-writer's body.

There are more elegant ways to express certain aspects of the generated application.  The coding style has been deliberately kept as simple and straight-forward as possible in order to facilitate easier understanding and adjustment of the generated code.
<br/>
<br/>
### The application folder

![alt text](/md_images/app_layout/AppLayout1.jpeg "Application file structure")
<br/>
Following the execution of the application generator, a folder containing the generated app's files is created as shown.
<br/>

### The appobj folder 

![alt text](/md_images/app_layout/AppLayout_appobj1.jpeg "Application appobj folder content")
<br/>
The appobj folder contains the generated application's configuration loader and the main application object.
<br/>


#### appobj.go
The entry point for go applications is always the main() function, but we seldom write the so-called 'main' part of the application in this monolithic function.  To that end, an AppObj struct is declared and the main thread of the application runs against it.  The content of main.go simply creates an AppObj struct, parses some flags and calls the AppObj.Run() method.

When the generated application is started, AppObj.Run() is responsible for:
- loading the specified config
- creating the runtime services
- performing auto-migration of database artifacts
- initializing the key for JWT/ECDSA support
- instantiating controllers
- initializting routes 
- starting the mux

The creation of the runtime services bears closer inspection before moving on.  Generated applications contain an internal 'service' for each entity declared in the source model files.  The AppObj is responsible for the instantiation of these services when the application is started via the AppObj.createServices() method.

A Services object containing each of the entity runtime services is created on the one-and-only instance of the AppObj.  A runtime service is first created to support access to the backend DBMS via the sqac ORM, then a service is started for each entity.  Entity services contain a reference to the ORM access handle, as well as an instance of the entity's validator class which is contained in the model-layer.

#### appconf.go
The code in appconf.go contains the functions used to load application configuration files, as well as functions containing so-called 'default' configuration.  It is possible to edit the DefaultConfig() function so that it holds values specific to the local test/development environment.  This prevents the need for maintaining a set of configuration files that the development staff need to keep in sync.
<br/>
<br/>
<br/>
### The controllers folder 

![alt text](/md_images/app_layout/AppLayout_controllers.jpeg "Application controllers folder content")
<br/>
The appobj folder contains the generated application's controllers.  A controller is created for each entity that has been declared in the model files, as well as a static controller that is used to handle the application's users.

Controllers act as a bridge between an entity's routes and its model layer.  Each entity mux route is assigned a method in their respective controller based on the intent of that route.  For example, to create a new new Library entity the following POST could be made:

```code

https://servername:port/library {JSON body} + POST

```  

The route for this call is defined in appobj.go as follows, where 'a' is the one-and-only instance of the AppObj:

```golang

a.router.HandleFunc("/library", requireUserMw.ApplyFn(a.libraryC.Create)).Methods("POST")

```
The '/library'-POST route is assigned a HandleFunc belonging to the instance of the LibraryController that has been created on the appobj.  a.libraryC.Create is called for the 'library' route when the http method equals 'POST'.  The route contains some additional code related to authentication and authorization of the requester but this can be ignored for now.  The handler function for a mux.route must conform to the standard go http.Handler interface:

```golang

    type Handler interface {
        ServeHTTP(ResponseWriter, *Request)
    }

```

This interface facilitates the passing of the incoming request header and body to the controller method, as well as the passing of the formatted response back to the router.  With this out of out the way, let's look at generated Controller method LibraryController.Create:

```golang

    // Create facilitates the creation of a new Library.  This method is bound
    // to the gorilla.mux router in main.go.
    //
    // POST /library
    func (lc *LibraryController) Create(w http.ResponseWriter, r *http.Request) {

	    var l models.Library
	    decoder := json.NewDecoder(r.Body)
	    if err := decoder.Decode(&l); err != nil {
	    	log.Println("Library Create:", err)
	    	respondWithError(w, http.StatusBadRequest, "libraryc: Invalid request payload")
	    	return
	    }
	    defer r.Body.Close()

	    // fill the model
	    library := models.Library{
	    	Name: l.Name,
	    	City: l.City,
	    }

	    // build a base urlString for the JSON Body self-referencing Href tag
	    urlString := buildHrefStringFromCRUDReq(r, true)

	    // call the Create method on the library model
	    err := lc.ls.Create(&library)
	    if err != nil {
	    	log.Println("Library Create:", err)
	    	respondWithError(w, http.StatusBadRequest, err.Error())
	    	return
	    }
	    library.Href = urlString + strconv.FormatUint(uint64(library.ID), 10)
	    respondWithJSON(w, http.StatusCreated, library)
}

```
The complete Library.Create(http.Handler) controller method is shown exactly as it has been generated.


Each section of the method is broken down in the following subsets of commented code:
```golang

        // declare a local variable of struct type models.Library to hold the decoded 
        // JSON body provided in the request.Body.
    	var l models.Library

        // create a new JSON decoder passing in the request.Body
        decoder := json.NewDecoder(r.Body)
        
        // call the Decoder.Decode(interface{}) method passing a reference to the locally
        // declared models.Library struct 'l'.  if the decoder is able to decode the JSON
        // contained in the request.Body, the member fields of 'l' will be populated.  if
        // the decoder fails to parse and map the incoming JSON to the models.Library 
        // struct, it will return an error.  The problem will be logged to stdout (for now)
        // on the server-instance, and a response conforming to the http.Handler interface
        // will be constructed and passed back to the router.  if the JSON was parsed 
        // successfully, a defer call is made to ensure that the request.Body will be 
        // closed upon exit of the method.
	    if err := decoder.Decode(&l); err != nil {
	    	log.Println("Library Create:", err)
	    	respondWithError(w, http.StatusBadRequest, "libraryc: Invalid request payload")
	    	return
        }
        defer r.Body.Close()
        
        // fill the model with the parsed content of the JSON body.  this step looks 
        // redundant, but can be thought of as a way to separate the incoming data 
        // from the response.  going forward from this point, 'l' is ignored and 
        // all data transformation occurs on the 'library' variable.
	    library := models.Library{
	    	Name: l.Name,
	    	City: l.City,
	    }

        // build a base urlString for the JSON Body self-referencing Href tag
        urlString := buildHrefStringFromCRUDReq(r, true)
        
        // call the Create method on the library model.  each controller contains an
        // instance of the Service for it's respective entity.  the Create method on 
        // the service is called, passing a reference to the 'library' data structure.
        // recall that the Service for an entity provides the link to that entity's 
        // model-layer by way of the entity's validator.  lc.ls.Create(&library) will
        // result in a call the model Validator Create() method for the Library 
        // entity, and in-turn, call to the enitity's model.Create() method where 
        // the data will be passed to the ORM-layer.  if the Create() call returns
        // an error, the problem will be logged to stdout (for now) on the server-
        // instance, and a response conforming to the http.Handler interface will be
        // constructed and passed back to the router.
	    err := lc.ls.Create(&library)
	    if err != nil {
		    log.Println("Library Create:", err)
		    respondWithError(w, http.StatusBadRequest, err.Error())
		    return
        }
        
        // if the call to the model-layer was successful, it indicates that a new 
        // Library entity was created in the DBMS.  the 'library' reference passsed
        // to the Create() method(s) in the model-layer will now contiain the new 
        // Library's information.  first, the ID for the new Library will be added
        // to the urlString and assigned to the library struct's Href member field.
        // Href is another injected field in the entity and fullfills the purpose
        // of providing a direct URI for the returned entity.  finally the populated
        // 'library' struct is formatted as a JSON response and passed back to the 
        // router along with an http status-code indicating success. 
        library.Href = urlString + strconv.FormatUint(uint64(library.ID), 10)
        respondWithJSON(w, http.StatusCreated, library)
    }

```

<br/>






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
    'db_dialect' refers to the backend database type that will be used by the generated application.
    Currently, the following db_dialects are supported by the sqac ORM runtime:
    
    |  Database               | JSON Value for db_dialect field    |
    |-------------------------|------------------------------------|
    | Postgres                | "db_dialect": "postgres"           |   
    | MSSQL (2008+)           | "db_dialect": "mssql"              |
    | SAP Hana                | "db_dialect": "hdb"                |
    | SQLite3                 | "db_dialect": "sqlite3"            |
    | MySQL / MariaDB         | "db_dialect": "mysql"              |


    'database' is a JSON block holding the access information for the database system.  Fill in
    what is needed for the type of database you are connecting to.  SQLite for example, does not
    have any user-access control etc.  Sample database configuration blocks have been included in
    file sample_configs.json.

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
    The default configuration may be edited in the generated appobj/appconf.go file to suit local 
    requirements.  The default application settings are shown in the server configuration file format.  
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
    "env": "dev",     
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

### Generate a Private Certificate Authority (CA) Certificate Key

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
Ensure that myCA.cer is fully-trusted in your local certificate store.  The process to do this will differ per operating system, so look online for instructions regarding 'trusting a self-signed CA certificate'.  You may also need to adjust the settings in test tools like Postman in order for them to accept self-signed certs.

### Add Certificates to the Configuration File
In order to publish the generated services over https, add the "srvcert.cer" and "svrcert.key" files to the 'cert_file' and 'key_file' keys respectively in the appropriate configuration file.  Additionally, the myCA.key file must be placed in the same directory as the "srvcert.*" files in order for go's https (TLS) server to operate correctly.

___
<br/>

## Automated Testing
Automated testing can be performed using the standard go test tooling.  Tests can be run using http
or https, and run against the port that the application is presently serving on.  Remember, the 
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

 At the moment, relationships are not included in the generated tests.


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
  - [x] fully implement nullable / pointer support
  - [x] add support for single-field unique constraints
  - [x] implement GetEntities to use the standard sqac.PublicDB interface
  - [ ] SAML integration with usr/login
  - [ ] implement self-documenting API 
  - [ ] consider the use of db-views as entity sources
  - [ ] add support for BLOB storage (S3?)
  - [ ] add service activation to the config
  - [x] add support for additional db platforms via the dialect
    - [ ] write a dialect for db2 community edition
    - [ ] write a dialect for ASE
    - [ ] write ASE driver?
    - [x] write a dialect for hana as a relational-db
    - [ ] hana hybrid model(...)
  - [ ] add server-side user creation / disallow open user creation route
    - [ ] web-based interface for API documentation?
  - [ ] add method-chaining to new 'cust' package to allow for code-regen
  - [ ] add opportunistic locking via etag concept / investigate rpc-based enqueue server
    - [ ] look at fast hash algorithms (murmur-2??)
  - [x] add Href to entities as a common self-referential field
  - [x] update ReadModel() to accept new model format with relations
  - [ ] update ReadModel() to handle multiple model files
  - [x] add code to support the links via child-href
  - [ ] add code to support expansion of child-href
    - [x] 	Href string  `rgen:"-" json:"Href,omitempty"`
	- [x]   Test string  `rgen:"-" json:"Test,omitempty"`
  - [ ] add code to support filtering of expansions
  - [ ] add scopes to config
    - [ ] use scopes in JWT to allow / disallow access to routes / actions
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
  - [x] add a flag for model file i.e.   $ go run main.go -m "/Users/tomthedog/config/mymodel.json
  - [x] look at how gorilla.mux handles routes like  ../product?Attr1='foo'&&Attr2
    - *see https://stackoverflow.com/questions/45378566/gorilla-mux-optional-query-values*
  - [x] add support for a dev config.json file
  - [ ] add support for LetsEncrypt
  - [x] add capability of generating keys for JWT via ecdsa256
  - [x] add automated default tests 
  - [x] run go fmt on each file immediately following generation?
  - [x] remove the gorilla csrf dependency; the use of JWT's in a stateless application obviates the need for CSRF protection. 
  - [x] run goimports on generated code  
  - [ ] add the capability of automatically running go get (look at go dep) for missing packages in the dependency list
  - [ ] add capability to generate self-signed certs for local ssl testing
  - [ ] create github repo for gnerated code via https://godoc.org/github.com/google/go-github/github#RepositoriesService
