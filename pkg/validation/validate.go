package validation

import (
	"fmt"
	"log"

	"github.com/moby/tool/src/moby"
)

// IsValid validates the Linuxkit configuration and
// return nil if valid or error if some error found
func IsValid(c moby.Moby) error {
	for index, file := range c.Files {
		log.Println(file)
		if file.Source != "" {
			return fmt.Errorf("Invalid configuration in 'files[%d].source'. You cannot use file source when building in remote build server", index)
		}
	}
	return nil
}
