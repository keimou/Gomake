package daosql

import (
	"errors"
	"fmt"
	"gmkdb/utils"
	"log"
	"reflect"
	"strings"
	"time"

	"go.uber.org/zap"
)

// INSERT INTO
// [table] USING [stable]
// (s, i) TAGS ('[s]', '[i]')
// (fv, ...) VALUES ([v], ... ) field Column

//	data: *struct{}
//	args: (string) '数据库名.'
//
// 新增单条数据:
func (s *SqlDao) Insert(data TABBASE, args ...any) (err error) {
	var dbName string
	for _, arg := range args {
		switch v := arg.(type) {
		case string: // table name: [dbName].[symbol_interval]
			slices := strings.Split(v, ".")
			dbName = slices[0]
			// fmt.Printf("Insert - dbName: %v\n", dbName)
		default:
		}
	}
	// if len(_dbName) > 0 && strings.TrimSpace(_dbName[0]) != "" {
	// 	dbName = _dbName[0]
	// }

	var tagSql, valSql sqlKeyValue
	columns := fieldColumns(data)
	for _, v := range columns {
		// rv := reflect.ValueOf(v.Value)
		// log.Printf("rv: %v  ", rv.Interface()) // 地址	addr
		// log.Print(rv.Elem().Interface())       // 值		value
		// log.Printf("v.Field: %v\n", v.Field)   // 字段名	field
		// value := reflect.ValueOf(v.Value).Elem().Interface()
		var value any
		val := reflect.ValueOf(v.Value)
		if val.Kind() == reflect.Ptr {
			value = val.Elem().Interface()
		} else {
			value = val.Interface()
		}

		// switch v.Kind {
		switch vt := value.(type) {
		// case reflect.String:
		case string:
			// fmt.Printf("s: %v\n", v)
			if utils.IsNull(value) && v.Note == "NULL" {
				// log.Println(" * str is nil(skip):", v.Field, v.Kind, v.Type)
				continue
			}
			value = fmt.Sprintf("'%s'", value)
		// case reflect.Int, reflect.Int64, reflect.Float64:
		case int, int64, float64:
			// fmt.Printf("i: %v\n", v)
			if utils.IsZero(value) && v.Note == "NULL" {
				// log.Println(" * num is 0(skip):", v.Field, v.Kind, v.Type)
				continue
			}
		case time.Time:
			// fmt.Printf("v: %v %v\n", v, v.UnixMilli())
			value = vt.UnixMilli()
		default:
			// fmt.Printf("d: %v\n", v)
		}
		if v.Note == "TAG" {
			// if !strings.Contains(v.Type, "INT") {
			// 	value = fmt.Sprintf("'%s'", value)
			// }
			utils.StrJoin(&tagSql.key, v.Field)
			utils.StrJoin(&tagSql.value, value)
		} else {
			utils.StrJoin(&valSql.key, v.Field)
			utils.StrJoin(&valSql.value, value)
			// if strings.Contains(v.Type, "INT") && utils.IsZero(value) {
			// 	log.Printf(" - %v(%v)int: %v\n", v.Field, v.Type, value)
			// } else if !strings.Contains(v.Type, "INT") && utils.IsNull(value) {
			// 	log.Printf(" - %v(%v)str: %v\n", v.Field, v.Type, value)
			// 	// } else if strings.Contains(v.Type, "CHAR") || strings.Contains(v.Type, "ARY") {
			// 	// 	// log.Printf(" * %v(%v)s: %v\n", v.Field, v.Type, value)
			// 	// 	utils.StrJoin(&valSql.key, v.Field)
			// 	// 	utils.StrJoin(&valSql.value, fmt.Sprintf("'%s'", value))
			// } else {
			// 	utils.StrJoin(&valSql.key, v.Field)
			// 	utils.StrJoin(&valSql.value, value)
			// }
		}
	}

	if tagSql.key != "" && valSql.key != "" {
		// INSERT INTO %v USING %v (%s) TAGS (%s) (%s) VALUES (%s)
		insertQuery := fmt.Sprintf("INSERT INTO %v USING %v (%s) TAGS (%s) (%s) VALUES(%s)",
			s.dbSTable(data.SubTableName(), dbName), s.dbSTable(data.TableName(), dbName),
			tagSql.key, tagSql.value, valSql.key, valSql.value)

		// fmt.Printf("\nexec sql:\n %v\n\n", insertQuery)
		_, err = s.taos.Exec(insertQuery)
	} else {
		err = errors.New("Insert input data is incorrect")
	}
	return
}

