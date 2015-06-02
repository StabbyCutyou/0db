0db
==============

It's 0% a database!

Disclaimer
==============

0db is currently in open-alpha. The interfaces, both internal and external, are likely to change. You should avoid using this in production for anything critical to the business until the initial v1.0.0 release

What is 0db
================

At the moment, not much. As an open-alpha offering which is still under development, only the core aspects of storing and accessing data are supported: Eventually Consistent Reads and Writes.

A word about Eventual Consistency
=================================

0db makes no gaurentees about the availability of your data at the time you request it. You should take care to design your system to be resilient to this aspect of 0db.

What will 0db be?
=================

0db will be a fully distributed system, replicating written data across nodes in the cluster, and coordinating across the cluster to serve read requests. It uses a novel take on the industry standard Paxos consensus algorithm, dubbed Slaxos. Slaxos is optimized to the specific set of needs that 0db exposes, enabling it to cut through the waste.

Interacting with 0db
====================

0db comes provided with a simple REST interface which accepts reads and writes over HTTP. Currently, 0db provides a single endpoint for key-value reads and writes, which is located at:

* 127.0.0.1:5050/v1/{key_name}

## Storing Data

You can store a value securely in 0db by issuing a POST with the data to store residing in the request body. Currently, 0db is a "last-write-wins" system, so take care not to overwrite your keys.

The response to this call with be a JSON object containing the data that was stored, along with the key it was stored at:

```javascript
{"key":"xxx", "data":"yyy"}
```

## Retrieving Data

You can retrieve a value from 0db by issuing a GET.

It will return to you a JSON object with the key and data available for that key.

```javascript
{"key":"xxx", "data":"yyy"}
```

As 0db is an eventually consistent system, there is no gaurentee that your data will be available to be read at the time your request it. In this case, you will receive a response like so:

```javascript
{"key":"xxx", "data":""}
```

## Errors

If the system cannot process the request, you'll receive a response with the appropriate status code, and a JSON object containing the error, like so

```javascript
{"error":"It failed because reasons"}
```

Roadmap
=======

* Official Benchmarks
* Protobuffs interface
* Clustering
* CRDTs
* Officially move to using gb
