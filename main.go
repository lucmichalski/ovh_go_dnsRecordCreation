package main 

import (
	"os"
	"strconv"
	"github.com/ovh/go-ovh/ovh"
	"fmt"
)

//ovhzoneRecord define the struct that contain all the argument passed to the post api call to create a DNS record
type ovhZoneRecord struct {
	Id			int		`json:"id,omitempty"`
	FieldType	string	`json:"fieldType"`
	Subdomain	string	`json:"subDomain"`
	Target		string	`json:"target,omitempty"`
	TTL 		int 	`json:"ttl,omitempty"`
}


func main() {
	//Creation of the ovh client
	client, err := ovh.NewEndpointClient("ovh-eu")
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}

	//Fetching all the needed info
	recordType := "A" //recordtype is fixed because for now, I mainly work with IPV4, but you can make the variable set by environment variable by setting the value to os.getenv("OVH_RECORDTYPE")
	domain := os.Getenv("OVH_DOMAIN")
	subDomain := os.Getenv("OVH_SUBDOMAIN")
	endpoint := os.Getenv("OVH_IP_ENDPOINT")
	actionRecord := os.Getenv("OVH_ACTION")

	fmt.Printf("Going to %s the %s subdomain for %s domain pointing to %s \n",actionRecord,subDomain,domain,endpoint)

	//Switch case for handling the action
	switch actionRecord := os.Getenv("OVH_ACTION"); actionRecord {
	case "CREATE":
		record,err := createARecord(client,domain,recordType,subDomain,endpoint)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
		}
		fmt.Printf("Success: %v\n", record)


	case "DELETE":
		err = deleteARecord(client,domain,recordType,endpoint)
		if err != nil {
			fmt.Printf("Error: %s",err)
		}
	
	}

}

//createARecord take as argument the client , zonename,diledtype,subdomain and target, and post a GET request on the url, and write the response to the record struc
func createARecord(ovhClient *ovh.Client, zoneName, fieldType, subdomain, target string) (*ovhZoneRecord, error) {
	url := "/domain/zone/"+ zoneName +"/record"

	parameters := ovhZoneRecord{
		FieldType: fieldType,
		Subdomain: subdomain,
		Target:		target,
	}

	record := ovhZoneRecord{}

	err := ovhClient.Post(url, &parameters, &record)
	if err != nil {
		return nil, fmt.Errorf("OVH API Call Failed: POST %s - %v \n with param %v", url, err, parameters)
	}

	return &record,nil
}

//getRecordID take as arg client, zonename, filedtype and subdomain, and query the id for the subdomain, the id is needed for the deletion of the record 
func getRecordId(ovhClient *ovh.Client, zoneName, fieldType, subdomain string) ([]int,error) {
	url := "/domain/zone/"+ zoneName + "/record?fieldType=" + fieldType + "&subDomain=" + subdomain
	ids := []int{}

	err := ovhClient.Get(url, &ids)
	if err != nil {
		return nil , fmt.Errorf("OVH API Call Failed: GET %s \n Error: %v", url, err)
	}

	return ids, err
}


//deleteARecord takes as argument the client, zoneName, fieldType and subdomain, and susing these arg it will first querry the id using the getRecordId func and make a post request with the id
func deleteARecord(ovhClient *ovh.Client, zoneName, fieldType, subdomain string) error {
	
	ids , _ := getRecordId(ovhClient, zoneName, fieldType, subdomain)
	for _, id := range ids {
		url := "/domain/zone/" + zoneName + "/record/" + strconv.Itoa(id)

		err := ovhClient.Delete(url, nil)
		if err != nil {
			return fmt.Errorf("OVH API Call Failed: DELETE %s \nError: %v",url,err)
		}
	}
	return nil
}
