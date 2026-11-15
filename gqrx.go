package gqrx

import (
	"fmt"
	"log"
)

var (
	Demods = []string{"OFF", "RAW", "AM", "AMS", "LSB", "USB", "CWL", "CWR", "CWU", "CW", "FM", "WFM", "WFM_ST", "WFM_ST_OIRT"}
)

const (
	GQRX_Get_Freq         = "f"
	GQRX_Set_Freq         = "F"
	GQRX_Get_Demod        = "m"
	GQRX_Set_Demod        = "M"
	GQRX_Get_Sig_Strength = "l STRENGTH"
	GQRX_Get_Sql          = "l SQL"
	GQRX_Set_Sql          = "L SQL"
	GQRX_Close_conn       = "Q"
	GQRX_Get_Mute         = "u MUTE"
	GQRX_Set_Mute         = "U MUTE"
)

type GQRX struct {
	Address string
	Port    string
	telnet  Telnet
}

func (g *GQRX) Init() {
	g.telnet = Telnet{
		Address: g.Address,
		Port:    g.Port,
		Message: make(chan string),
		Send:    make(chan string),
	}
	go func() {
		if err := g.telnet.init(); err != nil {
			log.Fatal(err)
		}
	}()
}

// GetFreq get current frequency in Hz as string
func (g *GQRX) GetFreq() string {
	g.telnet.Send <- fmt.Sprintf("%s", GQRX_Get_Freq)
	return <-g.telnet.Message
}

// SetFreq set listen frequency in Mhz as float
func (g *GQRX) SetFreq(input float64) {
	g.telnet.Send <- fmt.Sprintf("%s %d", GQRX_Set_Freq, int(input*1000000))
}

func (g *GQRX) SetDemod(demod string) {
	g.telnet.Send <- fmt.Sprintf("%s %s", GQRX_Set_Demod, demod)
}

func (g *GQRX) List() {
	g.telnet.Send <- fmt.Sprintf("L")
}
