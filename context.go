// Copyright 2021 Harran Ali <harran.m@gmail.com>. All rights reserved.
// Use of this source code is governed by MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/gocondor/core/logger"
)

type Context struct {
	Request   *Request
	Response  *Response
	logger    *logger.Logger
	Validator *Validator
}

func (c *Context) DebugAny(variable interface{}) {
	m := reflect.ValueOf(variable)
	if m.Kind() == reflect.Pointer {
		m = m.Elem()
	}
	formatted := fmt.Sprintf("Type: %T (%v) | value: %v", variable, m.Kind(), variable)
	fmt.Println(formatted)
	c.Response.HttpResponseWriter.Write([]byte(formatted))
}

func (c *Context) Next() {
	ResolveApp().Next(c)
}

func (c *Context) prepare(ctx *Context) {
	ctx.Request.HttpRequest.ParseMultipartForm(int64(app.Config.Request.MaxUploadFileSize))
}

func (c *Context) LogInfo(msg interface{}) {
	logger.ResolveLogger().Info(msg)
}

func (c *Context) LogError(msg interface{}) {
	logger.ResolveLogger().Error(msg)
}

func (c *Context) LogWarning(msg interface{}) {
	logger.ResolveLogger().Warning(msg)
}

func (c *Context) LogDebug(msg interface{}) {
	logger.ResolveLogger().Debug(msg)
}

func (c *Context) GetPathParam(key string) interface{} {
	return c.Request.httpPathParams.ByName(key)
}

func (c *Context) GetRequestParam(key string) interface{} {
	return c.Request.HttpRequest.FormValue(key)
}

func (c *Context) RequestParamExists(key string) interface{} {
	return c.Request.HttpRequest.Form.Has(key)
}

func (c *Context) GetHeader(key string) string {
	return c.Request.HttpRequest.Header.Get(key)
}

func (c *Context) GetUploadedFile(name string) *UploadedFileInfo {
	file, fileHeader, err := c.Request.HttpRequest.FormFile(name)
	if err != nil {
		panic(fmt.Sprintf("error with file,[%v]", err.Error()))
	}
	defer file.Close()
	ext := strings.TrimPrefix(path.Ext(fileHeader.Filename), ".")
	tmpFilePath := filepath.Join(os.TempDir(), fileHeader.Filename)
	tmpFile, err := os.Create(tmpFilePath)
	if err != nil {
		panic(fmt.Sprintf("error with file,[%v]", err.Error()))
	}
	buff := make([]byte, 100)
	for {
		n, err := file.Read(buff)
		if err != nil && err != io.EOF {
			panic("error with uploaded file")
		}
		if n == 0 {
			break
		}
		n, _ = tmpFile.Write(buff[:n])
	}
	tmpFileInfo, err := os.Stat(tmpFilePath)
	if err != nil {
		panic(fmt.Sprintf("error with file,[%v]", err.Error()))
	}
	defer tmpFile.Close()
	uploadedFileInfo := &UploadedFileInfo{
		FullPath:             tmpFilePath,
		Name:                 fileHeader.Filename,
		NameWithoutExtension: strings.TrimSuffix(fileHeader.Filename, path.Ext(fileHeader.Filename)),
		Extension:            ext,
		Size:                 int(tmpFileInfo.Size()),
	}
	return uploadedFileInfo
}

func (c *Context) MoveFile(sourceFilePath string, destFolderPath string, newFileName string) error {
	os.MkdirAll(destFolderPath, 644)
	srcFileInfo, err := os.Stat(sourceFilePath)
	if err != nil {
		return err
	}
	if !srcFileInfo.Mode().IsRegular() {
		return errors.New("can not copy file, not in a regular mode")
	}
	srcFile, err := os.Open(sourceFilePath)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	destFilePath := filepath.Join(destFolderPath, newFileName)
	fmt.Println(srcFileInfo.Name())
	destFile, err := os.OpenFile(destFilePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 744)
	if err != nil {
		return err
	}
	defer destFile.Close()
	buff := make([]byte, 100)
	for {
		n, err := srcFile.Read(buff)
		if err != nil && err != io.EOF {
			panic("error moving file")
		}
		if n == 0 {
			break
		}
		_, err = destFile.Write(buff)
		if err != nil {
			return err
		}
	}
	return nil
}

type UploadedFileInfo struct {
	FullPath             string
	Name                 string
	NameWithoutExtension string
	Extension            string
	Size                 int
}

func (c *Context) GetBaseDirPath() string {
	wd, _ := os.Getwd()
	return wd
}

func (c *Context) CastToString(value interface{}) string {
	if !basicType(value) {
		panic("can not cast to string")
	}
	return fmt.Sprintf("%s", value)
}

func (c *Context) CastToInt(value interface{}) int {
	if !basicType(value) {
		panic("can not cast to int")
	}
	s, ok := value.(string)
	if !ok {
		panic("error casting to int")
	}
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic("error casting to int")
	}
	iint := int(i)
	return iint
}

func (c *Context) CastToFloat(value interface{}) float64 {
	if !basicType(value) {
		panic("can not cast to int")
	}
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}
	var str string
	var ok bool
	if v.Kind() == reflect.Float64 {
		f, ok := value.(float64)
		if !ok {
			panic("error casting to float")
		}
		return f
	}
	if v.Kind() == reflect.Float32 {
		s := fmt.Sprintf("%v", value)
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			panic("error casting to float")
		}
		return f
	}
	if v.Kind() == reflect.String {
		str, ok = value.(string)
		if !ok {
			panic("error casting to float")
		}
	}
	if v.CanInt() {
		i, ok := value.(int)
		if !ok {
			panic("error casting to float")
		}
		str = fmt.Sprintf("%v.0", i)
	}
	f, err := strconv.ParseFloat(str, 64)
	if err != nil {
		panic("error casting to float")
	}
	return f
}

func basicType(value interface{}) bool {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}
	if !(v.Kind() == reflect.Array ||
		v.Kind() == reflect.Slice ||
		v.Kind() == reflect.Map ||
		v.Kind() == reflect.Struct ||
		v.Kind() == reflect.Interface ||
		v.Kind() == reflect.Func) {
		return true
	}
	return false
}
