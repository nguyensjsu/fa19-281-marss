
FROM golang:latest

MAINTAINER Apeksha 

# Expose port 8080 to the outside world
EXPOSE 8080

WORKDIR /app

COPY . .

#Getting required imports
RUN go get github.com/gorilla/mux
RUN go get go.mongodb.org/mongo-driver/bson
RUN go get go.mongodb.org/mongo-driver/bson/primitive
RUN go get go.mongodb.org/mongo-driver/mongo
RUN go get go.mongodb.org/mongo-driver/mongo/options
RUN go get github.com/aws/aws-sdk-go/aws
RUN go get github.com/aws/aws-sdk-go/aws/session
RUN go get github.com/aws/aws-sdk-go/service/s3
RUN go get github.com/aws/aws-sdk-go/service/s3/s3manager
RUN go get github.com/aws/aws-sdk-go/aws/awserr
RUN go get github.com/nu7hatch/gouuid
RUN go get github.com/gorilla/handlers

ENV AWS_ACCESS_KEY_ID=<your-access-key>
ENV AWS_SECRET_ACCESS_KEY=<your-secret-key>

# Build the Go app
RUN go build -o main .

# Command to run the executable
CMD ["./main"]
