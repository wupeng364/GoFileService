FROM golang

WORKDIR /go/bin/

COPY ./src /go/src
COPY ./bin/conf /go/bin/conf
COPY ./bin/datas /go/bin/datas
COPY ./bin/static /go/bin/static
COPY ./bin/app.ico /go/bin/app.ico

RUN cd /go/src/gofs \
        && go install

# CMD [  ]

ENTRYPOINT  ["./gofs"]

VOLUME ["/go/bin/datas", "/go/bin/conf"]