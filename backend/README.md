# Backend

overview of the backend services and their dependencies.

## Local Development
Clone the repository, cd into backend directory, and run:
```
make run
```
this will start all the services in the docker-compose.yml file.

```
make stop
```
will stop all the services.

## Services

### AuthService
Manages the user identity and generates auth and refresh tokens that are used to authenticate requests to other services.

#### Dependent on:
- [MongoDB](https://www.mongodb.com/)
- [Redis](https://redis.io/)
- [Postgres](https://www.postgresql.org/)

### FanoutService
Manages the fanout of events from the writeservice to the other services.

#### Dependent on:
- [Kafka](https://kafka.apache.org/)

### ReadService
Manages the read model of the application. It is responsible for querying the database and returning the data to the client.

#### Dependent on:
- [TimelineService](#timelineservice)
- [Postgres](https://www.postgresql.org/)

### SearchService
Manages the search index of the application. It is responsible for indexing and searching the data.

#### Dependent on:
- [Elasticsearch](https://www.elastic.co/)
- [Kafka](https://kafka.apache.org/)

### TimelineService
Manages the timeline of the application. It is responsible for creating and updating the timeline of a user.

#### Dependent on:
- [Redis](https://redis.io/)
- [Postgres](https://www.postgresql.org/)

### UserGraphService
Manages the user graph of the application. It is reponsible for creating and updating user followers and following.

#### Dependent on:
- [Redis](https://redis.io/)
- [Kafka](https://kafka.apache.org/)

### WriteService
Manages the write model of the application. It is responsible for creating and updating the data in the database.

#### Dependent on:
- [Postgres](https://www.postgresql.org/)
- [FanoutService](#fanoutservice)

## Database
The database is a Postgres database. It is used by the [ReadService](#readservice), [TimelineService](#timelineservice), and [WriteService](#writeservice).

## Search Index
The search index is an Elasticsearch index. It is used by the [SearchService](#searchservice).

## Message Broker
The message broker is a Kafka cluster. It is used by the [FanoutService](#fanoutservice), [SearchService](#searchservice) and [UserGraphService](#usergraphservice).

## Cache
The cache is a Redis cluster. It is used by the [AuthService](#authservice), [TimelineService](#timelineservice), and [UserGraphService](#usergraphservice).

## API Gateway
The API Gateway is an NGINX server. It is used to route requests to the correct service. For kubernetes an ingress controller is used instead.
- [NGINX](https://www.nginx.com/)
- [Kubernetes Ingress](https://kubernetes.io/docs/concepts/services-networking/ingress/)

## Monitoring
The monitoring is done with Prometheus and Grafana. It is used to monitor the health of the services and the cluster.
- [Prometheus](https://prometheus.io/)
- [Grafana](https://grafana.com/)

## Logging
The logging is done with the ELK stack. It is used to log the requests and errors of the services.
- [Elasticsearch](https://www.elastic.co/)
- [Logstash](https://www.elastic.co/logstash)
- [Kibana](https://www.elastic.co/kibana)

## Tracing
The tracing is done with Jaeger. It is used to trace the requests between the services.
- [Opentelemetry](https://opentelemetry.io/)