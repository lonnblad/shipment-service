# Deployment

In this chapter I will try to present my reasoning around deployment and scaling.

I will use AWS services as references.

## REST API Service

This service is designed to be running as a container deployed on a tool for container orchestration, like the AWS ECS service. It is also designed to handle path routing between API versions by itself instead of using a tool like NGINX for that.

With some modifications, it could instead be deployed using a serverless tool such as AWS Lambda together with AWS API Gateway.

### Why containers over serverless for a REST API?

This is a question with a long answer, which I believe comes down to one question; How do you value worst-case latency vs. cost?

Cold starts is the major issue for Serverless, AWS Lambda adds 150+ ms for every cold start, for golang, it's fairly common to be between 250-500 ms, but in some cases the cold starts could move up towards 1000 ms. It depends on your binary size and what your lambda do during initialization in the main function. When creating a REST API Service based on the AWS API Gateway and AWS Lambda, you could also be using a custom Authorizer lambda, doing this could make one API call end up with two could starts, with a fairly large total.

What about cost, using AWS API Gateway and AWS Lambda together is basically free when you don't have traffic, while AWS ECS will cost you money even with no traffic. If you value the extremely simple auto-scaling, the cost is lower for serverless and cold starts doesn't matter to much, use AWS Lambda instead of AWS ECS.

### AWS ECS - Fargate Spot

AWS ECS - Fargate is basically the serverless version of Container Orchestration and when using **Spot** instances, you greatly reduce the cost. When running a REST API as a stateless service, you typically will be able to finish all requests within seconds and since Fargate Spot will give you a shutdown notice up to 120 seconds prior to the container being forced to shutdown, this is a good option that will allow you to have a simple configuration, like AWS Lambda, but deploy your service as a long-running container.

### Scalability

**1000 monthly active users using the application once per day.**

If all these requests happen evenly distributed over one hour, that would mean ~20 requests per minute, I believe this could be served by one of the smaller Fargate Spot configurations, but should be deployed with one spare for redundancy, all in the same AWS region.

**10.000.000 monthly active users using the application 100 times per day.**

If all these requests happen evenly distributed over 24 hours, that would mean ~700'000 requests per minute. I believe this could be served by several instances of one of the smaller Fargate configurations distributed across multiple AWS Regions. Multiple regions would provide reduced latency, increased redundancy and ability to ramp up and down based on needs in the time zones served most by a given region.

These metrics need to be actively monitored from launch and used to pro-actively make choices on improvements to code, infrastructure, deployment patterns, etc.

**5 developers. vs. 30 developers and Business is just getting started. vs. Business is more mature.**

When scaling a business and a development organization together with a software architecture, I believe it's key to have clear areas of responsibility and autonomy through the whole tech stack within one or a few teams. I believe this is how you truly build a scalable organization and with that a scalable software archictecture. You can't do one without the other.

I also believe that to create great solutions, everyone needs to be working in the trenches together, from taking part in support cases, non-technical staff supporting developers when the platform is on fire or developers being on call for operations.

### Data Storage

Data Storage is a hard topic with solutions ranging from storing raw data in just a file on a service like Amazon S3, to using SQL or NOSQL databases like AWS RDS and AWS DynamoDB, to graph-databases like Dgraph or an event streaming platform like Apache Kafka to name some.

Access-patterns is key, reads, writes, updates and deletes. How much do you know about your queries beforehand, does distributed storage on multiple regions play a role.

As an example for this specific case, we query everything by just a Tenant ID or with a Shipment ID, so we need an index on Tenant ID and maybe a combinatory index on Shipment ID to optimize search queries. If we add an update or a delete, those indexes would still be fine. So far NOSQL using AWS DynamoDB is a good option, it's simple, distributed, scalable, it has global replication for redundancy and pretty fast.

But say we want to search on all shipments given a Tenant ID and either the Country Code of the sender or receiver and all of a sudden we want to add two more indexes for those combination. This will greatly impact the use of AWS DynamoDB, while SQL using PostgreSQL could probably handle the extra index like nothing.

Say that we instead want to query relationships between our different senders and receivers to be able to optimize pickups and deliveries, then a graph-database would be extremely efficient to provide those kinds of insights.

If our organization is investing in a more event-driven style of architecture, than Kafka could be a better solution to help facilitating communication between services as well as within a given service.

To sum up, I believe that it's better optimize on as few things as possible. Say that we want to have a CRUD REST API for shipments, where we also can run optimzation algorithms on pickups and deliveries and publish shipment events to an invoice-service.

Then I would have the CRUD API using AWS DynamoDB, I would use Dgraph to do network-queries on senders and receivers for transport optimization and I would use Kafka or maybe AWS SNS for publishing events to other services.
