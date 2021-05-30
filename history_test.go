package serato_parser

import (
	"log"
	"strconv"
	"testing"
)

func TestReadHistorySession(t *testing.T) {
	p := New(SERATO_DIR)
	sessions := p.GetHistorySessions()
	if len(sessions) == 0 {
		t.Errorf("GetHistorySessions() = %q, want %q", strconv.Itoa(len(sessions)), ">0")
		log.Println(sessions)
	}

	log.Println(sessions[0].Name())
	historyEntities := p.ReadHistorySession(sessions[0].Name())
	if len(historyEntities) == 0 {
		t.Errorf("ReadHistorySession() = %q, want %q", strconv.Itoa(len(historyEntities)), ">0")
		log.Println(historyEntities)
	}
}