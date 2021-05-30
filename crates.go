package serato_parser

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
)


// GetCrates returns files of all crates found in Serato Path
// TODO: Should we parse meta data of theses files, or only provide OS level elements. Weird to be a parser library and no parsing.
// TODO: Is there Serato Crate meta data?
func (p Parser) GetCrates() []os.FileInfo {
	var crateFiles []os.FileInfo
	var subcratesPaths []string

	seratoFolder := filepath.FromSlash(p.FilePath + "/Subcrates")
	subcratesPaths = append(subcratesPaths, seratoFolder)
	seratoFiles, _ := ioutil.ReadDir(seratoFolder)
	for _,seratoFile := range seratoFiles {
		fileExt := filepath.Ext(seratoFile.Name())
		if fileExt != ".crate" {
			continue
		}

		filePath := filepath.FromSlash(seratoFolder + "/" + seratoFile.Name())
		filePath,_ = filepath.Abs(filePath)
		crateFiles = append(crateFiles, seratoFile)
	}

	sort.Slice(crateFiles, func(i, j int) bool {
		return len(crateFiles[i].Name()) < len(crateFiles[j].Name())
	})

	return crateFiles
}


// GetCrateTracks takes a filename and returns all the tracks/entities inside the crate
func (p Parser) GetCrateTracks (fileName string) []MediaEntity {
	return readMediaEntities(filepath.FromSlash(p.FilePath + "/Subcrates/" + fileName))
}