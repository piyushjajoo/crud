/*
Copyright © 2021 Piyush Jajoo piyush.jajoo1991@gmail.com

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package tpl

// MainTemplate returns the template for main.go
func MainTemplate() []byte {
	return []byte(`package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"{{ .ModuleName }}/pkg/conf"
	"{{ .ModuleName }}/pkg/routes"
	"{{ .ModuleName }}/pkg/utils"

	"github.com/gorilla/mux"
)

func init() {
	utils.LoadEnvConfig(&conf.Env)
}

func main() {
	
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	// create a router
	r := mux.NewRouter()

	routes.Routes(r)
	
	// start the server
	srv := &http.Server{
		Addr: "0.0.0.0:8080",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	log.Println("http server started at port 8080")

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}

`)
}

// RoutesTemplate returns template for pkg/routes/routes.go
func RoutesTemplate() []byte {
	return []byte(`package routes

import (
	"github.com/gorilla/mux"
)

func Routes(r *mux.Router) {

}

`)
}

// ConstsTemplate returns template for pkg/consts/consts.go
func ConstsTemplate() []byte {
	return []byte(`package consts

const ()

`)
}

// ConfTemplate returns template for pkg/conf/conf.go
func ConfTemplate() []byte {
	return []byte(`package conf

// EnvConfig stores env vars
type EnvConfig struct {
}

// Env stores env vars
var Env EnvConfig

`)
}

// UtilsTemplate returns template for pkg/utils/utils.go
func UtilsTemplate() []byte {
	return []byte(`package utils

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/go-playground/validator"
)

// LoadEnvConfig loads the env vars into the provided struct and validates them based on tags
func LoadEnvConfig(spec interface{}) error {
	err := envconfig.Process("", spec)
	if err != nil {
		return err
	}
	err = validator.New().Struct(spec)
	if err != nil {
		return err
	}
	return nil
}

`)
}

// ModelsTemplate returns template for pkg/models/models.go
func ModelsTemplate() []byte {
	return []byte(`package models
`)
}

// DockerfileTemplate returns template for Dockerfile
func DockerfileTemplate() []byte {
	return []byte(`FROM debian
COPY ./{{ .ProjectDirName }} /{{ .ProjectDirName }}
ENTRYPOINT [ "/{{ .ProjectDirName }}" ]
`)
}

// BuildFileTemplate returns template for build.sh file
func BuildFileTemplate() []byte {
	return []byte(`#!/bin/sh

GOOS=linux go build .
docker build -t {{ .ProjectDirName }} .
`)
}