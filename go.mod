module clean

// +heroku goVersion go1.16
go 1.16

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-ozzo/ozzo-validation/v4 v4.3.0
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/google/uuid v1.2.0
	github.com/labstack/echo/v4 v4.3.0
	github.com/spf13/cobra v1.1.3
	github.com/spf13/viper v1.7.1
	go.uber.org/zap v1.17.0
	golang.org/x/crypto v0.0.0-20220214200702-86341886e292
	gorm.io/gorm v1.23.2
)

require (
	github.com/onsi/ginkgo v1.16.5 // indirect
	github.com/onsi/gomega v1.19.0 // indirect
	gorm.io/datatypes v1.0.6
	gorm.io/driver/mysql v1.3.2
	gorm.io/driver/sqlite v1.3.1
)
