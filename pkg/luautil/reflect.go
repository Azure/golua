package luautil

import (
	"reflect"

	"github.com/Azure/golua/lua"
)

func ValueOf(state *lua.State, any interface{}) lua.Value {
	if any == nil {
		return nil
	}
	return valueOf(state, any)
}

func valueOf(state *lua.State, value interface{}) lua.Value {
	if value, ok := value.(lua.Value); ok {
		return value
	}
	switch rv := reflect.ValueOf(value); rv.Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return lua.Float(rv.Uint())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return lua.Int(rv.Int())
	case reflect.Float32, reflect.Float64:
		return lua.Float(rv.Float())
	case reflect.String:
		return lua.String(rv.String())
	case reflect.Slice:
		return valueFromSlice(state, rv)
	case reflect.Bool:
		return lua.Bool(rv.Bool())
	case reflect.Map:
		return valueFromMap(state, rv)
	case reflect.Ptr:
		if rv.Elem().Kind() == reflect.Struct {
			return valueFromStruct(state, rv)
		}
		fallthrough
	default:
		return lua.UserData(rv.Interface())
	}
}

func check(state *lua.State, index int) reflect.Value {
	return reflect.ValueOf(state.ToUserData(index).Value())
}

func toGoValue(v lua.Value) reflect.Value {
	switch v := v.(type) {
	case *lua.Object:
		return reflect.ValueOf(v.Value())
	case lua.Table:
		return tableToGo(v)
	case lua.String:
		return reflect.ValueOf(string(v))
	case lua.Float:
		return reflect.ValueOf(float64(v))
	case lua.Int:
		return reflect.ValueOf(int64(v))
	case lua.Bool:
		return reflect.ValueOf(bool(v))
	}
	return reflect.ValueOf(nil)
}

func tableToGo(table lua.Table) reflect.Value {
	if length := table.Length(); length == 0 { // map
		gomap := make(map[interface{}]interface{})
		table.ForEach(func(key, val lua.Value) {
			k := toGoValue(key).Interface()
			v := toGoValue(val).Interface()
			gomap[k] = v
		})
		return reflect.ValueOf(gomap)
	} else { // slice
		slice := make([]interface{}, 0, length)
		for i := 1; i <= length; i++ {
			elem := toGoValue(table.Index(lua.Int(i)))
			slice = append(slice, elem.Interface())
		}
		return reflect.ValueOf(slice)
	}
}
