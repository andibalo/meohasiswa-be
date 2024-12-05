package config

import (
	"fmt"
	"github.com/andibalo/meowhasiswa-be/pkg/logger"
	"github.com/andibalo/meowhasiswa-be/pkg/trace"
	"github.com/spf13/viper"
)

const (
	AppAddress         = ":8082"
	EnvDevEnvironment  = "DEV"
	EnvProdEnvironment = "PROD"
	ServiceName        = "core-service"
)

type Config interface {
	Logger() logger.Logger

	AppVersion() string
	AppID() string
	AppName() string
	AppEnv() string
	AppAddress() string

	DBConnString() string
	TraceConfig() trace.Config

	HttpExternalServiceTimeout() int64
	HttpMaxUploadSizeMB() int

	GetNotifSvcCfg() NotifSvc

	GetFlags() Flag
	GetAuthCfg() Auth
	GetAWSCfg() AWS
	GetBrevoSvcCfg() BrevoSvc
	GetMailerCfg() Mailer
}

type AppConfig struct {
	logger   logger.Logger
	App      app
	Db       db
	Tracer   tracer
	Http     http
	NotifSvc NotifSvc
	Flag     Flag
	Auth     Auth
	Aws      AWS
	BrevoSvc BrevoSvc
	Mailer   Mailer
}

type app struct {
	AppEnv      string
	AppVersion  string
	Name        string
	Description string
	AppUrl      string
	AppID       string
}

type db struct {
	DSN      string
	User     string
	Password string
	Name     string
	Host     string
	Port     int
	MaxPool  int
}

type tracer struct {
	ServiceName          string
	CollectorURL         string
	CollectorEnvironment string
	Insecure             bool
	FragmentRatio        float64
}

type NotifSvc struct {
	URL   string
	Token string
}

type http struct {
	ServiceExternalTimeout int64
	MaxUploadSizeMB        int
}

type Flag struct {
	EnableTracer    bool
	EnableSendEmail bool
}

type Auth struct {
	UserSecretCodeExpiryMins int
	JWTSecret                string
	JWTStaticToken           string
}

type AWS struct {
	Region            string
	ACCESS_KEY_ID     string
	SECRET_ACCESS_KEY string
	S3
}

type S3 struct {
	DefaultBucket string
}

type BrevoSvc struct {
	APIKey                         string
	SendVerificationCodeTemplateId int64
	SendResetPasswordTemplateId    int64
}

type Mailer struct {
	DefaultSenderName  string
	DefaultSenderEmail string
}

