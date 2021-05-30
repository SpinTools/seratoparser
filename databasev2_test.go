package seratoparser

import (
"log"
	"strconv"
	"testing"
)

func TestReadDatabase(t *testing.T) {
	p := New(SeratoDir)
	mediaEntities := p.GetAllTracks()
	if len(mediaEntities) == 0 {
		t.Errorf("GetAllTracks() = %q, want %q", strconv.Itoa(len(mediaEntities)), ">0")
		log.Println(mediaEntities)
	}
}

