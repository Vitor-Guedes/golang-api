package types

import (
	"fmt"
	"time"
)

type JSONTime time.Time

func (to JSONTime) MarshalJSON() ([]byte, error) {
	formatted := time.Time(to).Format("2006-01-02 15:04:05")
	return []byte(fmt.Sprintf(`"%s"`, formatted)), nil
}