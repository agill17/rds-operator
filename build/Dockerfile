FROM alpine:3.8
RUN apk --update upgrade && apk add ca-certificates
RUN apk -v --update add python py-pip jq \
    && pip install awscli 

USER nobody

ADD build/_output/bin/rds-operator /usr/local/bin/rds-operator
