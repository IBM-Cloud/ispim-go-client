package ibmispim

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/google/go-querystring/query"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"reflect"
)

const (
	mediaType   = "application/json"
	providerURL = "providerURL"
	username    = "username"
)

// Client manages communication with ISIM  API.
type Client struct {
	// HTTP client used to communicate with the ISIM API.
	//Credentials *Credentials
	URL string
	//Username string
	//Password string
	Token                string
	Endpoint             string
	Insecure             bool
	client               *http.Client
	BaseURL              *url.URL
	IspimIdpProviders    IspimIdpProvidersService
	AdminDomains         AdminDomainsService
	IspimResources       IspimResourcesService
	IspimCredentials     IspimCredentialsService
	IspimSyncCredentials IspimSyncCredentialsService
}

const (
	//ErrCodeEmptyResponse ...
	ErrCodeEmptyResponse   = "EmptyResponseBody"
	ErrCodeFailedtoConnect = "FailedtoConnect"
	ErrCodeVolumeNotFound  = "VolumeNotFoundError"
)

type Response struct {
	*http.Response
}

// An ErrorResponse reports the error caused by an API request
type ErrorResponse struct {
	// HTTP response that caused this error
	Response *http.Response

	// Error message
	Message string `json:"STATUS"`

	// RequestID returned from the API, useful to contact support.
	RequestID string `json:"REQUEST_ID"`
}

func addOptions(s string, opt interface{}) (string, error) {
	v := reflect.ValueOf(opt)

	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	origURL, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	origValues := origURL.Query()

	newValues, err := query.Values(opt)
	if err != nil {
		return s, err
	}

	for k, v := range newValues {
		origValues[k] = v
	}

	origURL.RawQuery = origValues.Encode()
	return origURL.String(), nil
}

//func NewClient(username, password , token, endpoint string, insecure bool) (*Client, error) {
func NewClient(token, endpoint string, insecure bool) (*Client, error) {
	/* Make sure that you add the client service to the Op - Refer to line 132 */

	baseUrl := os.Getenv("ISPIM_URL")

	baseURL, _ := url.Parse(baseUrl)
	transCfg := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure}, // ignore expired SSL certificates
	}

	httpClient := http.DefaultClient
	httpClient.Transport = transCfg
	//c := &Client{Username: username, Password: password, Token: token, BaseURL: baseURL, client: http.DefaultClient}
	/* ISPIM uses a token based authentication . We obtain this during our first login or get it from the ISIM team */
	c := &Client{Token: token, BaseURL: baseURL, client: http.DefaultClient}
	c.IspimIdpProviders = &IspimIdpProvidersServiceOp{client: c}
	c.AdminDomains = &AdminDomainsServiceOp{client: c}
	c.IspimResources = &IspimResourcesServiceOp{client: c}
	c.IspimCredentials = &IspimCredentialsServiceOp{client: c}
	c.IspimSyncCredentials = &IspimSyncCredentialsServiceOp{client: c}

	return c, nil

}

// ClientOpt are options for New.
type ClientOpt func(*Client) error

// SetBaseURL is a client option for setting the base URL.
func SetBaseURL(bu string) ClientOpt {
	return func(c *Client) error {
		u, err := url.Parse(bu)
		if err != nil {
			return err
		}

		c.BaseURL = u
		return nil
	}
}

// NewRequest creates an API request. A relative URL can be provided in urlStr, which will be resolved to the
// BaseURL of the Client. Relative URLS should always be specified without a preceding slash. If specified, the
// value pointed to by body is JSON encoded and included in as the request body.
func (c *Client) NewRequest(ctx context.Context, method, urlStr string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)

	log.Printf("[DEBUG: ibmispim : ] Printing the body of this request %v\n", ctx)
	//log.Printf("[DEBUG: ibmispim : ] Printing the username being passed to the NewRequest %s",c.Username)
	buf := new(bytes.Buffer)
	if body != nil {
		err = json.NewEncoder(buf).Encode(body)

		log.Printf(buf.String())
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	} else {
		log.Printf("[+++++DEBUG+++++] Starting to  process the http request  %s", u.String())
	}

	req.Header.Add("Content-Type", mediaType)
	req.Header.Add("Accept", mediaType)
	req.Header.Add("X-HTTP-Method-Override", "SUBMIT-IN-BATCH")
	req.Header.Add("Authorization", "Basic "+c.Token)
	log.Printf("[DEBUG : ibmispim: ] Printing the request to be passed to the call %v", req)
	return req, nil
}

