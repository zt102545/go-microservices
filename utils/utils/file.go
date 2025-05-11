package utils

import (
	"os"
	"regexp"
	"strings"
)

func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil || os.IsExist(err)
}

func MakeFilename(prefix string, objectName string, resolution string, attribute string, extension string) string {
	if objectName != "" {
		reg := regexp.MustCompile(`(?i)[^\da-z_]`)
		objectName = "_" + strings.ToLower(reg.ReplaceAllString(objectName, ""))
	}

	resolutionReg := regexp.MustCompile(`^\d+x$`)
	if resolutionReg.MatchString(resolution) {
		resolution = "@" + resolution
	} else {
		resolution = ""
	}

	return prefix + objectName + attribute + resolution + "." + extension
}

func MakeImageUrl(dir string, name string, resolution string, theme string, screen string) string {
	theme_dir := ""

	if theme != "" {
		theme_dir = theme + "/"
	}

	screen_dir := ""

	if screen != "" {
		screen_dir = screen + "/"
	}

	return dir + "/" + theme_dir + resolution + "/" + screen_dir + name
}
