FROM golang:1.22

WORKDIR /app

COPY ./challengeA/challengeA.go /app/challengeA/challengeA.go
COPY ./challengeB/challengeB.go /app/challengeB/challengeB.go
COPY ./run.sh /app/

RUN go build -o challengeA ./challengeA/challengeA.go
RUN go build -o challengeB ./challengeB/challengeB.go

CMD ["/bin/sh", "./run.sh"]