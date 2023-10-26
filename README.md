# Moulin

'Moulin' is 'Mill' in French. A Dutch saying goes: *wie het eerst komt, wie het eerst maalt*. Who brings their wheat to the mill first, is the first to get his/her wheat milled, and is the first to eat.

For information how it works, read [here](https://github.com/dhrp/moulin/blob/master/how_it_works.md)

## Running the service

### Run the server

First make sure you have Redis running. For example like so:

```
docker run -p 6379:6379 --name redis -d redis
```

Then start the moulin server. In this example we link to the redis container that we just started.

```
docker run --link redis:redis -e REDIS_ADDRESS=redis:6379 --name moulin -d dhrp/moulin
```

Specify the REDIS_ADDRESS, to point to a Redis server (defaults to localhost)

### Run the client

The client is also built into the docker image. You can run it with 

```
docker run --entrypoint moulin-cli --link moulin -e MOULIN_SERVER=moulin:8042 dhrp/moulin work queue1 forever
```

Here we tell use `moulin-cli` to listen to a queue 'queue1' and execute each command that comes from the queue.


## Using Moulin

### How to use it

Moulin primarily meant to simply execute a given task on the shell; and as such it is language agnostic. Let's say you have a python script "download_files.py" which takes one argument, a URL to download. You can create a task like:

```
moulin-cli create queue1 "./download_files.py https://myserver.org/file2.bin"
```

You can then also create a worker like so:
```
moulin-cli work queue1 forever
```

It will listen on a socket from the moulin server, pop a task from the queue and processes it. When no more tasks are available it will wait for the next to show up.

### Reliability
Once a Moulin worker takes a task from the queue it moves to a 'in process' queue, and only moves to "succeeded" if the child process exits with error code 0. The worker is also expected to send a regular heartbeat to the server to make sure the process didn't die unexpectedly. After 5 minutes of no heartbeat the task is considered failed. 

See [how it works](how_it_works.md)


Functions available:

- **create**: Create a task to be executed.
- **progress**: Check the progress of a certain queue
- **peek**: Have a look at the latest N items in the queue



## Building

## Creating and starting the server
You can build this image with Docker, `docker build -t dhrp/moulin .` and then just run it like so:

```
make image
# docker build -t dhrp/moulin .
```

```
make build
make cli
```


## GRPC

When the communication protocol changes; GRPC definitions need to be updated. Generate a new protocol buffer definition with a version of the following command:
```
$ protoc -I helloworld/ helloworld/helloworld.proto --go_out=plugins=grpc:helloworld
```


### Uploading a file for task batch:

Currently we're supporting multipart file uploads; this is what you'd get when you'd
upload a file through a form. Probably it would make more sense to (just) support receiving
a text/plain body which has only the content of the file (lines with instructions).
```
http --verify no --form POST https://localhost:8042/v1/task_list/batch/ file@./kafkaproducer/test/testtextfile.txt -v
```


## further reading:

* [grpc on kubernetes](https://github.com/kelseyhightower/grpc-hello-service/tree/master/Tutorials/kubernetes)

* [grpc on k8s video](https://vimeo.com/190648663)
