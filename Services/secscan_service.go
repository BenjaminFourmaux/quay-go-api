package Services

import (
	"quay-go-api/Common"
	"quay-go-api/Common/Errors"
	"quay-go-api/Entities/Dto"
	"quay-go-api/Repositories"
	"quay-go-api/Services/Auth"
	logger "quay-go-api/Services/Logger"
	"strconv"
)

/*
Dev Notes:
Quay provides sec scan info (vulnerabilities score and CVEs like Docker Scout) via Clair (https://github.com/quay/clair) a static scanning tool.
Quay API sent manifest blob to Clair for scanning.
Clair saves in a db the scan result and Quay API can retrieve it or launch a new scan (via a queue). It's for that Quay API can respond to different status
*/

func GetRepositoryManifestSecScan(repositoryNamespaced string, manifestRef string, filters map[string]string, currentUser *Auth.AuthenticatedUser) (Dto.SecScanReport, error) {
	logger.Info("[SecScan Service] Get Manifest SecScan")
	logger.Debug("Repository name: %s", repositoryNamespaced)
	logger.Debug("Manifest ref: %s", manifestRef)
	logger.Debug("With filters: %v", filters)

	// Validating filters
	var filterIncludeCVEs bool
	if includeCVEs, ok := filters["include_cve"]; ok {
		isIncludeCVEs, err := strconv.ParseBool(includeCVEs)
		if err != nil {
			logger.Warning("Invalid filter include_cve value: %s", includeCVEs)
			return Dto.SecScanReport{}, Errors.InvalidParameterValue("include_cve", []string{"true", "false"})
		}
		filterIncludeCVEs = isIncludeCVEs
	}

	// Split repositoryNamespaced into namespace and name
	namespace, reponame, err := Common.SplitRepositoryNamespaced(repositoryNamespaced)
	if err != nil {
		logger.Warning("Invalid repository namespaced: %s", repositoryNamespaced)
		return Dto.SecScanReport{}, Errors.RepositoryInvalid(repositoryNamespaced)
	}

	// Check if the namespace (org or user) exists
	if namespace != nil {
		_, err = Repositories.GetUserOrOrganizationByName(*namespace)
		if err != nil {
			switch err.Error() {
			case "record not found":
				logger.Warning("No user or organization found with name: %s", *namespace)
				return Dto.SecScanReport{}, Errors.RepositoryNamespaceNotFound(*namespace)
			default:
				logger.Error("Error retrieving repository  from database: %s", err.Error())
				return Dto.SecScanReport{}, err
			}
		}
	}

	// Check if the repository exits
	repoExist, err := Repositories.FindRepositoryByNameAndNamespace(reponame, namespace)
	if err != nil {
		switch err.Error() {
		case "record not found":
			logger.Warning("No repository '%s' found", repositoryNamespaced)
			return Dto.SecScanReport{}, Errors.RepositoryNotFound(repositoryNamespaced)
		default:
			logger.Error("Error retrieving repository  from database: %s", err.Error())
			return Dto.SecScanReport{}, err
		}
	}

	// Get the manifest and check if exists
	_, err = Repositories.GetRepositoryManifestByDigest(repoExist.ID, manifestRef)
	if err != nil {
		switch err.Error() {
		case "record not found":
			logger.Warning("No manifest '%s' found in repository '%s'", manifestRef, repositoryNamespaced)
			return Dto.SecScanReport{}, Errors.ManifestNotFound(manifestRef, repositoryNamespaced)
		default:
			logger.Error("Error retrieving manifest from database: %s", err.Error())
			return Dto.SecScanReport{}, err
		}
	}

	// Get the sec scan report for the manifest
	// TODO: implement the logic to get real secscan report via the Clair endpoint
	secScanReport, err := QuayGetManifestSecScan(repositoryNamespaced, manifestRef, *currentUser)
	if err != nil {
		/*logger.Error("Error retrieving sec scan report from database: %s", err.Error())
		return Dto.SecScanReport{}, err*/

		// TODO: just fo testing. this feature is not complete yet. return a dummy report
		secScanReport = generateDummySecScanReport()
	}

	if !filterIncludeCVEs {
		secScanReport.CVEs = nil
	}

	return secScanReport, nil
}

func generateDummySecScanReport() Dto.SecScanReport {
	return Dto.SecScanReport{
		Vulnerabilities: Dto.Vulnerabilities{
			Critical:    1,
			High:        1,
			Medium:      0,
			Low:         0,
			Unspecified: 0,
		},
		CVEs: &[]Dto.VulnerabilityCVE{
			{
				CVEID:            "CVE-2021-12345",
				SeverityScore:    9.3,
				Fixable:          "1.0.1",
				PresentIn:        "layer",
				AffectedPackages: "package1, package2",
				AffectedVersions: "<1.0.1",
				PublishedDate:    Common.ParseTime("2021-01-01T00:00:00Z"),
			},
			{
				CVEID:            "CVE-2021-67890",
				SeverityScore:    7.4,
				Fixable:          "",
				PresentIn:        "base",
				AffectedPackages: "package3",
				AffectedVersions: "<2.0.0",
				PublishedDate:    Common.ParseTime("2021-02-01T00:00:00Z"),
			},
		},
	}
}
