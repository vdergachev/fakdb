package main

import (
	"strconv"
	"strings"
)

type OperationCode uint8

const (
	OC_SET OperationCode = iota + 1 // set k v
	OC_GET                          // get k
	OC_DEL                          // del k
	OC_INC                          // inc k
	OC_DEC                          // dec k
)

type OperationResultCode uint8

const (
	ORC_OK OperationResultCode = iota + 1
	ORC_ERR
	ORC_INT
)

type OperationErrorCode uint8

const (
	OEC_OP_INVALID OperationErrorCode = iota + 1
	OEC_ARG_MISSING
	OEC_ARG_NUM_INVALID
)

type OperationSpec interface {
	Parse(buf []byte) error
	Execute() (OperationResult, error)
}

type Operation struct {
	Code   OperationCode
	Keys   []string
	Values []string
}

type OperationResult struct {
	Code OperationResultCode
	Body string // error only ???
}

type OpError struct {
	Code OperationErrorCode
	Text string
}

func NewOpError(code OperationErrorCode, text string) *OpError {
	return &OpError{Code: code, Text: text}
}

func (oe OpError) Error() string {
	return strconv.Itoa(int(oe.Code)) + " " + oe.Text
}

func (o *Operation) Parse(buf []byte) error {
	var val = strings.TrimRight(string(buf), "\n") // TODO: Use reader instead of manual trimming
	var raw = strings.Split(val, " ")              // TODO: Remove multiple spaces
	var op = raw[0]
	var args = raw[1:]

	if strings.EqualFold(op, "SET") { // TODO: Use cache string,uint8 instead of strcmp
		o.Code = OC_SET
		if err := o.setKeysValues(args); err != nil {
			return err
		}

	} else if strings.EqualFold(op, "GET") {
		o.Code = OC_GET
		if err := o.validateKeys(args); err != nil {
			return err
		}
		o.Keys = args
	} else if strings.EqualFold(op, "DEL") {
		o.Code = OC_DEL
		if err := o.validateKeys(args); err != nil {
			return err
		}
		o.Keys = args
	} else {
		return NewOpError(OEC_OP_INVALID, op+" is invalid operation")
	}

	return nil
}

func (o *Operation) setKeysValues(args []string) error {
	// len keys == 0
	if len(args) == 0 {
		return NewOpError(OEC_ARG_MISSING, "Arguments are missing")
	}
	// len keys = len values
	if len(args)%2 != 0 {
		return NewOpError(OEC_ARG_NUM_INVALID, "Not enough arguments")
	}

	for i, v := range args {
		if i%2 == 0 {
			o.Keys = append(o.Keys, v)
		} else {
			o.Values = append(o.Values, v)
		}
	}

	return nil
}

func (o *Operation) validateKeys(args []string) error {
	// len keys == 0
	if len(args) == 0 {
		return NewOpError(OEC_ARG_MISSING, "Arguments are missing")
	}
	return nil
}

func (o *Operation) Execute() (OperationResult, error) { // TODO: Remove method
	return OperationResult{
		Code: ORC_OK,
	}, nil
}
