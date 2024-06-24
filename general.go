package manaba

import "fmt"

func e(funcName string, e error) error {
	return fmt.Errorf("%v: %v", funcName, e.Error())
}
