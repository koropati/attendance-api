package presence

import "time"

func GetCurrentMillis() int64 {
	return time.Now().UnixNano() / 1000000
}