func basicAuth(username, password string) string {
	log.Printf("[DEBUG: basicAuth] : Now building the auth function")
	auth := username + ":" + password
	log.Printf("[DEBUG: basicAuth] : Generated the token %s", auth)
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

// API request that is called from the client
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*Response, error) {

	log.Printf("[DEBUG]: ibmispim:  In the Do Method - Printing request %#v\n", req)

	resp, err := DoRequestWithClient(ctx, c.client, req)
	if err != nil {
		return nil, err
	}
	log.Printf("Printing the response post calls %#v", resp)
	log.Printf("Prrinting RESPONSE DATA %s", resp.Header.Get("Set-Cookie"))

	if err != nil {
		return nil, err
	}

	defer func() {
		if rerr := resp.Body.Close(); err == nil {
			err = rerr
		}

	}()
	//
	response := newResponse(resp)
	log.Printf("Printing the response prior to the body %#v", response)
	err = CheckResponse(resp)
	if err != nil {
		log.Printf("DEBUG Printing the value %s", response.Header.Get("JSESSIONID"))
		return response, err
	}

	bodyBytes, _ := ioutil.ReadAll(response.Body)
	bodyString := string(bodyBytes)
	//print raw response body for debugging purposes
	log.Printf("[DEBUG]: ibmispim: In the Client do call - Printing response %s", bodyString)
	//
	//
	// //log.Printf("[DEBUG] In the do call after checking - Printing response %s", bodyString)
	//
	if v != nil {
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, resp.Body)
			if err != nil {
				return nil, err
			}
		} else {
			log.Printf("[DEBUG] No error during copy of the resp body %s", string(bodyBytes))
			err = json.Unmarshal(bodyBytes, &v)
			if err != nil {
				log.Printf("[DEBUG]: ibmispim.go:  Failed to decode the response")
				return nil, err
			}
		}
	}

	log.Printf("[DEBUG] PRINTING THE RESPONSE TO BE PASSED BACK %#v\n", response)
	return response, err
}

// newResponse creates a new Response for the provided http.Response
func newResponse(r *http.Response) *Response {
	response := Response{Response: r}
	//response.populateRate()

	return &response
}

func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; c >= 200 && c <= 299 {
		return nil
	}
	//
	errorResponse := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && len(data) > 0 {
		err := json.Unmarshal(data, errorResponse)
		if err != nil {
			errorResponse.Message = string(data)

		}
	}
	// 	log.Printf("[DEBUG PRINTING ERROR RESPONSE] %s",err)
	return errorResponse
}

// DoRequest submits an HTTP request.
func DoRequest(ctx context.Context, req *http.Request) (*http.Response, error) {
	log.Printf("[DEBUG] In the DoRequest")
	return DoRequestWithClient(ctx, http.DefaultClient, req)
}

// DoRequestWithClient submits an HTTP request using the specified client.
func DoRequestWithClient(
	ctx context.Context,
	client *http.Client,
	req *http.Request) (*http.Response, error) {

	log.Printf("[DEBUG] In the DoRequestWithClient Method")
	log.Printf("[DEBUG] Printing the request %v", req)

	req = req.WithContext(ctx)
	return client.Do(req)
}

func (r *ErrorResponse) Error() string {
	if r.RequestID != "" {
		return fmt.Sprintf("%v %v: %d (request %q) %v",
			r.Response.Request.Method, r.Response.Request.URL, r.Response.StatusCode, r.RequestID, r.Message)
	}
	return fmt.Sprintf("%v %v: %d %v",
		r.Response.Request.Method, r.Response.Request.URL, r.Response.StatusCode, r.Message)
}

// String is a helper routine that allocates a new string value
// to store v and returns a pointer to it.
func String(v string) *string {
	p := new(string)
	*p = v
	return p
}

// Int is a helper routine that allocates a new int32 value
// to store v and returns a pointer to it, but unlike Int32
// its argument value is an int.
func Int(v int) *int {
	p := new(int)
	*p = v
	return p
}

// Bool is a helper routine that allocates a new bool value
// to store v and returns a pointer to it.
func Bool(v bool) *bool {
	p := new(bool)
	*p = v
	return p
}

// StreamToString converts a reader to a string
func StreamToString(stream io.Reader) string {
	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(stream)
	return buf.String()
}
