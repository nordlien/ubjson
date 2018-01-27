package ubjson

import (
	"reflect"
	"sort"
)

func init() {
	// Force deterministic ordering to reflective map key traversal.
	mapKeys = func(m reflect.Value) []reflect.Value {
		vs := m.MapKeys()
		sort.Slice(vs, func(i, j int) bool {
			return vs[i].String() < vs[j].String()
		})
		return vs
	}
}

// A Values struct holds one of each primitive type with a corresponding
// Encode/Decode UBJSON value method.
type Values struct {
	Int            int
	UInt8          uint8
	Int8           int8
	Int16          int16
	Int32          int32
	Int64          int64
	Float32        float32
	Float64        float64
	Bool           bool
	String         string
	Char           Char
	HighPrecNumber HighPrecNumber

	IntPtr            *int
	UInt8Ptr          *uint8
	Int8Ptr           *int8
	Int16Ptr          *int16
	Int32Ptr          *int32
	Int64Ptr          *int64
	Float32Ptr        *float32
	Float64Ptr        *float64
	BoolPtr           *bool
	StringPtr         *string
	CharPtr           *Char
	HighPrecNumberPtr *HighPrecNumber
}

type testCase struct {
	value  interface{}
	binary []byte
	block  string
}

var cases = map[string]testCase{
	"Int=0":   {int(0), []byte{'U', 0x00}, "[U][0]"},
	"Int=255": {int(255), []byte{'U', 0xFF}, "[U][255]"},
	"UInt=80": {uint8(0), []byte{'U', 0x00}, "[U][0]"},

	"Int=-128": {int(-128), []byte{'i', 0x80}, "[i][-128]"},
	"Int8=127": {int8(127), []byte{'i', 0x7F}, "[i][127]"},

	"Int=256":     {int(256), []byte{'I', 0x01, 0x00}, "[I][256]"},
	"Int16=32767": {int16(32767), []byte{'I', 0x7F, 0xFF}, "[I][32767]"},

	"Int=32768":        {int(32768), []byte{'l', 0x00, 0x00, 0x80, 0x00}, "[l][32768]"},
	"Int32=2147483647": {int32(2147483647), []byte{'l', 0x7F, 0xFF, 0xFF, 0xFF}, "[l][2147483647]"},

	"Int=2147483648": {int(2147483648), []byte{0: 'L', 0x00, 0x00, 0x00, 0x00, 0x80, 0x00, 0x00, 0x00},
		"[L][2147483648]"},
	"Int62=231457428363448": {int64(231457428363448), []byte{0: 'L', 0x00, 0x00, 0xD2, 0x82, 0x61, 0xCC, 0x58, 0xB8},
		"[L][231457428363448]"},

	"Float32=2147483648": {float32(2147483648), []byte{'d', 0x4F, 0x00, 0x00, 0x00},
		"[d][2.1474836e+09]"},
	"Float64=3.402823e38": {float64(3.402823e38), []byte{'D', 0x47, 0xEF, 0xFF, 0xFF, 0x96, 0x6A, 0xD9, 0x24},
		"[D][3.402823e+38]"},

	"HighPrecNumber=3.402823e38": {HighPrecNumber("3.402823e38"), append([]byte{'H', 0x55, 0x0B}, "3.402823e38"...),
		"[H][U][11][3.402823e38]"},

	"C=a": {Char('a'), []byte{'C', 0x61}, "[C][a]"},

	"string=string": {"string", append([]byte{'S', 0x55, 0x06}, "string"...), "[S][U][6][string]"},

	"Array-empty": {[0]int{}, []byte{0x5b, 0x23, 0x55, 0x0}, "[[][#][U][0]"},

	"Slice-empty": {[]int{}, []byte{0x5b, 0x23, 0x55, 0x0}, "[[][#][U][0]"},

	"Array-UInt8=byte-array": {[2]byte{0x4C, 0x7F}, []byte{'[', '$', 'U', '#', 'U', 0x02, 0x4C, 0x7F},
		"[[][$][U][#][U][2]\n\t[76]\n\t[127]"},

	"Array-UInt8=byte-slice": {[]byte{0x4C, 0x7F}, []byte{'[', '$', 'U', '#', 'U', 0x02, 0x4C, 0x7F},
		"[[][$][U][#][U][2]\n\t[76]\n\t[127]"},

	"Object=map": {
		map[string]interface{}{"a": uint8(5)},
		[]byte{'{', '#', 'U', 0x01,
			'U', 0x01, 'a', 'U', 0x05},
		"[{][#][U][1]\n\t[U][1][a][U][5]",
	},
	"Object-Int8=map": {
		map[string]int8{"a": 5},
		[]byte{'{', '$', 'i', '#', 'U', 0x01,
			'U', 0x01, 'a', 0x05},
		"[{][$][i][#][U][1]\n\t[U][1][a][5]",
	},
	"Object-Int8=struct": {
		struct {
			A int8
			B int8
		}{5, 8},
		[]byte{'{',
			'U', 0x01, 'A', 'i', 0x05,
			'U', 0x01, 'B', 'i', 0x08,
			'}'},
		"[{]\n\t[U][1][A][i][5]\n\t[U][1][B][i][8]\n[}]",
	},
	"Object-Int8=struct-tagged": {
		struct {
			A int8 `ubjson:"a"`
			a int  // Ignored - Exercises field index logic.
			B int8 `json:"wrong" ubjson:"b"`
		}{5, 0, 8},
		[]byte{'{',
			'U', 0x01, 'a', 'i', 0x05,
			'U', 0x01, 'b', 'i', 0x08,
			'}'},
		"[{]\n\t[U][1][a][i][5]\n\t[U][1][b][i][8]\n[}]",
	},

	"Object=complex-struct": {complexStruct, complexStructBinary, complexStructBlock},
	"Object=complex-map":    {complexMap, complexMapBinary, complexMapBlock},
}

