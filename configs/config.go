package configs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

const (
	activeWorker = "ACTIVE_WORKER"

	androidPOSAppMinimumVersion = "ANDROID_POS_APP_MINIMUM_VERSION"
	iosPOSAppMinimumVersion     = "IOS_POS_APP_MINIMUM_VERSION"

	appUrl   = "APP_URL"
	portApps = "PORT_APPS"

	serverName = "SERVER_NAME"

	configFileLocation = "CONF_ENV_LOCATION"
	dbConnectionString = "DB_CONNECTION_STRING"

	redisAddr     = "REDIS_ADDR"
	redisDB       = "REDIS_DB"
	redisPassword = "REDIS_PASSWORD"
	redisTimeOut  = "REDIS_TIME_OUT"

	jwtTimeOut = "JWT_TIME_OUT"

	whitelistedIps = "WHITELISTED_IPS"

	vultrAccessKey = "VULTR_ACCESS_KEY"
	vultrBucket    = "VULTR_BUCKET"
	vultrHostname  = "VULTR_HOSTNAME"
	vultrSecretKey = "VULTR_SECRET_KEY"
	vultrRegion    = "VULTR_REGION"
)

// TODO check mana yg masih dipakai
var (
	JwtActiveToken *string
)

// Config contains application configuration
type Config struct {
	// Actives
	ActiveWorker int

	// Minimum App versions
	AndroidPOSAppMinimumVersion string
	IosPOSAppMinimumVersion     string

	// DB
	DBConnectionString string

	// Misc
	AppURL             string
	PortApps           string
	ServerName         string
	WhitelistedIps     string
	ConfigFileLocation string

	// Redis
	RedisAddr     string
	RedisDB       int
	RedisPassword string
	RedisTimeOut  int

	JwtTimeOut int

	// Vultr
	VultrAccessKey string
	VultrBucket    string
	VultrHostname  string
	VultrSecretKey string
	VultrRegion    string
}

var config *Config

func getEnvOrDefault(env string, defaultVal string) string {
	e := os.Getenv(env)
	if e == "" {
		return defaultVal
	}

	return e
}

// GetConfiguration , get application configuration based on set environment
func GetConfiguration() (*Config, error) {
	if config != nil {
		return config, nil
	}

	dataENV, err := ioutil.ReadFile(getEnvOrDefault(configFileLocation, ".env"))
	if err != nil {
		fmt.Println("File reading error", err)
		return nil, fmt.Errorf("failed to locate env file: %v", err)
	}

	var result map[string]interface{}
	json.Unmarshal(dataENV, &result)

	redisDBi, err := strconv.Atoi(result[redisDB].(string))
	if err != nil {
		return nil, fmt.Errorf("failed to parse redis db: %v", err)
	}

	redisTimeOut, err := strconv.Atoi(result[redisTimeOut].(string))
	if err != nil {
		return nil, fmt.Errorf("failed to parse redis timeout: %v", err)
	}

	jwtTimeOut, err := strconv.Atoi(result[jwtTimeOut].(string))
	if err != nil {
		return nil, fmt.Errorf("failed to parse jwt timeout: %v", err)
	}

	activeWorker, err := strconv.Atoi(result[activeWorker].(string))
	if err != nil {
		return nil, fmt.Errorf("failed to parse active worker: %v", err)
	}

	config := &Config{
		ActiveWorker: activeWorker,

		AndroidPOSAppMinimumVersion: result[androidPOSAppMinimumVersion].(string),
		IosPOSAppMinimumVersion:     result[iosPOSAppMinimumVersion].(string),

		DBConnectionString: result[dbConnectionString].(string),

		AppURL:         result[appUrl].(string),
		PortApps:       result[portApps].(string),
		WhitelistedIps: result[whitelistedIps].(string),

		RedisAddr:     result[redisAddr].(string),
		RedisDB:       redisDBi,
		RedisPassword: result[redisPassword].(string),
		RedisTimeOut:  redisTimeOut,

		JwtTimeOut: jwtTimeOut,

		VultrAccessKey: result[vultrAccessKey].(string),
		VultrBucket:    result[vultrBucket].(string),
		VultrHostname:  result[vultrHostname].(string),
		VultrSecretKey: result[vultrSecretKey].(string),
		VultrRegion:    result[vultrRegion].(string),
	}

	return config, nil
}
