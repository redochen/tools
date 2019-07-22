package log

import (
	"strconv"
	"time"
)

type LogEntity struct {
	Id         int         `bson:"_id"`
	BusinessNo int         `bson:"bno"`
	Level      int         `bson:"le"`
	Time       int         `bson:"ti"`
	Category   int         `bson:"ca"`
	Content    interface{} `bson:"co"`
}

func NewLogEntity(id, level int, content interface{}) *LogEntity {
	ti, _ := strconv.Atoi(time.Now().Format("0601021504"))
	return &LogEntity{
		Id:         id,
		BusinessNo: businessNo,
		Level:      level,
		Time:       ti,
		Category:   category,
		Content:    content,
	}
}

type LastId struct {
	Table  string
	LastId int
}
