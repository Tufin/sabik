package client

import "github.com/tufin/sabik/common/env"

func CreateMiddleware() Middleware {

	if isEnable() {
		return newSabikMiddleware()
	}

	return &EmptyMiddleware{}
}

func isEnable() bool {

	return env.GetEnvWithDefault(EnvKeyEnable, "false") == "true"
}
