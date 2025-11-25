# alert-syncer-go

## General information

My solution consists of 3 elements:
1. mock-alerts-api - a very simple service which mocks upstream alert service. I didn't spend much time on the actual implementation. It's very basic.
1. alert-ingest-service - a service which pulls alerts from mock-alerts-api service, stores them in the db and exposes 3 API endpoints as described in the requirements.
1. PostgreSQL database - used to store alerts. alert-ingest-service integrates with this db through GORM
