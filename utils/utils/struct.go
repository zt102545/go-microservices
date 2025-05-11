package utils

import (
	"database/sql"
	"reflect"
	"time"
)

func CopyFields(src, dst interface{}) error {
	srcValue := reflect.ValueOf(src).Elem()
	dstValue := reflect.ValueOf(dst).Elem()

	for i := 0; i < srcValue.NumField(); i++ {
		srcField := srcValue.Field(i)
		dstField := dstValue.FieldByName(srcValue.Type().Field(i).Name)

		if dstField.IsValid() && dstField.CanSet() {
			if srcField.Type() == dstField.Type() {
				dstField.Set(srcField)
			} else if dstField.Type() == reflect.TypeOf(sql.NullInt64{}) && srcField.Type() == reflect.TypeOf(int64(0)) {
				dstField.Set(reflect.ValueOf(sql.NullInt64{Int64: srcField.Int(), Valid: true}))
			} else if dstField.Type() == reflect.TypeOf(sql.NullString{}) && srcField.Type() == reflect.TypeOf("") {
				dstField.Set(reflect.ValueOf(sql.NullString{String: srcField.String(), Valid: true}))
			} else if dstField.Type() == reflect.TypeOf(sql.NullTime{}) && srcField.Type() == reflect.TypeOf("") {
				// 解析字符串为 time.Time 类型
				timeValue, err := time.Parse(time.DateTime, srcField.String())
				if err != nil {
					continue
				}
				dstField.Set(reflect.ValueOf(sql.NullTime{Time: timeValue, Valid: true}))
			} else if dstField.Type() == reflect.TypeOf(time.Time{}) && srcField.Type() == reflect.TypeOf("") {
				// 解析字符串为 time.Time 类型
				timeValue, err := time.Parse(time.DateTime, srcField.String())
				if err != nil {
					continue
				}
				dstField.Set(reflect.ValueOf(timeValue))
			}
		}
	}

	return nil
}

func CopyFieldsBack(src, dst interface{}) error {
	srcValue := reflect.ValueOf(src).Elem()
	dstValue := reflect.ValueOf(dst).Elem()

	for i := 0; i < srcValue.NumField(); i++ {
		srcField := srcValue.Field(i)
		dstField := dstValue.FieldByName(srcValue.Type().Field(i).Name)

		if dstField.IsValid() && dstField.CanSet() {
			if srcField.Type() == dstField.Type() {
				dstField.Set(srcField)
			} else if srcField.Type() == reflect.TypeOf(sql.NullInt64{}) && dstField.Type() == reflect.TypeOf(int64(0)) {
				dstField.Set(reflect.ValueOf(srcField.Interface().(sql.NullInt64).Int64))
			} else if srcField.Type() == reflect.TypeOf(sql.NullString{}) && dstField.Type() == reflect.TypeOf("") {
				dstField.Set(reflect.ValueOf(srcField.Interface().(sql.NullString).String))
			} else if srcField.Type() == reflect.TypeOf(sql.NullTime{}) && dstField.Type() == reflect.TypeOf("") {
				dstField.Set(reflect.ValueOf(srcField.Interface().(sql.NullTime).Time.Format(time.DateTime)))
			} else if srcField.Type() == reflect.TypeOf(time.Time{}) && dstField.Type() == reflect.TypeOf("") {
				dstField.Set(reflect.ValueOf(srcField.Interface().(time.Time).Format(time.DateTime)))
			}
		}
	}

	return nil
}
