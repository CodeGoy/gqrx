package gqrx

import (
	"errors"
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
	GQRX_Get_AF           = "l AF" // get audio gain
	GQRX_Set_AF           = "L AF"
	GQRX_Get_Mute         = "u MUTE"
	GQRX_Set_Mute         = "U MUTE"
	GQRX_Get_DSP          = "u DSP"
	GQRX_Set_DSP          = "U DSP"
	GQRX_Close_conn       = "Q"
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

// GetDSP get DSP status
func (g *GQRX) GetDSP() (int, error) {
	g.telnet.Send <- fmt.Sprintf("%s", GQRX_Get_DSP)
	return toInt(<-g.telnet.Message)
}

// SetDSP set DSP status
func (g *GQRX) SetDSP(input int) {
	g.telnet.Send <- fmt.Sprintf("%s %d", GQRX_Set_DSP, input)
}

// GetFreq get current frequency in Hz as int
func (g *GQRX) GetFreq() (int, error) {
	g.telnet.Send <- fmt.Sprintf("%s", GQRX_Get_Freq)
	return toInt(<-g.telnet.Message)
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
func (g *GQRX) GetSigStrength() (float64, error) {
	g.telnet.Send <- fmt.Sprintf("%s", GQRX_Get_Sig_Strength)
	return toFloat(<-g.telnet.Message)
}

// GetLnaGain get LNA gain, range 0.0 - 50.0
func (g *GQRX) GetLnaGain() (float64, error) {
	g.telnet.Send <- fmt.Sprintf("%s", GQRX_Get_LNA_Gain)
	return toFloat(<-g.telnet.Message)
}

// SetLnaGain set LNA gain, range 0.0 - 50.0
func (g *GQRX) SetLnaGain(gain float64) {
	if gain > 50 || gain < 0 {
		fmt.Println("gain must be between 0 and 50")
		return
	}
	g.telnet.Send <- fmt.Sprintf("%s %f", GQRX_Set_LNA_Gain, gain)
}

// GetSql range -150.0db - 0.0db
func (g *GQRX) GetSql() (float64, error) {
	g.telnet.Send <- fmt.Sprintf("%s", GQRX_Get_Sql)
	return toFloat(<-g.telnet.Message)
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
func (g *GQRX) GetAF() (float64, error) {
	g.telnet.Send <- fmt.Sprintf("%s", GQRX_Get_AF)
	return toFloat(<-g.telnet.Message)
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
func (g *GQRX) GetMute() (int, error) {
	g.telnet.Send <- fmt.Sprintf("%s", GQRX_Get_Mute)
	return toInt(<-g.telnet.Message)
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
	g.telnet.Send <- fmt.Sprintf("%s", "?")
	for {
		fmt.Println(<-g.telnet.Message)
	}
}

func toInt(input string) (int, error) {
	if input == "RPRT 1" {
		return 0, errors.New("command failed")
	}
	parsedInt, err := strconv.ParseInt(input, 0, 64)
	if err != nil {
		return 0, err
	}
	return int(parsedInt), nil
}

func toFloat(input string) (float64, error) {
	if input == "RPRT 1" {
		return 0.0, errors.New("command failed")
	}
	return strconv.ParseFloat(input, 64)
}
