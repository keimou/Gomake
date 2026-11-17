package system

import "gmkdb/services"

type ApiGroup struct {
	// DBApi
	// JwtApi
	BaseApi
	// SystemApi
	// CasbinApi
	// AutoCodeApi
	// SystemApiApi
	// AuthorityApi
	// DictionaryApi
	// AuthorityMenuApi
	// OperationRecordApi
	// AutoCodeHistoryApi
	// DictionaryDetailApi
	// AuthorityBtnApi
}

var (
	// 	apiService              = service.ServiceGroupApp.SystemSvcGroup.ApiService
	jwtService = services.ServiceGroupApp.SystemSvcGroup.JwtService
	// 	menuService             = service.ServiceGroupApp.SystemSvcGroup.MenuService
	userService = services.ServiceGroupApp.SystemSvcGroup.UserService

// initDBService           = service.ServiceGroupApp.SystemSvcGroup.InitDBService
// casbinService           = service.ServiceGroupApp.SystemSvcGroup.CasbinService
// autoCodeService         = service.ServiceGroupApp.SystemSvcGroup.AutoCodeService
// baseMenuService         = service.ServiceGroupApp.SystemSvcGroup.BaseMenuService
// authorityService        = service.ServiceGroupApp.SystemSvcGroup.AuthorityService
// dictionaryService       = service.ServiceGroupApp.SystemSvcGroup.DictionaryService
// systemConfigService     = service.ServiceGroupApp.SystemSvcGroup.SystemConfigService
// operationRecordService  = service.ServiceGroupApp.SystemSvcGroup.OperationRecordService
// autoCodeHistoryService  = service.ServiceGroupApp.SystemSvcGroup.AutoCodeHistoryService
// dictionaryDetailService = service.ServiceGroupApp.SystemSvcGroup.DictionaryDetailService
// authorityBtnService     = service.ServiceGroupApp.SystemSvcGroup.AuthorityBtnService
)
