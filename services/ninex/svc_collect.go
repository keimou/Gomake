package ninex

import (
	"fmt"
	"gmkdb/global"
	"log"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type CollectService struct{}

// @Author: [keid]
// @Function: CollectExpert
// @Description: 收集高手信息 新增用户条例
// @Date: 2026-10-28
func (*CollectService) CollectExpert(uid uint) {
	// var uid uint
	user, datas := GetExpertData(uid)
	fmt.Printf("user: %v\n", user)
	fmt.Printf("datas: %v %v\n", len(datas), datas[len(datas)-1])

	err := global.GIA_DB.Create(&user).Error
	if err != nil {
		global.GIA_LOG.Error("新增用户失败!", zap.Error(err))
	}

	if len(datas) > 0 {
		// result := global.GIA_DB.CreateInBatches(&datas, 100)
		err := global.GIA_DB.Transaction(func(tx *gorm.DB) error {
			if err := tx.CreateInBatches(&datas, 200).Error; err != nil {
				return err // 自动回滚
			}
			return nil // 自动提交
		})
		if err != nil {
			log.Printf("批量插入失败: %v", err)
		}
	}

}

func (*CollectService) CheckHistory() {

}
