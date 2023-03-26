/*
Copyright © 2023 Oluwole Fadeyi (oluwolefadeyi@gmail.com)

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
//go:generate swag fmt
//go:generate swag init --parseDependency

package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/auth0-simple-exporter/cmd"
	"github.com/labstack/gommon/log"
)

// @title       Auth0 simple exporter
// @version     0.0.7
// @description A simple Prometheus exporter for Auth0 log events,(https://auth0.com/docs/api/management/v2#!/Logs/get_logs), for a simple
//way to monitor Auth0 from a Prometheus monitoring stack.

// @contact.name Oluwole Fadeyi (@tfadeyi)

// @license.name Apache 2.0
// @license.url  https://github.com/tfadeyi/auth0-simple-exporter/blob/main/LICENSE

// @host     localhost:9301
// @BasePath /
func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM)
	defer cancel()

	log.SetPrefix("auth0-exporter")

	cmd.Execute(ctx)
}
