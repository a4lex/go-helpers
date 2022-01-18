package wrapers

import (
	"flag"
	"fmt"
	"regexp"
	"time"

	telnet "github.com/a4lex/go-telnet"
)

var (
	TelnetTimeout = flag.Int("telnet-timeout", 30, "Telnet Timeout for waiting responce from device")
	TelnetRetries = flag.Int("telnet-retries", 3, "Telnet Retries connect to device")
)

type MyTelnet struct {
	*telnet.Conn
	connTimeout   time.Duration
	logger        func(string)
	comChainState bool
	lineData      string
}

func TelnetConnecct(network, addr string, connTimeout time.Duration, logger func(string)) (*MyTelnet, error) {
	if t, err := telnet.Dial(network, addr); err != nil {
		return nil, err
	} else {
		return &MyTelnet{t, connTimeout, logger, false, ""}, nil
	}
}

func (t *MyTelnet) GetCommandChainState() bool {
	return t.comChainState
}

func (t *MyTelnet) ResetCommandChainState() *MyTelnet {
	t.comChainState = true
	return t
}

func (t *MyTelnet) Expect(delim ...string) *MyTelnet {
	if t.comChainState && t.isSuccess(t.SetReadDeadline(time.Now().Add(t.connTimeout))) {
		t.isSuccess(t.SkipUntil(delim...))
	}
	return t
}

func (t *MyTelnet) SendLine(command string, args ...interface{}) *MyTelnet {
	if t.comChainState && t.isSuccess(t.SetWriteDeadline(time.Now().Add(t.connTimeout))) {
		_command := fmt.Sprintf(command+"\n", args...)
		buf := make([]byte, len(_command))
		copy(buf, _command)
		_, err := t.Write(buf)
		t.isSuccess(err)
	}
	return t
}

func (t *MyTelnet) ReadUntil(delim byte) *MyTelnet {
	if t.comChainState {
		if data, err := t.ReadBytes(delim); t.isSuccess(err) {
			t.lineData = string(data)
		}
	}
	return t
}

func (t *MyTelnet) FindAllStringSubmatch(re *regexp.Regexp) [][]string {
	if t.comChainState {
		return re.FindAllStringSubmatch(t.lineData, -1)
	}
	return [][]string{}
}

func (t *MyTelnet) FindAllString(re *regexp.Regexp) []string {
	if t.comChainState {
		return re.FindAllString(t.lineData, -1)
	}
	return []string{}
}

func (t *MyTelnet) isSuccess(err error) bool {
	if err != nil {
		t.logger(fmt.Sprintf("%s", err))
		t.comChainState = false
	}
	return t.comChainState
}
