# flaggio

Flaggio is a feature flag system that runs on your infrastructure. It supports single on/off as well as string and numeric flag values, user segmentation, and percentage rollout of features.



## How to run

#### External dependencies

* MongoDB 4+ (required)
* Redis (recommended)
* [Jaeger](https://github.com/jaegertracing/jaeger) (optional)

The easiest way is to run flaggio with docker:

```shell script
$ docker run --rm -p 8080:8080 -p 8081:8081 -t flaggio/flaggio:latest -database-uri <MONGO_URI>
```

Flaggio UI will then be available at http://localhost:8081.

## Concepts

### Flags

Flags consist of a key and a value (one of the variants), and they can be used to toggle parts of your application on or off, change the appearance of a UI element, and more.

### Variants

Variants are the values a flag can return. These can be a boolean, a number, or a string. Boolean values are useful for feature-toggling flags, whereas numbers and strings enable additional use cases.

### Rules

Rules define a set of constraints and a specific variant to return when all constraint requirements are met. For example, if the user is using Chome browser return `blue`.

### Constraints

Constraints define what field and values to look for on the user context. It can also be used to check if they belong to a certain segment. For example, the user's country should equal Brazil.

### User context

Thse are any values associated with a user. For example `age = 24`, `country = France`, `browser = Chrome`, `operationalSystem = Windows`, etc.

### Segments

Segments are a group of users that share a common set of properties. For example "Users from the UK",  "Age 30-40", "MacOS users", etc.

## Architecture

Flaggio is comprised of two APIs and a UI to manage the flags and segments, as well as being able to view the flag evaluations for each user.

### Admin API

This is a graphql API that is able to perform CRUD operations for flags and segments.

### Evaluation API

This is a REST JSON API which takes the user context and returns the flag value.

#### Request model

|value|type|required|description|
|-----|----|--------|-----------|
|userId|string|yes|an arbitrary ID that identifies a unique user|
|context|object|yes|a set of values associated with the user|
|debug|boolean|no|returns additional debugging information when `true`|

#### Example request

```json
{
  "userId": "john@doe.com",
  "context": {
    "name": "john",
    "age": 26,
    "browser": "Firefox"
  },
  "debug": false
}
```

#### Example response

```json
{
  "evaluations": [
    {
      "flagKey": "showHeader",
      "value": true
    },
    {
      "flagKey": "backgroundColor",
      "value": "#FFFFFF"
    }
  ]
}
```

## Configuration

The flaggio CLI accepts the following options:

 ```
   --database-uri value          Database URI [$DATABASE_URI]
   --redis-uri value             Redis URI [$REDIS_URI]
   --build-path value            UI build absolute path [$BUILD_PATH]
   --cors-allowed-origins value  CORS allowed origins separated by comma [$CORS_ALLOWED_ORIGINS]
   --cors-allowed-headers value  CORS allowed headers [$CORS_ALLOWED_HEADERS]
   --no-api                      Don't start the API server (default: false) [$NO_API]
   --no-admin                    Don't start the admin server (default: false) [$NO_ADMIN]
   --no-admin-ui                 Don't start the admin UI (default: false) [$NO_ADMIN_UI]
   --playground                  Enable graphql playground (default: false) [$PLAYGROUND]
   --api-addr value              Sets the bind address for the API (default: ":8080") [$API_ADDR]
   --admin-addr value            Sets the bind address for the admin (default: ":8081") [$ADMIN_ADDR]
   --log-formatter value         Sets the log formatter for the application. Valid values are: text, json (default: "json") [$LOG_FORMATTER]
   --log-level value             Sets the log level for the application (default: "info") [$LOG_LEVEL]
   --jaeger-agent-host value     The address of the jaeger agent (host:port) [$JAEGER_AGENT_HOST]
```

## License

Apache License 2.0