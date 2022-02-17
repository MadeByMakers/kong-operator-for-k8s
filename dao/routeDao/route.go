package routeDAO

import (
	"encoding/json"
	"strconv"

	datav1alpha1 "github.com/MadeByMakers/kong-operator-for-k8s/api/v1alpha1"
	httpClient "github.com/MadeByMakers/kong-operator-for-k8s/util"
)

func Create(route datav1alpha1.Route) datav1alpha1.Route {
	status, response := httpClient.Post(httpClient.GetBaseURL()+"/routes", route.Spec)

	// OK
	if status == 201 {
		var returnValue datav1alpha1.RouteSpec
		json.Unmarshal(response, &returnValue)

		route.Spec = returnValue
		route.Status = datav1alpha1.RouteStatus{
			Message: "OK",
		}
	} else {
		var stringValue string
		json.Unmarshal(response, &stringValue)

		route.Status = datav1alpha1.RouteStatus{
			Message: "ERROR (" + strconv.Itoa(status) + ")",
			Response: datav1alpha1.HttpStatus{
				Code: status,
				Body: stringValue,
			},
		}
	}

	return route
}

func Delete(route datav1alpha1.Route) datav1alpha1.Route {
	status, response := httpClient.Delete(httpClient.GetBaseURL() + "/routes/" + route.Spec.Id)

	// OK
	if status == 204 {
		route.Status = datav1alpha1.RouteStatus{
			Message: "DELETED",
		}
	} else {
		var stringValue string
		json.Unmarshal(response, &stringValue)

		route.Status = datav1alpha1.RouteStatus{
			Message: "ERROR (" + strconv.Itoa(status) + ")",
			Response: datav1alpha1.HttpStatus{
				Code: status,
				Body: stringValue,
			},
		}
	}

	return route
}

func Update(route datav1alpha1.Route) datav1alpha1.Route {
	status, response := httpClient.Patch(httpClient.GetBaseURL()+"/routes", route.Spec)

	// OK
	if status == 200 {
		var returnValue datav1alpha1.RouteSpec
		json.Unmarshal(response, &returnValue)

		route.Spec = returnValue
		route.Status = datav1alpha1.RouteStatus{
			Message: "OK",
		}
	} else {
		var stringValue string
		json.Unmarshal(response, &stringValue)

		route.Status = datav1alpha1.RouteStatus{
			Message: "ERROR (" + strconv.Itoa(status) + ")",
			Response: datav1alpha1.HttpStatus{
				Code: status,
				Body: stringValue,
			},
		}
	}

	return route
}

func Get(nameOrId string) *datav1alpha1.RouteSpec {
	status, response := httpClient.Get(httpClient.GetBaseURL() + "/routes/" + nameOrId)

	// OK
	if status == 200 {
		var returnValue datav1alpha1.RouteSpec
		json.Unmarshal(response, &returnValue)
		return &returnValue
	} else {
		var stringValue string
		json.Unmarshal(response, &stringValue)
	}

	return nil
}

func GetAll() []datav1alpha1.RouteSpec {
	status, response := httpClient.Get(httpClient.GetBaseURL() + "/routes/")

	// OK
	if status == 200 {
		var returnValue []datav1alpha1.RouteSpec
		json.Unmarshal(response, &returnValue)
		return returnValue
	} else {
		var stringValue string
		json.Unmarshal(response, &stringValue)
	}

	return nil
}
