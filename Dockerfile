FROM golang:latest
RUN apt-get update -qq && apt-get install -y build-essential
RUN mkdir -p /go/src/github.com/ThoughtWorksStudios/bobcat
WORKDIR /go/src/github.com/ThoughtWorksStudios/bobcat
ENV PATH=$GOPATH/bin:$PATH
COPY .bashrc /root
COPY . .
RUN go-wrapper download
