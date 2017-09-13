# Moulin. A Redis & Go based task delivery mechanism
#technotes/moulin #technotes/queue

# Introduction
This is part of a project to rethink how we can do queueing. This is based on the earlier design of Moulin, specifically where a client connects to the go-moulin component that connects to Redis. In the Go-Moulin section we describe in some detail how we use Redis as a reliable queue and datastore.

# How it works
*Redis* is responsible for storing the task body, and keeping all metadata

*Moulin* is the component responsible for moving the task through the different states of processing.

*Workers* connect to Moulin to receive tasks, send heartbeats and mark tasks as complete.

*A queue loader*, is responsible for reliably loading items into the queue. This is particularly important when we’re talking about large amounts of sequences (e.g. 1B tasks, Sequenomics).

*The PubSub server* is a component that listens on a Redis PubSub channel and is able to send the events coming by over a websocket to subscribers. For example the dashboard.

```
+---------+               +---------+      +----------+
| Queue   |   +--------+  |         |      | PubSub   |
| Loader  +---+(Kafka?)+-->  REDIS  +------+ Websocket|
|         |   +--------+  |         |      | Ser^er   |
+---------+               +^-^^-^-^-+      +----------+
                           | || | |
                           | || | |
                            | || | |
            +--------------+----+-+--------------+
            |              Go!Moulin             |
            +-+-------+-------+-------+-------+--+
              ^       ^       ^       ^       ^
              |       |       |       |       |
              |       |       |       |       |
            +-+--+  +-+--+  +-+--+  +-+--+  +-+--+
            | W1 |  | W2 |  | W3 |  | W4 |  | W5 |
            +----+  +----+  +----+  +----+  +----+
```

## Workers
Each worker connects to Moulin. For each connection from a worker the API opens a connection to Redis. As soon as Redis returns an item the API in turn, returns the item to the worker. Both Go and Redis easily scale to tens of thousands of connections per node (some RAM is used per connection), and scale of the Moulin API is easy to increase as it’s essentially stateless.

The connection from a given worker to Moulin uses GRPC with protobuff; which has the benefit of not only being fast and performant (which is not super needed), but also that the connection essentially can be seen as a two-way socket. This means that the worker can easily do a blocking wait on moulin.

## Moulin

```
                                              +---------+
   REDIS                                     completed set
                                              |    |    |
                                             score | member
                                              |    |    |
+ -----+          +--------+                  |    |    |
incoming          data store   +---------+    +---------+
| list |          |        |   running set    +---------+
|      |  +----+  |id1: { }|   |    |    |    failed set
|{....}| received |id2: { }|  score | member  |    +    |
|{....}|  +list+  |id3: { }|   |    |    |   score | member
|      |  |    |  |        |   |    |    |    +    |    +
+ -----+  +----+  +--------+   +----+----+    +----+----+


            +---------------+
            |     Moulin    |
            +---------------+
```

There are three phases for the queueing system:

### Task loading
This is by far the most complicated task.

1. Check if an expired task is available [running set] (and return it), at the same time increment the counter queue_id.task_id.retries [data store]. If the amount of retries exceeds a limit; push the task_id to a ‘failed’ set.
2. Check if a task was received [received list], but not returned yet (exception) (and return it)
3. Do a blocking wait on the incoming queue [incoming list], and push it to received [received list] when received.
4. Pop the item from the received queue, save it to the data store and add it to the running set, and return the task, /All in one transaction./ Return the object (only) if the ID’s match. Restart task loading if nil was returned.

### Task Heartbeat

1. Update the member in the running set with a score (time) that is 5 minutes in the future. Return “ok”

### Task Completion
1. Remove the member from the running set, and move it to the completed set.

## How Expiration works
For this we use the datatype /sorted set/. Each item in the set is given an integer for the score. We add the taskID’s as items with a score of unixtime + 300s. This way when we later check for all items that have score smaller than the current timestamp, we’ll get all expired items ([ZRANGEBYSCORE](https://redis.io/commands/zrangebyscore) )

## How we guarantee reliability
For the situation when a disconnect (anywhere) might happen, we make sure we always use atomic updates of the dataset, where, for example, a task is either popped AND saved to the running set, or nothing at all. A blocking wait cannot be done in eval(), so we first push it to a received set, and only then pop it.

Redis itself has clustering capabilities with replication. For a single node Redis can be set to write every transaction to disk (fsync on write). Or write every second, which is faster.

## Queue Loader
To be described

## PubSub server
To be described
