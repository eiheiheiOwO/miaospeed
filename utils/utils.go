package utils

import (
	"github.com/gofrs/uuid"
	jsoniter "github.com/json-iterator/go"
	"os"
)

func RandomUUID() string {
	uid, _ := uuid.NewV4()
	return uid.String()
}

func ToJSON(a any) string {
	r, _ := jsoniter.MarshalToString(a)
	return r
}

func ReadFile(path string) string {
	if path == "" {
		return ""
	}
	file, err := os.ReadFile(path)
	if err != nil {
		DWarnf("cannot read the file: %s, err: %s", path, err.Error())
		return ""
	}
	return string(file)
}
