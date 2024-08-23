# Intro for temporal
https://github.com/temporalio/money-transfer-project-template-go. The original readme is preserved at the bottom

## Why do we use temporal?
It provides functionality out of the box that allows us to model our business
processes as code and makes them
- Resilient to failure
- Able to be gracefully interrupted and resumed

In return temporal asks us that we write our code a certain way. Our code is
split up into two concepts
- Activity Code
- Workflow Code

Workflow code must be deterministic! This means given the same inputs and
results, it should do the same thing every time. Code that can't guaruntee that
doesn't exist in workflows. You can think of a workflow as something that
describes a flow chart or a state machine.

Activity code isn't required to look a certain way. In this code, you can run
do non deterministic things. Eg - An api call to an external service, Run
random number generators, do something based on the time of day.

You may want activity code to be idempotent depending on your use case as
activities can fail and be re ran or even succeed and be ran again if there are
failures between the worker and the server.

There is a lot of complexity to temporal and it is worth reading through the
documentation (multiple times) and trying to run and break things.

Give the following a read
- https://docs.temporal.io/temporal
- https://docs.temporal.io/evaluate/why-temporal

## Running things...

Install the temopral development server

```bash
brew install temporal
```

Run the temporal dev server. The data is ephemeral when starts this way.

```bash
temporal server start-dev
```

Launch the jaeger container

```bash
./run-jaeger.sh
```

Start the worker

```bash
go run worker/main.go
```

Trigger the workflow

```bash
go run start/main.go
```

Check that you can run all of the above and see a finished workflow in the temporal ui and some trace data in the jaeger UI
- temporal: http://localhost:8233/
- jaeger: http://localhost:16686/

If you are having troubles, try disconnecting from the VPN. That was causing issues with my worker and starter talking to the temporal server and jaeger.

## Exercises

Work through the following:
- Let the workflow run successfully. Look at the data stored in the webview. Does it make sense?
- Match that up to the trace data in Jaeger

Challenges:
- Add a span that is a child of the existing worker start up span. It should capture how long it takes to register the activities and the workflow. How is this displayed in Jaeger?
- Add an attribute to the worker resource. Validate that it shows up in all spans.
- Add an attribute to one or both of the start up spans. Validate that it only shows up in the spans you've added it to.
- Add the same tracing interceptor to the start code. Give the start code a different service name
- Investigate the trace. Can you see spans from the starter and the worker in the same trace? (You should be able to when everything is set up correctly)

Some hints:
- You can copy a lot of the set up code from `worker/main.go` to get the start code instrumented.
- The temporal client can also be configured with interceptors like the worker. It takes the following in its options `Interceptors: []interceptor.ClientInterceptor`

More things to work through:
- Make some activities that may fail. What what does this look like in the temporal ui and Jaeger when they fail? Hint - Import `math/rand` and use `rand.Intn` to return a number between 0 and 10. Let the activity succeed if the number is > 8.
- Introduce some non determinism into the workflow code and see if you can cause a non determinism error. Hint - This will 


### TODO - Add some skeleton code and examples for the below

Interceptors and context propagators can give SRE some super powers. Let's explore what they are, how we can write them and what we might want to use them for.

- https://docs.temporal.io/develop/go/observability#context-propagation
- https://pkg.go.dev/go.temporal.io/sdk/interceptor#Interceptor
- https://pkg.go.dev/go.temporal.io/sdk/workflow#ContextPropagator


# Original content below...

# Temporal Go Project Template

This is a simple project for demonstrating Temporal with the Go SDK.

The full 10-minute tutorial is here: https://learn.temporal.io/getting_started/go/first_program_in_go/

## Basic instructions

### Step 0: Temporal Server

Make sure [Temporal Server is running](https://docs.temporal.io/docs/server/quick-install/) first:

```bash
git clone https://github.com/temporalio/docker-compose.git
cd  docker-compose
docker-compose up
```

### Step 1: Clone this Repository

In another terminal instance, clone this repo and run this application.

```bash
git clone https://github.com/temporalio/money-transfer-project-template-go
cd money-transfer-project-template-go
```

### Step 2: Run the Workflow

```bash
go run start/main.go
```

Observe that Temporal Web reflects the workflow, but it is still in "Running" status. This is because there is no Workflow or Activity Worker yet listening to the `TRANSFER_MONEY_TASK_QUEUE` task queue to process this work.

### Step 3: Run the Worker

In YET ANOTHER terminal instance, run the worker. Notice that this worker hosts both Workflow and Activity functions.

```bash
go run worker/main.go
```

Now you can see the workflow run to completion. You can also see the worker polling for workflows and activities in the task queue at [http://localhost:8080/namespaces/default/task-queues/TRANSFER_MONEY_TASK_QUEUE](http://localhost:8080/namespaces/default/task-queues/TRANSFER_MONEY_TASK_QUEUE).

## What Next?

You can run the Workflow code a few more times with `go run start/main.go` to understand how it interacts with the Worker and Temporal Server.

Please [read the tutorial](https://learn.temporal.io/getting_started/go/first_program_in_go/) for more details.
