package http

import (
	"context"

	"github.com/eser/ajan/httpfx"
	"github.com/eser/ajan/httpfx/middlewares"
	"github.com/eser/ajan/httpfx/modules/healthcheck"
	"github.com/eser/ajan/httpfx/modules/openapi"
	"github.com/eser/ajan/httpfx/modules/profiling"
	"github.com/eser/ajan/logfx"
	"github.com/eser/ajan/metricsfx"
	"github.com/eser/aya.is-services/pkg/api/business/profiles"
	"github.com/eser/aya.is-services/pkg/api/business/stories"
	"github.com/eser/aya.is-services/pkg/api/business/users"
)

func Run(
	ctx context.Context,
	config *httpfx.Config,
	metricsProvider *metricsfx.MetricsProvider,
	logger *logfx.Logger,
	profilesService *profiles.Service,
	storiesService *stories.Service,
	usersService *users.Service,
) (func(), error) {
	routes := httpfx.NewRouter("/")
	httpService := httpfx.NewHttpService(config, routes, metricsProvider, logger)

	// http middlewares
	routes.Use(middlewares.ErrorHandlerMiddleware())
	routes.Use(middlewares.ResolveAddressMiddleware())
	routes.Use(middlewares.ResponseTimeMiddleware())
	routes.Use(middlewares.CorrelationIdMiddleware())
	routes.Use(middlewares.CorsMiddleware())
	routes.Use(middlewares.MetricsMiddleware(httpService.InnerMetrics))
	// routes.Use(AuthMiddleware(usersService))

	// http modules
	healthcheck.RegisterHttpRoutes(routes, config)
	openapi.RegisterHttpRoutes(routes, config)
	profiling.RegisterHttpRoutes(routes, config)

	// http routes
	RegisterHttpRoutesForUsers( //nolint:contextcheck
		routes,
		logger,
		usersService,
	)
	RegisterHttpRoutesForSite( //nolint:contextcheck
		routes,
		logger,
		profilesService,
	)
	RegisterHttpRoutesForProfiles( //nolint:contextcheck
		routes,
		logger,
		profilesService,
		storiesService,
	)
	RegisterHttpRoutesForStories( //nolint:contextcheck
		routes,
		logger,
		storiesService,
	)

	// run
	return httpService.Start(ctx) //nolint:wrapcheck
}
