package asynserver

import (
    "encoding/binary"
    "fmt"
    "io"
    "net"
    "os"
    "runtime"
    "strings"
    "time"
)

const (
    keepAliveTime = 10
    maxPackageLen = 1024 * 10
)

type Callback struct {
    //收到包回调
    IPCMessageCallback func(sess *Session, data []byte) bool
    //收到包回调
    MessageCallback func(sess *Session, data []byte) bool
    //连接关闭回调
    CloseCallback func(sess *Session)
    //调用者根据头部信息返回包体长度
    GetSizeCallback func(header []byte) int
    //目标地址host example: 192.168.1.1:8888
    Host string
}

func init() {
    runtime.GOMAXPROCS(runtime.NumCPU())
}

type Session struct {
    ip       net.IP
    Conn     net.Conn    //the tcp connection from client
    RecvChan chan []byte //data from client
    IPCChan  chan []byte //internet process connection message
    ErrChan  chan bool
}

func startSession(client *Session, handler Callback) {
    defer func() {
        client.Conn.Close()
    }()
    for {
        select {
        case msg, ok := <-client.RecvChan:
            if !ok {
                // the cmd channel is closed by the other side
                fmt.Println("The channel is closed by the other side")
                return
            }
            if !handler.MessageCallback(client, msg) {
                //if process error player need load again
                fmt.Println("dispatch error ")
                return
            }
        case msg, ok := <-client.IPCChan:
            if !ok {
                // the cmd channel is closed by the other side
                fmt.Println("The channel is closed by the other side")
                return
            }
            if !handler.IPCMessageCallback(client, msg) {
                //if process error player need load again
                fmt.Println("ipc dispatch error ")
            }
        case <-time.After(20 * time.Second):
            fmt.Println("Timeout and will close")
            runtime.Goexit()
        }
    }
}

func createSession(conn net.Conn) *Session {
    var client Session
    client.ip = net.ParseIP(strings.Split(conn.RemoteAddr().String(), ":")[0])
    client.RecvChan = make(chan []byte, 1024)
    client.Conn = conn
    client.IPCChan = make(chan []byte, 1024)
    return &client
}

func handleConn(sess *Session, handler Callback) {
    header := make([]byte, 2)
    defer func() {
        handler.CloseCallback(sess)
        sess.Conn.Close()
    }()
    //data := make([]byte, 2048)
    go startSession(sess, handler)
    for {
        n, err := io.ReadFull(sess.Conn, header)
        if n == 0 && err == io.EOF {
            //Opposite socket is closed
            fmt.Println("peer socket is closed")
            break
        } else if err != nil {
            //Sth wrong with this socket
            fmt.Println(err)
            break
        }
        //size := binary.LittleEndian.Uint16(header) + 4
        size := handler.GetSizeCallback(header)
        data := make([]byte, size)
        n, err = io.ReadFull(sess.Conn, data[0:size])
        if n == 0 && err == io.EOF {
            fmt.Println("peer socket is closed")
            break
        } else if err != nil {
            fmt.Println(err)
            break
        }
        sess.RecvChan <- data[0:size] //send data to Client to process
    }
}

func cliHandleConn(sess *Session) {
    header := make([]byte, 2)
    defer sess.Conn.Close()
    defer func() {
        sess.ErrChan <- true
    }()
    for {
        n, err := io.ReadFull(sess.Conn, header)
        if n == 0 && err == io.EOF {
            //Opposite socket is closed
            fmt.Println("peer socket is closed")
            break
        } else if err != nil {
            //Sth wrong with this socket
            fmt.Println(err)
            break
        }
        size := binary.LittleEndian.Uint16(header) + 4
        //size := handler.GetSizeCallback(header)
        data := make([]byte, size)
        n, err = io.ReadFull(sess.Conn, data[0:size])
        if n == 0 && err == io.EOF {
            fmt.Println("peer socket is closed")
            break
        } else if err != nil {
            fmt.Println(err)
            break
        }
        sess.RecvChan <- data //send data to Client to process
    }
}

func connectKeepAlive(conn net.Conn, t int) {
    var interval int
    if t > 0 {
        interval = t
    } else {
        interval = keepAliveTime
    }
    ticker := time.NewTicker(time.Second * time.Duration(interval))
    data := make([]byte, 6)
    for t := range ticker.C {
        _, err := conn.Write(data)
        if err != nil {
            fmt.Println(t, err)
            break
        }
    }
}

func ConnectToServer(host string, keepalive int) (bool, *Session) {
    conn, err := net.Dial("tcp", host)
    if err != nil {
        fmt.Println("error connect", host, err.Error())
        return false, nil
    }
    sess := createSession(conn)
    go cliHandleConn(sess)
    go connectKeepAlive(conn, keepalive)
    return true, sess
}

func StartServer(h Callback) {
    listener, err := net.Listen("tcp", h.Host)
    if err != nil {
        fmt.Println("fatal error listening:", err)
        os.Exit(1)
    }
    defer listener.Close()
    fmt.Println("asynServer start listening:", h.Host)
    for {
        conn, err := listener.Accept()
        if err != nil {
            fmt.Println("fail accept", err)
            continue
        }
        sess := createSession(conn)
        go handleConn(sess, h)
    }
}
