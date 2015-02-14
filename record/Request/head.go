package request

import (
	"fmt"
	"github.com/cgentry/gus/ecode"
	"github.com/cgentry/gus/record/stamp"
	"github.com/cgentry/gus/record/signature"
	"time"
)

const (
	TIMESTAMP_EXPIRATION = 2 * 60
)

var unixTimeZero = time.Unix(0, 0)


// Head implements the record.HeaderInterface
type Head struct {
	*stamp.Timestamp
	Domain string
	Id     string
	Sequence int
	*signature.Signature
}

func NewHead() Head {
	h := new(Head)
	h.Timestamp = stamp.New()
	h.Signature = signature.New()
	h.Signature.SetSignature("")
	return *h
}

func (h Head) Check() error {

	if h.Domain == "" {
		return ecode.ErrHeadNoDomain
	}
	if h.Id == "" {
		return ecode.ErrHeadNoId
	}
	if !h.IsTimeSet() {
		return ecode.ErrHeadNoTimestamp
	}

	// Note: stale time is always 2 minutes old. You can check for earlier times...
	window := h.Window(TIMESTAMP_EXPIRATION)
	if window != 0 {
		if window > 0 {
			return ecode.ErrHeadFuture
		}
		if window < 0 {
			return ecode.ErrHeadExpired
		}
	}
	return nil
}

func (h Head) String() string {
	return fmt.Sprintf("Domain: '%s', Id: '%s', Time: '%v', Signature: '%s', Set? %v",
		h.Domain, h.Id, h.Timestamp, h.Signature.GetSignature(), h.IsTimeSet())
}


