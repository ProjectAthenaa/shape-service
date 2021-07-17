package generation

import (
	"encoding/binary"
	"github.com/A-Solutionss/superpack"
	"log"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type VersionPayloadHolder struct {
	Enabled      bool
	Hashseedbase int
	Hashseed1    int
	Hashseed2    int
	Currval      int
	Seedcount    int
	Keyholder    map[int]int64
	GenOrder     []int
	BaseFile     string
	HashString   string
	SeedString   string
	ResString    string
	Alphabet     string
	Dheader      string
	GlobalKeys   []int
	BaseKeys     []int
	Calckey      int
	Encrounds    int
	P4VAL        int
	P1VAL        []int
	P2VAL        []int
	P3VAL        []int
	P5VAL        int
	P11VAL       []int
	P13VAL       []int
	P14VAL       []int
	P15VAL       []int
	P16VAL       []int
	P17VAL       [][]int
	P18VAL       []int
	P19VAL       []int
	P20VAL       []int
	P21VAL       []int
	P23VAL       []int
	P24VAL       []int
	P25VAL       []int
	P30VAL       []int
	P31VAL       []int
	P32VAL       []int
	P35VAL       []int
	P36VAL       []int
	P40VAL       []int
	P42VAL       []int
	P44VAL       []int
}

type KeyEvent struct {
	EventType  int
	keyCode    int
	TargetId   string
	TargetName string
	Timestamp  float64
}
type MouseEvent struct {
	DocumentRelativeX int
	DocumentRelativeY int
	Timestamp         float64
}
type MouseTarget struct {
	Button            int
	DocumentRelativeX int
	DocumentRelativeY int
	EventType         int
	TargetId          string
	TargetName        string
	TargetRelativeX   int
	TargetRelativeY   int
	Timestamp         float64
}

func (h *VersionPayloadHolder) GenerateHeaders() map[string]string {
	rawarrholder := [][]int{}
	encarr := []int{0}
	for _, indx := range h.GenOrder {
		switch indx {
		case 0:
			apparr := h.Payload0(int(h.Keyholder[0]))
			rawarrholder = append(rawarrholder, apparr)
		case 1:
			apparr := h.Payload1(int(h.Keyholder[1]))
			rawarrholder = append(rawarrholder, apparr)
		case 2:
			apparr := h.Payload2(int(h.Keyholder[2]))
			rawarrholder = append(rawarrholder, apparr)
		case 3:
			apparr := h.Payload3(int(h.Keyholder[3]))
			rawarrholder = append(rawarrholder, apparr)
		case 4:
			apparr := h.Payload4(int(h.Keyholder[4]))
			rawarrholder = append(rawarrholder, apparr)
		case 5:
			apparr := h.Payload5(int(h.Keyholder[5]))
			rawarrholder = append(rawarrholder, apparr)
		case 6:
			apparr := h.Payload6(int(h.Keyholder[6]))
			rawarrholder = append(rawarrholder, apparr)
		case 7:
			apparr := h.Payload7(int(h.Keyholder[7]))
			rawarrholder = append(rawarrholder, apparr)
		case 8:
			rawarrholder = append(rawarrholder, []int{})
		case 9:
			apparr := h.Payload9(int(h.Keyholder[9]))
			rawarrholder = append(rawarrholder, apparr)
		case 10:
			apparr := h.Payload10(int(h.Keyholder[10]))
			rawarrholder = append(rawarrholder, apparr)
		case 11:
			apparr := h.Payload11(int(h.Keyholder[11]))
			rawarrholder = append(rawarrholder, apparr)
		case 12:
			apparr := h.Payload12(int(h.Keyholder[12]))
			rawarrholder = append(rawarrholder, apparr)
		case 13:
			apparr := h.Payload13(int(h.Keyholder[13]))
			rawarrholder = append(rawarrholder, apparr)
		case 14:
			apparr := h.Payload14(int(h.Keyholder[14]))
			rawarrholder = append(rawarrholder, apparr)
		case 15:
			apparr := h.Payload15(int(h.Keyholder[15]))
			rawarrholder = append(rawarrholder, apparr)
		case 16:
			apparr := h.Payload16(int(h.Keyholder[16]))
			rawarrholder = append(rawarrholder, apparr)
		case 17:
			apparr := h.Payload17(int(h.Keyholder[17]))
			rawarrholder = append(rawarrholder, apparr)
		case 18:
			apparr := h.Payload18(int(h.Keyholder[18]))
			rawarrholder = append(rawarrholder, apparr)
		case 19:
			apparr := h.Payload19(int(h.Keyholder[19]))
			rawarrholder = append(rawarrholder, apparr)
		case 20:
			apparr := h.Payload20(int(h.Keyholder[20]))
			rawarrholder = append(rawarrholder, apparr)
		case 21:
			apparr := h.Payload21(int(h.Keyholder[21]))
			rawarrholder = append(rawarrholder, apparr)
		case 22:
			apparr := h.Payload22(int(h.Keyholder[22]))
			rawarrholder = append(rawarrholder, apparr)
		case 23:
			apparr := h.Payload23(int(h.Keyholder[23]))
			rawarrholder = append(rawarrholder, apparr)
		case 24:
			apparr := h.Payload24(int(h.Keyholder[24]))
			rawarrholder = append(rawarrholder, apparr)
		case 25:
			apparr := h.Payload25(int(h.Keyholder[25]))
			rawarrholder = append(rawarrholder, apparr)
		case 26:
			apparr := h.Payload26(int(h.Keyholder[26]))
			rawarrholder = append(rawarrholder, apparr)
		case 27:
			apparr := h.Payload27(int(h.Keyholder[27]))
			rawarrholder = append(rawarrholder, apparr)
		case 28:
			apparr := h.Payload28(int(h.Keyholder[28]))
			rawarrholder = append(rawarrholder, apparr)
		case 29:
			apparr := h.Payload29(int(h.Keyholder[29]))
			rawarrholder = append(rawarrholder, apparr)
		case 30:
			apparr := h.Payload30(int(h.Keyholder[30]))
			rawarrholder = append(rawarrholder, apparr)
		case 31:
			apparr := h.Payload31(int(h.Keyholder[31]))
			rawarrholder = append(rawarrholder, apparr)
		case 32:
			apparr := h.Payload32(int(h.Keyholder[32]))
			rawarrholder = append(rawarrholder, apparr)
		case 33:
			apparr := h.Payload33(int(h.Keyholder[33]))
			rawarrholder = append(rawarrholder, apparr)
		case 34:
			apparr := h.Payload34(int(h.Keyholder[34]))
			rawarrholder = append(rawarrholder, apparr)
		case 35:
			apparr := h.Payload35(int(h.Keyholder[35]))
			rawarrholder = append(rawarrholder, apparr)
		case 36:
			apparr := h.Payload36(int(h.Keyholder[36]))
			rawarrholder = append(rawarrholder, apparr)
		case 37:
			apparr := h.Payload37(int(h.Keyholder[37]))
			rawarrholder = append(rawarrholder, apparr)
		case 38:
			apparr := h.Payload38(int(h.Keyholder[38]))
			rawarrholder = append(rawarrholder, apparr)
		case 39:
			apparr := h.Payload39(int(h.Keyholder[39]))
			rawarrholder = append(rawarrholder, apparr)
		case 40:
			apparr := h.Payload40(int(h.Keyholder[40]))
			rawarrholder = append(rawarrholder, apparr)
		case 41:
			apparr := h.Payload41(int(h.Keyholder[41]))
			rawarrholder = append(rawarrholder, apparr)
		case 42:
			apparr := h.Payload42(int(h.Keyholder[42]))
			rawarrholder = append(rawarrholder, apparr)
		case 43:
			apparr := h.Payload43(int(h.Keyholder[43]))
			rawarrholder = append(rawarrholder, apparr)
		case 44:
			apparr := h.Payload44(int(h.Keyholder[44]))
			rawarrholder = append(rawarrholder, apparr)
		}
	}

	superpackobjarr, err := h.PayloadLast(h.SeedString, h.ResString)
	if err != nil {
		log.Fatal("superpack error")
	}
	rawarrholder = append(rawarrholder, superpackobjarr)

	for i, arr := range rawarrholder {
		if i != len(rawarrholder)-1 {
			encarr = append(encarr, 0)
		}
		for _, val := range SuperpackLengthEncrypt(len(arr)) {
			encarr = append(encarr, val)
		}
		for _, val := range arr {
			encarr = append(encarr, val)
		}
	}

	indata := Encrypt(encarr, h.GlobalKeys, h.BaseKeys, rand.Intn(4294967295), rand.Intn(4294967295), h.Calckey, h.Encrounds)

	outstr := AlphabetEncode(indata, h.Alphabet)
	outmap := make(map[string]string)
	outmap["X-GyJwza5Z-a"] = outstr
	if len(outstr) >= 7900 {
		outmap["X-GyJwza5Z-a"] = outstr[:7900]
		outmap["X-GyJwza5Z-a0"] = outstr[7900:]
	}
	outmap["X-GyJwza5Z-b"] = strconv.FormatInt(int64(StringHash(h.SeedString+outstr)), 36)
	outmap["X-GyJwza5Z-c"] = h.HashString
	outmap["X-GyJwza5Z-d"] = h.Dheader
	outmap["X-GyJwza5Z-f"] = h.SeedString
	outmap["X-GyJwza5Z-z"] = "q"

	return outmap
}
func AlphabetEncode(arrin []int, alphabetin string) string {
	outstr := strings.Builder{}
	lenmod := len(arrin) % 3
	roundlength := len(arrin) - lenmod
	i := 0
	for ; i < roundlength; {
		outstr.WriteString(string(alphabetin[arrin[i]>>2]) + string(alphabetin[((arrin[i]&3)<<4)|(arrin[i+1]>>4)]) + string(alphabetin[((arrin[i+1]&15)<<2)|(arrin[i+2]>>6)]) + string(alphabetin[(arrin[i+2]&63)]))
		i += 3
	}
	switch lenmod {
	case 2:
		outstr.WriteString(string(alphabetin[arrin[i]>>2]) + string(alphabetin[((arrin[i]&3)<<4)|(arrin[i+1]>>4)]) + string(alphabetin[(arrin[i+1]&15)<<2]))
	case 1:
		outstr.WriteString(string(alphabetin[arrin[i]>>2]) + string(alphabetin[((arrin[i]&3)<<4)]))
	}
	return outstr.String()
}
func Encrypt(data, globalkeysin, basekeysin []int, iv1in, iv2in, calckeyin, encrounds int) []int {
	iv1 := int(int32(iv1in))
	iv2 := int(int32(iv2in))
	calckeycopy := calckeyin
	outarr := []int{(iv1 >> 24) & 255, (iv1 >> 16) & 255, (iv1 >> 8) & 255, (iv1) & 255, (iv2 >> 24) & 255, (iv2 >> 16) & 255, (iv2 >> 8) & 255, (iv2) & 255}
	keyvar := []int{}
	for _, val := range globalkeysin {
		keyvar = append(keyvar, int(int32(val&4294967295)))
	}
	for _, val := range basekeysin {
		keyvar = append(keyvar, int(int32(val&4294967295)))
	}
	keyvar = append(keyvar, 0)
	keyvar = append(keyvar, 0)
	keyvar = append(keyvar, int(int32(iv1&4294967295)))
	keyvar = append(keyvar, int(int32(iv2&4294967295)))
	keycopy := keyvar[:]
	rounds := int(math.Ceil(float64(len(data)) / 64))
	roundindexer := 0
	for rounds > 0 {
		ncc := int(int32(calckeycopy))
		if ncc < 0 {
			keycopy[12] = ncc + 4294967296
		} else {
			keycopy[12] = ncc
		}
		keycopy[13] = int(math.Floor(float64(calckeycopy / 4294967296)))
		copyslice := []int{}
		for _, val := range keycopy {
			copyslice = append(copyslice, val)
		}
		mainindexer := 0
		for mainindexer < encrounds {
			copyslice[0] = int(int32(int(int32(copyslice[0]+copyslice[4])) & 4294967295))
			copyslice[12] = int(int32(int(int32(int(int32((copyslice[12]^copyslice[0])<<16&4294967295))|int(uint32(copyslice[12]^copyslice[0]))>>16)) & 4294967295))
			copyslice[8] = int(int32(int(int32(copyslice[8]+copyslice[12])) & 4294967295))
			copyslice[4] = int(int32(int(int32(int(int32((copyslice[4]^copyslice[8])<<12&4294967295))|int(uint32(copyslice[4]^copyslice[8]))>>20)) & 4294967295))
			copyslice[0] = int(int32(int(int32(copyslice[0]+copyslice[4])) & 4294967295))
			copyslice[12] = int(int32(int(int32(int(int32((copyslice[12]^copyslice[0])<<8&4294967295))|int(uint32(copyslice[12]^copyslice[0]))>>24)) & 4294967295))
			copyslice[8] = int(int32(int(int32(copyslice[8]+copyslice[12])) & 4294967295))
			copyslice[4] = int(int32(int(int32(int(int32((copyslice[4]^copyslice[8])<<7&4294967295))|int(uint32(copyslice[4]^copyslice[8]))>>25)) & 4294967295))
			copyslice[1] = int(int32(int(int32(copyslice[1]+copyslice[5])) & 4294967295))
			copyslice[13] = int(int32(int(int32(int(int32((copyslice[13]^copyslice[1])<<16&4294967295))|int(uint32(copyslice[13]^copyslice[1]))>>16)) & 4294967295))
			copyslice[9] = int(int32(int(int32(copyslice[9]+copyslice[13])) & 4294967295))
			copyslice[5] = int(int32(int(int32(int(int32((copyslice[5]^copyslice[9])<<12&4294967295))|int(uint32(copyslice[5]^copyslice[9]))>>20)) & 4294967295))
			copyslice[1] = int(int32(int(int32(copyslice[1]+copyslice[5])) & 4294967295))
			copyslice[13] = int(int32(int(int32(int(int32((copyslice[13]^copyslice[1])<<8&4294967295))|int(uint32(copyslice[13]^copyslice[1]))>>24)) & 4294967295))
			copyslice[9] = int(int32(int(int32(copyslice[9]+copyslice[13])) & 4294967295))
			copyslice[5] = int(int32(int(int32(int(int32((copyslice[5]^copyslice[9])<<7&4294967295))|int(uint32(copyslice[5]^copyslice[9]))>>25)) & 4294967295))
			copyslice[2] = int(int32(int(int32(copyslice[2]+copyslice[6])) & 4294967295))
			copyslice[14] = int(int32(int(int32(int(int32((copyslice[14]^copyslice[2])<<16&4294967295))|int(uint32(copyslice[14]^copyslice[2]))>>16)) & 4294967295))
			copyslice[10] = int(int32(int(int32(copyslice[10]+copyslice[14])) & 4294967295))
			copyslice[6] = int(int32(int(int32(int(int32((copyslice[6]^copyslice[10])<<12&4294967295))|int(uint32(copyslice[6]^copyslice[10]))>>20)) & 4294967295))
			copyslice[2] = int(int32(int(int32(copyslice[2]+copyslice[6])) & 4294967295))
			copyslice[14] = int(int32(int(int32(int(int32((copyslice[14]^copyslice[2])<<8&4294967295))|int(uint32(copyslice[14]^copyslice[2]))>>24)) & 4294967295))
			copyslice[10] = int(int32(int(int32(copyslice[10]+copyslice[14])) & 4294967295))
			copyslice[6] = int(int32(int(int32(int(int32((copyslice[6]^copyslice[10])<<7&4294967295))|int(uint32(copyslice[6]^copyslice[10]))>>25)) & 4294967295))
			copyslice[3] = int(int32(int(int32(copyslice[3]+copyslice[7])) & 4294967295))
			copyslice[15] = int(int32(int(int32(int(int32((copyslice[15]^copyslice[3])<<16&4294967295))|int(uint32(copyslice[15]^copyslice[3]))>>16)) & 4294967295))
			copyslice[11] = int(int32(int(int32(copyslice[11]+copyslice[15])) & 4294967295))
			copyslice[7] = int(int32(int(int32(int(int32((copyslice[7]^copyslice[11])<<12&4294967295))|int(uint32(copyslice[7]^copyslice[11]))>>20)) & 4294967295))
			copyslice[3] = int(int32(int(int32(copyslice[3]+copyslice[7])) & 4294967295))
			copyslice[15] = int(int32(int(int32(int(int32((copyslice[15]^copyslice[3])<<8&4294967295))|int(uint32(copyslice[15]^copyslice[3]))>>24)) & 4294967295))
			copyslice[11] = int(int32(int(int32(copyslice[11]+copyslice[15])) & 4294967295))
			copyslice[7] = int(int32(int(int32(int(int32((copyslice[7]^copyslice[11])<<7&4294967295))|int(uint32(copyslice[7]^copyslice[11]))>>25)) & 4294967295))
			copyslice[0] = int(int32(int(int32(copyslice[0]+copyslice[5])) & 4294967295))
			copyslice[15] = int(int32(int(int32(int(int32((copyslice[15]^copyslice[0])<<16&4294967295))|int(uint32(copyslice[15]^copyslice[0]))>>16)) & 4294967295))
			copyslice[10] = int(int32(int(int32(copyslice[10]+copyslice[15])) & 4294967295))
			copyslice[5] = int(int32(int(int32(int(int32((copyslice[5]^copyslice[10])<<12&4294967295))|int(uint32(copyslice[5]^copyslice[10]))>>20)) & 4294967295))
			copyslice[0] = int(int32(int(int32(copyslice[0]+copyslice[5])) & 4294967295))
			copyslice[15] = int(int32(int(int32(int(int32((copyslice[15]^copyslice[0])<<8&4294967295))|int(uint32(copyslice[15]^copyslice[0]))>>24)) & 4294967295))
			copyslice[10] = int(int32(int(int32(copyslice[10]+copyslice[15])) & 4294967295))
			copyslice[5] = int(int32(int(int32(int(int32((copyslice[5]^copyslice[10])<<7&4294967295))|int(uint32(copyslice[5]^copyslice[10]))>>25)) & 4294967295))
			copyslice[1] = int(int32(int(int32(copyslice[1]+copyslice[6])) & 4294967295))
			copyslice[12] = int(int32(int(int32(int(int32((copyslice[12]^copyslice[1])<<16&4294967295))|int(uint32(copyslice[12]^copyslice[1]))>>16)) & 4294967295))
			copyslice[11] = int(int32(int(int32(copyslice[11]+copyslice[12])) & 4294967295))
			copyslice[6] = int(int32(int(int32(int(int32((copyslice[6]^copyslice[11])<<12&4294967295))|int(uint32(copyslice[6]^copyslice[11]))>>20)) & 4294967295))
			copyslice[1] = int(int32(int(int32(copyslice[1]+copyslice[6])) & 4294967295))
			copyslice[12] = int(int32(int(int32(int(int32((copyslice[12]^copyslice[1])<<8&4294967295))|int(uint32(copyslice[12]^copyslice[1]))>>24)) & 4294967295))
			copyslice[11] = int(int32(int(int32(copyslice[11]+copyslice[12])) & 4294967295))
			copyslice[6] = int(int32(int(int32(int(int32((copyslice[6]^copyslice[11])<<7&4294967295))|int(uint32(copyslice[6]^copyslice[11]))>>25)) & 4294967295))
			copyslice[2] = int(int32(int(int32(copyslice[2]+copyslice[7])) & 4294967295))
			copyslice[13] = int(int32(int(int32(int(int32((copyslice[13]^copyslice[2])<<16&4294967295))|int(uint32(copyslice[13]^copyslice[2]))>>16)) & 4294967295))
			copyslice[8] = int(int32(int(int32(copyslice[8]+copyslice[13])) & 4294967295))
			copyslice[7] = int(int32(int(int32(int(int32((copyslice[7]^copyslice[8])<<12&4294967295))|int(uint32(copyslice[7]^copyslice[8]))>>20)) & 4294967295))
			copyslice[2] = int(int32(int(int32(copyslice[2]+copyslice[7])) & 4294967295))
			copyslice[13] = int(int32(int(int32(int(int32((copyslice[13]^copyslice[2])<<8&4294967295))|int(uint32(copyslice[13]^copyslice[2]))>>24)) & 4294967295))
			copyslice[8] = int(int32(int(int32(copyslice[8]+copyslice[13])) & 4294967295))
			copyslice[7] = int(int32(int(int32(int(int32((copyslice[7]^copyslice[8])<<7&4294967295))|int(uint32(copyslice[7]^copyslice[8]))>>25)) & 4294967295))
			copyslice[3] = int(int32(int(int32(copyslice[3]+copyslice[4])) & 4294967295))
			copyslice[14] = int(int32(int(int32(int(int32((copyslice[14]^copyslice[3])<<16&4294967295))|int(uint32(copyslice[14]^copyslice[3]))>>16)) & 4294967295))
			copyslice[9] = int(int32(int(int32(copyslice[9]+copyslice[14])) & 4294967295))
			copyslice[4] = int(int32(int(int32(int(int32((copyslice[4]^copyslice[9])<<12&4294967295))|int(uint32(copyslice[4]^copyslice[9]))>>20)) & 4294967295))
			copyslice[3] = int(int32(int(int32(copyslice[3]+copyslice[4])) & 4294967295))
			copyslice[14] = int(int32(int(int32(int(int32((copyslice[14]^copyslice[3])<<8&4294967295))|int(uint32(copyslice[14]^copyslice[3]))>>24)) & 4294967295))
			copyslice[9] = int(int32(int(int32(copyslice[9]+copyslice[14])) & 4294967295))
			copyslice[4] = int(int32(int(int32(int(int32((copyslice[4]^copyslice[9])<<7&4294967295))|int(uint32(copyslice[4]^copyslice[9]))>>25)) & 4294967295))
			mainindexer += 2
		}
		for thirdlastround := 0; thirdlastround < 16; thirdlastround++ {
			copyslice[thirdlastround] = (copyslice[thirdlastround] + keycopy[thirdlastround]) % 9007199254740992
		}
		roundoutarr := [64]int{}
		for secondlastroundindexer := 0; secondlastroundindexer < 16; secondlastroundindexer++ {
			roundoutarr[secondlastroundindexer*4] = int(int32(copyslice[secondlastroundindexer] & 255))
			roundoutarr[secondlastroundindexer*4+1] = int(int32((copyslice[secondlastroundindexer] >> 8) & 255))
			roundoutarr[secondlastroundindexer*4+2] = int(int32((copyslice[secondlastroundindexer] >> 16) & 255))
			roundoutarr[secondlastroundindexer*4+3] = int(int32((copyslice[secondlastroundindexer] >> 24) & 255))
		}
		for appendindexer := 0; (appendindexer < len(roundoutarr)) && roundindexer+appendindexer < len(data); appendindexer++ {
			outarr = append(outarr, roundoutarr[appendindexer]^data[roundindexer+appendindexer])
		}
		calckeycopy = (calckeycopy + 1) % 9007199254740992
		roundindexer += 64
		rounds -= 1
	}
	return outarr
}
func SuperpackLengthEncrypt(numberin int) []int {
	outarr := []int{}
	for numberin>>7 > 0 {
		outarr = append(outarr, (numberin&127)|128)
		numberin = numberin >> 7
	}
	outarr = append(outarr, numberin)
	return outarr
}
func StringHash(stringin string) int {
	outarr := []int{}
	for i := 0; i < len(stringin); i += 2 {
		if len(stringin)-i == 1 {
			outarr = append(outarr, int(rune(stringin[i]))<<16)
		} else {
			outarr = append(outarr, int(rune(stringin[i]))<<16|int(rune(stringin[i+1])))
		}
	}
	outint := 0
	for _, val := range outarr {
		outint = int(int32((outint << 5) - outint + val))
	}
	return outint
}
func Float64bytes(float float64) []int {
	bits := math.Float64bits(float)
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, bits)
	outarr := []int{}
	for _, v := range bytes {
		outarr = append(outarr, int(v))
	}
	return outarr
}
func BaseNumEnc(numraw interface{}) []int {
	outarr := []int{}
	switch numraw.(type) {
	case int:
		shiftval := 0
		numin := numraw.(int)
		if numin < 0 {
			numin = -numin
			shiftval = 64
		}
		if int(int32(numin)) != numin {
			bytearray := Float64bytes(float64(numraw.(int)))
			outarr = append(outarr, 128)
			for _, v := range bytearray {
				outarr = append(outarr, int(uint8(v)))
			}
		} else {
			if numin < 32 {
				outarr = append(outarr, numin)
			} else {
				outarr = append(outarr, (32+shiftval)|(numin&31))
				if numin < 4096 {
					outarr = append(outarr, (numin>>5)&127)
				} else {
					outarr = append(outarr, 128|((numin>>5)&127))
					if numin < 524288 {
						outarr = append(outarr, (numin>>12)&127)
					} else {
						outarr = append(outarr, 128|((numin>>12)&127))
						if numin < 67108864 {
							outarr = append(outarr, (numin>>19)&127)
						} else {
							outarr = append(outarr, 128|((numin>>19)&127))
							outarr = append(outarr, int(uint32(numin))>>26)
						}
					}
				}
			}
		}
	case float64:
		bytearray := Float64bytes(numraw.(float64))
		outarr = append(outarr, 128)
		for _, v := range bytearray {
			outarr = append(outarr, int(uint8(v)))
		}
	}
	return outarr
}

