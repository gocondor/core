package core

import "github.com/gocondor/core/env"

var testingAppC = AppConfig{
	AppENV:        env.GetVarOtherwiseDefault("APP_ENV", "testing"), // local | production | testing
	UseDotEnvFile: true,
}

var testingRequestC = RequestConfig{
	MaxUploadFileSize: 20000000, // 20MB
}
