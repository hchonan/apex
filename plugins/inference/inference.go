// Package inference adds file-based inference to detect runtimes when they are not explicitly specified.
package inference

import (
	"os"
	"path/filepath"

	"github.com/hchonan/apex/function"
	"github.com/hchonan/apex/plugins/golang"
	"github.com/hchonan/apex/plugins/java"
	"github.com/hchonan/apex/plugins/nodejs"
	"github.com/hchonan/apex/plugins/python"
	"github.com/hchonan/apex/plugins/ruby"
)

func init() {
	function.RegisterPlugin("inference", &Plugin{
		Files: map[string]string{
			"main.py":             python.Runtime,
			"index.js":            nodejs.Runtime,
			"main.go":             golang.Runtime,
			"target/apex.jar":     java.Runtime,
			"build/libs/apex.jar": java.Runtime,
			"lambda.rb":           ruby.Runtime,
		},
	})
}

// Plugin implementation.
type Plugin struct {
	Files map[string]string
}

// Open checks for files in the function directory to infer its runtime.
func (p *Plugin) Open(fn *function.Function) error {
	if fn.Runtime != "" {
		return nil
	}

	fn.Log.Debug("inferring runtime")

	for name, runtime := range p.Files {
		if _, err := os.Stat(filepath.Join(fn.Path, name)); err == nil {
			fn.Log.WithField("runtime", runtime).Debug("inferred runtime")
			fn.Runtime = runtime
			return nil
		}
	}

	return nil
}
