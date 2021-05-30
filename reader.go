package seratoparser

import (
        "bufio"
        "github.com/romana/rlog"
        "os"
        "path/filepath"
        "runtime"
        "strings"
)

// Parser holds the filepath of all databases
type Parser struct {
        FilePath string
}

// SeratoParser is the Parser for this Serato Database module. Exported for future option of override
var SeratoParser Parser

// New creates a new object with the provided Serato Database Path
func New(seratoPath string) Parser {
        SeratoParser = Parser{}
        SeratoParser.FilePath = strings.TrimSuffix(seratoPath, "/")
        return SeratoParser
}

func fileHeader(fileBuffer *bufio.Reader, key string, value string) bool {
        // get first key
        firstKey, _ := parseCString(fileBuffer)
        if firstKey == "vrsn" {
                // Skip over next \0
                parseCString(fileBuffer)
        }

        // match vrsn
        key16 := makeUtf16(key)
        if !matchUtf16(fileBuffer, key16) {
                rlog.Error("fileHeader: vrsn value mismatch |", key16)
                return false
        }

        // match version type
        value16 := makeUtf16(value)
        if !matchUtf16(fileBuffer, value16) {
                rlog.Error("fileHeader: vrsn type mismatch |", value16)
                return false
        }

        return true
}

func fileCrateColumns(fileBuffer *bufio.Reader) {
        /*
           Field #    TAG     Description
           ==============================

           Header
           01        vrsn     4 byte string/tag denoting version start
           02        <DATA>   68 Byte null padded string representing the DB version.
           Decodes to 81.0/Serato ScratchLive Crate
           Column Sorting data
           03        osrt     4 byte string/tag denoting sorting config start
           04        <DATA>   4 bytes / 32bit int
           05        tvcn     4 byte string/tag denoting column name
           06        <DATA>   32bit int + Variable length string
           07        brev     4 byte string/tag, ???
           08        <DATA>   5 bytes 0x 00 00 00 01 00

           Column Details - repeated for all columns
           09        ovct     4 byte string/tag, ???
           10        <DATA>   4 bytes, 32bit int
           11        tvcn     byte string/tag denoting column name
           12        <DATA>   32bit int + Variable length string
           13        tvcw     4 byte string/tag, column width?
           14        <DATA>   6 bytes of ???

           Song/Track details - repeated for each track
           XX        otrk     4 byte string/tag stores track length
           XX        ptrk     32bit int + variable length track name
        */

        for {
                nextTag, _ := parseFilePeek(fileBuffer, 4)
                if string(nextTag) == "otrk" || string(nextTag) == "oent" {
                        break
                }
                // do stuff here for crate columns
                _, eof := parseFileLen(fileBuffer, 1)
                if eof {
                        break
                }
        }
}

func volumeName(filePath string) string {
        volume := filepath.VolumeName(filePath)

        if volume == "" && runtime.GOOS == "darwin" {
                if strings.HasPrefix(filePath, "/Volumes/") {
                        splitFilePath := strings.Split(filePath, "/")
                        if len(splitFilePath) >= 2 {
                                return "/Volumes/" + splitFilePath[2]
                        }
                }
        }

        return volume
}

func readMediaEntities(fileName string) []MediaEntity {
        var mediaEntities []MediaEntity
        seratoFile, err := filepath.Abs(fileName)
        if err != nil || seratoFile == "" {
                rlog.Critical(err)
                return mediaEntities
        }

        // Only read files that exist, and then report errors if we cant read it
        _, err = os.Stat(seratoFile)
        if err != nil {
                rlog.Critical(err)
                return mediaEntities
        }
        if os.IsNotExist(err) {
                rlog.Critical(err)
                return mediaEntities
        }

        ioFile, err := os.Open(seratoFile)
        if err != nil {
                rlog.Critical(err)
                return mediaEntities
        }
        defer ioFile.Close()

        seratoVolume = volumeName(seratoFile)

        fileBuffer := bufio.NewReader(ioFile)
        fileExt := filepath.Ext(seratoFile)
        fileType := filepath.Base(seratoFile)
        if fileExt == ".crate" {
                if !fileHeader(fileBuffer, "81.0", "/Serato ScratchLive Crate") {
                        rlog.Critical("ReadFile: Unable to parse crate |", seratoFile)
                }
        } else if fileType == "database V2" {
                if !fileHeader(fileBuffer, "@2.0", "/Serato Scratch LIVE Database") {
                        rlog.Critical("ReadFile: Unable to parse database v2 |", seratoFile)
                }
        } else if fileType == "history.database" || fileExt == ".session" {
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
                                if name == "otrk" {
                                        mediaEntity := MediaEntity{}
                                        parseOtrk(&dataBuffer, &mediaEntity)
                                        mediaEntities = append(mediaEntities, mediaEntity)
                                }
                        }
                } else {
                        parseFileLen(fileBuffer, 1)
                }
        }

        return mediaEntities
}