func (h *VersionPayloadHolder) NumberEncrypt(numraw interface{}) []int {
	outarr := []int{}
	switch numraw.(type) {
	case int:
		shiftval := 0
		numin := numraw.(int)
		if numin < 0 {
			numin = - numin
			shiftval = 64
		}
		if int(int32(numin)) != numin {
			bytearray := Float64bytes(float64(numraw.(int)))
			outarr = append(outarr, 128)
			for _, v := range bytearray {
				outarr = append(outarr, int(uint8(v))^h.RandomFromSeed())
			}
		} else {
			if numin < 32 {
				outarr = append(outarr, numin^h.RandomFromSeed())
			} else {
				outarr = append(outarr, ((32+shiftval)|(numin&31))^h.RandomFromSeed())
				if numin < 4096 {
					outarr = append(outarr, ((numin>>5)&127)^h.RandomFromSeed())
				} else {
					outarr = append(outarr, (128|((numin>>5)&127))^h.RandomFromSeed())
					if numin < 524288 {
						outarr = append(outarr, ((numin>>12)&127)^h.RandomFromSeed())
					} else {
						outarr = append(outarr, (128|((numin>>12)&127))^h.RandomFromSeed())
						if numin < 67108864 {
							outarr = append(outarr, ((numin>>19)&127)^h.RandomFromSeed())
						} else {
							outarr = append(outarr, (128|((numin>>19)&127))^h.RandomFromSeed())
							outarr = append(outarr, (int(uint32(numin))>>26)^h.RandomFromSeed())
						}
					}
				}
			}
		}
	case float64:
		bytearray := Float64bytes(numraw.(float64))
		outarr = append(outarr, 128)
		for _, v := range bytearray {
			outarr = append(outarr, int(uint8(v))^h.RandomFromSeed())
		}
	}
	return outarr
}
func (h *VersionPayloadHolder) RandomFromSeed() int {
	if h.Seedcount%4 == 0 {
		h.Seedcount %= 4
		h.Currval = int(int32(h.Hashseed1 + h.Hashseed2))
		h.Hashseed2 ^= h.Hashseed1
		h.Hashseed1 = int(int32(h.Hashseed2 ^ ((h.Hashseed1 << 26) | (int(uint32(h.Hashseed1)) >> 6)) ^ (h.Hashseed2 << 9)))
		h.Hashseed2 = int(int32(int(int32(h.Hashseed2<<13)) | (int(uint32(h.Hashseed2)) >> 19)))
		h.Seedcount++
		return h.Currval & 255
	} else {
		h.Seedcount++
		return (h.Currval >> (8 * (h.Seedcount - 1))) & 255
	}
}
func (h *VersionPayloadHolder) SetRandInt(keyin int) []int {
	randint := int(int32(rand.Uint32()))
	h.Hashseed1 = h.Hashseedbase
	h.Hashseed2 = int(int32(randint ^ keyin))
	h.Seedcount = 0
	h.RandomFromSeed()
	h.RandomFromSeed()
	h.RandomFromSeed()
	h.RandomFromSeed()
	return []int{int(uint32(randint)) >> 24 & 255, int(uint32(randint)) >> 16 & 255, int(uint32(randint)) >> 8 & 255, int(uint32(randint)) & 255}
}

