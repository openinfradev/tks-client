package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"runtime/debug"

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
	json.NewEncoder(buf).Encode(in)
	json.NewDecoder(buf).Decode(out)
}

func CheckError(err error) {
	if err != nil {
		fmt.Println(err)
		debug.PrintStack()
		os.Exit(-1)
	}
}
