// Copyright 2021 Harran Ali <harran.m@gmail.com>. All rights reserved.
// Use of this source code is governed by MIT-style
// license that can be found in the LICENSE file.

package core

type EnvFile struct {
	UseDotEnvFile bool
}

type AppConfig struct {
	AppName             string
	AppEnv              string
	AppHttpHost         string
	AppHttpPort         string
	UseHttps            bool //-------//
	UseLetsEncrypt      bool
	LetsEncryptEmail    string
	HttpsHosts          string
	RedirectHttpToHttps string
	CertFilePath        string
	KeyFilePath         string
}

type RequestConfig struct {
	MaxUploadFileSize int
}

type JWTConfig struct {
	SecretKey string
	Lifetime  int
}
