// Package datasource has the packages and logic that Waypoint uses
// for sourcing data for remote runs.
package datasource

import (
	"context"
	"reflect"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/hcl/v2"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
)

// Sourcer is implemented by all data sourcers and is responsible for
// sourcing data, configuring projects, determining the default values
// for operations, and more.
type Sourcer interface {
	// ProjectSource translates the configuration into a default data
	// source for the project. This should also perform any validation
	// on the configuration.
	ProjectSource(hcl.Body, *hcl.EvalContext) (*pb.Job_DataSource, error)

	// Override reconfigures the given data source with the given overrides.
	Override(*pb.Job_DataSource, map[string]string) error

	// Get downloads the sourced data and returns the directory where
	// the data is stored, a cleanup function, and any errors that occurred.
	// The cleanup function may be nil.
	Get(
		ctx context.Context,
		log hclog.Logger,
		ui terminal.UI,
		source *pb.Job_DataSource,
		baseDir string,
	) (string, func() error, error)
}

var (
	// FromString maps a string key to a source implementation.
	FromString = map[string]func() Sourcer{
		"git":   newGitSource,
		"local": newLocalSource,
	}

	// FromType maps a server DataSource type to a source implementation.
	FromType = map[reflect.Type]func() Sourcer{
		reflect.TypeOf((*pb.Job_DataSource_Git)(nil)):   newGitSource,
		reflect.TypeOf((*pb.Job_DataSource_Local)(nil)): newLocalSource,
	}
)
