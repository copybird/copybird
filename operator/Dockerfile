FROM golang:1.11

ENV GO111MODULE=on

WORKDIR /copybird

COPY . .

RUN curl -sSL https://sdk.cloud.google.com | bash
ENV PATH $PATH:/root/google-cloud-sdk/bin

RUN go build -o operator .

CMD [ "/copybird/operator" ]