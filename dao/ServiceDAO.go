package dao

import (
	"encoding/json"
	"strconv"

	datav1alpha1 "github.com/MadeByMakers/kong-operator-for-k8s/api/v1alpha1"
	httpClient "github.com/MadeByMakers/kong-operator-for-k8s/util"
)

type ServiceDAO struct{}

func (this ServiceDAO) Delete(service datav1alpha1.Service) datav1alpha1.Service {
	status, response := httpClient.Delete(httpClient.GetBaseURL() + "/services/" + service.Spec.Id)

	// OK
	if status == 204 {
		service.Status = datav1alpha1.ServiceStatus{
			Message: "DELETED",
			Code:    200,
		}
	} else {
		var stringValue string
		json.Unmarshal(response, &stringValue)

		service.Status = datav1alpha1.ServiceStatus{
			Message: "ERROR (" + strconv.Itoa(status) + ")",
			Code:    status,
			Response: datav1alpha1.HttpStatus{
				Code: status,
				Body: stringValue,
			},
		}
	}

	return service
}

func (this ServiceDAO) Save(service datav1alpha1.Service) datav1alpha1.Service {
	status, response := httpClient.Put(httpClient.GetBaseURL()+"/services/"+service.Spec.Name, service.Spec)

	// OK
	if status == 200 || status == 201 {
		var returnValue datav1alpha1.ServiceSpec
		json.Unmarshal(response, &returnValue)

		service.Spec = returnValue
		service.Status = datav1alpha1.ServiceStatus{
			Message: "SAVED",
			Code:    200,
		}
	} else {
		var stringValue string
		json.Unmarshal(response, &stringValue)

		service.Status = datav1alpha1.ServiceStatus{
			Message: "ERROR (" + strconv.Itoa(status) + ")",
			Code:    status,
			Response: datav1alpha1.HttpStatus{
				Code: status,
				Body: stringValue,
			},
		}
	}

	return service
}

func (this ServiceDAO) Get(nameOrId string) *datav1alpha1.ServiceSpec {
	status, response := httpClient.Get(httpClient.GetBaseURL() + "/services/" + nameOrId)

	// OK
	if status == 200 {
		var returnValue datav1alpha1.ServiceSpec
		json.Unmarshal(response, &returnValue)
		return &returnValue
	}

	return nil
}

func (this ServiceDAO) GetAll() []datav1alpha1.ServiceSpec {
	status, response := httpClient.Get(httpClient.GetBaseURL() + "/services/")

	// OK
	if status == 200 {
		var returnValue []datav1alpha1.ServiceSpec
		json.Unmarshal(response, &returnValue)
		return returnValue
	} else {
		var stringValue string
		json.Unmarshal(response, &stringValue)
	}

	return nil
}
