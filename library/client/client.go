package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"

	//"case-study-kredit-plus/library/appcontext"
	"strings"

	"case-study-kredit-plus/library/types"

	"github.com/go-redis/redis"
)

// Method represents the enum for http call method
type Method string

// Enum value for http call method
const (
	POST   Method = "POST"
	PUT    Method = "PUT"
	DELETE Method = "DELETE"
	GET    Method = "GET"
	PATCH  Method = "PATCH"
)

// ResponseError represents struct of Authorization Type
type ResponseError struct {
	Code       string         `json:"code"`
	Message    string         `json:"message"`
	Fields     types.Metadata `json:"-"`
	StatusCode int            `json:"statusCode"`
	Error      error          `json:"error"`
}

// AuthorizationTypeStruct represents struct of Authorization Type
type AuthorizationTypeStruct struct {
	HeaderName      string
	HeaderType      string
	HeaderTypeValue string
	Token           string
}

// AuthorizationType represents the enum for http authorization type
type AuthorizationType AuthorizationTypeStruct

// Enum value for http authorization type
var (
	Basic        = AuthorizationType(AuthorizationTypeStruct{HeaderName: "Authorization", HeaderType: "Basic", HeaderTypeValue: "Basic "})
	Bearer       = AuthorizationType(AuthorizationTypeStruct{HeaderName: "Authorization", HeaderType: "Bearer", HeaderTypeValue: "Bearer "})
	AccessToken  = AuthorizationType(AuthorizationTypeStruct{HeaderName: "Access-Token", HeaderType: "Bearer", HeaderTypeValue: "Bearer "})
	Secret       = AuthorizationType(AuthorizationTypeStruct{HeaderName: "Secret", HeaderType: "Secret", HeaderTypeValue: ""})
	APPKey       = AuthorizationType(AuthorizationTypeStruct{HeaderName: "APP_KEY", HeaderType: "APP_KEY", HeaderTypeValue: ""})
	DeviceID     = AuthorizationType(AuthorizationTypeStruct{HeaderName: "DEVICE_ID", HeaderType: "DEVICE_ID", HeaderTypeValue: ""})
	FSID         = AuthorizationType(AuthorizationTypeStruct{HeaderName: "FSID", HeaderType: "Basic", HeaderTypeValue: ""})
	ClientID     = AuthorizationType(AuthorizationTypeStruct{HeaderName: "ClientID", HeaderType: "Basic", HeaderTypeValue: ""})
	ClientSecret = AuthorizationType(AuthorizationTypeStruct{HeaderName: "ClientSecret", HeaderType: "Basic", HeaderTypeValue: ""})
)

//
// Private constants
//

const apiURL = "https://127.0.0.1:8080"
const defaultHTTPTimeout = 80 * time.Second
const maxNetworkRetriesDelay = 5000 * time.Millisecond
const minNetworkRetriesDelay = 500 * time.Millisecond

//
// Private variables
//

var httpClient = &http.Client{Timeout: defaultHTTPTimeout}

// GenericHTTPClient represents an interface to generalize an object to implement HTTPClient
type GenericHTTPClient interface {
	Do(req *http.Request) (string, *ResponseError)
	CallClient(ctx *gin.Context, path string, method Method, request interface{}, result interface{}, isAcknowledgeNeeded bool) *ResponseError
	CallClientWithCaching(ctx *gin.Context, path string, method Method, request interface{}, result interface{}, isAcknowledgeNeeded bool) *ResponseError
	CallClientWithCachingInRedis(ctx *gin.Context, durationInSecond int, path string, method Method, request interface{}, result interface{}, isAcknowledgeNeeded bool) *ResponseError
	CallClientWithCircuitBreaker(ctx *gin.Context, path string, method Method, request interface{}, result interface{}, isAcknowledgeNeeded bool) *ResponseError
	CallClientWithoutLog(ctx *gin.Context, path string, method Method, request interface{}, result interface{}, isAcknowledgeNeeded bool) *ResponseError
	CallClientWithBaseURLGiven(ctx *gin.Context, url string, method Method, request interface{}, result interface{}, isAcknowledgeNeeded bool) *ResponseError
	CallClientWithCustomizedError(ctx *gin.Context, path string, method Method, queryParams interface{}, request interface{}, result interface{}, isAcknowledgeNeeded bool) *ResponseError
	CallClientWithCustomizedErrorAndCaching(ctx *gin.Context, path string, method Method, queryParams interface{}, request interface{}, result interface{}, isAcknowledgeNeeded bool) *ResponseError
	AddAuthentication(ctx *gin.Context, authorizationType AuthorizationType)
}

// HTTPClient represents the service http client
type HTTPClient struct {
	clientRequestLogStorage ClientRequestLogStorage
	clientCacheService      ClientCacheServiceInterface
	redisClient             *redis.Client
	APIURL                  string
	HTTPClient              *http.Client
	MaxNetworkRetries       int
	UseNormalSleep          bool
	AuthorizationTypes      []AuthorizationType
	ClientName              string
}

func (c *HTTPClient) shouldRetry(err error, res *http.Response, retry int) bool {
	if retry >= c.MaxNetworkRetries {
		return false
	}

	if err != nil {
		return true
	}

	return false
}

func (c *HTTPClient) sleepTime(numRetries int) time.Duration {
	if c.UseNormalSleep {
		return 0
	}

	// exponentially backoff by 2^numOfRetries
	delay := minNetworkRetriesDelay + minNetworkRetriesDelay*time.Duration(1<<uint(numRetries))
	if delay > maxNetworkRetriesDelay {
		delay = maxNetworkRetriesDelay
	}

	// generate random jitter to prevent thundering herd problem
	jitter := rand.Int63n(int64(delay / 4))
	delay -= time.Duration(jitter)

	if delay < minNetworkRetriesDelay {
		delay = minNetworkRetriesDelay
	}

	return delay
}

