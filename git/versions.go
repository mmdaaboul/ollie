package git

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/log"
)

var DATEVER = "DATEVER"
var SEMVER = "SEMVER"

func ParseVersion(version string) (int, int, int, int, string, string, error) {
	version = strings.TrimPrefix(version, "v") // Remove the 'v' prefix
	parts := strings.Split(version, ".")       // Split by dot

	if len(parts) < 3 {
		return 0, 0, 0, 0, "", "", errors.New("invalid version format")
	}

	// Determine if it's semver or datever
	var versionType string
	if len(parts[0]) == 4 {
		versionType = DATEVER
	} else {
		versionType = SEMVER
	}

	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, 0, 0, "", "", err
	}

	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, 0, 0, "", "", err
	}

	var patch, deployment int
	var tag string

	if versionType == SEMVER {
		// Handle patch and optional tag
		subParts := strings.SplitN(parts[2], "-", 2)
		deployment = 0
		patch, err = strconv.Atoi(subParts[0])
		if err != nil {
			return 0, 0, 0, 0, "", "", err
		}
		if len(subParts) > 1 {
			tag = subParts[1]
		}
	} else {
		subParts := strings.SplitN(parts[3], "-", 2)
		patch, err = strconv.Atoi(parts[2])
		if err != nil {
			return 0, 0, 0, 0, "", "", err
		}
		deployment, err = strconv.Atoi(subParts[0])
		if err != nil {
			fmt.Printf("Deployment failed: %s", err)
			return 0, 0, 0, 0, "", "", err
		}
		if len(subParts) > 1 {
			tag = subParts[1]
		}
	}

	return major, minor, patch, deployment, tag, versionType, nil
}

func VersionBump(version string, bump string, toProd, toStaging bool) (string, error) {
	major, minor, patch, deployment, tag, versionType, err := ParseVersion(version)

	if err != nil {
		log.Fatal(fmt.Sprintf("Version not in correct format: %s", err))
		return "", err
	}

	if toProd {
		tag = ""
	} else if toStaging {
		tag = "-rc"
	} else {
		tag = fmt.Sprintf("-%s", tag)
	}

	if versionType == DATEVER {
		deployment += 1
		return fmt.Sprintf("v%d.%d.%d.%d%s", major, minor, patch, deployment, tag), nil
	}

	switch bump {
	case "major":
		major += 1
	case "minor":
		minor += 1
	case "patch":
		patch += 1
	case "same":
		tag = IncrementString(tag)
	default:
		log.Fatal("Incorrect bump type")
		return "", errors.New("Incorrect bump type")
	}

	return fmt.Sprintf("v%d.%d.%d%s", major, minor, patch, tag), nil
}

func IncrementString(str string) string {
	// Find the last digit
	for i := len(str) - 1; i >= 0; i-- {
		if str[i] >= '0' && str[i] <= '9' {
			// Extract the numeric part
			numStr := str[i:]
			var num int
			fmt.Sscanf(numStr, "%d", &num)

			// Increment and convert back to string
			num++
			return str[:i] + fmt.Sprintf("%d", num)
		}
	}
	// No digits found, append 1
	return str + "1"
}
