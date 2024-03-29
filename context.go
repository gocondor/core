// Copyright 2021 Harran Ali <harran.m@gmail.com>. All rights reserved.
// Use of this source code is governed by MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"syscall"

	"github.com/gocondor/core/logger"
	"gorm.io/gorm"
)

type Context struct {
	Request          *Request
	Response         *Response
	GetValidator     func() *Validator
	GetJWT           func() *JWT
	GetGorm          func() *gorm.DB
	GetCache         func() *Cache
	GetHashing       func() *Hashing
	GetMailer        func() *Mailer
	GetEventsManager func() *EventsManager
	GetLogger        func() *logger.Logger
}

// TODO enhance
func (c *Context) DebugAny(variable interface{}) {
	var formatted string
	m := reflect.ValueOf(variable)
	if m.Kind() == reflect.Pointer {
		formatted = fmt.Sprintf("\n\nType: [pointer] %v (%v) \nMemory Address: %v \nValue: %v\n\n", m.Type(), m.Elem().Kind(), m, m.Elem())
	} else {
		formatted = fmt.Sprintf("\n\nType: %v (%v) \nValue: %v\n\n", m.Type(), m.Kind(), variable)
	}
	fmt.Println(formatted)
	c.Response.HttpResponseWriter.Write([]byte(formatted))
}

func (c *Context) Next() {
	ResolveApp().Next(c)
}

func (c *Context) prepare(ctx *Context) {
	ctx.Request.httpRequest.ParseMultipartForm(int64(app.Config.Request.MaxUploadFileSize))
}

func (c *Context) GetPathParam(key string) interface{} {
	return c.Request.httpPathParams.ByName(key)
}

func (c *Context) GetRequestParam(key string) interface{} {
	return c.Request.httpRequest.FormValue(key)
}

func (c *Context) RequestParamExists(key string) bool {
	return c.Request.httpRequest.Form.Has(key)
}

func (c *Context) GetHeader(key string) string {
	return c.Request.httpRequest.Header.Get(key)
}

func (c *Context) GetUploadedFile(name string) *UploadedFileInfo {
	file, fileHeader, err := c.Request.httpRequest.FormFile(name)
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

func (c *Context) MoveFile(sourceFilePath string, destFolderPath string) error {
	o := syscall.Umask(0)
	defer syscall.Umask(o)
	newFileName := filepath.Base(sourceFilePath)
	os.MkdirAll(destFolderPath, 766)
	srcFileInfo, err := os.Stat(sourceFilePath)
	if err != nil {
		return err
	}
	if !srcFileInfo.Mode().IsRegular() {
		return errors.New("can not move file, not in a regular mode")
	}
	srcFile, err := os.Open(sourceFilePath)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	destFilePath := filepath.Join(destFolderPath, newFileName)
	destFile, err := os.Create(destFilePath)
	if err != nil {
		return err
	}
	buff := make([]byte, 1024*8)
	for {
		n, err := srcFile.Read(buff)
		if err != nil && err != io.EOF {
			panic(fmt.Sprintf("error moving file %v", sourceFilePath))
		}
		if n == 0 {
			break
		}
		_, err = destFile.Write(buff[:n])
		if err != nil {
			return err
		}
	}
	destFile.Close()
	err = os.Remove(sourceFilePath)
	if err != nil {
		return err
	}
	return nil
}

func (c *Context) CopyFile(sourceFilePath string, destFolderPath string) error {
	o := syscall.Umask(0)
	defer syscall.Umask(o)
	newFileName := filepath.Base(sourceFilePath)
	os.MkdirAll(destFolderPath, 766)
	srcFileInfo, err := os.Stat(sourceFilePath)
	if err != nil {
		return err
	}
	if !srcFileInfo.Mode().IsRegular() {
		return errors.New("can not move file, not in a regular mode")
	}
	srcFile, err := os.Open(sourceFilePath)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	destFilePath := filepath.Join(destFolderPath, newFileName)
	destFile, err := os.Create(destFilePath)
	if err != nil {
		return err
	}
	buff := make([]byte, 1024*8)
	for {
		n, err := srcFile.Read(buff)
		if err != nil && err != io.EOF {
			panic(fmt.Sprintf("error moving file %v", sourceFilePath))
		}
		if n == 0 {
			break
		}
		_, err = destFile.Write(buff[:n])
		if err != nil {
			return err
		}
	}
	destFile.Close()

	return nil
}

func (c *Context) MapToJson(v any) string {
	r := reflect.ValueOf(v)
	if r.Kind() != reflect.Map {
		panic("parameter is not a map")
	}
	j, err := json.Marshal(v)
	if err != nil {
		panic(err.Error())
	}
	return string(j)
}

type UploadedFileInfo struct {
	FullPath             string
	Name                 string
	NameWithoutExtension string
	Extension            string
	Size                 int
}

func (c *Context) GetBaseDirPath() string {
	return basePath
}

func (c *Context) CastToString(value interface{}) string {
	if !basicType(value) {
		panic("can not cast to string")
	}

	return fmt.Sprintf("%v", value)
}

func (c Context) GetUserAgent() string {
	return c.Request.httpRequest.UserAgent()
}

func (c *Context) CastToInt(value interface{}) int {
	var i int
	if !basicType(value) {
		panic("can not cast to int")
	}
	i, ok := value.(int)
	if ok {
		return i
	}
	_i, ok := value.(int32)
	if ok {
		i := int(_i)
		return i
	}
	_ii, ok := value.(int64)
	if ok {
		i := int(_ii)
		return i
	}
	f, ok := value.(float32)
	if ok {
		i := int(f)
		return i
	}
	ff, ok := value.(float64)
	if ok {
		i := int(ff)
		return i
	}
	s, ok := value.(string)
	if ok {
		fff, err := strconv.ParseFloat(s, 64)
		if err != nil {
			panic("error casting to int")
		}
		i = int(fff)
		return i
	}
	panic("error casting to int")
}

func (c *Context) CastToFloat(value interface{}) float64 {
	if !basicType(value) {
		panic("can not cast to float")
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
