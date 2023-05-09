module github.com/gocondor/core

replace github.com/gocondor/core/logger => ./logger

replace github.com/gocondor/core/env => ./env

go 1.20

require (
	github.com/go-redis/redis/v8 v8.8.0
	github.com/google/uuid v1.3.0
	github.com/joho/godotenv v1.3.0
	github.com/julienschmidt/httprouter v1.3.0
	golang.org/x/net v0.9.0
	gorm.io/driver/mysql v1.0.5
	gorm.io/driver/sqlite v1.1.4
	gorm.io/gorm v1.21.6
)

require (
	github.com/cespare/xxhash/v2 v2.1.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/go-sql-driver/mysql v1.5.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.2 // indirect
	github.com/mattn/go-sqlite3 v1.14.5 // indirect
	go.opentelemetry.io/otel v0.19.0 // indirect
	go.opentelemetry.io/otel/metric v0.19.0 // indirect
	go.opentelemetry.io/otel/trace v0.19.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
