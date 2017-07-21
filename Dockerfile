FROM golang:latest
RUN apt-get update -qq && apt-get install -y build-essential
RUN mkdir -p /go/src/github.com/ThoughtWorksStudios/datagen
WORKDIR /go/src/github.com/ThoughtWorksStudios/datagen
ENV PATH=$GOPATH/bin:$PATH
COPY .bashrc /root
COPY . .
RUN go-wrapper download
RUN go-wrapper install
