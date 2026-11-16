package gqrx

import (
	"fmt"
	"log"
	"slices"
	"strconv"
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
	GQRX_Get_LNA_Gain     = "l LNA_GAIN"
	GQRX_Set_LNA_Gain     = "L LNA_GAIN"
	GQRX_Close_conn       = "Q"
	GQRX_Get_Mute         = "u MUTE"
	GQRX_Set_Mute         = "U MUTE"
	GQRX_Get_AF           = "l AF" // get audio gain
	GQRX_Set_AF           = "L AF"
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

// GetFreq get current frequency in Hz as int
func (g *GQRX) GetFreq() int {
	g.telnet.Send <- fmt.Sprintf("%s", GQRX_Get_Freq)
	parsedInt, err := strconv.ParseInt(<-g.telnet.Message, 10, 64)
	if err != nil {
		return 0
	}
	return int(parsedInt)
}

// SetFreq set listen frequency in Hz as int, 1Mhz = 1000000
func (g *GQRX) SetFreq(input int) {
	g.telnet.Send <- fmt.Sprintf("%s %d", GQRX_Set_Freq, input)
}

// SetDemod
func (g *GQRX) SetDemod(demod string) {
	if !slices.Contains(Demods, demod) {
		log.Printf("Demod %s not found in Demods", demod)
		return
	}
	g.telnet.Send <- fmt.Sprintf("%s %s", GQRX_Set_Demod, demod)
}

// GetDemod
func (g *GQRX) GetDemod() string {
	g.telnet.Send <- fmt.Sprintf("%s", GQRX_Get_Demod)
	return <-g.telnet.Message
}

// GetSigStrength
func (g *GQRX) GetSigStrength() float64 {
	g.telnet.Send <- fmt.Sprintf("%s", GQRX_Get_Sig_Strength)
	parsedFloat, err := strconv.ParseFloat(<-g.telnet.Message, 64)
	if err != nil {
		return 0
	}
	return parsedFloat
}

// GetLnaGain
func (g *GQRX) GetLnaGain() float64 {
	g.telnet.Send <- fmt.Sprintf("%s", GQRX_Get_LNA_Gain)
	parsedFloat, err := strconv.ParseFloat(<-g.telnet.Message, 64)
	if err != nil {
		return 0
	}
	return parsedFloat
}

// SetLnaGain
func (g *GQRX) SetLnaGain(gain float64) {
	if gain > 50 || gain < 0 {
		fmt.Println("gain must be between 0 and 50")
		return
	}
	g.telnet.Send <- fmt.Sprintf("%s %f", GQRX_Set_LNA_Gain, gain)
}

// GetSql
func (g *GQRX) GetSql() float64 {
	g.telnet.Send <- fmt.Sprintf("%s", GQRX_Get_Sql)
	parsedFloat, err := strconv.ParseFloat(<-g.telnet.Message, 64)
	if err != nil {
		return 0
	}
	return parsedFloat
}

// SetSql range -150.0db - 0.0db
func (g *GQRX) SetSql(input float64) {
	if input > 0.0 || input < -150.0 {
		fmt.Println("input must be between 0.0 and -150.0")
		return
	}
	g.telnet.Send <- fmt.Sprintf("%s %f", GQRX_Set_Sql, input)
}

// GetAF get audio gain
func (g *GQRX) GetAF() float64 {
	g.telnet.Send <- fmt.Sprintf("%s", GQRX_Get_AF)
	parsedFloat, err := strconv.ParseFloat(<-g.telnet.Message, 64)
	if err != nil {
		return 0
	}
	return parsedFloat
}

// SetAF set audio gain range -80.0 - 50.0
func (g *GQRX) SetAF(input float64) {
	if input > 50 || input < -80 {
		fmt.Println("input for audio gain must be between -80.0 and 50.0")
		return
	}
	g.telnet.Send <- fmt.Sprintf("%s %f", GQRX_Set_AF, input)
}

// GetMute get mute status, 0 = false(unmuted), 1 = true(muted)
func (g *GQRX) GetMute() int {
	g.telnet.Send <- fmt.Sprintf("%s", GQRX_Get_Mute)
	parsedInt, err := strconv.ParseInt(<-g.telnet.Message, 0, 64)
	if err != nil {
		return 0
	}
	return int(parsedInt)
}

// SetMute set audio mute, 0 = false(unmuted), 1 = true(muted)
func (g *GQRX) SetMute(input int) {
	if input < 0 || input > 1 {
		fmt.Println("input must be between 0 and 1, 0 = off(unmuted), 1 = on(muted)")
		return
	}
	g.telnet.Send <- fmt.Sprintf("%s %d", GQRX_Set_Mute, input)
}

// CloseConnection
func (g *GQRX) CloseConnection() {
	g.telnet.Send <- fmt.Sprintf("%s", GQRX_Close_conn)
}

// CommandHelp print command help
func (g *GQRX) CommandHelp(command string) {
	g.telnet.Send <- fmt.Sprintf("%s %s", command, "?")
	for {
		fmt.Println(<-g.telnet.Message)
	}
}
