package mapper

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Mapping represents a UID/GID mapping between container and host
type Mapping struct {
	ContainerUID int
	ContainerGID int
	HostUID      int
	HostGID      int
}

// ValidateInput validates and parses the input string into container and host UID/GID mappings
func ValidateInput(value string) (*Mapping, error) {
	min, max := 1, 65535
	mapping := &Mapping{}

	container, host := value, value
	if strings.Contains(value, "=") {
		parts := strings.Split(value, "=")
		container, host = parts[0], parts[1]
	}

	containerUIDStr, containerGIDStr := container, container
	if strings.Contains(container, ":") {
		parts := strings.Split(container, ":")
		containerUIDStr, containerGIDStr = parts[0], parts[1]
	}

	hostUIDStr, hostGIDStr := host, host
	if strings.Contains(host, ":") {
		parts := strings.Split(host, ":")
		hostUIDStr, hostGIDStr = parts[0], parts[1]
	}

	uidGIDValidate := func(value string) (int, error) {
		if value == "" {
			return -1, nil
		}
		num, err := strconv.Atoi(value)
		if err != nil || num < min || num > max {
			return 0, fmt.Errorf("value '%s' is not in range %d-%d", value, min, max)
		}
		return num, nil
	}

	var err error
	mapping.ContainerUID, err = uidGIDValidate(containerUIDStr)
	if err != nil {
		return nil, fmt.Errorf("container UID: %w", err)
	}

	mapping.ContainerGID, err = uidGIDValidate(containerGIDStr)
	if err != nil {
		return nil, fmt.Errorf("container GID: %w", err)
	}

	mapping.HostUID, err = uidGIDValidate(hostUIDStr)
	if err != nil {
		return nil, fmt.Errorf("host UID: %w", err)
	}

	mapping.HostGID, err = uidGIDValidate(hostGIDStr)
	if err != nil {
		return nil, fmt.Errorf("host GID: %w", err)
	}

	return mapping, nil
}

// CreateMap generates the LXC mapping strings
func CreateMap(idType string, idList [][2]int) []string {
	var result []string
	for i, ids := range idList {
		containerID, hostID := ids[0], ids[1]
		if i == 0 {
			result = append(result, fmt.Sprintf("lxc.idmap: %s 0 100000 %d", idType, containerID))
		} else if idList[i][0] != idList[i-1][0]+1 {
			previous := idList[i-1]
			rangeEnd := previous[0] + 1
			rangeHost := previous[0] + 100001
			rangeSize := (containerID - 1) - previous[0]
			result = append(result, fmt.Sprintf("lxc.idmap: %s %d %d %d", idType, rangeEnd, rangeHost, rangeSize))
		}
		result = append(result, fmt.Sprintf("lxc.idmap: %s %d %d 1", idType, containerID, hostID))
		if i == len(idList)-1 {
			rangeEnd := containerID + 1
			rangeHost := containerID + 100001
			rangeSize := 65535 - containerID
			result = append(result, fmt.Sprintf("lxc.idmap: %s %d %d %d", idType, rangeEnd, rangeHost, rangeSize))
		}
	}
	return result
}

// ValidateMappingString validates the format of a mapping string
func ValidateMappingString(value string) error {
	if matched, _ := regexp.MatchString(`^\d+(:\d+)?(=\d+(:\d+)?)?$`, value); !matched {
		return errors.New("invalid format")
	}
	return nil
}
