# ---- build stage ----
FROM golang:1.26-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -o /out/server ./cmd/server

# ---- run stage ----
FROM alpine:3.20
WORKDIR /app
# 时区/证书（连 PostgreSQL 不需要证书，但保留以备外连）
RUN apk add --no-cache tzdata ca-certificates
COPY --from=build /out/server /app/server
COPY configs /app/configs
EXPOSE 8080
# 容器内 DB 主机走 compose 服务名，通过环境变量覆盖 config.yaml：
#   GOADMIN_DATABASE_HOST=db  GOADMIN_JWT_SECRET=...  GOADMIN_ADMIN_PASSWORD=...
CMD ["/app/server"]
