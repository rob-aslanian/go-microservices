package hc_errors

import (
	"github.com/globalsign/mgo"
	"strings"
)

func handleMongo(mgoErr *mgo.LastError) *JsonError {
	if mgoErr.Code == 11000 {
		var startKey = "index: "
		iStart := strings.Index(mgoErr.Err, startKey)
		var indexName string
		if iStart >= 0 {
			startTrimed := mgoErr.Err[iStart+len(startKey) : len(mgoErr.Err)]
			iEnd := strings.Index(startTrimed, " ")
			if iEnd >= 0 {
				indexName = startTrimed[0:iEnd]
			}
		}
		return DUPLICATE_KEY_ERROR.WithData("index", indexName)
	}

	return &JsonError{
		Type:        MONGO_ERROR_TYPE,
		Description: mgoErr.Err,
		Data: map[string]interface{}{
			"code": mgoErr.Code,
		},
	}
}
