package asynserver

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
)

type Session struct {
	conn      net.Conn    //the tcp connection from client
	recvChan  chan []byte //data from client
	sendChan  chan []byte //data to client
	closeChan chan bool
	exitChan  chan bool
	exitGroup sync.WaitGroup
	uniqID    uint64
	server    *Server
	data      interface{}
	lock      sync.Mutex
	ok        bool
}

func (sess *Session) IsOK() bool {
	sess.lock.Lock()
	defer sess.lock.Unlock()
	return sess.ok
}

func (sess *Session) Close() {
	sess.conn.Close()
	sess.lock.Lock()
	defer sess.lock.Unlock()
	if sess.ok {
		sess.ok = false
		close(sess.closeChan)
	}
}

func (sess *Session) handleSend() {
	defer func() {
		sess.Close()
		close(sess.sendChan)
		if x := recover(); x != nil {
			fmt.Println("send Panic, the panic is", x)
		}
	}()
	for {
		select {
		case msg, ok := <-sess.sendChan:
			if !ok {
				break
			}
			_, err := sess.conn.Write(msg)
			if err != nil {
				break
			}
		case <-sess.closeChan:
			break
		}
	}
}

func (sess *Session) handleDispatch() {
	defer func() {
		sess.Close()
		close(sess.recvChan)
		//for msg := range sess.recvChan {
		//sess.server.MessageCallback(sess, msg)
		//}
		if x := recover(); x != nil {
			fmt.Println("dispatch Panic, the panic is", x)
		}
	}()
	for {
		//接受数据 调用回调
		select {
		case msg, ok := <-sess.recvChan:
			if !ok {
				fmt.Println("The channel is closed by the other side")
				return
			}
			if !sess.server.messageCallback(sess, msg) {
				fmt.Println("dispatch error ")
				return
			}
		case <-sess.closeChan:
		}
	}
}

func (sess *Session) RunInQueue(msg []byte) bool {
	defer func() {
		if x := recover(); x != nil {
			fmt.Println("Panic, the panic is", x)
		}
	}()
	if sess.IsOK() {
		select {
		case sess.recvChan <- msg:
			return true
		case <-sess.closeChan:
			return false
		default:
			return false
		}
	}
	return false
}

func (sess *Session) Send(msg []byte) bool {
	if sess.IsOK() {
		select {
		case sess.sendChan <- msg:
			return true
		default:
			return false
		}
	}
	return false
}

func (sess *Session) handleRecv() {
	defer func() {
		if x := recover(); x != nil {
			fmt.Println("Panic, the panic is", x)
		}
	}()
	header := make([]byte, 2)
	for {
		n, err := io.ReadFull(sess.conn, header)
		if n == 0 && err == io.EOF {
			//Opposite socket is closed
			fmt.Println("peer socket is closed")
			break
		} else if err != nil {
			//Sth wrong with this socket
			fmt.Println(err)
			break
		}
		size := binary.LittleEndian.Uint16(header)
		data := make([]byte, size)
		n, err = io.ReadFull(sess.conn, data[0:size])
		if n == 0 && err == io.EOF {
			fmt.Println("peer socket is closed")
			break
		} else if err != nil {
			fmt.Println(err)
			break
		}
		sess.recvChan <- data[0:size] //send data to Client to process
	}
}

func (sess *Session) SetData(m interface{}) {
	sess.lock.Lock()
	defer sess.lock.Unlock()
	sess.data = m
}

func (sess *Session) GetData() (m interface{}) {
	sess.lock.Lock()
	defer sess.lock.Unlock()
	m = sess.data
	return
}

func (sess *Session) Start() {
	defer func() {
		sess.exitGroup.Done()
		sess.server.delSession(sess)
		if x := recover(); x != nil {
			fmt.Println("start Panic, the panic is", x)
		}
	}()
	go sess.handleRecv()
	go sess.handleSend()
	sess.handleDispatch()
	sess.server.closeCallback(sess, nil)
}

func (sess *Session) WaitExit() {
	sess.exitGroup.Wait()
}

type Server struct {
	bindAddr  string
	listener  net.Listener
	sessIndex uint64
	sessions  map[uint64]*Session
	chanSize  uint
	lock      sync.Mutex
	//收到包回调
	messageCallback func(sess *Session, data []byte) bool
	//连接关闭回调
	closeCallback func(sess *Session, data []byte) bool
}

func NewServer(host string) {
	var server Server
	server.bindAddr = host
	server.sessions = make(map[uint64]*Session)
}

func (this *Server) SetCallback(m, c func(sess *Session, data []byte) bool) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.messageCallback = m
	this.closeCallback = c
}

func (this *Server) SetChanSize(s uint) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.chanSize = s
}

func (this *Server) getChanSize() uint {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.chanSize
}

func (this *Server) createSession(conn net.Conn) *Session {
	var client Session
	client.exitGroup.Add(1)
	client.conn = conn
	size := this.getChanSize()
	if size == 0 {
		size = 64
	}
	client.recvChan = make(chan []byte, size)
	client.ok = true
	client.server = this
	this.sessIndex++
	client.uniqID = this.sessIndex
	this.lock.Lock()
	defer this.lock.Unlock()
	this.sessions[this.sessIndex] = &client
	return &client
}

func (this *Server) delSession(sess *Session) {
	this.lock.Lock()
	defer this.lock.Unlock()
	delete(this.sessions, this.sessIndex)
}

func (this *Server) CloseAllSessions() {
	this.lock.Lock()
	defer this.lock.Unlock()
	for _, v := range this.sessions {
		v.Close()
	}
}

func (this *Server) Start() error {
	var err error
	if this.messageCallback == nil || this.closeCallback == nil {
		return err
	}
	this.listener, err = net.Listen("tcp", this.bindAddr)
	if err != nil {
		fmt.Println("fatal error listening:", err)
		os.Exit(1)
	}
	defer this.listener.Close()
	for {
		conn, err := this.listener.Accept()
		if err != nil {
			fmt.Println("fail accept", err)
			continue
		}
		sess := this.createSession(conn)
		go sess.Start()
	}
	return nil
}
