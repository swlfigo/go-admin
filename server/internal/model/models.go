package model

import (
	"time"

	"gorm.io/gorm"
)

type Base struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type User struct {
	Base
	Username       string     `gorm:"uniqueIndex:idx_user_username,where:deleted_at IS NULL;size:64;not null" json:"username"`
	Password       string     `gorm:"size:100;not null" json:"-"`
	Nickname       string     `gorm:"size:64" json:"nickname"`
	Avatar         string     `gorm:"size:255" json:"avatar"`
	Email          string     `gorm:"size:128" json:"email"`
	Phone          string     `gorm:"size:32" json:"phone"`
	Status         int        `gorm:"default:1" json:"status"` // 0停用 1启用
	LastLoginAt    *time.Time `json:"lastLoginAt"`
	LoginFailCount int        `gorm:"default:0" json:"-"`
	LockUntil      *time.Time `json:"-"`
	Roles          []Role     `gorm:"many2many:sys_user_role;" json:"roles"`
}

func (User) TableName() string { return "sys_user" }

type Role struct {
	Base
	Name   string `gorm:"size:64;not null" json:"name"`
	Code   string `gorm:"uniqueIndex:idx_role_code,where:deleted_at IS NULL;size:64;not null" json:"code"`
	Sort   int    `gorm:"default:0" json:"sort"`
	Status int    `gorm:"default:1" json:"status"`
	Remark string `gorm:"size:255" json:"remark"`
	Menus  []Menu `gorm:"many2many:sys_role_menu;" json:"menus"`
}

func (Role) TableName() string { return "sys_role" }

type Menu struct {
	Base
	ParentID  uint   `gorm:"default:0;index" json:"parentId"`
	Name      string `gorm:"size:64;not null" json:"name"`
	Type      string `gorm:"size:1;not null" json:"type"` // M目录 C菜单 F按钮
	Path      string `gorm:"size:128" json:"path"`
	Component string `gorm:"size:128" json:"component"`
	Perm      string `gorm:"size:128" json:"perm"`
	Icon      string `gorm:"size:64" json:"icon"`
	Sort      int    `gorm:"default:0" json:"sort"`
	Visible   int    `gorm:"default:1" json:"visible"`
	Status    int    `gorm:"default:1" json:"status"`
}

func (Menu) TableName() string { return "sys_menu" }

type DictType struct {
	Base
	Name   string `gorm:"size:64;not null" json:"name"`
	Type   string `gorm:"uniqueIndex:idx_dict_type,where:deleted_at IS NULL;size:64;not null" json:"type"`
	Status int    `gorm:"default:1" json:"status"`
	Remark string `gorm:"size:255" json:"remark"`
}

func (DictType) TableName() string { return "sys_dict_type" }

type DictData struct {
	Base
	DictType string `gorm:"size:64;not null;index" json:"dictType"`
	Label    string `gorm:"size:64;not null" json:"label"`
	Value    string `gorm:"size:64;not null" json:"value"`
	TagType  string `gorm:"size:32" json:"tagType"`
	Sort     int    `gorm:"default:0" json:"sort"`
	Status   int    `gorm:"default:1" json:"status"`
	Remark   string `gorm:"size:255" json:"remark"`
}

func (DictData) TableName() string { return "sys_dict_data" }

// OnlineSession = 在线用户列表 + refresh 轮换载体
type OnlineSession struct {
	ID              string    `gorm:"primaryKey;size:64" json:"id"` // access jti
	UserID          uint      `gorm:"index" json:"userId"`
	Username        string    `gorm:"size:64" json:"username"`
	LoginIP         string    `gorm:"size:64" json:"loginIp"`
	Browser         string    `gorm:"size:64" json:"browser"`
	OS              string    `gorm:"size:64" json:"os"`
	RefreshID       string    `gorm:"size:64;index" json:"-"`
	RefreshExpireAt time.Time `json:"-"`
	LoginAt         time.Time `json:"loginAt"`
	LastActiveAt    time.Time `json:"lastActiveAt"`
}

func (OnlineSession) TableName() string { return "sys_online_session" }

// OperationLog 操作日志（无软删，纯插入/清空）。
type OperationLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"size:64;index" json:"username"`
	Method    string    `gorm:"size:8" json:"method"`
	Path      string    `gorm:"size:255" json:"path"`
	IP        string    `gorm:"size:64" json:"ip"`
	Status    int       `json:"status"`
	LatencyMs int       `json:"latencyMs"`
	CreatedAt time.Time `gorm:"index" json:"createdAt"`
}

func (OperationLog) TableName() string { return "sys_operation_log" }
