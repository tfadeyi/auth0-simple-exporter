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
//go:generate swag init --parseDependency --generalInfo ./pkg/exporter/server.go --output ./pkg/docs --markdownFiles ./docs
//go:generate sloscribe init --to-file
//go:generate sloth generate -i ./slo_definitions/auth0-exporter.yaml -o example_rules.yml

package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/tfadeyi/auth0-simple-exporter/cmd"
	"github.com/tfadeyi/auth0-simple-exporter/pkg/logging"
)

// @sloth service auth0-exporter

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM)
	defer cancel()

	log := logging.NewPromLogger()
	ctx = logging.ContextWithLogger(ctx, log)

	cmd.Execute(ctx)
}
