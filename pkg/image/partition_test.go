package image

import (
	"fmt"
	"testing"

	"github.com/c2h5oh/datasize"
	"github.com/stretchr/testify/assert"
)

func TestValidatePartitionTableBootFlag(t *testing.T) {
	assert.NoError(t, validatePartitionTable([]Partition{
		{Start: startOffset, Size: uint64(datasize.MB * 100), Boot: true, FsType: "fat32"},
	}, uint64(datasize.MB*50)), "should be valid if single boot partition")

	assert.Error(t, validatePartitionTable([]Partition{}, 0), "should return error if no boot partition")
	assert.Error(t, validatePartitionTable([]Partition{
		{Start: startOffset, Size: uint64(datasize.MB * 100), Boot: true, FsType: "fat32"},
		{Start: startOffset + uint64(datasize.MB*100), Size: uint64(datasize.MB * 100), Boot: true, FsType: "fat32"},
	}, uint64(datasize.MB*50)), "should return error if have two boot partitions")
}

func TestValidatePartitionTableBootType(t *testing.T) {
	assert.NoError(t, validatePartitionTable([]Partition{
		{Start: startOffset, Size: uint64(datasize.MB * 100), Boot: true, FsType: "fat32"},
	}, uint64(datasize.MB*50)), "should be valid if boot file system type is fat32")

	assert.Error(t, validatePartitionTable([]Partition{
		{Start: startOffset, Size: uint64(datasize.MB * 100), Boot: true, FsType: "ext4"},
	}, uint64(datasize.MB*50)), "should be invalid if boot file system type is not fat32")
}

func TestValidatePartitionTableBootSize(t *testing.T) {
	assert.NoError(t, validatePartitionTable([]Partition{
		{Start: startOffset, Size: uint64(datasize.MB * 100), Boot: true, FsType: "fat32"},
	}, uint64(datasize.MB*50)), "should be valid if boot filesystem is large enough for boot files")

	assert.Error(t, validatePartitionTable([]Partition{
		{Start: startOffset, Size: uint64(datasize.MB * 100), Boot: true, FsType: "fat32"},
	}, uint64(datasize.MB*200)), "should be invalid if boot filesystem is not large enough for boot files")
}

func TestValidatePartitionTableAllowMaxFourPrimaryPartitions(t *testing.T) {
	assert.NoError(t, validatePartitionTable([]Partition{
		{Start: startOffset, Size: uint64(datasize.MB * 100), Boot: true, FsType: "fat32"},
		{Start: startOffset + uint64(datasize.MB*100), Size: uint64(datasize.MB * 100), FsType: "ext4"},
		{Start: startOffset + uint64(datasize.MB*200), Size: uint64(datasize.MB * 100), FsType: "ext4"},
		{Start: startOffset + uint64(datasize.MB*300), Size: uint64(datasize.MB * 100), FsType: "ext4"},
	}, uint64(datasize.MB*50)), "should be valid if four partitions")

	assert.Error(t, validatePartitionTable([]Partition{
		{Start: startOffset, Size: uint64(datasize.MB * 100), Boot: true, FsType: "fat32"},
		{Start: startOffset + uint64(datasize.MB*100), Size: uint64(datasize.MB * 100), FsType: "ext4"},
		{Start: startOffset + uint64(datasize.MB*200), Size: uint64(datasize.MB * 100), FsType: "ext4"},
		{Start: startOffset + uint64(datasize.MB*300), Size: uint64(datasize.MB * 100), FsType: "ext4"},
		{Start: startOffset + uint64(datasize.MB*400), Size: uint64(datasize.MB * 100), FsType: "ext4"},
	}, uint64(datasize.MB*50)), "should be valid if over four partitions")
}

func TestValidatePartitionFsTypes(t *testing.T) {
	for _, fsType := range []string{"fat32", "ext4"} {
		assert.NoError(t, validatePartitionTable([]Partition{
			{Start: startOffset, Size: uint64(datasize.MB * 100), Boot: true, FsType: "fat32"},
			{Start: startOffset + uint64(datasize.MB*100), Size: uint64(datasize.MB * 100), FsType: fsType},
		}, uint64(datasize.MB*50)), fmt.Sprintf("should support %s file system type partitions", fsType))
	}

	for _, fsType := range []string{"fat16", "ext2"} {
		assert.Error(t, validatePartitionTable([]Partition{
			{Start: startOffset, Size: uint64(datasize.MB * 100), Boot: true, FsType: "fat32"},
			{Start: startOffset + uint64(datasize.MB*100), Size: uint64(datasize.MB * 100), FsType: fsType},
		}, uint64(datasize.MB*50)), fmt.Sprintf("should not support %s file system type partitions", fsType))
	}
}

func TestGetTotalSize(t *testing.T) {
	assert.Equal(t,
		startOffset+uint64(datasize.MB*100)+endOffset,
		getTotalSize([]Partition{
			{Start: startOffset, Size: uint64(datasize.MB * 100), Boot: true, FsType: "fat32"},
		}),
		"Should include start offset into total size")

	assert.Equal(t,
		startOffset+uint64(datasize.MB*100)+uint64(datasize.MB*100)+endOffset,
		getTotalSize([]Partition{
			{Start: startOffset, Size: uint64(datasize.MB * 100), Boot: true, FsType: "fat32"},
			{Start: startOffset + uint64(datasize.MB*100), Size: uint64(datasize.MB * 100), Boot: true, FsType: "ext4"},
		}),
		"Should sum all partitions to total size")
}
