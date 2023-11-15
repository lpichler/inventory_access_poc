// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.15.0 DO NOT EDIT.
package api

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/chi/v5"
	"github.com/oapi-codegen/runtime"
	strictnethttp "github.com/oapi-codegen/runtime/strictmiddleware/nethttp"
)

// GetContentPackagesParams defines parameters for GetContentPackages.
type GetContentPackagesParams struct {
	// Page Page number for packages.
	Page *int `form:"page,omitempty" json:"page,omitempty"`

	// Matches Matches all packages whose name contains this parameter as a substring
	Matches *string `form:"matches,omitempty" json:"matches,omitempty"`

	// PatchesAvailable Filter for packages for which there are systems with patches available (true) or already up to date (false).
	PatchesAvailable *bool `form:"patches_available,omitempty" json:"patches_available,omitempty"`

	// Tag Filter package by tag
	Tag *string `form:"tag,omitempty" json:"tag,omitempty"`

	// SortKey The package key to sort by.
	SortKey *string `form:"sort_key,omitempty" json:"sort_key,omitempty"`

	// SortOrder Sorting ascending (true) or descending (false).
	SortOrder *bool `form:"sort_order,omitempty" json:"sort_order,omitempty"`
}

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (GET /content/packages)
	GetContentPackages(w http.ResponseWriter, r *http.Request, params GetContentPackagesParams)
}

// Unimplemented server implementation that returns http.StatusNotImplemented for each endpoint.

type Unimplemented struct{}

// (GET /content/packages)
func (_ Unimplemented) GetContentPackages(w http.ResponseWriter, r *http.Request, params GetContentPackagesParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.Handler) http.Handler

// GetContentPackages operation middleware
func (siw *ServerInterfaceWrapper) GetContentPackages(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params GetContentPackagesParams

	// ------------- Optional query parameter "page" -------------

	err = runtime.BindQueryParameter("form", true, false, "page", r.URL.Query(), &params.Page)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "page", Err: err})
		return
	}

	// ------------- Optional query parameter "matches" -------------

	err = runtime.BindQueryParameter("form", true, false, "matches", r.URL.Query(), &params.Matches)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "matches", Err: err})
		return
	}

	// ------------- Optional query parameter "patches_available" -------------

	err = runtime.BindQueryParameter("form", true, false, "patches_available", r.URL.Query(), &params.PatchesAvailable)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "patches_available", Err: err})
		return
	}

	// ------------- Optional query parameter "tag" -------------

	err = runtime.BindQueryParameter("form", true, false, "tag", r.URL.Query(), &params.Tag)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "tag", Err: err})
		return
	}

	// ------------- Optional query parameter "sort_key" -------------

	err = runtime.BindQueryParameter("form", true, false, "sort_key", r.URL.Query(), &params.SortKey)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "sort_key", Err: err})
		return
	}

	// ------------- Optional query parameter "sort_order" -------------

	err = runtime.BindQueryParameter("form", true, false, "sort_order", r.URL.Query(), &params.SortOrder)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "sort_order", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetContentPackages(w, r, params)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

type UnescapedCookieParamError struct {
	ParamName string
	Err       error
}

func (e *UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter '%s'", e.ParamName)
}

func (e *UnescapedCookieParamError) Unwrap() error {
	return e.Err
}

type UnmarshalingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshalingParamError) Error() string {
	return fmt.Sprintf("Error unmarshaling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshalingParamError) Unwrap() error {
	return e.Err
}

type RequiredParamError struct {
	ParamName string
}

func (e *RequiredParamError) Error() string {
	return fmt.Sprintf("Query argument %s is required, but not found", e.ParamName)
}

type RequiredHeaderError struct {
	ParamName string
	Err       error
}

func (e *RequiredHeaderError) Error() string {
	return fmt.Sprintf("Header parameter %s is required, but not found", e.ParamName)
}

func (e *RequiredHeaderError) Unwrap() error {
	return e.Err
}

type InvalidParamFormatError struct {
	ParamName string
	Err       error
}

func (e *InvalidParamFormatError) Error() string {
	return fmt.Sprintf("Invalid format for parameter %s: %s", e.ParamName, e.Err.Error())
}

func (e *InvalidParamFormatError) Unwrap() error {
	return e.Err
}

type TooManyValuesForParamError struct {
	ParamName string
	Count     int
}

func (e *TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("Expected one value for %s, got %d", e.ParamName, e.Count)
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{})
}

