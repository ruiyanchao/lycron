package lycron

import "reflect"

func IsNil(i interface{}) bool {
	if i ==nil{
		return  true
	}
	vi := reflect.ValueOf(i)//通过反射获取其对应的值
	if vi.Kind() == reflect.Ptr || vi.Kind() == reflect.Slice || vi.Kind()==reflect.Array  {
		return vi.IsNil()
	}
	return false
}
