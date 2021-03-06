/*
说明:此包有副作用的函数
主键:mongo自带的bsonID
*/
package xmdb

import (
	"log"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var sconn *mgo.Session

//获取session连接
func GetSConn() *mgo.Session {
	if sconn == nil {
		var err error
		sconn, err = mgo.Dial("mongodb://localhost:27017")
		if err != nil {
			log.Fatal("mongod 数据库连接失败:", err)
		}
		sconn.SetMode(mgo.Strong, true)
		sconn.Safe()
	}
	return sconn.Clone()
}

//辅助函数:建立session连接
func WithSConn(execute func(sconn *mgo.Session) error) error {
	sconn := GetSConn()
	defer sconn.Close()
	return execute(sconn)
}

//通过id查找一条数据
func FindId(dbName, cName string, des interface{}, id bson.ObjectId) error {
	return WithSConn(func(sconn *mgo.Session) error {
		return sconn.DB(dbName).C(cName).FindId(id).One(des)
	})
}

//通过条件查找一条数据
func FindOne(dbName, cName string, des, selector interface{}) error {
	return WithSConn(func(sconn *mgo.Session) error {
		return sconn.DB(dbName).C(cName).Find(selector).One(des)
	})
}

//通过条件查找多条数据
func Find(dbName, cName string, des, selector interface{}) error {
	return WithSConn(func(sconn *mgo.Session) error {
		return sconn.DB(dbName).C(cName).Find(selector).All(des)
	})
}

//通过条件找到的条数
func Count(dbName, cName string, selector interface{}) int {
	var (
		count int
		err   error
	)
	WithSConn(func(sconn *mgo.Session) error {
		count, err = sconn.DB(dbName).C(cName).Find(selector).Count()
		if err != nil {
			log.Fatal("get count error:", err)
		}
		return nil
	})
	return count
}

//插入数据(多条数据插入时，有主见重复的数据，依次插入每条数据，直到插入到重复的主键处报错，前面主键没有重复的数据已经插入)
func Insert(dbName string, cName string, src ...interface{}) {
	err := WithSConn(func(sconn *mgo.Session) error {
		return sconn.DB(dbName).C(cName).Insert(src...)
	})
	if err != nil {
		log.Fatal("insert error:", err)
	}
}

//更新数据
func Update(dbName string, cName string, selector, update interface{}) error {
	return WithSConn(func(sconn *mgo.Session) error {
		return sconn.DB(dbName).C(cName).Update(selector, update)
	})
}

//更新数据
func UpdateAll(dbName string, cName string, selector, update interface{}) error {
	return WithSConn(func(sconn *mgo.Session) error {
		_, err := sconn.DB(dbName).C(cName).UpdateAll(selector, update)
		return err
	})
}

//更新数据
func UpdateId(dbName string, cName string, id bson.ObjectId, update interface{}) error {
	return WithSConn(func(sconn *mgo.Session) error {
		return sconn.DB(dbName).C(cName).UpdateId(id, update)
	})
}

//组装pipe函数
func ComposePipe(cpfs ...func() bson.M) []bson.M {
	bms := []bson.M{}
	for _, f := range cpfs {
		bms = append(bms, f())
	}
	return bms
}

//填充pipe函数
func FillPipe(cpfs ...func([]bson.M) []bson.M) []bson.M {
	bms := []bson.M{}
	for _, f := range cpfs {
		bms = f(bms)
	}
	return bms
}

//管道筛选
func Pipe(dbName, cName string, des, pipe interface{}) error {
	return WithSConn(func(sconn *mgo.Session) error {
		return sconn.DB(dbName).C(cName).Pipe(pipe).All(des)
	})
}

//管道筛选
func PipeOne(dbName, cName string, des, pipe interface{}) error {
	return WithSConn(func(sconn *mgo.Session) error {
		return sconn.DB(dbName).C(cName).Pipe(pipe).One(des)
	})
}
