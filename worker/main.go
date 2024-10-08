package main

import (
	"context"
	"log"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/interceptor"
	"go.temporal.io/sdk/worker"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"

	"go.temporal.io/sdk/contrib/opentelemetry"
	temporalOpentelemetry "go.temporal.io/sdk/contrib/opentelemetry"

	"money-transfer-project-template-go/app"
)

// The name of our service
// We pass this in to the resource attributes and it will appear on all our spans
var serviceName = semconv.ServiceNameKey.String("moneyTransfer.worker")

// @@@SNIPSTART money-transfer-project-template-go-worker
func main() {

	ctx := context.Background()

	// Create the worker "Resource"
	// https://opentelemetry.io/docs/languages/js/resources/
	res, err := resource.New(ctx,
		resource.WithAttributes(
			// The service name used to display traces in backends
			serviceName,
			// We can pass in other attributes here that will appear
			// in all spans generated by this service
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	shutdownTracerProvider, err := app.NewHttpTracerProvider(ctx, res) // jaegerConn)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := shutdownTracerProvider(ctx); err != nil {
			log.Fatalf("failed to shutdown TracerProvider: %s", err)
		}
	}()

	// A lot of the above is abstracted away by our internal opentelemetry module

	// Let's Create a span here to check our set up is working.
	tracer := otel.Tracer("")

	// We will make a span that captures the time it takes to configure
	// the worker with temporal
	ctx, span := tracer.Start(ctx, "worker.setUp")

	// Create the tracing interceptor to let us know what the worker is doing.
	tracingInterceptor, err := opentelemetry.NewTracingInterceptor(temporalOpentelemetry.TracerOptions{})

	temporalConn, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create Temporal client.", err)
	}
	defer temporalConn.Close()

	w := worker.New(temporalConn, app.MoneyTransferTaskQueueName, worker.Options{
		Interceptors: []interceptor.WorkerInterceptor{tracingInterceptor},
	})

	// This worker hosts both Workflow and Activity functions.
	w.RegisterWorkflow(app.MoneyTransfer)
	w.RegisterActivity(app.Withdraw)
	w.RegisterActivity(app.Deposit)
	w.RegisterActivity(app.Refund)

	// The worker is now ready to run so lets end the span.
	// You always need to end spans that you create
	span.End()

	// Start listening to the Task Queue.
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("unable to start Worker", err)
	}
}

// @@@SNIPEND
