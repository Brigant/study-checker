package helpers

import (
	"encoding/json"
	"log"
	"time"
)

// Just convert error to json.
func ReturnErrorJSON(err error) string {
	type JError struct {
		Result     string
		Time       string
		ErrDetails string
	}

	currentTime := time.Now().Format("2006-01-02 15:04:05")

	jerror := JError{
		Result:     "Error",
		Time:       currentTime,
		ErrDetails: err.Error(),
	}

	res, e := json.Marshal(jerror)
	if e != nil {
		log.Println(e)
	}

	s := string(res)

	return s
}

// Just return pretty "OK" in json.
func ReturnOkJSON(result string) string {
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
