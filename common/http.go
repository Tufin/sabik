package common

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func StringHeader(header http.Header) (string, error) {

	ret, err := json.Marshal(header)
	if err != nil {
		log.Errorf("failed to marshal HTTP header with '%v'", err)
		return "", err
	}

	return string(ret), nil
}
