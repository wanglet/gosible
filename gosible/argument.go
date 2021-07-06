package gosible

import (
    "fmt"
    "reflect"
)

func marshal(i interface{}, tagName string, separator string) []string {
    res := []string{}
    v := reflect.ValueOf(i)
    t := v.Type()
    for i := 0; i < v.NumField(); i++ {
        vField := v.Field(i)
        tag := t.Field(i).Tag.Get(tagName)
        if tag == "" {
            continue
        }
        switch vField.Interface().(type) {
        case int:
            if vField.Interface() == 0 {
                continue
            }
        case string:
            if vField.Interface() == "" {
                continue
            }
        }
        kv := tag + separator + fmt.Sprintf("%v", vField.Interface())
        res = append(res, kv)
    }
    return res
}