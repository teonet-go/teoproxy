
# Copyright 2024 Kirill Scherba <kirill@scherba.ru>.  All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.
#
# Teonet fortune-gui server application docker file
#
# Docker build (local):
# ---------------------
#
#   # Run this docker build command from root teoproxy directory
#   docker build -t fortune-gui -f ./cmd/websocket/client-gui/serve/Dockerfile .
#
# Docker build (included private repositories):
# --------------------------------------------
#
#   docker build --build-arg github_user="${USERNAME}" --build-arg github_personal_token="${TOKEN_FOR_REPOSITORIES}" -t fortune-gui -f ./cmd/websocket/client-gui/serve/Dockerfile .
#
# Docker test run:
# ---------------
#
#   docker run --network host -it --rm fortune-gui serve
#   docker run --network host -it --rm fortune-gui serve -laddr=localhost:9093 -gzip
#
#   docker run --name=fortune-gui --network host -it -d --restart=always -v $HOME/.config/teonet:/root/.config/teonet fortune-gui serve -laddr=localhost:9092 -gzip -monitor=XXXXX
#   docker logs -f fortune-gui
#
# Publish to github:
# -----------------
#
#   docker login ghcr.io -u USERNAME -p TOKEN_FOR_PACKAGES
#   #
#   docker tag fortune-gui ghcr.io/teonet-go/fortune-gui:0.0.0
#   docker push ghcr.io/teonet-go/fortune-gui:0.0.0
#   #
#   docker tag fortune-gui ghcr.io/teonet-go/fortune-gui:0.0.0
#   docker push ghcr.io/teonet-go/fortune-gui:0.0.0
#

# Build stage
#
FROM golang:1.21.5-alpine AS builder
RUN apk add --no-cache git
WORKDIR /go/src/app
COPY . .

# Add the keys from arguments
# ARG github_user
# ENV github_user=$github_user
# ARG github_personal_token
# ENV github_personal_token=$github_personal_token

# Change github url
# RUN git config \
#     --global \
#     url."https://${github_user}:${github_personal_token}@github.com".insteadOf \
#     "https://github.com"

RUN go install fyne.io/fyne/v2/cmd/fyne@latest
RUN cd ./cmd/websocket/client-gui/serve && go get
RUN cd ./cmd/websocket/client-gui/serve && fyne package -os wasm --sourceDir ../
RUN cd ./cmd/websocket/client-gui/serve && go install -tags=prod .

# Final stage
#
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/bin/* /usr/local/bin/
CMD "serve"
LABEL Name=fortune-gui Version=0.0.1