func MouseGen() ([]MouseEvent, []MouseTarget, []MouseEvent) {
	outarr1 := []MouseEvent{}
	outarr2 := []MouseTarget{}
	outarr3 := []MouseEvent{}
	startts1 := 500 + (rand.Float64() * 100)
	startts20 := float64(1000 + rand.Intn(200))
	startts21 := float64(1800 + rand.Intn(200))
	startts3 := 800 + (rand.Float64() * 200)
	x1 := 100 + rand.Intn(300)
	x20 := 100 + rand.Intn(300)
	x21 := 100 + rand.Intn(300)
	x3 := 100 + rand.Intn(300)
	y1 := rand.Intn(400)
	y20 := 200 + rand.Intn(300)
	y21 := 200 + rand.Intn(300)
	y3 := 200 + rand.Intn(300)
	docrelx0 := 200 + rand.Intn(100)
	docrely0 := 100 + rand.Intn(100)
	docrelx1 := docrelx0 + rand.Intn(100)
	docrely1 := docrely0 + rand.Intn(200)

	for i := 0; i < 1+rand.Intn(3); i++ {
		appevent := MouseEvent{
			DocumentRelativeX: x1,
			DocumentRelativeY: y1,
			Timestamp:         float64(startts1),
		}
		outarr1 = append(outarr1, appevent)
		startts1 += 100 + rand.Float64()*50
		if rand.Float64() > 0.5 {
			x1 += 3 + rand.Intn(100)
		} else {
			x1 -= 3 + rand.Intn(100)
			if x1 < 0 {
				x1 = int(math.Abs(float64(x1)))
			}
		}
		if rand.Float64() > 0.5 {
			y1 += 3 + rand.Intn(100)
		} else {
			y1 -= 3 + rand.Intn(100)
			if y1 < 0 {
				y1 = int(math.Abs(float64(y1)))
			}
		}
	}
	for i := 0; i < 3; i++ {
		appobj := MouseTarget{
			Button:            0,
			DocumentRelativeX: docrelx0,
			DocumentRelativeY: docrely0,
			EventType:         i + 1,
			TargetId:          "username",
			TargetName:        "username",
			TargetRelativeX:   x20,
			TargetRelativeY:   y20,
			Timestamp:         startts20,
		}
		outarr2 = append(outarr2, appobj)
		startts20 += 50 + rand.Float64()*50
	}
	for i := 0; i < 3; i++ {
		appobj := MouseTarget{
			Button:            0,
			DocumentRelativeX: docrelx1,
			DocumentRelativeY: docrely1,
			EventType:         i + 1,
			TargetId:          "login",
			TargetName:        "",
			TargetRelativeX:   x21,
			TargetRelativeY:   y21,
			Timestamp:         startts21,
		}
		outarr2 = append(outarr2, appobj)
		startts21 += 50 + rand.Float64()*50
	}
	for i := 0; i < 50; i++ {
		appevent := MouseEvent{
			DocumentRelativeX: x3,
			DocumentRelativeY: y3,
			Timestamp:         float64(startts3),
		}
		outarr3 = append(outarr3, appevent)
		startts3 += 8 + rand.Float64()*4
		if rand.Float64() > 0.5 {
			x3 += 3 + rand.Intn(100)
		} else {
			x3 -= 3 + rand.Intn(100)
			if x3 < 0 {
				x3 = int(math.Abs(float64(x1)))
			}
		}
		if rand.Float64() > 0.5 {
			y3 += 3 + rand.Intn(100)
		} else {
			y3 -= 3 + rand.Intn(100)
			if y3 < 0 {
				y3 = int(math.Abs(float64(y1)))
			}
		}
	}
	return outarr1, outarr2, outarr3
}
func KeyGen() []KeyEvent {
	outarr := []KeyEvent{}
	starttime := float64(1800 + rand.Intn(200))
	for i := 0; i < 2; i++ {
		appobj := KeyEvent{
			EventType:  1,
			keyCode:    1,
			TargetId:   "username",
			TargetName: "username",
			Timestamp:  starttime,
		}
		outarr = append(outarr, appobj)
		starttime += 4 + rand.Float64()*4
	}
	starttime += 1 + rand.Float64()*1
	for i := 0; i < 2; i++ {
		appobj := KeyEvent{
			EventType:  1,
			keyCode:    1,
			TargetId:   "password",
			TargetName: "password",
			Timestamp:  starttime,
		}
		outarr = append(outarr, appobj)
		starttime += 4 + rand.Float64()*4
		starttime = math.Floor(starttime)
	}
	return outarr
}

