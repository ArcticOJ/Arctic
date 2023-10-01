FROM node:18-alpine AS avalanche-builder
WORKDIR /usr/src/avalanche

COPY avalanche .

RUN yarn install --frozen-lockfile --immutable

RUN yarn build

FROM golang:alpine AS arctic-builder
WORKDIR /usr/src/app

COPY . .

COPY --from=avalanche-builder /usr/src/avalanche/out avalanche/out

RUN go mod download

RUN go build -o ./out/arctic -tags ui,headless -ldflags "-s -w" .

FROM alpine
WORKDIR /arctic

COPY --from=arctic-builder /usr/src/app/out/arctic ./

ENTRYPOINT ["/arctic/arctic"]

