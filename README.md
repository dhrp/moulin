# Moulin

'Moulin' is 'Mill' in French. A Dutch saying goes: *wie het eerst komt, wie het eerst maalt*. Who brings their wheat to the mill first, is the first to get his/her wheat milled, and is the first to eat.

For information how it works, read [here](https://github.com/nerdalize/moulin/blob/master/how_it_works.md)

## Usage

You can build this image with Docker, `docker build -t moulin .` and then just run it. The client is also built into the container. You can run it with `docker run --entrypoint /go/bin/client dhrp/moulin [hostname] [port]`


## building

Generate a new protocol buffer definition with a version of the following command:
$ protoc -I helloworld/ helloworld/helloworld.proto --go_out=plugins=grpc:helloworld


## testing

### Starting Redis

For the tests we use redis on localhost, with the default port, and no password.
So you can do something like:
```
docker run -p 6379:6379 --name redis redis -d
```

### Starting the Kafka server is best done a docker-compose file:

clone the repo:
https://hub.docker.com/r/wurstmeister/kafka/

then run `docker-compose up -d`

Please note that when you want change the hostname in the docker-compose file you
may need to delete (fix?) zookeeper or Kafka datastore, as it won't correctly connect
otherwise.

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
