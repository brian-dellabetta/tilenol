package tilenol

import (
	"context"
	"errors"
	"fmt"

	"github.com/paulmach/orb/geojson"
)

var (
	MultipleSourcesErr = errors.New("Layers can only support a single backend source")
	NoSourcesErr       = errors.New("Layers must have a single backend source configured")
)

// SourceConfig represents a generic YAML source configuration object
type SourceConfig struct {
	// Elasticsearch is an optional YAML key for configuring an ElasticsearchConfig
	Elasticsearch *ElasticsearchConfig `yaml:"elasticsearch"`
	// PostGIS is an optional YAML key for configuring a PostGISConfig
	PostGIS *PostGISConfig `yaml:"postgis"`
}

// LayerConfig represents a general YAML layer configuration object
type LayerConfig struct {
	// Name is the effective name of the layer
	Name string `yaml:"name"`
	// Description is an optional short descriptor of the layer
	Description string `yaml:"description"`
	// Minzoom specifies the minimum z value for the layer
	Minzoom int `yaml:"minzoom"`
	// Maxzoom specifies the maximum z value for the layer
	Maxzoom int `yaml:"maxzoom"`
	// Source configures the underlying Source for the layer
	Source SourceConfig `yaml:"source"`
}

// Source is a generic interface for all feature data sources
type Source interface {
	// GetFeatures retrieves the GeoJSON FeatureCollection for the given request
	GetFeatures(context.Context, *TileRequest) (*geojson.FeatureCollection, error)
}

// Layer is a configured, hydrated tile server layer
type Layer struct {
	Name        string
	Description string
	Minzoom     int
	Maxzoom     int
	Source      Source
}

// CreateLayer creates a new Layer given a LayerConfig
func CreateLayer(layerConfig LayerConfig) (*Layer, error) {
	layer := &Layer{
		Name:        layerConfig.Name,
		Description: layerConfig.Description,
		Minzoom:     layerConfig.Minzoom,
		Maxzoom:     layerConfig.Maxzoom,
	}
	// TODO: How can we make this more generic?
	if layerConfig.Source.Elasticsearch != nil && layerConfig.Source.PostGIS != nil {
		return nil, MultipleSourcesErr
	}
	if layerConfig.Source.Elasticsearch == nil && layerConfig.Source.PostGIS == nil {
		return nil, NoSourcesErr
	}
	if layerConfig.Source.Elasticsearch != nil {
		source, err := NewElasticsearchSource(layerConfig.Source.Elasticsearch)
		if err != nil {
			return nil, err
		}
		layer.Source = source
		return layer, nil
	}
	if layerConfig.Source.PostGIS != nil {
		source, err := NewPostGISSource(layerConfig.Source.PostGIS)
		if err != nil {
			return nil, err
		}
		layer.Source = source
		return layer, nil
	}
	return nil, fmt.Errorf("Invalid layer source config for layer: %s", layerConfig.Name)
}
