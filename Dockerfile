# syntax=docker/dockerfile:1

FROM golang:1.16-alpine 

# COPY . /app/

WORKDIR /app

# Clone the conf files into the docker container
RUN git clone <https-link>

RUN go mod download 

RUN go build -o /empapp

EXPOSE 8000

CMD [ "./empapp" ]

# docker build --tag empapp . 
