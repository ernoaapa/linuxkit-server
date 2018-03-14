package image

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"sort"
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

	lines := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" && !strings.Contains(line, "deleted") {
			lines = append(lines, line)
		}
	}

	partitions := make([]struct {
		start int
		path  string
	}, len(lines))

	for i, line := range lines {
		parts := strings.SplitN(line, ":", 2)
		if parts[0] != "" {
			info := strings.SplitN(strings.TrimSpace(parts[1]), " ", 4)
			size, err := strconv.Atoi(info[3])
			if err != nil {
				log.Fatalf("Failed to parse kpartx output line: %s", line)
			}
			partitions[i] = struct {
				start int
				path  string
			}{
				start: size,
				path:  fmt.Sprintf("/dev/mapper/%s", strings.TrimSpace(parts[0])),
			}
		}
	}
	sort.Slice(partitions, func(i, j int) bool { return partitions[i].start < partitions[j].start })

	result := make([]string, len(partitions))
	for i, device := range partitions {
		result[i] = device.path
	}
	return result
}
