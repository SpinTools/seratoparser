package serato_parser

var seratoVolume string

type Crate struct {
	Id          	int64
	ParentId    	int64
	Name        	string
	Path        	string
	Filepath    	string
	Folders     	[]string
	Genre       	string
	Expanded		int
	Position		int
	Casing			int
	MakeSubCrates	bool
}

type CrateMediaEntities struct {
	Id		int64
	Tracks	[]MediaEntity
}

type MediaEntity struct {
	// META
	DVOL    string  // volume

	// UTFSTR
	PTRK    string  // filetrack
	PFIL    string  // filebase

	// INT1
	BMIS    bool  // missing
	BCRT    bool  // corrupt

	// INT4
	UADD    int  // timeadded

	// BYTE SLICE
	ULBL    []byte // color - track colour
}

type MediaDbEntity struct {
	Vol         string
	Filepath    string
	Missing     bool
	Corrupt     bool
	Added       int
	Color       string
}

type SessionHistoryEntities struct {
	Id		int64
	Tracks	[][]HistoryEntity
}

type HistoryEntity struct {
	RROW    int     // rrow
	RDIR    string  // rfullpath
	TTMS    int     // rstarttime
	TTME    int     // rendtime
	TDCK    int     // rdeck
	RDTE    string  // rdate*
	RSRT    int     // rstart*
	REND    int     // rend*
	TPTM    int     // rplaytime
	RSES    int     // rsessionId
	RPLY    int     // rplayed = 1
	RADD    int     // radded
	RUPD    int     // rupdatedAt
	RSWR    string  // rsoftware*
	RSWB    int     // rsoftwareBuild*
	RDEV    string  // rdevice
}

// TODO: Identify all fields of an ADAT object
var SeratoAdatMap = map[int]string{
	1   :   "RROW",  // rrow
	2   :   "RDIR",  // rfullpath
	28  :   "TTMS",  // rstarttime
	29  :   "TTME",  // rendtime
	31  :   "TDCK",  // rdeck

	41  :   "RDTE",  // rdate
	43  :   "RSRT",  // rstart
	44  :   "REND",  // rend

	45  :   "TPTM",  // rplaytime
	48  :   "RSES",  // rsessionId
	50  :   "RPLY",  // rplayed
	52  :   "RADD",  // radded
	53  :   "RUPD",  // rupdatedAt
	54  :   "RUNK",  // rr54unknownTimestamp

	57  :   "RSWR",  // rsoftware
	58  :   "RSWB",  // rsoftwareBuild

	63  :   "RDEV",  // rdevice
}