func InitConfig() *AppConfig {
	viper.SetConfigType("env")
	viper.SetConfigName(".env") // name of Config file (without extension)
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	l := logger.GetLogger(logger.Options{
		DefaultFields: map[string]string{
			"service.name":    ServiceName,
			"service.version": viper.GetString("APP_VERSION"),
			"service.env":     viper.GetString("APP_ENV"),
		},
		ContextFields: map[string]string{
			"path":        "path",
			"method":      "method",
			"status_code": "status_code",
			"status":      "status",
			"error":       "error",
			"user_id":     "x-user-id",
			"user_email":  "x-user-email",
			"client_ip":   "x-forwarded-for",
			"payload":     "payload",
			"x-client-id": "x-client-id",
			"topic":       "topic",
			"broker":      "broker",
			"trace.id":    "trace.id",
			"span.id":     "span.id",
		},
		Level:     logger.LevelInfo,
		HookLevel: logger.LevelError,
	})

	if err := viper.ReadInConfig(); err != nil {
		l.Warn("Env config file not found")
	}

	return &AppConfig{
		logger: l,
		App: app{
			AppEnv:      viper.GetString("APP_ENV"),
			AppVersion:  viper.GetString("APP_VERSION"),
			Name:        ServiceName,
			Description: "core service",
			AppUrl:      viper.GetString("APP_URL"),
			AppID:       viper.GetString("APP_ID"),
		},
		Db: db{
			DSN:      getRequiredString("DB_DSN"),
			User:     viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
			Host:     viper.GetString("DB_HOST"),
			Port:     viper.GetInt("DB_PORT"),
			Name:     viper.GetString("DB_NAME"),
			MaxPool:  viper.GetInt("DB_MAX_POOLING_CONNECTION"),
		},
		Http: http{
			MaxUploadSizeMB: viper.GetInt("MAX_UPLOAD_SIZE_MB"),
		},
		Tracer: tracer{
			ServiceName:          ServiceName,
			CollectorURL:         viper.GetString("OTEL_APM_SERVER_URL"),
			CollectorEnvironment: viper.GetString("OTEL_APM_ENV"),
			Insecure:             viper.GetBool("OTEL_APM_INSECURE"),
			FragmentRatio:        viper.GetFloat64("OTEL_JAEGER_FRACTION_RATIO"),
		},
		NotifSvc: NotifSvc{
			URL:   viper.GetString("NOTIF_SVC_URL"),
			Token: viper.GetString("NOTIF_SVC_TOKEN"),
		},
		Flag: Flag{
			EnableTracer:    viper.GetBool("ENABLE_TRACER"),
			EnableSendEmail: viper.GetBool("ENABLE_SEND_EMAIL"),
		},
		Auth: Auth{
			UserSecretCodeExpiryMins: viper.GetInt("USER_SECRET_CODE_EXPIRY_MINS"),
			JWTSecret:                viper.GetString("JWT_SECRET"),
			JWTStaticToken:           viper.GetString("JWT_STATIC_TOKEN"),
		},
		Aws: AWS{
			Region:            viper.GetString("AWS_REGION"),
			ACCESS_KEY_ID:     viper.GetString("AWS_ACCESS_KEY_ID"),
			SECRET_ACCESS_KEY: viper.GetString("AWS_SECRET_ACCESS_KEY"),
			S3: S3{
				DefaultBucket: viper.GetString("AWS_S3_DEFAULT_BUCKET"),
			},
		},
		BrevoSvc: BrevoSvc{
			APIKey:                         viper.GetString("BREVO_SVC_API_KEY"),
			SendVerificationCodeTemplateId: viper.GetInt64("BREVO_SVC_SEND_VERIFICATION_CODE_TEMPLATE_ID"),
			SendResetPasswordTemplateId:    viper.GetInt64("BREVO_SVC_SEND_RESET_PASSWORD_TEMPLATE_ID"),
		},
		Mailer: Mailer{
			DefaultSenderName:  viper.GetString("DEFAULT_SENDER_NAME"),
			DefaultSenderEmail: viper.GetString("DEFAULT_SENDER_EMAIL"),
		},
	}
}

func getRequiredString(key string) string {
	if viper.IsSet(key) {
		return viper.GetString(key)
	}

	panic(fmt.Errorf("KEY %s IS MISSING", key))
}

func (c *AppConfig) Logger() logger.Logger {
	return c.logger
}

func (c *AppConfig) AppVersion() string {
	return c.App.AppVersion
}

func (c *AppConfig) AppID() string {
	return c.App.AppID
}

func (c *AppConfig) AppName() string {
	return c.App.Name
}

func (c *AppConfig) AppEnv() string {
	return c.App.AppEnv
}

func (c *AppConfig) AppAddress() string {
	return AppAddress
}

func (c *AppConfig) DBConnString() string {
	return c.Db.DSN
}

func (c *AppConfig) TraceConfig() trace.Config {
	return trace.Config{
		ServiceName:          c.Tracer.ServiceName,
		CollectorURL:         c.Tracer.CollectorURL,
		CollectorEnvironment: c.Tracer.CollectorEnvironment,
		Insecure:             c.Tracer.Insecure,
		FragmentRatio:        c.Tracer.FragmentRatio,
	}
}

func (c *AppConfig) HttpExternalServiceTimeout() int64 {
	return c.Http.ServiceExternalTimeout
}

func (c *AppConfig) HttpMaxUploadSizeMB() int {

	return c.Http.MaxUploadSizeMB
}

func (c *AppConfig) GetNotifSvcCfg() NotifSvc {
	return c.NotifSvc
}

func (c *AppConfig) GetFlags() Flag {
	return c.Flag
}

func (c *AppConfig) GetAuthCfg() Auth {
	return c.Auth
}

func (c *AppConfig) GetAWSCfg() AWS {
	return c.Aws
}

func (c *AppConfig) GetBrevoSvcCfg() BrevoSvc {
	return c.BrevoSvc
}

func (c *AppConfig) GetMailerCfg() Mailer {
	return c.Mailer
}
