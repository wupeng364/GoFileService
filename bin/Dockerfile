FROM centos:centos8

WORKDIR /gofs/

# RUN yum install golang

COPY ./conf /gofs/conf
COPY ./datas /gofs/datas
COPY ./static /gofs/static
COPY ./app.ico /gofs/app.ico
COPY ./gofs /gofs/gofs

# CMD [  ]

ENTRYPOINT  ["./gofs"]

VOLUME ["/gofs/datas", "/gofs/conf"]