type complexType struct {
	Location            string
	Email               string
	Type                string
	Total_private_repos int
	Blog                string
	Gravatar_id         string
	URL                 string
	Company             string
	Hireable            bool
	Bio                 string
	Public_gists        int
	Html_url            string
	Id                  int
	Followers           int
	Following           int
	Created_at          string
	Owned_private_repos int
	Collaborators       int
	Login               string
	Name                string
	Public_repos        int
	Private_gists       int
	Disk_usage          int
	Plan                Plan
	Avatar_url          string
}

type Plan struct {
	Name          string
	Space         int
	Collaborators int
	Private_repos int
}

var complexStruct = complexType{
	Location:            "San Francisco",
	Email:               "octocat@github.com",
	Type:                "User",
	Total_private_repos: 100,
	Blog:                "https://github.com/blog",
	Gravatar_id:         "somehexcode",
	URL:                 "https://api.github.com/users/octocat",
	Company:             "GitHub",
	Hireable:            false,
	Bio:                 "There once was...",
	Public_gists:        1,
	Html_url:            "https://github.com/octocat",
	Id:                  1,
	Followers:           20,
	Following:           0,
	Created_at:          "2008-01-14T04:33:35Z",
	Owned_private_repos: 100,
	Collaborators:       8,
	Login:               "octocat",
	Name:                "monalisa octocat",
	Public_repos:        2,
	Private_gists:       81,
	Disk_usage:          10000,
	Plan: Plan{
		Name:          "Medium",
		Space:         400,
		Collaborators: 10,
		Private_repos: 20,
	},
	Avatar_url: "https://github.com/images/error/octocat_happy.gif",
}

