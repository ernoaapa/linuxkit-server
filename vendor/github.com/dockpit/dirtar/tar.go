package dirtar

import (
	"archive/tar"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

// tar the given directory, paths inside the archive
// are relative to the directory
func Tar(dir string, w io.Writer) error {
	tw := tar.NewWriter(w)

	//walk and tar all files in a dir
	//@from http://stackoverflow.com/questions/13611100/how-to-write-a-directory-not-just-the-files-in-it-to-a-tar-gz-file-in-golang
	visit := func(fpath string, fi os.FileInfo, err error) error {

		//cancel walk if something went wrong
		if err != nil {
			return err
		}

		//skip root
		if fpath == dir {
			return nil
		}

		//dont 'add' dirs to archive
		if fi.IsDir() {
			return nil
		}

		f, err := os.Open(fpath)
		if err != nil {
			return err
		}
		defer f.Close()

		//use relative path inside archive
		rel, err := filepath.Rel(dir, fpath)
		if err != nil {
			return err
		}

		//create header from file info struct
		hdr, err := tar.FileInfoHeader(fi, rel)
		if err != nil {
			return err
		}

		//write header to archive using rel dir
		hdr.Name = rel
		err = tw.WriteHeader(hdr)
		if err != nil {
			return err
		}

		//copy content into archive
		if _, err = io.Copy(tw, f); err != nil {
			return err
		}

		return nil
	}

	//walk the context and create archive
	if err := filepath.Walk(dir, visit); err != nil {
		return err
	}

	return nil
}

// untar the archive into a given directory
func Untar(dir string, r io.Reader) error {

	//check if dir exists and return fi
	dirfi, err := os.Stat(dir)
	if err != nil {
		return err
	}

	//check if not empty
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	if len(files) != 0 {
		return fmt.Errorf("Directory to untar into '%s' is not empty", dir)
	}

	//create a new reader
	tr := tar.NewReader(r)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			// end of tar archive
			break
		}
		if err != nil {
			return err
		}

		switch hdr.Typeflag {
		case tar.TypeDir:
			continue
		default:
			//create and open files
			//@todo this assumes the archives dir seperators are the same?
			path := filepath.Join(dir, hdr.Name)

			//make directory if doesnt exist with the same permissions as the root dir
			os.MkdirAll(filepath.Dir(path), dirfi.Mode())

			//create the actual files
			f, err := os.Create(path)
			if err != nil {
				return err
			}
			defer f.Close()

			//copy tar content into file, effectively untarring
			if _, err := io.Copy(f, tr); err != nil {
				return err
			}
		}
	}

	return nil
}
