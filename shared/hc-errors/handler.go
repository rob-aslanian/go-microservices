package hc_errors

import (
	"github.com/globalsign/mgo"
	"reflect"
)

func Handle(err error) error {
	switch e := err.(type) {
	case JsonError:
		return err
	case *JsonError:
		return err

	case *mgo.LastError:
		return handleMongo(e)

	}

	if reflect.TypeOf(err).String() == "*status.statusError" || reflect.TypeOf(err).String() == "status.statusError" {
		jsonErr, ok := UnwrapJsonErrorFromRPCError(err)
		if ok {
			return jsonErr
		}
		return err
	}


	return NewGeneralError(err)
}
