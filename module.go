package echomodule

import (
	"context"
	"errors"

	"go.viam.com/rdk/components/generic"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
	"go.viam.com/utils/rpc"
)

var (
	EchoModel        = resource.NewModel("bashar-prod", "echo-module", "echo-model")
	errUnimplemented = errors.New("unimplemented")
)

func init() {
	resource.RegisterComponent(generic.API, EchoModel,
		resource.Registration[resource.Resource, *Config]{
			Constructor: newEchoModuleEchoModel,
		},
	)
}

type Config struct {
	/*
		Put config attributes here. There should be public/exported fields
		with a `json` parameter at the end of each attribute.

		Example config struct:
			type Config struct {
				Pin   string `json:"pin"`
				Board string `json:"board"`
				MinDeg *float64 `json:"min_angle_deg,omitempty"`
			}

		If your model does not need a config, replace *Config in the init
		function with resource.NoNativeConfig
	*/
}

// Validate ensures all parts of the config are valid and important fields exist.
// Returns implicit required (first return) and optional (second return) dependencies based on the config.
// The path is the JSON path in your robot's config (not the `Config` struct) to the
// resource being validated; e.g. "components.0".
func (cfg *Config) Validate(path string) ([]string, []string, error) {
	var (
		reqDeps []string
		optDeps []string
	)

	return reqDeps, optDeps, nil
}

type echoModuleEchoModel struct {
	resource.AlwaysRebuild

	name resource.Name

	logger logging.Logger
	cfg    *Config

	cancelCtx  context.Context
	cancelFunc func()
}

func newEchoModuleEchoModel(ctx context.Context, deps resource.Dependencies, rawConf resource.Config, logger logging.Logger) (resource.Resource, error) {
	conf, err := resource.NativeConfig[*Config](rawConf)
	if err != nil {
		return nil, err
	}

	return NewEchoModel(ctx, deps, rawConf.ResourceName(), conf, logger)

}

func NewEchoModel(ctx context.Context, deps resource.Dependencies, name resource.Name, conf *Config, logger logging.Logger) (resource.Resource, error) {

	cancelCtx, cancelFunc := context.WithCancel(context.Background())

	s := &echoModuleEchoModel{
		name:       name,
		logger:     logger,
		cfg:        conf,
		cancelCtx:  cancelCtx,
		cancelFunc: cancelFunc,
	}
	return s, nil
}

func (s *echoModuleEchoModel) Name() resource.Name {
	return s.name
}

func (s *echoModuleEchoModel) NewClientFromConn(ctx context.Context, conn rpc.ClientConn, remoteName string, name resource.Name, logger logging.Logger) (resource.Resource, error) {
	panic("not implemented")
}

func (s *echoModuleEchoModel) DoCommand(ctx context.Context, cmd map[string]any) (map[string]any, error) {
	name, ok := cmd["name"]
	if !ok {
		return nil, errors.New("name not specified")
	}

	switch name {
	case "number":
		rawNumber, ok := cmd["number"]
		if !ok {
			return nil, errors.New("number not specified")
		}

		number, ok := rawNumber.(float64)
		if !ok {
			return nil, errors.New("unable to parse number as float")
		}

		return map[string]any{ "number": number }, nil
	default:
		return nil, errors.New("name not defined")
	}
}

func (s *echoModuleEchoModel) Close(context.Context) error {
	// Put close code here
	s.cancelFunc()
	return nil
}
