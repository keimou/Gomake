package daosql

import (
	"errors"
	"fmt"
	"gmkdb/utils"
	"log"
	"reflect"
	"strings"

	"go.uber.org/zap"
)

// 查询数据 data: *struct{} args: subT(string); st(int64 <=0->DESC,0:now)
func (s *SqlDao) Select(data TABBASE, args ...any) (err error) {
	// fmt.Printf("s.dbName: %v\n", s.dbName)
	// query := "SELECT"
	// query += fmt.Sprintf(` %s`, "ts,end_time,open_price,close_price,volume,trade_num,is_final,s,i")
	// query += fmt.Sprintf(` FROM %s`, s.dbSTable(dest.SubTableName()))
	// // query += fmt.Sprintf(` WHERE start_time = %d`, 123)
	// query += " ORDER BY ts DESC"
	// query += fmt.Sprintf(` LIMIT %d`, 3)

	var queryField, sqlWhere, dbTable, dbName string
	var openTime int64
	var dest []any

	// var order bool
	// fmt.Printf("args: %v\n", args)
	for _, arg := range args {
		switch v := arg.(type) {
		case string: // table name
			// dbTable = v
			slices := strings.Split(v, ".")
			dbName = slices[0]
			if len(slices) == 2 {
				dbTable = slices[1]
			}
		case int64: // v>0 => ts>=v (正序前排), v=0 取最大最新值(倒序第一), v<0 => ts<=|v| (倒序前排)
			openTime = v
		// case bool:
		// 	order = v
		default:
			// r.Data = v // data
		}
	}
	dbTable = data.SubTableName(dbTable)
	columns := fieldColumns(data)
	for _, v := range columns {
		dest = append(dest, v.Value)
		if v.Note != "NULL" { // 设置 查询字段 和数据存储地址
			utils.StrJoin(&queryField, v.Field)
		} else {
			// CASE WHEN value IS NULL THEN '' ELSE value END
			value := "0"
			if v.Kind == reflect.String {
				value = "''"
			}
			nullField := fmt.Sprintf("CASE WHEN %s IS NULL THEN %s ELSE %s END", v.Field, value, v.Field)
			utils.StrJoin(&queryField, nullField)
		}
		if v.Note != "TAG" { // 获取不为nil的查询参数
			var val string
			switch vtp := v.Value.(type) {
			case *int:
				if *vtp != 0 {
					val = fmt.Sprintf("%d", *vtp)
					// utils.StrJoin(&sqlValue, fmt.Sprintf(" WHERE %s = %d", v.Field, *vtp), " And")
				}
			case *int64:
				if *vtp != 0 {
					val = fmt.Sprintf("%d", *vtp)
					// utils.StrJoin(&sqlValue, fmt.Sprintf(" WHERE %s = %d", v.Field, *vtp), " And")
				}
			case *float64:
				if *vtp != 0 {
					val = fmt.Sprintf("%v", *vtp)
					// utils.StrJoin(&sqlValue, fmt.Sprintf(" WHERE %s = %v", v.Field, *vtp), " And")
				}
			case *string:
				if strings.TrimSpace(*vtp) != "" {
					val = *vtp
					// utils.StrJoin(&sqlValue, fmt.Sprintf(" WHERE %s = %s", v.Field, *vtp), " And")
				}
			case *bool:
				if *vtp {
					val = fmt.Sprintf("%v", *vtp)
					// utils.StrJoin(&sqlValue, fmt.Sprintf(" WHERE %s = %v", v.Field, *vtp), " And")
				}
			default:
			}
			if val != "" {
				utils.StrJoin(&sqlWhere, fmt.Sprintf(" %s = %v", v.Field, val), " And")
			}
		}
	}

	if queryField != "" {
		// INSERT INTO %v USING %v (%s) TAGS (%s) (%s) VALUES (%s)
		// query := fmt.Sprintf("SELECT %v FROM %s ORDER BY ts ASC LIMIT 1",
		// 	queryField, dbTable)
		query := "SELECT"
		query += fmt.Sprintf(" %v FROM %s",
			queryField, s.dbSTable(dbTable, dbName))

		if openTime != 0 {
			tag, ord := ">=", "ASC"
			if openTime < 0 {
				tag, ord = "<=", "DESC"
				openTime = -openTime
			}
			utils.StrJoin(&sqlWhere, fmt.Sprintf(" ts %s %v ORDER BY ts %s", tag, openTime, ord), " And")
			query += " WHERE" + sqlWhere
			// query += fmt.Sprintf(" WHERE ts %s %v ORDER BY ts %s", tag, openTime, ord)
		} else {
			if sqlWhere != "" {
				query += " WHERE" + sqlWhere
			}
			query += " ORDER BY ts DESC"
		}
		// query += " LIMIT 1"

		// fmt.Printf("\nQuery sql:\n %v\n\n", query)
		row := s.taos.QueryRow(query)
		err = row.Scan(dest...)
		// if err != nil {
		// 	global.GIA_LOG.Debug("SqlDao select Scan Err", zap.Error(err), zap.Int64(dbTable, openTime))
		// }
		// fmt.Printf("dest: %v\n", dest)
	} else {
		err = errors.ErrUnsupported
		// global.GIA_LOG.Debug("SqlDao select Scan Err", zap.Error(err), zap.Int64(dbName+"."+dbTable, openTime))
		log.Panicln("SqlDao select Scan Err", zap.Error(err), zap.Int64(dbName+"."+dbTable, openTime))
	}
	return
}

