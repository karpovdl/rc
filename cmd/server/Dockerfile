from karpovdl/golang:1.14.4

LABEL author="Denis Karpov" \
      site="github.com/karpovdl" \
      email="karpovdl@hotmail.com" \
      version="1.0.0" \
      release-date="2020-07-05" \
      golang="1.14.4"

ENV JWT_TOKEN="*kz@T23Wa5&,qJD0.&pm"

RUN adduser --disabled-password --gecos '' api
USER api

WORKDIR /go/src/app
COPY ./server .

RUN go get github.com/karpovdl/rc/cmd/server/server
RUN go get ./...

EXPOSE 8080
CMD ["server"]
