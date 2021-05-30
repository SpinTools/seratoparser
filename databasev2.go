package serato_parser

import "path/filepath"

func (p Parser) GetAllTracks () []MediaEntity {
	return readMediaEntities(filepath.FromSlash(p.FilePath + "/database V2"))
}