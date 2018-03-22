package image

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/ernoaapa/linuxkit-server/pkg/utils"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func Build(sourceDir string, w io.Writer) error {
	tmpfile, err := ioutil.TempFile("", "img")
	if err != nil {
		return errors.Wrapf(err, "Failed to create temporary file")
	}
	defer tmpfile.Close()
	defer os.Remove(tmpfile.Name())

	if err := buildImage(sourceDir, tmpfile.Name()); err != nil {
		return err
	}

	if _, err := io.Copy(w, tmpfile); err != nil {
		return errors.Wrapf(err, "Failed to copy img to writer")
	}

	return nil
}

func buildImage(sourceDir, filename string) error {
	if err := createZeroFile(filename, 1024*1024*200); err != nil {
		return err
	}
	if err := createFat32Partition(filename); err != nil {
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

	var primary string
	var secondary string
	switch l := len(devices); l {
	case 0:
		return fmt.Errorf("%s don't have any paritions. There must be single partition", filename)
	case 1:
		primary = devices[0]
	case 2:
		primary = devices[0]
		secondary = devices[1]
	default:
		return fmt.Errorf("%s contains too many paritions, but we support only 1-2 partition img files", filename)
	}

	if err := formatFat32(primary); err != nil {
		return errors.Wrapf(err, "Failed to format device %s as Fat32", primary)
	}

	if secondary != "" {
		if err := formatExt4(secondary); err != nil {
			return errors.Wrapf(err, "Failed to format device %s as ext4", secondary)
		}
	}

	buildDir, err := ioutil.TempDir("/mnt", "")
	if err != nil {
		return errors.Wrapf(err, "Failed to create temporary build directory")
	}
	defer os.RemoveAll(buildDir)

	if err := mountDevice(primary, buildDir); err != nil {
		return errors.Wrapf(err, "Failed to mount primary device %s to dir %s", primary, buildDir)
	}
	defer unmountDevice(buildDir)

	if err := utils.Copy(sourceDir, buildDir, true); err != nil {
		return errors.Wrapf(err, "Failed to copy files from %s to %s", sourceDir, buildDir)
	}

	return nil
}

func createZeroFile(path string, size int64) error {
	log.Debugf("Create %d bytes empty zero file to %s", size, path)
	f, err := os.Create(path)
	if err != nil {
		return errors.Wrapf(err, "Failed to create empty file to %s", path)
	}
	defer f.Close()

	if err := f.Truncate(size); err != nil {
		return errors.Wrapf(err, "Failed to fill file %s to size %d", path, size)
	}
	return nil
}

func createFat32Partition(path string) error {
	log.Debugf("Create Fat32 partition to %s", path)
	if err := runParted(path, "mklabel", "msdos"); err != nil {
		return errors.Wrapf(err, "Failed to execute 'parted mklabel' to %s", path)
	}
	if err := runParted("--script", "--align=opt", path, "mkpart", "primary", "fat32", "2048s", "70%"); err != nil {
		return errors.Wrapf(err, "Failed to execute 'parted mkpart' to make boot partition to %s", path)
	}
	if err := runParted("--script", "--align=opt", path, "mkpart", "primary", "ext4", "70%", "100%"); err != nil {
		return errors.Wrapf(err, "Failed to execute 'parted mkpart' to make secondary partition to %s", path)
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
