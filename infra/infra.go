package infra

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Infra interface {
	Config() *viper.Viper
	SetMode() string
	GormDB() *gorm.DB
	Cipher(typeCrypto string) cipher.Stream
	GoMail() *gomail.Dialer
	SendGrid() *sendgrid.Client
	Migrate(values ...interface{})
	Port() string
}

type infra struct {
	configFile string
}

func New(configFile string) Infra {
	return &infra{configFile: configFile}
}

var (
	vprOnce sync.Once
	vpr     *viper.Viper
)

func (i *infra) Config() *viper.Viper {
	vprOnce.Do(func() {
		viper.SetConfigFile(i.configFile)
		if err := viper.ReadInConfig(); err != nil {
			logrus.Fatalf("[infra][Config][viper.ReadInConfig] %v", err)
		}

		vpr = viper.GetViper()
	})

	return vpr
}

var (
	modeOnce    sync.Once
	mode        string
	development = "development"
	production  = "production"
)

func (i *infra) SetMode() string {
	modeOnce.Do(func() {
		env := i.Config().Sub("environment").GetString("mode")
		if env == development {
			mode = gin.DebugMode
		} else if env == production {
			mode = gin.ReleaseMode
		} else {
			logrus.Fatalf("[infa][SetMode] %v", errors.New("environment not setup"))
		}

		gin.SetMode(mode)
	})

	return mode
}

var (
	grmOnce sync.Once
	grm     *gorm.DB
)

func (i *infra) GormDB() *gorm.DB {
	grmOnce.Do(func() {
		config := i.Config().Sub("database")
		user := config.GetString("user")
		pass := config.GetString("pass")
		host := config.GetString("host")
		port := config.GetString("port")
		name := config.GetString("name")

		// dns := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, pass, host, port, name)
		// db, err := gorm.Open(postgres.Open(dns), &gorm.Config{})
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, pass, host, port, name)
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			logrus.Fatalf("[infra][GormDB][gorm.Open] %v", err)
		}

		sqlDB, err := db.DB()
		if err != nil {
			logrus.Fatalf("[infra][GormDB][db.DB] %v", err)
		}

		if err := sqlDB.Ping(); err != nil {
			logrus.Fatalf("[infra][GormDB][sqlDB.Ping] %v", err)
		}

		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetConnMaxLifetime(time.Hour)

		grm = db
	})

	return grm
}

var (
	myCryptoOnce sync.Once
	myCrypto     cipher.Stream
)

func (i *infra) Cipher(typeCrypto string) cipher.Stream {
	myCryptoOnce.Do(func() {
		bytes := []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}
		secretKey := i.Config().Sub("secret")
		keyCrypto := secretKey.GetString("crypto")

		block, err := aes.NewCipher([]byte(keyCrypto))
		if err != nil {
			log.Printf("[Error][Initial Cipher] E: %v", err)
		}
		if typeCrypto == "encrypt" {
			cfb := cipher.NewCFBEncrypter(block, bytes)
			myCrypto = cfb
		} else {
			cfb := cipher.NewCFBDecrypter(block, bytes)
			myCrypto = cfb
		}
	})
	return myCrypto
}

var (
	goMailOnce sync.Once
	goMail     *gomail.Dialer
)

func (i *infra) GoMail() *gomail.Dialer {
	goMailOnce.Do(func() {
		config := i.Config().Sub("smtp")
		user := config.GetString("user")
		pass := config.GetString("pass")
		host := config.GetString("host")
		port := config.GetInt("port")

		dialer := gomail.NewDialer(
			host,
			port,
			user,
			pass,
		)

		goMail = dialer
	})
	return goMail
}

var (
	sendGridOnce sync.Once
	sendGrid     *sendgrid.Client
)

func (i *infra) SendGrid() *sendgrid.Client {
	sendGridOnce.Do(func() {
		config := i.Config().Sub("smtp")
		pass := config.GetString("pass")

		client := sendgrid.NewSendClient(pass)

		sendGrid = client
	})
	return sendGrid
}

var (
	migrateOnce sync.Once
)

func (i *infra) Migrate(values ...interface{}) {
	migrateOnce.Do(func() {
		if i.SetMode() == gin.DebugMode {
			if err := i.GormDB().Debug().AutoMigrate(values...); err != nil {
				logrus.Fatalf("[infra][Migrate][GormDB.Debug.AutoMigrate] %v", err)
			}
		} else if i.SetMode() == gin.ReleaseMode {
			if err := i.GormDB().AutoMigrate(values...); err != nil {
				logrus.Fatalf("[infra][Migrate][GormDB.AutoMigrate] %v", err)
			}
		}
	})
}

var (
	portOnce sync.Once
	port     string
)

func (i *infra) Port() string {
	portOnce.Do(func() {
		port = i.Config().Sub("server").GetString("port")
	})

	return ":" + port
}
