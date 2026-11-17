package system

import (
	"time"
)

type SysAuthority struct {
	AuthorityId   uint   `json:"authorityId" gorm:"not null;unique;primary_key;size:90"` // 角色ID
	AuthorityName string `json:"authorityName" gorm:""`                                  // 角色名
	ParentId      uint   `json:"parentId" gorm:""`                                       // 父角色ID

	Children        []SysAuthority  `json:"children" gorm:"-"`
	DataAuthorityId []*SysAuthority `json:"dataAuthorityId" gorm:"many2many:sys_data_authority_id;"`
	SysBaseMenus    []SysBaseMenu   `json:"menus" gorm:"many2many:sys_authority_menus;"`
	Users           []SysUser       `json:"-" gorm:"many2many:sys_user_authority;"`
	DefaultRouter   string          `json:"defaultRouter" gorm:"comment:默认菜单;default:dashboard"` // 默认菜单(默认dashboard)

	CreatedAt time.Time  // 创建时间
	UpdatedAt time.Time  // 更新时间
	DeletedAt *time.Time `sql:"index"`
}

func (SysAuthority) TableName() string {
	return "sys_authorities"
}
