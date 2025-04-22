package library

import (
	"log"
	"time"

	"case-study-kredit-plus/configs"
	"case-study-kredit-plus/library/appcontext"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type Credential struct {
	ID       string `json:"ID"`
	Username string `json:"Username"`
	Email    string `json:"Email"`
	Type     string `json:"Type"`

	FsId         string `json:"fsid"`
	ClientId     string `json:"clientid"`
	ClientSecret string `json:"clientsecret"`
	RefreshToken string `json:"refreshtoken"`
}

type CredentialMobile struct {
	ID       uint   `json:"ID"`
	Username string `json:"Username"`
	Email    string `json:"Email"`
	Type     string `json:"Type"`

	FsId         string `json:"fsid"`
	ClientId     string `json:"clientid"`
	ClientSecret string `json:"clientsecret"`
	RefreshToken string `json:"refreshtoken"`
}

const JwtSalt = "secret"

func JwtSignString(c Credential) (string, error) {
	config, _ := configs.GetConfiguration()

	sign := jwt.New(jwt.GetSigningMethod("HS256"))
	claims := sign.Claims.(jwt.MapClaims)

	claims["ID"] = c.ID
	claims["Email"] = c.Email
	claims["LoginTime"] = UTCPlus7()
	claims["Exp"] = UTCPlus7().Add(time.Duration(config.JwtTimeOut) * time.Second)
	claims["Type"] = c.Type

	// redisClient := redis.NewClient(&redis.Options{
	// 	Addr:     config.RedisAddr,
	// 	Password: config.RedisPassword,
	// 	DB:       config.RedisDB,
	// })

	token, err := sign.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}

	// if errRedis := redisClient.Set(
	// 	token,
	// 	fmt.Sprintf(`{"id":%s}`, c.ID),
	// 	time.Second*time.Duration(config.RedisTimeOut),
	// ).Err(); errRedis != nil {
	// 	log.Printf(`
	// 	======================================================================
	// 	Error Storing Caching in "Auth":
	// 	Error: %v,
	// 	======================================================================
	// 	`, errRedis)
	// 	return "", errRedis
	// }

	return token, nil
}

func JwtSignMobileString(c Credential) (string, error) {
	config, _ := configs.GetConfiguration()

	sign := jwt.New(jwt.GetSigningMethod("HS256"))
	claims := sign.Claims.(jwt.MapClaims)

	claims["ID"] = c.ID
	claims["Email"] = c.Email
	claims["LoginTime"] = UTCPlus7()
	claims["Exp"] = UTCPlus7().Add(time.Duration(config.JwtTimeOut) * time.Second)
	claims["Type"] = c.Type

	// redisClient := redis.NewClient(&redis.Options{
	// 	Addr:     config.RedisAddr,
	// 	Password: config.RedisPassword,
	// 	DB:       config.RedisDB,
	// })

	token, err := sign.SignedString([]byte("secretmobile"))
	if err != nil {
		return "", err
	}

	// if errRedis := redisClient.Set(
	// 	token,
	// 	fmt.Sprintf(`{\"id\":%s}`, c.ID),
	// 	time.Second*time.Duration(config.RedisTimeOut),
	// ).Err(); errRedis != nil {
	// 	log.Printf(`
	// 	======================================================================
	// 	Error Storing Caching in "Auth":
	// 	Error: %v,
	// 	======================================================================
	// 	`, errRedis)
	// 	return "", errRedis
	// }

	return token, nil
}

func extractClaims(tokenStr string) (jwt.MapClaims, bool) {
	hmacSecretString := "secret" // Value
	hmacSecret := []byte(hmacSecretString)
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// check token signing method etc
		return hmacSecret, nil
	})

	if err != nil {
		return nil, false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, true
	} else {
		log.Printf("Invalid JWT Token")
		return nil, false
	}
}

func extractMobileClaims(tokenStr string) (jwt.MapClaims, bool) {
	hmacSecretString := "secretmobile" // Value
	hmacSecret := []byte(hmacSecretString)
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// check token signing method etc
		return hmacSecret, nil
	})

	if err != nil {
		return nil, false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, true
	} else {
		log.Printf("Invalid JWT Token")
		return nil, false
	}
}

func GetJWTClaims(ctx *gin.Context, token string) (jwt.MapClaims, bool) {
	var claims jwt.MapClaims
	var ok bool
	if token == "" {
		JwtActiveToken := appcontext.SessionID(ctx)
		claims, ok = extractClaims(*JwtActiveToken)
	} else {
		JwtActiveToken := token
		claims, ok = extractClaims(JwtActiveToken)
	}

	return claims, ok
}

func GetJWTMobileClaims(ctx *gin.Context, token string) (jwt.MapClaims, bool) {
	var claims jwt.MapClaims
	var ok bool
	if token == "" {
		JwtActiveToken := appcontext.SessionID(ctx)
		claims, ok = extractMobileClaims(*JwtActiveToken)
	} else {
		JwtActiveToken := token
		claims, ok = extractMobileClaims(JwtActiveToken)
	}

	return claims, ok
}

func GetJWTClaimsMock() jwt.MapClaims {
	var ctx *gin.Context
	claims, _ := GetJWTClaims(ctx, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJCdXNpbmVzc0lEIjoxLCJFbWFpbCI6ImppbW15QHNlcGFyaW5kby5jb20iLCJFeHAiOjE1NzYyMzg1NTAsIklEIjoyLCJSb2xlSUQiOjEsIlRyaXBJRCI6MCwiVXNlcm5hbWUiOiJqaW1teSJ9.RvdZ6I7VTSspCnsvQflBgwrCVwUtENGu846CQqgcSh4")
	return claims
}

func SetJwtClaimsMock() {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJCdXNpbmVzc0lEIjoxLCJFbWFpbCI6ImppbW15QHNlcGFyaW5kby5jb20iLCJFeHAiOjE1NzYyMzg1NTAsIklEIjoyLCJSb2xlSUQiOjEsIlRyaXBJRCI6MCwiVXNlcm5hbWUiOiJqaW1teSJ9.RvdZ6I7VTSspCnsvQflBgwrCVwUtENGu846CQqgcSh4"
	configs.JwtActiveToken = &token
}
