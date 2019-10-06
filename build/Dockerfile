FROM alpine:3.8
RUN apk --update upgrade && apk add ca-certificates
RUN apk -v --update add python py-pip jq \
    && pip install awscli

ENV OPERATOR=/usr/local/bin/rds-operator \
    USER_UID=1001 \
    USER_NAME=rds-operator

# install operator binary
COPY build/_output/bin/rds-operator ${OPERATOR}

COPY build/bin /usr/local/bin
RUN  /usr/local/bin/user_setup

ENTRYPOINT ["/usr/local/bin/entrypoint"]

USER ${USER_UID}
