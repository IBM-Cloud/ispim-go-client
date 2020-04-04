package ibmispim

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

const adminDomainsPath = "ispim/rest/organizationcontainers/admindomains?"

type AdminDomainsService interface {
	Get(context.Context, string) ([]AdminDomain, *Response, error)
}

type AdminDomainsServiceOp struct {
	client *Client
}

var _ AdminDomainsService = &AdminDomainsServiceOp{}

type AdminDomain struct {
	Links Links `json:"_links"`
}

type Links struct {
	Self Self `json:"self"`
}

type Self struct {
	Href  string `json:"href"`
	Title string `json:"title"`
}

type AdminDomainRoot struct {
	AdminDomain *AdminDomain `json:"AdminDomain"`
}

type AdminDomainsRoot struct {
	AdminDomains []AdminDomain `json:"AdminDomains"`
}

func (s *AdminDomainsServiceOp) Get(ctx context.Context, servername string) ([]AdminDomain, *Response, error) {
	//path := fmt.Sprintf("%s/%v", serversBasePath, servername)

	//
	//var filter="filter=(OU="
	//var endFilter=")"
	// filter="(%26(ou=A1A*)(description=OS*))"
	//var filterStart="filter=%26(OU="
	//var cdir="AIA*"
	//var filterend=")(description="
	//var appendToFilter="))"
	var filter = "filter=(%26"

	//getPath := fmt.Sprintf("%s%s%s%s", adminDomainsPath,filter,servername,endFilter)
	getPath := fmt.Sprintf("%s%s%s", adminDomainsPath, filter, servername)

	log.Printf("[DEBUG] - The path to get the admindomains is %s", getPath)

	req, err := s.client.NewRequest(ctx, http.MethodGet, getPath, nil)
	if err != nil {
		return nil, nil, err
	} else {
		log.Printf("[DEBUG]: adminDomains: Successfully executed the call")
	}

	//root := AdminDomain{}
	var root []AdminDomain
	log.Printf("[ DEBUG: Printing the root %+v\n", &root)
	resp, err := s.client.Do(ctx, req, &root)
	if err != nil {
		log.Printf("%+v", root)
		log.Printf("[DEBUG]: admindomains.go: If we are here , the do call has failed")
		log.Fatal(err)
		return nil, resp, err
	}

	log.Printf("[++++++DEBUG++++++] Printing the root object %+v\n", &root)
	return root, resp, err
}
