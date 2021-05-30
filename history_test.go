package seratoparser

import (
	"log"
	"strconv"
	"testing"
)

func TestReadHistorySession(t *testing.T) {
	p := New(SeratoDir)
	sessions := p.GetHistorySessions()
	if len(sessions) == 0 {
		t.Errorf("GetHistorySessions() = %q, want %q", strconv.Itoa(len(sessions)), ">0")
		log.Println(sessions)
	}
	
	historyEntities := p.ReadHistorySession(sessions[0].Name())
	if len(historyEntities) == 0 {
		t.Errorf("ReadHistorySession() = %q, want %q", strconv.Itoa(len(historyEntities)), ">0")
		log.Println(historyEntities)
	}
}