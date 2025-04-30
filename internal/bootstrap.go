package internal

import (
	"context"
	"fmt"
	"testing"

	"github.com/ankorstore/yokai/fxconfig"
	"github.com/ankorstore/yokai/fxcore"
	"github.com/ankorstore/yokai/fxhttpclient"
	"github.com/ankorstore/yokai/fxhttpserver"
	"github.com/ankorstore/yokai/fxsql"
	"github.com/ekkinox/yokai-mcp/pkg/mcp"
	"go.uber.org/fx"
)

func init() {
	RootDir = fxcore.RootDir(1)
}

// RootDir is the application root directory.
var RootDir string

// Bootstrapper can be used to load modules, options, dependencies, routing and bootstraps the application.
var Bootstrapper = fxcore.NewBootstrapper().WithOptions(
	// modules registration
	fxhttpserver.FxHttpServerModule,
	fxhttpclient.FxHttpClientModule,
	fxsql.FxSQLModule,
	mcp.MCPServerModule,
	// dependencies registration
	Register(),
)

// Run starts the application, with a provided [context.Context].
func Run(ctx context.Context) {
	Bootstrapper.WithContext(ctx).RunApp()
}

// RunTest starts the application in test mode, with an optional list of [fx.Option].
func RunTest(tb testing.TB, options ...fx.Option) {
	tb.Helper()

	Bootstrapper.RunTestApp(
		tb,
		// config lookup
		fxconfig.AsConfigPath(fmt.Sprintf("%s/configs/", RootDir)),
		// test options
		fx.Options(options...),
	)
}
