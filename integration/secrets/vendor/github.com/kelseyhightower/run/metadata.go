package run

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
	"time"
)

const (
	userAgent = "golang-run/0.0.1"
)

// MetadataError reports an error and the endpoint that caused it.
type MetadataError struct {
	Response *http.Response
}

func (e *MetadataError) Error() string {
	if e.Response.StatusCode == 404 {
		return fmt.Sprintf("run/metadata: %s not found", e.Response.Request.URL)
	}

	return fmt.Sprintf("run/metadata: http error %s", e.Response.Status)
}

// AccessToken holds a GCP access token.
type AccessToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

// ProjectID returns the active project ID from the metadata service.
func ProjectID() (string, error) {
	endpoint := "http://metadata.google.internal/computeMetadata/v1/project/project-id"

	data, err := httpRequest(endpoint)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// NumericProjectID returns the active project ID from the metadata service.
func NumericProjectID() (string, error) {
	endpoint := "http://metadata.google.internal/computeMetadata/v1/project/numeric-project-id"

	data, err := httpRequest(endpoint)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// Token returns the default service account token.
func Token(scopes []string) (*AccessToken, error) {
	s := strings.Join(scopes, ",")

	endpoint := fmt.Sprintf("http://metadata.google.internal/computeMetadata/v1/instance/service-accounts/default/token?scopes=%s", s)
	data, err := httpRequest(endpoint)
	if err != nil {
		return nil, err
	}

	var accessToken AccessToken
	err = json.Unmarshal(data, &accessToken)
	if err != nil {
		return nil, err
	}

	return &accessToken, nil
}

// IDToken returns an id token based on the service url.
func IDToken(serviceURL string) (string, error) {
	endpoint := fmt.Sprintf("http://metadata.google.internal/computeMetadata/v1/instance/service-accounts/default/identity?audience=%s", serviceURL)

	idToken, err := httpRequest(endpoint)
	if err != nil {
		return "", fmt.Errorf("metadata.Get: failed to query id_token: %+v", err)
	}
	return string(idToken), nil
}

// Region returns the name of the Cloud Run region.
func Region() (string, error) {
	endpoint := "http://metadata.google.internal/computeMetadata/v1/instance/zone"

	data, err := httpRequest(endpoint)
	if err != nil {
		return "", err
	}

	region := strings.TrimSuffix(path.Base(string(data)), "-1")
	return region, nil
}

func httpRequest(endpoint string) ([]byte, error) {
	request, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("User-Agent", userAgent)
	request.Header.Add("Metadata-Flavor", "Google")

	timeout := time.Duration(3) * time.Second
	client := http.Client{Timeout: timeout}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		return nil, &MetadataError{Response: response}
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	return data, nil
}