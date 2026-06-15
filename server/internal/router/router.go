package router

import (
	"go-admin/internal/auth"
	"go-admin/internal/config"
	"go-admin/internal/handler"
	"go-admin/internal/middleware"
	"go-admin/internal/repository"
	"go-admin/internal/service"
	"go-admin/pkg/captcha"
	"go-admin/pkg/jwt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type sessionAdapter struct{ repo *repository.SessionRepo }

func (s sessionAdapter) Exists(jti string) bool {
	_, err := s.repo.GetByID(jti)
	return err == nil
}
func (s sessionAdapter) Touch(jti string) error { return s.repo.Touch(jti) }

func Setup(cfg *config.Config, db *gorm.DB, jm *jwt.Manager) *gin.Engine {
	r := gin.Default()

	// repositories
	userRepo := repository.NewUserRepo(db)
	roleRepo := repository.NewRoleRepo(db)
	menuRepo := repository.NewMenuRepo(db)
	dictRepo := repository.NewDictRepo(db)
	sessRepo := repository.NewSessionRepo(db)
	logRepo := repository.NewLogRepo(db)

	// services
	authSvc := service.NewAuthService(userRepo, sessRepo, jm, cfg.Login)
	menuSvc := service.NewMenuService(menuRepo, userRepo)
	userSvc := service.NewUserService(userRepo)
	roleSvc := service.NewRoleService(roleRepo)
	menuAdminSvc := service.NewMenuAdminService(menuRepo)
	dictSvc := service.NewDictService(dictRepo)
	logSvc := service.NewLogService(logRepo)

	// enforcer + handlers
	enforcer := auth.NewDBEnforcer(db)
	authH := handler.NewAuthHandler(authSvc, menuSvc, captcha.NewStore(), cfg.Captcha.Enabled)
	userH := handler.NewUserHandler(userSvc)
	roleH := handler.NewRoleHandler(roleSvc)
	menuH := handler.NewMenuHandler(menuAdminSvc)
	dictH := handler.NewDictHandler(dictSvc)
	onlineH := handler.NewOnlineHandler(sessRepo)
	profileH := handler.NewProfileHandler(userSvc)
	logH := handler.NewLogHandler(logSvc)
	monitorH := handler.NewMonitorHandler()
	checker := sessionAdapter{repo: sessRepo}

	api := r.Group("/api")
	{
		api.GET("/captcha", authH.Captcha)
		api.POST("/auth/login", middleware.RateLimitByIP(cfg.Login.RateLimitPerMinute), authH.Login)
		api.POST("/auth/refresh", authH.Refresh)

		authed := api.Group("")
		authed.Use(middleware.Auth(jm, checker))
		authed.Use(middleware.OperationLog(logRepo))
		{
			authed.POST("/auth/logout", authH.Logout)
			authed.GET("/auth/me", authH.Me)
			authed.GET("/menus/me", authH.MenusMe)
			authed.PUT("/profile", profileH.UpdateProfile)
			authed.PUT("/profile/password", profileH.ChangePassword)

			perm := func(p string) gin.HandlerFunc { return middleware.RequirePerm(enforcer, p) }

			// 用户
			authed.GET("/users", perm("system:user:list"), userH.List)
			authed.POST("/users", perm("system:user:create"), userH.Create)
			authed.PUT("/users/:id", perm("system:user:update"), userH.Update)
			authed.DELETE("/users/:id", perm("system:user:delete"), userH.Delete)
			authed.PUT("/users/:id/roles", perm("system:user:assignrole"), userH.AssignRoles)
			authed.PUT("/users/:id/password", perm("system:user:resetpwd"), userH.ResetPassword)
			authed.PUT("/users/:id/unlock", perm("system:user:update"), userH.Unlock)

			// 角色
			authed.GET("/roles", perm("system:role:list"), roleH.List)
			authed.GET("/roles/:id", perm("system:role:list"), roleH.Get)
			authed.POST("/roles", perm("system:role:create"), roleH.Create)
			authed.PUT("/roles/:id", perm("system:role:update"), roleH.Update)
			authed.DELETE("/roles/:id", perm("system:role:delete"), roleH.Delete)
			authed.PUT("/roles/:id/menus", perm("system:role:assignmenu"), roleH.AssignMenus)

			// 菜单
			authed.GET("/menus", perm("system:menu:list"), menuH.Tree)
			authed.POST("/menus", perm("system:menu:create"), menuH.Create)
			authed.PUT("/menus/:id", perm("system:menu:update"), menuH.Update)
			authed.DELETE("/menus/:id", perm("system:menu:delete"), menuH.Delete)

			// 字典
			authed.GET("/dict/types", perm("system:dict:list"), dictH.ListTypes)
			authed.POST("/dict/types", perm("system:dict:create"), dictH.CreateType)
			authed.PUT("/dict/types/:id", perm("system:dict:update"), dictH.UpdateType)
			authed.DELETE("/dict/types/:id", perm("system:dict:delete"), dictH.DeleteType)
			authed.GET("/dict/data", perm("system:dict:list"), dictH.DataByType)
			authed.POST("/dict/data", perm("system:dict:create"), dictH.CreateData)
			authed.PUT("/dict/data/:id", perm("system:dict:update"), dictH.UpdateData)
			authed.DELETE("/dict/data/:id", perm("system:dict:delete"), dictH.DeleteData)

			// 在线用户
			authed.GET("/online", perm("monitor:online:list"), onlineH.List)
			authed.DELETE("/online/:id", perm("monitor:online:kick"), onlineH.Kick)

			// 操作日志
			authed.GET("/logs/operation", perm("monitor:log:list"), logH.List)
			authed.DELETE("/logs/operation", perm("monitor:log:delete"), logH.Clear)

			// 服务监控
			authed.GET("/monitor/server", perm("monitor:server:list"), monitorH.Server)
		}
	}
	return r
}
