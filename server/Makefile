.PHONY: test test-v build run vet fmt pg-up pg-down

# 集成测试使用独立的 go_admin_test 库（首次运行自动创建并 seed），绝不触碰开发库 go_admin，
# 因此 migrate 重建表不会清空开发数据。-p 1 串行执行各包，避免 migrate 重建表与其它包测试并发打架。
test:
	go test -p 1 -count=1 ./...

test-v:
	go test -p 1 -count=1 -v ./...

build:
	go build ./...

vet:
	go vet ./...

fmt:
	gofmt -w .

run:
	go run ./cmd/server

# 本地用 Docker 起停 PostgreSQL
pg-up:
	docker start go-admin-pg || docker run -d --name go-admin-pg -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=go_admin -p 5432:5432 postgres:16

pg-down:
	docker stop go-admin-pg
