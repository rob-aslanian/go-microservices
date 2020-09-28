FROM alpine:3.8

WORKDIR /opt

ADD ./build ./
ADD ./graphiql/libs graphiql/libs/

ENTRYPOINT ["./graphql_server"]
