FROM golang:1.10.3
WORKDIR /go/src/github.com/alexesDev/cas
COPY . .
RUN cd cmd/cascli && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /cli .

FROM scratch
ENV PATH=/
COPY --from=0 /cli /
