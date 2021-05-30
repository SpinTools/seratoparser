package seratoparser

import "path/filepath"

// GetAllTracks returns all the tracks/entities inside the Database
func (p Parser) GetAllTracks() []MediaEntity {
        return readMediaEntities(filepath.FromSlash(p.FilePath + "/database V2"))
}
