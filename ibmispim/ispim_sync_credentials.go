package ibmispim

import (
	"context"
	"fmt"
	"github.com/IBM-Cloud/ispim-go-client/errors"
	"log"
	"net/http"
	"strconv"
)

const IspimSyncCredentialPath = "/ispim/rest/credentials"

//const IspimSyncCredentialPostPath="/ispim/rest/services/"

type IspimSyncCredentialsService interface {
	Get(context.Context, int) (*IspimSyncCredential, *Response, error)
	Create(context.Context, *IspimSyncCredentialRequest, string) (*IspimSyncResponse, *Response, error)
	//Update(context.Context, int, *IspimIdpCredentialUpdateRequest) (*IspimIdpCredential, *Response, error)
	//Delete(context.Context, *IspimIdpCredentialDeleteRequest) (*IspimIdpCredential, *Response, error)

}

type IspimSyncCredentialsServiceOp struct {
	client *Client
}

/// START of REQUEST Structure

type IspimSyncCredentialRequest struct {
	//IdpUrl string `json:"idpurl"`
	Justification          string                 `json:"justification"`
	ResetPassword          string                 `json:"resetPassword"`
	ResetPasswordOnCheckin string                 `json:"resetPasswordOnCheckin"`
	IspimSyncCredentials   []IspimSyncCredentials `json:"credentials"`
}

type SyncCredentialRequest struct {
	CredentialLink             string                     `json:"credentialLink"`
	IspimSyncCredentialRequest IspimSyncCredentialRequest `json:"ispimcredentialrequest"`
}

type IspimSyncCredentials struct {
	Href string `json:"href"`
}

/// End of the Request

///Start of the response
type IspimSyncResponse []struct {
	IspimSyncResponseLinks IspimSyncResponseLinks `json:"_links"`
	RequestID              string                 `json:"requestId"`
}
type IspimSyncResponseRequest struct {
	Href string `json:"href"`
}
type IspimSyncResponseLinks struct {
	Request Request `json:"request"`
}

/// End of the response

type IspimSyncCredential struct {
	IspimSyncCredentialResponseList []IspimSyncCredentialResponseList `json:"responseList"`
	OverAllStatus                   string                            `json:"overAllStatus"`
}

type IspimSyncCredentialEntity struct {
	Links     Links  `json:"_links"`
	RequestID string `json:"requestId"`
}

type IspimSyncCredentialResponse struct {
	IspimSyncCredentialEntity IspimSyncCredentialEntity `json:"entity"`
	Metadata                  Metadata                  `json:"metadata"`
	Status                    int                       `json:"status"`
}
type IspimSyncCredentialResponseList struct {
	RequestAction               string                      `json:"requestAction"`
	IspimSyncCredentialResponse IspimSyncCredentialResponse `json:"response"`
}

///End of the response

var _ IspimSyncCredentialsService = &IspimSyncCredentialsServiceOp{}

// Create a new Provider with a given configuration.
func (ispim *IspimSyncCredentialsServiceOp) Create(ctx context.Context, ismpimr *IspimSyncCredentialRequest, idpurl string) (*IspimSyncResponse, *Response, error) {
	log.Printf("[DEBUG]: IspimIdpCredential: Printing in the Sync Credentials Create service")
	if ismpimr == nil {
		return nil, nil, errors.NewArgError("ispimr - input request", "cannot be nil")
	}

	//log.Printf("[DEBUG] Printing the url id %s",ismpimr.CredentialLink)

	//path := IspimSyncCredentialPostPath

	log.Printf("[DEBUG ] - The idpurl for the identityprovider is %s", idpurl)

	credPath := "/credentials"

	getPath := fmt.Sprintf("%s%s", idpurl, credPath)

	req, err := ispim.client.NewRequest(ctx, http.MethodPost, getPath, ismpimr)
	if err != nil {
		return nil, nil, err
	} else {
		log.Printf("[DEBUG]: IspimIdpCredential:  Successfully executed the NewRequest Call for IspimIdpCredential")

	}

	root := IspimSyncResponse{}

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

// Get the Configuration for the SyncCredential created

func (ispims *IspimSyncCredentialsServiceOp) Get(ctx context.Context, ispimId int) (*IspimSyncCredential, *Response, error) {
	requestNumber := strconv.Itoa(ispimId)

	log.Printf("[DEBUG] Printing the from the GET CALL requestNumber %s", requestNumber)

	getPath := fmt.Sprintf("%s%s", IspimIdpCredentialPath, requestNumber)

	log.Printf("[DEBUG] In the get method - the path is %s", getPath)

	req, err := ispims.client.NewRequest(ctx, http.MethodGet, getPath, nil)
	if err != nil {
		return nil, nil, err
	} else {
		log.Printf("[DEBUG]: ISIM-Provider:  Successfully executed the GetRequest Call for ISIMProvider")

	}

	root := IspimSyncCredential{}
	resp, err := ispims.client.Do(ctx, req, &root)
	if err != nil {
		log.Printf("[DEBUG]. If we are here , the do call has failed")
		log.Fatal(err)
		return nil, resp, err
	}

	return &root, resp, err
}

// Delete - Add code to disconnect credentials - We do a lookup and then perform the disconnect -
