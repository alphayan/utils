package models

import (
	"time"
)

// BaseModel can used for others models
//
// swagger:model
type BaseModel struct {
	Id         int64     `orm:"pk;auto;"`                     //主键id
	CreateDate time.Time `orm:"auto_now_add;type(datetime);"` //创建时间
	ModifyDate time.Time `orm:"auto_now;type(timestamp);"`    //修改时间
}
