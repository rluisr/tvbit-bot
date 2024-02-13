/*
 *
 * tvbit-bot
 * Copyright (C) 2022  rluisr(Takuya Hasegawa)
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 * /
 */

package external

import (
	"context"

	"github.com/rluisr/tvbit-bot/pkg/adapter/controllers"
	tvbitBybit "github.com/rluisr/tvbit-bot/pkg/external/bybit"
	"github.com/rluisr/tvbit-bot/pkg/external/logging"
	"github.com/rluisr/tvbit-bot/pkg/external/mysql"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	otelresource "go.opentelemetry.io/otel/sdk/resource"
	oteltrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

var (
	tvController *controllers.TVController
)

// OTLP Exporter
func newOTLPExporter(ctx context.Context) (oteltrace.SpanExporter, error) {
	// Change default HTTPS -> HTTP
	insecureOpt := otlptracehttp.WithInsecure()
	return otlptracehttp.New(ctx, insecureOpt)
}

// TracerProvider is an OpenTelemetry TracerProvider.
// It provides Tracers to instrumentation so it can trace operational flow through a system.
func newTraceProvider(serviceName string, exp oteltrace.SpanExporter) *oteltrace.TracerProvider {
	// Ensure default SDK resources and the required service name are set.
	r, err := otelresource.Merge(
		otelresource.Default(),
		otelresource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(serviceName),
		),
	)

	if err != nil {
		panic(err)
	}

	return oteltrace.NewTracerProvider(
		oteltrace.WithBatcher(exp),
		oteltrace.WithResource(r),
	)
}

func Init(source string) (err error) {
	log, err := logging.New(source)
	if err != nil {
		return err
	}

	ctx := context.Background()
	exp, err := newOTLPExporter(ctx)
	if err != nil {
		return err
	}
	tp := newTraceProvider(source, exp) // logger for middleware

	otel.SetTracerProvider(tp)

	rwDB, roDB, err := mysql.Connect()
	if err != nil {
		return err
	}

	httpClient := NewHTTPClient()
	bybitClient := tvbitBybit.Init(httpClient)

	tvController = controllers.NewTVController(log, rwDB, roDB, bybitClient)

	return nil
}
