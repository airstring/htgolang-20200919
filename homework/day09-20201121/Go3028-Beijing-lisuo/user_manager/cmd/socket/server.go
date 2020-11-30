package socket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
)

const (
	proto   = "tcp"
	addr    = ":8081"
	headLen = 5
)

var helpMsg = `
+-------+---------------------+
|  CMD  |      Function       |
+-------+---------------------+
| help  | ShowHelp            |
| add   | AddUser             |
| show  | ShowCurrentUserList |
| mod   | ModifyUser          |
| del   | DelUser             |
| get   | QueryUser           |
+-------+---------------------+
`

// Head  represents operation and status
type Head struct {
	Operation string
	Message   string
	Status    int
}

/*

+-------+---------------------+
|  CMD  |      Function       |
+-------+---------------------+
| get   | QueryUser           |
| h     | ShowHelp            |
| show  | ShowCurrentUserList |
| q     | utils.Quit          |
| del   | DelUser             |
| help  | ShowHelp            |
| cls   | utils.ClearScreen   |
| quit  | utils.Quit          |
| mycmd | DoMap               |
| rot   | Rotate              |
| add   | AddUser             |
| mod   | ModifyUser          |
| save  | SaveUsers           |
| Q     | utils.Quit          |
| exit  | utils.Quit          |
+-------+---------------------+

*/

// Server for remote user manager
func Server() {
	listener, err := net.Listen(proto, addr)
	if err != nil {
		panic(err)
	}
	conn, errA := listener.Accept()
	if errA != nil {
		panic(errA)
	}

	res := ReadHead(conn)
	fmt.Println("Head: ", res)
	if res.Operation == "help" {
		res.Message = helpMsg
		res.Status = 200
		showClientHelp(conn, &res)
	}

	conn.Close()

}

func showClientHelp(c net.Conn, h *Head) {
	WriteHead(c, *h)
}

// ============== protocol =============

// WriteHead wrap WriteHeadLen and WriteHeadBody
func WriteHead(c net.Conn, h Head) {
	WriteHeadLen(c, h)
	WriteHeadBody(c, h)
}

// WriteHeadLen send json head len to client
func WriteHeadLen(c net.Conn, h Head) {
	bt, err := json.Marshal(h)
	if err != nil {
		c.Close()
		panic(err)
	}
	contentLen := len(string(bt))
	lenStr := fmt.Sprintf("%05d", contentLen)
	_, errW := c.Write([]byte(lenStr))
	if errW != nil {
		c.Close()
		panic(errW)
	}
}

// WriteHeadBody send json head to server
func WriteHeadBody(c net.Conn, h Head) {
	b, _ := json.Marshal(h)
	_, errW := c.Write(b)
	if errW != nil {
		c.Close()
		panic(errW)
	}
}

// ReadHead read json response head from server
func ReadHead(c net.Conn) Head {
	conLen := readHeadLen(c)
	var d = make([]byte, conLen)
	buf := bytes.NewBuffer(d)
	_, errR := c.Read(buf.Bytes())
	if errR != nil {
		c.Close()
		panic(errR)
	}
	responseBytes := buf.Bytes()
	var response = Head{}
	errUnmarshal := json.Unmarshal(responseBytes, &response)
	if errUnmarshal != nil {
		panic(errUnmarshal)
	}
	return response
}

func readHeadLen(c net.Conn) int {
	var buf = make([]byte, headLen)
	_, errRead := c.Read(buf)
	if errRead != nil {
		c.Close()
		panic(errRead)
	}
	len, err := strconv.Atoi(string(buf))
	if err != nil {
		panic(err)
	}
	return len
}
