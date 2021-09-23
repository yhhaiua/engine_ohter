package db

import "time"

type EventType int

//数据库操作类型
const (
	save EventType = iota
	update
	remove
)

//数据库缓存队列类型
const (
	PRE5SECOND  = 5 * time.Second		//5秒保存一次
	PRE30SECOND = 30 * time.Second		//30秒保存一次
	PRE1MINUTE = time.Minute			//1分钟保存一次
	PRE5MINUTE = 5 * time.Minute		//5分钟保存一次
)
//数据表结构需要的函数
type IEntity interface {
	TableName() string									//表名
	GetId() int64										//表的主键id
	SetId(id int64)										//设置表的主键id
	GetCron() time.Duration								//保存策略
}

type ISince interface {
	GetSince() bool										//是否实现了自增接口
}