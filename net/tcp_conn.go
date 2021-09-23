package net

import (
	"github.com/yhhaiua/engine/buffer"
	"github.com/yhhaiua/engine/handler"
	"io"
	"net"
	"sync"
)

const DataLength  = 100

type TCPConn struct {
	sync.Mutex
	conn      net.Conn
	receive   *buffer.ByteBuf
	listener  SocketListener
	hd        handler.Handler
	data      chan []byte
	connected bool
	closeData bool
}

func newTcpConn(conn net.Conn,listener SocketListener) *TCPConn {
	t := new(TCPConn)
	t.conn = conn
	t.listener = listener
	t.receive = buffer.NewByteBuf()
	t.data = make(chan []byte,DataLength)
	t.connected = true
	hd,err := handler.NewLengthDecoder()
	if err != nil{
		logger.Errorf("new TcpConn err: %s",err.Error())
		return nil
	}
	t.hd = hd
	return t
}

func (t *TCPConn) start() {

	t.listener.OnConnected(t)
	go t.read()
	go t.run()
}

func (t *TCPConn)read()  {

	defer func() {
		if r := recover();r != nil{
			logger.TraceErr(r)
			t.close()
			t.listener.OnDisconnected(t)
		}
	}()
	for  {
		err := t.receive.ReadFrom(t.conn)
		if err == io.EOF {
			logger.Infof("远程连接：%s,关闭",t.conn.RemoteAddr().String());
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

func (t *TCPConn)run()  {
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
			if !ok || msg == nil {
				logger.Infof("run destroy:%s",t.String())
				end = true
				break
			}
			if msg != nil && t.connected{
				_,err := t.conn.Write(msg)
				if err != nil{
					logger.Errorf(" Write:%s",err.Error())
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
func (t *TCPConn) Destroy() {
	t.Lock()
	defer t.Unlock()

	t.doDestroy()
}

func (t *TCPConn)WriteAndFlush(msg []byte)()  {
	t.Lock()
	defer t.Unlock()
	if t.closeData {
		return
	}
	t.doWrite(msg);
}

func (t *TCPConn)Close()  {
	t.Lock()
	defer t.Unlock()
	if t.closeData {
		return
	}
	t.doWrite(nil)
	t.closeData = true
}
func (t *TCPConn)String() string  {
	return t.conn.RemoteAddr().String()
}
func (t *TCPConn)close(){
	if !t.connected{
		return
	}
	t.connected = false
	_ = t.conn.Close()
}

func (t *TCPConn) doDestroy() {
	t.close()
	if !t.closeData {
		close(t.data)
		t.closeData = true
	}

}

func (t *TCPConn) doWrite(msg []byte) {
	if len(t.data) == cap(t.data) {
		logger.Errorf("close conn: channel full")
		t.doDestroy()
		return
	}
	t.data <- msg
}