// Package cxwh contains an example Dialogflow CX webhook
package cxwh

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	cx "cloud.google.com/go/dialogflow/cx/apiv3beta1/cxpb"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
)

// webhookResponse handles webhook calls.
func webhookResponse(request cx.WebhookRequest) (cx.WebhookResponse, error) {

	// Create session parameters that are populated in the response as an example.
	p := map[string]*structpb.Value{
		"key": {Kind: &structpb.Value_StringValue{StringValue: "value"}},
	}

	// Example reply from webhook
	messages := []*cx.ResponseMessage{
		{
			Message: &cx.ResponseMessage_Text_{
				Text: &cx.ResponseMessage_Text{
					Text: []string{"hi from the webhook!"},
				},
			},
		},
	}

	// Build and return the response.
	response := cx.WebhookResponse{
		FulfillmentResponse: &cx.WebhookResponse_FulfillmentResponse{
			Messages:      messages,
			MergeBehavior: cx.WebhookResponse_FulfillmentResponse_REPLACE,
		},
		SessionInfo: &cx.SessionInfo{
			Parameters: p,
		},
	}
	return response, nil
}

// handleError handles internal errors.
func handleError(w http.ResponseWriter, err error) {
	fmt.Printf("ERROR: %v", err)
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, "ERROR: %v", err)
}

// HandleWebhookRequest handles WebhookRequest and sends the WebhookResponse.
func HandleWebhookRequest(w http.ResponseWriter, r *http.Request) {
	var request cx.WebhookRequest
	var response cx.WebhookResponse
	var err error

	// Read input JSON
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		handleError(w, err)
		return
	}

	unmarshal := protojson.UnmarshalOptions{
		AllowPartial:   true,
		DiscardUnknown: true,
	}

	// Unmarshal from the protobuf struct
	if err = unmarshal.Unmarshal(body, &request); err != nil {
		handleError(w, err)
		return
	}

	log.Printf("Request: %+v", request)

	// Execute the response
	response, err = webhookResponse(request)

	if err != nil {
		handleError(w, err)
		return
	}
	log.Printf("Response: %+v", response)

	marshal := protojson.MarshalOptions{
		AllowPartial: true,
	}

	// Marshal from the protobuf struct
	bytes, err := marshal.Marshal(&response)
	if err != nil {
		handleError(w, err)
		return
	}

	// Send response
	if _, err = w.Write(bytes); err != nil {
		handleError(w, err)
		return
	}

}
