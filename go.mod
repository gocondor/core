module github.com/gocondor/core

replace github.com/gocondor/core/logger => ./logger

replace github.com/gocondor/core/env => ./env

go 1.20

require (
	github.com/brianvoe/gofakeit/v6 v6.21.0
	github.com/go-redis/redis/v8 v8.11.5
	github.com/golang-jwt/jwt/v5 v5.0.0
	github.com/google/uuid v1.3.0
	github.com/joho/godotenv v1.5.1
	github.com/julienschmidt/httprouter v1.3.0
	golang.org/x/crypto v0.9.0
	gorm.io/gorm v1.25.2
)

require (
	github.com/asaskevich/govalidator v0.0.0-20230301143203-a9d515a09cc2 // indirect
	github.com/go-sql-driver/mysql v1.7.0 // indirect
	github.com/stretchr/testify v1.8.2 // indirect
	golang.org/x/net v0.10.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	gorm.io/driver/mysql v1.5.1 // indirect
)

require (
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/go-ozzo/ozzo-validation v3.6.0+incompatible
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
)
