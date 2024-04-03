package api

import (
	"fmt"
	"strings"

	"github.com/darkqiank/who-dat/lib"
	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func sendJSONResponse(ctx *fasthttp.RequestCtx, statusCode int, data interface{}, err error) {
	response := Response{
		Success: err == nil,
	}

	if err != nil {
		response.Error = err.Error()
	}
	response.Data = data

	// Convert response to JSON
	jsonResponse, jsonErr := jsoniter.Marshal(response)
	if jsonErr != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetContentType("application/json")
		fmt.Fprintf(ctx, `{"success": false, "error": "Error marshalling response JSON"}`)
		return
	}

	// Set content type and status code
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(statusCode)

	// Send the response
	ctx.Write(jsonResponse)
}

// MainHandler handles Whois info for a single domain
func MainHandler(ctx *fasthttp.RequestCtx) {
	// Make sure it's a GET request
	if !ctx.IsGet() {
		sendJSONResponse(ctx, fasthttp.StatusMethodNotAllowed, nil, fmt.Errorf("please use a GET request"))
		return
	}

	// Extract domain from URL path
	path := string(ctx.Path())
	path = strings.TrimPrefix(path, "/")

	// Extract domain from URL path
	if path == "" {
		sendJSONResponse(ctx, fasthttp.StatusBadRequest, nil, fmt.Errorf("domain not specified"))
		return
	}

	// Initialize disableReferral as true
	disableReferral := true

	// Check if path starts with "ref/" and adjust accordingly
	if strings.HasPrefix(path, "ref/") {
		disableReferral = false
		path = strings.TrimPrefix(path, "ref/")
	}

	// Make sure path is not empty after trimming "ref/"
	if path == "" {
		sendJSONResponse(ctx, fasthttp.StatusBadRequest, nil, fmt.Errorf("domain not specified"))
		return
	}

	// Get Whois data
	whois, err := lib.GetWhois(path, disableReferral)
	if err != nil {
		sendJSONResponse(ctx, fasthttp.StatusInternalServerError, whois, err)
		return
	}

	// Get EMPTY data
	if whois.Domain == nil {
		sendJSONResponse(ctx, fasthttp.StatusNotFound, nil, fmt.Errorf("WHOIS DATA EMPTY"))
		return
	}

	// Success response
	sendJSONResponse(ctx, fasthttp.StatusOK, whois, nil)
}
