package net

import (
	"github.com/yhhaiua/engine/buffer"
	"github.com/yhhaiua/engine/handler"
	"github.com/gorilla/websocket"
	"io"
	"sync"
)

type WSConn struct {
	sync.Mutex
	conn      *websocket.Conn
	receive   *buffer.ByteBuf
	listener  SocketListener
	hd        handler.Handler
	data      chan []byte
	connected bool
	ip        string
	closeData bool
}

func newWSConn(conn *websocket.Conn,listener SocketListener,ip string) *WSConn {
	t := new(WSConn)
	t.conn = conn
	t.listener = listener
	t.receive = buffer.NewByteBuf()
	t.data = make(chan []byte,DataLength)
	t.connected = true
	hd,err := handler.NewWsDecoder()
	if err != nil{
		logger.Errorf("new WSConn err: %v",err)
		return nil
	}
	t.hd = hd
	t.ip = ip
	return t
}
func (t *WSConn) start() {

	t.listener.OnConnected(t)
	go t.read()
	go t.run()
}

func (t *WSConn)read()  {

	defer func() {
		if r := recover();r != nil{
			logger.TraceErr(r)
			t.close()
			t.listener.OnDisconnected(t)
		}
	}()
	for  {
		_, b, err := t.conn.ReadMessage()
		if err == io.EOF {
			logger.Infof("远程连接：%s,关闭", t.conn.RemoteAddr().String());
			t.close()
			t.listener.OnDisconnected(t)
			return
		}
		if err != nil{
			logger.Errorf("read err: %s",err.Error())
			t.close()
			t.listener.OnDisconnected(t)
			return
		}
		_, _ = t.receive.Write(b)
		msg,err := t.hd.Decode(t.receive)
		if err != nil{
			logger.Errorf("msg err: %s",err.Error())
			continue
		}
		if msg != nil{
			t.listener.OnData(t,msg)
		}
	}
}

func (t *WSConn)run()  {
	defer func() {
		if r := recover();r != nil{
			logger.TraceErr(r)
			t.close()
		}
	}()

	end := false
	for  {
		select {
		case msg, ok := <-t.data:
			if !ok || msg == nil{
				logger.Infof("run destroy:%s",t.String())
				end = true
				break
			}
			if msg != nil && t.connected{
				err := t.conn.WriteMessage(websocket.BinaryMessage, msg)
				if err != nil{
					logger.Errorf(" WriteMessage:%s",err.Error())
					end = true
					break
				}
			}
		}
		if end{
			t.close()
			t.Lock()
			t.closeData = true
			t.Unlock()
			break
		}
	}

}

//Destroy 目标通知断开后销毁
func (t *WSConn) Destroy() {
	t.Lock()
	defer t.Unlock()

	t.doDestroy()
}

//WriteAndFlush 向目标发送数据
func (t *WSConn)WriteAndFlush(msg []byte)()  {
	t.Lock()
	defer t.Unlock()
	if t.closeData {
		return
	}
	t.doWrite(msg);
}

//Close 主动关闭连接（调用前先向目标发送关闭信息）
func (t *WSConn)Close()  {
	t.Lock()
	defer t.Unlock()
	if t.closeData {
		return
	}
	t.doWrite(nil)
	t.closeData = true
}
func (t *WSConn)String() string  {
	return t.ip
}
func (t *WSConn)close(){

	if !t.connected{
		return
	}
	t.connected = false
	_ = t.conn.Close()
}

func (t *WSConn) doDestroy() {
	t.close()
	if !t.closeData {
		close(t.data)
		t.closeData = true
	}

}

func (t *WSConn) doWrite(msg []byte) {
	if len(t.data) == cap(t.data) {
		logger.Errorf("close conn: channel full")
		t.doDestroy()
		return
	}
	t.data <- msg
}