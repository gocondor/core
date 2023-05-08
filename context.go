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
	"strings"

	"github.com/google/uuid"
)

type Context struct {
	Request  *Request
	Response *Response
	logger   *Logger
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
	ctx.Request.HttpRequest.ParseMultipartForm(20000000)
}

func (c *Context) LogInfo(msg interface{}) {
	ResolveLogger().Info(msg)
}

func (c *Context) LogError(msg interface{}) {
	ResolveLogger().Error(msg)
}

func (c *Context) LogWarning(msg interface{}) {
	ResolveLogger().Warning(msg)
}

func (c *Context) LogDebug(msg interface{}) {
	ResolveLogger().Debug(msg)
}

func (c *Context) GetPathParam(key string) string {
	return c.Request.httpPathParams.ByName(key)
}

func (c *Context) GetRequestParam(key string) string {
	return c.Request.HttpRequest.FormValue(key)
}

func (c *Context) RequestParamExists(key string) bool {
	return c.Request.HttpRequest.Form.Has(key)
}

func (c *Context) GetUploadedFile(name string) *UploadedFileInfo {
	file, fileHeader, _ := c.Request.HttpRequest.FormFile(name)
	defer file.Close()
	ext := strings.TrimPrefix(path.Ext(fileHeader.Filename), ".")
	tmpFilePath := filepath.Join(os.TempDir(), uuid.NewString()+"."+ext)
	tmpFile, _ := os.Create(tmpFilePath)
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
	tmpFileInfo, _ := os.Stat(tmpFilePath)
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
	destFilePath := filepath.Join(destFolderPath, srcFileInfo.Name())
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
