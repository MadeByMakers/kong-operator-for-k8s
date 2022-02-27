package dao

import (
	"bytes"
	"encoding/json"

	datav1alpha1 "github.com/MadeByMakers/kong-operator-for-k8s/api/v1alpha1"
)

type Plugin4Rest struct {
	Id        string                 `json:"id,omitempty"`
	Name      string                 `json:"name,omitempty"`
	Config    map[string]interface{} `json:"config,omitempty"`
	Tags      []string               `json:"tags,omitempty"`
	Route     *datav1alpha1.ObjectId `json:"route"`
	Service   *datav1alpha1.ObjectId `json:"service"`
	Protocols []string               `json:"protocols,omitempty"`
	Consumer  []string               `json:"consumer,omitempty"`
	Enabled   bool                   `json:"enabled,omitempty"`
}

func (plugin *Plugin4Rest) ToSpec() datav1alpha1.PluginSpec {
	retorno := datav1alpha1.PluginSpec{}

	configJson := new(bytes.Buffer)
	json.NewEncoder(configJson).Encode(plugin.Config)

	retorno.Config = configJson.String()

	retorno.Id = plugin.Id
	retorno.Name = plugin.Name
	retorno.Tags = plugin.Tags
	retorno.Route = plugin.Route
	retorno.Service = plugin.Service
	retorno.Protocols = plugin.Protocols
	retorno.Consumer = plugin.Consumer
	retorno.Enabled = plugin.Enabled

	return retorno
}

func (this Plugin4Rest) FromSpec(spec datav1alpha1.PluginSpec) Plugin4Rest {
	var x map[string]interface{}
	json.Unmarshal([]byte(spec.Config), &x)

	this.Config = x

	this.Id = spec.Id
	this.Name = spec.Name
	this.Tags = spec.Tags
	this.Route = spec.Route
	this.Service = spec.Service
	this.Protocols = spec.Protocols
	this.Consumer = spec.Consumer
	this.Enabled = spec.Enabled

	if this.Protocols == nil {
		this.Protocols = []string{"grpc", "grpcs", "http", "https"}
	}

	return this
}
