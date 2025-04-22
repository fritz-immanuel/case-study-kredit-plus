package appcontext

import (
	"fmt"
	"reflect"

	"github.com/gin-gonic/gin"
)

type contextKey string

const (
	// KeyURLPath represents the url path key in http server context
	KeyURLPath contextKey = "URLPath"

	// KeyHTTPMethodName represents the method name key in http server context
	KeyHTTPMethodName contextKey = "HTTPMethodName"

	// KeySessionID represents the current logged-in SessionID
	KeySessionID contextKey = "SessionID"

	//KeyWarehouseIDs represents the list access of warehouseId
	KeyOutletID contextKey = "OutletID"

	// KeyUserID represents the current logged-in UserID
	KeyUserID contextKey = "UserID"

	// KeyUserName represents the current logged-in UserID
	KeyUserName contextKey = "UserName"

	// KeyLoginToken represents the current logged-in token
	KeyLoginToken contextKey = "LoginToken"

	// KeyCustomerID represents the current logged-in UserID's CustomerID from customer-payfazz
	KeyBusinessID contextKey = "BusinessID"

	// KeyBusinessShiftID represents the current logged-in UserID's CustomerID from customer-payfazz
	KeyBusinessShiftID contextKey = "BusinessShiftID"

	// KeySupervisorUserID represents the current logged-in UserID's SupervisorUserID
	KeySupervisorUserID contextKey = "SupervisorUserID"

	// KeyKitchenID represents the current logged-in KitchenID's CustomerID from customer-payfazz
	KeyKitchenID contextKey = "KitchenID"

	// Type represents the current logged-in UserID's CustomerID from customer-payfazz
	KeyType contextKey = "Type"

	// KeyVersionCode represents the current version code of request
	KeyVersionCode contextKey = "VersionCode"

	// KeyCurrentXAccessToken represents the current access token of request
	KeyCurrentXAccessToken contextKey = "CurrentAccessToken"

	KeyRequestStatus contextKey = "KeyRequestStatus"

	// KeyRequestHeader represents the header of the request
	KeyRequestHeader contextKey = "KeyRequestHeader"

	// KeyRequestBody represents the body of the request
	KeyRequestBody contextKey = "KeyRequestBody"
)

// RequestStatus gets request status from context
func RequestStatus(ctx *gin.Context) *string {
	requestStatus := (*ctx).Value(KeyRequestStatus)
	if requestStatus != nil {
		v := requestStatus.(string)
		return &v
	}
	return nil
}

// RequestHeader gets client request header
func RequestHeader(ctx *gin.Context) string {
	requestHeader := (*ctx).Value(KeyRequestHeader)
	if requestHeader != nil {
		v := requestHeader.(string)
		return v
	}
	return ""
}

// RequestBody gets client request body
func RequestBody(ctx *gin.Context) interface{} {
	requestBody := (*ctx).Value(KeyRequestBody)
	if requestBody != nil {
		v := requestBody.(interface{})
		return v
	}
	return nil
}

// URLPath gets the data url path from the context
func URLPath(ctx *gin.Context) *string {
	urlPath := ctx.Value(fmt.Sprintf("%s", KeyURLPath))
	if urlPath != nil {
		v := urlPath.(string)
		return &v
	}
	return nil
}

// HTTPMethodName gets the data http method from the context
func HTTPMethodName(ctx *gin.Context) *string {
	httpMethodName := ctx.Value(fmt.Sprintf("%s", KeyHTTPMethodName))
	if httpMethodName != nil {
		v := httpMethodName.(string)
		return &v
	}
	return nil
}

// SessionID gets the data session id from the context
func SessionID(ctx *gin.Context) *string {
	sessionID := ctx.Value(fmt.Sprintf("%s", KeySessionID))
	if sessionID != nil {
		v := sessionID.(string)
		return &v
	}
	return nil
}

// UserID gets current userId logged in from the context
func UserID(ctx *gin.Context) *string {
	userID := ctx.Value(fmt.Sprintf("%v", KeyUserID))
	if userID != nil {
		if reflect.ValueOf(userID).Kind().String() == "string" {
			v := userID.(string)
			return &v
		} else {
			v := fmt.Sprintf("%v", userID)
			return &v
		}
	}
	return nil
}

// UserName gets current userName logged in from the context
func UserName(ctx *gin.Context) *string {
	userID := ctx.Value(fmt.Sprintf("%v", KeyUserName))
	if userID != nil {
		if reflect.ValueOf(userID).Kind().String() == "string" {
			v := userID.(string)
			return &v
		} else {
			v := fmt.Sprintf("%v", userID)
			return &v
		}
	}
	return nil
}

// TypeID gets current TypeID logged in from the context
func Type(ctx *gin.Context) *string {
	typeData := ctx.Value(fmt.Sprintf("%s", KeyType))

	if typeData != nil {
		v := typeData.(string)
		return &v
	}
	return nil
}

// OutletID gets current logged-in UserID's OutletID from context
func OutletID(ctx *gin.Context) int {
	outletID := ctx.Value(fmt.Sprintf("%v", KeyOutletID))
	if outletID != nil {
		v := int(outletID.(float64))
		return v
	}
	return 0
}

// BusinessID gets current prefered BusinessID of UserID
func BusinessID(ctx *gin.Context) int {
	businessID := ctx.Value(fmt.Sprintf("%s", KeyBusinessID))
	if businessID != nil {
		v := int(businessID.(float64))
		return v
	}
	return 0
}

// BusinessShiftID gets current prefered BusinessShiftID of UserID
func BusinessShiftID(ctx *gin.Context) int {
	businessShiftID := ctx.Value(fmt.Sprintf("%s", KeyBusinessShiftID))
	if businessShiftID != nil {
		v := int(businessShiftID.(float64))
		return v
	}
	return 0
}

// SupervisorUserID gets current prefered SupervisorUserID of UserID
func SupervisorUserID(ctx *gin.Context) int {
	SupervisorUserID := ctx.Value(fmt.Sprintf("%s", KeySupervisorUserID))
	if SupervisorUserID != nil {
		v := int(SupervisorUserID.(float64))
		return v
	}
	return 0
}

// KitchenID gets current prefered KitchenID of UserID
func KitchenID(ctx *gin.Context) int {
	KitchenID := ctx.Value(fmt.Sprintf("%s", KeyKitchenID))
	if KitchenID != nil {
		v := int(KitchenID.(float64))
		return v
	}
	return 0
}

// VersionCode gets current version code of request
func VersionCode(ctx *gin.Context) int {
	versionCode := ctx.Value(fmt.Sprintf("%s", KeyVersionCode))
	if versionCode != nil {
		v := int(versionCode.(float64))
		return v
	}
	return 0
}

// CurrentXAccessToken gets current x access token code of request
func CurrentXAccessToken(ctx *gin.Context) string {
	currentAccessToken := ctx.Value(fmt.Sprintf("%s", KeyCurrentXAccessToken))
	if currentAccessToken != nil {
		v := currentAccessToken.(string)
		return v
	}
	return ""
}
