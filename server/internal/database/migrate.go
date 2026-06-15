package database

import (
	"go-admin/internal/model"
	"go-admin/pkg/password"

	"gorm.io/gorm"
)

// MigrateAndSeed 用默认超管账号 admin/admin123 迁移并种子（测试用）。
func MigrateAndSeed(db *gorm.DB) error {
	return MigrateAndSeedWithAdmin(db, "admin", "admin123")
}

// MigrateAndSeedWithAdmin 迁移并按指定超管账号/初始密码种子（主程序从 config 传入）。
func MigrateAndSeedWithAdmin(db *gorm.DB, adminUser, adminPass string) error {
	if adminUser == "" {
		adminUser = "admin"
	}
	if adminPass == "" {
		adminPass = "admin123"
	}
	if err := db.AutoMigrate(
		&model.User{}, &model.Role{}, &model.Menu{},
		&model.DictType{}, &model.DictData{}, &model.OnlineSession{},
		&model.OperationLog{},
	); err != nil {
		return err
	}
	return seed(db, adminUser, adminPass)
}

func seed(db *gorm.DB, adminUser, adminPass string) error {
	// 预置角色（幂等：FirstOrCreate by code）
	roles := []model.Role{
		{Name: "超级管理员", Code: "super_admin", Sort: 1, Status: 1, Remark: "系统内置，拥有全部权限"},
		{Name: "伪超级管理员", Code: "admin", Sort: 2, Status: 1, Remark: "管理普通用户"},
		{Name: "普通用户", Code: "user", Sort: 3, Status: 1, Remark: "基础访问"},
	}
	for i := range roles {
		if err := db.Where("code = ?", roles[i].Code).
			FirstOrCreate(&roles[i]).Error; err != nil {
			return err
		}
	}

	// 超管账号（幂等：仅首次创建时写入初始密码，之后不覆盖用户改过的密码）
	hashed, err := password.Hash(adminPass)
	if err != nil {
		return err
	}
	admin := model.User{Username: adminUser, Nickname: "超级管理员", Status: 1}
	if err := db.Where("username = ?", adminUser).
		Attrs(model.User{Password: hashed}).
		FirstOrCreate(&admin).Error; err != nil {
		return err
	}
	// 绑定 super_admin 角色
	var superRole model.Role
	if err := db.Where("code = ?", "super_admin").First(&superRole).Error; err != nil {
		return err
	}
	if err := db.Model(&admin).Association("Roles").Replace(&superRole); err != nil {
		return err
	}

	// 示范字典（幂等 by type）
	dictTypes := []model.DictType{
		{Name: "用户状态", Type: "sys_user_status", Status: 1},
		{Name: "用户性别", Type: "sys_user_sex", Status: 1},
	}
	for i := range dictTypes {
		if err := db.Where("type = ?", dictTypes[i].Type).
			FirstOrCreate(&dictTypes[i]).Error; err != nil {
			return err
		}
	}
	if err := seedMenus(db); err != nil {
		return err
	}
	return nil
}

// seedMenus 幂等地种入菜单树并把全部菜单挂到 super_admin。
func seedMenus(db *gorm.DB) error {
	// 用稳定的 perm/name 做幂等键
	upsert := func(m *model.Menu) (*model.Menu, error) {
		key := db.Where("name = ? AND type = ?", m.Name, m.Type)
		if err := key.FirstOrCreate(m).Error; err != nil {
			return nil, err
		}
		return m, nil
	}

	dashboard := &model.Menu{ParentID: 0, Name: "工作台", Type: "C", Path: "/dashboard", Component: "dashboard/index", Icon: "dashboard", Sort: 1, Visible: 1, Status: 1}
	system := &model.Menu{ParentID: 0, Name: "系统管理", Type: "M", Path: "/system", Icon: "system", Sort: 2, Visible: 1, Status: 1}
	if _, err := upsert(dashboard); err != nil {
		return err
	}
	if _, err := upsert(system); err != nil {
		return err
	}

	type leaf struct {
		name, path, comp, perm string
		buttons                map[string]string // label->perm
	}
	leaves := []leaf{
		{"用户管理", "/system/user", "system/user/index", "system:user:list", map[string]string{
			"新增": "system:user:create", "修改": "system:user:update", "删除": "system:user:delete",
			"重置密码": "system:user:resetpwd", "分配角色": "system:user:assignrole"}},
		{"角色管理", "/system/role", "system/role/index", "system:role:list", map[string]string{
			"新增": "system:role:create", "修改": "system:role:update", "删除": "system:role:delete",
			"分配权限": "system:role:assignmenu"}},
		{"菜单管理", "/system/menu", "system/menu/index", "system:menu:list", map[string]string{
			"新增": "system:menu:create", "修改": "system:menu:update", "删除": "system:menu:delete"}},
		{"字典管理", "/system/dict", "system/dict/index", "system:dict:list", map[string]string{
			"新增": "system:dict:create", "修改": "system:dict:update", "删除": "system:dict:delete"}},
	}
	sort := 1
	for _, lf := range leaves {
		menu := &model.Menu{ParentID: system.ID, Name: lf.name, Type: "C", Path: lf.path, Component: lf.comp, Perm: lf.perm, Sort: sort, Visible: 1, Status: 1}
		if _, err := upsert(menu); err != nil {
			return err
		}
		for label, perm := range lf.buttons {
			btn := &model.Menu{ParentID: menu.ID, Name: lf.name + "-" + label, Type: "F", Perm: perm, Status: 1}
			if _, err := upsert(btn); err != nil {
				return err
			}
		}
		sort++
	}

	online := &model.Menu{ParentID: 0, Name: "在线用户", Type: "C", Path: "/monitor/online", Component: "monitor/online/index", Perm: "monitor:online:list", Sort: 3, Visible: 1, Status: 1}
	if _, err := upsert(online); err != nil {
		return err
	}
	kick := &model.Menu{ParentID: online.ID, Name: "在线用户-强退", Type: "F", Perm: "monitor:online:kick", Status: 1}
	if _, err := upsert(kick); err != nil {
		return err
	}

	opLog := &model.Menu{ParentID: 0, Name: "操作日志", Type: "C", Path: "/monitor/log", Component: "monitor/log/index", Perm: "monitor:log:list", Sort: 4, Visible: 1, Status: 1}
	if _, err := upsert(opLog); err != nil {
		return err
	}
	clearLog := &model.Menu{ParentID: opLog.ID, Name: "操作日志-清空", Type: "F", Perm: "monitor:log:delete", Status: 1}
	if _, err := upsert(clearLog); err != nil {
		return err
	}

	server := &model.Menu{ParentID: 0, Name: "服务监控", Type: "C", Path: "/monitor/server", Component: "monitor/server/index", Perm: "monitor:server:list", Sort: 5, Visible: 1, Status: 1}
	if _, err := upsert(server); err != nil {
		return err
	}

	// 全部菜单挂到 super_admin
	var all []model.Menu
	if err := db.Find(&all).Error; err != nil {
		return err
	}
	var superRole model.Role
	if err := db.Where("code = ?", "super_admin").First(&superRole).Error; err != nil {
		return err
	}
	return db.Model(&superRole).Association("Menus").Replace(&all)
}
