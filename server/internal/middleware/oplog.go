package middleware

import (
	"time"

	"go-admin/internal/model"
	"go-admin/internal/repository"

	"github.com/gin-gonic/gin"
)

// OperationLog 异步记录变更类（POST/PUT/DELETE）请求的操作日志。
func OperationLog(repo *repository.LogRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		if method != "POST" && method != "PUT" && method != "DELETE" {
			c.Next()
			return
		}
		start := time.Now()
		c.Next()

		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}
		entry := model.OperationLog{
			Username:  c.GetString(CtxUsername),
			Method:    method,
			Path:      path,
			IP:        c.ClientIP(),
			Status:    c.Writer.Status(),
			LatencyMs: int(time.Since(start).Milliseconds()),
			CreatedAt: time.Now(),
		}
		go func() { _ = repo.Create(&entry) }()
	}
}
