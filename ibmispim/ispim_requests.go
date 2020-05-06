package ibmispim

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

const IspimRequestPath = "/ispim/rest/requests"

type IspimRequestsService interface {
	Get(context.Context, int) (*IspimRequestResponse, *Response, error)
}

type IspimRequestsServiceOp struct {
	client *Client
}

/// START of REQUEST Structure

/// End of the Request

///Start of the response
type IspimRequestResponse struct {
	Links                  Links      `json:"_links"`
	IspimRequestAttributes Attributes `json:"_attributes"`
}
type Children struct {
	Href string `json:"href"`
}

type Status struct {
	Key  string `json:"key"`
	Text string `json:"text"`
}
type IspimRequestAttributes struct {
	Result        []string `json:"result"`
	Completedate  []string `json:"completedate"`
	Status        []Status `json:"status"`
	Subject       []string `json:"subject"`
	Requestername []string `json:"requestername"`
	Scheduleddate []string `json:"scheduleddate"`
	Lastmodified  []string `json:"lastmodified"`
	Resultdetail  []string `json:"resultdetail"`
	Type          []string `json:"type"`
	Justification []string `json:"justification"`
	Starteddate   []string `json:"starteddate"`
	Submitteddate []string `json:"submitteddate"`
}

///End of the response

var _ IspimRequestsService = &IspimRequestsServiceOp{}

// Get the Configuration for the IspimRequestsServiceOp created

func (ispims *IspimRequestsServiceOp) Get(ctx context.Context, ispimId int) (*IspimRequestResponse, *Response, error) {
	requestNumber := strconv.Itoa(ispimId)

	log.Printf("[DEBUG] Printing the from the GET CALL requestNumber %s", requestNumber)

	getPath := fmt.Sprintf("%s%s", IspimRequestPath, requestNumber)

	log.Printf("[DEBUG] In the get method - the path is %s", getPath)

	req, err := ispims.client.NewRequest(ctx, http.MethodGet, getPath, nil)
	if err != nil {
		return nil, nil, err
	} else {
		log.Printf("[DEBUG]: ISIM-Provider:  Successfully executed the GetRequest Call for IspimRequest")

	}

	root := IspimRequestResponse{}
	resp, err := ispims.client.Do(ctx, req, &root)
	if err != nil {
		log.Printf("[DEBUG]. If we are here , the do call has failed")
		log.Fatal(err)
		return nil, resp, err
	}

	return &root, resp, err
}

// Delete - Add code to disconnect credentials - We do a lookup and then perform the disconnect -
