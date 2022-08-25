// Package cxwh contains an example Dialogflow CX webhook
package cxwh

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	cx "google.golang.org/genproto/googleapis/cloud/dialogflow/cx/v3beta1"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
)

// webhookRequest is used to unmarshal a WebhookRequest JSON object. Note that
// not all members need to be defined--just those that you need to process.
// As an alternative, you could use the types provided by the Dialogflow protocol buffers:
// https://pkg.go.dev/google.golang.org/genproto/googleapis/cloud/dialogflow/cx/v3#WebhookRequest

// webhookResponse handles webhook calls.
func webhookResponse(request cx.WebhookRequest) (cx.WebhookResponse, error) {

	// Create session parameters that are populated in the response.
	// The "cancel-period" parameter is referenced by the agent.
	// This example hard codes the value 2, but a real system
	// might look up this value in a database.
	p := map[string]*structpb.Value{
		"key": {Kind: &structpb.Value_StringValue{StringValue: "value"}},
	}

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

	// Send response
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
