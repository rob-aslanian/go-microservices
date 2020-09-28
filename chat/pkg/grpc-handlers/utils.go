package grpc_handlers

import (
	"fmt"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/chatRPC"
	"gitlab.lan/Rightnao-site/microservices/shared/hc-errors"
)

func recoverHandler(err *error) {
	if r := recover(); r != nil {
		switch e := r.(type) {
		case error:
			er := hc_errors.Handle(e)
			*err = er
		case string:
			er := hc_errors.JsonError{Type: hc_errors.GENERAL_ERROR_TYPE, Description: e}
			*err = er
		default:
			fmt.Println("recoverd withoud error: ", e)
		}
	}
}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}

var EMPTY = &chatRPC.Empty{}
