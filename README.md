# Shipment Service - backend

This repository is designed to showcase my take on how to design a micro-service.

It will talk about the lifecycle, documentation and testing and give some thoughts around issues that won't be covered, like authentication, distributed tracing and deployment.

The service will have a REST API and is designed around being a multi-tenant solution.

**Table of Contents**

- [Getting Started](#getting-started)
- [Configuration](#configuration)
  - [Run the Shipment-Service](#run-the-shipment-service)
- [File Structure](#file-structure)
- [Choices](#choices)
- [Thoughts](#thoughts)
- [Deployment](/docs/deployment.md)
- [Sequence Diagrams](/docs/sequence-diagrams.md)

## Getting Started

### Required Tools

- golang v1.14 or newer
- make

### Run the Shipment-Service

> `> make run-local`

This make target will start the service locally on serve the REST API on port 8080.

A few links to the running service

- [Health Status](http://localhost:8080/healthy/status)
- [Swagger API](http://localhost:8080/v1/docs/swagger/index.html)

### Run Unit Tests

> `> make test`

### Run Behaviour Specifications as Tests

> `> make run-behaviour-test`

This make target will use [godog](https://github.com/cucumber/godog) to execute the defined behaviours as tests against the service.

The service is required to be running on port 8080.

### Generate Open API spec. and UML diagrams

> `> make gen`

This make target will use [swaggo](https://github.com/swaggo/swag) and [go-swagger](https://github.com/go-swagger) to generate and validate the Open API spec. It will also use [gopuml](https://github.com/lonnblad/gopuml) to generate `.svg` files from the Plant UML diagrams.

### Run Linter

> `> make lint`

This make target will use [golangci-lint](https://github.com/golangci/golangci-lint) to lint the go-code.

### Run Autoformatter

> `> make fmt`

## Configuration

### REST API Service

These environment variables can be used to configure the REST API Service.

```
- ENVIRONMENT       Name of the deployment environment. Defaults to "local", needs to be one of local, sandbox, staging, production.
- SERVICE_NAME      Name of the service. Defaults to "shipment-service".
- SERVICE_VERSION   The version of the service. Defaults to "dev".
- REST_PORT         The port to serve the REST API on. Defaults to 8080.
- REST_URL          The Base URL which the API can be reached on. Defaults to http://localhost:8080.
- SHUTDOWN_TIMEOUT  The timeout before forcing the service to shutdown. Defaults to 20 seconds.
```

## File Structure

```
└─ shipment-service-backend
   ├─ behaviour         # Behaviour specifications
   ├─ boundaries        # Entrypoints into the Shipment Service
   │  └─ rest              # The REST server boundary
   │     ├─ utils             # Utils pkg for REST interfaces
   │     └─ v1                # v1 of the REST interface
   ├─ businesslogic     # The Businesslogic of the Shipment Service
   │  ├─ models             # Internal data models
   │  └─ price              # Price Calculation pkg
   ├─ cmd               # All binaries
   │  ├─ rest-api          # The REST API
   │  └─ rest-api-test     # The behaviour test
   ├─ config            # Configuration pkg
   ├─ docs              # Dedicated documentation
   │  └─ diagrams          # UML Diagrams
   ├─ storage           # Storage interfaces and data structures
   │  └─ go-memdb          # go-memdb implementation of the ShipmentStorage
   └─ trace             # A utility trace pkg
```

## Choices

### Why is there a Tenant ID in the endpoints?

I see this as being a multi tenant service and with this comes mainly two factors in play, security and scalability.

Adding a Tenant ID to all data gives you the opportunity to secure the data by tenant even down to database queries, not all Databases have support for this, but as an example, AWS DynamoDB and PostgreSQL have. It will also provide you with the base for a mechanic where you can choice which regions to use when replicating the data for a given tenant, as an example, it might be that a US customer doesn't want their data stored in Russia or China and this can be guarded with the help of a Tenant ID.

What about scalability, by using a Tenant ID, you can easier co-locate data in the world both in terms of database partitions, but also on region. Say that a given tenant is a Swedish customer, we happen to use Stockholm, but for some reason, this tenant was routed to Ireland instead for two out of 1000 shipments in a month, then a cleanup job could move those two shipments to the Stockholm and with that co-locate all data belonging to that tenant.

When scaling to many regions across the globe, it's not really worth the cost of storing all data in every region, like with AWS DynamoDB's Global Table, but instead distribute the data for a tenant based on access patterns and redundancy, using region-lookup tables to show which region that has a copy of the given tenant.

You can even provide services based on tenant, like extra redundancy or offline capabilities that checks in asynchronously with the cloud on new shipments, a great feature for customers with poor internet connection.

### Behaviour Driven Development

In the [behaviour folder](/behaviour), you can find behaviour specifications for the create shipment endpoint.

These should be viewed as an example to show how one could work with documenting the expected behaviour of a service. In a real world scenario, they should cover all endpoints and more price and validation cases.

For further reading on godog as the tool for running behaviour specifications or for BDD (Behaviour Driven Development) in general. These are some good links.

- [godog](https://github.com/cucumber/godog)
- [cucumber.io](https://cucumber.io/)

### The price package

In [price.go](/businesslogic/price/price.go), you can find the implementation of the price rules.

This is also the package which got real unit testing, instead of just using the Behaviour specification as tests. The reasoing behind this is because this is a business critical equation, which if it calculates the wrong thing will make us loose money. In this case, the price rules are simple so we could test them fairly easy using a Behaviour specification, but in the case where the complexity is greater and far more complex, I believe it's good to test this as it's own package.

### Storage (in-memory)

In [storage.go](/storage/storage.go) you will find a general ShipmentStorage interface{}, being used in [rest-api/main.go](/cmd/rest-api/main.go). There is currently only one implementation [go-memdb](/storage/go-memdb/memdb.go), which is an in-mem database package. However, since this structure uses interfaces, we can simply add an implementation of the ShipmentStorage for AWS DynamoDB or Mongo.

## Thoughts

### gRPC vs. REST vs. GraphQL

I don't have strong opinions for any of these technologies or against one of them. What I have come to realize is that it depends on your developers and your consumers needs.

Personally, I believe that if the developers spend the time in understanding gRPC for inter-service communication when using containers, it greatly improves the speed at which those services can communicate compared to REST or GraphQL. I have been in a situation, where some developers didn't take the time, which triggered multiple lambdas spinning up, trying to connect, failing, spinning up again and so on, ultimately ending with a Denial of Service of a core function.

However in a setup with services using gRPC and then adding an API Gateway to front those services, you can have the API Gateway being a thin bridge from whatever technology you want that translates into gRPC and back. Could be REST, multi-service GraphQL queries, WebSocket, or just plain gRPC.

### Event Source

In this Shipment Service example, I store the shipment data as plain data structures, nothing strange really.

However what I have come to realize is that Event Source is a pretty neat thing in terms of managing a ledger of the lifecycle of a data entity. A simple example, a purchase order, starts with creating the order, the order then need to be payed for, then need to be fulfilled and then need to be shipped. All those are asynchronous events that happen in a order, but instead of updating a plain order data object, you can update a ledger with all those events and then aggregate the ledger into the data object.

This is Martin Fowler talking about [Event Sourcing](https://youtu.be/STKCRSUsyP0?t=1271), I believe the whole talk is pretty good talking about the "The Many Meanings of Event-Driven Architecture".
