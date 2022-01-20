package wrapers

import (
	"fmt"
	"regexp"
	"time"

	telnet "github.com/a4lex/go-telnet"
)

type MyTelnet struct {
	*telnet.Conn
	network     string
	addr        string
	connTimeout time.Duration
	logger      func(string)
	connected   bool
	lineData    string
}

func TelnetConnect(network, addr string, connTimeout time.Duration, logger func(string)) (*MyTelnet, error) {
	if conn, err := telnet.Dial(network, addr); err != nil {
		return nil, err
	} else {
		return &MyTelnet{conn, network, addr, connTimeout, logger, true, ""}, nil
	}
}

func (t *MyTelnet) Reconnect() *MyTelnet {
	if conn, err := telnet.Dial(t.network, t.addr); err == nil {
		t.connected = true
		t.Conn = conn
	}
	return t
}

func (t *MyTelnet) Close() *MyTelnet {
	if t.connected {
		t.connected = false
		t.Close()
	}
	return t
}

func (t *MyTelnet) IsConnected() bool {
	return t.connected
}

func (t *MyTelnet) Expect(delim ...string) *MyTelnet {
	if t.connected && t.isSuccess(t.SetReadDeadline(time.Now().Add(t.connTimeout))) {
		t.isSuccess(t.SkipUntil(delim...))
	}
	return t
}

func (t *MyTelnet) SendLine(command string, args ...interface{}) *MyTelnet {
	if t.connected && t.isSuccess(t.SetWriteDeadline(time.Now().Add(t.connTimeout))) {
		_command := fmt.Sprintf(command+"\n", args...)
		buf := make([]byte, len(_command))
		copy(buf, _command)
		_, err := t.Write(buf)
		t.isSuccess(err)
	}
	return t
}

func (t *MyTelnet) ReadUntil(delim byte) *MyTelnet {
	if t.connected {
		if data, err := t.ReadBytes(delim); t.isSuccess(err) {
			t.lineData = string(data)
		}
	}
	return t
}

func (t *MyTelnet) FindAllStringSubmatch(re *regexp.Regexp) [][]string {
	if t.connected {
		return re.FindAllStringSubmatch(t.lineData, -1)
	}
	return [][]string{}
}

func (t *MyTelnet) FindAllString(re *regexp.Regexp) []string {
	if t.connected {
		return re.FindAllString(t.lineData, -1)
	}
	return []string{}
}

func (t *MyTelnet) isSuccess(err error) bool {
	if err != nil {
		t.logger(fmt.Sprintf("%s", err))
		t.connected = false
	}
	return t.connected
}
