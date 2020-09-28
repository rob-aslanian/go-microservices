package account

import "time"

// Birthday ...
type Birthday struct {
	Birthday   time.Time  `bson:"birthday"`
	Permission Permission `bson:"permission"`
}
