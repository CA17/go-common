package excel

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"

	"github.com/ca17/go-common/sqltype"
)

func WriteToFile(sheet string, records []interface{}, filepath string) error {
	xlsx := excelize.NewFile()
	index := xlsx.NewSheet(sheet)

	for i, t := range records {
		WriteRow(t, i, xlsx, sheet)
	}

	xlsx.SetActiveSheet(index)
	return xlsx.SaveAs(filepath)

}

func WriteRow(t interface{}, i int, xlsx *excelize.File, sheet string) {
	d := reflect.TypeOf(t).Elem()
	for j := 0; j < d.NumField(); j++ {
		// 设置表头
		if i == 0 {
			column := strings.Split(d.Field(j).Tag.Get("xlsx"), "-")[0]
			name := strings.Split(d.Field(j).Tag.Get("xlsx"), "-")[1]
			xlsx.SetCellValue(sheet, fmt.Sprintf("%s%d", column, i+1), name)
		}
		// 设置内容
		column := strings.Split(d.Field(j).Tag.Get("xlsx"), "-")[0]
		ctype := d.Field(j).Type.String()
		switch ctype {
		case "string":
			xlsx.SetCellValue(sheet, fmt.Sprintf("%s%d", column, i+2), reflect.ValueOf(t).Elem().Field(j).String())
		case "int":
			xlsx.SetCellValue(sheet, fmt.Sprintf("%s%d", column, i+2), reflect.ValueOf(t).Elem().Field(j).Int())
		case "int32":
			xlsx.SetCellValue(sheet, fmt.Sprintf("%s%d", column, i+2), reflect.ValueOf(t).Elem().Field(j).Int())
		case "int64":
			xlsx.SetCellValue(sheet, fmt.Sprintf("%s%d", column, i+2), reflect.ValueOf(t).Elem().Field(j).Int())
		case "bool":
			xlsx.SetCellValue(sheet, fmt.Sprintf("%s%d", column, i+2), reflect.ValueOf(t).Elem().Field(j).Bool())
		case "float32":
			xlsx.SetCellValue(sheet, fmt.Sprintf("%s%d", column, i+2), reflect.ValueOf(t).Elem().Field(j).Float())
		case "float64":
			xlsx.SetCellValue(sheet, fmt.Sprintf("%s%d", column, i+2), reflect.ValueOf(t).Elem().Field(j).Float())
		case "time.Time":
			xlsx.SetCellValue(sheet, fmt.Sprintf("%s%d", column, i+2), reflect.ValueOf(t).Elem().Field(j).Interface().(time.Time).Format("2006-01-02 15:04:05"))
		case "sqltype.JsonNullTime":
			f := reflect.ValueOf(t).Elem().Field(j).Interface().(sqltype.JsonNullTime)
			xlsx.SetCellValue(sheet, fmt.Sprintf("%s%d", column, i+2), f.String())
		default:
			xlsx.SetCellValue(sheet, fmt.Sprintf("%s%d", column, i+2), reflect.ValueOf(t).Elem().Field(j).String())

		}
	}
}
