package Services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"quay-go-api/Common/Errors"
	"quay-go-api/Entities/Dto"
	"quay-go-api/Services/Auth"
	logger "quay-go-api/Services/Logger"
	"time"
)

/*
QuayGetManifest Call the official Quay API, /api/v1/repository/{repository}/manifest/{manifestRef} to get the manifest information
*/
func QuayGetManifest(repositoryNamespaced string, manifestRef string, currentUser Auth.AuthenticatedUser) (Dto.Manifest, error) {
	// Get the quay url
	if quayUrl := os.Getenv("QUAY_URL"); quayUrl == "" {
		logger.Warning("QUAY_URL not set")
		return Dto.Manifest{}, Errors.QuayUrlNotSet()
	}

	// Prepare the request to Quay
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v1/repository/%s/manifest/%s", os.Getenv("QUAY_URL"), repositoryNamespaced, manifestRef), nil)
	if err != nil {
		logger.Error("Error creating request: %s", err.Error())
		return Dto.Manifest{}, err
	}

	// Add authorization header
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", currentUser.Token)) // Like OBO in Microsoft authentication workflow lol

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("Error making request: %s", err.Error())
		return Dto.Manifest{}, err
	}
	defer resp.Body.Close()

	// Ensure success
	if resp.StatusCode != http.StatusOK {
		logger.Error("Quay API returned status: %d", resp.StatusCode)
		bodyBytes, _ := io.ReadAll(resp.Body)
		return Dto.Manifest{}, Errors.QuayApiError(resp.StatusCode, string(bodyBytes))
	}

	// Parse response into Dto
	var body map[string]json.RawMessage
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		logger.Error("Error decoding response: %s", err.Error())
		return Dto.Manifest{}, err
	}

	var digest string
	_ = json.Unmarshal(body["manifest"], &digest)
	var IsManifestList bool
	_ = json.Unmarshal(body["is_manifest_list"], &IsManifestList)
	var manifestData string
	_ = json.Unmarshal(body["manifest_data"], &manifestData)
	var configMediaType string
	_ = json.Unmarshal(body["config_media_type"], &configMediaType)
	var rawLayers []json.RawMessage
	_ = json.Unmarshal(body["layers"], &rawLayers)

	var layers []Dto.ManifestLayer
	for _, raw := range rawLayers {
		var layer map[string]json.RawMessage
		_ = json.Unmarshal(raw, &layer)

		var index int
		_ = json.Unmarshal(layer["index"], &index)
		var compressedSize int
		_ = json.Unmarshal(layer["compressed_size"], &compressedSize)
		var isRemote bool
		_ = json.Unmarshal(layer["is_remote"], &isRemote)
		var urls *[]string
		_ = json.Unmarshal(layer["urls"], &urls)
		var command []string
		_ = json.Unmarshal(layer["command"], &command)
		var comment string
		_ = json.Unmarshal(layer["comment"], &comment)
		var author *string
		_ = json.Unmarshal(layer["author"], &author)
		var blobDigest string
		_ = json.Unmarshal(layer["blob_digest"], &blobDigest)
		var created string
		_ = json.Unmarshal(layer["created_datetime"], &created)
		createdDatetime, _ := time.Parse(time.RFC1123Z, created)

		layers = append(layers, Dto.ManifestLayer{
			index,
			compressedSize,
			isRemote,
			urls,
			command,
			comment,
			author,
			blobDigest,
			createdDatetime,
		})
	}

	// reconstruct the manifest object
	manifest := Dto.Manifest{
		Digest:          digest,
		IsManifestList:  IsManifestList,
		ManifestData:    manifestData,
		ConfigMediaType: &configMediaType,
		Layers:          layers,
	}

	return manifest, nil
}

func QuayGetManifestSecScan(repositoryNamespaced string, manifestRef string, currentUser Auth.AuthenticatedUser) (Dto.SecScanReport, error) {
	// Get the quay url
	if quayUrl := os.Getenv("QUAY_URL"); quayUrl == "" {
		logger.Warning("QUAY_URL not set")
		return Dto.SecScanReport{}, Errors.QuayUrlNotSet()
	}

	// Prepare the request to Quay
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v1/repository/%s/manifest/%s/security", os.Getenv("QUAY_URL"), repositoryNamespaced, manifestRef), nil)
	if err != nil {
		logger.Error("Error creating request: %s", err.Error())
		return Dto.SecScanReport{}, err
	}

	// Add authorization header
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", currentUser.Token)) // Like OBO in Microsoft authentication workflow lol

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("Error making request: %s", err.Error())
		return Dto.SecScanReport{}, err
	}
	defer resp.Body.Close()

	// Ensure success
	if resp.StatusCode != http.StatusOK {
		logger.Error("Quay API returned status: %d", resp.StatusCode)
		bodyBytes, _ := io.ReadAll(resp.Body)
		return Dto.SecScanReport{}, Errors.QuayApiError(resp.StatusCode, string(bodyBytes))
	}

	var secScanReport Dto.SecScanReport
	if err := json.NewDecoder(resp.Body).Decode(&secScanReport); err != nil {
		logger.Error("Error decoding response: %s", err.Error())
		return Dto.SecScanReport{}, err
	}

	return secScanReport, nil
}
