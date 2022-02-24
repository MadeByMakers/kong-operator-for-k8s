package dao

import (
	"encoding/json"
	"strconv"

	datav1alpha1 "github.com/MadeByMakers/kong-operator-for-k8s/api/v1alpha1"
	httpClient "github.com/MadeByMakers/kong-operator-for-k8s/util"
)

type PluginDAO struct{}

func (this PluginDAO) Delete(plugin datav1alpha1.Plugin) datav1alpha1.Plugin {
	status, response := httpClient.Delete(httpClient.GetBaseURL() + "/plugins/" + plugin.Spec.Id)

	// OK
	if status == 204 {
		plugin.Status = datav1alpha1.PluginStatus{
			Code:    200,
			Message: "DELETED",
		}
	} else {
		var stringValue string
		json.Unmarshal(response, &stringValue)

		plugin.Status = datav1alpha1.PluginStatus{
			Code:    status,
			Message: "ERROR (" + strconv.Itoa(status) + ")",
			Response: datav1alpha1.HttpStatus{
				Code: status,
				Body: stringValue,
			},
		}
	}

	return plugin
}

func (this PluginDAO) Save(plugin datav1alpha1.Plugin) datav1alpha1.Plugin {
	status, response := httpClient.Put(httpClient.GetBaseURL()+"/plugins/"+plugin.Spec.Name, plugin.Spec)

	if plugin.Spec.Service.Name != "" {
		service := ServiceDAO{}.Get(plugin.Spec.Service.Name)

		if service != nil {
			plugin.Spec.Service.Id = service.Id
		} else {
			plugin.Status = datav1alpha1.PluginStatus{
				Message: "Service '" + plugin.Spec.Service.Name + "' not found",
				Code:    404,
				Response: datav1alpha1.HttpStatus{
					Code: 404,
				},
			}

			return plugin
		}
	}

	if plugin.Spec.Route.Name != "" {
		route := RouteDAO{}.Get(plugin.Spec.Route.Name)

		if route != nil {
			plugin.Spec.Route.Id = route.Id
		} else {
			plugin.Status = datav1alpha1.PluginStatus{
				Message: "Route '" + plugin.Spec.Route.Name + "' not found",
				Code:    404,
				Response: datav1alpha1.HttpStatus{
					Code: 404,
				},
			}

			return plugin
		}
	}

	// OK
	if status == 200 || status == 201 {
		var returnValue datav1alpha1.PluginSpec
		json.Unmarshal(response, &returnValue)

		if plugin.Spec.Service.Name != "" {
			returnValue.Service.Name = plugin.Spec.Service.Name
		}

		if plugin.Spec.Route.Name != "" {
			returnValue.Route.Name = plugin.Spec.Route.Name
		}

		plugin.Spec = returnValue
		plugin.Status = datav1alpha1.PluginStatus{
			Code:    200,
			Message: "SAVED",
		}
	} else {
		var stringValue string
		json.Unmarshal(response, &stringValue)

		plugin.Status = datav1alpha1.PluginStatus{
			Code:    status,
			Message: "ERROR (" + strconv.Itoa(status) + ")",
			Response: datav1alpha1.HttpStatus{
				Code: status,
				Body: stringValue,
			},
		}
	}

	return plugin
}

func (this PluginDAO) Get(nameOrId string) *datav1alpha1.PluginSpec {
	status, response := httpClient.Get(httpClient.GetBaseURL() + "/plugins/" + nameOrId)

	// OK
	if status == 200 {
		var returnValue datav1alpha1.PluginSpec
		json.Unmarshal(response, &returnValue)
		return &returnValue
	}

	return nil
}

func (this PluginDAO) GetAll() []datav1alpha1.PluginSpec {
	status, response := httpClient.Get(httpClient.GetBaseURL() + "/plugins/")

	// OK
	if status == 200 {
		var returnValue []datav1alpha1.PluginSpec
		json.Unmarshal(response, &returnValue)
		return returnValue
	}

	return nil
}
