// Copyright 2021 Harran Ali <harran.m@gmail.com>. All rights reserved.
// Use of this source code is governed by MIT-style
// license that can be found in the LICENSE file.

package core

type EnvFileConfig struct {
	UseDotEnvFile bool
}

type RequestConfig struct {
	MaxUploadFileSize int
}

type JWTConfig struct {
	SecretKey string
	Lifetime  int
}

type GormConfig struct {
	EnableGorm bool
}

type CacheConfig struct {
	EnableCache bool
}
