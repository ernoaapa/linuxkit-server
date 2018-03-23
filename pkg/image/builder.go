package image

import (
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/ernoaapa/linuxkit-server/pkg/utils"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func Build(sourceDir string, partitions []Partition, w io.Writer) error {
	tmpfile, err := ioutil.TempFile("", "img")
	if err != nil {
		return errors.Wrapf(err, "Failed to create temporary file")
	}
	defer tmpfile.Close()
	defer os.Remove(tmpfile.Name())

	if err := buildImage(sourceDir, partitions, tmpfile.Name()); err != nil {
		return err
	}

	if _, err := io.Copy(w, tmpfile); err != nil {
		return errors.Wrapf(err, "Failed to copy img to writer")
	}

	return nil
}

func buildImage(sourceDir string, partitions []Partition, filename string) error {
	bootSize, err := utils.GetDirSize(sourceDir)
	if err != nil {
		return errors.Wrapf(err, "Failed to resolve image source size")
	}

	if err := validatePartitionTable(partitions, uint64(bootSize)); err != nil {
		return err
	}

	if err := createZeroFile(filename, getTotalSize(partitions)); err != nil {
		return err
	}
	if err := createPartitions(filename, partitions); err != nil {
		return err
	}

	if err := addDevMappings(filename); err != nil {
		return err
	}
	defer removeDevMappings(filename)

	devices, err := getDevMappings(filename)
	if err != nil {
		return err
	}

	for i, device := range devices {
		partition := partitions[i]

		switch partition.FsType {
		case "fat32":
			if err := formatFat32(device); err != nil {
				return errors.Wrapf(err, "Failed to format device %s as Fat32", device)
			}

		case "ext4":
			if err := formatExt4(device); err != nil {
				return errors.Wrapf(err, "Failed to format device %s as ext4", device)
			}
		}

		if partition.Boot {
			if err := writeBootPartition(device, sourceDir); err != nil {
				return err
			}
		}
	}

	return nil
}

func writeBootPartition(device, sourceDir string) error {
	buildDir, err := ioutil.TempDir("/mnt", "")
	if err != nil {
		return errors.Wrapf(err, "Failed to create temporary build directory")
	}
	defer os.RemoveAll(buildDir)

	if err := mountDevice(device, buildDir); err != nil {
		return errors.Wrapf(err, "Failed to mount root partition %s to dir %s", device, buildDir)
	}
	defer unmountDevice(buildDir)

	if err := utils.Copy(sourceDir, buildDir, true); err != nil {
		return errors.Wrapf(err, "Failed to copy files from %s to %s", sourceDir, buildDir)
	}

	return nil
}

func createZeroFile(path string, size uint64) error {
	log.Debugf("Create %d bytes empty zero file to %s", size, path)
	f, err := os.Create(path)
	if err != nil {
		return errors.Wrapf(err, "Failed to create empty file to %s", path)
	}
	defer f.Close()

	if err := f.Truncate(int64(size)); err != nil {
		return errors.Wrapf(err, "Failed to fill file %s to size %d", path, size)
	}
	return nil
}

func createPartitions(path string, table []Partition) error {
	log.Debugf("Create partitions to %s", path)
	if err := runParted(path, "mklabel", "msdos"); err != nil {
		return errors.Wrapf(err, "Failed to execute 'parted mklabel' to %s", path)
	}
	for i, partition := range table {
		args := []string{"--script", "--align=opt", path, "mkpart", "primary", partition.FsType, strconv.FormatUint(partition.Start, 10) + "B", strconv.FormatUint(partition.Start+partition.Size, 10) + "B"}
		if err := runParted(args...); err != nil {
			log.Debugf("Failed to execute (parted %s): %s", strings.Join(args, " "), err)
			return errors.Wrapf(err, "Failed to create %d/%d partition to %s", i, len(table), path)
		}
	}
	return nil
}

func runParted(args ...string) error {
	cmd := exec.Command("parted", args...)
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func formatFat32(device string) error {
	log.Debugf("Format device %s", device)
	cmd := exec.Command("mkfs.vfat", "-F", "32", device)
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func formatExt4(device string) error {
	log.Debugf("Format device %s", device)
	cmd := exec.Command("mkfs.ext4", device)
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func mountDevice(device, path string) error {
	log.Debugf("Mount device %s to path %s", device, path)
	cmd := exec.Command("mount", device, path)
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func unmountDevice(path string) error {
	log.Debugf("Unmount device at path %s", path)
	cmd := exec.Command("umount", path)
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func appendSuffix(str, suffix string) string {
	if strings.HasSuffix(str, suffix) {
		return str
	}
	return str + suffix
}
