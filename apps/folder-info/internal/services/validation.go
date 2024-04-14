package services

import "strings"

// isPathAllowed determines if path is prefixed with on of the Core.rootDirs
func (c Core) isPathAllowed(path string) bool {
	for _, rootDir := range c.rootDirs {
		if strings.HasPrefix(path, rootDir) {
			return true
		}
	}

	return false
}
