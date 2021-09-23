package gxml

import (
	"testing"
	"time"
)

type TestUser struct {
	Id   int 								`xml:"id"`
	MapType int 							`xml:"mapType"`
	TypeName string							`xml:"typeName"`
	Name string								`xml:"name"`
	MapId string							`xml:"mapId" `
	PortMonster string						`xml:"portMonster"`
}

//func  (t *TestUser)AfterLoad()  {
//	logger.Debug("读取配置:%d",t.Id)
//}
func TestInitialize(t *testing.T)  {

	path := "C:\\trunk\\data\\portGroup.xml"

	m := make(map[int]*TestUser)
	Initialize(path,&TestUser{},m)
	r := len(m)
	logger.Infof("%d",r)
	time.Sleep(time.Minute)
}