var complexMap = map[string]interface{}{
	"Location":            "San Francisco",
	"Email":               "octocat@github.com",
	"Type":                "User",
	"Total_private_repos": uint8(100),
	"Blog":                "https://github.com/blog",
	"Gravatar_id":         "somehexcode",
	"URL":                 "https://api.github.com/users/octocat",
	"Company":             "GitHub",
	"Hireable":            false,
	"Bio":                 "There once was...",
	"Public_gists":        uint8(1),
	"Html_url":            "https://github.com/octocat",
	"Id":                  uint8(1),
	"Followers":           uint8(20),
	"Following":           uint8(0),
	"Created_at":          "2008-01-14T04:33:35Z",
	"Owned_private_repos": uint8(100),
	"Collaborators":       uint8(8),
	"Login":               "octocat",
	"Name":                "monalisa octocat",
	"Public_repos":        uint8(2),
	"Private_gists":       uint8(81),
	"Disk_usage":          int16(10000),
	"Plan": map[string]interface{}{
		"Name":          "Medium",
		"Space":         int16(400),
		"Collaborators": uint8(10),
		"Private_repos": uint8(20),
	},
	"Avatar_url": "https://github.com/images/error/octocat_happy.gif",
}

var complexMapBlock = `[{][#][U][25]
	[U][10][Avatar_url][S][U][49][https://github.com/images/error/octocat_happy.gif]
	[U][3][Bio][S][U][17][There once was...]
	[U][4][Blog][S][U][23][https://github.com/blog]
	[U][13][Collaborators][U][8]
	[U][7][Company][S][U][6][GitHub]
	[U][10][Created_at][S][U][20][2008-01-14T04:33:35Z]
	[U][10][Disk_usage][I][10000]
	[U][5][Email][S][U][18][octocat@github.com]
	[U][9][Followers][U][20]
	[U][9][Following][U][0]
	[U][11][Gravatar_id][S][U][11][somehexcode]
	[U][8][Hireable][F]
	[U][8][Html_url][S][U][26][https://github.com/octocat]
	[U][2][Id][U][1]
	[U][8][Location][S][U][13][San Francisco]
	[U][5][Login][S][U][7][octocat]
	[U][4][Name][S][U][16][monalisa octocat]
	[U][19][Owned_private_repos][U][100]
	[U][4][Plan][{][#][U][4]
		[U][13][Collaborators][U][10]
		[U][4][Name][S][U][6][Medium]
		[U][13][Private_repos][U][20]
		[U][5][Space][I][400]
	[U][13][Private_gists][U][81]
	[U][12][Public_gists][U][1]
	[U][12][Public_repos][U][2]
	[U][19][Total_private_repos][U][100]
	[U][4][Type][S][U][4][User]
	[U][3][URL][S][U][36][https://api.github.com/users/octocat]`

var complexStructBlock = `[{]
	[U][8][Location][S][U][13][San Francisco]
	[U][5][Email][S][U][18][octocat@github.com]
	[U][4][Type][S][U][4][User]
	[U][19][Total_private_repos][U][100]
	[U][4][Blog][S][U][23][https://github.com/blog]
	[U][11][Gravatar_id][S][U][11][somehexcode]
	[U][3][URL][S][U][36][https://api.github.com/users/octocat]
	[U][7][Company][S][U][6][GitHub]
	[U][8][Hireable][F]
	[U][3][Bio][S][U][17][There once was...]
	[U][12][Public_gists][U][1]
	[U][8][Html_url][S][U][26][https://github.com/octocat]
	[U][2][Id][U][1]
	[U][9][Followers][U][20]
	[U][9][Following][U][0]
	[U][10][Created_at][S][U][20][2008-01-14T04:33:35Z]
	[U][19][Owned_private_repos][U][100]
	[U][13][Collaborators][U][8]
	[U][5][Login][S][U][7][octocat]
	[U][4][Name][S][U][16][monalisa octocat]
	[U][12][Public_repos][U][2]
	[U][13][Private_gists][U][81]
	[U][10][Disk_usage][I][10000]
	[U][4][Plan][{]
		[U][4][Name][S][U][6][Medium]
		[U][5][Space][I][400]
		[U][13][Collaborators][U][10]
		[U][13][Private_repos][U][20]
	[}]
	[U][10][Avatar_url][S][U][49][https://github.com/images/error/octocat_happy.gif]
[}]`

