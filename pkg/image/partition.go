package image

import (
	"errors"
	"fmt"
	"strings"
)

// The default partition sector size
var sectorSize = uint64(512)
var startOffset = sectorSize * 2048
var endOffset = sectorSize

var supportedTypes = []string{"fat32", "ext4"}

func validatePartitionTable(table []Partition, bootSize uint64) error {
	if len(table) > 4 {
		return errors.New("Only max 4 partitions are supported")
	}

	for _, partition := range table {
		if !contains(supportedTypes, partition.FsType) {
			return fmt.Errorf("%s is unsupported partition file system type. Supported types are %s", partition.FsType, strings.Join(supportedTypes, ","))
		}
	}

	bootPartitions := getBootPartitions(table)
	if len(bootPartitions) == 0 {
		return errors.New("Partition table must have at least one boot partition")
	} else if len(bootPartitions) > 1 {
		return errors.New("Partition table can contain only one partition with boot=true")
	}

	boot := bootPartitions[0]

	if boot.FsType != "fat32" {
		return errors.New("Boot partition file system type must be fat32")
	}

	if boot.Size < bootSize {
		return errors.New("Boot partition is not large enough for the boot files")
	}

	return nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func getBootPartitions(table []Partition) (result []Partition) {
	for _, p := range table {
		if p.Boot {
			result = append(result, p)
		}
	}
	return result
}

func getTotalSize(table []Partition) uint64 {
	var last Partition
	for _, p := range table {
		if last.Start < p.Start {
			last = p
		}
	}
	return last.Start + last.Size + sectorSize
}
