# Serato Parser

[![Go Report Card](https://goreportcard.com/badge/github.com/SpinTools/seratoparser)](https://goreportcard.com/report/github.com/SpinTools/seratoparser)


A GoLang library for reading Serato database files.

Data Types Supported:
- [x] Database V2
- [x] Crates
- [ ] History Database
- [x] History Sessions

## Installation

This package can be installed with the go get command:

```bash
go get -u github.com/SpinTools/seratoparser
```

## Usage

```go
func main() {
    // Provide Serato Folder
    p := SeratoParser.New("/Users/Stoyvo/Music/_Serato_")
    
    // Get All Tracks in Serato Database
    Tracks := p.GetAllTracks()
    log.Println("Database V2:", Tracks)
    
    // Get all Crates
    crates := p.GetCrates()
    log.Println("Crates:", crates)
    
    // Read crate and get all Tracks
    mediaEntities := p.GetCrateTracks(crates[0].Name())
    log.Println("Crate Tracks:", mediaEntities)

    // Get all session files
    sessions := p.GetHistorySessions()
    log.Println("History Sessions:", sessions)
    
    // Read History Session
    historyEntities := p.ReadHistorySession(sessions[0].Name())
    log.Println("History Tracks:", historyEntities)
}
```

## Contributing
Pull requests are welcome, update tests as appropriate.

## License
[MIT](https://github.com/SpinTools/seratoparser/LICENSE)
