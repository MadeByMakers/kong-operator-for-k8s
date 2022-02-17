package dao

import (
	"encoding/json"
	"strconv"

	datav1alpha1 "github.com/MadeByMakers/kong-operator-for-k8s/api/v1alpha1"
	httpClient "github.com/MadeByMakers/kong-operator-for-k8s/util"
)

type PluginDAO struct{}

func (this PluginDAO) Create(plugin datav1alpha1.Plugin) datav1alpha1.Plugin {
	status, response := httpClient.Post(httpClient.GetBaseURL()+"/plugins", plugin.Spec)

	// OK
	if status == 201 {
		var returnValue datav1alpha1.PluginSpec
		json.Unmarshal(response, &returnValue)

		plugin.Spec = returnValue
		plugin.Status = datav1alpha1.PluginStatus{
			Message: "CREATED",
			Code:    200,
		}
	} else {
		var stringValue string
		json.Unmarshal(response, &stringValue)

		plugin.Status = datav1alpha1.PluginStatus{
			Message: "ERROR (" + strconv.Itoa(status) + ")",
			Code:    status,
			Response: datav1alpha1.HttpStatus{
				Code: status,
				Body: stringValue,
			},
		}
	}

	return plugin
}

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

func (this PluginDAO) Update(plugin datav1alpha1.Plugin) datav1alpha1.Plugin {
	status, response := httpClient.Patch(httpClient.GetBaseURL()+"/plugins", plugin.Spec)

	// OK
	if status == 200 {
		var returnValue datav1alpha1.PluginSpec
		json.Unmarshal(response, &returnValue)

		plugin.Spec = returnValue
		plugin.Status = datav1alpha1.PluginStatus{
			Code:    200,
			Message: "UPDATED",
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
	} else {
		var stringValue string
		json.Unmarshal(response, &stringValue)
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
	} else {
		var stringValue string
		json.Unmarshal(response, &stringValue)
	}

	return nil
}
