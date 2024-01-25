// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: streaming_service.proto

package proto

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"google.golang.org/protobuf/types/known/anypb"
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
	_ = anypb.Any{}
	_ = sort.Sort
)

// Validate checks the field values on ReadFileRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *ReadFileRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ReadFileRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// ReadFileRequestMultiError, or nil if none found.
func (m *ReadFileRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *ReadFileRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if !_ReadFileRequest_Name_Pattern.MatchString(m.GetName()) {
		err := ReadFileRequestValidationError{
			field:  "Name",
			reason: "value does not match regex pattern \"^(/)?([^/\\x00]+(/)?)+$\"",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return ReadFileRequestMultiError(errors)
	}

	return nil
}

// ReadFileRequestMultiError is an error wrapping multiple validation errors
// returned by ReadFileRequest.ValidateAll() if the designated constraints
// aren't met.
type ReadFileRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ReadFileRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ReadFileRequestMultiError) AllErrors() []error { return m }

// ReadFileRequestValidationError is the validation error returned by
// ReadFileRequest.Validate if the designated constraints aren't met.
type ReadFileRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ReadFileRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ReadFileRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ReadFileRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ReadFileRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ReadFileRequestValidationError) ErrorName() string { return "ReadFileRequestValidationError" }

// Error satisfies the builtin error interface
func (e ReadFileRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sReadFileRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ReadFileRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ReadFileRequestValidationError{}

var _ReadFileRequest_Name_Pattern = regexp.MustCompile("^(/)?([^/\x00]+(/)?)+$")

// Validate checks the field values on ReadFileReply with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *ReadFileReply) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ReadFileReply with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in ReadFileReplyMultiError, or
// nil if none found.
func (m *ReadFileReply) ValidateAll() error {
	return m.validate(true)
}

func (m *ReadFileReply) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Stream

	if len(errors) > 0 {
		return ReadFileReplyMultiError(errors)
	}

	return nil
}

// ReadFileReplyMultiError is an error wrapping multiple validation errors
// returned by ReadFileReply.ValidateAll() if the designated constraints
// aren't met.
type ReadFileReplyMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ReadFileReplyMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ReadFileReplyMultiError) AllErrors() []error { return m }

// ReadFileReplyValidationError is the validation error returned by
// ReadFileReply.Validate if the designated constraints aren't met.
type ReadFileReplyValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ReadFileReplyValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ReadFileReplyValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ReadFileReplyValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ReadFileReplyValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ReadFileReplyValidationError) ErrorName() string { return "ReadFileReplyValidationError" }

// Error satisfies the builtin error interface
func (e ReadFileReplyValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sReadFileReply.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ReadFileReplyValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ReadFileReplyValidationError{}

// Validate checks the field values on LsRequest with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *LsRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on LsRequest with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in LsRequestMultiError, or nil
// if none found.
func (m *LsRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *LsRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if !_LsRequest_Dir_Pattern.MatchString(m.GetDir()) {
		err := LsRequestValidationError{
			field:  "Dir",
			reason: "value does not match regex pattern \"^(/)?([^/\\x00]+(/)?)+$\"",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return LsRequestMultiError(errors)
	}

	return nil
}

// LsRequestMultiError is an error wrapping multiple validation errors returned
// by LsRequest.ValidateAll() if the designated constraints aren't met.
type LsRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m LsRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m LsRequestMultiError) AllErrors() []error { return m }

// LsRequestValidationError is the validation error returned by
// LsRequest.Validate if the designated constraints aren't met.
type LsRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e LsRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e LsRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e LsRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e LsRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e LsRequestValidationError) ErrorName() string { return "LsRequestValidationError" }

// Error satisfies the builtin error interface
func (e LsRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sLsRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = LsRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = LsRequestValidationError{}

var _LsRequest_Dir_Pattern = regexp.MustCompile("^(/)?([^/\x00]+(/)?)+$")

// Validate checks the field values on LsReply with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *LsReply) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on LsReply with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in LsReplyMultiError, or nil if none found.
func (m *LsReply) ValidateAll() error {
	return m.validate(true)
}

func (m *LsReply) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if len(errors) > 0 {
		return LsReplyMultiError(errors)
	}

	return nil
}

// LsReplyMultiError is an error wrapping multiple validation errors returned
// by LsReply.ValidateAll() if the designated constraints aren't met.
type LsReplyMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m LsReplyMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m LsReplyMultiError) AllErrors() []error { return m }

// LsReplyValidationError is the validation error returned by LsReply.Validate
// if the designated constraints aren't met.
type LsReplyValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e LsReplyValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e LsReplyValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e LsReplyValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e LsReplyValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e LsReplyValidationError) ErrorName() string { return "LsReplyValidationError" }

// Error satisfies the builtin error interface
func (e LsReplyValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sLsReply.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = LsReplyValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = LsReplyValidationError{}

// Validate checks the field values on MetaRequest with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *MetaRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on MetaRequest with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in MetaRequestMultiError, or
// nil if none found.
func (m *MetaRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *MetaRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if !_MetaRequest_Name_Pattern.MatchString(m.GetName()) {
		err := MetaRequestValidationError{
			field:  "Name",
			reason: "value does not match regex pattern \"^(/)?([^/\\x00]+(/)?)+$\"",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return MetaRequestMultiError(errors)
	}

	return nil
}

// MetaRequestMultiError is an error wrapping multiple validation errors
// returned by MetaRequest.ValidateAll() if the designated constraints aren't met.
type MetaRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m MetaRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m MetaRequestMultiError) AllErrors() []error { return m }

// MetaRequestValidationError is the validation error returned by
// MetaRequest.Validate if the designated constraints aren't met.
type MetaRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e MetaRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e MetaRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e MetaRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e MetaRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e MetaRequestValidationError) ErrorName() string { return "MetaRequestValidationError" }

// Error satisfies the builtin error interface
func (e MetaRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sMetaRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = MetaRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = MetaRequestValidationError{}

var _MetaRequest_Name_Pattern = regexp.MustCompile("^(/)?([^/\x00]+(/)?)+$")

// Validate checks the field values on MetaReply with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *MetaReply) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on MetaReply with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in MetaReplyMultiError, or nil
// if none found.
func (m *MetaReply) ValidateAll() error {
	return m.validate(true)
}

func (m *MetaReply) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Size

	// no validation rules for Mode

	// no validation rules for IsDir

	if len(errors) > 0 {
		return MetaReplyMultiError(errors)
	}

	return nil
}

// MetaReplyMultiError is an error wrapping multiple validation errors returned
// by MetaReply.ValidateAll() if the designated constraints aren't met.
type MetaReplyMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m MetaReplyMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m MetaReplyMultiError) AllErrors() []error { return m }

// MetaReplyValidationError is the validation error returned by
// MetaReply.Validate if the designated constraints aren't met.
type MetaReplyValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e MetaReplyValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e MetaReplyValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e MetaReplyValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e MetaReplyValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e MetaReplyValidationError) ErrorName() string { return "MetaReplyValidationError" }

// Error satisfies the builtin error interface
func (e MetaReplyValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sMetaReply.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = MetaReplyValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = MetaReplyValidationError{}
