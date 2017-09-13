# Moulin

'Moulin' is 'Mill' in French. A Dutch saying goes: *wie het eerst komt, wie het eerst maalt*. Who brings their wheat to the mill first, is the first to get his/her wheat milled, and is the first to eat.

For information how it works, read [here](https://github.com/nerdalize/moulin/blob/master/how_it_works.md)

## Usage

You can build this image with Docker, `docker build -t moulin .` and then just run it. The client is also built into the container. You can run it with `docker run --entrypoint /go/bin/client dhrp/moulin [hostname] [port]`


## building

Generate a new protocol buffer definition with a version of the following command:
$ protoc -I helloworld/ helloworld/helloworld.proto --go_out=plugins=grpc:helloworld


## further reading:

* [grpc on kubernetes](https://github.com/kelseyhightower/grpc-hello-service/tree/master/Tutorials/kubernetes)

* [grpc on k8s video](https://vimeo.com/190648663)