// 批量查询数据， datas: &[]*struct{} args: subT(string); st(int64); limit(int>0); order(bool)
func (s *SqlDao) Selects(datas any, args ...any) (err error) {
	ptrVal := reflect.ValueOf(datas)
	if ptrVal.Kind() != reflect.Ptr || ptrVal.Elem().Kind() != reflect.Slice {
		// || ptrVal.Type().Elem().Elem().Kind() != reflect.Ptr  // 且为指针结构体, 已解决该问题
		return errors.New("参数必须为切片指针")
	}

	sliceVal := ptrVal.Elem()          // 获取切片指针
	elemType := sliceVal.Type().Elem() // 获取切片类型, 为添加新实例提供模型

	var elem1st reflect.Value // 扫描数据库表 字段,表名
	if sliceVal.Len() > 0 {   // 如果原切片非空,提取第一个元素(获取子表名,标签等信息)
		elem1st = sliceVal.Index(0) // .Interface().(TABBASE)
		// 清空原切片 -> 结尾赋值已实现(考虑无数据报错情况, 继续保留)
		newEmptySlice := reflect.MakeSlice(sliceVal.Type(), 0, 0)
		sliceVal.Set(newEmptySlice)
	} else { // 否则创建一个新实例
		if elemType.Kind() == reflect.Ptr {
			// structType := elemType.Elem()     // 获取指针指向的实际类型
			elem1st = reflect.New(elemType.Elem()) // 创建结构体实例的指针(等价于 &struct{})
		} else {
			elem1st = reflect.New(elemType).Elem() // 创建结构体实例(等价于 struct{})
		}
	}

	// 读取参数
	var sqlField, sqlWhere, dbTable, dbName string
	var openTime int64
	var order bool
	limit := 1000 // 默认查询数量
	for _, arg := range args {
		switch v := arg.(type) {
		case string: // table name: [dbName].[symbol_interval]
			slices := strings.Split(v, ".")
			if len(slices) == 2 {
				dbName = slices[0]
				dbTable = slices[1]
			}
		case int64: // ts: 1234567890
			openTime = v
		case int:
			limit = v
		case *int:
			limit = *v
		case bool: // 是否倒序
			order = v
		default:
			// r.Data = v // data
		}
	}
	if limit > 3000 { // 限制查询数量
		limit = 3000
	} else if limit < 1 {
		limit = 1
	}
	// 获取相关查询信息
	if elem1st.IsValid() {
		elem := elem1st.Interface().(TABBASE)
		dbTable = elem.SubTableName(dbTable) // 获取sql表结构字段
		columns := fieldColumns(elem)
		for _, v := range columns { // time.Time -> struct
			if v.Note != "NULL" { // taos 读取忽略, db 字段忽略
				utils.StrJoin(&sqlField, v.Field)
			} else {
				// CASE WHEN value IS NULL THEN '' ELSE value END
				value := "0" // 为nil值赋初始值(数值型)
				if v.Kind == reflect.String {
					value = "''"
				}
				nullField := fmt.Sprintf("CASE WHEN %s IS NULL THEN %s ELSE %s END", v.Field, value, v.Field)
				utils.StrJoin(&sqlField, nullField)
			}
			if v.Note != "TAG" { // 获取不为nil的查询参数
				// fmt.Printf("\nv: %v\n", v)
				if !utils.IsNil(v.Value) {
					if value, ok := utils.RemPtr(v.Value); !ok {
						utils.StrJoin(&sqlWhere, fmt.Sprintf(" %s = %v", v.Field, value), " And")
					}
				}
			}
		}
	}

	// // var newElem reflect.Value
	// if elemType.Kind() == reflect.Ptr {
	// 	fmt.Println("获取指针指向的实际类型")
	// 	elemType = elemType.Elem() // 获取指针指向的实际类型
	// 	// newElem = reflect.New(structType) // 创建结构体实例的指针(等价于 &struct{})
	// 	// } else {
	// 	// 	newElem = reflect.New(elemType).Elem() // 创建结构体实例(等价于 struct{})
	// }
	// newElem := reflect.New(elemType).Elem()
	// // newElemPtr := reflect.New(elemType) // 类型为 *struct
	// // newElem := reflect.New(elemType).Elem() // .Interface().(TABBASE) // 创建实例
	// // 追加并更新原切片
	// newSlice := reflect.Append(sliceVal, newElem)
	// sliceVal.Set(newSlice)
	// dest, ok := newElem.Interface().(TABBASE)
	// if !ok {
	// 	return errors.New("参数非有效的数据库表")
	// }
	// fmt.Printf("dest: %v\n", dest.SubTableName())
	// var valSql sqlKeyValue
	// var args []any
	// columns := fieldColumns(dest)
	// // fmt.Printf("columns: %v\n", columns)
	// for _, v := range columns {
	// 	var value any
	// 	val := reflect.ValueOf(v.Value)
	// 	if val.Kind() == reflect.Ptr {
	// 		value = val.Elem().Interface()
	// 		if v.Note != "-" {
	// 			utils.StrJoin(&valSql.key, v.Field)
	// 			utils.StrJoin(&valSql.value, value)
	// 			if v.Type == "int64" {
	// 				fmt.Printf("v.Field: %v\n", v.Field)
	// 			}
	// 			args = append(args, v.Value)
	// 		}
	// 	}
	// }
	// fmt.Printf("args: %v\n", args)

	// fmt.Printf("dbTable: %v\n", dbTable)
	// fmt.Printf("sqlField: %v\n", sqlField)
	if dbTable != "" && sqlField != "" {
		// SELECT * FROM [table] ORDER BY ts DESC LIMIT 1
		query := "SELECT"
		query += fmt.Sprintf(" %v FROM %s", // ORDER BY ts DESC
			sqlField, s.dbSTable(dbTable, dbName))
		if openTime > 0 {
			tag := ">="
			if order {
				tag = "<="
			}
			utils.StrJoin(&sqlWhere, fmt.Sprintf(" ts %s %v", tag, openTime), " And")
			// query += fmt.Sprintf(" WHERE ts >= %v", openTime)
			// WHERE ts > 'A' AND ts <= 'B'	// 从 A 到 B 时间范围
			// AND col3 LIKE '%nny'			// 并且 col3 的字符串是'nny'结尾
			// WHERE ts >= NOW - 10m		// 查询过去 10 分钟的记录
		}
		if sqlWhere != "" {
			query += " WHERE" + sqlWhere
		}
		query += " ORDER BY ts"
		if order {
			query += " DESC"
		}
		query += fmt.Sprintf(" LIMIT %d", limit)

		// fmt.Printf("\n *** %v\n\n", query)
		rows, e := s.taos.Query(query)
		if e != nil {
			// fmt.Printf("Query err: %v\n", e)
			return e
		}
		defer rows.Close()
		// // 获取查询结果的列信息
		// columns, err := rows.Columns()
		// if err != nil {
		// 	return fmt.Errorf("获取列信息失败: %v", err)
		// }
		// fmt.Printf("columns: %v\n", columns)
		// for i, v := range columns {
		// 	fmt.Printf("%v: %v\n", i, v)
		// }
		// // if true {
		// // 	return err
		// // }

		tempSlice := reflect.MakeSlice(sliceVal.Type(), 0, 0)
		for rows.Next() {
			// // 新建切片模型
			var newElem reflect.Value
			if elemType.Kind() == reflect.Ptr {
				newElem = reflect.New(elemType.Elem())
			} else {
				newElem = reflect.New(elemType)
			}

			// 关联切片模型到入参
			var dest []any
			columns := fieldColumns(newElem.Interface().(TABBASE))
			for _, v := range columns {
				dest = append(dest, v.Value)
				// if v.Note != "NULL" {
				// 	dest = append(dest, v.Value)
				// } else {
				// 	// DefaultInt64
				// 	dest = append(dest, v.Value)
				// }
			}

			// fmt.Printf("len(dest): %v\n", len(dest))
			// 扫描行数据
			err = rows.Scan(dest...)
			// err = rows.Scan(&ts, &endTime, &open, &close, &volume, &num, &isFinal)
			if err != nil {
				fmt.Printf(" * Scan err: %v\n", err)
			}

			// 将新元素追加到切片
			if elemType.Kind() == reflect.Ptr {
				// sliceVal.Set(reflect.Append(sliceVal, newElem))
				tempSlice = reflect.Append(tempSlice, newElem)
			} else {
				// sliceVal.Set(reflect.Append(sliceVal, newElem.Elem()))
				tempSlice = reflect.Append(tempSlice, newElem.Elem())
			}
		}
		if true {
			// 正序
			sliceVal.Set(tempSlice)
		} else {
			// 倒序添加：将新元素添加到开头
			newEmptySlice := reflect.MakeSlice(sliceVal.Type(), 0, 0)
			for i := tempSlice.Len() - 1; i >= 0; i-- {
				newEmptySlice = reflect.Append(newEmptySlice, tempSlice.Index(i))
			}
			sliceVal.Set(newEmptySlice)
		}
		return rows.Err()
	} else {
		return errors.ErrUnsupported
	}
}

