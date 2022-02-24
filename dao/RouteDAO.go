package dao

import (
	"encoding/json"
	"strconv"

	datav1alpha1 "github.com/MadeByMakers/kong-operator-for-k8s/api/v1alpha1"
	httpClient "github.com/MadeByMakers/kong-operator-for-k8s/util"
)

type RouteDAO struct{}

func (this RouteDAO) Delete(route datav1alpha1.Route) datav1alpha1.Route {
	status, response := httpClient.Delete(httpClient.GetBaseURL() + "/routes/" + route.Spec.Id)

	// OK
	if status == 204 {
		route.Status = datav1alpha1.RouteStatus{
			Message: "DELETED",
			Code:    200,
		}
	} else {
		var stringValue string
		json.Unmarshal(response, &stringValue)

		route.Status = datav1alpha1.RouteStatus{
			Message: "ERROR (" + strconv.Itoa(status) + ")",
			Code:    status,
			Response: datav1alpha1.HttpStatus{
				Code: status,
				Body: stringValue,
			},
		}
	}

	return route
}

func (this RouteDAO) Save(route datav1alpha1.Route) datav1alpha1.Route {

	if route.Spec.Service.Id == "" && route.Spec.Service.Name != "" {
		service := ServiceDAO{}.Get(route.Spec.Service.Name)

		if service != nil {
			route.Spec.Service.Id = service.Id
		} else {
			route.Status = datav1alpha1.RouteStatus{
				Message: "Service '" + route.Spec.Service.Name + "' not found",
				Code:    404,
				Response: datav1alpha1.HttpStatus{
					Code: 404,
				},
			}

			return route
		}

	}

	status, response := httpClient.Put(httpClient.GetBaseURL()+"/routes/"+route.Spec.Name, route.Spec)

	// OK
	if status == 200 || status == 201 {
		var returnValue datav1alpha1.RouteSpec
		json.Unmarshal(response, &returnValue)

		returnValue.Service.Name = route.Spec.Service.Name

		route.Spec = returnValue
		route.Status = datav1alpha1.RouteStatus{
			Message: "SAVED",
			Code:    200,
		}
	} else {
		var stringValue string
		json.Unmarshal(response, &stringValue)

		route.Status = datav1alpha1.RouteStatus{
			Message: "ERROR (" + strconv.Itoa(status) + ")",
			Code:    status,
			Response: datav1alpha1.HttpStatus{
				Code: status,
				Body: stringValue,
			},
		}
	}

	return route
}

func (this RouteDAO) Get(nameOrId string) *datav1alpha1.RouteSpec {
	status, response := httpClient.Get(httpClient.GetBaseURL() + "/routes/" + nameOrId)

	// OK
	if status == 200 {
		var returnValue datav1alpha1.RouteSpec
		json.Unmarshal(response, &returnValue)
		return &returnValue
	}

	return nil
}

func (this RouteDAO) GetAll() []datav1alpha1.RouteSpec {
	status, response := httpClient.Get(httpClient.GetBaseURL() + "/routes/")

	// OK
	if status == 200 {
		var returnValue []datav1alpha1.RouteSpec
		json.Unmarshal(response, &returnValue)
		return returnValue
	}

	return nil
}
