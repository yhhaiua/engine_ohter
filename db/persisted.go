package db

import (
	"github.com/yhhaiua/engine/util/concurrent"
	"sync"
	"time"
)

type TimingPersisted struct {
	elements 	concurrent.ConcurrentMap
	name 		string
	timer		*time.Ticker
	shut		chan bool
	wg			*sync.WaitGroup
	stop 		bool
	m			*sync.RWMutex
}

func newTimingPersisted(cron time.Duration) *TimingPersisted  {
	persisted := new(TimingPersisted)
	persisted.name = cron.String() + "db saver"
	persisted.elements =concurrent.New()
	persisted.shut = make(chan bool)
	persisted.timer = time.NewTicker(cron)
	persisted.m = new(sync.RWMutex)
	go persisted.run()
	return persisted
}

func (t *TimingPersisted)shutDown()  {
	t.stop = true
	logger.Warnf("关闭程序保存数据开始:%s",t.name)
	t.wg = &sync.WaitGroup{}
	t.wg.Add(1)
	t.shut <- true
	t.wg.Wait()
	logger.Warnf("关闭程序保存数据完成:%s",t.name)
}
func (t *TimingPersisted)put(element *Element)  {
	if element == nil || t.stop{
		return
	}
	t.m.RLock()
	t.elements.Set(element.getIdentity(),element)
	t.m.RUnlock()
}

func (t *TimingPersisted)run()  {
	defer func() {
		if r := recover();r != nil{
			logger.TraceErr(r)
			go t.run()	//出现退出异常，再次启动
		}
	}()
	for{
		select {
			//数据处理
		case <- t.timer.C:
			//定时器处理
			t.timerProcessing()
		case <- t.shut:
			ret := t.timerProcessing()	//对数据进行保存
			logger.Warnf("关闭程序保存数据:%s,保存数据:%d条",t.name,ret)
			t.wg.Done()
			return
		}
	}
}
func (t *TimingPersisted)clearElements()map[interface{}]interface{}{
	t.m.Lock()
	ret := t.elements.Items()
	t.elements.Clear()
	t.m.Unlock()
	return ret
}

func (t *TimingPersisted)timerProcessing() int  {
	defer func() {
		if r := recover();r != nil{
			logger.TraceErr(r)
		}
	}()
	count := t.elements.Count()
	if count == 0{
		return 0
	}
	el := t.clearElements()
	for _, v := range el {
		element := v.(*Element)
		switch element.event {
		case save:
			accessor.Create(element.dbObject)
		case update:
			accessor.Save(element.dbObject)
		case remove:
			accessor.Delete(element.dbObject.GetId(),element.dbObject)
		}
	}
	return len(el)
}