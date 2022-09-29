/*
 Copyright 2022 The CloudEvents Authors
 SPDX-License-Identifier: Apache-2.0
*/

package model

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"

	"gopkg.in/yaml.v3"
)

func FromYaml(recursive bool, files ...string) ([]Features, error) {
	feats := make([]Features, 0)
	for _, pathName := range files {
		var feat []Features
		var err error

		switch {
		case pathName == "-":
			feat, err = decode(os.Stdin)
		case isURL(pathName):
			feat, err = readURL(pathName)
		default:
			feat, err = readPath(pathName, recursive)
		}

		if err != nil {
			return nil, err
		}
		feats = append(feats, feat...)
	}
	return feats, nil
}

func isURL(pathname string) bool {
	if _, err := os.Lstat(pathname); err == nil {
		return false
	}
	uri, err := url.ParseRequestURI(pathname)
	return err == nil && uri.Scheme != ""
}

func readURL(uri string) ([]Features, error) {
	resp, err := http.Get(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return decode(resp.Body)
}

func readPath(pathName string, recursive bool) ([]Features, error) {
	info, err := os.Stat(pathName)
	if err != nil {
		return nil, err
	}

	if info.IsDir() {
		return readDir(pathName, recursive)
	}
	return readFile(pathName)
}

func readFile(pathName string) ([]Features, error) {
	file, err := os.Open(pathName)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	return decode(file)
}

func readDir(pathName string, recursive bool) ([]Features, error) {
	list, err := ioutil.ReadDir(pathName)
	if err != nil {
		return nil, err
	}

	feats := make([]Features, 0)
	for _, f := range list {
		name := path.Join(pathName, f.Name())
		var fs []Features

		switch {
		case f.IsDir() && recursive:
			fs, err = readDir(name, recursive)
		case !f.IsDir():
			fs, err = readFile(name)
		}

		if err != nil {
			return nil, err
		}
		feats = append(feats, fs...)
	}
	return feats, nil
}

func decode(reader io.Reader) ([]Features, error) {
	decoder := yaml.NewDecoder(reader)
	events := make([]Features, 0)
	var err error
	for {
		out := Features{}
		err = decoder.Decode(&out)
		if err != nil {
			break
		}
		events = append(events, out)
	}
	if err != io.EOF {
		return nil, err
	}
	return events, nil
}
