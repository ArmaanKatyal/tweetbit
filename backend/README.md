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