func (h *VersionPayloadHolder) MousePayload(keyin int) []int {
	outarr := h.SetRandInt(keyin)
	mouse1, mouse2, mouse3 := MouseGen()
	outarr = append(outarr, len(mouse1)^h.RandomFromSeed())
	for _, event := range mouse1 {
		for _, val := range h.NumberEncrypt(event.Timestamp) {
			outarr = append(outarr, val)
		}
		for _, val := range h.NumberEncrypt(event.DocumentRelativeX) {
			outarr = append(outarr, val)
		}
		for _, val := range h.NumberEncrypt(event.DocumentRelativeY) {
			outarr = append(outarr, val)
		}
	}

	outarr = append(outarr, len(mouse2)^h.RandomFromSeed())
	for _, event := range mouse2 {
		for _, val := range h.NumberEncrypt(event.TargetRelativeX) {
			outarr = append(outarr, val)
		}
		for _, char := range event.TargetId {
			outarr = append(outarr, int(char)^h.RandomFromSeed())
		}
		for _, val := range h.NumberEncrypt(event.Timestamp) {
			outarr = append(outarr, val)
		}
		if len(event.TargetName) == 0 {
			outarr = append(outarr, 0^h.RandomFromSeed())
		} else {
			for _, char := range event.TargetName {
				outarr = append(outarr, int(char)^h.RandomFromSeed())
			}
		}
		for _, val := range h.NumberEncrypt(event.DocumentRelativeX) {
			outarr = append(outarr, val)
		}
		for _, val := range h.NumberEncrypt(event.Button) {
			outarr = append(outarr, val)
		}
		for _, val := range h.NumberEncrypt(event.DocumentRelativeY) {
			outarr = append(outarr, val)
		}
		for _, val := range h.NumberEncrypt(event.TargetRelativeY) {
			outarr = append(outarr, val)
		}
		for _, val := range h.NumberEncrypt(event.EventType) {
			outarr = append(outarr, val)
		}
		outarr = append(outarr, 3^h.RandomFromSeed())
	}

	outarr = append(outarr, len(mouse3)^h.RandomFromSeed())
	for _, event := range mouse3 {
		for _, val := range h.NumberEncrypt(event.DocumentRelativeX) {
			outarr = append(outarr, val)
		}
		outarr = append(outarr, 3^h.RandomFromSeed())
		for _, val := range h.NumberEncrypt(event.Timestamp) {
			outarr = append(outarr, val)
		}
		for _, val := range h.NumberEncrypt(event.DocumentRelativeY) {
			outarr = append(outarr, val)
		}
	}
	return outarr
}
func (h *VersionPayloadHolder) KeyPayload(keyin int) []int {
	outarr := h.SetRandInt(keyin)
	keyArr := KeyGen()
	outarr = append(outarr, len(keyArr)^h.RandomFromSeed())
	for i := len(keyArr) - 1; i >= 0; i-- {

		for _, val := range h.NumberEncrypt(keyArr[i].Timestamp) {
			outarr = append(outarr, val)
		}
		for _, val := range h.NumberEncrypt(keyArr[i].keyCode) {
			outarr = append(outarr, val)
		}

		outarr = append(outarr, 4^h.RandomFromSeed())
		outarr = append(outarr, h.RandomFromSeed())

		switch len(keyArr) - i - 1 {
		case 0:
			for _, char := range keyArr[i].TargetName {
				outarr = append(outarr, int(char)^h.RandomFromSeed())
			}
		case 1:
			outarr = append(outarr, 129^h.RandomFromSeed())
			outarr = append(outarr, h.RandomFromSeed())
		case 2:
			for _, char := range keyArr[i].TargetName {
				outarr = append(outarr, int(char)^h.RandomFromSeed())
			}
		case 3:
			outarr = append(outarr, 1^h.RandomFromSeed())
		}

		for _, val := range h.NumberEncrypt(keyArr[i].EventType) {
			outarr = append(outarr, val)
		}
		outarr = append(outarr, 1^h.RandomFromSeed())
		switch len(keyArr) - i - 1 {
		case 0:
			outarr = append(outarr, 129^h.RandomFromSeed())
			outarr = append(outarr, h.RandomFromSeed())
		case 1:
			outarr = append(outarr, 129^h.RandomFromSeed())
			outarr = append(outarr, h.RandomFromSeed())
		case 2:
			outarr = append(outarr, 129^h.RandomFromSeed())
			outarr = append(outarr, 1^h.RandomFromSeed())
		case 3:
			outarr = append(outarr, 129^h.RandomFromSeed())
			outarr = append(outarr, 1^h.RandomFromSeed())
		}
	}
	return outarr
}

