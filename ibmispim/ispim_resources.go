package ibmispim

import (
	"context"
	"github.com/IBM-Cloud/ispim-go-client/errors"
	"log"
	"net/http"
	//"strconv"
	"fmt"
)

const ispimresourcePath = "ispim/rest/resources"
const ispimresourceGetPath = "ispim/rest/resources?attributes=*&filter="

type IspimResourcesService interface {
	Get(context.Context, string) ([]GetIspimResourceResponse, *Response, error)
	Create(context.Context, *IspimResourceRequest) (*IspimResource, *Response, error)
	//Update(context.Context, int, *IspimResourceUpdateRequest) (*IspimResource, *Response, error)
	Delete(context.Context, string) (*Response, error)
}

type IspimResourcesServiceOp struct {
	client *Client
}

//Start of Create Resource Request

type IspimResourceRequest struct {
	CreateList []CreateList `json:"createList"`
}

type IspimResourceContainer struct {
	Href string `json:"href"`
	// Self Self `json:"self"`
}

type IspimResourceLinks struct {
	IspimResourceContainer IspimResourceContainer `json:"container"`
}

type IspimResourceAttributes struct {
	UID        []string `json:"uid"`
	Name       []string `json:"name"`
	Alias      []string `json:"alias"`
	Tag        []string `json:"tag"`
	Type       []string `json:"type,omitempty"`
	Properties []string `json:"properties,omitempty"`
}

type CreateList struct {
	IspimResourceLinks      IspimResourceLinks      `json:"_links"`
	IspimResourceAttributes IspimResourceAttributes `json:"_attributes"`
}

// End of Create Resource Request

// Start of response from the get call
type GetIspimResourceResponse struct {
	Links                      Links                      `json:"_links"`
	GetIspimResourceAttributes GetIspimResourceAttributes `json:"_attributes"`
}

type GetIspimResourceLinks struct {
	IspimResourceResponseContainer IspimResourceResponseContainer `json:"container"`
	Self                           Self                           `json:"self"`
}

type IspimResourceResponseContainer struct {
	Href  string `json:"href"`
	Title string `json:"title"`
}

type GetIspimResourceAttributes struct {
	UID   []string `json:"uid"`
	Alias []string `json:"alias"`
	Tag   []string `json:"tag"`
	Name  []string `json:"name"`
	Type  []string `json:"type"`
}

// End of response from the Get CALL

type IspimResource struct {
	ResponseList  []ResponseList `json:"responseList"`
	OverAllStatus string         `json:"overAllStatus"`
}

type IspimResourceResponse struct {
	Entity   Entity   `json:"entity"`
	Metadata Metadata `json:"metadata"`
	Status   int      `json:"status"`
}
type IspimResourceSelfLink struct {
	Href  string `json:"href"`
	Title string `json:"title"`
}

//type ResponseList struct {
//	RequestAction string   `json:"requestAction"`
//	IspimResourceResponse      IspimResourceResponse `json:"response"`
//	IspimResourceSelfLink      IspimResourceSelfLink `json:"selfLink"`
//}

// End of Response

type IspimDeleteResourceRequest struct {
}

type IspimDeleteResourceResponse struct {
}

// End of the common structure
type IspimResourceRoot struct {
	IspimResource *IspimResource `json:"IspimResource"`
}

type IspimResourcesRoot struct {
	IspimResources []IspimResource `json:"IspimResources"`
}

type GetIspimResourceResponseRoot struct {
	GetIspimResourceResponse *GetIspimResourceResponse `json:"GetIspimResourceResponse"`
}

type GetIspimResourcesResponse struct {
	GetIspimResourcesResponse []GetIspimResourcesResponse `json:"GetIspimResourceResponses"`
}

var _ IspimResourcesService = &IspimResourcesServiceOp{}

// Create a new Resource with a given configuration.
func (ispim *IspimResourcesServiceOp) Create(ctx context.Context, ismpimr *IspimResourceRequest) (*IspimResource, *Response, error) {
	log.Printf("[DEBUG]: ismpimresources.go: Printing in the ISPIM Resource Create service")
	if ismpimr == nil {
		return nil, nil, errors.NewArgError("ovmr - input request", "cannot be nil")
	}

	path := ispimresourcePath

	req, err := ispim.client.NewRequest(ctx, http.MethodPost, path, ismpimr)
	if err != nil {
		return nil, nil, err
	} else {
		log.Printf("[DEBUG]: IspimResource:  Successfully executed the NewRequest Call for IspimResource")

	}
	root := IspimResource{}

	resp, err := ispim.client.Do(ctx, req, &root)
	//log.Printf("We should have got the response here.. maybe we should try to parse it %s",resp.Body)
	if err != nil {
		//return nil, resp, err
		log.Fatal(err)

		return nil, resp, err
	}

	log.Printf("[DEBUG]: isimProvider: Printing the response %d", resp.StatusCode)
	log.Printf("[DEBUG]: isimProvider: Printing the response from the call %#v\n", root)

	return &root, resp, err
}

// Get

func (ispims *IspimResourcesServiceOp) Get(ctx context.Context, ispimId string) ([]GetIspimResourceResponse, *Response, error) {

	// The ISPIMID is actually the filter that needs to be passed to get the resource. Since we need to get only one result
	// It will have to be unique . So we have to pass in a combination of aliases preferably the
	// fqdn, cdir and ip address

	log.Printf("[DEBUG] Printing the from the GET CALL requestNumber %s", ispimId)

	getPath := fmt.Sprintf("%s%s", ispimresourceGetPath, ispimId)

	log.Printf("[DEBUG] In the get method - the path is %s", getPath)

	req, err := ispims.client.NewRequest(ctx, http.MethodGet, getPath, nil)
	if err != nil {
		return nil, nil, err
	} else {
		log.Printf("[DEBUG]: ISIM-Provider:  Successfully executed the GetRequest Call for IspimResource")

	}

	//root := GetIspimResourceResponse{}
	var root []GetIspimResourceResponse

	log.Printf("[DEBUG] : ispimresource.go: Printing the response from the get call %#v", root)
	resp, err := ispims.client.Do(ctx, req, &root)
	if err != nil {
		log.Printf("[DEBUG]. If we are here , the do call has failed")
		log.Fatal(err)
		return nil, resp, err
	}

	return root, resp, err
}

// Delete
func (ispims *IspimResourcesServiceOp) Delete(ctx context.Context, ispim_resource_id string) (*Response, error) {

	log.Printf("[DEBUG] Printing the from the GET CALL requestNumber %s", ispim_resource_id)

	getPath := fmt.Sprintf("%s%s", ispimresourceGetPath, ispim_resource_id)

	log.Printf("[DEBUG] In the delete method - the path is %s", getPath)

	req, err := ispims.client.NewRequest(ctx, http.MethodDelete, getPath, nil)
	if err != nil {
		return nil, err
	} else {
		log.Printf("[DEBUG]: ISIM-Provider:  Successfully executed the GetRequest Call for IspimResource")

	}

	return ispims.client.Do(ctx, req, nil)

}
