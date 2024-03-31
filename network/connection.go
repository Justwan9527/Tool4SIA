package network

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net"
)

var (
	buf    bytes.Buffer
	logger = log.New(&buf, "logger: ", log.Lshortfile)
)

// 连接模块
type Connection struct {

	//当前连接的socket tcp
	Conn *net.TCPConn

	//连接的ID
	ConnID uint32

	//当前连接的状态
	isClosed bool

	//告知当前连接已经退出的/停止 channel
	ExitChan chan bool

	//无缓冲的管道，用于读，写go协成之间的消息通信
	msgChan chan []byte

	//channel for reading
	MsgReadChan chan []byte

	//该连接的处理方法router
	//Router ziface.IRouter

}

// 初始化连接模块的方法
func NewConnection(conn *net.TCPConn, connID uint32) *Connection {
	c := &Connection{
		Conn:        conn,
		ConnID:      connID,
		isClosed:    false,
		msgChan:     make(chan []byte, 1),
		ExitChan:    make(chan bool, 1),
		MsgReadChan: make(chan []byte, 1),
	}

	return c
}

func (c *Connection) StartReader() {
	fmt.Println("Reader goroutine is running")
	//logger.Println("Reader goroutine is running")

	//defer fmt.Println("connID=", c.ConnID, "Reader is exit, remote addrr=", c.RemoteAddr().String())
	defer c.Stop()

	for {

		//read 512 from tcp connection
		//nMsgLen, err := io.ReadFull(c.GetTCPConnection(), data);
		var data []byte = make([]byte, 512)
		if nMsgLen, err := c.Conn.Read(data); err != nil {
			fmt.Println("ReadFull error:", err)
			if err.Error() != "EOF" {
				break
			}

		} else {
			fmt.Println("Start to read msg len:", nMsgLen)
			c.MsgReadChan <- data
			fmt.Println("End to read msg")

		}

	}

	fmt.Println("Reader goroutine exits")

}

// 写消息的go协成，专门
func (c *Connection) StartWriter() {
	fmt.Println("Writer goroutine is running")
	//defer fmt.Println(c.Conn.RemoteAddr().String(), "conn Writer exit!")

	//不断的阻塞的等待channel的消息，进行写给客户端
	for {
		select {
		case data := <-c.msgChan:
			//有数据要写给客户端
			if nMsgLen, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send data error:", err)
				return
			} else {
				fmt.Println("Send data Len:", nMsgLen, "data:", string(data))
			}

		case <-c.ExitChan:
			//代表Reader已经退出，此时Writer也要退出
			fmt.Println("StartWriter exit !!!!")
			return

		}
	}

}

// 启动连接 让当前的连接准备开始工作
func (c *Connection) Start() {
	//fmt.Println("Conn Strar()...ConnID=", c.ConnID)

	//启动当前连接的写业务
	go c.StartWriter()

	//启动从当前连接的读业务
	go c.StartReader()

}

// 停止连接 结束当前连接的工作
func (c *Connection) Stop() {
	//fmt.Println("Conn Stop()...ConnID=", c.ConnID)
	if c.isClosed {
		return
	}

	c.isClosed = true

	c.Conn.Close()

	//告知Writer关闭
	c.ExitChan <- true

	//回收资源
	close(c.ExitChan)
	close(c.msgChan)
}

// 获取当前连接绑定的socket conn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn

}

// 获取当前连接模块的连接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID

}

// 获取远程客户端的TCP状态 IP port
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// 发送数据 将数据发送给远程的客户端
func (c *Connection) Send(data []byte) error {

	_, err := c.Conn.Write(data)

	return err
}

// get connection status
func (c *Connection) GetConnStatus() bool {
	return c.isClosed
}

// 提供一个SendMsg方法 将我们要发送给客户端的数据， 先进行封包，再发送
func (c *Connection) SendMsg(data []byte) error {

	fmt.Println("!!!!Start to send msg data = ", string(data))
	if c.isClosed {
		return errors.New("Connection is closed when send msg")
	}

	/*
		//将数据发送给客户端
		if _, err := c.Conn.Write(binaryMsg); err != nil {
			fmt.Println("Write msg id=", msgId, "error=", err)
			return errors.New("connn Write error")
		}
	*/

	//将数据发送给客户端
	fmt.Println("Start to send msg")
	c.msgChan <- data
	fmt.Println("End to send msg")
	return nil
}
