package seratoparser

import (
        "bufio"
        "github.com/romana/rlog"
        "io/fs"
        "io/ioutil"
        "os"
        "path/filepath"
        "sort"
        "strings"
)

// HistoryPath is the path inside the _Serato_ folder that contains History data
var HistoryPath = "/History"

// SessionPath is the path inside the _Serato_/History folder that contains all the played Sessions.
var SessionPath = HistoryPath + "/Sessions"

// GetHistorySessions returns a list of all Serato History session files in the users Serato directory.
func (p Parser) GetHistorySessions() []fs.FileInfo {
        var sessionFiles []fs.FileInfo
        var err error

        //historySessionDir := currentUser.HomeDir + "/Music/_Serato_/History/Sessions"
        historySessionDir := filepath.FromSlash(p.FilePath + SessionPath)
        sessionFiles, err = ioutil.ReadDir(historySessionDir)

        // Remove .DS_STORE files.
        for i := 0; i < len(sessionFiles); i++ {
                if !strings.HasSuffix(sessionFiles[i].Name(), ".session") {
                        sessionFiles = append(sessionFiles[:i], sessionFiles[i+1:]...)
                        i--
                }
        }

        if err != nil || len(sessionFiles) == 0 {
                rlog.Critical(err)
                return sessionFiles
        }

        sort.Slice(sessionFiles, func(i, j int) bool {
                return sessionFiles[i].ModTime().Unix() > sessionFiles[j].ModTime().Unix()
        })

        return sessionFiles
}

// ReadHistorySession returns all track entities within the provided filepath.
func (p Parser) ReadHistorySession(fileName string) []HistoryEntity {
        var historyEntities []HistoryEntity

        historySessionFilepath := filepath.FromSlash(p.FilePath + SessionPath + "/" + fileName)
        seratoFile, err := filepath.Abs(historySessionFilepath)
        if err != nil || seratoFile == "" {
                rlog.Critical(err)
                return historyEntities
        }

        // Only read files that exist, and then report errors if we cant read it
        _, err = os.Stat(seratoFile)
        if err != nil {
                rlog.Critical(err)
                return historyEntities
        }
        if os.IsNotExist(err) {
                rlog.Critical(err)
                return historyEntities
        }

        ioFile, err := os.Open(seratoFile)
        if err != nil {
                rlog.Critical(err)
                return historyEntities
        }
        defer ioFile.Close()

        seratoVolume = volumeName(seratoFile)

        fileBuffer := bufio.NewReader(ioFile)
        fileExt := filepath.Ext(seratoFile)
        if fileExt == ".session" {
                if !fileHeader(fileBuffer, "<1.0", "/Serato Scratch LIVE Review") {
                        rlog.Critical("ReadFile: Unable to parse history |", seratoFile)
                }
        }

        defer func() {
                if r := recover(); r != nil {
                        rlog.Warn("Recovered in fileReader", r)
                }
        }()

        for {
                nextTag1, eof := parseFilePeek(fileBuffer, 1)
                nextTag4, _ := parseFilePeek(fileBuffer, 4)
                if eof || string(nextTag1) == "" {
                        break
                } else if string(nextTag4) == "osrt" {
                        fileCrateColumns(fileBuffer)
                } else if string(nextTag4) == "otrk" || string(nextTag4) == "oent" { //|| string(nextTag4) == "oses" {
                        for {
                                name, data, eof := parseField(fileBuffer)

                                // break if we don't need to be here
                                //  TODO: What is oses?
                                if eof || (name != "otrk" && name != "oent") { //&& name != "oses") {
                                        break
                                }

                                dataBuffer := *bufio.NewReader(strings.NewReader(data))
                                if name == "oent" {
                                        for {
                                                dataName, dataValue, eof := parseField(&dataBuffer)
                                                if eof || dataName != "adat" {
                                                        break
                                                }

                                                historyEntity := HistoryEntity{}
                                                parseAdat(dataValue, &historyEntity)
                                                historyEntities = append(historyEntities, historyEntity)
                                        }
                                } else if name == "oses" {
                                        for {
                                                dataName, dataValue, eof := parseField(&dataBuffer)
                                                if eof || dataName != "adat" {
                                                        break
                                                }

                                                rlog.Println(dataName)
                                                rlog.Println(dataValue)
                                                rlog.Println("--------")

                                                historyEntity := HistoryEntity{}
                                                parseAdat(dataValue, &historyEntity)
                                        }
                                }
                        }
                } else {
                        parseFileLen(fileBuffer, 1)
                }
        }

        return historyEntities
}
