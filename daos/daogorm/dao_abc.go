package daogorm

import (
	"fmt"
	"gmkdb/utils"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MysqlDao interface {
	InsBatches(datas any, uptables ...string) error // 批量新增
	Insert(data any) error                          // 新增
	Select(data any, args ...string) error          // 查询
	Selects(datas any, args ...string) error        // 查询多个数据 (未完成)
	Update(data any, args ...map[string]any) error  // 更新
	Delete(datas any, args ...string) error         // 删除
	Restore(data any) error                         // 恢复 还原数据
}

// @Author: [keid]
// @Function: InsBatches
// @Description: Insert data in batches (uptables 指定若冲突则需更新的db字段)
// @Date: 2025-03-02
func (m *gormDB) InsBatches(datas any, uptables ...string) error {
	if err := checkValues(datas); err != nil {
		return err
	}
	// Start the transaction
	tx := m.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Conflict setting
	var onConflict clause.OnConflict
	if len(uptables) > 0 && uptables[0] != "*" { // When the primary key is violated, the non-primary key field is updated
		onConflict.DoUpdates = clause.AssignmentColumns(uptables) // Update only the fields that need to be updated
	} else {
		onConflict.UpdateAll = true // Automatically update all non-primary key fields
	}
	// Insert data in batches (500 entries per batch) CreateInBatches
	if err := tx.Clauses(onConflict).Create(datas).Error; err != nil {
		tx.Rollback()
		return err
	}
	// Commit the transaction
	return tx.Commit().Error
}

// @Author: [keid]
// @Function: Insert
// @Description: 单个插入数据, 冲突时返回错误
// @Date: 2025-03-02
func (m *gormDB) Insert(data any) error {
	if err := checkValue(data); err != nil {
		return err
	}

	// reflectValue := reflect.Indirect(reflect.ValueOf(data))
	// switch reflectValue.Kind() {
	// }
	// var onConflict clause.OnConflict
	// if len(uptables) > 0 {
	// 	if uptables[0] == "*" {
	// 		onConflict.UpdateAll = true
	// 	} else {
	// 		onConflict.DoUpdates = clause.AssignmentColumns(uptables)
	// 	}
	// }

	return m.db.Create(data).Error
}

// @Author: [keid]
// @Function: Selects
// @Description: 查询数据,关联结构字段查询
// @Date: 2025-03-02
func (m *gormDB) Select(data any, args ...string) error {
	if err := checkValue(data); err != nil {
		return err
	}
	// map[string]any 不会忽略零值更新

	// 遍历结构体内非零值进行查询
	sql := m.db.Where(data)
	// Where 若包含主键和其它字段, 会全部关联查询, 若主键和其它字段信息和数据匹配不一致, 则会导致查询不到任何值

	// // 获取表内所有外键名
	// stmt := &gorm.Statement{DB: m.db}
	// if err := stmt.Parse(data); err != nil {
	// 	return err
	// }
	// for _, rel := range stmt.Schema.Relationships.Relations {
	// 	fmt.Printf("rel.Name: %v\n", rel.Name) // 外键名
	// }
	// // stmt.Schema.PrimaryFieldDBNames// 主键字段
	// // stmt.Schema.PrioritizedPrimaryField.Name // 主键名

	// 预加载关联 FieldKey
	for _, fk := range args {
		if utils.HasField(data, fk) {
			sql = sql.Preload(fk)
		} else if fk == "unscoped" {
			sql = sql.Unscoped()
		} else {
			fmt.Println("Not FieldKey", fk)
		}
	}
	return sql.First(data).Error
}

// @Author: [keid]
// @Function: Selects
// @Description: 查询多个数据 (未完成)
// @Date: 2025-03-02
func (m *gormDB) Selects(datas any, args ...string) error {
	if err := checkValues(datas); err != nil {
		return err
	}

	querysql := datas
	for _, fk := range args {
		querysql = fk
	}
	return m.db.Where(querysql).Find(datas).Error
	// return db.Model(datas).Where(querysql).FindInBatches(datas, 100, func(tx *gorm.DB, batch int) error { return nil }).Error
}

// @Author: [keid]
// @Function: Update
// @Description: 更新数据
// @Date: 2025-03-02
func (m *gormDB) Update(data any, mapkv ...map[string]any) error {
	if err := checkValue(data); err != nil {
		return err
	}
	if len(mapkv) > 0 {
		// 显式更新map内所有字段
		return m.db.Model(data).Updates(mapkv[0]).Error
		// Model(&User{})指定模型; Model(&user) 主键不为空, 则按主键查找或修改
	} else {
		// Updates()已保护模型和主键,不需要Model()重复声明; Updates() 会忽略data内的零值更新
		return m.db.Updates(data).Error
		// 或使用 Select() 可显式更新零值字段
		// db.Model(data).Select("ID", "Name", "Time").Updates(data) // Select("*") 更新所有字段
	}
}

// 更新多个数据 (未完成)
func (m *gormDB) Updates(data any, args ...map[string]any) error {
	if err := checkValues(data); err != nil {
		return err
	}
	if len(args) > 0 {
		// 显式更新map内所有字段
		return m.db.Model(data).Updates(args[0]).Error
		// Model(&User{})指定模型; Model(&user) 主键不为空, 则按主键查找或修改
	} else {
		// Updates()已保护模型和主键,不需要Model()重复声明; Updates() 会忽略data内的零值更新
		return m.db.Updates(data).Error
		// 或使用 Select() 可显式更新零值字段
		// db.Model(data).Select("ID", "Name", "Time").Updates(data) // Select("*") 更新所有字段
	}
}

// 删, args: unscoped 强制删除
func (m *gormDB) Delete(datas any, args ...string) (err error) {
	if err := checkValues(datas); err != nil {
		return err
	}
	sql := m.db
	for _, fk := range args {
		if fk == "unscoped" {
			sql = sql.Unscoped()
		}
	}

	// Where是批量删除符合data里面非零数据的条件, 若定位参数 非主键或非不唯一的数据请小心处理
	// m.db.Where(data).Delete(data).Error // 没有Where()的data必须包含非零主键, 否则无法定位
	return sql.Delete(datas).Error
}

// 还原数据 restore
func (m *gormDB) Restore(data any) (err error) {
	if err := checkValue(data); err != nil {
		return err
	}

	// v := reflect.ValueOf(data)
	// if v.Kind() == reflect.Ptr {
	// 	v2 := v.Elem()
	// 	fmt.Printf("v2: %v\n", v2)
	// 	if v2.Kind() == reflect.Slice || v2.Kind() == reflect.Array {
	// 		// v.Type().Elem().Elem().Kind() != reflect.Struct
	// 		v3 := v.Type().Elem().Elem()
	// 		fmt.Printf("v3: %v\n", v3)
	// 		if v3.Kind() == reflect.Struct {
	// 			fmt.Printf("v3.Type(): %v\n", v3)
	// 		}
	// 	} else {
	// 		fmt.Printf("v2: %v\n", v2)
	// 	}
	// }

	// 校验数据是否存在
	// err = m.db.Unscoped().Where(data).Find(data).Error
	err = m.Select(data, "unscoped")
	if err != nil {
		return
	}
	// 还原数据
	err = m.db.Unscoped().Model(data).Update("deleted_at", gorm.DeletedAt{}).Error
	return
}
