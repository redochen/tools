package log

import (
	"encoding/json"
	"fmt"
	. "github.com/redochen/tools/config"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"strconv"
)

type MongodbWriter struct {
	Session *mgo.Session
	Level   int
	StartId int
	EndId   int
}

func NewMongodbWriter() LoggerInterface {
	session, err := mgo.Dial(mongodbUrl)
	if err != nil {
		fmt.Println("无法打开mongodb连接")
		return nil
	}
	return &MongodbWriter{
		Level:   LevelTrace,
		Session: session,
	}
}

func (w *MongodbWriter) SetId() {
	c := w.Session.DB("log").C("lastid")
	var d LastId
	change := mgo.Change{
		Update:    bson.M{"$inc": bson.M{"lastid": 200}},
		ReturnNew: true,
	}
	_, err := c.Find(bson.M{"table": tableName}).Apply(change, &d)
	if err != nil {
		fmt.Println("get last id err:", err.Error())
		c.Insert(&LastId{
			LastId: 200,
			Table:  tableName,
		})
		w.StartId = 1
		w.EndId = 200
	} else {
		if d.LastId > 2e10 {
			c.Update(bson.M{"table": tableName}, bson.M{"$set": bson.M{"lastid": 200}})
			w.StartId = 1
			w.EndId = 200
		} else {
			w.StartId = d.LastId - 200 + 1
			w.EndId = d.LastId
		}

	}
}

func (w *MongodbWriter) Init(level string) error {
	var err error
	w.Level, err = strconv.Atoi(level)
	return err
}

func (w *MongodbWriter) WriteMsg(msg string, level int) error {
	if level > w.Level {
		return nil
	}

	if w.StartId == 0 || w.StartId >= w.EndId {
		w.SetId()
	} else {
		w.StartId++
	}
	var content interface{}
	json.Unmarshal([]byte(msg), &content)
	if content == nil {
		content = msg
	}
	entity := NewLogEntity(w.StartId, level, content)
	c := w.Session.DB("log").C(tableName)
	err := c.Insert(entity)
	return err
}

func (w *MongodbWriter) Flush() {
	return
}

func (w *MongodbWriter) Destroy() {
	return
}

func registerMongo() {
	if nil == Conf || !Conf.IsValid() {
		return
	}

	isUseMongodb = Conf.BoolEx(LogSection, "usemongodb", false)
	if !isUseMongodb {
		return
	}

	mongodbUrl = Conf.StringEx(LogSection, "mongodb", "")
	if mongodbUrl == "" {
		fmt.Println("缺少日志库连接字符串")
		isUseMongodb = false
		return
	}

	category = Conf.IntEx(LogSection, "category", 0)
	businessNo = Conf.IntEx(LogSection, "businessno", 0)
	tableName = Conf.StringEx(LogSection, "tablename", "log")

	Register("mongo", NewMongodbWriter)
}
