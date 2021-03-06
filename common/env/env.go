package env

import (
	"os"

	log "github.com/sirupsen/logrus"
)

const (
	KeyDomain  = "TUFIN_DOMAIN"
	KeyProject = "TUFIN_PROJECT"
)

func GetEnvWithDefault(variable, defaultValue string) string {

	ret := os.Getenv(variable)
	if ret == "" {
		ret = defaultValue
	}
	log.Infof("'%s': '%s'", variable, ret)

	return ret
}

func GetEnvOrExit(variable string) string {

	ret := os.Getenv(variable)
	if ret == "" {
		log.Fatalf("Please, set '%s'", variable)
	}
	log.Infof("'%s': '%s'", variable, ret)

	return ret
}

func GetEnvSensitive(variable string) string {

	ret := os.Getenv(variable)
	if ret != "" {
		log.Infof("'%s': [sensitive]", variable)
	}

	return ret
}

func GetEnvSensitiveOrExit(variable string) string {

	ret := GetEnvSensitive(variable)
	if ret == "" {
		log.Fatalf("Please, set '%s'", variable)
	}

	return ret
}
