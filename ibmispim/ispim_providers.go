package ibmispim

import (
	"context"
	"fmt"
	"github.com/IBM-Cloud/ispim-go-client/errors"
	"log"
	"net/http"
)

const IspimIdpProviderPath = "/ispim/rest/services"
const IspimIdpProviderGetPath = "ispim/rest/services/idp?"

type IspimIdpProvidersService interface {
	Get(context.Context, string) ([]GetIspimProviderResponse, *Response, error)
	Create(context.Context, *IspimIdpProviderRequest) (*IspimIdpProvider, *Response, error)
}

type IspimIdpProvidersServiceOp struct {
	client *Client
}

/// START of REQUEST Structure

type IspimIdpProviderRequest struct {
	BatchAction string       `json:"batchAction"`
	EntityList  []EntityList `json:"entityList"`
}

type EntityList struct {
	IdpLinks   IdpLinks   `json:"_links"`
	Attributes Attributes `json:"_attributes"`
}

type IdpLinks struct {
	Container Container `json:"container"`
	//Owner Owner `json:"owner,omitempty"`
}

type Container struct {
	Href string `json:"href"`
}

type Owner struct {
	Href string `json:"href"`
}

type Attributes struct {
	Profile               []string `json:"profile"`
	Erservicename         []string `json:"erservicename,omitempty"`
	Description           []string `json:"description,omitempty"`
	Erauthenticatemode    []string `json:"erauthenticatemode,omitempty"`
	Erurl                 []string `json:"erurl,omitempty"`
	Erposixauthmethod     []string `json:"erposixauthmethod,omitempty"`
	Erserviceuid          []string `json:"erserviceuid,omitempty"`
	Erpassword            []string `json:"erpassword,omitempty"`
	Erposixuseshadow      []string `json:"erposixuseshadow,omitempty"`
	Erposixhomedirremove  []string `json:"erposixhomedirremove,omitempty"`
	Erposixfailedlogincmd []string `json:"erposixfailedlogincmd,omitempty"`
	Erposixsudoerspath    []string `json:"erposixsudoerspath,omitempty"`
	Eritdiurl             []string `json:"eritdiurl,omitempty"`
	Erposixpassphrase     []string `json:"erposixpassphrase,omitempty"`
	Erposixpkfile         []string `json:"erposixpkfile,omitempty"`
	Erposixusesudo        []string `json:"erposixusesudo,omitempty"`
	Eruid                 []string `json:"eruid,omitempty"`
}

/// End of the Request

type IspimIdpProvider struct {
	ResponseList  []ResponseList `json:"responseList"`
	OverAllStatus string         `json:"overAllStatus"`
}

type Erparent struct {
	Href string `json:"href"`
}
type IdpProviderSelf struct {
	Href  string `json:"href"`
	Title string `json:"title"`
}
type IdpProviderLinks struct {
	Owner    []Owner    `json:"owner"`
	Erparent []Erparent `json:"erparent"`
	Self     Self       `json:"self"`
}

type Entity struct {
	Links      Links      `json:"_links"`
	Attributes Attributes `json:"_attributes"`
}
type Metadata struct {
}
type IspimIdpResponse struct {
	Entity   Entity   `json:"entity"`
	Metadata Metadata `json:"metadata"`
	Status   int      `json:"status"`
}
type SelfLink struct {
	Href  string `json:"href"`
	Title string `json:"title"`
}

/*type IdpProviderResponseList struct {
	RequestAction string   `json:"requestAction"`
	IspimIdpResponse      IspimIdpResponse `json:"response"`
	SelfLink      SelfLink `json:"selfLink"`
}*/

type ResponseList struct {
	RequestAction    string           `json:"requestAction"`
	IspimIdpResponse IspimIdpResponse `json:"response"`
	SelfLink         SelfLink         `json:"selfLink"`
}

type GetIspimProviderResponse struct {
	Links Links `json:"_links"`
	//GetIspimResourceAttributes GetIspimResourceAttributes `json:"_attributes"`
}

type GetIspimProviderResponseRoot struct {
	GetIspimResourceResponse *GetIspimResourceResponse `json:"GetIspimProviderResponse"`
}

type GetIspimProvidersResponse struct {
	GetIspimProvidersResponse []GetIspimProviderResponse `json:"GetIspimProviderResponses"`
}

var _ IspimIdpProvidersService = &IspimIdpProvidersServiceOp{}

// Create a new Provider with a given configuration.
func (ispim *IspimIdpProvidersServiceOp) Create(ctx context.Context, ismpimr *IspimIdpProviderRequest) (*IspimIdpProvider, *Response, error) {
	log.Printf("[DEBUG]: IspimIdpProvider: Printing in the ISIM Idp Provider Create service")
	if ismpimr == nil {
		return nil, nil, errors.NewArgError("ispim - input request", "cannot be nil")
	}

	path := IspimIdpProviderPath

	req, err := ispim.client.NewRequest(ctx, http.MethodPost, path, ismpimr)
	if err != nil {
		return nil, nil, err
	} else {
		log.Printf("[DEBUG]: IspimIdpProvider:  Successfully executed the NewRequest Call for IspimIdpProvider")

	}
	root := IspimIdpProvider{}

	resp, err := ispim.client.Do(ctx, req, &root)
	//log.Printf("We should have got the response here.. maybe we should try to parse it %s",resp.Body)
	if err != nil {
		//return nil, resp, err
		log.Fatal(err)

		return nil, resp, err
	}

	log.Printf("[DEBUG]: isimProvider: Printing the response %d", resp.StatusCode)
	//log.Printf("[DEBUG]: isimProvider: Printing the response from the call %#v\n", root)

	return &root, resp, err
}

// Get the Configuration for the Provider created

func (ispims *IspimIdpProvidersServiceOp) Get(ctx context.Context, providerFilter string) ([]GetIspimProviderResponse, *Response, error) {
	//requestNumber := strconv.Itoa(ispimId)

	//var filter="filter=(%26"
	var filter = "filter="
	var attributes = "&attributes=erservicename,erurl,erauthenticatemode,erserviceuid"

	getPath := fmt.Sprintf("%s%s%s%s", IspimIdpProviderGetPath, filter, providerFilter, attributes)
	log.Printf("[DEBUG] Printing the from the GET CALL requestNumber %s", providerFilter)

	log.Printf("[DEBUG] In the get method - the path is %s", getPath)

	req, err := ispims.client.NewRequest(ctx, http.MethodGet, getPath, nil)

	if err != nil {
		return nil, nil, err
	} else {
		log.Printf("[DEBUG]: ISIM-Provider:  Successfully executed the GetRequest Call for ISIMProvider")

	}

	var root []GetIspimProviderResponse
	resp, err := ispims.client.Do(ctx, req, &root)
	if err != nil {
		log.Printf("[DEBUG]. If we are here , the do call has failed")
		log.Fatal(err)
		return nil, resp, err
	}

	return root, resp, err
}

// Exists