// 带默认值的字符串类型
type DefaultString string

func (s *DefaultString) Scan(value interface{}) error {
	if value == nil {
		*s = "" // NULL 时赋值为空字符串
		return nil
	}
	switch v := value.(type) {
	case string:
		*s = DefaultString(v)
	case []byte:
		*s = DefaultString(v)
	default:
		return fmt.Errorf("不支持的字符串类型: %T", value)
	}
	return nil
}

// 带默认值的浮点类型
type DefaultFloat64 float64

func (f *DefaultFloat64) Scan(value interface{}) error {
	if value == nil {
		*f = 0.0 // NULL 时赋值为 0
		return nil
	}
	switch v := value.(type) {
	case float64:
		*f = DefaultFloat64(v)
	case float32:
		*f = DefaultFloat64(v)
	case int64:
		*f = DefaultFloat64(v)
	default:
		return fmt.Errorf("不支持的浮点类型: %T", value)
	}
	return nil
}

// 带默认值的浮点类型
type DefaultInt64 int64

func (f *DefaultInt64) Scan(value interface{}) error {
	if value == nil {
		*f = 0.0 // NULL 时赋值为 0
		return nil
	}
	switch v := value.(type) {
	case int:
		*f = DefaultInt64(v)
	case int32:
		*f = DefaultInt64(v)
	case int64:
		*f = DefaultInt64(v)
	default:
		return fmt.Errorf("不支持的浮点类型: %T", value)
	}
	return nil
}
