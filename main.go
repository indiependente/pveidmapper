package main

import (
	"errors"
	"flag"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// validateInput validates and parses the input string into container and host UID/GID mappings.
func validateInput(value string) (int, int, int, int, error) {
	min, max := 1, 65535
	var containerUID, containerGID, hostUID, hostGID int

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

	containerUID, err := uidGIDValidate(containerUIDStr)
	if err != nil {
		return 0, 0, 0, 0, fmt.Errorf("container UID: %w", err)
	}

	containerGID, err = uidGIDValidate(containerGIDStr)
	if err != nil {
		return 0, 0, 0, 0, fmt.Errorf("container GID: %w", err)
	}

	hostUID, err = uidGIDValidate(hostUIDStr)
	if err != nil {
		return 0, 0, 0, 0, fmt.Errorf("host UID: %w", err)
	}

	hostGID, err = uidGIDValidate(hostGIDStr)
	if err != nil {
		return 0, 0, 0, 0, fmt.Errorf("host GID: %w", err)
	}

	return containerUID, containerGID, hostUID, hostGID, nil
}

// createMap generates the LXC mapping strings.
func createMap(idType string, idList [][2]int) []string {
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

func main() {
	// Parse user input
	var ids stringArray
	flag.Var(&ids, "id", "containeruid[:containergid][=hostuid[:hostgid]]")
	flag.Parse()

	if len(ids) == 0 {
		fmt.Println("Error: No IDs provided. Use -id to specify mappings.")
		flag.Usage()
		return
	}

	var uidList, gidList [][2]int
	for _, id := range ids {
		containerUID, containerGID, hostUID, hostGID, err := validateInput(id)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		if containerUID != -1 {
			uidList = append(uidList, [2]int{containerUID, hostUID})
		}
		if containerGID != -1 {
			gidList = append(gidList, [2]int{containerGID, hostGID})
		}
	}

	sort.Slice(uidList, func(i, j int) bool { return uidList[i][0] < uidList[j][0] })
	sort.Slice(gidList, func(i, j int) bool { return gidList[i][0] < gidList[j][0] })

	uidMap := createMap("u", uidList)
	gidMap := createMap("g", gidList)

	// Output the mappings
	fmt.Println("# Add to /etc/pve/lxc/<container_id>.conf:")
	for _, line := range uidMap {
		fmt.Println(line)
	}
	for _, line := range gidMap {
		fmt.Println(line)
	}

	fmt.Println("\n# Add to /etc/subuid:")
	for _, uid := range uidList {
		fmt.Printf("root:%d:1\n", uid[1])
	}

	fmt.Println("\n# Add to /etc/subgid:")
	for _, gid := range gidList {
		fmt.Printf("root:%d:1\n", gid[1])
	}
}

// stringArray implements the flag.Value interface for string slices.
type stringArray []string

func (s *stringArray) String() string {
	return strings.Join(*s, ", ")
}

func (s *stringArray) Set(value string) error {
	if matched, _ := regexp.MatchString(`^\d+(:\d+)?(=\d+(:\d+)?)?$`, value); !matched {
		return errors.New("invalid format")
	}
	*s = append(*s, value)
	return nil
}
