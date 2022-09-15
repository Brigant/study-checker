package helpers

import (
	"encoding/json"
	"log"
	"time"
)

// TODO: implement global helper func, which will be able to recive writer pointer and write
// in it status code and writes json answer (error or some stucture)

// just convert error to json
func ReturnErrorJson(err error) string {
	type JError struct {
		Result     string
		ErrorTime  string
		ErrDetails string
	}
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	je := JError{
		Result:     "Error",
		ErrorTime:  currentTime,
		ErrDetails: err.Error(),
	}
	res, e := json.Marshal(je)
	if e != nil {
		log.Println(e)
	}
	s := string(res)

	return s
}
