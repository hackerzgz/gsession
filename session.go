package gsession // import "github.com/hackez/gsession"

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/hackez/gsession/ticker"
	"github.com/hackez/randstr"
)

type SessionPool struct {
	tokenLen int

	t        *ticker.TimeTicker
	sequence int

	m sync.Mutex
}

func New(m ticker.Method, l int) *SessionPool {
	return &SessionPool{
		tokenLen: l,
		t:        ticker.New(m),
		sequence: 0,
	}
}

func (sp *SessionPool) Generate(uid string) (Session, error) {
	token, err := randstr.NewString(randstr.CharDigit|randstr.CharLowerCase|randstr.CharUpperCase, sp.tokenLen)
	if err != nil {
		return "", err
	}

	// MD5(uid + token)
	buf := bytes.NewBufferString(uid)
	buf.WriteString(token)

	hashed := hash(buf.String())
	buf.Reset()
	buf.WriteString(hex.EncodeToString(hashed))

	// timestamp
	stamp := sp.t.Get()
	buf.WriteString(strconv.Itoa(int(stamp)))

	sp.m.Lock()
	if sp.sequence > 999 {
		sp.sequence = 0
	} else {
		sp.sequence++
	}

	seq := sp.sequence
	sp.m.Unlock()

	// sequence
	buf.WriteString(fmt.Sprintf("%04d", seq))

	if buf.Len() != 40 {
		return "", fmt.Errorf("something internal error inside Generate: %d", buf.Len())
	}

	return Session(buf.String()), nil
}

func (sp *SessionPool) Free() {
	sp.t.Stop()
	return
}

// Session a key mapping to a uid
// struct of session:
// | SHA256(uid+token) |  timestamp | sequence  |
// |     32 bits       |    5 bits  |  3 bits   |
type Session string

// LoginTime return format probable time the session was generated
func (s Session) LoginTime() (string, error) {
	if len(s) != 40 {
		return "", fmt.Errorf("not a valid session")
	}

	stime := string(s)[32:37]
	t, err := strconv.Atoi(stime)
	if err != nil {
		return "", err
	}

	if t != 0 {
		for {
			if t < ticker.StampTail {
				t *= 10
			} else {
				t /= 10
				break
			}
		}
	}

	now := time.Now().Unix()
	now /= ticker.StampTail
	now *= ticker.StampTail
	now += int64(t)

	return time.Unix(now, 0).Format(time.Stamp), nil
}