//	datas: []*struct{}
//	args: (string) '数据库名.'
//
// 批量新增数据,
func (s *SqlDao) Inserts(datas any, args ...any) (err error) {
	var dbName string
	for _, arg := range args {
		switch v := arg.(type) {
		case string: // table name: [dbName].[symbol_interval]
			slices := strings.Split(v, ".")
			dbName = slices[0]
		default:
		}
	}

	elements, err := processSlice(datas)
	if err != nil {
		return
	}

	// var insertQuery string
	insertQuery := "INSERT INTO"
	for _, elem := range elements {
		var tagSql, valSql sqlKeyValue
		columns := fieldColumns(elem)
		for _, v := range columns {
			val := reflect.ValueOf(v.Value)
			if val.Kind() == reflect.Ptr {
				val = val.Elem()
			}
			value := val.Interface()
			// switch v.Kind {
			// switch val.Kind() {
			switch vt := value.(type) {
			// case reflect.String:
			case string:
				if utils.IsNull(value) && v.Note == "NULL" {
					continue
				}
				value = fmt.Sprintf("'%s'", value)
			// case reflect.Int, reflect.Int64, reflect.Float64:
			case int, int64, float64:
				if utils.IsZero(value) && v.Note == "NULL" {
					continue
				}
			case time.Time:
				// fmt.Printf("v: %v %v\n", v, v.UnixMilli())
				value = vt.UnixMilli()
			default:
				// fmt.Printf("d: %v\n", v)
			}
			if v.Note == "TAG" {
				utils.StrJoin(&tagSql.key, v.Field)
				utils.StrJoin(&tagSql.value, value)
			} else {
				utils.StrJoin(&valSql.key, v.Field)
				utils.StrJoin(&valSql.value, value)
			}
		}

		insertQuery += fmt.Sprintf(" %v USING %v (%s) TAGS (%s) (%s) VALUES(%s)",
			s.dbSTable(elem.SubTableName(), dbName), s.dbSTable(elem.TableName(), dbName),
			tagSql.key, tagSql.value, valSql.key, valSql.value)

		// INSERT INTO d1001 USING meters TAGS('beijing', 2) (ts, current, voltage, phase) VALUES (now, 10.2, 219, 0.32);
		// INSERT INTO %v USING %v (%s) TAGS (%s) (%s) VALUES (%s_v1), (%s_v2), (%s_v3);
		// 一个表对应一个VALUES, 若要每组数据都标识对应VALUES的字段, 只有整表往下接
	}

	insertQuery += ";" // 执行句尾, 可加可不加
	// fmt.Printf("\n\n %v \n\n", insertQuery)
	_, err = s.taos.Exec(insertQuery)
	if err != nil {
		err = errors.Join(err, errors.New(insertQuery))
	}

	return
}

// 检查输入数据, 将其转换
func processSlice(input any) (output []TABBASE, err error) {
	// 使用反射检查输入类型是否为 []*struct
	inputVal := reflect.ValueOf(input)
	// fmt.Printf("inputVal.Kind(): %v\n", inputVal.Kind())
	// fmt.Printf("inputVal.IsValid(): %v\n", inputVal.IsValid())
	// fmt.Printf("inputVal.IsZero(): %v\n", inputVal.IsZero())

	err = errors.New("SqlDB Inserts input data is incorrect")
	if valType := inputVal.Kind(); valType == reflect.Slice { // || valType == reflect.Array
		// fmt.Printf("inputVal.Len(): %v\n", inputVal.Len())
		// if valType == reflect.Array && inputVal.IsZero() { // 0判断有效值数量
		// 	return
		// } else if valType == reflect.Slice && (inputVal.IsNil() || !(inputVal.Len() > 0)) { // nil判断初始化, len判断有效值数量
		// 	return
		// }
		valType = inputVal.Type().Elem().Kind()
		if valType == reflect.Ptr {
			valType = inputVal.Type().Elem().Elem().Kind()
		}
		if valType == reflect.Struct {
			err = nil
		}
	}

	if err != nil { // inputVal.Kind() != reflect.Slice || inputVal.Type().Elem().Kind() != reflect.Ptr || inputVal.Type().Elem().Elem().Kind() != reflect.Struct {
		// err = errors.New("input data is incorrect") // "input is not a slice of pointers to structs"
		return
	}

	// output := make([]TABBASE, inputVal.Len())
	for i := 0; i < inputVal.Len(); i++ {
		element := inputVal.Index(i)
		if !element.IsValid() || element.IsZero() { // element.IsNil()
			// global.GIA_LOG.Debug("SqlDB Inserts processSlice invalid element at index", zap.Int("index", i), zap.Any("elem", element))
			log.Println("SqlDB Inserts processSlice invalid element at index", zap.Int("index", i), zap.Any("elem", element))
			// return
			continue
		}
		// // 获取指针指向的实际值(非必要)
		// elem := element.Elem()
		// fmt.Printf("element: %v\n", element)

		// 通过反射获取结构体指针指向的结构体
		// 使用类型断言检查是否实现了TABBASE接口
		if tabbase, ok := element.Interface().(TABBASE); ok {
			output = append(output, tabbase)
		} else {
			output = nil
			fmt.Printf("SqlDB Inserts input data tabbase: %v\n", tabbase)
			err = errors.New("SqlDB Inserts input data struct is incorrect:" + element.String())
			break
		}
	}

	if !(len(output) > 0) && err == nil {
		err = errors.New("SqlDB Inserts no valid data entry")
	}

	return
}
