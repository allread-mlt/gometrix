FROM golang:alpine

ARG USERNAME=vscode

RUN apk add --no-cache git curl unzip

WORKDIR /gometrix

RUN go install github.com/go-delve/delve/cmd/dlv@latest \
    && go install -v golang.org/x/tools/gopls@latest
