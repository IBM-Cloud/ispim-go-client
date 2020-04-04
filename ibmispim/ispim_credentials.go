package ibmispim

import (
	"context"
	"fmt"
	"github.com/IBM-Cloud/ispim-go-client/errors"
	"log"
	"net/http"
	//"strconv"
)

const IspimIdpCredentialPath = "/ispim/rest/credentials"
const IspimIdpCredentialGetPath = "/ispim/rest/credentials/?attributes=*&embedded=resource.name&filter="

type IspimCredentialsService interface {
	Get(context.Context, string) ([]GetIspimCredentialResponse, *Response, error)
	Create(context.Context, *IspimCredentialRequest) (*IspimCredential, *Response, error)
	//Update(context.Context, int, *IspimIdpCredentialUpdateRequest) (*IspimIdpCredential, *Response, error)
	Delete(context.Context, string) (*Response, error)
}

type IspimCredentialsServiceOp struct {
	client *Client
}

/// START of REQUEST Structure

type IspimCredentialRequest struct {
	IspimCredentialCreateList []IspimCredentialCreateList `json:"createList"`
}
type IspimCredentialContainer struct {
	Href string `json:"href"`
}
type IspimCredentialResource struct {
	Href string `json:"href"`
}

type IspimCredentialLinks struct {
	IspimCredentialContainer IspimCredentialContainer `json:"container"`
	IspimCredentialResource  IspimCredentialResource  `json:"resource"`
}
type IspimCredentialAttributes struct {
	ErCredentialName       []string `json:"erCredentialName"`
	Password               []string `json:"password"`
	Description            []string `json:"description"`
	ErTag                  []string `json:"erTag"`
	ErpPwdRotationInterval []string `json:"erpPwdRotationInterval"`
	UseDefaultSetting      []string `json:"useDefaultSetting"`
}
type IspimCredentialSettings struct {
	AccessMode          string `json:"accessMode"`
	IsShortTermPassword string `json:"isShortTermPassword"`
	CheckoutDuration    string `json:"checkoutDuration"`
	IsCheckoutSearch    string `json:"isCheckoutSearch"`
	IsPasswordViewable  string `json:"isPasswordViewable"`
}
type IspimCredentialObjects struct {
	IspimCredentialSettings []IspimCredentialSettings `json:"settings"`
}
type IspimCredentialCreateList struct {
	IspimCredentialLinks      IspimCredentialLinks      `json:"_links"`
	IspimCredentialAttributes IspimCredentialAttributes `json:"_attributes"`
	IspimCredentialObjects    IspimCredentialObjects    `json:"_objects,omitempty"`
}

/// End of the Request

///Start of the response

type IspimCredential struct {
	IspimCredentialResponseList []IspimCredentialResponseList `json:"responseList"`
	OverAllStatus               string                        `json:"overAllStatus"`
}
type Request struct {
	Href string `json:"href"`
}

type IspimCredentialEntity struct {
	Links     Links  `json:"_links"`
	RequestID string `json:"requestId"`
}

type IspimCredentialResponse struct {
	IspimCredentialEntity IspimCredentialEntity `json:"entity"`
	Metadata              Metadata              `json:"metadata"`
	Status                int                   `json:"status"`
}
type IspimCredentialResponseList struct {
	RequestAction           string                  `json:"requestAction"`
	IspimCredentialResponse IspimCredentialResponse `json:"response"`
}

/// Start of the response from the get call

type GetIspimCredentialResponse struct {
	Links      Links      `json:"_links"`
	Embedded   Embedded   `json:"_embedded"`
	Attributes Attributes `json:"_attributes"`
	Objects    Objects    `json:"_objects"`
}
type IDProvider struct {
	Href  string `json:"href"`
	Title string `json:"title"`
}

type Resource struct {
	Links      Links      `json:"_links"`
	Attributes Attributes `json:"_attributes"`
}
type Embedded struct {
	Resource Resource `json:"resource"`
}
type Settings struct {
	IsCheckoutSearch    string `json:"isCheckoutSearch"`
	IsShortTermPassword string `json:"isShortTermPassword"`
	CheckoutDuration    string `json:"checkoutDuration"`
	AccessMode          string `json:"accessMode"`
	IsExclusive         string `json:"isExclusive"`
	IsPasswordViewable  string `json:"isPasswordViewable"`
}
type Objects struct {
	Settings []Settings `json:"settings"`
}

///

var _ IspimCredentialsService = &IspimCredentialsServiceOp{}

// Create a new Provider with a given configuration.
func (ispim *IspimCredentialsServiceOp) Create(ctx context.Context, ismpimr *IspimCredentialRequest) (*IspimCredential, *Response, error) {
	log.Printf("[DEBUG]: IspimIdpCredential: Printing in ISPIM Credential Create service")
	if ismpimr == nil {
		return nil, nil, errors.NewArgError("ovmr - input request", "cannot be nil")
	}

	path := IspimIdpCredentialPath

	req, err := ispim.client.NewRequest(ctx, http.MethodPost, path, ismpimr)
	if err != nil {
		return nil, nil, err
	} else {
		log.Printf("[DEBUG]: IspimIdpCredential:  Successfully executed the NewRequest Call for IspimIdpCredential")

	}

	root := IspimCredential{}

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

// Get the Configuration for the Provider created

func (ispims *IspimCredentialsServiceOp) Get(ctx context.Context, ispimId string) ([]GetIspimCredentialResponse, *Response, error) {
	//requestNumber := strconv.Itoa(ispimId)

	log.Printf("[DEBUG] Printing the filter from  the GET CALL Filter %s", ispimId)

	getPath := fmt.Sprintf("%s%s", IspimIdpCredentialGetPath, ispimId)

	log.Printf("[DEBUG] In the get method - the path is %s", getPath)

	req, err := ispims.client.NewRequest(ctx, http.MethodGet, getPath, nil)
	if err != nil {
		return nil, nil, err
	} else {
		log.Printf("[DEBUG]: ISIM-Provider:  Successfully executed the GetRequest Call for ISIMProvider")

	}

	var root []GetIspimCredentialResponse
	resp, err := ispims.client.Do(ctx, req, &root)
	if err != nil {
		log.Printf("[DEBUG]. If we are here , the do call has failed")
		log.Fatal(err)
		return nil, resp, err
	}

	return root, resp, err
}

// Delete
func (ispims *IspimCredentialsServiceOp) Delete(ctx context.Context, ispim_credential_id string) (*Response, error) {

	log.Printf("[DEBUG] Printing the from the GET CALL requestNumber %s", ispim_credential_id)

	getPath := fmt.Sprintf("%s%s", IspimIdpCredentialPath, ispim_credential_id)

	log.Printf("[DEBUG] In the delete method - the path is %s", getPath)

	req, err := ispims.client.NewRequest(ctx, http.MethodDelete, getPath, nil)
	if err != nil {
		return nil, err
	} else {
		log.Printf("[DEBUG]: ISIM-Provider:  Successfully executed the Delete  Call for IspimCredential")

	}

	return ispims.client.Do(ctx, req, nil)

}
