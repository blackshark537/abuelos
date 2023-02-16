FROM --platform=amd64 golang:alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY ./src ./src
COPY *.go ./

RUN go build -o /dataprod
#RUN CGO_ENABLED=0 go build -ldflags=”-s -w” -o /dataprod

FROM --platform=amd64 golang:alpine
WORKDIR /root/
COPY --from=builder /dataprod /

CMD [ "/dataprod", "serve" ]