func (h *VersionPayloadHolder) Payload0(keyin int) []int {
	outarr := h.SetRandInt(keyin)
	for _, val := range []int{0} {
		outarr = append(outarr, val^h.RandomFromSeed())
	}
	return outarr
}
func (h *VersionPayloadHolder) Payload1(keyin int) []int {
	outarr := h.SetRandInt(keyin)
	for _, val := range h.P1VAL {
		outarr = append(outarr, val^h.RandomFromSeed())
	}
	return outarr
}
func (h *VersionPayloadHolder) Payload2(keyin int) []int {
	outarr := h.SetRandInt(keyin)
	for _, val := range h.P2VAL {
		outarr = append(outarr, val^h.RandomFromSeed())
	}
	return outarr
}
func (h *VersionPayloadHolder) Payload3(keyin int) []int {
	outarr := h.SetRandInt(keyin)
	for _, val := range h.P3VAL {
		outarr = append(outarr, val^h.RandomFromSeed())
	}
	return outarr
}
func (h *VersionPayloadHolder) Payload4(keyin int) []int {
	outarr := h.SetRandInt(keyin)
	for _, val := range []int{h.P4VAL} {
		outarr = append(outarr, val^h.RandomFromSeed())
	}
	return outarr
}
func (h *VersionPayloadHolder) Payload5(keyin int) []int {
	outarr := h.SetRandInt(keyin)
	for _, val := range []int{h.P5VAL} {
		outarr = append(outarr, val^h.RandomFromSeed())
	}
	return outarr
}
func (h *VersionPayloadHolder) Payload6(keyin int) []int {
	outarr := h.SetRandInt(keyin)
	for _, val := range []int{3, 0, 0, 0} {
		outarr = append(outarr, val^h.RandomFromSeed())
	}
	return outarr
}
func (h *VersionPayloadHolder) Payload7(keyin int) []int {
	outarr := h.SetRandInt(keyin)
	for _, val := range []int{127, 127, 63} {
		outarr = append(outarr, val^h.RandomFromSeed())
	}
	return outarr
}

