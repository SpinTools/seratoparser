package seratoparser

import (
	"log"
	"strconv"
	"testing"
)

func TestReadCrates(t *testing.T) {
	p := New(SeratoDir)
	crates := p.GetCrates()
	if len(crates) == 0 {
		t.Errorf("GetCrates() = %q, want %q", strconv.Itoa(len(crates)), ">0")
		log.Println(crates)
	}

	foundTracks := false
	for _,crate := range crates {
		mediaEntities := p.GetCrateTracks(crate.Name())
		if len(mediaEntities) > 0 {
			foundTracks = true
			break
		}
	}
	if !foundTracks {
		t.Errorf("GetCrateTracks() = %q, want %q", "0", ">0")
	}
}
