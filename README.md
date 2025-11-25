# alert-syncer-go

## General information

My solution consists of 3 main components:
1. mock-alerts-api - a very simple service which mocks upstream alert service. I didn't spend much time on the actual implementation. It's very basic.
1. alert-ingest-service - a service which pulls alerts from mock-alerts-api service, stores them in the db and exposes 3 API endpoints as described in the requirements.
1. PostgreSQL database - used to store alerts. alert-ingest-service integrates with this db through GORM. I decided to use Postrges mainly because of time constraints to complete the task. In real, production ready environment I'll have to evaluate if there are better solutions. I'll consider using NoSQL database for this use case. In general I'll consider using SQL database as a safe choice. It should handle even high traffic without bigger issues.

I tried to follow a uniform approach for upstream integrations. In general I treat all dependencies (here db and mock-alerts-api) the same way. There's an abstract client implemented which (should) provide a basic functionalities:
1. Retries - right now hardcoded to 3 retries and simple backoff strategy. This should be configurable for each upstream dependency.
1. Client side observability - not implemented (only logged) but in general I would like to emit latency, failed and successful requests. This will allow us to build dashboards which will clearly show performance of our dependencies and build internal alerts on top of that.
1. Error handling - simplify the error handling for alert-ingest-service clients. We translate upstream errors to the standarized errors. We also try to detect a non-retryable errors for a "circuit breaker"-like approach.

For server side operations I also tried to unify the approach for APIs which are exposed by the alert-ingest-service. The generic wrapper is responsible for:
1. Rate limiting - this should be moved to a middlware and be applied before the request hit our endpoints. For now we have a dummy implementation.
1. Server side observability - not implemented (only logged) but in general I would like to emit latency, failed and successful requests. Together with client side metrics it will give us a full picture of how alert-ingest-service is performing. Example: it will be very easy to correlate server side errors with dependency (db or mock-alerts-api) errors, or elevated server side latency with eleveted upstream latency. Or, actual internal errors if there are no upstream errors at the same time.

## What's missing

Simple answer: A lot!!!

You'll notice there's a lot of TODOs in the code. Mainly:
1. No test coverage. Usually I aim for at 90% test coverage. In this case I run out of time. Also, I have to admit that Go is not my main programming language. In many places I should relay more on interfaces than on concrete implementations. This will improve the testability of this code. I realized that late and as a result unit test are missing. I'll need few weeks to switch from Java/Python to Go.
1. No integration tests. There are no, end-to-end tests of the system. Ideally I would deploy all the components of the system to a dedicated infrastructure (on-prem or in cloud), build a mechanism which allows me to prepare test data, run tests against this data and then remove the data.
1. Simplified error handling - I run out of time to fully implement it. Current implementation shows more or less what I wanted to achieve. Especially for the /sync API, where the API call triggers calls to the upstream, I should propagate a translated outcome of the upstream call all the way to the client. Examples: if upstream reposrts non-retryable error I should also return non-retryable error to my client. Or, if upstream call failed for all retries with InternalServerError I should also retrun with InternalServerError.
1. Observability. You'll notice I'm logging latencies and errors. I would rather emit metrics for those. We can build OPS dahsboards and base alerts on top of them.
1. Authentication. There's no authn/authz mechanism. Ideally I would have 2 lines of defense. One in higher-level components like API Gateway and inside the service itself. I run out of time to implement even a placeholder.
1. Structured API responses for alert-ingest-service. Right now the responses are not structured. Ideally I would have a response which always contains a specific set of fields. Regardless of the outcome (success/failure) these fields are always present but only some of them are actually propagated with data. For example, we can have a structure with 3 main fields: response, error and additional_info. In case of critical error, only error is propagated. In case of non-critical error, response is propagated and additional_info contains degradation warning. This is to simplify integration for clients.
1. Rate limiting middlware. Right now rate limitting is integrated into REST handlers. Ideally rate limitting should happen earlier before request hits the handler.
1. Alert deduplication logic. Current implementation doesn't handle cases where the same alert is propagated multiple times from the upstream.

## How to start
1. Pull the repo
1. From the repo's root directory run `docker compose up` command - this should start all 3 components in separate docker containers. I'm assuming that docker cli is installed.

## How to test
1. Start all 3 container as described above
1. There are very basic 4 tests implemented inside the Taskfile.yml. You can execute it by running `task test-all` command from the repo's root directory. I'm assuming that task cli is installed.
1. additinal manual tests can be performed using curl. Few examples:
 * curl "http://localhost:8080/health"
 * curl "http://localhost:8080/sync"
 * curl "http://localhost:8080/alerts?limit=10"
 * curl "http://localhost:8080/alerts?severity=low&limit=10"
