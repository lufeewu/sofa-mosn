// Code generated by protoc-gen-validate
// source: envoy/config/filter/fault/v2/fault.proto
// DO NOT EDIT!!!

package v2

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gogo/protobuf/types"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = types.DynamicAny{}
)

// Validate checks the field values on FaultDelay with the rules defined in the
// proto definition for this message. If any rules are violated, an error is returned.
func (m *FaultDelay) Validate() error {
	if m == nil {
		return nil
	}

	if _, ok := FaultDelay_FaultDelayType_name[int32(m.GetType())]; !ok {
		return FaultDelayValidationError{
			Field:  "Type",
			Reason: "value must be one of the defined enum values",
		}
	}

	if m.GetPercent() > 100 {
		return FaultDelayValidationError{
			Field:  "Percent",
			Reason: "value must be less than or equal to 100",
		}
	}

	switch m.FaultDelaySecifier.(type) {

	case *FaultDelay_FixedDelay:

		if d := m.GetFixedDelay(); d != nil {
			dur := *d

			gt := time.Duration(0*time.Second + 0*time.Nanosecond)

			if dur <= gt {
				return FaultDelayValidationError{
					Field:  "FixedDelay",
					Reason: "value must be greater than 0s",
				}
			}

		}

	default:
		return FaultDelayValidationError{
			Field:  "FaultDelaySecifier",
			Reason: "value is required",
		}

	}

	return nil
}

// FaultDelayValidationError is the validation error returned by
// FaultDelay.Validate if the designated constraints aren't met.
type FaultDelayValidationError struct {
	Field  string
	Reason string
	Cause  error
	Key    bool
}

// Error satisfies the builtin error interface
func (e FaultDelayValidationError) Error() string {
	cause := ""
	if e.Cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.Cause)
	}

	key := ""
	if e.Key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sFaultDelay.%s: %s%s",
		key,
		e.Field,
		e.Reason,
		cause)
}

var _ error = FaultDelayValidationError{}
