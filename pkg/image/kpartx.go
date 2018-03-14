package image

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

func addDevMappings(path string) error {
	cmd := exec.Command("kpartx", "-s", "-a", path)
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func getDevMappings(path string) ([]string, error) {
	out, err := exec.Command("kpartx", "-l", path).Output()
	if err != nil {
		return []string{}, err
	}
	return mustParseDevMappings(string(out)), nil
}

func removeDevMappings(path string) error {
	log.Debugf("Unmap %s as a device", path)
	cmd := exec.Command("kpartx", "-d", path)
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

// mustParseDevMappings parse 'kpartx -l <img file>' output and resolves
// what are the '/dev/mapper' paths for the partitions
func mustParseDevMappings(raw string) []string {
	scanner := bufio.NewScanner(strings.NewReader(raw))

	partitions := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" && !strings.Contains(line, "deleted") {
			partitions = append(partitions, line)
		}
	}

	result := make([]string, len(partitions))
	for _, line := range partitions {
		parts := strings.SplitN(line, ":", 2)
		if parts[0] != "" {
			info := strings.SplitN(strings.TrimSpace(parts[1]), " ", 4)
			index, err := strconv.Atoi(info[0])
			if err != nil {
				log.Fatalf("Failed to parse kpartx output line: %s", line)
			}
			result[index] = fmt.Sprintf("/dev/mapper/%s", strings.TrimSpace(parts[0]))
		}
	}
	return result
}
