package quick

import (
	"fmt"
)

func ErrorStatusCodeNotOk(code int) error {
	return fmt.Errorf("unexpected status code %d of the http response", code)
}
