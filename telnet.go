package gqrx

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

type Telnet struct {
	Address string
	Port    string
	Message chan string
	Send    chan string
}

const (
	IAC  byte = 255
	DONT byte = 254
	DO   byte = 253
	WONT byte = 252
	WILL byte = 251
	SB   byte = 250
	SE   byte = 240
)

func (t *Telnet) init() error {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%s", t.Address, t.Port), 10*time.Second)
	if err != nil {
		return fmt.Errorf("init::connection error: %v", err)

	}
	defer conn.Close()
	go t.listen(conn)
	for {
		line := <-t.Send
		_, err := fmt.Fprintf(conn, "%s\r\n", line)
		if err != nil {
			return fmt.Errorf("init::connection error: %v", err)
		}
	}
}

func (t *Telnet) listen(conn net.Conn) {
	rns := ""
	reader := bufio.NewReader(conn)
	for {
		b, err := reader.ReadByte()
		if err != nil {
			if err != io.EOF {
				fmt.Printf("\nRead error: %v\n", err)
			}
			os.Exit(0)
		}
		if b == IAC {
			cmd, err := reader.ReadByte()
			if err != nil {
				fmt.Printf("listen::command error: %v\n", err)
				return
			}
			switch cmd {
			case DO, DONT, WILL, WONT:
				opt, err := reader.ReadByte()
				if err != nil {
					fmt.Printf("listen::command error: %v\n", err)
					return
				}
				if cmd == DO {
					conn.Write([]byte{IAC, WONT, opt})
				} else if cmd == WILL {
					conn.Write([]byte{IAC, DONT, opt})
				}
			case SB:
				for {
					subByte, err := reader.ReadByte()
					if err != nil {
						fmt.Printf("listen::command error: %v\n", err)
						return
					}
					if subByte == IAC {
						next, _ := reader.ReadByte()
						if next == SE {
							break
						}
					}
				}
			case IAC:
				fmt.Print(string(IAC))
			}
			continue
		}
		if b == '\n' {
			if rns == "RPRT 0" {
				rns = ""
				continue
			}
			t.Message <- rns
			rns = ""
		} else {
			rns += fmt.Sprintf("%c", b)
		}
	}
}
