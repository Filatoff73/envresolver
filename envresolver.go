package envResolver

import (
	"errors"
	"fmt"
	"os"
	"reflect"
)

//Change all env variables beginnig on 'prefix' in 'input'. If 'ignoreNotFound' then no return if env variable not exist
func Resolve(prefix string, input interface{}, ignoreNotFound bool) error {

	val := reflect.ValueOf(input).Elem()
	if reflect.TypeOf(val).Kind() != reflect.Struct || !val.IsValid() {
		return errors.New("Work only with not nil structs")
	}
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)

		if valueField.CanInterface() {
			value := valueField.Interface()
			if reflect.TypeOf(value).Kind() == reflect.Ptr {
				err := Resolve(prefix, value, ignoreNotFound)
				if err != nil {
					return err
				}
			} else if reflect.TypeOf(value).Kind() == reflect.String {
				str, _ := value.(string)
				if len(str) >= len(prefix) && str[0:len(prefix)] == prefix {
					temp, isExist := os.LookupEnv(str[len(prefix):])
					if !isExist {
						if !ignoreNotFound {
							return errors.New(fmt.Sprintf("env variable %s dosn't exist!", str))
						}
					} else {
						if valueField.CanSet() {
							fmt.Printf("Resolve env variable %s to %s", str, temp)
							valueField.SetString(temp)
						}
					}
				}
			}
		}

	}
	return nil
}
