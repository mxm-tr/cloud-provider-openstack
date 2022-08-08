package openstack

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	defaultMetadataVersion = "2012-08-10"
	metadataURLTemplate    = "http://169.254.169.254/openstack/%s/meta_data.json"
)

// IMetadata implements GetInstanceID
type IMetadata interface {
	GetInstanceID() (string, error)
}

type metadata struct {
	UUID string
}

// MetadataService instance of IMetadata
var MetadataService IMetadata

// GetMetadataProvider retrieves instance of IMetadata
func GetMetadataProvider() (IMetadata, error) {

	if MetadataService == nil {
		MetadataService = &metadata{}
	}
	return MetadataService, nil
}

func getMetadata(metadataURL string) ([]byte, error) {
	resp, err := http.Get(metadataURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, err
	}

	md, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return md, nil
}

func getMetaDataInfo() (metadata, error) {
	metadataURL := fmt.Sprintf(metadataURLTemplate, defaultMetadataVersion)
	var m metadata
	md, err := getMetadata(metadataURL)
	if err != nil {
		return m, err
	}
	err = json.Unmarshal(md, &m)
	if err != nil {
		return m, err
	}
	return m, nil
}

// GetInstanceID from metadata service
func (m *metadata) GetInstanceID() (string, error) {
	md, err := getMetaDataInfo()
	if err != nil {
		return "", err
	}
	return md.UUID, nil
}