//payload 8 empty
func (h *VersionPayloadHolder) Payload9(keyin int) []int {
	outarr := h.SetRandInt(keyin)
	for _, val := range []int{0, 0} {
		outarr = append(outarr, val^h.RandomFromSeed())
	}
	return outarr
}
func (h *VersionPayloadHolder) Payload10(keyin int) []int {
	outarr := h.SetRandInt(keyin)
	for _, val := range []int{117, 110, 100, 101, 102, 105, 110, 101, 100, 0, 31} {
		outarr = append(outarr, val^h.RandomFromSeed())
	}
	return outarr
}
func (h *VersionPayloadHolder) Payload11(keyin int) []int {
	outarr := h.SetRandInt(keyin)
	for _, val := range h.P11VAL {
		outarr = append(outarr, val^h.RandomFromSeed())
	}
	return outarr
}
func (h *VersionPayloadHolder) Payload12(keyin int) []int {
	outarr := h.SetRandInt(keyin)
	for _, val := range []int{56, 23} {
		outarr = append(outarr, val^h.RandomFromSeed())
	}
	return outarr
}
func (h *VersionPayloadHolder) Payload13(keyin int) []int {
	outarr := h.SetRandInt(keyin)
	for _, val := range h.P13VAL {
		outarr = append(outarr, val^h.RandomFromSeed())
	}
	return outarr
}
func (h *VersionPayloadHolder) Payload14(keyin int) []int {
	outarr := h.SetRandInt(keyin)
	for _, val := range h.P14VAL {
		outarr = append(outarr, val^h.RandomFromSeed())
	}
	return outarr
}
func (h *VersionPayloadHolder) Payload15(keyin int) []int {
	outarr := h.SetRandInt(keyin)
	for _, val := range h.P15VAL {
		outarr = append(outarr, val^h.RandomFromSeed())
	}
	return outarr
}
func (h *VersionPayloadHolder) Payload16(keyin int) []int {
	outarr := h.SetRandInt(keyin)
	for _, val := range h.P16VAL {
		outarr = append(outarr, val^h.RandomFromSeed())
	}
	return outarr
}
func (h *VersionPayloadHolder) Payload17(keyin int) []int {
	outarr := h.SetRandInt(keyin)
	chose := rand.Intn(len(h.P17VAL))
	for _, val := range h.P17VAL[chose] {
		outarr = append(outarr, val^h.RandomFromSeed())
	}
	return outarr
}
func (h *VersionPayloadHolder) Payload18(keyin int) []int {
	parr := []int{}
	outarr := h.SetRandInt(keyin)
	lenmod := len(h.P18VAL) % 7
	looprange := (len(h.P18VAL) - lenmod) / 7
	for j := 0; j < looprange; j++ {
		val := 0
		for i := 0; i < 7; i++ {
			val |= h.P18VAL[(j*7)+i] << i
		}
		parr = append(parr, val)
		outarr = append(outarr, val^h.RandomFromSeed())
	}
	lastval := 0
	for i := 0; i < lenmod; i++ {
		lastval |= h.P18VAL[(looprange*7)+i] << i
	}
	parr = append(parr, lastval)
	outarr = append(outarr, lastval^h.RandomFromSeed())
	return outarr
}
func (h *VersionPayloadHolder) Payload19(keyin int) []int {
	outarr := h.SetRandInt(keyin)
	for _, val := range h.P19VAL {
		outarr = append(outarr, val^h.RandomFromSeed())
	}
	return outarr
}
func (h *VersionPayloadHolder) Payload20(keyin int) []int {
	outarr := h.SetRandInt(keyin)
	for _, val := range h.P20VAL {
		outarr = append(outarr, val^h.RandomFromSeed())
	}
	return outarr
}
func (h *VersionPayloadHolder) Payload21(keyin int) []int {
	outarr := h.SetRandInt(keyin)
	for _, val := range h.P21VAL {
		outarr = append(outarr, val^h.RandomFromSeed())
	}
	return outarr
}
func (h *VersionPayloadHolder) Payload22(keyin int) []int {
	outarr := h.SetRandInt(keyin)
	for _, val := range []int{56, 135, 171, 185, 4} {
		outarr = append(outarr, val^h.RandomFromSeed())
	}
	return outarr
}
func (h *VersionPayloadHolder) Payload23(keyin int) []int {
	outarr := h.SetRandInt(keyin)
	for _, val := range h.P23VAL {
		outarr = append(outarr, val^h.RandomFromSeed())
	}
	return outarr
}
func (h *VersionPayloadHolder) Payload24(keyin int) []int {
	outarr := h.SetRandInt(keyin)
	for _, val := range h.P24VAL {
		outarr = append(outarr, val^h.RandomFromSeed())
	}
	return outarr
}
func (h *VersionPayloadHolder) Payload25(keyin int) []int {
	outarr := h.SetRandInt(keyin)
	for _, val := range h.P25VAL {
		outarr = append(outarr, val^h.RandomFromSeed())
	}
	return outarr
}
func (h *VersionPayloadHolder) Payload26(keyin int) []int {
	outarr := h.SetRandInt(keyin)
	for _, val := range []int{83, 117, 110, 32, 65, 117, 103, 32, 48, 53, 32, 49, 57, 52, 53, 32, 49, 57, 58, 49, 54, 58, 48, 48, 32, 71, 77, 84, 45, 48, 52, 48, 48, 32, 40, 69, 97, 115, 116, 101, 114, 110, 32, 68, 97, 121, 108, 105, 103, 104, 116, 32, 84, 105, 109, 101, 41, 0} {
		outarr = append(outarr, val^h.RandomFromSeed())
	}
	return outarr
}
func (h *VersionPayloadHolder) Payload27(keyin int) []int {
	outarr := h.SetRandInt(keyin)
	for _, val := range []int{0} {
		outarr = append(outarr, val^h.RandomFromSeed())
	}
	return outarr
}
func (h *VersionPayloadHolder) Payload28(keyin int) []int {
	outarr := h.SetRandInt(keyin)
	for _, val := range []int{0} {
		outarr = append(outarr, val^h.RandomFromSeed())
	}
	return outarr
}
func (h *VersionPayloadHolder) Payload29(keyin int) []int {
	outarr := h.SetRandInt(keyin)
	parr := []int{}
	for _, val := range h.NumberEncrypt(float64(int(time.Now().UTC().UnixNano() / 1e6))) {
		parr = append(parr, val)
		outarr = append(outarr, val)
	}
	return outarr
}
func (h *VersionPayloadHolder) Payload30(keyin int) []int {
	outarr := h.SetRandInt(keyin)
	for _, val := range h.P30VAL {
		outarr = append(outarr, val^h.RandomFromSeed())
	}
	return outarr
}
func (h *VersionPayloadHolder) Payload31(keyin int) []int {
	outarr := h.SetRandInt(keyin)
	for _, val := range h.P31VAL {
		outarr = append(outarr, val^h.RandomFromSeed())
	}
	return outarr
}
func (h *VersionPayloadHolder) Payload32(keyin int) []int {
	outarr := h.SetRandInt(keyin)
	for _, val := range h.P32VAL {
		outarr = append(outarr, val^h.RandomFromSeed())
	}
	return outarr
}
func (h *VersionPayloadHolder) Payload33(keyin int) []int {
	outarr := h.SetRandInt(keyin)
	parr := []int{}
	for _, val := range h.NumberEncrypt(32 + rand.Intn(4064)) {
		outarr = append(outarr, val)
		parr = append(parr, val)
	}
	for _, val := range h.NumberEncrypt(float64(int(time.Now().UTC().UnixNano() / 1e6))) {
		outarr = append(outarr, val)
		parr = append(parr, val)
	}
	for _, val := range h.NumberEncrypt(32 + rand.Intn(4064)) {
		outarr = append(outarr, val)
		parr = append(parr, val)
	}
	outarr = append(outarr, 0^h.RandomFromSeed())
	return outarr
}
func (h *VersionPayloadHolder) Payload34(keyin int) []int {
	outarr := h.SetRandInt(keyin)
	for _, val := range []int{0} {
		outarr = append(outarr, val^h.RandomFromSeed())
	}
	return outarr
}
func (h *VersionPayloadHolder) Payload35(keyin int) []int {
	h.Enabled = true
	outarr := h.SetRandInt(keyin)
	for _, val := range h.P35VAL {
		outarr = append(outarr, val^h.RandomFromSeed())
	}
	h.Enabled = false
	return outarr
}
func (h *VersionPayloadHolder) Payload36(keyin int) []int {
	outarr := h.SetRandInt(keyin)
	for _, val := range h.P36VAL {
		outarr = append(outarr, val^h.RandomFromSeed())
	}
	return outarr
}
func (h *VersionPayloadHolder) Payload37(keyin int) []int {
	outarr := h.SetRandInt(keyin)
	for _, val := range []int{127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 3} {
		outarr = append(outarr, val^h.RandomFromSeed())
	}
	return outarr
}
func (h *VersionPayloadHolder) Payload38(keyin int) []int {
	outarr := h.SetRandInt(keyin)
	for _, val := range []int{97, 192, 128, 103, 192, 128, 104, 192, 128, 107, 192, 128, 101, 192, 128, 102, 192, 128, 106, 192, 128, 105, 192, 128, 110, 192, 128, 100, 192, 128, 98, 192, 128, 108, 192, 128, 109, 192, 128, 111, 192, 128, 0} {
		outarr = append(outarr, val^h.RandomFromSeed())
	}
	return outarr
}
func (h *VersionPayloadHolder) Payload39(keyin int) []int {
	outarr := h.SetRandInt(keyin)
	for _, val := range []int{42, 31} {
		outarr = append(outarr, val^h.RandomFromSeed())
	}
	return outarr
}
func (h *VersionPayloadHolder) Payload40(keyin int) []int {
	outarr := h.SetRandInt(keyin)
	for _, val := range h.P40VAL {
		outarr = append(outarr, val^h.RandomFromSeed())
	}
	return outarr
}
func (h *VersionPayloadHolder) Payload41(keyin int) []int {
	outarr := h.SetRandInt(keyin)
	for _, val := range []int{53, 183, 143, 158, 18} {
		outarr = append(outarr, val^h.RandomFromSeed())
	}
	return outarr
}
func (h *VersionPayloadHolder) Payload42(keyin int) []int {
	outarr := h.SetRandInt(keyin)
	for _, val := range h.P42VAL {
		outarr = append(outarr, val^h.RandomFromSeed())
	}
	return outarr
}
func (h *VersionPayloadHolder) Payload43(keyin int) []int {
	outarr := h.SetRandInt(keyin)
	for _, val := range []int{128, 255, 255, 255, 255, 255, 255, 239, 67} {
		outarr = append(outarr, val^h.RandomFromSeed())
	}
	return outarr
}
func (h *VersionPayloadHolder) Payload44(keyin int) []int {
	outarr := h.SetRandInt(keyin)
	for _, val := range h.P44VAL {
		outarr = append(outarr, val^h.RandomFromSeed())
	}
	return outarr
}

func (h *VersionPayloadHolder) PayloadLast(seedin, hashstring string) ([]int, error) {
	rin := rand.Intn(32768) + 32768
	charsum := (65536 * int(seedin[10])) + (256 * int(seedin[11])) + int(seedin[12])
	out1 := 127 + rin - charsum
	out2 := rin + charsum

	return superpack.Encrypt([]interface{}{"nquunymwfzhl", []int{out1, out2}},
		[]interface{}{"qcfocgbborau", struct {
			Response  string `json:"response"`
			StyleHash int    `json:"style_hash"`
		}{
			Response:  hashstring,
			StyleHash: 1037328191,
		}})
}
