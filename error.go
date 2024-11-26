package onig

import "C"
import "fmt"

func errorFromCode(code C.int) error {
	return fmt.Errorf("error from oniguruma: %d", int(code))
}