type ChiServerOptions struct {
	BaseURL          string
	BaseRouter       chi.Router
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, r chi.Router) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseRouter: r,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, r chi.Router, baseURL string) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseURL:    baseURL,
		BaseRouter: r,
	})
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options ChiServerOptions) http.Handler {
	r := options.BaseRouter

	if r == nil {
		r = chi.NewRouter()
	}
	if options.ErrorHandlerFunc == nil {
		options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandlerFunc:   options.ErrorHandlerFunc,
	}

	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/content/packages", wrapper.GetContentPackages)
	})

	return r
}

type GetContentPackagesRequestObject struct {
	Params GetContentPackagesParams
}

type GetContentPackagesResponseObject interface {
	VisitGetContentPackagesResponse(w http.ResponseWriter) error
}

// StrictServerInterface represents all server handlers.
type StrictServerInterface interface {

	// (GET /content/packages)
	GetContentPackages(ctx context.Context, request GetContentPackagesRequestObject) (GetContentPackagesResponseObject, error)
}

type StrictHandlerFunc = strictnethttp.StrictHttpHandlerFunc
type StrictMiddlewareFunc = strictnethttp.StrictHttpMiddlewareFunc

type StrictHTTPServerOptions struct {
	RequestErrorHandlerFunc  func(w http.ResponseWriter, r *http.Request, err error)
	ResponseErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

func NewStrictHandler(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc) ServerInterface {
	return &strictHandler{ssi: ssi, middlewares: middlewares, options: StrictHTTPServerOptions{
		RequestErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		},
		ResponseErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		},
	}}
}

func NewStrictHandlerWithOptions(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc, options StrictHTTPServerOptions) ServerInterface {
	return &strictHandler{ssi: ssi, middlewares: middlewares, options: options}
}

type strictHandler struct {
	ssi         StrictServerInterface
	middlewares []StrictMiddlewareFunc
	options     StrictHTTPServerOptions
}

// GetContentPackages operation middleware
func (sh *strictHandler) GetContentPackages(w http.ResponseWriter, r *http.Request, params GetContentPackagesParams) {
	var request GetContentPackagesRequestObject

	request.Params = params

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.GetContentPackages(ctx, request.(GetContentPackagesRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetContentPackages")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(GetContentPackagesResponseObject); ok {
		if err := validResponse.VisitGetContentPackagesResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/5RT32vbTgz/V4SeWjBpv9+95W3r2OhgW6B9G6PIPsU+er67SXKCKf7fx9lpSrum0Kfk",
	"dNLnl08P6OM24foBHWsjPptPEdd4y2rA0eXkoylsk8COgndkPrbw7ef1D9CGI4lPCjXbnjkCNQ2rgnKZ",
	"kNTDTfYNf/4EFB1sfTAWdhVoEmM3FzO17JYBH8HHHUdLMq6wQvMWGNe4SVfQpKgpcPk1jvbUCB8311jh",
	"jkUX3f+tLleXOFWYMkfKHtf4YXW5+h8rzGSdFqcXB5iLTM09tTwXW7Z/U/jKBo9NQDvygerAcxqZrOnm",
	"k45q3GvRnDILldlrt0xfLUybR6KiQqhnY1Fc/3rJt6GWIQ59zXIgWeYKti8NfwaWESuM1JdsSnxYoTYd",
	"91Tk25hL3UfjlgWnqXpJ8b3oLm5CeLK275IyFNA5YvJRwTqvcFQLpECgQ60mPrYn9PQL+GuSDnOvKPoy",
	"P4xnfufDvvNNB9axMJAcc4a9t26J/9lHOTMZ+BySAAVhciMMGSyBI2M421JQPj+d44x2d0R7zUGdUmCK",
	"b1g4yId6BKNTGS0378jntuMj8j2PxVPZIKjHU3bK9d09j+/juUkyLzeVxXbl31OkpfOx+HaUM3cSx/J2",
	"hr8rFNacopb9i0MIU4U69D3JiGs8LM7xReA0TX8DAAD///VW2n6sBAAA",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
