package gotool

import ( 
	"reflect"
)

//Field 字段名与类型
type Field struct {
	Name string `FieldName:"Name"`
    Type string
}

//GetFields 获取o的可导出字段(包括匿名)
func GetFields(o interface{}) []Field {
	var fields []Field
	t := reflect.TypeOf(o)
	v := reflect.ValueOf(o)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if v.Field(i).CanInterface() {
			if field.Anonymous { //是否为匿名字段
				fields = append(fields, GetFields(v.Field(i).Interface())...)
				continue
			}
			name := field.Tag.Get("FieldName") //Tag指定Name
			if name == "" {
				name = field.Name
			}
			fields = append(fields, Field{Name: name, Type: field.Type.String()})
		}
	}
	return fields
}

//GetTag 获取结构体o的FieldName的Tag值
func GetTag(o interface{}, FieldName, Tagkey string) (string, bool) {
	field, ok := reflect.TypeOf(o).FieldByName(FieldName)
	return field.Tag.Get(Tagkey), ok
}