//TODO faked these from output - verify with another tool
var complexMapBinary = []byte{
	0x7b, 0x23, 0x55, 0x19, 0x55, 0xa, 0x41, 0x76, 0x61, 0x74, 0x61, 0x72, 0x5f, 0x75, 0x72, 0x6c, 0x53, 0x55, 0x31, 0x68, 0x74, 0x74, 0x70, 0x73, 0x3a, 0x2f, 0x2f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x73, 0x2f, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x2f, 0x6f, 0x63, 0x74, 0x6f, 0x63, 0x61, 0x74, 0x5f, 0x68, 0x61, 0x70, 0x70, 0x79, 0x2e, 0x67, 0x69, 0x66, 0x55, 0x3, 0x42, 0x69, 0x6f, 0x53, 0x55, 0x11, 0x54, 0x68, 0x65, 0x72, 0x65, 0x20, 0x6f, 0x6e, 0x63, 0x65, 0x20, 0x77, 0x61, 0x73, 0x2e, 0x2e, 0x2e, 0x55, 0x4, 0x42, 0x6c, 0x6f, 0x67, 0x53, 0x55, 0x17, 0x68, 0x74, 0x74, 0x70, 0x73, 0x3a, 0x2f, 0x2f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x62, 0x6c, 0x6f, 0x67, 0x55, 0xd, 0x43, 0x6f, 0x6c, 0x6c, 0x61, 0x62, 0x6f, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x73, 0x55, 0x8, 0x55, 0x7, 0x43, 0x6f, 0x6d, 0x70, 0x61, 0x6e, 0x79, 0x53, 0x55, 0x6, 0x47, 0x69, 0x74, 0x48, 0x75, 0x62, 0x55, 0xa, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x53, 0x55, 0x14, 0x32, 0x30, 0x30, 0x38, 0x2d, 0x30, 0x31, 0x2d, 0x31, 0x34, 0x54, 0x30, 0x34, 0x3a, 0x33, 0x33, 0x3a, 0x33, 0x35, 0x5a, 0x55, 0xa, 0x44, 0x69, 0x73, 0x6b, 0x5f, 0x75, 0x73, 0x61, 0x67, 0x65, 0x49, 0x27, 0x10, 0x55, 0x5, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x53, 0x55, 0x12, 0x6f, 0x63, 0x74, 0x6f, 0x63, 0x61, 0x74, 0x40, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x55, 0x9, 0x46, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x72, 0x73, 0x55, 0x14, 0x55, 0x9, 0x46, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x69, 0x6e, 0x67, 0x55, 0x0, 0x55, 0xb, 0x47, 0x72, 0x61, 0x76, 0x61, 0x74, 0x61, 0x72, 0x5f, 0x69, 0x64, 0x53, 0x55, 0xb, 0x73, 0x6f, 0x6d, 0x65, 0x68, 0x65, 0x78, 0x63, 0x6f, 0x64, 0x65, 0x55, 0x8, 0x48, 0x69, 0x72, 0x65, 0x61, 0x62, 0x6c, 0x65, 0x46, 0x55, 0x8, 0x48, 0x74, 0x6d, 0x6c, 0x5f, 0x75, 0x72, 0x6c, 0x53, 0x55, 0x1a, 0x68, 0x74, 0x74, 0x70, 0x73, 0x3a, 0x2f, 0x2f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6f, 0x63, 0x74, 0x6f, 0x63, 0x61, 0x74, 0x55, 0x2, 0x49, 0x64, 0x55, 0x1, 0x55, 0x8, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x55, 0xd, 0x53, 0x61, 0x6e, 0x20, 0x46, 0x72, 0x61, 0x6e, 0x63, 0x69, 0x73, 0x63, 0x6f, 0x55, 0x5, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x53, 0x55, 0x7, 0x6f, 0x63, 0x74, 0x6f, 0x63, 0x61, 0x74, 0x55, 0x4, 0x4e, 0x61, 0x6d, 0x65, 0x53, 0x55, 0x10, 0x6d, 0x6f, 0x6e, 0x61, 0x6c, 0x69, 0x73, 0x61, 0x20, 0x6f, 0x63, 0x74, 0x6f, 0x63, 0x61, 0x74, 0x55, 0x13, 0x4f, 0x77, 0x6e, 0x65, 0x64, 0x5f, 0x70, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x5f, 0x72, 0x65, 0x70, 0x6f, 0x73, 0x55, 0x64, 0x55, 0x4, 0x50, 0x6c, 0x61, 0x6e, 0x7b, 0x23, 0x55, 0x4, 0x55, 0xd, 0x43, 0x6f, 0x6c, 0x6c, 0x61, 0x62, 0x6f, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x73, 0x55, 0xa, 0x55, 0x4, 0x4e, 0x61, 0x6d, 0x65, 0x53, 0x55, 0x6, 0x4d, 0x65, 0x64, 0x69, 0x75, 0x6d, 0x55, 0xd, 0x50, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x5f, 0x72, 0x65, 0x70, 0x6f, 0x73, 0x55, 0x14, 0x55, 0x5, 0x53, 0x70, 0x61, 0x63, 0x65, 0x49, 0x1, 0x90, 0x55, 0xd, 0x50, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x5f, 0x67, 0x69, 0x73, 0x74, 0x73, 0x55, 0x51, 0x55, 0xc, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x5f, 0x67, 0x69, 0x73, 0x74, 0x73, 0x55, 0x1, 0x55, 0xc, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x5f, 0x72, 0x65, 0x70, 0x6f, 0x73, 0x55, 0x2, 0x55, 0x13, 0x54, 0x6f, 0x74, 0x61, 0x6c, 0x5f, 0x70, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x5f, 0x72, 0x65, 0x70, 0x6f, 0x73, 0x55, 0x64, 0x55, 0x4, 0x54, 0x79, 0x70, 0x65, 0x53, 0x55, 0x4, 0x55, 0x73, 0x65, 0x72, 0x55, 0x3, 0x55, 0x52, 0x4c, 0x53, 0x55, 0x24, 0x68, 0x74, 0x74, 0x70, 0x73, 0x3a, 0x2f, 0x2f, 0x61, 0x70, 0x69, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2f, 0x6f, 0x63, 0x74, 0x6f, 0x63, 0x61, 0x74,
}
var complexStructBinary = []byte{
	0x7b, 0x55, 0x8, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x55, 0xd, 0x53, 0x61, 0x6e, 0x20, 0x46, 0x72, 0x61, 0x6e, 0x63, 0x69, 0x73, 0x63, 0x6f, 0x55, 0x5, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x53, 0x55, 0x12, 0x6f, 0x63, 0x74, 0x6f, 0x63, 0x61, 0x74, 0x40, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x55, 0x4, 0x54, 0x79, 0x70, 0x65, 0x53, 0x55, 0x4, 0x55, 0x73, 0x65, 0x72, 0x55, 0x13, 0x54, 0x6f, 0x74, 0x61, 0x6c, 0x5f, 0x70, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x5f, 0x72, 0x65, 0x70, 0x6f, 0x73, 0x55, 0x64, 0x55, 0x4, 0x42, 0x6c, 0x6f, 0x67, 0x53, 0x55, 0x17, 0x68, 0x74, 0x74, 0x70, 0x73, 0x3a, 0x2f, 0x2f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x62, 0x6c, 0x6f, 0x67, 0x55, 0xb, 0x47, 0x72, 0x61, 0x76, 0x61, 0x74, 0x61, 0x72, 0x5f, 0x69, 0x64, 0x53, 0x55, 0xb, 0x73, 0x6f, 0x6d, 0x65, 0x68, 0x65, 0x78, 0x63, 0x6f, 0x64, 0x65, 0x55, 0x3, 0x55, 0x52, 0x4c, 0x53, 0x55, 0x24, 0x68, 0x74, 0x74, 0x70, 0x73, 0x3a, 0x2f, 0x2f, 0x61, 0x70, 0x69, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2f, 0x6f, 0x63, 0x74, 0x6f, 0x63, 0x61, 0x74, 0x55, 0x7, 0x43, 0x6f, 0x6d, 0x70, 0x61, 0x6e, 0x79, 0x53, 0x55, 0x6, 0x47, 0x69, 0x74, 0x48, 0x75, 0x62, 0x55, 0x8, 0x48, 0x69, 0x72, 0x65, 0x61, 0x62, 0x6c, 0x65, 0x46, 0x55, 0x3, 0x42, 0x69, 0x6f, 0x53, 0x55, 0x11, 0x54, 0x68, 0x65, 0x72, 0x65, 0x20, 0x6f, 0x6e, 0x63, 0x65, 0x20, 0x77, 0x61, 0x73, 0x2e, 0x2e, 0x2e, 0x55, 0xc, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x5f, 0x67, 0x69, 0x73, 0x74, 0x73, 0x55, 0x1, 0x55, 0x8, 0x48, 0x74, 0x6d, 0x6c, 0x5f, 0x75, 0x72, 0x6c, 0x53, 0x55, 0x1a, 0x68, 0x74, 0x74, 0x70, 0x73, 0x3a, 0x2f, 0x2f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6f, 0x63, 0x74, 0x6f, 0x63, 0x61, 0x74, 0x55, 0x2, 0x49, 0x64, 0x55, 0x1, 0x55, 0x9, 0x46, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x72, 0x73, 0x55, 0x14, 0x55, 0x9, 0x46, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x69, 0x6e, 0x67, 0x55, 0x0, 0x55, 0xa, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x53, 0x55, 0x14, 0x32, 0x30, 0x30, 0x38, 0x2d, 0x30, 0x31, 0x2d, 0x31, 0x34, 0x54, 0x30, 0x34, 0x3a, 0x33, 0x33, 0x3a, 0x33, 0x35, 0x5a, 0x55, 0x13, 0x4f, 0x77, 0x6e, 0x65, 0x64, 0x5f, 0x70, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x5f, 0x72, 0x65, 0x70, 0x6f, 0x73, 0x55, 0x64, 0x55, 0xd, 0x43, 0x6f, 0x6c, 0x6c, 0x61, 0x62, 0x6f, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x73, 0x55, 0x8, 0x55, 0x5, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x53, 0x55, 0x7, 0x6f, 0x63, 0x74, 0x6f, 0x63, 0x61, 0x74, 0x55, 0x4, 0x4e, 0x61, 0x6d, 0x65, 0x53, 0x55, 0x10, 0x6d, 0x6f, 0x6e, 0x61, 0x6c, 0x69, 0x73, 0x61, 0x20, 0x6f, 0x63, 0x74, 0x6f, 0x63, 0x61, 0x74, 0x55, 0xc, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x5f, 0x72, 0x65, 0x70, 0x6f, 0x73, 0x55, 0x2, 0x55, 0xd, 0x50, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x5f, 0x67, 0x69, 0x73, 0x74, 0x73, 0x55, 0x51, 0x55, 0xa, 0x44, 0x69, 0x73, 0x6b, 0x5f, 0x75, 0x73, 0x61, 0x67, 0x65, 0x49, 0x27, 0x10, 0x55, 0x4, 0x50, 0x6c, 0x61, 0x6e, 0x7b, 0x55, 0x4, 0x4e, 0x61, 0x6d, 0x65, 0x53, 0x55, 0x6, 0x4d, 0x65, 0x64, 0x69, 0x75, 0x6d, 0x55, 0x5, 0x53, 0x70, 0x61, 0x63, 0x65, 0x49, 0x1, 0x90, 0x55, 0xd, 0x43, 0x6f, 0x6c, 0x6c, 0x61, 0x62, 0x6f, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x73, 0x55, 0xa, 0x55, 0xd, 0x50, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x5f, 0x72, 0x65, 0x70, 0x6f, 0x73, 0x55, 0x14, 0x7d, 0x55, 0xa, 0x41, 0x76, 0x61, 0x74, 0x61, 0x72, 0x5f, 0x75, 0x72, 0x6c, 0x53, 0x55, 0x31, 0x68, 0x74, 0x74, 0x70, 0x73, 0x3a, 0x2f, 0x2f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x73, 0x2f, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x2f, 0x6f, 0x63, 0x74, 0x6f, 0x63, 0x61, 0x74, 0x5f, 0x68, 0x61, 0x70, 0x70, 0x79, 0x2e, 0x67, 0x69, 0x66, 0x7d,
}
