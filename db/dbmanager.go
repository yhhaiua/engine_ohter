package db

import "github.com/yhhaiua/engine/util"

var storage map[string]*util.AtomicLog

//Init 初始化数据库数据
func Init(config *MysqlConfig,args ...interface{})  {
	//连接数据库
	if !accessor.LoadConfig(config){
		return
	}
	//自动映射表
	accessor.AutoMigrate(args...)
	//开启数据保存缓存
	globalDbServer.InitDBCache()
	//初始化自增数据
	initStorage(args...)
}

//Create 创建数据库表结构 主键id存在
func Create(entity IEntity)  {
	if entity.GetId() == 0 {
		logger.Errorf("入库主键异常:%s",entity.TableName())
		return
	}
	globalDbServer.AddDBCache(save,entity)
}

//CreateAuto 创建数据库表结构 主键id自动生成
func CreateAuto(entity IEntity,operateId,serverId int)  {
	if entity.GetId() == 0 {
		at := storage[entity.TableName()]
		id := util.BuildPrimaryKey(operateId,serverId,at.IncrementAndGet())
		entity.SetId(id)
	}
	globalDbServer.AddDBCache(save,entity)
}

//Update 更新表数据
func Update(entity IEntity)  {
	if entity.GetId() == 0 {
		logger.Errorf("入库主键异常:%s",entity.TableName())
		return
	}
	globalDbServer.AddDBCache(update,entity)
}

//Delete 删除表数据
func Delete(entity IEntity)  {
	globalDbServer.AddDBCache(remove,entity)
}

//First 查找单条数据
func First(id interface{},entity IEntity)()  {
	accessor.First(id,entity)
}

//FindAll 查找所有数据 entity 切片
func FindAll(entity interface{})()  {
	accessor.FindAll(entity)
}

//FindCond 条件查询 dest 切片  //accessor.db.Where("name = ? AND age >= ?", "jinzhu", "22").Find(&users)
func FindCond(dest interface{},query interface{}, args ...interface{})  {
	accessor.FindCond(dest,query,args...)
}

//ShutDown 缓存数据关闭
func ShutDown()  {
	globalDbServer.ShutDown()
}
func initStorage(args ...interface{})  {
	storage = make(map[string]*util.AtomicLog)
	for _,v := range args{
		ie := v.(IEntity)
		if is,ok := ie.(ISince);!ok || is.GetSince() == false{
			continue
		}
		at := &util.AtomicLog{}
		maxId := getMaxId(ie)
		if maxId > 0{
			at.Set(util.ParseAutoincrement(maxId))
		}else{
			at.Set(0)
		}

		storage[ie.TableName()] = at
	}
}

func getMaxId(entity IEntity)int64  {
	hql := "select max(id) from " + entity.TableName();
	var maxId int64
	var count int64
	result := accessor.db.Table(entity.TableName()).Count(&count)
	if result.Error != nil {
		logger.Errorf("Count error:%s",result.Error.Error())
	}
	if count == 0{
		return 0
	}
	result = result.Raw(hql).Scan(&maxId)
	if result.Error != nil {
		logger.Errorf("getMaxId error:%s",result.Error.Error())
	}
	return maxId
}