// Do calls the api http request and parse the response into v
func (c *HTTPClient) Do(req *http.Request) (string, *ResponseError) {
	var res *http.Response
	var err error

	for retry := 0; ; {
		res, err = c.HTTPClient.Do(req)

		if !c.shouldRetry(err, res, retry) {
			break
		}

		sleepDuration := c.sleepTime(retry)
		retry++

		time.Sleep(sleepDuration)
	}
	if err != nil {
		return "", &ResponseError{
			Code:    "",
			Message: "",
			Fields:  nil,
			Error:   err,
		}
	}
	defer res.Body.Close()

	resBody, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return "", &ResponseError{
			Code:       string(res.StatusCode),
			Message:    "",
			Fields:     nil,
			StatusCode: res.StatusCode,
			Error:      err,
		}
	}

	errResponse := &ResponseError{
		Code:       string(res.StatusCode),
		Message:    "",
		Fields:     nil,
		StatusCode: res.StatusCode,
		Error:      nil,
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		err = json.Unmarshal([]byte(string(resBody)), errResponse)
		if err != nil {
			errResponse.Error = err
		}
		errResponse.Error = fmt.Errorf("Error while calling %s: %v", req.URL.String(), errResponse.Message)

		return "", errResponse
	}

	return string(resBody), errResponse
}

// CallClient do call client
func (c *HTTPClient) CallClient(ctx *gin.Context, path string, method Method, request interface{}, result interface{}, isAcknowledgeNeeded bool) *ResponseError {
	var jsonData []byte
	var err error
	var response string
	var errDo *ResponseError

	if request != nil && request != "" {
		jsonData, err = json.Marshal(request)
		if err != nil {
			errDo = &ResponseError{
				Error: err,
			}
			return errDo
		}
	}

	urlPath, err := url.Parse(fmt.Sprintf("%s/%s", c.APIURL, path))
	if err != nil {
		errDo = &ResponseError{
			Error: err,
		}
		return errDo
	}

	req, err := http.NewRequest(string(method), urlPath.String(), bytes.NewBuffer(jsonData))
	if err != nil {
		errDo = &ResponseError{
			Error: err,
		}
		return errDo
	}

	for _, authorizationType := range c.AuthorizationTypes {
		if authorizationType.HeaderType != "APIKey" {
			req.Header.Add(authorizationType.HeaderName, fmt.Sprintf("%s%s", authorizationType.HeaderTypeValue, authorizationType.Token))
		}
	}

	req.Header.Add("Content-Type", "application/json")

	response, errDo = c.Do(req)
	if errDo != nil && (errDo.Error != nil || errDo.Message != "") {
		return errDo
	}

	if response != "" && result != nil {
		err = json.Unmarshal([]byte(response), result)
		if err != nil {
			errDo = &ResponseError{
				Error: err,
			}
			return errDo
		}
	}

	return errDo
}

func (c *HTTPClient) CallClientFormEncode(ctx *gin.Context, path string, method Method, request url.Values, result interface{}, isAcknowledgeNeeded bool) *ResponseError {
	var response string
	var errDo *ResponseError

	urlPath, err := url.Parse(fmt.Sprintf("%s/%s", c.APIURL, path))
	if err != nil {
		errDo = &ResponseError{
			Error: err,
		}
		return errDo
	}

	req, err := http.NewRequest(string(method), urlPath.String(), strings.NewReader(request.Encode()))
	if err != nil {
		errDo = &ResponseError{
			Error: err,
		}
		return errDo
	}

	for _, authorizationType := range c.AuthorizationTypes {
		if authorizationType.HeaderType != "APIKey" {
			req.Header.Add(authorizationType.HeaderName, fmt.Sprintf("%s%s", authorizationType.HeaderTypeValue, authorizationType.Token))
		}
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	response, errDo = c.Do(req)
	if errDo != nil && (errDo.Error != nil || errDo.Message != "") {
		return errDo
	}
	if response != "" && result != nil {
		err = json.Unmarshal([]byte(response), result)
		if err != nil {
			errDo = &ResponseError{
				Error: err,
			}
			return errDo
		}
	}

	return errDo
}

// AddAuthentication do add authentication
func (c *HTTPClient) AddAuthentication(ctx *gin.Context, authorizationType AuthorizationType) {
	isExist := false
	for key, singleAuthorizationType := range c.AuthorizationTypes {
		if singleAuthorizationType.HeaderType == authorizationType.HeaderType {
			c.AuthorizationTypes[key].Token = authorizationType.Token
			isExist = true
			break
		}
	}

	if !isExist {
		c.AuthorizationTypes = append(c.AuthorizationTypes, authorizationType)
	}
}

// NewHTTPClient creates the new http client
func NewHTTPClient(
	config HTTPClient,
) *HTTPClient {
	if config.HTTPClient == nil {
		config.HTTPClient = httpClient
	}

	if config.APIURL == "" {
		config.APIURL = apiURL
	}

	return &HTTPClient{
		APIURL:             config.APIURL,
		HTTPClient:         config.HTTPClient,
		MaxNetworkRetries:  config.MaxNetworkRetries,
		UseNormalSleep:     config.UseNormalSleep,
		AuthorizationTypes: config.AuthorizationTypes,
		ClientName:         config.ClientName,
	}
}
