/*
*  Copyright (c) WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
*
*  WSO2 Inc. licenses this file to you under the Apache License,
*  Version 2.0 (the "License"); you may not use this file except
*  in compliance with the License.
*  You may obtain a copy of the License at
*
*    http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing,
* software distributed under the License is distributed on an
* "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
* KIND, either express or implied.  See the License for the
* specific language governing permissions and limitations
* under the License.
 */

package formatter

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
	"text/template"

	"github.com/wso2/product-apim-tooling/import-export-cli/templates"
)

// TableFormatKey is the identifier used for table
const TableFormatKey = "table"

// DetailedFormatKey is the identifier used for detail view
const DetailedFormatKey = "detail"

// Format is an alias for a string used for formatting
type Format string

// IsTable returns true if format string is prefixed with table
func (f Format) IsTable() bool {
	return strings.HasPrefix(string(f), TableFormatKey)
}

// IsDetailedFormat returns true if format string is prefixed with detail
func (f Format) IsDetailedFormat() bool {
	return strings.HasPrefix(string(f), DetailedFormatKey)
}

// Context keeps data about a format operation
type Context struct {
	// Output is used to write the output
	Output io.Writer
	// Format is used to keep format string
	Format Format

	// internal usage
	finalFormat string
	buffer      *bytes.Buffer
}

// NewContext creates a context with initialized fields
func NewContext(output io.Writer, format string) *Context {
	return &Context{Output: output, Format: Format(format), buffer: &bytes.Buffer{}}
}

// preFormat will clean format string
func (ctx *Context) preFormat() {
	format := string(ctx.Format)

	if ctx.Format.IsTable() {
		// if table is found skip it and take the rest
		format = format[len(TableFormatKey):]
	}

	if ctx.Format.IsDetailedFormat() {
		// if detail is found skip it and take the rest
		format = format[len(DetailedFormatKey):]
	}

	format = strings.TrimSpace(format)
	// this is done to avoid treating \t \n as template strings. This replaces them as special characters
	replacer := strings.NewReplacer(`\t`, "\t", `\n`, "\n")
	format = replacer.Replace(format)
	ctx.finalFormat = format
}

// parseTemplate will create a new template with basic functions
func (ctx Context) parseTemplate() (*template.Template, error) {
	tmpl, err := templates.NewBasicFormatter("").Parse(ctx.finalFormat)
	if err != nil {
		return tmpl, fmt.Errorf("Template parsing error: %v\n", err)
	}
	return tmpl, nil
}

// postFormat will output to writer
func (ctx *Context) postFormat(template *template.Template, headers interface{}) {
	if ctx.Format.IsTable() {
		// create a tab writer using Output
		w := tabwriter.NewWriter(ctx.Output, 20, 1, 3, ' ', 0)
		// print headers
		_ = template.Funcs(templates.HeaderFuncs).Execute(w, headers)
		_, _ = w.Write([]byte{'\n'})
		// write buffer to the w
		// in this case anything in buffer will be rendered by tabwiter to the Output
		// buffer contains data to be written
		_, _ = ctx.buffer.WriteTo(w)
		// flush will perform actual write to the writer
		_ = w.Flush()
	} else if ctx.Format.IsDetailedFormat() {
		// create a tab writer using Output
		w := tabwriter.NewWriter(ctx.Output, 20, 1, 3, ' ', 0)
		// write buffer to the w
		// in this case anything in buffer will be rendered by tabwiter to the Output
		// buffer contains data to be written
		_, _ = ctx.buffer.WriteTo(w)
		// flush will perform actual write to the writer
		_ = w.Flush()
	} else {
		// just write it as normal
		_, _ = ctx.buffer.WriteTo(ctx.Output)
	}
}

// Renderer is used to render a particular resource using templates
type Renderer func(io.Writer, *template.Template) error

// Write writes data using r and headers
func (ctx *Context) Write(r Renderer, headers interface{}) error {
	// prepare formatting
	ctx.preFormat()
	// parse template
	tmpl, err := ctx.parseTemplate()
	if err != nil {
		return err
	}
	// using renderer provided render collection
	// Note: See the renderer implementation in cmd/apis.go for more
	if err = r(ctx.buffer, tmpl); err != nil {
		return err
	}
	// write results to writer
	ctx.postFormat(tmpl, headers)
	return nil
}
