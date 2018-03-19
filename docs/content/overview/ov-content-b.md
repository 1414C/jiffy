---
title: "Jiffy Application Overview"
date: 2018-02-05T13:26:10-07:00
weight: 20
draft: true
---

### What does a generated Jiffy application look like?

Generated Jiffy applications can be pointed at the DBMS of your choice without the need to recompile the binary (architecture differences not withstanding).  This means that a developer can build a model, fully test it locally using SQLite and then redirect the appplication to a formal testing environment running SAP Hana, or any of the other supported database systems.

Applications are generated based on model files which are encoded as simple JSON.  The concepts of entity and resource-id form the cornerstones upon which the model, application and RESTful end-points are built.

Entities can be thought of anything that needs to be modelled; Order, Customer, Invoice, Truck, ..., ... Each entity is mandated to have an ID field, which is analagous to a primary-key or row-id in the backend database.  ID is used as the primary resource identifier for an entity, and is setup by default as an auto-incrementing column in the database.  ID is implemented as go-type uint64 and is inserted into the model entity definition during application generation.

Accessing an entity via the generated CRUD interface is very simple.  For example, a customer entity could be defined in the model and then accessed via the application as follows:

1. Create a customer entity:
    - https://servername:port/customer  + {JSON body}

2. Update a customer entity:
    - https://servername:port/customer/:id  + {JSON body}

3. Read a customer entity:
    - https://servername:port/customer/:id

4. Delete a customer entity:
    - https://servername:port/customer/:id

5. Read all customer entities:
    - https://servername:port/customers


Additional routes can also be generated based on the model file, including custom filters for GET operations, static end-points for common GET operations, HasOne, HasMany and BelongsTo relationships:

1. Use a filter to Get customers where the last name is 'Smith':
    - https://servername:port/customers/?last_name=Smith

2. Use a generated static end-point to Get customers where credit score is less than 4:
    - https://servername:port/customers/credit_score(LT 4)

3. Use a generated relationship to retrieve all orders for a specific customer (10023):
    - https://servername:port/customer/10023/orders

4. Use a generated relationship to retrieve a specific order (99000022) for the specified customer (10023):
    - https://servername:port/customer/10023/order/99000022

5. Use a generated belongsTo relationship to retrieve the customer for a specific order (990000222):
    - https://servername:port/order/99000022/customer


A set of commands can be appended to an operation's URL to perform some common activities.  The commands can be appended to the URL in any order.

1. Get a count of all customer entities:
    - https://servername:port/customer/$count

2. Limit the number of returned customer entities to 3.  The default ordering for this example would be ascending based on the entity ID field.
    - https://servername:port/customer/$limit=3

3. Offset the database selection by 2 records, top-down, using the default order; (ascending based on the entity ID field):
    - https://servername:port/customer/$offset=2

4. Select records in descending order based on the entity ID field:
    - https://servername:port/customer/$desc

5. Select records in descending order using the customer name field as the order-by criteria:
    - https://servername:port/customer/$orderby=name$desc

6. Limit the number of selected records to 3 and select in descending order based on the entity ID field:
    - https://servername:port/customer/$limit=3/$desc

7. Limit the number of selected records to 3 and select in ascending order using the customer name field as the order-by criteria:
    - https://servername:port/customer/$limit=3;$orderby=name$asc

8. Limit the number of selected records to 3 with an offset of 2 and select the records in ascending order using the customer name field as the sort criteria:
    - https://servername:port/customer/$limit=3$offset=2$orderby=name$asc

9. Limit the return of a static filter end-point to 3 records:
    - https://servername:port/customers/credit_score(LT4)/$limit=3


More details regarding application modlelling are contained in later sections of this documentation.