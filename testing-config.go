package core

import "github.com/gocondor/core/env"

var testingAppC = AppConfig{
	AppEnv: env.GetVarOtherwiseDefault("APP_ENV", "testing"), // local | production | testing
}

var testingRequestC = RequestConfig{
	MaxUploadFileSize: 20000000, // 20MB
}
