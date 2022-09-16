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
		Time       string
		ErrDetails string
	}
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	je := JError{
		Result:     "Error",
		Time:       currentTime,
		ErrDetails: err.Error(),
	}
	res, e := json.Marshal(je)
	if e != nil {
		log.Println(e)
	}
	s := string(res)

	return s
}

// just return pretty "OK" in json
func ReturnOkJson(result string) string {
	type JError struct {
		Result string
		Time   string
	}
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	je := JError{
		Result: result,
		Time:   currentTime,
	}
	res, e := json.Marshal(je)
	if e != nil {
		log.Println(e)
	}
	s := string(res)

	return s
}
