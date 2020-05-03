
FROM golang:latest

MAINTAINER Sindhu

# Expose port 8080 to the outside world
EXPOSE 3000

#RUN mkdir -p /home/centos/.aws
#ADD credentials /home/centos/.aws/credentials
#ADD ./build/libs/starbucks-all.jar /srv/starbucks-all.jar

WORKDIR /app

COPY . .

#Getting required imports
RUN go get github.com/aws/aws-sdk-go/aws
RUN go get github.com/aws/aws-sdk-go/aws/session
RUN go get github.com/aws/aws-sdk-go/service/sns
RUN go get github.com/gorilla/mux
RUN go get github.com/nu7hatch/gouuid
RUN go get go.mongodb.org/mongo-driver/bson
RUN go get go.mongodb.org/mongo-driver/bson/primitive
RUN go get go.mongodb.org/mongo-driver/mongo
RUN go get go.mongodb.org/mongo-driver/mongo/options
RUN go get github.com/gorilla/handlers
RUN export AWS_ACCESS_KEY_ID=AKIAYLO2AOZ5NT3LRPT7
RUN export AWS_SECRET_ACCESS_KEY=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

ENV AWS_ACCESS_KEY_ID=**************
ENV AWS_SECRET_ACCESS_KEY=*****************

# Build the Go app
RUN go build -o main .

# Command to run the executable
CMD ["./main"]
