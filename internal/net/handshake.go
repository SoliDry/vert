package net

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	vauth "github.com/SoliDry/vert/internal/auth"
)

const (
	ProtocolVersion = "\x01"
	ServerVersion   = "0.1.1"
	MaxPacketSize   = 4096
)

/* the table of a character set and collation with IDs
+--------------------+---------------------+-----+
| CHARACTER_SET_NAME | COLLATION_NAME      | ID  |
+--------------------+---------------------+-----+
| big5               | big5_chinese_ci     |   1 |
| dec8               | dec8_swedish_ci     |   3 |
| cp850              | cp850_general_ci    |   4 |
| hp8                | hp8_english_ci      |   6 |
| koi8r              | koi8r_general_ci    |   7 |
| latin1             | latin1_swedish_ci   |   8 |
| latin2             | latin2_general_ci   |   9 |
| swe7               | swe7_swedish_ci     |  10 |
| ascii              | ascii_general_ci    |  11 |
| ujis               | ujis_japanese_ci    |  12 |
| sjis               | sjis_japanese_ci    |  13 |
| hebrew             | hebrew_general_ci   |  16 |
| tis620             | tis620_thai_ci      |  18 |
| euckr              | euckr_korean_ci     |  19 |
| koi8u              | koi8u_general_ci    |  22 |
| gb2312             | gb2312_chinese_ci   |  24 |
| greek              | greek_general_ci    |  25 |
| cp1250             | cp1250_general_ci   |  26 |
| gbk                | gbk_chinese_ci      |  28 |
| latin5             | latin5_turkish_ci   |  30 |
| armscii8           | armscii8_general_ci |  32 |
| utf8               | utf8_general_ci     |  33 |
| ucs2               | ucs2_general_ci     |  35 |
| cp866              | cp866_general_ci    |  36 |
| keybcs2            | keybcs2_general_ci  |  37 |
| macce              | macce_general_ci    |  38 |
| macroman           | macroman_general_ci |  39 |
| cp852              | cp852_general_ci    |  40 |
| latin7             | latin7_general_ci   |  41 |
| cp1251             | cp1251_general_ci   |  51 |
| utf16              | utf16_general_ci    |  54 |
| utf16le            | utf16le_general_ci  |  56 |
| cp1256             | cp1256_general_ci   |  57 |
| cp1257             | cp1257_general_ci   |  59 |
| utf32              | utf32_general_ci    |  60 |
| binary             | binary              |  63 |
| geostd8            | geostd8_general_ci  |  92 |
| cp932              | cp932_japanese_ci   |  95 |
| eucjpms            | eucjpms_japanese_ci |  97 |
| gb18030            | gb18030_chinese_ci  | 248 |
| utf8mb4            | utf8mb4_0900_ai_ci  | 255 |
+--------------------+---------------------+-----+
*/

// capability flags
const (
	clientConnectWithDb uint64 = 0x00000001
	//clientCompress = 0x00000002
)

//GetHandshakeSeq handshaking with connected client
func GetHandshakeSeq(conn Connection, dataIn []byte) ([]byte, error) {
	dataInSplit := strings.Split(string(dataIn), "\n")
	if len(dataInSplit) < 4 {
		return nil, fmt.Errorf("error getting handshake data: less than %d required elements", len(dataInSplit))
	}

	// The capability flags are used by the client and server to indicate which features they support and want to use.
	cf, err := strconv.ParseUint(dataInSplit[0], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing capability flags: %w", err)
	}

	capabilityFlags := cf
	maxPacketSize := dataInSplit[1]
	characterSet := dataInSplit[2]
	userName := dataInSplit[3]
	password := dataInSplit[4]

	if len(characterSet) > 0 && characterSet != conn.CharSet {
		conn.CharSet = characterSet
	}

	mpSize, err := strconv.ParseUint(maxPacketSize, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("error parsing maxPacketSize: %w", err)
	}

	conn.MaxPacketSize = MaxPacketSize
	if len(maxPacketSize) > 0 && conn.MaxPacketSize != uint32(mpSize) {
		conn.MaxPacketSize = uint32(mpSize)
	}

	if (capabilityFlags & clientConnectWithDb) != 0 {

	}

	dataOut := ProtocolVersion + "\n" + ServerVersion + "\n" + strconv.FormatInt(conn.ID, 10) + "\n" +
		conn.CharSet + "\n" + strconv.FormatUint(uint64(conn.MaxPacketSize), 10)

	user := vauth.NewUser(userName, password, false)
	ok, err := user.Authenticate()
	if err != nil {
		return nil, fmt.Errorf("login error: %w", err)
	}

	if !ok {
		return nil, errors.New("user doesn't exist")
	}

	return []byte(dataOut), nil
}
