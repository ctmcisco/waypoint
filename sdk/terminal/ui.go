package terminal

import (
	"fmt"
	"io"

	"github.com/fatih/color"
)

// UI is the primary interface for interacting with a user via the CLI.
//
// NOTE(mitchellh): This is an interface and not a struct directly so that
// we can support other user interaction patterns in the future more easily.
// Most importantly what I'm thinking of is when we support multiple "apps"
// in a single config file, we can build a UI that locks properly and so on
// without changing the API.
type UI interface {
	// Output outputs a message directly to the terminal. The remaining
	// arguments should be interpolations for the format string. After the
	// interpolations you may add Options.
	Output(string, ...interface{})

	// Status returns a live-updating status that can be used for single-line
	// status updates that typically have a spinner or some similar style.
	// While a Status is live (Close isn't called), Output should NOT be called.
	Status() Status
}

// BasicUI
type BasicUI struct{}

// Output implements UI
func (ui *BasicUI) Output(msg string, raw ...interface{}) {
	// Build our args and options
	var args []interface{}
	var opts []Option
	for _, r := range raw {
		if opt, ok := r.(Option); ok {
			opts = append(opts, opt)
		} else {
			args = append(args, r)
		}
	}

	// Build our message
	msg = fmt.Sprintf(msg, args...)

	// Build our config and set our options
	cfg := &config{Original: msg, Message: msg, Writer: color.Output}
	for _, opt := range opts {
		opt(cfg)
	}

	// Write it
	fmt.Fprintln(cfg.Writer, cfg.Message)
}

// Status implements UI
func (ui *BasicUI) Status() Status {
	return newSpinnerStatus()
}

type config struct {
	// Original is the original message, this should NOT be modified.
	Original string

	// Message is the message to write.
	Message string

	// Writer is where the message will be written to.
	Writer io.Writer
}

// Option controls output styling.
type Option func(*config)

// WithHeaderStyle styles the output like a header denoting a new section
// of execution. This should only be used with single-line output. Multi-line
// output will not look correct.
func WithHeaderStyle() Option {
	return func(c *config) {
		c.Message = colorHeader.Sprintf("==> %s", c.Message)
	}
}

// WithErrorStyle styles the output as an error message.
func WithErrorStyle() Option {
	return func(c *config) {
		c.Message = colorError.Sprint(c.Original)
	}
}

var (
	colorHeader = color.New(color.FgGreen, color.Bold)
	colorError  = color.New(color.FgRed)
)
