package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"runtime/debug"
	"time"

	"github.com/google/uuid"
)

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func NilUUID() uuid.UUID {
	nilId, _ := uuid.Parse("")
	return nilId
}

func ModelToJson(in any) string {
	a, _ := json.Marshal(in)
	n := len(a)        //Find the length of the byte array
	s := string(a[:n]) //convert to string
	return s
}

func Transcode(in, out interface{}) {
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(in)
	if err != nil {
		fmt.Println(err)
	}
	err = json.NewDecoder(buf).Decode(out)
	if err != nil {
		fmt.Println(err)
	}
}

func CheckError(err error) {
	if err != nil {
		fmt.Println(err)
		debug.PrintStack()
		os.Exit(-1)
	}
}

func PanicWithError(message string) {
	fmt.Println(message)
	os.Exit(-1)
}

func ParseTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}
