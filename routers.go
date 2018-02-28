package src

import (
	"net/http"
	"fmt"
	"github.com/gorilla/mux"
  "github.com/auth0/go-jwt-middleware"
  "github.com/dgrijalva/jwt-go"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc

		// TODO: Remove CreateUser from validation before deliver the product
		if route.Name != "LogIn" && route.Name != "CreateUser" {
			handler = Logger(jwtMiddleware.Handler(handler), route.Name)
		} else {
			handler = Logger(handler, route.Name)
		}

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

var jwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
  ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
    return mySigningKey, nil
  },
  SigningMethod: jwt.SigningMethodHS256,
})

var mySigningKey = []byte("secret")

var routes = Routes {
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},

	Route{
		"FilterReturn",
		"POST",
		"/filter_return",
		filterDateReturn,
	},

	Route{
		"FilterPickup",
		"POST",
		"/filter_pickup",
		filterDatePickup,
	},

	Route{
		"LogIn",
		"POST",
		"/login",
		LogIn,
	},

	Route{
		"CreateCategory",
		"POST",
		"/categories",
		CreateCategory,
	},

	Route{
		"DeleteCategory",
		"DELETE",
		"/categories/{category_uuid}",
		DeleteCategory,
	},

	Route{
		"GetCategories",
		"GET",
		"/categories",
		GetCategories,
	},

	Route{
		"GetCategoryByUUID",
		"GET",
		"/categories/{category_uuid}",
		GetCategoryByUUID,
	},

	Route{
		"UpdateCategory",
		"PUT",
		"/categories/{category_uuid}",
		UpdateCategory,
	},

	Route{
		"CreateClient",
		"POST",
		"/client",
		CreateClient,
	},

	Route{
		"DeleteClient",
		"DELETE",
		"/client/{client_uuid}",
		DeleteClient,
	},

	Route{
		"GetClientByUUID",
		"GET",
		"/client/{client_uuid}",
		GetClientByUUID,
	},

	Route{
		"IndexClients",
		"GET",
		"/client",
		IndexClients,
	},

	Route{
		"UpdateClient",
		"PUT",
		"/client/{client_uuid}",
		UpdateClient,
	},

	Route{
		"CreateHistory",
		"POST",
		"/history",
		CreateHistory,
	},

	Route{
		"DeleteHistory",
		"DELETE",
		"/history/{history_uuid}",
		DeleteHistory,
	},

	Route{
		"GetHistoryByUUID",
		"GET",
		"/history/{history_uuid}",
		GetHistoryByUUID,
	},

	Route{
		"IndexHistory",
		"GET",
		"/history",
		IndexHistory,
	},

	Route{
		"UpdateHistory",
		"PUT",
		"/history/{history_uuid}",
		UpdateHistory,
	},

	Route{
		"CreateModel",
		"POST",
		"/model",
		CreateModel,
	},

	Route{
		"DeleteModel",
		"DELETE",
		"/model/{model_uuid}",
		DeleteModel,
	},

	Route{
		"GetModelByUUID",
		"GET",
		"/model/{model_uuid}",
		GetModelByUUID,
	},

	Route{
		"IndexModels",
		"GET",
		"/model",
		IndexModels,
	},

	Route{
		"UpdateModel",
		"PUT",
		"/model/{model_uuid}",
		UpdateModel,
	},

	Route{
		"CreateRateType",
		"POST",
		"/rateType",
		CreateRateType,
	},

	Route{
		"DeleteRateType",
		"DELETE",
		"/rateType/{rateType_uuid}",
		DeleteRateType,
	},

	Route{
		"GetRateTypeByUUID",
		"GET",
		"/rateType/{rateType_uuid}",
		GetRateTypeByUUID,
	},

	Route{
		"IndexRateType",
		"GET",
		"/rateType",
		IndexRateType,
	},

	Route{
		"UpdateRateType",
		"PUT",
		"/rateType/{rateType_uuid}",
		UpdateRateType,
	},

	Route{
		"CreateRateTypeCategory",
		"POST",
		"/rateTypeCategory",
		CreateRateTypeCategory,
	},

	Route{
		"DeleteRateTypeCategory",
		"DELETE",
		"/rateTypeCategory/{rateTypeCategory_uuid}",
		DeleteRateTypeCategory,
	},

	Route{
		"GetRateTypeCategoryByUUID",
		"GET",
		"/rateTypeCategory/{rateTypeCategory_uuid}",
		GetRateTypeCategoryByUUID,
	},

	Route{
		"IndexRateTypeCategory",
		"GET",
		"/rateTypeCategory",
		IndexRateTypeCategory,
	},

	Route{
		"UpdateRateTypeCategory",
		"PUT",
		"/rateTypeCategory/{rateTypeCategory_uuid}",
		UpdateRateTypeCategory,
	},

	Route{
		"CreateRates",
		"POST",
		"/rates",
		CreateRates,
	},

	Route{
		"DeleteRates",
		"DELETE",
		"/rates/{rates_uuid}",
		DeleteRates,
	},

	Route{
		"GetRatesByUUID",
		"GET",
		"/rates/{rates_uuid}",
		GetRatesByUUID,
	},

	Route{
		"IndexRates",
		"GET",
		"/rates",
		IndexRates,
	},

	Route{
		"UpdateRates",
		"PUT",
		"/rates/{rates_uuid}",
		UpdateRates,
	},

	Route{
		"CreateReservation",
		"POST",
		"/reservation",
		CreateReservation,
	},

	Route{
		"DeleteReservation",
		"DELETE",
		"/reservation/{reservation_uuid}",
		DeleteReservation,
	},

	Route{
		"GetReservationByUUID",
		"GET",
		"/reservation/{reservation_uuid}",
		GetReservationByUUID,
	},

	Route{
		"IndexReservations",
		"GET",
		"/reservation",
		IndexReservations,
	},

	Route{
		"UpdateReservation",
		"PUT",
		"/reservation/{reservation_uuid}",
		UpdateReservation,
	},

	Route{
		"CreateUser",
		"POST",
		"/user",
		CreateUser,
	},

	Route{
		"DeleteUser",
		"DELETE",
		"/user/{user_uuid}",
		DeleteUser,
	},

	Route{
		"GetUserByUUID",
		"GET",
		"/user/{user_uuid}",
		GetUserByUUID,
	},

	Route{
		"IndexUser",
		"GET",
		"/user",
		IndexUser,
	},

	Route{
		"UpdateUser",
		"PUT",
		"/user/{user_uuid}",
		UpdateUser,
	},

}
