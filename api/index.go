package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/darkqiank/who-dat/lib"
	jsoniter "github.com/json-iterator/go"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func sendJSONResponse(w http.ResponseWriter, statusCode int, data interface{}, err error) {
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
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"success": false, "error": "Error marshalling response JSON"}`)
		return
	}

	// Set content type and status code
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	// Send the response
	w.Write(jsonResponse)
}

// MainHandler handles Whois info for a single domain
func MainHandler(w http.ResponseWriter, r *http.Request) {
	// Make sure it's a GET request
	if r.Method != http.MethodGet {
		sendJSONResponse(w, http.StatusMethodNotAllowed, nil, fmt.Errorf("please use a GET request"))
		return
	}

	// Extract domain from URL path
	path := strings.TrimPrefix(r.URL.Path, "/")
	if path == "" {
		sendJSONResponse(w, http.StatusBadRequest, nil, fmt.Errorf("domain not specified"))
		return
	}

	// Get Whois data
	whois, err := lib.GetWhois(path)
	if err != nil {
		sendJSONResponse(w, http.StatusInternalServerError, whois, err)
		return
	}

	// Get EMPTY data
	if whois.Domain == nil {
		sendJSONResponse(w, http.StatusNotFound, nil, fmt.Errorf("WHOIS DATA EMPTY"))
		return
	}

	// Success response
	sendJSONResponse(w, http.StatusOK, whois, nil)
}
