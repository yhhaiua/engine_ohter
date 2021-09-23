package db

import (
	"testing"
	"time"
)

type User struct {
	Model
	Name string
}

func (u *User) TableName() string {
	return "User"
}
func (u *User) GetId() int64 {
	return u.ID
}

func (u *User) SetId(id int64) {
	u.ID = id
}

func (u *User) GetCron() time.Duration {
	return PRE5SECOND
}
func (u *User)GetSince() bool  {
	return true
}

func TestAccessor_LoadConfig(t *testing.T) {
	config := &MysqlConfig{}
	config.Path = "127.0.0.1:3306"
	config.Config = "charset=utf8mb4&parseTime=True&loc=Local"
	config.Dbname ="wbcharge"
	config.Password ="123456"
	config.Username = "root"
	Init(config,&User{})
	var users []*User
	accessor.FindAll(&users)
	select {
	}
}
