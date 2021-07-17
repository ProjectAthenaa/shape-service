package shape

import (
	"fmt"
	"github.com/valyala/fasthttp"
	"log"
	"main/shape/generation"
	"math"
	"math/rand"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	globalOpArrayRegex             = regexp.MustCompile("var \\w{2}=\\[function(.*?})]")
	globalStringArrayRegex         = regexp.MustCompile("var (\\w)=\\[\"(.*?)\"]")
	globalTupleArrayRegex          = regexp.MustCompile("var (\\w)=\\[\\[\\[(.*?)]]]")
	globalKeyArrayRegex            = regexp.MustCompile("var (\\w)=\\[([.-]?\\d.*?\\d)]")
	globalHeapCopyRegex            = regexp.MustCompile("var (\\w)=\\[\\{(.*?)}]")
	opVarRegex                     = regexp.MustCompile("\\((\\w{2})\\)")
	tupleObjectRegex               = regexp.MustCompile("(\\w):\\[([\\d,]*?)]")
	tupleObjectIndexReferenceRegex = regexp.MustCompile("(\\w):([\\d]+),")
	funcArrayRegex                 = regexp.MustCompile("var (\\w)=\\[(?:\\w,){26}\\w]")
	functionDefsRegex              = regexp.MustCompile("var (\\w)=ReferenceError,(\\w)=TypeError,(\\w)=Object,(\\w)=RegExp,(\\w)=Number,(\\w)=String,(\\w)=Array,(\\w)=\\w.bind,(\\w)=\\w.call,(\\w)=\\w.bind\\(\\w,\\w\\),(\\w)=\\w.apply,(\\w)=\\w\\(\\w\\),(\\w)=\\[].push,(\\w)=\\[].pop,(\\w)=\\[].slice,(\\w)=\\[].splice,(\\w)=\\[].join,(\\w)=\\[].map,(\\w)=\\w\\(\\w\\),(\\w)=\\w\\(\\w\\),(\\w)=\\w\\(\\w\\),(\\w)=\\w\\(\\w\\),(\\w)=\\{}.hasOwnProperty,(\\w)=\\w\\(\\w\\),(\\w)=JSON.stringify,(\\w)=\\w.getOwnPropertyDescriptor,(\\w)=\\w.defineProperty,(\\w)=\\w.fromCharCode,(\\w)=Math.min,(\\w)=Math.floor,(\\w)=\\w.create,(\\w)=\"\".indexOf,(\\w)=\"\".charAt,(\\w)=\\w\\(\\w\\),(\\w)=\\w\\(\\w\\),(\\w)=typeof Uint8Array===\"function\"\\?Uint8Array:\\w;")
	funcObjectRegex                = regexp.MustCompile("this\\.(\\w)=\\w{2}\\(\\);this\\.(\\w)=\\w{2}\\(\\);this\\.(\\w)=\\w{2}\\(\\);this\\.(\\w)=void 0;this\\.(\\w)=\\w{2};this\\.(\\w)=\\w{2};this\\.(\\w)=\\w{2};this\\.(\\w)=\\w{2}==null\\?\\w:Object\\(\\w{2}\\);this\\.(\\w)=\\w{2};this\\.(\\w)=0")
	objectDefinitionRegex          = regexp.MustCompile("Object.defineProperty\\(\\w{2},\"(\\w)\"")
	XYCordRegex                    = regexp.MustCompile("globalTupleArray\\[this\\.(\\w)]\\[(\\w{2})\\[this\\.(\\w)\\+\\+]]")
	runFuncDeclarationRegex        = regexp.MustCompile("function (\\w{2})\\((?:\\w{2},){3}\\w{2}\\)\\{\"use strict\"")
	runFuncCopyIdentifiersRegex    = regexp.MustCompile("return \\w{2}\\(\\w{2},\\w{2},\\w{2},\\w{2}\\.(\\w),\\w{2}\\.(\\w),\\w{2}\\.(\\w),\\w{2}\\.(\\w),\\w{2}\\.(\\w)\\)")
	breakObjRegex                  = regexp.MustCompile("\\W(\\w{2})!==\\w\\W")
	throwMatchRegex                = regexp.MustCompile("throw ")

	b0matchregex       = regexp.MustCompile("b0:\\{(.*?}.*?}.*?)}")
	b1matchregex       = regexp.MustCompile("b1:\\{(.*?}.*?}.*?)}")
	rawstringdeobregex = regexp.MustCompile("var \\w{2}=(\\w{2})\\+\".*?\\)}")
	throwifregex       = regexp.MustCompile("if\\(([^{}]*?)\\)\\{([^{}]*?(?:throw)[^{}]*?)}")
	jumpifregex        = regexp.MustCompile("if\\(([^{}]*?)\\)\\{([^{}]*?(?:Xcoord)[^{}]*?)}")
	breakmatchregex    = regexp.MustCompile("breakobject=")
	forvarinregex      = regexp.MustCompile("for\\(var [A-Za-z]{2} in[^)]+\\)\\{[^}]+}")

	XCordIncrementRegex     = regexp.MustCompile("^varin\\.Xcoord\\+=(\\d+)$")
	lengthDecrementRegex    = regexp.MustCompile("^varin\\.stack1\\.length-=(\\d+)(?:\\+([A-Za-z]{2}))?$")
	varSetRegex             = regexp.MustCompile("^var ([A-Za-z]{2})=")
	heapSetRegex            = regexp.MustCompile("^varin\\.heap\\[(.*?)]\\.selector=")
	stackSetStackRegex      = regexp.MustCompile("^varin\\.stack1\\[varin\\.stack1\\.length(?:-(\\d+))?]=")
	stackSetVarRegex        = regexp.MustCompile("^varin\\.stack1\\[([A-Za-z]{2})(?:\\+(\\d+))?]=")
	objectFromStackSetRegex = regexp.MustCompile("^varin\\.stack1\\[varin\\.stack1\\.length-2]\\[varin\\.stack1\\[varin\\.stack1\\.length-1]]=")
	XCoordSetRegex          = regexp.MustCompile("^varin\\.Xcoord=")
	YCoordSetRegex          = regexp.MustCompile("^varin\\.Ycoord=")
	jumpHolderSetRegex      = regexp.MustCompile("^varin\\.explicitJumpHolder=")
	stackPushRegex          = regexp.MustCompile("varin\\.stack1\\.push\\(([^)]*?)\\)")
	keyArrayMatchRegex      = regexp.MustCompile("=globalKeyArray\\[mainnumarr\\[varin\\.Xcoord(?:\\+(\\d+))?]]$")

	stackLengthGetRegex        = regexp.MustCompile("=varin\\.stack1\\.length(?:-(\\d))?$")
	globalStringArrGetRegex    = regexp.MustCompile("=globalStringArray\\[mainnumarr\\[varin\\.Xcoord]<<8\\|mainnumarr\\[varin\\.Xcoord\\+(\\d+)]]$")
	globalStringArrOneGetRegex = regexp.MustCompile("=globalStringArray\\[mainnumarr\\[varin\\.Xcoord\\+(\\d+)]<<8\\|mainnumarr\\[varin\\.Xcoord\\+(\\d+)]]$")
	mainnumarrGetRegex         = regexp.MustCompile("=mainnumarr\\[varin\\.Xcoord]$")
	mainnumarrOneGetRegex      = regexp.MustCompile("=mainnumarr\\[varin\\.Xcoord\\+(\\d+)]$")
	mainnumarrTwoGetRegex      = regexp.MustCompile("=mainnumarr\\[varin\\.Xcoord(?:\\+(\\d+))?]<<8\\|mainnumarr\\[varin\\.Xcoord\\+(\\d+)]$")
	mainnumarrThreeGetRegex    = regexp.MustCompile("=mainnumarr\\[varin\\.Xcoord]<<16\\|\\(mainnumarr\\[varin\\.Xcoord\\+1]<<8\\|mainnumarr\\[varin\\.Xcoord\\+2]\\)$")
	stackGetRegex              = regexp.MustCompile("=varin\\.stack1\\[varin\\.stack1\\.length-(\\d+)]$")
	varGetRegex                = regexp.MustCompile("=([A-Za-z]{2})$")
	runFuncGetRegex            = regexp.MustCompile("=runFuncDeclaration\\(([^(),]*?),([^(),]*?),([^(),]*?),varin\\.heap\\)")
	heapGetRegex               = regexp.MustCompile("=varin\\.heap\\[([A-Za-z]{2})]\\.selector$")
	jumpHolderXGetRegex        = regexp.MustCompile("=varin\\.explicitJumpHolder\\.Xcoord$")
	jumpHolderYGetRegex        = regexp.MustCompile("=varin\\.explicitJumpHolder\\.Ycoord$")
	arrCreateRegex             = regexp.MustCompile("=\\[]$")
	objCreateRegex             = regexp.MustCompile("=\\{}$")
	strCreateRegex             = regexp.MustCompile("=\"\"$")
	stackObjectGetRegex        = regexp.MustCompile("=varin\\.stack1\\[varin\\.stack1\\.length-2]\\[varin\\.stack1\\[varin\\.stack1\\.length-1]]$")
	stackObjectFuncGetRegex    = regexp.MustCompile("=varin\\.stack1\\[varin\\.stack1\\.length-2]\\[varin\\.stack1\\[varin\\.stack1\\.length-1]]\\(\\)$")
	stackModRegex              = regexp.MustCompile("=varin\\.stack1\\[varin\\.stack1\\.length-2]%varin\\.stack1\\[varin\\.stack1\\.length-1]$")
	varAddRegex                = regexp.MustCompile("=([A-Za-z]{2})\\+([A-Za-z]{2})$")
	stackAddRegex              = regexp.MustCompile("=varin\\.stack1\\[varin\\.stack1\\.length-2]\\+varin\\.stack1\\[varin\\.stack1\\.length-1]$")
	varSubtractRegex           = regexp.MustCompile("=([A-Za-z]{2})-([A-Za-z]{2})$")
	stackSubtractRegex         = regexp.MustCompile("=varin\\.stack1\\[varin\\.stack1\\.length-2]-varin\\.stack1\\[varin\\.stack1\\.length-1]$")
	stackDivisionRegex         = regexp.MustCompile("=varin\\.stack1\\[varin\\.stack1\\.length-2]/varin\\.stack1\\[varin\\.stack1\\.length-1]$")
	stackMultiplyRegex         = regexp.MustCompile("=varin\\.stack1\\[varin\\.stack1\\.length-2]\\*varin\\.stack1\\[varin\\.stack1\\.length-1]$")
	varLeftShiftRegex          = regexp.MustCompile("=([A-Za-z]{2})<<([A-Za-z]{2})$")
	stackLeftShiftRegex        = regexp.MustCompile("=varin\\.stack1\\[varin\\.stack1\\.length-2]<<varin\\.stack1\\[varin\\.stack1\\.length-1]$")
	varRightShiftRegex         = regexp.MustCompile("=([A-Za-z]{2})>>>([A-Za-z]{2})$")
	stackRightShiftTwoRegex    = regexp.MustCompile("=varin\\.stack1\\[varin\\.stack1\\.length-2]>>varin\\.stack1\\[varin\\.stack1\\.length-1]$")
	stackRightShiftThreeRegex  = regexp.MustCompile("=varin\\.stack1\\[varin\\.stack1\\.length-2]>>>varin\\.stack1\\[varin\\.stack1\\.length-1]$")
	varAndRegex                = regexp.MustCompile("=([A-Za-z]{2})&([A-Za-z]{2})$")
	stackAndRegex              = regexp.MustCompile("=varin\\.stack1\\[varin\\.stack1\\.length-2]&varin\\.stack1\\[varin\\.stack1\\.length-1]$")
	stackXorRegex              = regexp.MustCompile("=varin\\.stack1\\[varin\\.stack1\\.length-2]\\^varin\\.stack1\\[varin\\.stack1\\.length-1]$")
	varEqualRegex              = regexp.MustCompile("=([A-Za-z]{2})==([A-Za-z]{2})$")
	stackEqualRegex            = regexp.MustCompile("=varin\\.stack1\\[varin\\.stack1\\.length-2]==varin\\.stack1\\[varin\\.stack1\\.length-1]$")
	stackTripleEqualRegex      = regexp.MustCompile("=varin\\.stack1\\[varin\\.stack1\\.length-2]===varin\\.stack1\\[varin\\.stack1\\.length-1]$")
	trueGetRegex               = regexp.MustCompile("=true$")
	falseGetRegex              = regexp.MustCompile("=false$")
	functionRunRegex           = regexp.MustCompile("=([A-Za-z]{2})\\(([^)]*?)\\)")
	greaterEqualVarBoolRegex   = regexp.MustCompile("=([A-Za-z]{2})>=([A-Za-z]{2})$")
	lessEqualVarBoolRegex      = regexp.MustCompile("=([A-Za-z]{2})<=([A-Za-z]{2})$")
	lessStackBoolRegex         = regexp.MustCompile("=varin\\.stack1\\[varin\\.stack1\\.length-2]<varin\\.stack1\\[varin\\.stack1\\.length-1]$")
	greaterStackBoolRegex      = regexp.MustCompile("=varin\\.stack1\\[varin\\.stack1\\.length-2]>varin\\.stack1\\[varin\\.stack1\\.length-1]$")
	lessvarBoolRegex           = regexp.MustCompile("=([A-Za-z]{2})<([A-Za-z]{2})$")
	stringCreatedRegex         = regexp.MustCompile("=[A-Za-z]\\[[A-Za-z]{2}]=([A-Za-z]{2})$")
	stack3PushRegex            = regexp.MustCompile("varin\\.stack3\\.push\\(\\{\\w+:([A-Za-z]{2}),Ycoord:([A-Za-z]{2}),\\w+:\\d+}\\)")

	rawStackMatch          = regexp.MustCompile("(\\W|^)varin\\.stack1\\[varin\\.stack1\\.length-(\\d+)]")
	rawVarMatch            = regexp.MustCompile("(\\W|^)([A-Za-z]{2})(\\W|$)")
	varParenthesisAddition = regexp.MustCompile("\\((\\d+)\\+([A-Za-z]{2})\\)")
	stringFromArrayIndex   = regexp.MustCompile("=globalStringArray\\[([A-Za-z]{2})]")
	stringSetRegex         = regexp.MustCompile("([A-Za-z]{2})\\+=String")

	routeMatch         = regexp.MustCompile("walk_\\d+route\\d!")
	keyMatchRegex      = regexp.MustCompile("(?:-)?globalKeyArray\\[(\\d+)]")
	dateheapregex      = regexp.MustCompile("Number\\(new heapin\\[\\d+]\\[Date]\\)")
	numbernewheapregex = regexp.MustCompile("Number\\(new heapin\\[\\d+]\\)-0")

	testregex1          = regexp.MustCompile("\\(heapin\\[\\d+]\\[global]\\)")
	lessThanRoundsRegex = regexp.MustCompile("[^<]< ?(\\d+)")
	sixteenKeyRegex     = regexp.MustCompile("\\((\\w{16}),")

	baseMainMatchRegex   = regexp.MustCompile("init\\((.*?),document")
	basehashmatch        = regexp.MustCompile("seed=([^&]*?)&")
	dheadermatch         = regexp.MustCompile("value:\"([^\"]+)\"")
	windowdatematch      = regexp.MustCompile("new \\w\\[Date]")
	P2SEPERATORMATCH     = regexp.MustCompile("=(\\d+)\\^")
	P4MATCH              = regexp.MustCompile("selector=(\\d+);")
	P15IFMATCH           = regexp.MustCompile("if\\((.*?);\\)")
	P15WHILEMATCH        = regexp.MustCompile("while\\(([^;]*?);\\)\\{[^{]*?}")
	P17HEAPSETMATCH      = regexp.MustCompile("selector=(\\d+);")
	P17DEBUGMATCH        = regexp.MustCompile("heap\\[(\\d+)]\\.selector=heapin\\[\\d+]\\[\"debugInfo\"]")
	P17PARAMMATCH        = regexp.MustCompile("heap\\[(\\d+)]\\.selector=heapin\\[\\d+]\\[\"params\"]")
	P19HEAPSETMATCH      = regexp.MustCompile("\\(heapin\\[\\d+]\\[heapin\\[\\d+]],")
	P23XORMATCH          = regexp.MustCompile("=(\\d+)\\^varin")
	P25SHIFTMATCH        = regexp.MustCompile("<< 8")
	P31VALMATCH          = regexp.MustCompile("selector=(\\d+);")
	P32ITEMMATCH         = regexp.MustCompile("heap\\[(\\d+)]\\.selector=Object\\.apply\\.call\\(heapin\\[\\d+]\\[getItem]")
	P32CAPACITYMATCH     = regexp.MustCompile("heap\\[(\\d+)]\\.selector=heapin\\[\\d+]\\[\"NO_CAPACITY\"]")
	P36URLMATCH          = regexp.MustCompile("heap\\[(\\d+)]\\.selector=Object\\.apply\\.call\\(heapin\\[\\d+]\\[slice]")
	P36DATAMATCH         = regexp.MustCompile("heap\\[(\\d+)]\\.selector=varin\\.stack1\\[0]")
	P40DOCUMENTMATCH     = regexp.MustCompile("heap\\[(\\d+)]\\.selector=heapin\\[\\d+]\\[document]")
	P40DOCUMENTBODYMATCH = regexp.MustCompile("heap\\[(\\d+)]\\.selector=heapin\\[\\d+]\\[documentBody]")
	P40GLOBALMATCH       = regexp.MustCompile("heap\\[(\\d+)]\\.selector=heapin\\[\\d+]\\[global]")
	P40NAVIGATORMATCH    = regexp.MustCompile("heap\\[(\\d+)]\\.selector=heapin\\[\\d+]\\[\"navigator\"]")
	P40CRYPTOMATCH       = regexp.MustCompile("heap\\[(\\d+)]\\.selector=heapin\\[\\d+]\\[\"crypto\"]")
	P40EXTERNALMATCH     = regexp.MustCompile("heap\\[(\\d+)]\\.selector=heapin\\[\\d+]\\[\"external\"]")
	P42DPIMATCH          = regexp.MustCompile("heap\\[(\\d+)]\\.selector=heapin\\[\\d+];")
	P42IND0MATCH         = regexp.MustCompile("heapin\\[\\d+]\\[0]")
	P42IND1MATCH         = regexp.MustCompile("heapin\\[\\d+]\\[1]")
	P42HEAPSETMATCH      = regexp.MustCompile("heap\\[\\d+]\\.selector=([1-9](?:\\d+)?);")
	P44TOARRAYMATCH      = regexp.MustCompile("\\[\"toArray\"]\\(\\);varin\\.heap\\[(\\d+)]")

	apixelmatch = regexp.MustCompile("heap\\[(\\d+)]\\.selector=heapin\\[\\d+]\\[\"a\"]")

	ifsliceregex = regexp.MustCompile("if\\(.*?\\)\\{")
	xorkeymatch  = regexp.MustCompile("\\^(\\d{5,})\\)")
)

type GlobalHolder struct {
	Stack3Holder        []Jumper          `json:"stack_3_holder"`
	RanKeys             map[int]bool      `json:"ran_keys"`
	Negatives           map[int]bool      `json:"negatives"`
	Encrounds           int               `json:"encrounds"`
	PayloadOrder        []int             `json:"payload_order"`
	IteratedFiles       map[string]bool   `json:"iterated_files"`
	Routecount          int               `json:"routecount"`
	KeyCount            map[string]int    `json:"key_count"`
	KeyHolderRaw        map[int][]string  `json:"key_holder_raw"`
	KeyHolder           map[int]int64     `json:"key_holder"`
	EncryptionKeyHolder []int             `json:"encryption_key_holder"`
	Tagholder           map[string]string `json:"tagholder"`
	Initializerholder   []FuncInitializer `json:"initializerholder"`
	GlobalStringArray   []string          `json:"global_string_array"`
	GlobalTupleArray    [][][]int         `json:"global_tuple_array"`
	GlobalKeyArray      []int64           `json:"global_key_array"`
	GlobalHeapCopyArray []InitVarObj      `json:"global_heap_copy_array"`
	GlobalOpArray       []Op              `json:"global_op_array"`
	MainNumArray        []uint8           `json:"main_num_array"`
	ObjDefCount         int               `json:"obj_def_count"`
	ArrDefCount         int               `json:"arr_def_count"`
	StrDefCount         int               `json:"str_def_count"`
	IdentifierMap       map[string]string `json:"identifier_map"`
	Basekeys            []int             `json:"basekeys"`
	Seedstring          string            `json:"seedstring"`
	Hashstring          string            `json:"hashstring"`
	Alphabet            string            `json:"alphabet"`
	Base                string            `json:"base"`
	Dheader             string            `json:"dheader"`
	BaseKeys            []int             `json:"base_keys"`
	IV1                 int               `json:"iv_1"`
	IV2                 int               `json:"iv_2"`
	Calckey             int               `json:"calckey"`
	P1VAL               []int             `json:"p_1_val"`
	P2VAL               []int             `json:"p_2_val"`
	P3VAL               []int             `json:"p_3_val"`
	P4VAL               int               `json:"p_4_vval"`
	P5VAL               int               `json:"p_5_vval"`
	P11VAL              []int             `json:"p_11_val"`
	P13VAL              []int             `json:"p_13_val"`
	P14VAL              []int             `json:"p_14_val"`
	P15VAL              []int             `json:"p_15_val"`
	P16VAL              []int             `json:"p_16_val"`
	P17VAL              [][]int             `json:"p_17_val"`
	P18VAL              []int             `json:"p_18_vval"`
	P19VAL              []int             `json:"p_19_val"`
	P20VAL              []int             `json:"p_20_val"`
	P21VAL              []int             `json:"p_21_vval"`
	P23VAL              []int             `json:"p_23_val"`
	P24VAL              []int             `json:"p_24_val"`
	P25VAL              []int             `json:"p_25_val"`
	P30VAL              []int             `json:"p_30_val"`
	P31VAL              []int             `json:"p_31_val"`
	P32VAL              []int             `json:"p_32_val"`
	P35VAL              []int             `json:"p_35_val"`
	P36VAL              []int             `json:"p_38_val"`
	P40VAL              []int             `json:"p_40_val"`
	P42VAL              []int             `json:"p_42_val"`
	P44VAL              []int             `json:"p_44_val"`
}

type Jumper struct {
	x int
	y int
}
type Op struct {
	Type        int      `json:"type"`
	Block1Lines []string `json:"block_1_lines"`
	Block2Lines []string `json:"block_2_lines"`
	Conditional string   `json:"conditional"`
	CondBody    []string `json:"cond_body"`
	Path2       bool     `json:"path_2"`
}
type FuncInitializer struct {
	X         int      `json:"x"`
	Y         int      `json:"y"`
	CopyIndex int      `json:"copy_index"`
	Heap      []string `json:"heap"`
}
type InitVarObj struct {
	LocalObjInitIndexes    []int `json:"local_obj_init_indexes"`
	ArrayCopyIndexes       []int `json:"array_copy_indexes"`
	ArgumentsIndexMapArray []int `json:"arguments_index_map_array"`
	RunSelfReference       int   `json:"run_self_reference"`
	MainSelfReference      int   `json:"main_self_reference"`
}

func elseSlicer(stringin string) string {
	var outstr strings.Builder
	ifmatch := ifsliceregex.FindStringSubmatch(stringin)
	lock := false
	if len(ifmatch) > 0 {
		startindex := strings.Index(stringin, ifmatch[0])
		outstr.WriteString(stringin[0:startindex])
		startindex += len(ifmatch[0])
		breakindex := startindex
		breaks := 1
		for breaks > 0 {
			switch string(stringin[breakindex]) {
			case "\"":
				lock = !lock
				breakindex++
			case "{":
				if !lock {
					breaks++
				}
				breakindex++
			case "}":
				if !lock {
					breaks--
				}
				breakindex++
			default:
				breakindex++
			}
		}
		breaks = 1
		newstart := breakindex + 5
		newbreak := newstart
		lock = false
		for breaks > 0 {
			switch string(stringin[newbreak]) {
			case "\"":
				lock = !lock
				newbreak++
			case "{":
				if !lock {
					breaks++
				}
				newbreak++
			case "}":
				if !lock {
					breaks--
				}
				newbreak++
			default:
				newbreak++
			}
		}
		outstr.WriteString(elseSlicer(stringin[newstart : newbreak-1]))
		if newbreak < len(stringin)-1 {
			outstr.WriteString(elseSlicer(stringin[newbreak:]))
		}
	} else {
		return stringin
	}
	return outstr.String()
}
func P1elseSlicer(stringin string) string {
	var outstr strings.Builder
	ifmatch := ifsliceregex.FindStringSubmatch(stringin)
	lock := false
	if len(ifmatch) > 0 {
		startindex := strings.Index(stringin, ifmatch[0])
		outstr.WriteString(stringin[0:startindex])
		startindex += len(ifmatch[0])
		breakindex := startindex
		breaks := 1
		for breaks > 0 {
			switch string(stringin[breakindex]) {
			case "\"":
				lock = !lock
				breakindex++
			case "{":
				if !lock {
					breaks++
				}
				breakindex++
			case "}":
				if !lock {
					breaks--
				}
				breakindex++
			default:
				breakindex++
			}
		}
		for _, match := range P2SEPERATORMATCH.FindAllStringSubmatch(stringin[startindex:breakindex-1], -1) {
			outstr.WriteString(match[0])
		}
		breaks = 1
		newstart := breakindex + 5
		newbreak := newstart
		lock = false
		for breaks > 0 {
			switch string(stringin[newbreak]) {
			case "\"":
				lock = !lock
				newbreak++
			case "{":
				if !lock {
					breaks++
				}
				newbreak++
			case "}":
				if !lock {
					breaks--
				}
				newbreak++
			default:
				newbreak++
			}
		}
		outstr.WriteString(P1elseSlicer(stringin[newstart : newbreak-1]))
		if newbreak < len(stringin)-1 {
			outstr.WriteString(P1elseSlicer(stringin[newbreak:]))
		}
	} else {
		return stringin
	}
	return outstr.String()
}
func replaceAllOccurrences(arr []int) []int {
	nums := make(map[int]bool)
	for _, entry := range arr {
		var timesFound int
		for _, num := range arr {
			if num == entry {
				timesFound++
			}
		}
		if timesFound > 1 {
			nums[entry] = true
		}
	}
	for k, v := range nums {
		if v {
			arr = replaceOccurence(k, arr)
		}
	}
	return arr
}
func replaceOccurence(number int, arr []int) []int {
	var newarr []int
	index := findLastOccurenceIndex(number, arr)
	for i, num := range arr {
		if i < index && num != number {
			newarr = append(newarr, num)
		} else if i >= index {
			newarr = append(newarr, num)
		}
	}
	return newarr
}
func findLastOccurenceIndex(number int, arr []int) int {
	var lastindex int
	for i, num := range arr {
		if num == number {
			lastindex = i
		}
	}
	return lastindex
}
func (g *GlobalHolder) DeobGlobals(stringin string) {
	stringarrmatch := globalStringArrayRegex.FindStringSubmatch(stringin)
	for _, stringval := range strings.Split(stringarrmatch[2], "\",\"") {
		appstr, _ := strconv.Unquote("\"" + stringval + "\"")
		g.GlobalStringArray = append(g.GlobalStringArray, appstr)
	}
	stringin = strings.ReplaceAll(stringin, fmt.Sprintf("=%s[", stringarrmatch[1]), "=globalStringArray[")
	tuplearrmatch := globalTupleArrayRegex.FindStringSubmatch(stringin)
	for _, tuplearr := range strings.Split(tuplearrmatch[2], "]],[[") {
		appslice := [][]int{}
		for _, tuplepair := range strings.Split(tuplearr, "],[") {
			pair := strings.Split(tuplepair, ",")
			x, _ := strconv.Atoi(pair[0])
			y, _ := strconv.Atoi(pair[1])
			apppair := []int{x, y}
			appslice = append(appslice, apppair)
		}
		g.GlobalTupleArray = append(g.GlobalTupleArray, appslice)
	}
	stringin = strings.ReplaceAll(stringin, fmt.Sprintf("=%s[", tuplearrmatch[1]), "=globalTupleArray[")
	keyarrmatch := globalKeyArrayRegex.FindStringSubmatch(stringin)
	for i, key := range strings.Split(keyarrmatch[2], ",") {
		var digitval int64
		if strings.Contains(key, "e") {
			blocks := strings.Split(key, "e")
			base, _ := strconv.Atoi(blocks[0])
			endian, _ := strconv.Atoi(blocks[1])
			digitval = int64(base * int(math.Pow(10.0, float64(endian))))
		} else
		{
			digitval, _ = strconv.ParseInt(key, 0, 64)
		}
		g.GlobalKeyArray = append(g.GlobalKeyArray, digitval)
		g.RanKeys[i] = true
	}
	stringin = strings.ReplaceAll(stringin, fmt.Sprintf("=%s[", keyarrmatch[1]), "=globalKeyArray[")
	for i, identifier := range runFuncCopyIdentifiersRegex.FindStringSubmatch(stringin)[1:] {
		switch i {
		case 0:
			g.IdentifierMap[identifier] = "localObjInitIndexes"
		case 1:
			g.IdentifierMap[identifier] = "arrayCopyIndexes"
		case 2:
			g.IdentifierMap[identifier] = "argumentsIndexMapArray"
		case 3:
			g.IdentifierMap[identifier] = "runSelfReference"
		case 4:
			g.IdentifierMap[identifier] = "mainSelfReference"
		}
	}
	heapcopyarrmatch := globalHeapCopyRegex.FindStringSubmatch(stringin)
	for _, copyobject := range strings.Split(heapcopyarrmatch[2], "},{") {
		appobj := InitVarObj{
			LocalObjInitIndexes:    nil,
			ArrayCopyIndexes:       nil,
			ArgumentsIndexMapArray: nil,
			RunSelfReference:       -1,
			MainSelfReference:      -1,
		}
		for _, instantiation := range tupleObjectRegex.FindAllStringSubmatch(copyobject, -1) {
			switch g.IdentifierMap[instantiation[1]] {
			case "localObjInitIndexes":
				appslice := []int{}
				for _, val := range strings.Split(instantiation[2], ",") {
					if len(val) != 0 {
						appval, _ := strconv.Atoi(val)
						appslice = append(appslice, appval)
					}
				}
				appobj.LocalObjInitIndexes = appslice
			case "arrayCopyIndexes":
				appslice := []int{}
				for _, val := range strings.Split(instantiation[2], ",") {
					if len(val) != 0 {
						appval, _ := strconv.Atoi(val)
						appslice = append(appslice, appval)
					}
				}
				appobj.ArrayCopyIndexes = appslice
			case "argumentsIndexMapArray":
				appslice := []int{}
				for _, val := range strings.Split(instantiation[2], ",") {
					if len(val) != 0 {
						appval, _ := strconv.Atoi(val)
						appslice = append(appslice, appval)
					}
				}
				if len(appslice) == 0 {
					appslice = nil
				}
				appobj.ArgumentsIndexMapArray = appslice
			}
		}
		for _, instantiation := range tupleObjectIndexReferenceRegex.FindAllStringSubmatch(copyobject, -1) {
			switch g.IdentifierMap[instantiation[1]] {
			case "runSelfReference":
				initval, _ := strconv.Atoi(instantiation[2])
				appobj.RunSelfReference = initval
			case "mainSelfReference":
				initval, _ := strconv.Atoi(instantiation[2])
				appobj.MainSelfReference = initval
			}
		}
		g.GlobalHeapCopyArray = append(g.GlobalHeapCopyArray, appobj)
	}
	stringin = strings.ReplaceAll(stringin, fmt.Sprintf("=%s[", heapcopyarrmatch[1]), "=globalHeapCopyArray[")
	funcarrmatch := funcArrayRegex.FindStringSubmatch(stringin)
	stringin = strings.ReplaceAll(stringin, fmt.Sprintf("(%s)", funcarrmatch[1]), "(funcArray)")
	for i, val := range functionDefsRegex.FindStringSubmatch(stringin)[1:] {
		replacementmatch := regexp.MustCompile(fmt.Sprintf("([^\\w.])%s\\(", val))
		switch i {
		case 0:
			for _, rep := range replacementmatch.FindAllStringSubmatch(stringin, -1) {
				stringin = strings.ReplaceAll(stringin, rep[0], rep[1]+"ReferenceError(")
			}
		case 1:
			for _, rep := range replacementmatch.FindAllStringSubmatch(stringin, -1) {
				stringin = strings.ReplaceAll(stringin, rep[0], rep[1]+"TypeError(")
			}
		case 2:
			for _, rep := range replacementmatch.FindAllStringSubmatch(stringin, -1) {
				stringin = strings.ReplaceAll(stringin, rep[0], rep[1]+"Object(")
			}
		case 3:
			for _, rep := range replacementmatch.FindAllStringSubmatch(stringin, -1) {
				stringin = strings.ReplaceAll(stringin, rep[0], rep[1]+"RegExp(")
			}
		case 4:
			for _, rep := range replacementmatch.FindAllStringSubmatch(stringin, -1) {
				stringin = strings.ReplaceAll(stringin, rep[0], rep[1]+"Number(")
			}
		case 5:
			for _, rep := range replacementmatch.FindAllStringSubmatch(stringin, -1) {
				stringin = strings.ReplaceAll(stringin, rep[0], rep[1]+"String(")
			}
		case 6:
			for _, rep := range replacementmatch.FindAllStringSubmatch(stringin, -1) {
				stringin = strings.ReplaceAll(stringin, rep[0], rep[1]+"Array(")
			}
		case 7:
			for _, rep := range replacementmatch.FindAllStringSubmatch(stringin, -1) {
				stringin = strings.ReplaceAll(stringin, rep[0], rep[1]+"Object.bind(")
			}
		case 8:
			for _, rep := range replacementmatch.FindAllStringSubmatch(stringin, -1) {
				stringin = strings.ReplaceAll(stringin, rep[0], rep[1]+"Object.call(")
			}
		case 9:
			for _, rep := range replacementmatch.FindAllStringSubmatch(stringin, -1) {
				stringin = strings.ReplaceAll(stringin, rep[0], rep[1]+"Object.call.bind(Object.bind, Object.call)(")
			}
		case 10:
			for _, rep := range replacementmatch.FindAllStringSubmatch(stringin, -1) {
				stringin = strings.ReplaceAll(stringin, rep[0], rep[1]+"Object.apply(")
			}
		case 11:
			for _, rep := range replacementmatch.FindAllStringSubmatch(stringin, -1) {
				stringin = strings.ReplaceAll(stringin, rep[0], rep[1]+"Object.apply.call(")
			}
		case 12:
			for _, rep := range replacementmatch.FindAllStringSubmatch(stringin, -1) {
				stringin = strings.ReplaceAll(stringin, rep[0], rep[1]+"[].push(")
			}
			stringin = strings.Replace(stringin, fmt.Sprintf(":%s}", val), ":[].push}", 1)
		case 13:
			for _, rep := range replacementmatch.FindAllStringSubmatch(stringin, -1) {
				stringin = strings.ReplaceAll(stringin, rep[0], rep[1]+"[].pop(")
			}
			stringin = strings.Replace(stringin, fmt.Sprintf(":%s}", val), ":[].pop}", 1)
		case 14:
			for _, rep := range replacementmatch.FindAllStringSubmatch(stringin, -1) {
				stringin = strings.ReplaceAll(stringin, rep[0], rep[1]+"[].slice(")
			}
			stringin = strings.Replace(stringin, fmt.Sprintf(":%s}", val), ":[].slice}", 1)
		case 15:
			for _, rep := range replacementmatch.FindAllStringSubmatch(stringin, -1) {
				stringin = strings.ReplaceAll(stringin, rep[0], rep[1]+"[].splice(")
			}
			stringin = strings.Replace(stringin, fmt.Sprintf(":%s}", val), ":[].splice}", 1)
		case 16:
			for _, rep := range replacementmatch.FindAllStringSubmatch(stringin, -1) {
				stringin = strings.ReplaceAll(stringin, rep[0], rep[1]+"[].join(")
			}
		case 17:
			for _, rep := range replacementmatch.FindAllStringSubmatch(stringin, -1) {
				stringin = strings.ReplaceAll(stringin, rep[0], rep[1]+"[].map(")
			}
		case 18:
			for _, rep := range replacementmatch.FindAllStringSubmatch(stringin, -1) {
				stringin = strings.ReplaceAll(stringin, rep[0], rep[1]+"[].push.call(")
			}
		case 19:
			for _, rep := range replacementmatch.FindAllStringSubmatch(stringin, -1) {
				stringin = strings.ReplaceAll(stringin, rep[0], rep[1]+"[].slice.call(")
			}
		case 20:
			for _, rep := range replacementmatch.FindAllStringSubmatch(stringin, -1) {
				stringin = strings.ReplaceAll(stringin, rep[0], rep[1]+"[].join.call(")
			}
		case 21:
			for _, rep := range replacementmatch.FindAllStringSubmatch(stringin, -1) {
				stringin = strings.ReplaceAll(stringin, rep[0], rep[1]+"[].map.call(")
			}
		case 22:
			for _, rep := range replacementmatch.FindAllStringSubmatch(stringin, -1) {
				stringin = strings.ReplaceAll(stringin, rep[0], rep[1]+"{}.hasOwnProperty(")
			}
		case 23:
			for _, rep := range replacementmatch.FindAllStringSubmatch(stringin, -1) {
				stringin = strings.ReplaceAll(stringin, rep[0], rep[1]+"{}.hasOwnProperty.call(")
			}
		case 24:
			for _, rep := range replacementmatch.FindAllStringSubmatch(stringin, -1) {
				stringin = strings.ReplaceAll(stringin, rep[0], rep[1]+"JSON.stringify(")
			}
		case 25:
			for _, rep := range replacementmatch.FindAllStringSubmatch(stringin, -1) {
				stringin = strings.ReplaceAll(stringin, rep[0], rep[1]+"Object.getOwnPropertyDescriptor(")
			}
		case 26:
			for _, rep := range replacementmatch.FindAllStringSubmatch(stringin, -1) {
				stringin = strings.ReplaceAll(stringin, rep[0], rep[1]+"Object.defineProperty(")
			}
		case 27:
			for _, rep := range replacementmatch.FindAllStringSubmatch(stringin, -1) {
				stringin = strings.ReplaceAll(stringin, rep[0], rep[1]+"String.fromCharCode(")
			}
		case 28:
			for _, rep := range replacementmatch.FindAllStringSubmatch(stringin, -1) {
				stringin = strings.ReplaceAll(stringin, rep[0], rep[1]+"Math.min(")
			}
		case 29:
			for _, rep := range replacementmatch.FindAllStringSubmatch(stringin, -1) {
				stringin = strings.ReplaceAll(stringin, rep[0], rep[1]+"Math.floor(")
			}
		case 30:
			for _, rep := range replacementmatch.FindAllStringSubmatch(stringin, -1) {
				stringin = strings.ReplaceAll(stringin, rep[0], rep[1]+"Object.create(")
			}
		case 31:
			for _, rep := range replacementmatch.FindAllStringSubmatch(stringin, -1) {
				stringin = strings.ReplaceAll(stringin, rep[0], rep[1]+"\"\".indexOf(")
			}
		case 32:
			for _, rep := range replacementmatch.FindAllStringSubmatch(stringin, -1) {
				stringin = strings.ReplaceAll(stringin, rep[0], rep[1]+"\"\".charAt(")
			}
		case 33:
			for _, rep := range replacementmatch.FindAllStringSubmatch(stringin, -1) {
				stringin = strings.ReplaceAll(stringin, rep[0], rep[1]+"\"\".indexOf.call(")
			}
		case 34:
			for _, rep := range replacementmatch.FindAllStringSubmatch(stringin, -1) {
				stringin = strings.ReplaceAll(stringin, rep[0], rep[1]+"\"\".charAt.call(")
			}
		case 35:
			for _, rep := range replacementmatch.FindAllStringSubmatch(stringin, -1) {
				stringin = strings.ReplaceAll(stringin, rep[0], rep[1]+"Uint8Array(")
			}
		}
	}
	functiondefsmatch := globalOpArrayRegex.FindAllStringSubmatch(stringin, -1)[1][1]
	xymatch := XYCordRegex.FindStringSubmatch(stringin)
	xreplaceregex := regexp.MustCompile(fmt.Sprintf("\\.%s(\\W)", xymatch[3]))
	for _, replacement := range xreplaceregex.FindAllStringSubmatch(functiondefsmatch, -1) {
		functiondefsmatch = strings.ReplaceAll(functiondefsmatch, replacement[0], ".Xcoord"+replacement[1])
	}
	yreplaceregex := regexp.MustCompile(fmt.Sprintf("\\.%s(\\W)", xymatch[1]))
	for _, replacement := range yreplaceregex.FindAllStringSubmatch(functiondefsmatch, -1) {
		functiondefsmatch = strings.ReplaceAll(functiondefsmatch, replacement[0], ".Ycoord"+replacement[1])
	}
	xjumpsetterregex := regexp.MustCompile(fmt.Sprintf("(\\W)%s:", xymatch[3]))
	for _, replacement := range xjumpsetterregex.FindAllStringSubmatch(functiondefsmatch, -1) {
		functiondefsmatch = strings.ReplaceAll(functiondefsmatch, replacement[0], replacement[1]+"Xcoord:")
	}
	yjumpsetterregex := regexp.MustCompile(fmt.Sprintf("(\\W)%s:", xymatch[1]))
	for _, replacement := range yjumpsetterregex.FindAllStringSubmatch(functiondefsmatch, -1) {
		functiondefsmatch = strings.ReplaceAll(functiondefsmatch, replacement[0], replacement[1]+"Ycoord:")
	}
	g.MainNumArray = g.StringDecode(regexp.MustCompile(fmt.Sprintf("var %s=\\w\\(\"(.*?)\"\\)", xymatch[2])).FindStringSubmatch(stringin)[1])
	functiondefsmatch = strings.ReplaceAll(functiondefsmatch, fmt.Sprintf("[%s[", xymatch[2]), "[mainnumarr[")
	functiondefsmatch = strings.ReplaceAll(functiondefsmatch, fmt.Sprintf("=%s[", xymatch[2]), "=mainnumarr[")
	functiondefsmatch = strings.ReplaceAll(functiondefsmatch, fmt.Sprintf("(%s[", xymatch[2]), "(mainnumarr[")
	functiondefsmatch = strings.ReplaceAll(functiondefsmatch, fmt.Sprintf("|%s[", xymatch[2]), "|mainnumarr[")
	runfuncvar := runFuncDeclarationRegex.FindStringSubmatch(stringin)[1]
	functiondefsmatch = strings.ReplaceAll(functiondefsmatch, fmt.Sprintf("=%s(", runfuncvar), "=runFuncDeclaration(")
	for i, val := range objectDefinitionRegex.FindAllStringSubmatch(stringin, -1) {
		switch i {
		case 2: //getter
			for _, replacement := range regexp.MustCompile(fmt.Sprintf("\\.%s\\((\\w{2})\\)", val)).FindAllStringSubmatch(functiondefsmatch, -1) {
				functiondefsmatch = strings.ReplaceAll(functiondefsmatch, replacement[0], "["+replacement[1]+"].selector")
			}
		case 3: //setter
			for _, replacement := range regexp.MustCompile(fmt.Sprintf("\\.%s\\(([^(),]+),([^(),]+)\\)", val)).FindAllStringSubmatch(functiondefsmatch, -1) {
				functiondefsmatch = strings.ReplaceAll(functiondefsmatch, replacement[0], "["+replacement[1]+"].selector="+replacement[2])
			}
		case 5: //pop
			for _, replacement := range regexp.MustCompile(fmt.Sprintf("\\.%s\\(\\)", val)).FindAllStringSubmatch(functiondefsmatch, -1) {
				functiondefsmatch = strings.ReplaceAll(functiondefsmatch, replacement[0], ".pop()")
			}
		case 6: //push
			for _, replacement := range regexp.MustCompile(fmt.Sprintf("\\.%s\\((.*?)\\)", val)).FindAllStringSubmatch(functiondefsmatch, -1) {
				functiondefsmatch = strings.ReplaceAll(functiondefsmatch, replacement[0], ".push("+replacement[1]+")")
			}
		case 7: //slice
			for _, replacement := range regexp.MustCompile(fmt.Sprintf("\\.%s\\((.*?)\\)", val)).FindAllStringSubmatch(functiondefsmatch, -1) {
				functiondefsmatch = strings.ReplaceAll(functiondefsmatch, replacement[0], ".slice("+replacement[1]+")")
			}
		default:
			continue
		}
	}
	for i, val := range funcObjectRegex.FindStringSubmatch(stringin)[1:] {
		replacementmatch := regexp.MustCompile(fmt.Sprintf("\\.%s(\\W)", val)).FindAllStringSubmatch(functiondefsmatch, -1)
		switch i {
		case 0:
			for _, replacement := range replacementmatch {
				functiondefsmatch = strings.ReplaceAll(functiondefsmatch, replacement[0], ".stack1"+replacement[1])
			}
		case 1:
			for _, replacement := range replacementmatch {
				functiondefsmatch = strings.ReplaceAll(functiondefsmatch, replacement[0], ".stack2"+replacement[1])
			}
		case 2:
			for _, replacement := range replacementmatch {
				functiondefsmatch = strings.ReplaceAll(functiondefsmatch, replacement[0], ".stack3"+replacement[1])
			}
		case 6:
			for _, replacement := range replacementmatch {
				functiondefsmatch = strings.ReplaceAll(functiondefsmatch, replacement[0], ".heap"+replacement[1])
			}
		case 7:
			for _, replacement := range replacementmatch {
				functiondefsmatch = strings.ReplaceAll(functiondefsmatch, replacement[0], ".thisWindowOrCopy"+replacement[1])
			}
		case 8:
			for _, replacement := range replacementmatch {
				functiondefsmatch = strings.ReplaceAll(functiondefsmatch, replacement[0], ".thisRaw"+replacement[1])
			}
		case 9:
			for _, replacement := range replacementmatch {
				functiondefsmatch = strings.ReplaceAll(functiondefsmatch, replacement[0], ".explicitJumpHolder"+replacement[1])
			}
		}
	}
	for _, val := range regexp.MustCompile(fmt.Sprintf("(\\W)%s(\\W)", breakObjRegex.FindStringSubmatch(stringin)[1])).FindAllStringSubmatch(functiondefsmatch, -1) {
		functiondefsmatch = strings.Replace(functiondefsmatch, val[0], val[1]+"breakobject"+val[2], 1)
	}
	for _, funcdef := range strings.Split(functiondefsmatch, ",function") {
		funcvar := opVarRegex.FindStringSubmatch(funcdef)[1]
		funcbody := funcdef[strings.Index(funcdef, "{")+1 : strings.LastIndex(funcdef, "}")]
		repregex := regexp.MustCompile(fmt.Sprintf("(^|\\W)%s(\\W|$)", funcvar))
		for _, replacementmatch := range repregex.FindAllStringSubmatch(funcbody, -1) {
			funcbody = strings.ReplaceAll(funcbody, replacementmatch[0], replacementmatch[1]+"varin"+replacementmatch[2])
		}
		appobj := Op{}
		ifjump := jumpifregex.FindStringSubmatch(funcbody)
		ifthrow := throwifregex.FindStringSubmatch(funcbody)
		if len(ifjump) != 0 {
			appobj.Type = 2
			block := strings.Split(funcbody, ifjump[0])
			appobj.Block1Lines = strings.Split(block[0], ";")
			appobj.Block2Lines = strings.Split(block[1], ";")
			appobj.Conditional = ifjump[1]
			appobj.CondBody = strings.Split(ifjump[2], ";")
			appobj.Path2 = true
		} else
		if len(ifthrow) != 0 {
			appobj.Type = 2
			block := strings.Split(funcbody, ifthrow[0])
			appobj.Block1Lines = strings.Split(block[0], ";")
			appobj.Block2Lines = strings.Split(block[1], ";")
			appobj.Conditional = ifthrow[1]
			appobj.CondBody = strings.Split(ifthrow[2], ";")
			appobj.Path2 = false
		} else
		if len(breakmatchregex.FindString(funcbody)) != 0 {
			appobj.Type = 0
			forvarin := forvarinregex.FindStringSubmatch(funcbody)
			if len(forvarin) != 0 {
				block := strings.Split(funcbody, forvarin[0])
				appobj.Block1Lines = strings.Split(block[0], ";")
				appobj.Block1Lines = append(appobj.Block1Lines, forvarin[0])
				appobj.Block1Lines = append(appobj.Block1Lines, strings.Split(block[1], ";")...)
			} else {
				blocks := []string{}
				linkers := []string{}
				for len(rawstringdeobregex.FindStringSubmatch(funcbody)) != 0 {
					b0match := b0matchregex.FindStringSubmatch(funcbody)
					b1match := b1matchregex.FindStringSubmatch(funcbody)
					if len(b0match) != 0 {
						bodyblocks := strings.Split(funcbody, b0match[0])
						rawstringdeob := rawstringdeobregex.FindStringSubmatch(b0match[1])
						subblocks := strings.Split(b0match[1], rawstringdeob[0])
						blocks = append(blocks, bodyblocks[0]+subblocks[0])
						linkers = append(linkers, rawstringdeob[0])
						funcbody = subblocks[1] + ";" + bodyblocks[1]
					} else
					if len(b1match) != 0 {
						bodyblocks := strings.Split(funcbody, b1match[0])
						rawstringdeob := rawstringdeobregex.FindStringSubmatch(b1match[1])
						subblocks := strings.Split(b0match[1], rawstringdeob[0])
						blocks = append(blocks, bodyblocks[0]+subblocks[0])
						linkers = append(linkers, rawstringdeob[0])
						funcbody = subblocks[1] + ";" + bodyblocks[1]
					} else {
						rawstringdeob := rawstringdeobregex.FindStringSubmatch(funcbody)
						bodyblocks := strings.Split(funcbody, rawstringdeob[0])
						blocks = append(blocks, bodyblocks[0])
						linkers = append(linkers, rawstringdeob[0])
						funcbody = bodyblocks[1]
					}
				}
				for i, block := range blocks {
					appobj.Block1Lines = append(appobj.Block1Lines, strings.Split(block, ";")...)
					appobj.Block1Lines = append(appobj.Block1Lines, linkers[i])
				}
				appobj.Block1Lines = append(appobj.Block1Lines, strings.Split(funcbody, ";")...)
			}
		} else {
			if len(throwMatchRegex.FindStringSubmatch(funcbody)) != 0 {
				appobj.Type = 0
			} else {
				appobj.Type = 1
			}
			forvarin := forvarinregex.FindStringSubmatch(funcbody)
			if len(forvarin) != 0 {
				block := strings.Split(funcbody, forvarin[0])
				appobj.Block1Lines = strings.Split(block[0], ";")
				appobj.Block1Lines = append(appobj.Block1Lines, forvarin[0])
				appobj.Block1Lines = append(appobj.Block1Lines, strings.Split(block[1], ";")...)
			} else {
				blocks := []string{}
				linkers := []string{}
				for len(rawstringdeobregex.FindStringSubmatch(funcbody)) != 0 {
					b0match := b0matchregex.FindStringSubmatch(funcbody)
					b1match := b1matchregex.FindStringSubmatch(funcbody)
					if len(b0match) != 0 {
						bodyblocks := strings.Split(funcbody, b0match[0])
						rawstringdeob := rawstringdeobregex.FindStringSubmatch(b0match[1])
						subblocks := strings.Split(b0match[1], rawstringdeob[0])
						blocks = append(blocks, bodyblocks[0]+subblocks[0])
						linkers = append(linkers, rawstringdeob[0])
						funcbody = subblocks[1] + ";" + bodyblocks[1]
					} else
					if len(b1match) != 0 {
						bodyblocks := strings.Split(funcbody, b1match[0])
						rawstringdeob := rawstringdeobregex.FindStringSubmatch(b1match[1])
						subblocks := strings.Split(b1match[1], rawstringdeob[0])
						blocks = append(blocks, bodyblocks[0]+subblocks[0])
						linkers = append(linkers, rawstringdeob[0])
						funcbody = subblocks[1] + ";" + bodyblocks[1]
					} else {
						rawstringdeob := rawstringdeobregex.FindStringSubmatch(funcbody)
						bodyblocks := strings.Split(funcbody, rawstringdeob[0])
						blocks = append(blocks, bodyblocks[0])
						linkers = append(linkers, rawstringdeob[0])
						funcbody = bodyblocks[1]
					}
				}
				for i, block := range blocks {
					appobj.Block1Lines = append(appobj.Block1Lines, strings.Split(block, ";")...)
					appobj.Block1Lines = append(appobj.Block1Lines, linkers[i])
				}
				appobj.Block1Lines = append(appobj.Block1Lines, strings.Split(funcbody, ";")...)
			}
		}
		g.GlobalOpArray = append(g.GlobalOpArray, appobj)
	}
}
func (g *GlobalHolder) StringDeobWrapper(string1, string2 string) string {
	stringarr1 := g.StringDecode(string1)
	stringarr2 := g.StringDecode(string2)
	firstbyte := (int(stringarr1[0]) + int(stringarr2[0])) & 255
	var outstr strings.Builder
	for i := 1; i < len(stringarr2); i++ {
		outstr.WriteString(string(rune(int(stringarr1[i]) ^ int(stringarr2[i]) ^ firstbyte)))
	}
	return outstr.String()
}
func (g *GlobalHolder) OpStepper(x, y int) (int, int, int) {
	tupleout := g.GlobalTupleArray[y][g.MainNumArray[x]]
	x++
	return x, tupleout[0], tupleout[1]
}
func (g *GlobalHolder) StringDecode(stringin string) []uint8 {
	runearray := []rune(stringin)
	alphabet := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"
	runelength := len(runearray)
	outarr := make([]uint8, int(math.Floor(float64(runelength)*3/4)))
	for counter1, counter2 := 0, 0; counter1 < runelength; counter1, counter2 = counter1+4, counter2+3 {
		outarr[counter2] = uint8((strings.Index(alphabet, string(runearray[counter1])) << 2) | (strings.Index(alphabet, string(runearray[counter1+1])) >> 4))
		if counter1+2 < runelength {
			outarr[counter2+1] = uint8(((strings.Index(alphabet, string(runearray[counter1+1])))&15)<<4 | (strings.Index(alphabet, string(runearray[counter1+2])) >> 2))
		}
		if counter1+3 < runelength {
			outarr[counter2+2] = uint8(((strings.Index(alphabet, string(runearray[counter1+2])))&3)<<6 | (strings.Index(alphabet, string(runearray[counter1+3]))))
		}
	}
	return outarr
}
func (g *GlobalHolder) parseLine(line string, x, y int, stackin, heapin []string, jumper Jumper, roundVarMap map[string]string) (string, int, int, []string, []string, Jumper, map[string]string) {
	for _, val := range varParenthesisAddition.FindAllStringSubmatch(line, -1) {
		addbase, _ := strconv.Atoi(val[1])
		addvar, _ := strconv.Atoi(roundVarMap[val[2]])
		line = strings.Replace(line, val[0], strconv.Itoa(addbase+addvar), 1)
	}
	XCordIncrement := XCordIncrementRegex.FindStringSubmatch(line)
	lengthDecrement := lengthDecrementRegex.FindStringSubmatch(line)
	varSet := varSetRegex.FindStringSubmatch(line)
	heapSet := heapSetRegex.FindStringSubmatch(line)
	stackSetStack := stackSetStackRegex.FindStringSubmatch(line)
	stackSetVar := stackSetVarRegex.FindStringSubmatch(line)
	objectFromStackSet := objectFromStackSetRegex.FindStringSubmatch(line)
	XCoordSet := XCoordSetRegex.FindStringSubmatch(line)
	YCoordSet := YCoordSetRegex.FindStringSubmatch(line)
	jumpHolderSet := jumpHolderSetRegex.FindStringSubmatch(line)
	stackPush := stackPushRegex.FindStringSubmatch(line)
	stringDeob := rawstringdeobregex.FindStringSubmatch(line)
	stack3Push := stack3PushRegex.FindStringSubmatch(line)
	if len(stringDeob) != 0 {
		deobstring1 := roundVarMap[stringDeob[1]]
		string2index, _ := strconv.Atoi(roundVarMap[stringFromArrayIndex.FindStringSubmatch(line)[1]])
		deobstring2 := g.GlobalStringArray[string2index]
		outstr := "\"" + g.StringDeobWrapper(deobstring1, deobstring2) + "\""
		strvar := stringSetRegex.FindStringSubmatch(line)
		roundVarMap[strvar[1]] = outstr
		line = ""
	} else
	if len(XCordIncrement) != 0 {
		addnum, _ := strconv.Atoi(XCordIncrement[1])
		x += addnum
		line = ""
	} else
	if len(lengthDecrement) != 0 {
		subnum, _ := strconv.Atoi(lengthDecrement[1])
		if len(lengthDecrement) == 3 {
			appnum, _ := strconv.Atoi(roundVarMap[lengthDecrement[2]])
			subnum += appnum
		}
		stackin = stackin[:len(stackin)-subnum]
		line = ""
	} else
	if len(varSet) != 0 {
		globalStringArrGet := globalStringArrGetRegex.FindStringSubmatch(line)
		globalStringArrOneGet := globalStringArrOneGetRegex.FindStringSubmatch(line)
		mainnumarrGet := mainnumarrGetRegex.FindStringSubmatch(line)
		mainnumarrOneGet := mainnumarrOneGetRegex.FindStringSubmatch(line)
		mainnumarrTwoGet := mainnumarrTwoGetRegex.FindStringSubmatch(line)
		mainnumarrThreeGet := mainnumarrThreeGetRegex.FindStringSubmatch(line)
		stackGet := stackGetRegex.FindStringSubmatch(line)
		varGet := varGetRegex.FindStringSubmatch(line)
		runFuncGet := runFuncGetRegex.FindStringSubmatch(line)
		heapGet := heapGetRegex.FindStringSubmatch(line)
		arrCreate := arrCreateRegex.FindStringSubmatch(line)
		objCreate := objCreateRegex.FindStringSubmatch(line)
		strCreate := strCreateRegex.FindStringSubmatch(line)
		varAdd := varAddRegex.FindStringSubmatch(line)
		varSubtract := varSubtractRegex.FindStringSubmatch(line)
		varLeftShift := varLeftShiftRegex.FindStringSubmatch(line)
		varRightShift := varRightShiftRegex.FindStringSubmatch(line)
		varAnd := varAndRegex.FindStringSubmatch(line)
		functionRun := functionRunRegex.FindStringSubmatch(line)
		stringCreated := stringCreatedRegex.FindStringSubmatch(line)
		stackLengthGet := stackLengthGetRegex.FindStringSubmatch(line)
		keyArrayMatch := keyArrayMatchRegex.FindStringSubmatch(line)
		if len(keyArrayMatch) != 0 {
			indexnum := x
			if len(keyArrayMatch) == 2 {
				addnum, _ := strconv.Atoi(keyArrayMatch[1])
				indexnum += addnum
			}
			keyindexnum := g.MainNumArray[indexnum]
			roundVarMap[varSet[1]] = fmt.Sprintf("globalKeyArray[%d]", keyindexnum)
			line = strings.Replace(line, keyArrayMatch[0], "="+fmt.Sprintf("globalKeyArray[%d]", keyindexnum), 1)
			//line = ""
		} else
		if len(globalStringArrGet) != 0 {
			addnum, _ := strconv.Atoi(globalStringArrGet[1])
			stringVal := g.GlobalStringArray[(int(g.MainNumArray[x])<<8)|int(g.MainNumArray[x+addnum])]
			roundVarMap[varSet[1]] = stringVal
			//line = strings.Replace(line, globalStringArrGet[0], "="+stringVal, 1)
			line = ""
		} else
		if len(globalStringArrOneGet) != 0 {
			addnumone, _ := strconv.Atoi(globalStringArrOneGet[1])
			addnumtwo, _ := strconv.Atoi(globalStringArrOneGet[2])
			stringVal := g.GlobalStringArray[(int(g.MainNumArray[x+addnumone])<<8)|int(g.MainNumArray[x+addnumtwo])]
			roundVarMap[varSet[1]] = stringVal
			//line = strings.Replace(line, globalStringArrOneGet[0], "="+stringVal, 1)
			line = ""
		} else
		if len(mainnumarrGet) != 0 {
			stringnum := strconv.Itoa(int(g.MainNumArray[x]))
			roundVarMap[varSet[1]] = stringnum
			//line = strings.Replace(line, mainnumarrGet[0], "="+stringnum, 1)
			line = ""
		} else
		if len(mainnumarrOneGet) != 0 {
			numone, _ := strconv.Atoi(mainnumarrOneGet[1])
			stringnum := strconv.Itoa(int(g.MainNumArray[x+numone]))
			roundVarMap[varSet[1]] = stringnum
			//line = strings.Replace(line, mainnumarrOneGet[0], "="+stringnum, 1)
			line = ""
		} else
		if len(mainnumarrTwoGet) != 0 {
			numone := 0
			if len(mainnumarrTwoGet) == 3 {
				numone, _ = strconv.Atoi(mainnumarrTwoGet[1])
			}
			numtwo, _ := strconv.Atoi(mainnumarrTwoGet[len(mainnumarrTwoGet)-1])
			stringnum := strconv.Itoa((int(g.MainNumArray[x+numone]) << 8) | int(g.MainNumArray[x+numtwo]))
			roundVarMap[varSet[1]] = stringnum
			//line = strings.Replace(line, mainnumarrTwoGet[0], "="+stringnum, 1)
			line = ""
		} else
		if len(mainnumarrThreeGet) != 0 {
			stringnum := strconv.Itoa(int(g.MainNumArray[x])<<16 | ((int(g.MainNumArray[x+1]) << 8) | int(g.MainNumArray[x+2])))
			roundVarMap[varSet[1]] = stringnum
			//line = strings.Replace(line, mainnumarrThreeGet[0], "="+stringnum, 1)
			line = ""
		} else
		if len(stackGet) != 0 {
			subnum, _ := strconv.Atoi(stackGet[1])
			roundVarMap[varSet[1]] = stackin[len(stackin)-subnum]
			//line = strings.Replace(line, stackGet[0], "="+stackin[len(stackin)-subnum], 1)
			line = ""
		} else
		if len(varGet) != 0 {
			roundVarMap[varSet[1]] = roundVarMap[varGet[1]]
			//line = strings.Replace(line, varGet[0], "="+roundVarMap[varGet[1]], 1)
			line = ""
		} else
		if len(runFuncGet) != 0 {
			var argarr []string
			for _, arg := range runFuncGet[1:] {
				rawStack := rawStackMatch.FindStringSubmatch(arg)
				rawVar := rawVarMatch.FindStringSubmatch(arg)
				if len(rawStack) != 0 {
					subnum, _ := strconv.Atoi(rawStack[2])
					argarr = append(argarr, stackin[len(stackin)-subnum])
				} else if len(rawVar) != 0 {
					argarr = append(argarr, roundVarMap[rawVar[2]])
				}
			}
			roundVarMap[varSet[1]] = fmt.Sprintf("func_%s_%s", argarr[1], argarr[2])
			copyindexinit, _ := strconv.Atoi(argarr[0])
			xcordinit, _ := strconv.Atoi(argarr[1])
			ycordinit, _ := strconv.Atoi(argarr[2])
			g.Initializerholder = append(g.Initializerholder, FuncInitializer{
				X:         xcordinit,
				Y:         ycordinit,
				CopyIndex: copyindexinit,
				Heap:      heapin[:],
			})
			line = strings.Replace(line, runFuncGet[0], fmt.Sprintf("=runFuncDeclaration(%s,heap)", strings.Join(argarr, ",")), 1)
		} else
		if len(heapGet) != 0 {
			heapindex, _ := strconv.Atoi(roundVarMap[heapGet[1]])
			//roundVarMap[varSet[1]] = heapin[heapindex]
			roundVarMap[varSet[1]] = fmt.Sprintf("heapin[%d]", heapindex)
			//line = strings.Replace(line, heapGet[0], "="+heapin[heapindex], 1)
			line = ""
		} else
		if len(arrCreate) != 0 {
			roundVarMap[varSet[1]] = fmt.Sprintf("arrobj%d", g.ArrDefCount)
			line = fmt.Sprintf("arrobj%d=[]", g.ArrDefCount)
			g.ArrDefCount++
		} else
		if len(objCreate) != 0 {
			roundVarMap[varSet[1]] = fmt.Sprintf("objobj%d", g.ObjDefCount)
			line = fmt.Sprintf("objobj%d=[]", g.ObjDefCount)
			g.ObjDefCount++
		} else
		if len(strCreate) != 0 {
			roundVarMap[varSet[1]] = fmt.Sprintf("strobj%d", g.StrDefCount)
			line = fmt.Sprintf("strobj%d=[]", g.StrDefCount)
			g.StrDefCount++
		} else
		if len(varAdd) != 0 {
			firstnum, err1 := strconv.Atoi(roundVarMap[varAdd[1]])
			secondnum, err2 := strconv.Atoi(roundVarMap[varAdd[2]])
			if err1 != nil {
				sumstring := roundVarMap[varAdd[1]] + "+" + strconv.Itoa(secondnum)
				roundVarMap[varSet[1]] = sumstring
			} else
			if err2 != nil {
				sumstring := strconv.Itoa(firstnum) + "+" + roundVarMap[varAdd[2]]
				roundVarMap[varSet[1]] = sumstring
			} else
			{
				sumstring := strconv.Itoa(firstnum + secondnum)
				roundVarMap[varSet[1]] = sumstring
			}
			//line = strings.Replace(line, varAdd[0], "="+sumstring, 1)
			line = ""
		} else
		if len(varSubtract) != 0 {
			firstnum, err1 := strconv.Atoi(roundVarMap[varSubtract[1]])
			secondnum, err2 := strconv.Atoi(roundVarMap[varSubtract[2]])
			if err1 != nil {
				sumstring := roundVarMap[varSubtract[1]] + "-" + strconv.Itoa(secondnum)
				roundVarMap[varSet[1]] = sumstring
			} else
			if err2 != nil {
				sumstring := strconv.Itoa(firstnum) + "-" + roundVarMap[varSubtract[2]]
				roundVarMap[varSet[1]] = sumstring
			} else
			{
				sumstring := strconv.Itoa(firstnum - secondnum)
				roundVarMap[varSet[1]] = sumstring
			}
			//line = strings.Replace(line, varSubtract[0], "="+sumstring, 1)
			line = ""
		} else
		if len(varLeftShift) != 0 {
			firstnum, err1 := strconv.Atoi(roundVarMap[varLeftShift[1]])
			secondnum, err2 := strconv.Atoi(roundVarMap[varLeftShift[2]])
			if err1 != nil {
				sumstring := roundVarMap[varLeftShift[1]] + "<<" + strconv.Itoa(secondnum)
				roundVarMap[varSet[1]] = sumstring
			} else
			if err2 != nil {
				sumstring := strconv.Itoa(firstnum) + "<<" + roundVarMap[varLeftShift[2]]
				roundVarMap[varSet[1]] = sumstring
			} else
			{
				sumstring := strconv.Itoa(firstnum << secondnum)
				roundVarMap[varSet[1]] = sumstring
			}
			//line = strings.Replace(line, varLeftShift[0], "="+sumstring, 1)
			line = ""
		} else
		if len(varRightShift) != 0 {
			firstnum, err1 := strconv.Atoi(roundVarMap[varRightShift[1]])
			secondnum, err2 := strconv.Atoi(roundVarMap[varRightShift[2]])
			if err1 != nil {
				sumstring := roundVarMap[varRightShift[1]] + ">>" + strconv.Itoa(secondnum)
				roundVarMap[varSet[1]] = sumstring
			} else
			if err2 != nil {
				sumstring := strconv.Itoa(firstnum) + ">>" + roundVarMap[varRightShift[2]]
				roundVarMap[varSet[1]] = sumstring
			} else
			{
				sumstring := strconv.Itoa(int(uint32(firstnum)) >> int(uint32(secondnum)))
				roundVarMap[varSet[1]] = sumstring
			}
			//line = strings.Replace(line, varRightShift[0], "="+sumstring, 1)
			line = ""
		} else
		if len(varAnd) != 0 {
			firstnum, err1 := strconv.Atoi(roundVarMap[varAnd[1]])
			secondnum, err2 := strconv.Atoi(roundVarMap[varAnd[2]])
			if err1 != nil {
				sumstring := roundVarMap[varAnd[1]] + "&" + strconv.Itoa(secondnum)
				roundVarMap[varSet[1]] = sumstring
			} else
			if err2 != nil {
				sumstring := strconv.Itoa(firstnum) + "&" + roundVarMap[varAnd[2]]
				roundVarMap[varSet[1]] = sumstring
			} else
			{
				sumstring := strconv.Itoa(firstnum & secondnum)
				roundVarMap[varSet[1]] = sumstring
			}
			//line = strings.Replace(line, varAnd[0], "="+sumstring, 1)
			line = ""
		} else
		if len(functionRun) != 0 {
			arguments := []string{}
			for _, arg := range strings.Split(functionRun[2], ",") {
				stackmatch := rawStackMatch.FindStringSubmatch(arg)
				varmatch := rawVarMatch.FindStringSubmatch(arg)
				if len(stackmatch) != 0 {
					subval, _ := strconv.Atoi(stackmatch[2])
					arguments = append(arguments, stackin[len(stackin)-subval])
				} else if len(varmatch) != 0 {
					arguments = append(arguments, roundVarMap[varmatch[2]])
				}
			}
			roundVarMap[varSet[1]] = varSet[1]
			line = strings.Replace(line, functionRun[0], fmt.Sprintf("=%s(%s)", roundVarMap[functionRun[1]], strings.Join(arguments, ",")), 1)
		} else
		if len(stringCreated) != 0 {
			roundVarMap[varSet[1]] = roundVarMap[stringCreated[1]]
			line = ""
		} else
		if len(stackLengthGet) != 0 {
			setlen := len(stackin)
			if len(stackLengthGet) > 1 {
				subval, _ := strconv.Atoi(stackLengthGet[1])
				setlen -= subval
			}
			setval := strconv.Itoa(setlen)
			roundVarMap[varSet[1]] = setval
			//line = strings.Replace(line, stackLengthGet[0], "="+setval, 1)
			line = ""
		} else
		{
			blocks := strings.Split(line, "=")
			for _, repval := range rawStackMatch.FindAllStringSubmatch(blocks[1], -1) {
				subval, _ := strconv.Atoi(repval[2])
				blocks[1] = strings.Replace(blocks[1], repval[0], repval[1]+stackin[len(stackin)-subval], 1)
			}
			for _, repval := range rawStackMatch.FindAllStringSubmatch(blocks[1], -1) {
				subval, _ := strconv.Atoi(repval[2])
				blocks[1] = strings.Replace(blocks[1], repval[0], repval[1]+stackin[len(stackin)-subval], 1)
			}
			for _, repval := range rawVarMatch.FindAllStringSubmatch(blocks[1], -1) {
				blocks[1] = strings.Replace(blocks[1], repval[0], repval[1]+roundVarMap[repval[2]]+repval[3], 1)
			}
			for _, repval := range rawVarMatch.FindAllStringSubmatch(blocks[1], -1) {
				blocks[1] = strings.Replace(blocks[1], repval[0], repval[1]+roundVarMap[repval[2]]+repval[3], 1)
			}
			roundVarMap[varSet[1]] = blocks[1]
			line = strings.Join(blocks, "=")
		}
	} else
	if len(heapSet) != 0 {
		heapindex, _ := strconv.Atoi(roundVarMap[heapSet[1]])
		for len(heapin) <= heapindex {
			heapin = append(heapin, "")
		}
		stackGet := stackGetRegex.FindStringSubmatch(line)
		varGet := varGetRegex.FindStringSubmatch(line)
		if len(stackGet) != 0 {
			subnum, _ := strconv.Atoi(stackGet[1])
			heapin[heapindex] = stackin[len(stackin)-subnum]
			line = strings.Replace(line, stackGet[0], "="+stackin[len(stackin)-subnum], 1)
		} else
		if len(varGet) != 0 {
			heapin[heapindex] = roundVarMap[varGet[1]]
			line = strings.Replace(line, varGet[0], "="+roundVarMap[varGet[1]], 1)
		} else {
			blocks := strings.Split(line, "=")
			for _, repval := range rawStackMatch.FindAllStringSubmatch(blocks[1], -1) {
				subval, _ := strconv.Atoi(repval[2])
				blocks[1] = strings.Replace(blocks[1], repval[0], repval[1]+stackin[len(stackin)-subval], 1)
			}
			for _, repval := range rawStackMatch.FindAllStringSubmatch(blocks[1], -1) {
				subval, _ := strconv.Atoi(repval[2])
				blocks[1] = strings.Replace(blocks[1], repval[0], repval[1]+stackin[len(stackin)-subval], 1)
			}
			for _, repval := range rawVarMatch.FindAllStringSubmatch(blocks[1], -1) {
				blocks[1] = strings.Replace(blocks[1], repval[0], repval[1]+roundVarMap[repval[2]]+repval[3], 1)
			}
			for _, repval := range rawVarMatch.FindAllStringSubmatch(blocks[1], -1) {
				blocks[1] = strings.Replace(blocks[1], repval[0], repval[1]+roundVarMap[repval[2]]+repval[3], 1)
			}
			heapin[heapindex] = blocks[1]
			line = strings.Join(blocks, "=")
		}
		line = strings.Replace(line, fmt.Sprintf("[%s]", heapSet[1]), "["+roundVarMap[heapSet[1]]+"]", 1)
	} else
	if len(stackSetStack) != 0 || len(stackSetVar) != 0 {
		var stackindex int
		switch len(stackSetStack) {
		case 0:
			varnum, _ := strconv.Atoi(roundVarMap[stackSetVar[1]])
			stackindex = varnum
			if len(stackSetVar) == 3 {
				addnum, _ := strconv.Atoi(stackSetVar[2])
				stackindex += addnum
			}
		case 1:
			stackindex = len(stackin)
		case 2:
			subindex, _ := strconv.Atoi(stackSetStack[1])
			stackindex = len(stackin) - subindex
		}
		if stackindex == len(stackin) {
			stackin = append(stackin, "")
			repbase := strings.Split(line, "=")[1]
			repblock := rawStackMatch.FindStringSubmatch(repbase)
			if len(repblock) > 0 {
				original := repblock[0]
				subnum, _ := strconv.Atoi(repblock[2])
				subnum++
				line = strings.ReplaceAll(line, original, repblock[1]+fmt.Sprintf("varin.stack1[varin.stack1.length-%d]", subnum))
			}
		}
		varGet := varGetRegex.FindStringSubmatch(line)
		stackGet := stackGetRegex.FindStringSubmatch(line)
		runFuncGet := runFuncGetRegex.FindStringSubmatch(line)
		heapGet := heapGetRegex.FindStringSubmatch(line)
		arrCreate := arrCreateRegex.FindStringSubmatch(line)
		objCreate := objCreateRegex.FindStringSubmatch(line)
		stackMod := stackModRegex.FindStringSubmatch(line)
		varAdd := varAddRegex.FindStringSubmatch(line)
		stackAdd := stackAddRegex.FindStringSubmatch(line)
		varSubtract := varSubtractRegex.FindStringSubmatch(line)
		stackSubtract := stackSubtractRegex.FindStringSubmatch(line)
		stackDivision := stackDivisionRegex.FindStringSubmatch(line)
		stackMultiply := stackMultiplyRegex.FindStringSubmatch(line)
		stackLeftShift := stackLeftShiftRegex.FindStringSubmatch(line)
		varRightShift := varRightShiftRegex.FindStringSubmatch(line)
		stackRightShiftTwo := stackRightShiftTwoRegex.FindStringSubmatch(line)
		stackRightShiftThree := stackRightShiftThreeRegex.FindStringSubmatch(line)
		varAnd := varAndRegex.FindStringSubmatch(line)
		stackAnd := stackAndRegex.FindStringSubmatch(line)
		stackXor := stackXorRegex.FindStringSubmatch(line)
		varEqual := varEqualRegex.FindStringSubmatch(line)
		stackEqual := stackEqualRegex.FindStringSubmatch(line)
		stackTripleEqual := stackTripleEqualRegex.FindStringSubmatch(line)
		trueGet := trueGetRegex.FindStringSubmatch(line)
		falseGet := falseGetRegex.FindStringSubmatch(line)
		functionRun := functionRunRegex.FindStringSubmatch(line)
		greaterEqualVarBool := greaterEqualVarBoolRegex.FindStringSubmatch(line)
		lessEqualVarBool := lessEqualVarBoolRegex.FindStringSubmatch(line)
		lessStackBool := lessStackBoolRegex.FindStringSubmatch(line)
		greaterStackBool := greaterStackBoolRegex.FindStringSubmatch(line)
		lessvarBool := lessvarBoolRegex.FindStringSubmatch(line)
		stackObjectGet := stackObjectGetRegex.FindStringSubmatch(line)
		stackObjectFuncGet := stackObjectFuncGetRegex.FindStringSubmatch(line)
		stringCreated := stringCreatedRegex.FindStringSubmatch(line)
		stackLengthGet := stackLengthGetRegex.FindStringSubmatch(line)
		if len(varGet) != 0 {
			stackin[stackindex] = roundVarMap[varGet[1]]
			//line = strings.Replace(line, varGet[0], "="+roundVarMap[varGet[1]], 1)
			line = ""
		} else
		if len(stackGet) != 0 {
			subnum, _ := strconv.Atoi(stackGet[1])
			stackin[stackindex] = stackin[len(stackin)-subnum]
			//line = strings.Replace(line, stackGet[0], "="+stackin[len(stackin)-subnum], 1)
			line = ""
		} else
		if len(runFuncGet) != 0 {
			var argarr []string
			for _, arg := range runFuncGet[1:] {
				rawStack := rawStackMatch.FindStringSubmatch(arg)
				rawVar := rawVarMatch.FindStringSubmatch(arg)
				if len(rawStack) != 0 {
					subnum, _ := strconv.Atoi(rawStack[2])
					argarr = append(argarr, stackin[len(stackin)-subnum])
				} else if len(rawVar) != 0 {
					argarr = append(argarr, roundVarMap[rawVar[2]])
				}
			}
			stackin[stackindex] = fmt.Sprintf("func_%s_%s", argarr[1], argarr[2])
			copyindexinit, _ := strconv.Atoi(argarr[0])
			xcordinit, _ := strconv.Atoi(argarr[1])
			ycordinit, _ := strconv.Atoi(argarr[2])
			g.Initializerholder = append(g.Initializerholder, FuncInitializer{
				X:         xcordinit,
				Y:         ycordinit,
				CopyIndex: copyindexinit,
				Heap:      heapin[:],
			})
			line = strings.Replace(line, runFuncGet[0], fmt.Sprintf("=runFuncDeclaration(%s,heap)", strings.Join(argarr, ",")), 1)
		} else
		if len(heapGet) != 0 {
			heapindex, _ := strconv.Atoi(roundVarMap[heapGet[1]])
			//stackin[stackindex] = heapin[heapindex]
			stackin[stackindex] = fmt.Sprintf("heapin[%d]", heapindex)
			//line = strings.Replace(line, heapGet[0], "="+heapin[heapindex], 1)
			line = ""
		} else
		if len(arrCreate) != 0 {
			stackin[stackindex] = fmt.Sprintf("arrobj%d", g.ArrDefCount)
			line = fmt.Sprintf("arrobj%d=[]", g.ArrDefCount)
			g.ArrDefCount++
		} else
		if len(objCreate) != 0 {
			stackin[stackindex] = fmt.Sprintf("objobj%d", g.ObjDefCount)
			line = fmt.Sprintf("objobj%d=[]", g.ObjDefCount)
			g.ObjDefCount++
		} else
		if len(stackMod) != 0 {
			firstnum, err1 := strconv.Atoi(stackin[len(stackin)-2])
			secondnum, err2 := strconv.Atoi(stackin[len(stackin)-1])
			if err1 != nil {
				sumstring := stackin[len(stackin)-2] + "%" + strconv.Itoa(secondnum)
				stackin[stackindex] = sumstring
			} else
			if err2 != nil {
				sumstring := strconv.Itoa(firstnum) + "%" + stackin[len(stackin)-1]
				stackin[stackindex] = sumstring
			} else
			{
				sumstring := strconv.Itoa(firstnum % secondnum)
				stackin[stackindex] = sumstring
			}
			line = ""
		} else
		if len(varAdd) != 0 {
			firstnum, err1 := strconv.Atoi(roundVarMap[varAdd[1]])
			secondnum, err2 := strconv.Atoi(roundVarMap[varAdd[2]])
			if err1 != nil {
				sumstring := roundVarMap[varAdd[1]] + "+" + strconv.Itoa(secondnum)
				stackin[stackindex] = sumstring
			} else
			if err2 != nil {
				sumstring := strconv.Itoa(firstnum) + "+" + roundVarMap[varAdd[2]]
				stackin[stackindex] = sumstring
			} else
			{
				sumstring := strconv.Itoa(firstnum + secondnum)
				stackin[stackindex] = sumstring
			}
			//line = strings.Replace(line, varAdd[0], sumstring, 1)
			line = ""
		} else
		if len(stackAdd) != 0 {
			firstnum, err1 := strconv.Atoi(stackin[len(stackin)-2])
			secondnum, err2 := strconv.Atoi(stackin[len(stackin)-1])
			if err1 != nil {
				sumstring := stackin[len(stackin)-2] + "+" + strconv.Itoa(secondnum)
				stackin[stackindex] = sumstring
			} else
			if err2 != nil {
				sumstring := strconv.Itoa(firstnum) + "+" + stackin[len(stackin)-1]
				stackin[stackindex] = sumstring
			} else
			{
				sumstring := strconv.Itoa(firstnum + secondnum)
				stackin[stackindex] = sumstring
			}
			//line = strings.Replace(line, stackAdd[0], sumstring, 1)
			line = ""
		} else
		if len(varSubtract) != 0 {
			firstnum, err1 := strconv.Atoi(roundVarMap[varSubtract[1]])
			secondnum, err2 := strconv.Atoi(roundVarMap[varSubtract[2]])
			if err1 != nil {
				sumstring := roundVarMap[varSubtract[1]] + "-" + strconv.Itoa(secondnum)
				stackin[stackindex] = sumstring
			} else
			if err2 != nil {
				sumstring := strconv.Itoa(firstnum) + "-" + roundVarMap[varSubtract[2]]
				stackin[stackindex] = sumstring
			} else
			{
				sumstring := strconv.Itoa(firstnum - secondnum)
				stackin[stackindex] = sumstring
			}
			//line = strings.Replace(line, varSubtract[0], sumstring, 1)
			line = ""
		} else
		if len(stackSubtract) != 0 {
			firstnum, err1 := strconv.Atoi(stackin[len(stackin)-2])
			secondnum, err2 := strconv.Atoi(stackin[len(stackin)-1])
			if err1 != nil {
				sumstring := stackin[len(stackin)-2] + "-" + strconv.Itoa(secondnum)
				stackin[stackindex] = sumstring
			} else
			if err2 != nil {
				sumstring := strconv.Itoa(firstnum) + "-" + stackin[len(stackin)-1]
				stackin[stackindex] = sumstring
			} else
			{
				sumstring := strconv.Itoa(firstnum - secondnum)
				stackin[stackindex] = sumstring
			}
			//line = strings.Replace(line, stackSubtract[0], sumstring, 1)
			line = ""
		} else
		if len(stackDivision) != 0 {
			firstnum, err1 := strconv.Atoi(stackin[len(stackin)-2])
			secondnum, err2 := strconv.Atoi(stackin[len(stackin)-1])
			if err1 != nil {
				sumstring := stackin[len(stackin)-2] + "/" + strconv.Itoa(secondnum)
				stackin[stackindex] = sumstring
			} else
			if err2 != nil {
				sumstring := strconv.Itoa(firstnum) + "/" + stackin[len(stackin)-1]
				stackin[stackindex] = sumstring
			} else
			if firstnum == -1 && secondnum == 0 {
				stackin[stackindex] = "-Infinity"
			} else
			if firstnum == 1 && secondnum == 0 {
				stackin[stackindex] = "Infinity"
			} else
			{
				sumstring := strconv.Itoa(firstnum / secondnum)
				stackin[stackindex] = sumstring
			}
			//line = strings.Replace(line, stackDivision[0], sumstring, 1)
			line = ""
		} else
		if len(stackMultiply) != 0 {
			firstnum, err1 := strconv.Atoi(stackin[len(stackin)-2])
			secondnum, err2 := strconv.Atoi(stackin[len(stackin)-1])
			if err1 != nil {
				sumstring := stackin[len(stackin)-2] + "*" + strconv.Itoa(secondnum)
				stackin[stackindex] = sumstring
			} else
			if err2 != nil {
				sumstring := strconv.Itoa(firstnum) + "*" + stackin[len(stackin)-1]
				stackin[stackindex] = sumstring
			} else
			{
				sumstring := strconv.Itoa(firstnum * secondnum)
				stackin[stackindex] = sumstring
			}
			//line = strings.Replace(line, stackMultiply[0], sumstring, 1)
			line = ""
		} else
		if len(stackLeftShift) != 0 {
			firstnum, err1 := strconv.Atoi(stackin[len(stackin)-2])
			secondnum, err2 := strconv.Atoi(stackin[len(stackin)-1])
			if err1 != nil {
				sumstring := fmt.Sprintf("%s << %d", stackin[len(stackin)-2], secondnum)
				stackin[stackindex] = sumstring
			} else
			if err2 != nil {
				sumstring := fmt.Sprintf("%d << %s", firstnum, stackin[len(stackin)-1])
				stackin[stackindex] = sumstring
			} else
			{
				sumstring := strconv.Itoa(firstnum << secondnum)
				stackin[stackindex] = sumstring
			}
			//line = strings.Replace(line, stackLeftShift[0], sumstring, 1)
			line = ""
		} else
		if len(varRightShift) != 0 {
			firstnum, err1 := strconv.Atoi(roundVarMap[varRightShift[1]])
			secondnum, err2 := strconv.Atoi(roundVarMap[varRightShift[2]])
			if err1 != nil {
				sumstring := roundVarMap[varRightShift[1]] + ">>" + strconv.Itoa(secondnum)
				stackin[stackindex] = sumstring
			} else
			if err2 != nil {
				sumstring := strconv.Itoa(firstnum) + ">>" + roundVarMap[varRightShift[2]]
				stackin[stackindex] = sumstring
			} else
			{
				sumstring := strconv.Itoa(int(uint32(firstnum)) >> int(uint32(secondnum)))
				stackin[stackindex] = sumstring
			}
			//line = strings.Replace(line, varRightShift[0], sumstring, 1)
			line = ""
		} else
		if len(stackRightShiftTwo) != 0 {
			firstnum, err1 := strconv.Atoi(stackin[len(stackin)-2])
			secondnum, err2 := strconv.Atoi(stackin[len(stackin)-1])
			if err1 != nil {
				sumstring := stackin[len(stackin)-2] + ">>" + strconv.Itoa(secondnum)
				stackin[stackindex] = sumstring
			} else
			if err2 != nil {
				sumstring := strconv.Itoa(firstnum) + ">>" + stackin[len(stackin)-1]
				stackin[stackindex] = sumstring
			} else
			{
				sumstring := strconv.Itoa(firstnum >> secondnum)
				stackin[stackindex] = sumstring
			}
			//line = strings.Replace(line, stackRightShiftTwo[0], sumstring, 1)
			line = ""
		} else
		if len(stackRightShiftThree) != 0 {
			firstnum, err1 := strconv.Atoi(stackin[len(stackin)-2])
			secondnum, err2 := strconv.Atoi(stackin[len(stackin)-1])
			if err1 != nil {
				sumstring := stackin[len(stackin)-2] + ">>" + strconv.Itoa(secondnum)
				stackin[stackindex] = sumstring
			} else
			if err2 != nil {
				sumstring := strconv.Itoa(firstnum) + ">>" + stackin[len(stackin)-1]
				stackin[stackindex] = sumstring
			} else
			{
				sumstring := strconv.Itoa(int(uint32(firstnum)) >> int(uint32(secondnum)))
				stackin[stackindex] = sumstring
			}
			//line = strings.Replace(line, stackRightShiftThree[0], sumstring, 1)
			line = ""
		} else
		if len(varAnd) != 0 {
			firstnum, err1 := strconv.Atoi(roundVarMap[varAnd[1]])
			secondnum, err2 := strconv.Atoi(roundVarMap[varAnd[2]])
			if err1 != nil {
				sumstring := roundVarMap[varAnd[1]] + "&" + strconv.Itoa(secondnum)
				stackin[stackindex] = sumstring
			} else
			if err2 != nil {
				sumstring := strconv.Itoa(firstnum) + "&" + roundVarMap[varAnd[2]]
				stackin[stackindex] = sumstring
			} else
			{
				sumstring := strconv.Itoa(firstnum & secondnum)
				stackin[stackindex] = sumstring
			}
			//line = strings.Replace(line, varAnd[0], sumstring, 1)
			line = ""
		} else
		if len(stackAnd) != 0 {
			firstnum, err1 := strconv.Atoi(stackin[len(stackin)-2])
			secondnum, err2 := strconv.Atoi(stackin[len(stackin)-1])
			if err1 != nil {
				sumstring := stackin[len(stackin)-2] + "&" + strconv.Itoa(secondnum)
				stackin[stackindex] = sumstring
			} else
			if err2 != nil {
				sumstring := strconv.Itoa(firstnum) + "&" + stackin[len(stackin)-1]
				stackin[stackindex] = sumstring
			} else
			{
				sumstring := strconv.Itoa(firstnum & secondnum)
				stackin[stackindex] = sumstring
			}
			//line = strings.Replace(line, stackAnd[0], sumstring, 1)
			line = ""
		} else
		if len(stackXor) != 0 {
			firstnum, err1 := strconv.Atoi(stackin[len(stackin)-2])
			secondnum, err2 := strconv.Atoi(stackin[len(stackin)-1])
			if err1 != nil {
				sumstring := stackin[len(stackin)-2] + "^" + strconv.Itoa(secondnum)
				stackin[stackindex] = sumstring
			} else
			if err2 != nil {
				sumstring := strconv.Itoa(firstnum) + "^" + stackin[len(stackin)-1]
				stackin[stackindex] = sumstring
			} else
			{
				sumstring := strconv.Itoa(firstnum ^ secondnum)
				stackin[stackindex] = sumstring
			}
			//line = strings.Replace(line, stackXor[0], sumstring, 1)
			line = ""
		} else
		if len(varEqual) != 0 {
			sumstring := strconv.FormatBool(roundVarMap[varEqual[1]] == roundVarMap[varEqual[2]])
			stackin[stackindex] = sumstring
			//line = strings.Replace(line, varEqual[0], sumstring, 1)
			line = ""
		} else
		if len(stackEqual) != 0 {
			sumstring := strconv.FormatBool(stackin[len(stackin)-2] == stackin[len(stackin)-1])
			stackin[stackindex] = sumstring
			line = strings.Replace(line, stackEqual[0], sumstring, 1)
			//line = ""
		} else
		if len(stackTripleEqual) != 0 {
			sumstring := strconv.FormatBool(stackin[len(stackin)-2] == stackin[len(stackin)-1])
			stackin[stackindex] = sumstring
			line = strings.Replace(line, stackTripleEqual[0], sumstring, 1)
			//line = ""
		} else
		if len(trueGet) != 0 {
			stackin[stackindex] = "true"
		} else
		if len(falseGet) != 0 {
			stackin[stackindex] = "false"
		} else
		if len(functionRun) != 0 {
			arguments := []string{}
			for _, arg := range strings.Split(functionRun[2], ",") {
				stackmatch := rawStackMatch.FindStringSubmatch(arg)
				varmatch := rawVarMatch.FindStringSubmatch(arg)
				if len(stackmatch) != 0 {
					subval, _ := strconv.Atoi(stackmatch[2])
					arguments = append(arguments, stackin[len(stackin)-subval])
				} else if len(varmatch) != 0 {
					arguments = append(arguments, roundVarMap[varmatch[2]])
				}
			}
			stackin[stackindex] = fmt.Sprintf("varin.stack1[%d]", stackindex)
			line = strings.Replace(line, functionRun[0], fmt.Sprintf("=%s(%s)", roundVarMap[functionRun[1]], strings.Join(arguments, ",")), 1)
		} else
		if len(greaterEqualVarBool) != 0 {
			firstnum, err1 := strconv.Atoi(roundVarMap[greaterEqualVarBool[1]])
			secondnum, err2 := strconv.Atoi(roundVarMap[greaterEqualVarBool[2]])
			if err1 != nil {
				sumstring := roundVarMap[greaterEqualVarBool[1]] + ">=" + strconv.Itoa(secondnum)
				line = strings.Replace(line, greaterEqualVarBool[0], sumstring, 1)
				stackin[stackindex] = sumstring
			} else
			if err2 != nil {
				sumstring := strconv.Itoa(firstnum) + ">=" + roundVarMap[greaterEqualVarBool[2]]
				line = strings.Replace(line, greaterEqualVarBool[0], sumstring, 1)
				stackin[stackindex] = sumstring
			} else
			{
				sumstring := strconv.FormatBool(firstnum >= secondnum)
				line = strings.Replace(line, greaterEqualVarBool[0], sumstring, 1)
				stackin[stackindex] = sumstring
			}
		} else
		if len(lessEqualVarBool) != 0 {
			firstnum, err1 := strconv.Atoi(roundVarMap[lessEqualVarBool[1]])
			secondnum, err2 := strconv.Atoi(roundVarMap[lessEqualVarBool[2]])
			if err1 != nil {
				sumstring := roundVarMap[lessEqualVarBool[1]] + "<=" + strconv.Itoa(secondnum)
				line = strings.Replace(line, lessEqualVarBool[0], sumstring, 1)
				stackin[stackindex] = sumstring
			} else
			if err2 != nil {
				sumstring := strconv.Itoa(firstnum) + "<=" + roundVarMap[lessEqualVarBool[2]]
				line = strings.Replace(line, lessEqualVarBool[0], sumstring, 1)
				stackin[stackindex] = sumstring
			} else
			{
				sumstring := strconv.FormatBool(firstnum <= secondnum)
				line = strings.Replace(line, lessEqualVarBool[0], sumstring, 1)
				stackin[stackindex] = sumstring
			}
		} else
		if len(lessStackBool) != 0 {
			firstnum, err1 := strconv.Atoi(stackin[len(stackin)-2])
			secondnum, err2 := strconv.Atoi(stackin[len(stackin)-1])
			if err1 != nil {
				sumstring := stackin[len(stackin)-2] + "<" + strconv.Itoa(secondnum)
				line = strings.Replace(line, lessStackBool[0], sumstring, 1)
				stackin[stackindex] = sumstring
			} else
			if err2 != nil {
				sumstring := strconv.Itoa(firstnum) + "<" + stackin[len(stackin)-1]
				line = strings.Replace(line, lessStackBool[0], sumstring, 1)
				stackin[stackindex] = sumstring
			} else
			{
				sumstring := strconv.FormatBool(firstnum < secondnum)
				line = strings.Replace(line, lessStackBool[0], sumstring, 1)
				stackin[stackindex] = sumstring
			}
		} else
		if len(greaterStackBool) != 0 {
			firstnum, err1 := strconv.Atoi(stackin[len(stackin)-2])
			secondnum, err2 := strconv.Atoi(stackin[len(stackin)-1])
			if err1 != nil {
				sumstring := stackin[len(stackin)-2] + ">" + strconv.Itoa(secondnum)
				line = strings.Replace(line, greaterStackBool[0], sumstring, 1)
				stackin[stackindex] = sumstring
			} else
			if err2 != nil {
				sumstring := strconv.Itoa(firstnum) + ">" + stackin[len(stackin)-1]
				line = strings.Replace(line, greaterStackBool[0], sumstring, 1)
				stackin[stackindex] = sumstring
			} else
			{
				sumstring := strconv.FormatBool(firstnum > secondnum)
				line = strings.Replace(line, greaterStackBool[0], sumstring, 1)
				stackin[stackindex] = sumstring
			}
		} else
		if len(lessvarBool) != 0 {
			firstnum, err1 := strconv.Atoi(roundVarMap[lessvarBool[1]])
			secondnum, err2 := strconv.Atoi(roundVarMap[lessvarBool[2]])
			if err1 != nil {
				sumstring := roundVarMap[lessvarBool[1]] + "<" + strconv.Itoa(secondnum)
				line = strings.Replace(line, lessvarBool[0], sumstring, 1)
				stackin[stackindex] = sumstring
			} else
			if err2 != nil {
				sumstring := strconv.Itoa(firstnum) + "<" + roundVarMap[lessvarBool[2]]
				line = strings.Replace(line, lessvarBool[0], sumstring, 1)
				stackin[stackindex] = sumstring
			} else
			{
				sumstring := strconv.FormatBool(firstnum < secondnum)
				line = strings.Replace(line, lessvarBool[0], sumstring, 1)
				stackin[stackindex] = sumstring
			}
		} else
		if len(stackObjectGet) != 0 {
			setobj := stackin[len(stackin)-2] + "[" + stackin[len(stackin)-1] + "]"
			stackin[stackindex] = setobj
			line = strings.Replace(line, stackObjectGet[0], "="+setobj, 1)
		} else
		if len(stackObjectFuncGet) != 0 {
			setobj := fmt.Sprintf("varin.stack1[%d]", stackindex)
			line = strings.Replace(line, stackObjectFuncGet[0], "="+stackin[len(stackin)-2]+"["+stackin[len(stackin)-1]+"]()", 1)
			stackin[stackindex] = setobj
		} else
		if len(stringCreated) != 0 {
			stackin[stackindex] = roundVarMap[stringCreated[1]]
			line = strings.Replace(line, stringCreated[0], "="+roundVarMap[stringCreated[1]], 1)
		} else
		if len(stackLengthGet) != 0 {
			setlen := len(stackin)
			if len(stackLengthGet) > 1 {
				subval, _ := strconv.Atoi(stackLengthGet[1])
				setlen -= subval
			}
			setval := strconv.Itoa(setlen)
			stackin[stackindex] = setval
			//line = strings.Replace(line, stackLengthGet[0], "="+setval, 1)
			line = ""
		} else
		{
			blocks := strings.Split(line, "=")
			for _, repval := range rawStackMatch.FindAllStringSubmatch(blocks[1], -1) {
				subval, _ := strconv.Atoi(repval[2])
				blocks[1] = strings.Replace(blocks[1], repval[0], repval[1]+stackin[len(stackin)-subval], 1)
			}
			for _, repval := range rawStackMatch.FindAllStringSubmatch(blocks[1], -1) {
				subval, _ := strconv.Atoi(repval[2])
				blocks[1] = strings.Replace(blocks[1], repval[0], repval[1]+stackin[len(stackin)-subval], 1)
			}
			for _, repval := range rawVarMatch.FindAllStringSubmatch(blocks[1], -1) {
				if _, ok := roundVarMap[repval[2]]; ok {
					blocks[1] = strings.Replace(blocks[1], repval[0], repval[1]+roundVarMap[repval[2]]+repval[3], 1)
				}
			}
			for _, repval := range rawVarMatch.FindAllStringSubmatch(blocks[1], -1) {
				if _, ok := roundVarMap[repval[2]]; ok {
					blocks[1] = strings.Replace(blocks[1], repval[0], repval[1]+roundVarMap[repval[2]]+repval[3], 1)
				}
			}
			stackin[stackindex] = blocks[1]
			line = strings.Join(blocks, "=")
		}
	} else
	if len(objectFromStackSet) != 0 {
		getindex, _ := strconv.Atoi(stackGetRegex.FindStringSubmatch(line)[1])
		line = stackin[len(stackin)-2] + "[" + stackin[len(stackin)-1] + "]=" + stackin[len(stackin)-getindex]
	} else
	if len(XCoordSet) != 0 {
		varGet := varGetRegex.FindStringSubmatch(line)
		stackGet := stackGetRegex.FindStringSubmatch(line)
		jumpHolderXGet := jumpHolderXGetRegex.FindStringSubmatch(line)
		if len(varGet) != 0 {
			xset, _ := strconv.Atoi(roundVarMap[varGet[1]])
			x = xset
			//line = strings.Replace(line, varGet[0], "="+roundVarMap[varGet[1]], 1)
			line = ""
		} else
		if len(stackGet) != 0 {
			subnum, _ := strconv.Atoi(stackGet[1])
			xset, _ := strconv.Atoi(stackin[len(stackin)-subnum])
			x = xset
			//line = strings.Replace(line, stackGet[0], "="+stackin[len(stackin)-subnum], 1)
			line = ""
		} else
		if len(jumpHolderXGet) != 0 {
			x = jumper.x
			//line = strings.Replace(line, jumpHolderXGet[0], "="+strconv.Itoa(x), 1)
			line = ""
		}
	} else
	if len(YCoordSet) != 0 {
		varGet := varGetRegex.FindStringSubmatch(line)
		stackGet := stackGetRegex.FindStringSubmatch(line)
		jumpHolderYGet := jumpHolderYGetRegex.FindStringSubmatch(line)
		if len(varGet) != 0 {
			xset, _ := strconv.Atoi(roundVarMap[varGet[1]])
			y = xset
			//line = strings.Replace(line, varGet[0], "="+roundVarMap[varGet[1]], 1)
			line = ""
		} else
		if len(stackGet) != 0 {
			subnum, _ := strconv.Atoi(stackGet[1])
			xset, _ := strconv.Atoi(stackin[len(stackin)-subnum])
			y = xset
			//line = strings.Replace(line, stackGet[0], "="+stackin[len(stackin)-subnum], 1)
			line = ""
		} else
		if len(jumpHolderYGet) != 0 {
			y = jumper.y
			//line = strings.Replace(line, jumpHolderYGet[0], "="+strconv.Itoa(y), 1)
			line = ""
		}
	} else
	if len(stackPush) != 0 {
		stackin = append(stackin, stackPush[1])
		line = ""
	} else
	if len(jumpHolderSet) != 0 {
		jumper.x = x
		jumper.y = y
		//line = strings.Replace(line, "varin.Xcoord", strconv.Itoa(x), 1)
		//line = strings.Replace(line, "varin.Ycoord", strconv.Itoa(y), 1)
		line = ""
	} else
	if len(stack3Push) != 0 {
		xcord, _ := strconv.Atoi(roundVarMap[stack3Push[1]])
		ycord, _ := strconv.Atoi(roundVarMap[stack3Push[2]])
		found := false
		for _, jo := range g.Stack3Holder {
			if xcord == jo.x && ycord == jo.y {
				found = true
			}
		}
		if !found {
			g.Stack3Holder = append(g.Stack3Holder, Jumper{x: xcord, y: ycord})
		}
	}
	{
		blocks := strings.Split(line, "=")
		var repblock string
		var base string
		if len(blocks) > 1 {
			repblock = blocks[1]
			base = repblock
		} else {
			repblock = line
			base = repblock
		}
		for _, repval := range rawStackMatch.FindAllStringSubmatch(repblock, -1) {
			subval, _ := strconv.Atoi(repval[2])
			repblock = strings.Replace(repblock, repval[0], repval[1]+stackin[len(stackin)-subval], 1)
		}
		for _, repval := range rawStackMatch.FindAllStringSubmatch(repblock, -1) {
			subval, _ := strconv.Atoi(repval[2])
			repblock = strings.Replace(repblock, repval[0], repval[1]+stackin[len(stackin)-subval], 1)
		}
		for _, repval := range rawVarMatch.FindAllStringSubmatch(repblock, -1) {
			repblock = strings.Replace(repblock, repval[0], repval[1]+roundVarMap[repval[2]]+repval[3], 1)
		}
		for _, repval := range rawVarMatch.FindAllStringSubmatch(repblock, -1) {
			repblock = strings.Replace(repblock, repval[0], repval[1]+roundVarMap[repval[2]]+repval[3], 1)
		}
		line = strings.ReplaceAll(line, base, repblock)
	}
	if line != "" {
		line += ";"
	}
	return line, x, y, stackin, heapin, jumper, roundVarMap
}
func (g *GlobalHolder) walker(x, y int, stackin, heapin []string, jumper Jumper, tagline string, tagarrayin map[string]bool) map[string]bool {
	log.SetFlags(0)
	outstr := strings.Builder{}
	escapes := 1
	nextop := 0
	for escapes > 0 {
		roundVarMap := make(map[string]string)
		x, y, nextop = g.OpStepper(x, y)
		currop := g.GlobalOpArray[nextop]
		switch currop.Type {
		case 2:
			lastx := x
			lasty := y
			switch currop.Path2 {
			case true:
				var path2string strings.Builder
				tagline1 := "walk_" + strconv.Itoa(g.Routecount) + "route1!"
				tagline2 := "walk_" + strconv.Itoa(g.Routecount) + "route2!"
				g.Routecount++

				for _, line := range currop.Block1Lines {
					line, x, y, stackin, heapin, jumper, roundVarMap = g.parseLine(line, x, y, stackin, heapin, jumper, roundVarMap)
					outstr.WriteString(line)
				}
				wl, _, _, _, _, _, _ := g.parseLine(currop.Conditional, x, y, stackin, heapin, jumper, roundVarMap)
				//path2
				xsubpath1 := x
				ysubpath1 := y
				for _, line := range currop.CondBody {
					line, xsubpath1, ysubpath1, stackin, heapin, jumper, roundVarMap = g.parseLine(line, xsubpath1, ysubpath1, stackin, heapin, jumper, roundVarMap)
					path2string.WriteString(line)
				}
				//path1
				for _, line := range currop.Block2Lines {
					line, x, y, stackin, heapin, jumper, roundVarMap = g.parseLine(line, x, y, stackin, heapin, jumper, roundVarMap)
					outstr.WriteString(line)
				}

				checktag := strconv.Itoa(lastx) + "," + strconv.Itoa(lasty) + "," + strconv.Itoa(x) + "," + strconv.Itoa(y) + "," + strconv.Itoa(xsubpath1) + "," + strconv.Itoa(ysubpath1) + "," + strconv.Itoa(jumper.x) + "," + strconv.Itoa(jumper.y)
				if walked, ok := tagarrayin[checktag]; ok && walked {
					g.Tagholder[tagline] = outstr.String()
					outmap := map[string]bool{
						checktag: true,
					}
					return outmap
				}
				tagarrayin[checktag] = true
				path1walked := g.walker(x, y, stackin, heapin, jumper, tagline1, tagarrayin)
				path2walked := g.walker(xsubpath1, ysubpath1, stackin, heapin, jumper, tagline2, tagarrayin)
				mapunion := make(map[string]bool)
				onewalked := false
				twowalked := false
				if _, ok := path2walked[checktag]; ok {
					delete(path2walked, checktag)
					twowalked = true
				}
				if _, ok := path1walked[checktag]; ok {
					delete(path1walked, checktag)
					onewalked = true
				}
				for k, v := range path2walked {
					mapunion[k] = v
				}
				for k, v := range path1walked {
					mapunion[k] = v
				}
				if !onewalked && !twowalked {
					outstr.WriteString(fmt.Sprintf("if(%s){%s;%s}else{%s}", wl, path2string.String(), tagline2, tagline1))
				} else
				if onewalked && !twowalked {
					outstr.WriteString(fmt.Sprintf("while(!(%s)){%s}%s;%s", wl, tagline1, path2string.String(), tagline2))
				} else
				if !onewalked && twowalked {
					outstr.WriteString(fmt.Sprintf("while(%s){%s;%s}%s", wl, path2string.String(), tagline2, tagline1))
				} else
				if onewalked && twowalked {
					outstr.WriteString(fmt.Sprintf("while(true){if(%s){%s;%s}else{%s}}", wl, path2string.String(), tagline2, tagline1))
				}
				g.Tagholder[tagline] = outstr.String()
				escapes--
				return mapunion
			case false:

				for _, line := range currop.Block1Lines {
					line, x, y, stackin, heapin, jumper, roundVarMap = g.parseLine(line, x, y, stackin, heapin, jumper, roundVarMap)
					outstr.WriteString(line)
				}
				for _, line := range currop.Block2Lines {
					line, x, y, stackin, heapin, jumper, roundVarMap = g.parseLine(line, x, y, stackin, heapin, jumper, roundVarMap)
					outstr.WriteString(line)
				}

			}
		case 1:

			for _, line := range currop.Block1Lines {
				line, x, y, stackin, heapin, jumper, roundVarMap = g.parseLine(line, x, y, stackin, heapin, jumper, roundVarMap)
				outstr.WriteString(line)
			}

		case 0:

			for _, line := range currop.Block1Lines {
				line, x, y, stackin, heapin, jumper, roundVarMap = g.parseLine(line, x, y, stackin, heapin, jumper, roundVarMap)
				outstr.WriteString(line)
			}

			escapes--
		}
	}
	g.Tagholder[tagline] = outstr.String()
	return map[string]bool{}
}
func (g *GlobalHolder) genFile(x, y int, heapin []string) {
	g.Tagholder = map[string]string{}
	g.Initializerholder = []FuncInitializer{}
	g.walker(x, y, []string{}, heapin, Jumper{x: 0, y: 0}, "basescript", map[string]bool{})
	outstring := g.Tagholder["basescript"]
	for len(routeMatch.FindStringSubmatch(outstring)) > 0 {
		for _, match := range routeMatch.FindAllStringSubmatch(outstring, -1) {
			outstring = strings.ReplaceAll(outstring, match[0], strings.Join(strings.Split(g.Tagholder[match[0]], "$"), "dollarsign"))
		}
	}
	outstring = strings.Join(strings.Split(outstring, "dollarsign"), "$")
	for _, val := range g.Initializerholder[:] {
		g.FileLinker(val.CopyIndex, val.X, val.Y, val.Heap)
	}
}
func (g *GlobalHolder) keyResolve() {
	g.genFile(0, 0, []string{})
	for payload, keys := range g.KeyHolderRaw {
		if payload == 100 {
			for _, key := range keys {
				keyInt, _ := strconv.Atoi(key)
				g.EncryptionKeyHolder = append(g.EncryptionKeyHolder, keyInt)
			}
			g.EncryptionKeyHolder = replaceAllOccurrences(g.EncryptionKeyHolder)
		} else {
			for _, key := range keys {
				if g.KeyCount[key] < 5 {
					keyInt, _ := strconv.Atoi(key)
					g.KeyHolder[payload] = g.GlobalKeyArray[keyInt]
				}
			}
		}
	}
}
func (g *GlobalHolder) BaseScrape(src string) {
	blocks := strings.Split(baseMainMatchRegex.FindStringSubmatch(src)[1], "\"")
	g.Alphabet = blocks[3]
	g.Seedstring = blocks[1]
	basekeyarr := []int{}
	for _, num := range strings.Split(blocks[4][5:len(blocks[4])-1], ",") {
		var digitval int64
		if strings.Contains(num, "e") {
			subblocks := strings.Split(num, "e")
			base, _ := strconv.Atoi(subblocks[0])
			endian, _ := strconv.Atoi(subblocks[1])
			digitval = int64(base * int(math.Pow(10.0, float64(endian))))
		} else
		{
			digitval, _ = strconv.ParseInt(num, 0, 64)
		}
		basekeyarr = append(basekeyarr, int(digitval))
	}
	g.Basekeys = basekeyarr
	g.Hashstring = basehashmatch.FindStringSubmatch(src)[1]
}
func (g *GlobalHolder) getResString() (string, error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(res)

	req.SetRequestURI(fmt.Sprintf("https://ponos.zeronaught.com/0?a=22a94427081eb8b3faade27031c844aeedb00212&b=%s&c=1037328191", g.Seedstring))

	if err := fasthttp.Do(req, res); err != nil {
		return "", err
	}

	return string(res.Body()), nil
}

func Generator(src, basefile string) *GlobalHolder {
	gholder := GlobalHolder{
		Stack3Holder:        []Jumper{},
		RanKeys:             make(map[int]bool),
		Negatives:           make(map[int]bool),
		Encrounds:           0,
		PayloadOrder:        []int{},
		IteratedFiles:       make(map[string]bool),
		Routecount:          0,
		KeyCount:            make(map[string]int),
		KeyHolderRaw:        make(map[int][]string),
		KeyHolder:           make(map[int]int64),
		EncryptionKeyHolder: []int{},
		Tagholder:           make(map[string]string),
		Initializerholder:   []FuncInitializer{},
		GlobalStringArray:   []string{},
		GlobalTupleArray:    [][][]int{},
		GlobalKeyArray:      []int64{},
		GlobalHeapCopyArray: []InitVarObj{},
		GlobalOpArray:       []Op{},
		MainNumArray:        []uint8{},
		ObjDefCount:         0,
		ArrDefCount:         0,
		StrDefCount:         0,
		IdentifierMap:       make(map[string]string),
		Basekeys:            []int{},
		Seedstring:          "",
		Alphabet:            "",
		Base:                src,
		P1VAL:               []int{},
		P2VAL:               []int{},
		P3VAL:               []int{},
		P11VAL:              []int{},
		P13VAL:              []int{},
		P14VAL:              []int{},
		P15VAL:              []int{},
		P16VAL:              []int{},
		P17VAL:              [][]int{},
		P18VAL:              []int{},
		P20VAL:              []int{},
		P21VAL:              []int{},
		P23VAL:              []int{},
		P31VAL:              []int{},
		P32VAL:              []int{},
		P35VAL:              []int{},
		P36VAL:              []int{},
		P40VAL:              []int{},
		P42VAL:              []int{},
		P44VAL:              []int{},
	}
	gholder.DeobGlobals(src)
	gholder.BaseScrape(basefile)
	gholder.keyResolve()
	gholder.PayloadOrder = replaceAllOccurrences(gholder.PayloadOrder)
	return &gholder
}
func (g *GlobalHolder) FileLinker(copyindex, x, y int, heapin []string) {
	starttag := "//" + strings.Join(heapin, ",") + "\n"
	ctag := strconv.Itoa(x) + "_" + strconv.Itoa(y)
	if _, ok := g.IteratedFiles[ctag]; ok {
		return
	}
	g.Tagholder = map[string]string{}
	g.Initializerholder = []FuncInitializer{}
	runHeap := []string{}
	g.Routecount = 0
	g.Stack3Holder = []Jumper{}
	inspectElement := g.GlobalHeapCopyArray[copyindex]
	if len(inspectElement.LocalObjInitIndexes) != 0 {
		lastindex := inspectElement.LocalObjInitIndexes[len(inspectElement.LocalObjInitIndexes)-1] + 1
		for len(runHeap) < lastindex {
			runHeap = append(runHeap, "")
		}
		for _, idx := range inspectElement.LocalObjInitIndexes {
			runHeap[idx] = ""
		}
	}
	if len(inspectElement.ArrayCopyIndexes) != 0 {
		lastindex := inspectElement.ArrayCopyIndexes[len(inspectElement.ArrayCopyIndexes)-1] + 1
		for len(runHeap) < lastindex {
			runHeap = append(runHeap, "")
		}
		for _, idx := range inspectElement.ArrayCopyIndexes {
			runHeap[idx] = heapin[idx]
		}
	}
	if len(inspectElement.ArgumentsIndexMapArray) != 0 {
		lastindex := inspectElement.ArgumentsIndexMapArray[len(inspectElement.ArgumentsIndexMapArray)-1] + 1
		for len(runHeap) < lastindex {
			runHeap = append(runHeap, "")
		}
		for i := 0; i < len(inspectElement.ArgumentsIndexMapArray); i++ {
			runHeap[inspectElement.ArgumentsIndexMapArray[i]] = fmt.Sprintf("arguments[%s]", strconv.Itoa(i))
		}
	}
	if inspectElement.RunSelfReference != -1 {
		for len(runHeap) < inspectElement.RunSelfReference+1 {
			runHeap = append(runHeap, "")
		}
		runHeap[inspectElement.RunSelfReference] = "\"RunSelfReferenceArgs\""
	}
	if inspectElement.MainSelfReference != -1 {
		for len(runHeap) < inspectElement.MainSelfReference+1 {
			runHeap = append(runHeap, "")
		}
		runHeap[inspectElement.MainSelfReference] = "\"MainSelfReference\""
	}

	g.walker(x, y, []string{}, runHeap, Jumper{x: 0, y: 0}, "basescript", map[string]bool{})

	outstring := g.Tagholder["basescript"]

	traversed3 := []Jumper{}
	for len(g.Stack3Holder) > 0 {
		jumperobj := g.Stack3Holder[len(g.Stack3Holder)-1]
		for _, jo := range traversed3 {
			if jumperobj.x == jo.x && jumperobj.y == jo.y {
				g.Stack3Holder = g.Stack3Holder[:len(g.Stack3Holder)-1]
				continue
			}
		}
		g.walker(jumperobj.x, jumperobj.y, []string{}, runHeap, Jumper{x: 0, y: 0}, "walk_"+strconv.Itoa(jumperobj.x)+"route3!", map[string]bool{})
		outstring += "walk_" + strconv.Itoa(jumperobj.x) + "route3!"
		g.Stack3Holder = g.Stack3Holder[:len(g.Stack3Holder)-1]
		traversed3 = append(traversed3, Jumper{x: jumperobj.x, y: jumperobj.y})
	}

	for len(routeMatch.FindStringSubmatch(outstring)) > 0 {
		for _, match := range routeMatch.FindAllStringSubmatch(outstring, -1) {
			outstring = strings.ReplaceAll(outstring, match[0], strings.Join(strings.Split(g.Tagholder[match[0]], "$"), "dollarsign"))
		}
	}
	outstring = strings.Join(strings.Split(outstring, "dollarsign"), "$")

	keyarrmatches := keyMatchRegex.FindAllStringSubmatch(outstring, -1)

	if len(keyarrmatches) > 0 {
		mapunion := make(map[string]bool)
		for _, match := range keyarrmatches {
			mapunion[match[1]] = true
		}
		for key, _ := range mapunion {
			strpavl, _ := strconv.Atoi(key)
			g.RanKeys[strpavl] = false
		}
		if len(mapunion) == 1 {
			extra := xorkeymatch.FindStringSubmatch(outstring)
			if len(extra) > 0 {
				extrakeys, _ := strconv.Atoi(extra[1])
				g.GlobalKeyArray = append(g.GlobalKeyArray, int64(extrakeys))
				mapunion[strconv.Itoa(len(g.GlobalKeyArray)-1)] = true
			}
		}
		if len(mapunion) > 1 {
			if len(mapunion) > 4 {
				for _, v := range keyarrmatches {
					g.KeyCount[v[1]]++
					g.KeyHolderRaw[100] = append(g.KeyHolderRaw[100], v[1])
					if strings.Contains(v[0], "-") {
						negint, _ := strconv.Atoi(v[1])
						g.Negatives[negint] = true
					}
				}
				for _, match := range lessThanRoundsRegex.FindAllStringSubmatch(outstring, -1) {
					if match[1] != "16" && match[1] != "0" {
						roundint, _ := strconv.Atoi(match[1])
						g.Encrounds = roundint
					}
					if g.Encrounds == 0 {
						g.Encrounds = 16
					}
				}
			} else
			if strings.Contains(outstring, "[document][querySelectorAll]") {
				g.PayloadOrder = append(g.PayloadOrder, 0)
				for k, _ := range mapunion {
					g.KeyCount[k]++
					g.KeyHolderRaw[0] = append(g.KeyHolderRaw[0], k)
				}
			} else
			if strings.Contains(outstring, "[\"stacktrace\"]") {
				payloadmap := make(map[int][]int)
				outstring = P1elseSlicer(outstring)
				check1index := strings.Index(outstring, "][\"stack\"],h")
				check2index := strings.Index(outstring, "][name],h")
				check3index := strings.Index(outstring, "][message],h")
				check4index := strings.Index(outstring, "=130^")
				check5index := strings.Index(outstring, "=131^")
				check6index := strings.Index(outstring, "=132^")
				check7index := strings.Index(outstring, "=133^")

				payloadmap[check1index] = []int{84, 121, 112, 101, 69, 114, 114, 111, 114, 58, 32, 67, 97, 110, 110, 111, 116, 32, 114, 101, 97, 100, 32, 112, 114, 111, 112, 101, 114, 116, 121, 32, 39, 48, 39, 32, 111, 102, 32, 110, 117, 108, 108, 10, 32, 32, 32, 32, 97, 116, 32, 85, 82, 76, 0}
				payloadmap[check2index] = []int{84, 121, 112, 101, 69, 114, 114, 111, 114, 0}
				payloadmap[check3index] = []int{67, 97, 110, 110, 111, 116, 32, 114, 101, 97, 100, 32, 112, 114, 111, 112, 101, 114, 116, 121, 32, 39, 48, 39, 32, 111, 102, 32, 110, 117, 108, 108, 0}
				payloadmap[check4index] = []int{130}
				payloadmap[check5index] = []int{131}
				payloadmap[check6index] = []int{132}
				payloadmap[check7index] = []int{133}
				numindex := []int{check4index, check5index, check6index, check7index}
				sort.Ints(numindex)
				for i, val := range numindex {
					if val > check1index && numindex[i-1] < check1index {
						payloadmap[numindex[i-1]] = []int{}
						break
					}
				}

				a := []int{check1index, check2index, check3index, check4index, check5index, check6index, check7index}
				sort.Ints(a)
				for _, val := range a {

					g.P1VAL = append(g.P1VAL, payloadmap[val]...)
				}

				g.PayloadOrder = append(g.PayloadOrder, 1)
				for k, _ := range mapunion {
					g.KeyCount[k]++
					g.KeyHolderRaw[1] = append(g.KeyHolderRaw[1], k)
				}
			} else
			if strings.Contains(outstring, "operaVersion") { //payload3
				//osc:= outstring
				payloadmap := make(map[int][]int)
				outstring = elseSlicer(outstring)
				doNotTrackblocks := strings.Split(outstring, "[\"doNotTrack\"]")
				operaVersionblocks := strings.Split(outstring, "[\"operaVersion\"]")
				oscpublocks := strings.Split(outstring, "[\"oscpu\"]")
				buildIDblocks := strings.Split(outstring, "[\"buildID\"]")
				cpuClassblocks := strings.Split(outstring, "[\"cpuClass\"]")
				doNotTrackval, _ := strconv.Atoi(P2SEPERATORMATCH.FindStringSubmatch(doNotTrackblocks[len(doNotTrackblocks)-1])[1])
				operaVersionval, _ := strconv.Atoi(P2SEPERATORMATCH.FindStringSubmatch(operaVersionblocks[len(operaVersionblocks)-1])[1])
				oscpuval, _ := strconv.Atoi(P2SEPERATORMATCH.FindStringSubmatch(oscpublocks[len(oscpublocks)-1])[1])
				buildIDval, _ := strconv.Atoi(P2SEPERATORMATCH.FindStringSubmatch(buildIDblocks[len(buildIDblocks)-1])[1])
				cpuClassval, _ := strconv.Atoi(P2SEPERATORMATCH.FindStringSubmatch(cpuClassblocks[len(cpuClassblocks)-1])[1])

				ColorDepthindex := strings.Index(outstring, "[\"colorDepth\"];")
				vendorindex := strings.Index(outstring, "[\"vendor\"];")
				vendorsubindex := strings.Index(outstring, "[\"vendorSub\"];")
				doNotTrackindex := strings.Index(outstring, "[\"doNotTrack\"];")
				pixelDepthindex := strings.Index(outstring, "[\"pixelDepth\"];")
				availHeightindex := strings.Index(outstring, "[\"availHeight\"];")
				outerWidthindex := strings.Index(outstring, "[\"outerWidth\"];")
				operaVersionindex := strings.Index(outstring, "[\"operaVersion\"];")
				oscpuindex := strings.Index(outstring, "[\"oscpu\"];")
				appCodeNameindex := strings.Index(outstring, "[\"appCodeName\"];")
				productindex := strings.Index(outstring, "[\"product\"];")
				widthindex := strings.Index(outstring, "[width];")
				userAgentindex := strings.Index(outstring, "[\"userAgent\"];")
				screenXindex := strings.Index(outstring, "[\"screenX\"];")
				appVersionindex := strings.Index(outstring, "[\"appVersion\"];")
				appNameindex := strings.Index(outstring, "[\"appName\"];")
				productSubindex := strings.Index(outstring, "[\"productSub\"];")
				buildIDindex := strings.Index(outstring, "[\"buildID\"];")
				hardwareConcurrencyindex := strings.Index(outstring, "[\"hardwareConcurrency\"];")
				cpuClassindex := strings.Index(outstring, "[\"cpuClass\"];")
				screenYindex := strings.Index(outstring, "[\"screenY\"];")
				availWidthindex := strings.Index(outstring, "[\"availWidth\"];")
				innerHeightindex := strings.Index(outstring, "[\"innerHeight\"];")
				innerWidthindex := strings.Index(outstring, "[\"innerWidth\"];")
				maxTouchPointsindex := strings.Index(outstring, "[\"maxTouchPoints\"];")
				devicePixelRatioindex := strings.Index(outstring, "[\"devicePixelRatio\"];")
				platformindex := strings.Index(outstring, "[\"platform\"];")
				heightindex := strings.Index(outstring, "[height];")
				outerHeightindex := strings.Index(outstring, "[\"outerHeight\"];")

				shiftval := make(map[int]int)
				locationbarindex := strings.LastIndex(outstring, "[\"locationbar\"];")
				toolbarindex := strings.LastIndex(outstring, "[\"toolbar\"];")
				webdriverindex := strings.LastIndex(outstring, "[\"webdriver\"];")
				isSecureContextindex := strings.LastIndex(outstring, "[\"isSecureContext\"];")
				shiftval[locationbarindex] = 1
				shiftval[toolbarindex] = 1
				shiftval[webdriverindex] = 0
				shiftval[isSecureContextindex] = 1
				th := []int{locationbarindex, toolbarindex, webdriverindex, isSecureContextindex}
				sort.Ints(th)
				sv := 0
				lastshiftindex := th[3]
				for i, val := range th {
					sv |= shiftval[val] << i
				}

				payloadmap[lastshiftindex] = []int{sv}
				payloadmap[ColorDepthindex] = []int{24}
				payloadmap[vendorindex] = []int{71, 111, 111, 103, 108, 101, 32, 73, 110, 99, 46, 0}
				payloadmap[vendorsubindex] = []int{0}
				payloadmap[doNotTrackindex] = []int{doNotTrackval}
				payloadmap[pixelDepthindex] = []int{24}

				payloadmap[outerWidthindex] = []int{56, 48}
				payloadmap[operaVersionindex] = []int{operaVersionval}
				payloadmap[oscpuindex] = []int{oscpuval}
				payloadmap[appCodeNameindex] = []int{77, 111, 122, 105, 108, 108, 97, 0}
				payloadmap[productindex] = []int{71, 101, 99, 107, 111, 0}

				payloadmap[userAgentindex] = []int{77, 111, 122, 105, 108, 108, 97, 47, 53, 46, 48, 32, 40, 87, 105, 110, 100, 111, 119, 115, 32, 78, 84, 32, 49, 48, 46, 48, 59, 32, 87, 105, 110, 54, 52, 59, 32, 120, 54, 52, 41, 32, 65, 112, 112, 108, 101, 87, 101, 98, 75, 105, 116, 47, 53, 51, 55, 46, 51, 54, 32, 40, 75, 72, 84, 77, 76, 44, 32, 108, 105, 107, 101, 32, 71, 101, 99, 107, 111, 41, 32, 67, 104, 114, 111, 109, 101, 47, 57, 49, 46, 48, 46, 52, 52, 55, 50, 46, 49, 50, 51, 32, 83, 97, 102, 97, 114, 105, 47, 53, 51, 55, 46, 51, 54, 0}
				payloadmap[screenXindex] = []int{45, 1}
				payloadmap[appVersionindex] = []int{53, 46, 48, 32, 40, 87, 105, 110, 100, 111, 119, 115, 32, 78, 84, 32, 49, 48, 46, 48, 59, 32, 87, 105, 110, 54, 52, 59, 32, 120, 54, 52, 41, 32, 65, 112, 112, 108, 101, 87, 101, 98, 75, 105, 116, 47, 53, 51, 55, 46, 51, 54, 32, 40, 75, 72, 84, 77, 76, 44, 32, 108, 105, 107, 101, 32, 71, 101, 99, 107, 111, 41, 32, 67, 104, 114, 111, 109, 101, 47, 57, 49, 46, 48, 46, 52, 52, 55, 50, 46, 49, 50, 51, 32, 83, 97, 102, 97, 114, 105, 47, 53, 51, 55, 46, 51, 54, 0}
				payloadmap[appNameindex] = []int{78, 101, 116, 115, 99, 97, 112, 101, 0}
				payloadmap[productSubindex] = []int{50, 48, 48, 51, 48, 49, 48, 55, 0}
				payloadmap[buildIDindex] = []int{buildIDval}
				payloadmap[hardwareConcurrencyindex] = []int{12}
				payloadmap[cpuClassindex] = []int{cpuClassval}
				payloadmap[screenYindex] = []int{0}

				payloadmap[innerHeightindex] = []int{38, 24}
				payloadmap[innerWidthindex] = []int{41, 11}
				payloadmap[maxTouchPointsindex] = []int{0}
				payloadmap[devicePixelRatioindex] = []int{128, 0, 0, 0, 192, 30, 133, 243, 63}
				payloadmap[platformindex] = []int{87, 105, 110, 51, 50, 0}

				payloadmap[outerHeightindex] = []int{59, 27}
				//1040 40,32
				//1920 60,32
				//1920 60,32
				//1080 56,31
				if rand.Intn(2) == 1{
					payloadmap[availHeightindex] = []int{54, 27}
					payloadmap[widthindex] = []int{44, 98}
					payloadmap[availWidthindex] = []int{57, 96}
					payloadmap[heightindex] = []int{54, 27}
				}else{
					payloadmap[availHeightindex] = []int{40,32}
					payloadmap[widthindex] = []int{60,32}
					payloadmap[availWidthindex] = []int{60,32}
					payloadmap[heightindex] = []int{56,31}
				}

				tholder := []int{doNotTrackindex, operaVersionindex, oscpuindex, buildIDindex, cpuClassindex}
				sort.Ints(tholder)
				payloadmap[tholder[len(tholder)-1]] = append(payloadmap[tholder[len(tholder)-1]], 1)

				a := []int{lastshiftindex, vendorindex, ColorDepthindex, doNotTrackindex, vendorsubindex, pixelDepthindex, availHeightindex, outerWidthindex, operaVersionindex, oscpuindex, appCodeNameindex, productindex, widthindex, userAgentindex, screenXindex, appVersionindex, appNameindex, productSubindex, buildIDindex, hardwareConcurrencyindex, cpuClassindex, screenYindex, availWidthindex, innerHeightindex, innerWidthindex, maxTouchPointsindex, devicePixelRatioindex, platformindex, heightindex, outerHeightindex}
				sort.Ints(a)
				for _, val := range a {
					g.P2VAL = append(g.P2VAL, payloadmap[val]...)
				}

				g.PayloadOrder = append(g.PayloadOrder, 2)
				for k, _ := range mapunion {
					g.KeyCount[k]++
					g.KeyHolderRaw[2] = append(g.KeyHolderRaw[2], k)
				}
			} else
			if strings.Contains(outstring, "video/webm; codecs=") && !strings.Contains(outstring, "\"Andale Mono\"") { //payload4
				payloadmap := make(map[int]int)
				check0index := strings.Index(outstring, "[\"video/mp4; codecs=\"hvc1x\"\"]")
				check1index := strings.Index(outstring, "[\"audio/webz\"]")
				check2index := strings.Index(outstring, "[\"video/mp4; codecs=\"avc3.42001E\"\"]")
				check3index := strings.Index(outstring, "[\"video/mp4; codecs=\"avc1.42E034\"\"]")
				check4index := strings.Index(outstring, "[\"audio/webm; codecs=\"vp8\"\"]")
				check5index := strings.Index(outstring, "[\"audio/mpeg; codecs=\"mp3\"\"]")
				check6index := strings.Index(outstring, "[\"video/mp4; codecs=\"flac\"\"]")
				check7index := strings.Index(outstring, "[\"audio/x-m4a; codecs=\"vp9, mp4a.40.2\"\"]")
				check8index := strings.Index(outstring, "[\"audio/flac\"]")
				check9index := strings.Index(outstring, "[\"video/webm\"]")
				check10index := strings.Index(outstring, "[\"video/mp4; codecs=\"avc3.42E01E, mp4a.40.29\"\"]")
				check11index := strings.Index(outstring, "[\"audio/x-m4a; codecs=\"vp8, mp4a.40\"\"]")
				check12index := strings.Index(outstring, "[\"audio/x-m4a; codecs=\"mp3\"\"]")
				check13index := strings.Index(outstring, "[\"video/mp4; codecs=\"ac-3\"\"]")
				check14index := strings.Index(outstring, "[\"video/mp4; codecs=\"avc1.42F01E\"\"]")
				check15index := strings.Index(outstring, "[\"video/mp4; codecs=\"hev1\"\"]")
				check16index := strings.Index(outstring, "[\"video/ogg; codecs=\"vp8\"\"]")
				check17index := strings.Index(outstring, "[\"video/mp4; codecs=\"mp4a.67\"\"]")
				check18index := strings.Index(outstring, "[\"video/x-\"]")
				check19index := strings.Index(outstring, "[\"video/webm; codecs=\"vp09.02.10.08\"\"]")
				check20index := strings.Index(outstring, "[\"video/x-m4v; codecs=\"avc1.42AC23\"\"]")
				check21index := strings.Index(outstring, "[\"video/mp4; codecs=\"opus\"\"]")
				check22index := strings.Index(outstring, "[\"audio/aac; codecs=\"flac\"\"]")
				check23index := strings.Index(outstring, "[\"video/mp4; codecs=\"lavc1337\"\"]")
				check24index := strings.Index(outstring, "[\"audio/mpeg; codecs=\"vp9\"\"]")
				check25index := strings.Index(outstring, "[\"video/mp4; codecs=\"avc1.42E009\"\"]")
				check26index := strings.Index(outstring, "[\"video/webm; codecs=\"av01.0.04M.08\"\"]")
				check27index := strings.Index(outstring, "[\"video/mp4; codecs=\"avc1.42011E\"\"]")
				check28index := strings.Index(outstring, "[\"video/mp4\"]")
				check29index := strings.Index(outstring, "[\"video/mp4; codecs=\"mp4a.40.5\"\"]")
				check30index := strings.Index(outstring, "[\"video/ogg; codecs=\"flac\"\"]")
				check31index := strings.Index(outstring, "[\"video/ogg; codecs=\"opus\"\"]")
				check32index := strings.Index(outstring, "[\"video/mp4; codecs=\"avc1.4D001E\"\"]")
				check33index := strings.Index(outstring, "[\"video/mp4; codecs=\"mp4a.40.02\"\"]")
				payloadmap[check0index] = 0
				payloadmap[check1index] = 0
				payloadmap[check2index] = 2
				payloadmap[check3index] = 2
				payloadmap[check4index] = 0
				payloadmap[check5index] = 2
				payloadmap[check6index] = 2
				payloadmap[check7index] = 0
				payloadmap[check8index] = 2
				payloadmap[check9index] = 1
				payloadmap[check10index] = 2
				payloadmap[check11index] = 0
				payloadmap[check12index] = 0
				payloadmap[check13index] = 0
				payloadmap[check14index] = 2
				payloadmap[check15index] = 0
				payloadmap[check16index] = 2
				payloadmap[check17index] = 2
				payloadmap[check18index] = 0
				payloadmap[check19index] = 2
				payloadmap[check20index] = 1
				payloadmap[check21index] = 2
				payloadmap[check22index] = 0
				payloadmap[check23index] = 0
				payloadmap[check24index] = 0
				payloadmap[check25index] = 1
				payloadmap[check26index] = 2
				payloadmap[check27index] = 0
				payloadmap[check28index] = 1
				payloadmap[check29index] = 2
				payloadmap[check30index] = 2
				payloadmap[check31index] = 2
				payloadmap[check32index] = 2
				payloadmap[check33index] = 2

				a := []int{check0index, check1index, check2index, check3index, check4index, check5index, check6index, check7index, check8index, check9index, check10index, check11index, check12index, check13index, check14index, check15index, check16index, check17index, check18index, check19index, check20index, check21index, check22index, check23index, check24index, check25index, check26index, check27index, check28index, check29index, check30index, check31index, check32index, check33index}
				sort.Ints(a)
				for _, val := range a {
					g.P3VAL = append(g.P3VAL, payloadmap[val])
				}

				g.PayloadOrder = append(g.PayloadOrder, 3)
				for k, _ := range mapunion {
					g.KeyCount[k]++
					g.KeyHolderRaw[3] = append(g.KeyHolderRaw[3], k)
				}
			} else
			if strings.Contains(outstring, "createEvent") && strings.Contains(outstring, "isTrusted") { //payload5
				g.P4VAL, _ = strconv.Atoi(P4MATCH.FindAllStringSubmatch(outstring, -1)[3][1])
				g.PayloadOrder = append(g.PayloadOrder, 4)
				for k, _ := range mapunion {
					g.KeyCount[k]++
					g.KeyHolderRaw[4] = append(g.KeyHolderRaw[4], k)
				}
			} else
			if strings.Contains(starttag, "[Error]") { //payload6
				g.P5VAL, _ = strconv.Atoi(P2SEPERATORMATCH.FindStringSubmatch(outstring)[1])
				g.PayloadOrder = append(g.PayloadOrder, 5)
				for k, _ := range mapunion {
					g.KeyCount[k]++
					g.KeyHolderRaw[5] = append(g.KeyHolderRaw[5], k)
				}
			} else
			if strings.Contains(starttag, "[global][\"navigator\"]") && !strings.Contains(starttag, "[\"plugins\"]") { //payload7
				g.PayloadOrder = append(g.PayloadOrder, 6)
				for k, _ := range mapunion {
					g.KeyCount[k]++
					g.KeyHolderRaw[6] = append(g.KeyHolderRaw[6], k)
				}
			} else
			if strings.Contains(starttag, "[global][document]") && strings.Contains(outstring, "[map]") { //payload8
				g.PayloadOrder = append(g.PayloadOrder, 7)
				for k, _ := range mapunion {
					g.KeyCount[k]++
					g.KeyHolderRaw[7] = append(g.KeyHolderRaw[7], k)
				}
			} else
			//8 empty
			if strings.Contains(starttag, "-Infinity") { //payload10
				g.PayloadOrder = append(g.PayloadOrder, 8)
				for k, _ := range mapunion {
					g.KeyCount[k]++
					g.KeyHolderRaw[8] = append(g.KeyHolderRaw[8], k)
				}
			} else
			if strings.Contains(starttag, "[Math][round][bind]") { //payload10
				g.PayloadOrder = append(g.PayloadOrder, 9)
				for k, _ := range mapunion {
					g.KeyCount[k]++
					g.KeyHolderRaw[9] = append(g.KeyHolderRaw[9], k)
				}
			} else
			if strings.Contains(outstring, "documentElement") { //payload11
				g.PayloadOrder = append(g.PayloadOrder, 10)
				for k, _ := range mapunion {
					g.KeyCount[k]++
					g.KeyHolderRaw[10] = append(g.KeyHolderRaw[10], k)
				}
			} else
			if strings.Contains(outstring, "\"version\"") && strings.Contains(starttag, "[global][document]") && !strings.Contains(outstring, "ableToCollect") { //payload12
				g.P11VAL = []int{47, 1}
				if strings.Index(outstring, "if(0") > -1 {
					g.P11VAL = append(g.P11VAL, 12, 160, 18, 0, 176, 32)
				} else {
					g.P11VAL = append(g.P11VAL, 130, 6, 0, 164, 2, 24)
				}
				g.P11VAL = append(g.P11VAL, 5)
				g.PayloadOrder = append(g.PayloadOrder, 11)
				for k, _ := range mapunion {
					g.KeyCount[k]++
					g.KeyHolderRaw[11] = append(g.KeyHolderRaw[11], k)
				}
			} else
			if len(testregex1.FindStringSubmatch(outstring)) > 0 { //payload13
				g.PayloadOrder = append(g.PayloadOrder, 12)
				for k, _ := range mapunion {
					g.KeyCount[k]++
					g.KeyHolderRaw[12] = append(g.KeyHolderRaw[12], k)
				}
			} else
			if strings.Contains(outstring, "[\"events\"];") { //payload14
				outstring = elseSlicer(outstring)
				payloadmap := make(map[int]int)
				check1index := strings.Index(outstring, "[\"events\"]")
				check2index := strings.Index(outstring, "<< 1")
				check3index := strings.Index(outstring, "[\"visibilityEventCount\"]")
				payloadmap[check1index] = 0
				payloadmap[check2index] = 3
				payloadmap[check3index] = 0
				a := []int{check1index, check2index, check3index}
				sort.Ints(a)
				for _, val := range a {
					g.P13VAL = append(g.P13VAL, payloadmap[val])
				}
				g.PayloadOrder = append(g.PayloadOrder, 13)
				for k, _ := range mapunion {
					g.KeyCount[k]++
					g.KeyHolderRaw[13] = append(g.KeyHolderRaw[13], k)
				}
			} else
			if strings.Contains(outstring, "keyCode") { //payload15
				var shiftval int
				if strings.Contains(outstring, "[\"instanceOfUIEvent\"]|0") {
					shiftval = 2
				} else {
					shiftval = 1
				}
				check1index := strings.Index(outstring, "[\"timestamp\"]")
				check2index := strings.Index(outstring, "[\"eventType\"]")
				check3index := strings.Index(outstring, "[\"keyCode\"]")
				check4index := strings.Index(outstring, "<< 1")
				check5index := strings.Index(outstring, "[\"targetId\"]")
				check6index := strings.Index(outstring, "[\"targetName\"]")
				check7index := strings.Index(outstring, "[\"modifierKeys\"]")
				if check5index > check6index {
					holder := check5index
					check5index = check6index
					check6index = holder
				}

				a := []int{check1index, check2index, check3index, check4index, check5index, check6index, check7index}
				sort.Ints(a)
				float1 := 1000 + rand.Float64()*500
				float2 := float1 + 5 + rand.Float64()*10
				float3 := float2 + 5 + rand.Float64()*10
				float4 := float3 + 5 + rand.Float64()*10
				g.P14VAL = append(g.P14VAL, 4)
				timestamparr := []float64{float1, float2, float3, float4}

				if P15IFMATCH.FindStringSubmatch(outstring)[1] != "0" {
					for i := 0; i < 4; i++ {
						for _, val := range a {
							switch val {
							case check1index:
								g.P14VAL = append(g.P14VAL, generation.BaseNumEnc(timestamparr[i])...)
							case check2index:
								if i%2 == 0 {
									g.P14VAL = append(g.P14VAL, 1)
								} else {
									g.P14VAL = append(g.P14VAL, 2)
								}
							case check3index:
								g.P14VAL = append(g.P14VAL, 1)
							case check4index:
								g.P14VAL = append(g.P14VAL, shiftval)
							case check5index:
								switch i {
								case 0:
									g.P14VAL = append(g.P14VAL, 117, 115, 101, 114, 110, 97, 109, 101, 0)
								case 1:
									g.P14VAL = append(g.P14VAL, 129, 0)
								case 2:
									g.P14VAL = append(g.P14VAL, 112, 97, 115, 115, 119, 111, 114, 100, 0)
								case 3:
									g.P14VAL = append(g.P14VAL, 129, 1)
								}
							case check6index:
								switch i {
								case 0:
									g.P14VAL = append(g.P14VAL, 129, 0)
								case 1:
									g.P14VAL = append(g.P14VAL, 129, 0)
								case 2:
									g.P14VAL = append(g.P14VAL, 129, 1)
								case 3:
									g.P14VAL = append(g.P14VAL, 129, 1)
								}
							case check7index:
								g.P14VAL = append(g.P14VAL, 4, 0)
							}
						}
					}
				} else {
					for i := 3; i >= 0; i-- {
						for _, val := range a {
							switch val {
							case check1index:
								g.P14VAL = append(g.P14VAL, generation.BaseNumEnc(timestamparr[i])...)
							case check2index:
								if i%2 == 0 {
									g.P14VAL = append(g.P14VAL, 1)
								} else {
									g.P14VAL = append(g.P14VAL, 2)
								}
							case check3index:
								g.P14VAL = append(g.P14VAL, 1)
							case check4index:
								g.P14VAL = append(g.P14VAL, shiftval)
							case check5index:
								switch i {
								case 1:
									g.P14VAL = append(g.P14VAL, 117, 115, 101, 114, 110, 97, 109, 101, 0)
								case 0:
									g.P14VAL = append(g.P14VAL, 129, 1)
								case 3:
									g.P14VAL = append(g.P14VAL, 112, 97, 115, 115, 119, 111, 114, 100, 0)
								case 2:
									g.P14VAL = append(g.P14VAL, 129, 0)
								}
							case check6index:
								switch i {
								case 3:
									g.P14VAL = append(g.P14VAL, 129, 0)
								case 2:
									g.P14VAL = append(g.P14VAL, 129, 0)
								case 1:
									g.P14VAL = append(g.P14VAL, 129, 1)
								case 0:
									g.P14VAL = append(g.P14VAL, 129, 1)
								}
							case check7index:
								g.P14VAL = append(g.P14VAL, 4, 0)
							}
						}
					}
				}

				g.PayloadOrder = append(g.PayloadOrder, 14)
				for k, _ := range mapunion {
					g.KeyCount[k]++
					g.KeyHolderRaw[14] = append(g.KeyHolderRaw[14], k)
				}
			} else
			if strings.Contains(outstring, "\"pixels\"") && strings.Contains(outstring, "\"count\"") { //payload16
				var pixelbytes []int
				var mval string
				apixelindex := apixelmatch.FindStringSubmatch(outstring)[1]
				outstring = strings.ReplaceAll(outstring, fmt.Sprintf("heapin[%s]", apixelindex), "objectA")
				ifblocks := strings.Split(outstring, "objectA[\"pixels\"]")
				if len(ifblocks[2]) > len(ifblocks[4]) {
					mval = P15WHILEMATCH.FindStringSubmatch(ifblocks[2])[1]
				} else {
					mval = P15WHILEMATCH.FindStringSubmatch(ifblocks[4])[1]
				}
				if mval == "0" {
					pixelbytes = []int{63, 7, 0, 0, 0}
				} else {
					pixelbytes = []int{0, 0, 0, 63, 7}
				}
				check1index := strings.Index(outstring, "[\"hash\"]")
				check2index := strings.Index(outstring, "objectA[\"count\"]")
				check3index := strings.Index(outstring, "objectA[\"pixels\"]")
				check4index := strings.Index(outstring, "][\"count\"]")
				check5index := strings.Index(outstring, "][\"pixels\"]")
				payloadmap := make(map[int][]int)
				payloadmap[check1index] = []int{123, 183, 245, 163, 1}
				payloadmap[check2index] = []int{1}
				payloadmap[check3index] = append([]int{1, 4}, pixelbytes...)
				payloadmap[check4index] = []int{0}
				payloadmap[check5index] = []int{0}
				a := []int{check1index, check2index, check3index, check4index, check5index}
				sort.Ints(a)
				g.P15VAL = []int{4}
				for i := 0; i < 4; i++ {
					for _, val := range a {
						g.P15VAL = append(g.P15VAL, payloadmap[val]...)
					}
				}

				g.PayloadOrder = append(g.PayloadOrder, 15)
				for k, _ := range mapunion {
					g.KeyCount[k]++
					g.KeyHolderRaw[15] = append(g.KeyHolderRaw[15], k)
				}

			} else
			if strings.Contains(outstring, "[document][\"referrer\"]") { //payload17
				payloadmap := make(map[int]int)
				sm := P31VALMATCH.FindAllStringSubmatch(outstring, -1)[1]
				xorval, _ := strconv.Atoi(sm[1])
				smindex := strings.Index(outstring, sm[0])
				lengthindex := strings.LastIndex(outstring, "[\"history\"][length]")
				payloadmap[smindex] = xorval
				payloadmap[lengthindex] = 3

				a := []int{smindex, lengthindex}
				sort.Ints(a)
				for _, val := range a {
					g.P16VAL = append(g.P16VAL, payloadmap[val])
				}

				g.PayloadOrder = append(g.PayloadOrder, 16)
				for k, _ := range mapunion {
					g.KeyCount[k]++
					g.KeyHolderRaw[16] = append(g.KeyHolderRaw[16], k)
				}
			} else
			if strings.Contains(outstring, "ableToCollect") && !strings.Contains(outstring, "\"Andale Mono\"") { //payload18
				outstring = elseSlicer(outstring)
				payloadmap := make(map[int][]int)
				antialiasstring := strings.Split(outstring, "antialias")[1]
				antialiasval, _ := strconv.Atoi(P17HEAPSETMATCH.FindStringSubmatch(antialiasstring)[1])

				debugobjectmatch := P17DEBUGMATCH.FindStringSubmatch(outstring)[1]
				paramobjectmatch := P17PARAMMATCH.FindStringSubmatch(outstring)[1]

				debugrenderindex := strings.LastIndex(outstring, fmt.Sprintf("heapin[%s][\"renderer\"]", debugobjectmatch))
				debugvendorindex := strings.LastIndex(outstring, fmt.Sprintf("heapin[%s][\"vendor\"]", debugobjectmatch))
				paramrenderindex := strings.LastIndex(outstring, fmt.Sprintf("heapin[%s][\"renderer\"]", paramobjectmatch))
				paramvendorindex := strings.LastIndex(outstring, fmt.Sprintf("heapin[%s][\"vendor\"]", paramobjectmatch))
				vendorholders := [][]int{
					{71,111,111,103,108,101,32,73,110,99,46,32,40,73,110,116,101,108,41,0},
					{71,111,111,103,108,101,32,73,110,99,46,32,40,65,77,68,41,0},
					{71,111,111,103,108,101,32,73,110,99,46,32,40,65,77,68,41,0},
					{71,111,111,103,108,101,32,73,110,99,46,0},
					{71,111,111,103,108,101,32,73,110,99,46,32,40,78,86,73,68,73,65,41,0},
					{71,111,111,103,108,101,32,73,110,99,46,32,40,78,86,73,68,73,65,41,0},
					{71,111,111,103,108,101,32,73,110,99,46,32,40,78,86,73,68,73,65,41,0},
				}
				renderrholders := [][]int{
					{65,78,71,76,69,32,40,73,110,116,101,108,44,32,73,110,116,101,108,40,82,41,32,85,72,68,32,71,114,97,112,104,105,99,115,32,54,51,48,32,68,105,114,101,99,116,51,68,49,49,32,118,115,95,53,95,48,32,112,115,95,53,95,48,44,32,68,51,68,49,49,45,50,54,46,50,48,46,49,48,48,46,55,50,54,49,41,0},
					{65,78,71,76,69,32,40,65,77,68,44,32,82,97,100,101,111,110,32,82,88,32,53,55,48,32,83,101,114,105,101,115,32,68,105,114,101,99,116,51,68,49,49,32,118,115,95,53,95,48,32,112,115,95,53,95,48,44,32,68,51,68,49,49,45,50,55,46,50,48,46,49,50,48,51,51,46,50,48,48,55,41,0},
					{65,78,71,76,69,32,40,65,77,68,44,32,65,77,68,32,82,97,100,101,111,110,40,84,77,41,32,82,88,32,86,101,103,97,32,49,49,32,71,114,97,112,104,105,99,115,32,68,105,114,101,99,116,51,68,49,49,32,118,115,95,53,95,48,32,112,115,95,53,95,48,44,32,68,51,68,49,49,45,50,55,46,50,48,46,49,48,50,54,46,49,48,48,51,41,0},
					{71,111,111,103,108,101,32,83,119,105,102,116,83,104,97,100,101,114,0},
					{65,78,71,76,69,32,40,78,86,73,68,73,65,44,32,78,86,73,68,73,65,32,71,101,70,111,114,99,101,32,71,84,88,32,55,53,48,32,68,105,114,101,99,116,51,68,49,49,32,118,115,95,53,95,48,32,112,115,95,53,95,48,44,32,68,51,68,49,49,45,50,55,46,50,49,46,49,52,46,54,49,57,50,41,0},
					{65,78,71,76,69,32,40,78,86,73,68,73,65,44,32,78,86,73,68,73,65,32,71,101,70,111,114,99,101,32,71,84,88,32,49,48,53,48,32,84,105,32,68,105,114,101,99,116,51,68,49,49,32,118,115,95,53,95,48,32,112,115,95,53,95,48,44,32,68,51,68,49,49,45,50,55,46,50,49,46,49,52,46,53,54,55,49,41,0},
					{65,78,71,76,69,32,40,78,86,73,68,73,65,44,32,78,86,73,68,73,65,32,71,101,70,111,114,99,101,32,71,84,88,32,49,54,54,48,32,83,85,80,69,82,32,68,105,114,101,99,116,51,68,49,49,32,118,115,95,53,95,48,32,112,115,95,53,95,48,44,32,68,51,68,49,49,45,50,55,46,50,49,46,49,52,46,53,54,55,49,41,0},
				}
				payloadmap[paramrenderindex] = []int{87, 101, 98, 75, 105, 116, 32, 87, 101, 98, 71, 76, 0}
				payloadmap[paramvendorindex] = []int{87, 101, 98, 75, 105, 116, 0}
				antialiasIndex := strings.Index(outstring, "[\"antialias\"]")
				shaderPrecisionsIndex := strings.Index(outstring, "[\"shaderPrecisions\"]")
				maxTextureSizeIndex := strings.Index(outstring, "[\"maxTextureSize\"]")
				maxVertexUniformVectorsIndex := strings.Index(outstring, "[\"maxVertexUniformVectors\"]")
				maxFragmentUniformVectorsIndex := strings.Index(outstring, "[\"maxFragmentUniformVectors\"]")
				maxVertexTextureImageUnitsIndex := strings.Index(outstring, "[\"maxVertexTextureImageUnits\"]")
				contextPropertiesIndex := strings.Index(outstring, "[\"contextProperties\"]")
				supportedExtensionsIndex := strings.Index(outstring, "[\"supportedExtensions\"]")
				blueBitsIndex := strings.Index(outstring, "[\"blueBits\"]")
				depthBitsIndex := strings.Index(outstring, "[\"depthBits\"]")
				greenBitsIndex := strings.Index(outstring, "[\"greenBits\"]")
				redBitsIndex := strings.Index(outstring, "[\"redBits\"]")
				stencilBitsIndex := strings.Index(outstring, "[\"stencilBits\"]")
				maxAnisotropyIndex := strings.Index(outstring, "[\"maxAnisotropy\"]")
				maxRenderbufferSizeIndex := strings.Index(outstring, "[\"maxRenderbufferSize\"]")
				maxVertexAttribsIndex := strings.Index(outstring, "[\"maxVertexAttribs\"]")
				versionIndex := strings.Index(outstring, "[\"version\"]")
				shadingLanguageVersionIndex := strings.Index(outstring, "[\"shadingLanguageVersion\"]")
				maxVaryingVectorsIndex := strings.Index(outstring, "[\"maxVaryingVectors\"]")
				maxCubeMapTextureSizeIndex := strings.Index(outstring, "[\"maxCubeMapTextureSize\"]")
				alphaBitsIndex := strings.Index(outstring, "[\"alphaBits\"]")
				maxCombinedTextureImageUnitsIndex := strings.Index(outstring, "[\"maxCombinedTextureImageUnits\"]")
				dimensionsIndex := strings.Index(outstring, "[\"dimensions\"]")
				maxTextureImageUnitsIndex := strings.Index(outstring, "[\"maxTextureImageUnits\"]")
				payloadmap[antialiasIndex] = []int{antialiasval}
				payloadmap[maxTextureSizeIndex] = []int{32, 128, 4}
				payloadmap[maxVertexUniformVectorsIndex] = []int{63, 127}
				payloadmap[maxFragmentUniformVectorsIndex] = []int{32, 32}
				payloadmap[maxVertexTextureImageUnitsIndex] = []int{16}
				payloadmap[contextPropertiesIndex] = []int{103, 133, 132, 162, 14}
				payloadmap[blueBitsIndex] = []int{8}
				payloadmap[depthBitsIndex] = []int{24}
				payloadmap[greenBitsIndex] = []int{8}
				payloadmap[redBitsIndex] = []int{8}
				payloadmap[stencilBitsIndex] = []int{0}
				payloadmap[maxAnisotropyIndex] = []int{16}
				payloadmap[maxRenderbufferSizeIndex] = []int{32, 128, 4}
				payloadmap[maxVertexAttribsIndex] = []int{16}
				payloadmap[versionIndex] = []int{87, 101, 98, 71, 76, 32, 49, 46, 48, 32, 40, 79, 112, 101, 110, 71, 76, 32, 69, 83, 32, 50, 46, 48, 32, 67, 104, 114, 111, 109, 105, 117, 109, 41, 0}
				payloadmap[shadingLanguageVersionIndex] = []int{87, 101, 98, 71, 76, 32, 71, 76, 83, 76, 32, 69, 83, 32, 49, 46, 48, 32, 40, 79, 112, 101, 110, 71, 76, 32, 69, 83, 32, 71, 76, 83, 76, 32, 69, 83, 32, 49, 46, 48, 32, 67, 104, 114, 111, 109, 105, 117, 109, 41, 0}
				payloadmap[maxVaryingVectorsIndex] = []int{30}
				payloadmap[maxCubeMapTextureSizeIndex] = []int{32, 128, 4}
				payloadmap[alphaBitsIndex] = []int{8}
				payloadmap[maxCombinedTextureImageUnitsIndex] = []int{32, 1}
				payloadmap[dimensionsIndex] = []int{3, 49, 49, 0, 49, 49, 48, 50, 52, 0, 51, 50, 55, 54, 55, 51, 50, 55, 54, 55, 0}
				payloadmap[maxTextureImageUnitsIndex] = []int{16}

				th := []int{shaderPrecisionsIndex, supportedExtensionsIndex}
				sort.Ints(th)
				lastblock := outstring
				for i := 1; i >= 0; i-- {
					val := th[i]
					switch val {
					case shaderPrecisionsIndex:
						blocks := strings.Split(lastblock, "shaderPrecisions")
						lastblock = blocks[0]
						if strings.Index(blocks[len(blocks)-1], "while(") > -1 {
							payloadmap[shaderPrecisionsIndex] = []int{12, 50, 51, 49, 50, 55, 49, 50, 55, 0, 129, 2, 129, 2, 129, 2, 129, 2, 129, 2, 48, 51, 49, 51, 48, 0, 129, 3, 129, 3, 129, 3, 129, 3, 129, 3}
						} else {
							payloadmap[shaderPrecisionsIndex] = []int{12, 48, 51, 49, 51, 48, 0, 129, 1, 129, 1, 129, 1, 129, 1, 129, 1, 50, 51, 49, 50, 55, 49, 50, 55, 0, 129, 2, 129, 2, 129, 2, 129, 2, 129, 2}
						}
					case supportedExtensionsIndex:
						blocks := strings.Split(lastblock, "supportedExtensions")
						lastblock = blocks[0]
						if strings.Index(blocks[len(blocks)-1], "while(") > -1 {
							payloadmap[supportedExtensionsIndex] = []int{33, 1, 65, 78, 71, 76, 69, 95, 105, 110, 115, 116, 97, 110, 99, 101, 100, 95, 97, 114, 114, 97, 121, 115, 0, 69, 88, 84, 95, 98, 108, 101, 110, 100, 95, 109, 105, 110, 109, 97, 120, 0, 69, 88, 84, 95, 99, 111, 108, 111, 114, 95, 98, 117, 102, 102, 101, 114, 95, 104, 97, 108, 102, 95, 102, 108, 111, 97, 116, 0, 69, 88, 84, 95, 100, 105, 115, 106, 111, 105, 110, 116, 95, 116, 105, 109, 101, 114, 95, 113, 117, 101, 114, 121, 0, 69, 88, 84, 95, 102, 108, 111, 97, 116, 95, 98, 108, 101, 110, 100, 0, 69, 88, 84, 95, 102, 114, 97, 103, 95, 100, 101, 112, 116, 104, 0, 69, 88, 84, 95, 115, 104, 97, 100, 101, 114, 95, 116, 101, 120, 116, 117, 114, 101, 95, 108, 111, 100, 0, 69, 88, 84, 95, 116, 101, 120, 116, 117, 114, 101, 95, 99, 111, 109, 112, 114, 101, 115, 115, 105, 111, 110, 95, 98, 112, 116, 99, 0, 69, 88, 84, 95, 116, 101, 120, 116, 117, 114, 101, 95, 99, 111, 109, 112, 114, 101, 115, 115, 105, 111, 110, 95, 114, 103, 116, 99, 0, 69, 88, 84, 95, 116, 101, 120, 116, 117, 114, 101, 95, 102, 105, 108, 116, 101, 114, 95, 97, 110, 105, 115, 111, 116, 114, 111, 112, 105, 99, 0, 87, 69, 66, 75, 73, 84, 95, 69, 88, 84, 95, 116, 101, 120, 116, 117, 114, 101, 95, 102, 105, 108, 116, 101, 114, 95, 97, 110, 105, 115, 111, 116, 114, 111, 112, 105, 99, 0, 69, 88, 84, 95, 115, 82, 71, 66, 0, 75, 72, 82, 95, 112, 97, 114, 97, 108, 108, 101, 108, 95, 115, 104, 97, 100, 101, 114, 95, 99, 111, 109, 112, 105, 108, 101, 0, 79, 69, 83, 95, 101, 108, 101, 109, 101, 110, 116, 95, 105, 110, 100, 101, 120, 95, 117, 105, 110, 116, 0, 79, 69, 83, 95, 102, 98, 111, 95, 114, 101, 110, 100, 101, 114, 95, 109, 105, 112, 109, 97, 112, 0, 79, 69, 83, 95, 115, 116, 97, 110, 100, 97, 114, 100, 95, 100, 101, 114, 105, 118, 97, 116, 105, 118, 101, 115, 0, 79, 69, 83, 95, 116, 101, 120, 116, 117, 114, 101, 95, 102, 108, 111, 97, 116, 0, 79, 69, 83, 95, 116, 101, 120, 116, 117, 114, 101, 95, 102, 108, 111, 97, 116, 95, 108, 105, 110, 101, 97, 114, 0, 79, 69, 83, 95, 116, 101, 120, 116, 117, 114, 101, 95, 104, 97, 108, 102, 95, 102, 108, 111, 97, 116, 0, 79, 69, 83, 95, 116, 101, 120, 116, 117, 114, 101, 95, 104, 97, 108, 102, 95, 102, 108, 111, 97, 116, 95, 108, 105, 110, 101, 97, 114, 0, 79, 69, 83, 95, 118, 101, 114, 116, 101, 120, 95, 97, 114, 114, 97, 121, 95, 111, 98, 106, 101, 99, 116, 0, 87, 69, 66, 71, 76, 95, 99, 111, 108, 111, 114, 95, 98, 117, 102, 102, 101, 114, 95, 102, 108, 111, 97, 116, 0, 87, 69, 66, 71, 76, 95, 99, 111, 109, 112, 114, 101, 115, 115, 101, 100, 95, 116, 101, 120, 116, 117, 114, 101, 95, 115, 51, 116, 99, 0, 87, 69, 66, 75, 73, 84, 95, 87, 69, 66, 71, 76, 95, 99, 111, 109, 112, 114, 101, 115, 115, 101, 100, 95, 116, 101, 120, 116, 117, 114, 101, 95, 115, 51, 116, 99, 0, 87, 69, 66, 71, 76, 95, 99, 111, 109, 112, 114, 101, 115, 115, 101, 100, 95, 116, 101, 120, 116, 117, 114, 101, 95, 115, 51, 116, 99, 95, 115, 114, 103, 98, 0, 87, 69, 66, 71, 76, 95, 100, 101, 98, 117, 103, 95, 114, 101, 110, 100, 101, 114, 101, 114, 95, 105, 110, 102, 111, 0, 87, 69, 66, 71, 76, 95, 100, 101, 98, 117, 103, 95, 115, 104, 97, 100, 101, 114, 115, 0, 87, 69, 66, 71, 76, 95, 100, 101, 112, 116, 104, 95, 116, 101, 120, 116, 117, 114, 101, 0, 87, 69, 66, 75, 73, 84, 95, 87, 69, 66, 71, 76, 95, 100, 101, 112, 116, 104, 95, 116, 101, 120, 116, 117, 114, 101, 0, 87, 69, 66, 71, 76, 95, 100, 114, 97, 119, 95, 98, 117, 102, 102, 101, 114, 115, 0, 87, 69, 66, 71, 76, 95, 108, 111, 115, 101, 95, 99, 111, 110, 116, 101, 120, 116, 0, 87, 69, 66, 75, 73, 84, 95, 87, 69, 66, 71, 76, 95, 108, 111, 115, 101, 95, 99, 111, 110, 116, 101, 120, 116, 0, 87, 69, 66, 71, 76, 95, 109, 117, 108, 116, 105, 95, 100, 114, 97, 119, 0}
						} else {
							payloadmap[supportedExtensionsIndex] = []int{33, 1, 87, 69, 66, 71, 76, 95, 109, 117, 108, 116, 105, 95, 100, 114, 97, 119, 0, 87, 69, 66, 75, 73, 84, 95, 87, 69, 66, 71, 76, 95, 108, 111, 115, 101, 95, 99, 111, 110, 116, 101, 120, 116, 0, 87, 69, 66, 71, 76, 95, 108, 111, 115, 101, 95, 99, 111, 110, 116, 101, 120, 116, 0, 87, 69, 66, 71, 76, 95, 100, 114, 97, 119, 95, 98, 117, 102, 102, 101, 114, 115, 0, 87, 69, 66, 75, 73, 84, 95, 87, 69, 66, 71, 76, 95, 100, 101, 112, 116, 104, 95, 116, 101, 120, 116, 117, 114, 101, 0, 87, 69, 66, 71, 76, 95, 100, 101, 112, 116, 104, 95, 116, 101, 120, 116, 117, 114, 101, 0, 87, 69, 66, 71, 76, 95, 100, 101, 98, 117, 103, 95, 115, 104, 97, 100, 101, 114, 115, 0, 87, 69, 66, 71, 76, 95, 100, 101, 98, 117, 103, 95, 114, 101, 110, 100, 101, 114, 101, 114, 95, 105, 110, 102, 111, 0, 87, 69, 66, 71, 76, 95, 99, 111, 109, 112, 114, 101, 115, 115, 101, 100, 95, 116, 101, 120, 116, 117, 114, 101, 95, 115, 51, 116, 99, 95, 115, 114, 103, 98, 0, 87, 69, 66, 75, 73, 84, 95, 87, 69, 66, 71, 76, 95, 99, 111, 109, 112, 114, 101, 115, 115, 101, 100, 95, 116, 101, 120, 116, 117, 114, 101, 95, 115, 51, 116, 99, 0, 87, 69, 66, 71, 76, 95, 99, 111, 109, 112, 114, 101, 115, 115, 101, 100, 95, 116, 101, 120, 116, 117, 114, 101, 95, 115, 51, 116, 99, 0, 87, 69, 66, 71, 76, 95, 99, 111, 108, 111, 114, 95, 98, 117, 102, 102, 101, 114, 95, 102, 108, 111, 97, 116, 0, 79, 69, 83, 95, 118, 101, 114, 116, 101, 120, 95, 97, 114, 114, 97, 121, 95, 111, 98, 106, 101, 99, 116, 0, 79, 69, 83, 95, 116, 101, 120, 116, 117, 114, 101, 95, 104, 97, 108, 102, 95, 102, 108, 111, 97, 116, 95, 108, 105, 110, 101, 97, 114, 0, 79, 69, 83, 95, 116, 101, 120, 116, 117, 114, 101, 95, 104, 97, 108, 102, 95, 102, 108, 111, 97, 116, 0, 79, 69, 83, 95, 116, 101, 120, 116, 117, 114, 101, 95, 102, 108, 111, 97, 116, 95, 108, 105, 110, 101, 97, 114, 0, 79, 69, 83, 95, 116, 101, 120, 116, 117, 114, 101, 95, 102, 108, 111, 97, 116, 0, 79, 69, 83, 95, 115, 116, 97, 110, 100, 97, 114, 100, 95, 100, 101, 114, 105, 118, 97, 116, 105, 118, 101, 115, 0, 79, 69, 83, 95, 102, 98, 111, 95, 114, 101, 110, 100, 101, 114, 95, 109, 105, 112, 109, 97, 112, 0, 79, 69, 83, 95, 101, 108, 101, 109, 101, 110, 116, 95, 105, 110, 100, 101, 120, 95, 117, 105, 110, 116, 0, 75, 72, 82, 95, 112, 97, 114, 97, 108, 108, 101, 108, 95, 115, 104, 97, 100, 101, 114, 95, 99, 111, 109, 112, 105, 108, 101, 0, 69, 88, 84, 95, 115, 82, 71, 66, 0, 87, 69, 66, 75, 73, 84, 95, 69, 88, 84, 95, 116, 101, 120, 116, 117, 114, 101, 95, 102, 105, 108, 116, 101, 114, 95, 97, 110, 105, 115, 111, 116, 114, 111, 112, 105, 99, 0, 69, 88, 84, 95, 116, 101, 120, 116, 117, 114, 101, 95, 102, 105, 108, 116, 101, 114, 95, 97, 110, 105, 115, 111, 116, 114, 111, 112, 105, 99, 0, 69, 88, 84, 95, 116, 101, 120, 116, 117, 114, 101, 95, 99, 111, 109, 112, 114, 101, 115, 115, 105, 111, 110, 95, 114, 103, 116, 99, 0, 69, 88, 84, 95, 116, 101, 120, 116, 117, 114, 101, 95, 99, 111, 109, 112, 114, 101, 115, 115, 105, 111, 110, 95, 98, 112, 116, 99, 0, 69, 88, 84, 95, 115, 104, 97, 100, 101, 114, 95, 116, 101, 120, 116, 117, 114, 101, 95, 108, 111, 100, 0, 69, 88, 84, 95, 102, 114, 97, 103, 95, 100, 101, 112, 116, 104, 0, 69, 88, 84, 95, 102, 108, 111, 97, 116, 95, 98, 108, 101, 110, 100, 0, 69, 88, 84, 95, 100, 105, 115, 106, 111, 105, 110, 116, 95, 116, 105, 109, 101, 114, 95, 113, 117, 101, 114, 121, 0, 69, 88, 84, 95, 99, 111, 108, 111, 114, 95, 98, 117, 102, 102, 101, 114, 95, 104, 97, 108, 102, 95, 102, 108, 111, 97, 116, 0, 69, 88, 84, 95, 98, 108, 101, 110, 100, 95, 109, 105, 110, 109, 97, 120, 0, 65, 78, 71, 76, 69, 95, 105, 110, 115, 116, 97, 110, 99, 101, 100, 95, 97, 114, 114, 97, 121, 115, 0}
						}
					}
				}

				a := []int{antialiasIndex, shaderPrecisionsIndex, maxTextureSizeIndex, dimensionsIndex, maxTextureImageUnitsIndex, paramrenderindex, paramvendorindex, debugrenderindex, debugvendorindex, maxVertexUniformVectorsIndex, maxFragmentUniformVectorsIndex, maxVertexTextureImageUnitsIndex, contextPropertiesIndex, supportedExtensionsIndex, blueBitsIndex, depthBitsIndex, greenBitsIndex, redBitsIndex, stencilBitsIndex, maxAnisotropyIndex, maxRenderbufferSizeIndex, maxVertexAttribsIndex, versionIndex, shadingLanguageVersionIndex, maxVaryingVectorsIndex, maxCubeMapTextureSizeIndex, alphaBitsIndex, maxCombinedTextureImageUnitsIndex}
				sort.Ints(a)

				for j, _ := range vendorholders{
					apparr := []int{}
					for i, val := range a {
						if i == debugvendorindex{
							apparr = append(apparr, vendorholders[j]...)
						}else
						if i == debugrenderindex{
							apparr = append(apparr, renderrholders[j]...)
						}else
						{
							apparr = append(apparr, payloadmap[val]...)
						}
					}
					g.P17VAL = append(g.P17VAL, apparr)
				}
				g.PayloadOrder = append(g.PayloadOrder, 17)
				for k, _ := range mapunion {
					g.KeyCount[k]++
					g.KeyHolderRaw[17] = append(g.KeyHolderRaw[17], k)
				}
			} else
			if strings.Contains(outstring, "webdriver") && !strings.Contains(outstring, "\"Andale Mono\"") { //payload19
				P40DOCUMENTval := P40DOCUMENTMATCH.FindStringSubmatch(outstring)[1]
				P40DOCUMENTBODYval := P40DOCUMENTBODYMATCH.FindStringSubmatch(outstring)[1]
				P40GLOBALval := P40GLOBALMATCH.FindStringSubmatch(outstring)[1]
				P40NAVIGATORval := P40NAVIGATORMATCH.FindStringSubmatch(outstring)[1]
				P40CRYPTOval := P40CRYPTOMATCH.FindStringSubmatch(outstring)[1]
				P40EXTERNALval := P40EXTERNALMATCH.FindStringSubmatch(outstring)[1]
				outstring = strings.ReplaceAll(outstring, fmt.Sprintf("heapin[%s]", P40DOCUMENTval), "DOCUMENTOBJECT")
				outstring = strings.ReplaceAll(outstring, fmt.Sprintf("heapin[%s]", P40DOCUMENTBODYval), "DOCUMENTBODYOBJECT")
				outstring = strings.ReplaceAll(outstring, fmt.Sprintf("heapin[%s]", P40GLOBALval), "GLOBALOBJECT")
				outstring = strings.ReplaceAll(outstring, fmt.Sprintf("heapin[%s]", P40NAVIGATORval), "NAVIGATOROBJECT")
				outstring = strings.ReplaceAll(outstring, fmt.Sprintf("heapin[%s]", P40CRYPTOval), "CRYPTOOBJECT")
				outstring = strings.ReplaceAll(outstring, fmt.Sprintf("heapin[%s]", P40EXTERNALval), "EXTERNALOBJECT")
				payloadmap := make(map[int]int)
				check0index := strings.Index(outstring, "GLOBALOBJECT[\"requestAnimationFrame\"]")
				check1index := strings.Index(outstring, "DOCUMENTBODYOBJECT[innerText]")
				check2index := strings.Index(outstring, "NAVIGATOROBJECT[\"credentials\"]")
				check3index := strings.Index(outstring, "GLOBALOBJECT[\"XDomainRequest\"]")
				check4index := strings.Index(outstring, "DOCUMENTOBJECT[\"__fxdriver_unwrapped\"]")
				check5index := strings.Index(outstring, "GLOBALOBJECT[\"_Selenium_IDE_Recorder\"]")
				check6index := strings.Index(outstring, "GLOBALOBJECT[\"external\"]")
				check7index := strings.Index(outstring, "DOCUMENTOBJECT[\"_Selenium_IDE_Recorder\"]")
				check8index := strings.Index(outstring, "GLOBALOBJECT[\"sessionStorage\"]")
				check9index := strings.Index(outstring, "GLOBALOBJECT[\"MutationObserver\"]")
				check10index := strings.Index(outstring, "GLOBALOBJECT[\"webkitResolveLocalFileSystemURL\"]")
				check11index := strings.Index(outstring, "DOCUMENTOBJECT[all]")
				check12index := strings.Index(outstring, "GLOBALOBJECT[\"webkitRequestAnimationFrame\"]")
				check13index := strings.Index(outstring, "GLOBALOBJECT[\"localStorage\"]")
				check14index := strings.Index(outstring, "NAVIGATOROBJECT[\"requestMediaKeySystemAccess\"]")
				check15index := strings.Index(outstring, "DOCUMENTBODYOBJECT[\"contextMenu\"]")
				check16index := strings.Index(outstring, "GLOBALOBJECT[\"globalStorage\"]")
				check17index := strings.Index(outstring, "DOCUMENTBODYOBJECT[\"webkitRequestFullScreen\"]")
				check18index := strings.Index(outstring, "GLOBALOBJECT[\"postMessage\"]")
				check19index := strings.Index(outstring, "NAVIGATOROBJECT[\"webdriver\"]")
				check20index := strings.Index(outstring, "GLOBALOBJECT[\"BluetoothUUID\"]")
				check21index := strings.Index(outstring, "GLOBALOBJECT[\"_phantom\"]")
				check22index := strings.Index(outstring, "DOCUMENTOBJECT[\"__webdriver_script_fn\"]")
				check23index := strings.Index(outstring, "GLOBALOBJECT[fireEvent]")
				check24index := strings.Index(outstring, "DOCUMENTBODYOBJECT[\"mozRequestFullScreen\"]")
				check25index := strings.Index(outstring, "DOCUMENTOBJECT[\"images\"]")
				check26index := strings.Index(outstring, "GLOBALOBJECT[File]")
				check27index := strings.Index(outstring, "GLOBALOBJECT[\"TouchEvent\"]")
				check28index := strings.Index(outstring, "NAVIGATOROBJECT[\"bluetooth\"]")
				check29index := strings.Index(outstring, "NAVIGATOROBJECT[\"storage\"]")
				check30index := strings.Index(outstring, "GLOBALOBJECT[\"PushManager\"]")
				check31index := strings.Index(outstring, "GLOBALOBJECT[event]")
				check32index := strings.Index(outstring, "DOCUMENTOBJECT[characterSet]")
				check33index := strings.Index(outstring, "DOCUMENTOBJECT[\"layers\"]")
				check34index := strings.Index(outstring, "DOCUMENTOBJECT[\"$cdc_asdjflasutopfhvcZLmcfl_\"]")
				check35index := strings.Index(outstring, "GLOBALOBJECT[\"callPhantom\"]")
				check36index := strings.Index(outstring, "GLOBALOBJECT[\"ActiveXObject\"]")
				check37index := strings.Index(outstring, "GLOBALOBJECT[\"SharedWorker\"]")
				check38index := strings.Index(outstring, "GLOBALOBJECT[\"sidebar\"]")
				check39index := strings.Index(outstring, "GLOBALOBJECT[attachEvent]")
				check40index := strings.Index(outstring, "CRYPTOOBJECT[\"subtle\"]")
				check41index := strings.Index(outstring, "DOCUMENTOBJECT[\"compatMode\"]")
				check42index := strings.Index(outstring, "GLOBALOBJECT[\"ApplePaySession\"]")
				check43index := strings.Index(outstring, "DOCUMENTOBJECT[charset]")
				check44index := strings.Index(outstring, "GLOBALOBJECT[\"mozRequestAnimationFrame\"]")
				check45index := strings.Index(outstring, "GLOBALOBJECT[detachEvent]")
				check46index := strings.Index(outstring, "GLOBALOBJECT[\"phantom\"]")
				check47index := strings.Index(outstring, "DOCUMENTOBJECT[documentMode]")
				check48index := strings.Index(outstring, "GLOBALOBJECT[\"Notification\"]")
				check49index := strings.Index(outstring, "NAVIGATOROBJECT[\"vibrate\"]")
				check50index := strings.Index(outstring, "GLOBALOBJECT[\"webkitRTCPeerConnection\"]")
				check51index := strings.Index(outstring, "DOCUMENTBODYOBJECT[\"requestFullScreen\"]")
				check52index := strings.Index(outstring, "GLOBALOBJECT[\"mozRTCPeerConnection\"]")
				check53index := strings.Index(outstring, "GLOBALOBJECT[\"netscape\"]")
				check54index := strings.Index(outstring, "GLOBALOBJECT[\"registerProtocolHandler\"]")
				check55index := strings.Index(outstring, "EXTERNALOBJECT[\"Sequentum\"]")
				check56index := strings.Index(outstring, "GLOBALOBJECT[\"__fxdriver_unwrapped\"]")
				check57index := strings.Index(outstring, "GLOBALOBJECT[\"createPopup\"]")
				check58index := strings.Index(outstring, "GLOBALOBJECT[frameElement]")
				payloadmap[check0index] = 1
				payloadmap[check1index] = 0
				payloadmap[check2index] = 1
				payloadmap[check3index] = 0
				payloadmap[check4index] = 0
				payloadmap[check5index] = 0
				payloadmap[check6index] = 1
				payloadmap[check7index] = 0
				payloadmap[check8index] = 1
				payloadmap[check9index] = 1
				payloadmap[check10index] = 1
				payloadmap[check11index] = 1
				payloadmap[check12index] = 1
				payloadmap[check13index] = 1
				payloadmap[check14index] = 1
				payloadmap[check15index] = 0
				payloadmap[check16index] = 0
				payloadmap[check17index] = 0
				payloadmap[check18index] = 1
				payloadmap[check19index] = 1
				payloadmap[check20index] = 1
				payloadmap[check21index] = 0
				payloadmap[check22index] = 0
				payloadmap[check23index] = 0
				payloadmap[check24index] = 0
				payloadmap[check25index] = 1
				payloadmap[check26index] = 1
				payloadmap[check27index] = 1
				payloadmap[check28index] = 1
				payloadmap[check29index] = 1
				payloadmap[check30index] = 1
				payloadmap[check31index] = 1
				payloadmap[check32index] = 1
				payloadmap[check33index] = 0
				payloadmap[check34index] = 0
				payloadmap[check35index] = 0
				payloadmap[check36index] = 0
				payloadmap[check37index] = 1
				payloadmap[check38index] = 0
				payloadmap[check39index] = 0
				payloadmap[check40index] = 1
				payloadmap[check41index] = 1
				payloadmap[check42index] = 0
				payloadmap[check43index] = 1
				payloadmap[check44index] = 0
				payloadmap[check45index] = 0
				payloadmap[check46index] = 0
				payloadmap[check47index] = 0
				payloadmap[check48index] = 1
				payloadmap[check49index] = 1
				payloadmap[check50index] = 1
				payloadmap[check51index] = 0
				payloadmap[check52index] = 0
				payloadmap[check53index] = 0
				payloadmap[check54index] = 0
				payloadmap[check55index] = 0
				payloadmap[check56index] = 0
				payloadmap[check57index] = 0
				payloadmap[check58index] = 1
				a := []int{check0index, check1index, check2index, check3index, check4index, check5index, check6index, check7index, check8index, check9index, check10index, check11index, check12index, check13index, check14index, check15index, check16index, check17index, check18index, check19index, check20index, check21index, check22index, check23index, check24index, check25index, check26index, check27index, check28index, check29index, check30index, check31index, check32index, check33index, check34index, check35index, check36index, check37index, check38index, check39index, check40index, check41index, check42index, check43index, check44index, check45index, check46index, check47index, check48index, check49index, check50index, check51index, check52index, check53index, check54index, check55index, check56index, check57index, check58index}
				sort.Ints(a)
				for _, val := range a {
					g.P18VAL = append(g.P18VAL, payloadmap[val])
				}
				g.PayloadOrder = append(g.PayloadOrder, 18)
				for k, _ := range mapunion {
					g.KeyCount[k]++
					g.KeyHolderRaw[18] = append(g.KeyHolderRaw[18], k)
				}
			} else
			if strings.Contains(outstring, "[description]") { //payload20
				payloadmap := make(map[int][]int)
				blocks := strings.Split(outstring, "[Math][\"random\"]()")
				lastblock := blocks[len(blocks)-1]
				lastheapsetindex := strings.LastIndex(lastblock, P19HEAPSETMATCH.FindStringSubmatch(lastblock)[0])
				lastlengthindex := strings.LastIndex(lastblock, "[length]")
				firstlengthindex := strings.Index(lastblock, "[length]")
				lastshiftindex := strings.LastIndex(lastblock, "<<1")
				if strings.Index(lastblock, "if(0;)") > 0 {
					payloadmap[lastheapsetindex] = []int{50, 183, 232, 136, 16, 61, 243, 219, 211, 13, 115, 188, 197, 152, 7}
				} else {
					payloadmap[lastheapsetindex] = []int{115, 188, 197, 152, 7, 61, 243, 219, 211, 13, 50, 183, 232, 136, 16}
				}
				payloadmap[lastlengthindex] = []int{3}
				payloadmap[firstlengthindex] = []int{3}
				payloadmap[lastshiftindex] = []int{0}

				a := []int{firstlengthindex, lastheapsetindex, lastlengthindex, lastshiftindex}
				sort.Ints(a)
				for _, val := range a {
					g.P19VAL = append(g.P19VAL, payloadmap[val]...)
				}

				g.PayloadOrder = append(g.PayloadOrder, 19)
				for k, _ := range mapunion {
					g.KeyCount[k]++
					g.KeyHolderRaw[19] = append(g.KeyHolderRaw[19], k)
				}
			} else
			if strings.Count(outstring, "[\"next\"]") == 4 && strings.Contains(outstring, "[push]") && !strings.Contains(outstring, "stack3.push") && !strings.Contains(outstring, "console") { //payload21
				g.PayloadOrder = append(g.PayloadOrder, 20)
				float1 := float64(int(time.Now().UTC().UnixNano()/1e6 + int64(rand.Uint32())))
				float2 := float64(int(time.Now().UTC().UnixNano()/1e6 + int64(rand.Uint32())))
				float3 := float64(int(time.Now().UTC().UnixNano()/1e6 + int64(rand.Uint32())))
				g.P20VAL = append(g.P20VAL, 3)
				for _, val := range generation.BaseNumEnc(float1) {
					g.P20VAL = append(g.P20VAL, val)
				}
				for _, val := range generation.BaseNumEnc(float2) {
					g.P20VAL = append(g.P20VAL, val)
				}
				for _, val := range generation.BaseNumEnc(float3) {
					g.P20VAL = append(g.P20VAL, val)
				}
				for k, _ := range mapunion {
					g.KeyCount[k]++
					g.KeyHolderRaw[20] = append(g.KeyHolderRaw[20], k)
				}
			} else
			if strings.Contains(outstring, "numOrientationEvents") { //payload22
				payloadmap := make(map[int][]int)
				stdDevAlphaindex := strings.Index(outstring, "stdDevAlpha")
				payloadmap[stdDevAlphaindex] = []int{0}
				stdDevIntervalindex := strings.Index(outstring, "stdDevInterval")
				payloadmap[stdDevIntervalindex] = []int{0}
				avgAlphaindex := strings.Index(outstring, "avgAlpha")
				payloadmap[avgAlphaindex] = []int{128, 0, 0, 0, 0, 0, 0, 0, 0}
				stdDevGammaindex := strings.Index(outstring, "stdDevGamma")
				payloadmap[stdDevGammaindex] = []int{0}
				avgBetaindex := strings.Index(outstring, "avgBeta")
				payloadmap[avgBetaindex] = []int{128, 0, 0, 0, 0, 0, 0, 0, 0}
				numOrientationEventsindex := strings.Index(outstring, "numOrientationEvents")
				payloadmap[numOrientationEventsindex] = []int{1}
				stdDevBetaindex := strings.Index(outstring, "stdDevBeta")
				payloadmap[stdDevBetaindex] = []int{0}
				avgIntervalindex := strings.Index(outstring, "avgInterval")
				payloadmap[avgIntervalindex] = []int{0}
				avgGammaindex := strings.Index(outstring, "avgGamma")
				payloadmap[avgGammaindex] = []int{128, 0, 0, 0, 0, 0, 0, 0, 0}
				orderArr := []int{stdDevAlphaindex, stdDevIntervalindex, avgAlphaindex, stdDevGammaindex, avgBetaindex, numOrientationEventsindex, stdDevBetaindex, avgIntervalindex, avgGammaindex}
				sort.Ints(orderArr)
				for _, val := range orderArr {
					g.P21VAL = append(g.P21VAL, payloadmap[val]...)
				}
				g.PayloadOrder = append(g.PayloadOrder, 21)
				for k, _ := range mapunion {
					g.KeyCount[k]++
					g.KeyHolderRaw[21] = append(g.KeyHolderRaw[21], k)
				}
			} else
			if strings.Contains(outstring, "[map]") && strings.Contains(outstring, "[join]") { //payload23
				g.PayloadOrder = append(g.PayloadOrder, 22)
				for k, _ := range mapunion {
					g.KeyCount[k]++
					g.KeyHolderRaw[22] = append(g.KeyHolderRaw[22], k)
				}
			} else
			if strings.Contains(outstring, "[JSON][stringify]") { //payload24
				outstring = elseSlicer(outstring)
				xormatch := P23XORMATCH.FindStringSubmatch(outstring)
				appval, _ := strconv.Atoi(xormatch[1])
				payloadmap := make(map[int]int)
				xorindex := strings.Index(outstring, xormatch[0])
				globalindex := strings.Index(outstring, "[global]")
				payloadmap[xorindex] = appval
				payloadmap[globalindex] = 0
				a := []int{xorindex, globalindex}
				sort.Ints(a)
				for _, val := range a {
					g.P23VAL = append(g.P23VAL, payloadmap[val])
				}
				g.PayloadOrder = append(g.PayloadOrder, 23)
				for k, _ := range mapunion {
					g.KeyCount[k]++
					g.KeyHolderRaw[23] = append(g.KeyHolderRaw[23], k)
				}
			} else
			if strings.Contains(outstring, "[\"1.0.0\"]") { //payload25
				payloadmap := make(map[int]int)
				check1index := strings.Index(outstring, "[\"1.0.0\"]")
				check2index := strings.Index(outstring, "[\"6.1.0\"]")
				check3index := strings.Index(outstring, "[\"6.2.0\"]")
				check4index := strings.Index(outstring, "[\"6.3.0\"]")
				check5index := strings.Index(outstring, "[\"7.0.0\"]")
				check6index := strings.Index(outstring, "[\"8.0.0\"]")
				check7index := strings.Index(outstring, "[\"9.0.0\"]")
				check8index := strings.Index(outstring, "[\"10.0.0\"]")
				check9index := strings.Index(outstring, "[\"11.0.0\"]")
				check10index := strings.Index(outstring, "[\"12.0.0\"]")
				check11index := strings.Index(outstring, "[\"12.1.0\"]")
				check12index := strings.Index(outstring, "[\"13.0.0\"]")
				payloadmap[check1index] = 1
				payloadmap[check2index] = 1
				payloadmap[check3index] = 1
				payloadmap[check4index] = 0
				payloadmap[check5index] = 1
				payloadmap[check6index] = 1
				payloadmap[check7index] = 1
				payloadmap[check8index] = 1
				payloadmap[check9index] = 1
				payloadmap[check10index] = 1
				payloadmap[check11index] = 1
				payloadmap[check12index] = 1
				a := []int{check1index, check2index, check3index, check4index, check5index, check6index, check7index, check8index, check9index, check10index, check11index, check12index}
				sort.Ints(a)
				val1 := 0
				for i := 0; i < 7; i++ {
					val1 |= (payloadmap[a[i]] << i)
				}
				val2 := 0
				for i := 0; i < 5; i++ {
					val2 |= (payloadmap[a[i+7]] << i)
				}
				g.P24VAL = []int{val1, val2}
				g.PayloadOrder = append(g.PayloadOrder, 24)
				for k, _ := range mapunion {
					g.KeyCount[k]++
					g.KeyHolderRaw[24] = append(g.KeyHolderRaw[24], k)
				}
			} else
			if len(sixteenKeyRegex.FindStringSubmatch(outstring)) > 0 { //payload26
				g.P25VAL = []int{}
				g.PayloadOrder = append(g.PayloadOrder, 25)
				for k, _ := range mapunion {
					g.KeyCount[k]++
					g.KeyHolderRaw[25] = append(g.KeyHolderRaw[25], k)
				}
				hashkey := sixteenKeyRegex.FindStringSubmatch(outstring)[1]

				hashindex := strings.Index(outstring, hashkey)
				holderindexes := []int{hashindex}

				for _, match := range P25SHIFTMATCH.FindAllIndex([]byte(outstring), -1) {
					holderindexes = append(holderindexes, match[0])
				}
				sort.Ints(holderindexes)
				for _, val := range holderindexes {
					switch val {
					case hashindex:
						for _, sv := range hashkey {
							g.P25VAL = append(g.P25VAL, int(sv))
						}
						g.P25VAL = append(g.P25VAL, 0)
					default:
						for _, sv := range generation.BaseNumEnc((20 + rand.Intn(236)) << 8) {
							g.P25VAL = append(g.P25VAL, sv)
						}
					}
				}

			} else
			if len(windowdatematch.FindStringSubmatch(outstring)) > 0 && !strings.Contains(outstring, "NINT16") { //payload27
				g.PayloadOrder = append(g.PayloadOrder, 26)
				for k, _ := range mapunion {
					g.KeyCount[k]++
					g.KeyHolderRaw[26] = append(g.KeyHolderRaw[26], k)
				}
			} else
			if strings.Contains(outstring, "[error]") && strings.Contains(outstring, "[filter]") { //payload28
				g.PayloadOrder = append(g.PayloadOrder, 27)
				for k, _ := range mapunion {
					g.KeyCount[k]++
					g.KeyHolderRaw[27] = append(g.KeyHolderRaw[27], k)
				}
			} else
			if strings.Contains(outstring, "[\"touches\"]") { //payload29
				g.PayloadOrder = append(g.PayloadOrder, 28)
				for k, _ := range mapunion {
					g.KeyCount[k]++
					g.KeyHolderRaw[28] = append(g.KeyHolderRaw[28], k)
				}
			} else
			if len(dateheapregex.FindStringSubmatch(outstring)) > 0 { //payload30
				g.PayloadOrder = append(g.PayloadOrder, 29)
				for k, _ := range mapunion {
					g.KeyCount[k]++
					g.KeyHolderRaw[29] = append(g.KeyHolderRaw[29], k)
				}
			} else
			if strings.Contains(outstring, "[\"getContext\"]") && strings.Contains(outstring, "[\"emoji\"]") { //payload31
				if strings.Index(outstring, "if(0;)") > 0 {
					g.P30VAL = []int{8, 63, 7, 63, 6, 63, 5, 37, 4, 63, 7, 48, 1, 63, 7, 32, 4}
				} else {
					g.P30VAL = []int{8, 32, 4, 63, 7, 48, 1, 63, 7, 37, 4, 63, 5, 63, 6, 63, 7}
				}
				g.PayloadOrder = append(g.PayloadOrder, 30)
				for k, _ := range mapunion {
					g.KeyCount[k]++
					g.KeyHolderRaw[30] = append(g.KeyHolderRaw[30], k)
				}
			} else
			if strings.Contains(outstring, "toSource") { //payload32
				var appval int
				var appstringindex int
				count := 0
				for _, match := range P31VALMATCH.FindAllStringSubmatch(outstring, -1) {
					if match[1] != "0" && match[1] != "3" {
						if count == 0 {
							appval, _ = strconv.Atoi(match[1])
							appstringindex = strings.Index(outstring, match[0])
							count++
						}
					}
				}
				payloadmap := make(map[int][]int)
				//if (!heapin[0]; 9091
				check1index := strings.LastIndex(outstring, "toString")
				check2index := strings.LastIndex(outstring, "toSource")
				payloadmap[check1index] = []int{35, 200, 216, 148, 1}
				payloadmap[check2index] = []int{0}
				payloadmap[appstringindex] = []int{appval}
				a := []int{check1index, check2index, appstringindex}
				sort.Ints(a)
				for _, val := range a {
					g.P31VAL = append(g.P31VAL, payloadmap[val]...)
				}
				g.PayloadOrder = append(g.PayloadOrder, 31)
				for k, _ := range mapunion {
					g.KeyCount[k]++
					g.KeyHolderRaw[31] = append(g.KeyHolderRaw[31], k)
				}
			} else
			if strings.Contains(outstring, "STORAGE_OK") && !strings.Contains(outstring, "\"Andale Mono\"") { //payload33
				payloadmap := make(map[int][]int)
				P32ITEMindex := strings.Index(outstring, fmt.Sprintf("](heapin[%s],", P32ITEMMATCH.FindStringSubmatch(outstring)[1]))
				P32CAPACITYindex := strings.Index(outstring, fmt.Sprintf("](heapin[%s],", P32CAPACITYMATCH.FindStringSubmatch(outstring)[1]))
				payloadmap[P32ITEMindex] = []int{}
				for _, val := range g.Seedstring {
					payloadmap[P32ITEMindex] = append(payloadmap[P32ITEMindex], int(val))
				}
				payloadmap[P32ITEMindex] = append(payloadmap[P32ITEMindex], 0)
				payloadmap[P32CAPACITYindex] = []int{3}
				a := []int{P32ITEMindex, P32CAPACITYindex}
				sort.Ints(a)
				for _, val := range a {
					g.P32VAL = append(g.P32VAL, payloadmap[val]...)
				}

				g.PayloadOrder = append(g.PayloadOrder, 32)
				for k, _ := range mapunion {
					g.KeyCount[k]++
					g.KeyHolderRaw[32] = append(g.KeyHolderRaw[32], k)
				}
			} else
			if len(numbernewheapregex.FindStringSubmatch(outstring)) > 0 { //payload34
				g.PayloadOrder = append(g.PayloadOrder, 33)
				for k, _ := range mapunion {
					g.KeyCount[k]++
					g.KeyHolderRaw[33] = append(g.KeyHolderRaw[33], k)
				}
			} else
			if strings.Contains(outstring, "[\"candidate\"]") { //payload35
				g.PayloadOrder = append(g.PayloadOrder, 34)
				for k, _ := range mapunion {
					g.KeyCount[k]++
					g.KeyHolderRaw[34] = append(g.KeyHolderRaw[34], k)
				}
			} else
			if strings.Contains(outstring, "[\"protocol\"]") && strings.Contains(outstring, "[\"hostname\"]") && strings.Contains(outstring, "[\"port\"]") && !strings.Contains(outstring, "\"Andale Mono\"") { //payload36
				payloadmap := make(map[int][]int)
				check1index := strings.Index(outstring, "[\"protocol\"]")
				check2index := strings.Index(outstring, "[\"hostname\"]")
				check3index := strings.Index(outstring, "[\"port\"]")
				payloadmap[check1index] = []int{104, 116, 116, 112, 115, 58, 0}
				payloadmap[check2index] = []int{108, 111, 103, 105, 110, 46, 116, 97, 114, 103, 101, 116, 46, 99, 111, 109, 0}
				payloadmap[check3index] = []int{0}
				a := []int{check1index, check2index, check3index}
				sort.Ints(a)
				for _, val := range a {
					g.P35VAL = append(g.P35VAL, payloadmap[val]...)
				}
				g.PayloadOrder = append(g.PayloadOrder, 35)
				for k, _ := range mapunion {
					g.KeyCount[k]++
					g.KeyHolderRaw[35] = append(g.KeyHolderRaw[35], k)
				}
			} else
			if strings.Contains(outstring, "[\"toDataURL\"]") { //payload37
				P36URLVAL := P36URLMATCH.FindStringSubmatch(outstring)[1]
				P36DATAVAL := P36DATAMATCH.FindStringSubmatch(outstring)[1]
				P36URLindex := strings.Index(outstring, fmt.Sprintf("(heapin[%s]", P36URLVAL))
				P36DATAindex := strings.Index(outstring, fmt.Sprintf("(heapin[%s]", P36DATAVAL))
				a := []int{P36URLindex, P36DATAindex}
				payloadmap := map[int][]int{
					P36URLindex:  {100, 97, 116, 97, 58, 105, 109, 97, 103, 101, 47, 112, 110, 103, 59, 98, 97, 115, 101, 54, 52, 44, 105, 86, 66, 79, 82, 119, 48, 75, 71, 103, 111, 65, 65, 65, 65, 78, 83, 85, 104, 69, 85, 103, 65, 65, 65, 83, 119, 65, 0},
					P36DATAindex: {56, 217, 169, 187, 15},
				}
				sort.Ints(a)
				for _, val := range a {
					g.P36VAL = append(g.P36VAL, payloadmap[val]...)
				}
				g.PayloadOrder = append(g.PayloadOrder, 36)
				for k, _ := range mapunion {
					g.KeyCount[k]++
					g.KeyHolderRaw[36] = append(g.KeyHolderRaw[36], k)
				}
			} else
			if strings.Contains(starttag, "[Function][prototype][toString]") { //payload38
				g.PayloadOrder = append(g.PayloadOrder, 37)
				for k, _ := range mapunion {
					g.KeyCount[k]++
					g.KeyHolderRaw[37] = append(g.KeyHolderRaw[37], k)
				}
			} else
			if strings.Contains(outstring, "[Object][\"keys\"]") && strings.Contains(outstring, "[global][location]") { //payload39
				g.PayloadOrder = append(g.PayloadOrder, 38)
				for k, _ := range mapunion {
					g.KeyCount[k]++
					g.KeyHolderRaw[38] = append(g.KeyHolderRaw[38], k)
				}
			} else
			if strings.Contains(outstring, "[Object][getOwnPropertyNames]") && !strings.Contains(outstring, "[join]") { //payload40
				g.PayloadOrder = append(g.PayloadOrder, 39)
				for k, _ := range mapunion {
					g.KeyCount[k]++
					g.KeyHolderRaw[39] = append(g.KeyHolderRaw[39], k)
				}
			} else
			if strings.Contains(outstring, "[Object][getOwnPropertyNames]") && strings.Contains(outstring, "[join]") { //payload41
				payloadmap := make(map[int][]int)
				check1index := strings.Index(outstring, "[join]")
				check2index := strings.Index(outstring, "[length]")
				payloadmap[check1index] = []int{40, 226, 255, 144, 22}
				payloadmap[check2index] = []int{43, 1, 83, 81, 82, 84, 50, 0, 83, 81, 82, 84, 49, 95, 50, 0, 80, 73, 0, 76, 79, 71, 50, 69, 0, 76, 79, 71, 49, 48, 69, 0, 76, 78, 50, 0, 76, 78, 49, 48, 0, 69, 0, 116, 114, 117, 110, 99, 0, 116, 97, 110, 104, 0, 116, 97, 110, 0, 115, 113, 114, 116, 0, 115, 105, 110, 104, 0, 115, 105, 110, 0, 115, 105, 103, 110, 0, 114, 111, 117, 110, 100, 0, 114, 97, 110, 100, 111, 109, 0, 112, 111, 119, 0, 109, 105, 110, 0, 109, 97, 120, 0, 108, 111, 103, 49, 48, 0, 108, 111, 103, 50, 0, 108, 111, 103, 49, 112, 0, 108, 111, 103, 0, 105, 109, 117, 108, 0, 104, 121, 112, 111, 116, 0, 102, 114, 111, 117, 110, 100, 0, 102, 108, 111, 111, 114, 0, 101, 120, 112, 0, 99, 111, 115, 104, 0, 99, 111, 115, 0, 99, 108, 122, 51, 50, 0, 101, 120, 112, 109, 49, 0, 99, 98, 114, 116, 0, 99, 101, 105, 108, 0, 97, 116, 97, 110, 50, 0, 97, 116, 97, 110, 104, 0, 97, 116, 97, 110, 0, 97, 115, 105, 110, 104, 0, 97, 115, 105, 110, 0, 97, 99, 111, 115, 104, 0, 97, 99, 111, 115, 0, 97, 98, 115, 0}
				a := []int{check1index, check2index}
				sort.Ints(a)
				for _, val := range a {
					g.P40VAL = append(g.P40VAL, payloadmap[val]...)
				}

				g.PayloadOrder = append(g.PayloadOrder, 40)
				for k, _ := range mapunion {
					g.KeyCount[k]++
					g.KeyHolderRaw[40] = append(g.KeyHolderRaw[40], k)
				}
			} else
			if strings.Contains(outstring, "[push]") && strings.Contains(outstring, "[join]") { //payload42
				g.PayloadOrder = append(g.PayloadOrder, 41)
				for k, _ := range mapunion {
					g.KeyCount[k]++
					g.KeyHolderRaw[41] = append(g.KeyHolderRaw[41], k)
				}
			} else
			if strings.Contains(outstring, "[\"matchMedia\"]") && !strings.Contains(outstring, "\"Andale Mono\"") { //payload43
				P42HEAPVAL, _ := strconv.Atoi(P42HEAPSETMATCH.FindAllStringSubmatch(outstring, -1)[2][1])
				outstring = elseSlicer(outstring)
				payloadmap := make(map[int][]int)
				P42IND0 := P42IND0MATCH.FindStringSubmatch(outstring)
				P42IND0index := strings.Index(outstring, P42IND0[0])
				P42IND1 := P42IND1MATCH.FindStringSubmatch(outstring)
				P42IND1index := strings.Index(outstring, P42IND1[0])
				P42HEAPSETindex := strings.Index(outstring, P42HEAPSETMATCH.FindAllStringSubmatch(outstring, -1)[1][0])
				P42DPIindex := strings.Index(outstring, fmt.Sprintf("](heapin[%s],", P42DPIMATCH.FindStringSubmatch(outstring)[1]))
				overflow_block_index := strings.Index(outstring, "overflow-block")
				pointerindex := strings.Index(outstring, "\"pointer")
				max_monochromeindex := strings.Index(outstring, "max-monochrome")
				inverted_colorsindex := strings.Index(outstring, "inverted-colors")
				any_hoverindex := strings.Index(outstring, "any-hover")
				updateindex := strings.Index(outstring, "update")
				max_heightindex := strings.Index(outstring, "max-height")
				scriptingindex := strings.Index(outstring, "scripting")
				prefers_contrastindex := strings.Index(outstring, "prefers-contrast")
				prefers_color_schemeindex := strings.Index(outstring, "prefers-color-scheme")
				orientationindex := strings.Index(outstring, "orientation")
				hoverindex := strings.Index(outstring, ",\"hover")
				forced_colorsindex := strings.Index(outstring, "forced-colors")
				overflow_inlineindex := strings.Index(outstring, "overflow-inline")
				max_widthindex := strings.Index(outstring, "max-width")
				max_colorindex := strings.Index(outstring, "max-color\"")
				display_modeindex := strings.Index(outstring, "display-mode")
				color_gamutindex := strings.Index(outstring, "color-gamut")
				prefers_reduced_motionindex := strings.Index(outstring, "prefers-reduced-motion")
				gridindex := strings.Index(outstring, "grid")
				any_pointerindex := strings.Index(outstring, "any-pointer")
				scanindex := strings.Index(outstring, "scan")
				light_levelindex := strings.Index(outstring, "light-level")
				prefers_reduced_transparencyindex := strings.Index(outstring, "prefers-reduced-transparency")
				max_color_indexindex := strings.Index(outstring, "max-color-index")

				a := []int{P42DPIindex, P42IND0index, P42IND1index, P42HEAPSETindex, overflow_block_index, pointerindex, max_monochromeindex, inverted_colorsindex, any_hoverindex, updateindex, max_heightindex, scriptingindex, prefers_contrastindex, prefers_color_schemeindex, orientationindex, hoverindex, forced_colorsindex, overflow_inlineindex, max_widthindex, max_colorindex, display_modeindex, color_gamutindex, prefers_reduced_motionindex, gridindex, any_pointerindex, scanindex, light_levelindex, prefers_reduced_transparencyindex, max_color_indexindex}
				sort.Ints(a)

				payloadmap[P42HEAPSETindex] = []int{P42HEAPVAL}
				payloadmap[P42DPIindex] = []int{54, 3}
				payloadmap[P42IND0index] = []int{55, 170, 230, 158, 14}
				payloadmap[P42IND1index] = []int{128, 0, 0, 0, 0, 0, 0, 224, 65}
				payloadmap[overflow_block_index] = []int{4, 0}
				payloadmap[pointerindex] = []int{3}
				payloadmap[max_monochromeindex] = []int{0}
				payloadmap[inverted_colorsindex] = []int{2, 0}
				payloadmap[any_hoverindex] = []int{2}
				payloadmap[updateindex] = []int{3, 0}
				payloadmap[max_heightindex] = []int{38, 24}
				payloadmap[scriptingindex] = []int{3, 0}
				payloadmap[prefers_contrastindex] = []int{3, 0}
				payloadmap[prefers_color_schemeindex] = []int{3}
				payloadmap[orientationindex] = []int{2}
				payloadmap[hoverindex] = []int{2}
				payloadmap[forced_colorsindex] = []int{2}
				payloadmap[overflow_inlineindex] = []int{2, 0}
				payloadmap[max_widthindex] = []int{57, 10}
				payloadmap[max_colorindex] = []int{8}
				payloadmap[display_modeindex] = []int{4}
				payloadmap[color_gamutindex] = []int{3}
				payloadmap[prefers_reduced_motionindex] = []int{2}
				payloadmap[gridindex] = []int{2}
				payloadmap[any_pointerindex] = []int{3}
				payloadmap[scanindex] = []int{2, 0}
				payloadmap[light_levelindex] = []int{3, 0}
				payloadmap[prefers_reduced_transparencyindex] = []int{2, 0}
				payloadmap[max_color_indexindex] = []int{0}

				lastblock := outstring

				for i := len(a) - 1; i >= 0; i-- {
					val := a[i]
					switch val {
					case hoverindex: //
						blocks := strings.Split(lastblock, ",\"hover")
						lastblock = blocks[0]
						if strings.Index(blocks[len(blocks)-1], "]<0") > -1 {
							payloadmap[hoverindex] = append(payloadmap[hoverindex], 2)
						} else {
							payloadmap[hoverindex] = append(payloadmap[hoverindex], 1)
						}
					case scanindex:
						blocks := strings.Split(lastblock, "scan")
						lastblock = blocks[0]
					case forced_colorsindex: //
						blocks := strings.Split(lastblock, "forced-colors")
						lastblock = blocks[0]
						if strings.Index(blocks[len(blocks)-1], "]<0") > -1 {
							payloadmap[forced_colorsindex] = append(payloadmap[forced_colorsindex], 1)
						} else {
							payloadmap[forced_colorsindex] = append(payloadmap[forced_colorsindex], 2)
						}
					case any_pointerindex: //
						blocks := strings.Split(lastblock, "any-pointer")
						lastblock = blocks[0]
						if strings.Index(blocks[len(blocks)-1], "]<0") > -1 {
							payloadmap[any_pointerindex] = append(payloadmap[any_pointerindex], 4)
						} else {
							payloadmap[any_pointerindex] = append(payloadmap[any_pointerindex], 1)
						}
					case inverted_colorsindex:
						blocks := strings.Split(lastblock, "inverted-colors")
						lastblock = blocks[0]
					case pointerindex: //
						blocks := strings.Split(lastblock, "\"pointer")
						lastblock = blocks[0]
						if strings.Index(blocks[len(blocks)-1], "]<0") > -1 {
							payloadmap[pointerindex] = append(payloadmap[pointerindex], 4)
						} else {
							payloadmap[pointerindex] = append(payloadmap[pointerindex], 1)
						}
					case orientationindex: //
						blocks := strings.Split(lastblock, "orientation")
						lastblock = blocks[0]
						if strings.Index(blocks[len(blocks)-1], "]<0") > -1 {
							payloadmap[orientationindex] = append(payloadmap[orientationindex], 2)
						} else {
							payloadmap[orientationindex] = append(payloadmap[orientationindex], 1)
						}
					case light_levelindex:
						blocks := strings.Split(lastblock, "light-level")
						lastblock = blocks[0]
					case overflow_inlineindex:
						blocks := strings.Split(lastblock, "overflow-inline")
						lastblock = blocks[0]
					case updateindex:
						blocks := strings.Split(lastblock, "orientation")
						lastblock = blocks[0]
					case any_hoverindex: //
						blocks := strings.Split(lastblock, "any-hover")
						lastblock = blocks[0]
						if strings.Index(blocks[len(blocks)-1], "]<0") > -1 {
							payloadmap[any_hoverindex] = append(payloadmap[any_hoverindex], 2)
						} else {
							payloadmap[any_hoverindex] = append(payloadmap[any_hoverindex], 1)
						}
					case scriptingindex:
						blocks := strings.Split(lastblock, "scripting")
						lastblock = blocks[0]
					case gridindex: //
						blocks := strings.Split(lastblock, "grid")
						lastblock = blocks[0]
						if strings.Index(blocks[len(blocks)-1], "]<0") > -1 {
							payloadmap[gridindex] = append(payloadmap[gridindex], 1)
						} else {
							payloadmap[gridindex] = append(payloadmap[gridindex], 2)
						}
					case prefers_color_schemeindex: //
						blocks := strings.Split(lastblock, "prefers-color-scheme")
						lastblock = blocks[0]
						if strings.Index(blocks[len(blocks)-1], "]<0") > -1 {
							payloadmap[prefers_color_schemeindex] = append(payloadmap[prefers_color_schemeindex], 4)
						} else {
							payloadmap[prefers_color_schemeindex] = append(payloadmap[prefers_color_schemeindex], 1)
						}
					case color_gamutindex: //
						blocks := strings.Split(lastblock, "color-gamut")
						lastblock = blocks[0]
						if strings.Index(blocks[len(blocks)-1], "]<0") > -1 {
							payloadmap[color_gamutindex] = append(payloadmap[color_gamutindex], 1)
						} else {
							payloadmap[color_gamutindex] = append(payloadmap[color_gamutindex], 4)
						}
					case prefers_contrastindex:
						blocks := strings.Split(lastblock, "prefers-contrast")
						lastblock = blocks[0]
					case prefers_reduced_motionindex: //
						blocks := strings.Split(lastblock, "prefers-reduced-motion")
						lastblock = blocks[0]
						if strings.Index(blocks[len(blocks)-1], "]<0") > -1 {
							payloadmap[prefers_reduced_motionindex] = append(payloadmap[prefers_reduced_motionindex], 1)
						} else {
							payloadmap[prefers_reduced_motionindex] = append(payloadmap[prefers_reduced_motionindex], 2)
						}
					case overflow_block_index:
						blocks := strings.Split(lastblock, "overflow-block")
						lastblock = blocks[0]
					case display_modeindex: //
						blocks := strings.Split(lastblock, "display-mode")
						lastblock = blocks[0]
						if strings.Index(blocks[len(blocks)-1], "]<0") > -1 {
							payloadmap[display_modeindex] = append(payloadmap[display_modeindex], 8)
						} else {
							payloadmap[display_modeindex] = append(payloadmap[display_modeindex], 1)
						}
					case prefers_reduced_transparencyindex:
						blocks := strings.Split(lastblock, "prefers-reduced-transparency")
						lastblock = blocks[0]
					}
				}

				for _, val := range a {
					g.P42VAL = append(g.P42VAL, payloadmap[val]...)
				}

				g.PayloadOrder = append(g.PayloadOrder, 42)
				for k, _ := range mapunion {
					g.KeyCount[k]++
					g.KeyHolderRaw[42] = append(g.KeyHolderRaw[42], k)
				}
			} else
			if strings.Contains(outstring, "0xFFFFFFFFFFFFFBFF") { //payload44
				g.PayloadOrder = append(g.PayloadOrder, 43)
				for k, _ := range mapunion {
					g.KeyCount[k]++
					g.KeyHolderRaw[43] = append(g.KeyHolderRaw[43], k)
				}
			} else
			if strings.Contains(outstring, "documentRelativeX") { //payload45
				var first int
				toarrayholder := []string{}
				for _, match := range P44TOARRAYMATCH.FindAllStringSubmatch(outstring, 3) {
					toarrayholder = append(toarrayholder, match[1])
				}

				ctag1 := fmt.Sprintf("[%s][length]", toarrayholder[0]) //all mouse
				ctag2 := fmt.Sprintf("[%s][length]", toarrayholder[1]) //filtered mouse if(0;) descending
				ctag3 := fmt.Sprintf("[%s][length]", toarrayholder[2]) //event

				b1index := strings.LastIndex(outstring, ctag1)
				b2index := strings.LastIndex(outstring, ctag2)
				b3index := strings.LastIndex(outstring, ctag3)
				b1flipped := true
				b2flipped := true
				b3flipped := true

				var block1, block2, block3 string

				th := []int{b1index, b2index, b3index}
				sort.Ints(th)
				lastblock := outstring[th[0]:]

				for i := 2; i > -1; i-- {
					val := th[i]
					switch val {
					case b1index:
						blocks := strings.Split(lastblock, ctag1)
						lastblock = blocks[0]
						block1 = blocks[len(blocks)-1]
						if strings.Index(block1, "]<0") > -1 {
							b1flipped = false
						}

					case b2index:
						blocks := strings.Split(lastblock, ctag2)
						lastblock = blocks[0]
						block2 = blocks[len(blocks)-1]
						if strings.Index(block2, "]<0") > -1 {
							b1flipped = false
						}

					case b3index:
						blocks := strings.Split(lastblock, ctag3)
						lastblock = blocks[0]
						block3 = blocks[len(blocks)-1]
						if strings.Index(block3, "]<0") > -1 {
							b1flipped = false
						}
					}
				}

				check1index := strings.Index(block3, "[\"timestamp\"]")
				check2index := strings.Index(block3, "[\"targetRelativeX\"]")
				check3index := strings.Index(block3, "[\"documentRelativeX\"]")
				check4index := strings.Index(block3, "[\"targetName\"]")
				check5index := strings.Index(block3, "[\"targetRelativeY\"]")
				check6index := strings.Index(block3, "[\"button\"]")
				check7index := strings.Index(block3, "[\"eventType\"]")
				check8index := strings.Index(block3, "<< 1")
				check9index := strings.Index(block3, "[\"documentRelativeY\"]")
				check10index := strings.Index(block3, "[\"targetId\"]")
				if check4index > check10index {
					first = check10index
				} else {
					first = check4index
				}

				check11index := strings.Index(block1, "[\"documentRelativeY\"]") //all
				check12index := strings.Index(block1, "[\"documentRelativeX\"]")
				check13index := strings.Index(block1, "[\"timestamp\"]")
				check14index := strings.Index(block1, "<< 1")

				check15index := strings.Index(block2, "[\"documentRelativeY\"]") //filtered
				check16index := strings.Index(block2, "[\"documentRelativeX\"]")
				check17index := strings.Index(block2, "[\"timestamp\"]")
				check18index := strings.Index(block2, "<< 1")

				blockresolutions := []int{b1index, b2index, b3index}
				sort.Ints(blockresolutions)

				block1vals := []int{check11index, check12index, check13index, check14index}                                                                            //all
				block2vals := []int{check15index, check16index, check17index, check18index}                                                                            //filtered
				block3vals := []int{check1index, check2index, check3index, check4index, check5index, check6index, check7index, check8index, check9index, check10index} //keyevents
				sort.Ints(block1vals)
				sort.Ints(block2vals)
				sort.Ints(block3vals)

				filtered, events, all := generation.MouseGen()

				for _, val := range blockresolutions {
					switch val {
					case b1index:
						g.P44VAL = append(g.P44VAL, generation.BaseNumEnc(len(all))...)
						if b1flipped {
							for i := len(all) - 1; i >= 0; i-- {
								payload := all[i]
								for _, si := range block1vals {
									switch si {
									case check11index:
										for _, sv := range generation.BaseNumEnc(payload.DocumentRelativeY) {
											g.P44VAL = append(g.P44VAL, sv)
										}
									case check12index:
										for _, sv := range generation.BaseNumEnc(payload.DocumentRelativeX) {
											g.P44VAL = append(g.P44VAL, sv)
										}
									case check13index:
										for _, sv := range generation.BaseNumEnc(payload.Timestamp) {
											g.P44VAL = append(g.P44VAL, sv)
										}
									case check14index:
										g.P44VAL = append(g.P44VAL, 3)
									}
								}
							}
						} else {
							for _, payload := range all {
								for _, si := range block1vals {
									switch si {
									case check11index:
										for _, sv := range generation.BaseNumEnc(payload.DocumentRelativeY) {
											g.P44VAL = append(g.P44VAL, sv)
										}
									case check12index:
										for _, sv := range generation.BaseNumEnc(payload.DocumentRelativeX) {
											g.P44VAL = append(g.P44VAL, sv)
										}
									case check13index:
										for _, sv := range generation.BaseNumEnc(payload.Timestamp) {
											g.P44VAL = append(g.P44VAL, sv)
										}
									case check14index:
										g.P44VAL = append(g.P44VAL, 3)
									}
								}
							}
						}
					case b2index:
						g.P44VAL = append(g.P44VAL, generation.BaseNumEnc(len(filtered))...)
						if b2flipped {
							for i := len(filtered) - 1; i >= 0; i-- {
								payload := filtered[i]
								for _, si := range block2vals {
									switch si {
									case check15index:
										for _, sv := range generation.BaseNumEnc(payload.DocumentRelativeY) {
											g.P44VAL = append(g.P44VAL, sv)
										}
									case check16index:
										for _, sv := range generation.BaseNumEnc(payload.DocumentRelativeX) {
											g.P44VAL = append(g.P44VAL, sv)
										}
									case check17index:
										for _, sv := range generation.BaseNumEnc(payload.Timestamp) {
											g.P44VAL = append(g.P44VAL, sv)
										}
									case check18index:
										g.P44VAL = append(g.P44VAL, 3)
									}
								}
							}
						} else {
							for _, payload := range filtered {
								for _, si := range block2vals {
									switch si {
									case check15index:
										for _, sv := range generation.BaseNumEnc(payload.DocumentRelativeY) {
											g.P44VAL = append(g.P44VAL, sv)
										}
									case check16index:
										for _, sv := range generation.BaseNumEnc(payload.DocumentRelativeX) {
											g.P44VAL = append(g.P44VAL, sv)
										}
									case check17index:
										for _, sv := range generation.BaseNumEnc(payload.Timestamp) {
											g.P44VAL = append(g.P44VAL, sv)
										}
									case check18index:
										g.P44VAL = append(g.P44VAL, 3)
									}
								}
							}
						}
					case b3index:
						g.P44VAL = append(g.P44VAL, generation.BaseNumEnc(len(events))...)
						if b3flipped {
							for i := len(events) - 1; i >= 0; i-- {
								payload := events[i]
								for _, si := range block3vals {
									switch si {
									case check1index:
										//[\"timestamp\"]
										for _, sv := range generation.BaseNumEnc(payload.Timestamp) {
											g.P44VAL = append(g.P44VAL, sv)
										}
									case check2index:
										//[\"targetRelativeX\"]
										for _, sv := range generation.BaseNumEnc(payload.TargetRelativeX) {
											g.P44VAL = append(g.P44VAL, sv)
										}
									case check3index:
										//[\"documentRelativeX\"]
										for _, sv := range generation.BaseNumEnc(payload.DocumentRelativeX) {
											g.P44VAL = append(g.P44VAL, sv)
										}
									case check4index:
										//[\"targetName\"]
										if check4index == first {
											if i == 5 {
												g.P44VAL = append(g.P44VAL, 0)
											} else
											if i == 2 {
												for _, char := range payload.TargetName {
													g.P44VAL = append(g.P44VAL, int(char))
												}
												g.P44VAL = append(g.P44VAL, 0)
											} else {
												if i > 2 {
													g.P44VAL = append(g.P44VAL, 129, 0)
												} else {
													g.P44VAL = append(g.P44VAL, 129, 2)
												}
											}
										} else {
											if i == 5 {
												g.P44VAL = append(g.P44VAL, 0)
											} else
											if i > 2 {
												g.P44VAL = append(g.P44VAL, 129, 1)
											} else {
												g.P44VAL = append(g.P44VAL, 129, 2)
											}
										}
									case check5index:
										//[\"targetRelativeY\"]
										for _, sv := range generation.BaseNumEnc(payload.TargetRelativeY) {
											g.P44VAL = append(g.P44VAL, sv)
										}
									case check6index:
										//[\"button\"]
										for _, sv := range generation.BaseNumEnc(payload.Button) {
											g.P44VAL = append(g.P44VAL, sv)
										}
									case check7index:
										//[\"eventype\"]
										for _, sv := range generation.BaseNumEnc(payload.EventType) {
											g.P44VAL = append(g.P44VAL, sv)
										}
									case check8index:
										//<< 1
										g.P44VAL = append(g.P44VAL, 3)
									case check9index:
										//[\"documentRelativeY\"]
										for _, sv := range generation.BaseNumEnc(payload.DocumentRelativeY) {
											g.P44VAL = append(g.P44VAL, sv)
										}
									case check10index:
										//[\"targetId\"]
										if first == check10index {
											if i == 5 || i == 2 {
												for _, char := range payload.TargetId {
													g.P44VAL = append(g.P44VAL, int(char))
												}
												g.P44VAL = append(g.P44VAL, 0)
											} else {
												if i > 2 {
													g.P44VAL = append(g.P44VAL, 129, 0)
												} else {
													g.P44VAL = append(g.P44VAL, 129, 2)
												}
											}
										} else {
											if i == 5 {
												for _, char := range payload.TargetId {
													g.P44VAL = append(g.P44VAL, int(char))
												}
												g.P44VAL = append(g.P44VAL, 0)
											} else {
												if i > 2 {
													g.P44VAL = append(g.P44VAL, 129, 1)
												} else {
													g.P44VAL = append(g.P44VAL, 129, 2)
												}
											}
										}
									}
								}
							}
						} else {
							for i, payload := range events {
								for _, si := range block3vals {
									switch si {
									case check1index:
										//[\"timestamp\"]
										for _, sv := range generation.BaseNumEnc(payload.Timestamp) {
											g.P44VAL = append(g.P44VAL, sv)
										}
									case check2index:
										//[\"targetRelativeX\"]
										for _, sv := range generation.BaseNumEnc(payload.TargetRelativeX) {
											g.P44VAL = append(g.P44VAL, sv)
										}
									case check3index:
										//[\"documentRelativeX\"]
										for _, sv := range generation.BaseNumEnc(payload.DocumentRelativeX) {
											g.P44VAL = append(g.P44VAL, sv)
										}
									case check4index:
										//[\"targetName\"]
										if check4index == first {
											if i == 3 {
												g.P44VAL = append(g.P44VAL, 0)
											} else
											if i == 0 {
												for _, char := range payload.TargetName {
													g.P44VAL = append(g.P44VAL, int(char))
												}
												g.P44VAL = append(g.P44VAL, 0)
											} else {
												if i < 2 {
													g.P44VAL = append(g.P44VAL, 129, 0)
												} else {
													g.P44VAL = append(g.P44VAL, 129, 1)
												}
											}
										} else {
											if i == 3 {
												g.P44VAL = append(g.P44VAL, 0)
											} else
											if i < 2 {
												g.P44VAL = append(g.P44VAL, 129, 0)
											} else {
												g.P44VAL = append(g.P44VAL, 129, 2)
											}
										}
									case check5index:
										//[\"targetRelativeY\"]
										for _, sv := range generation.BaseNumEnc(payload.TargetRelativeY) {
											g.P44VAL = append(g.P44VAL, sv)
										}
									case check6index:
										//[\"button\"]
										for _, sv := range generation.BaseNumEnc(payload.Button) {
											g.P44VAL = append(g.P44VAL, sv)
										}
									case check7index:
										//[\"eventype\"]
										for _, sv := range generation.BaseNumEnc(payload.EventType) {
											g.P44VAL = append(g.P44VAL, sv)
										}
									case check8index:
										//<< 1
										g.P44VAL = append(g.P44VAL, 3)
									case check9index:
										//[\"documentRelativeY\"]
										for _, sv := range generation.BaseNumEnc(payload.DocumentRelativeY) {
											g.P44VAL = append(g.P44VAL, sv)
										}
									case check10index:
										//[\"targetId\"]
										if first == check10index {
											if i == 3 || i == 0 {
												for _, char := range payload.TargetId {
													g.P44VAL = append(g.P44VAL, int(char))
												}
												g.P44VAL = append(g.P44VAL, 0)
											} else {
												if i < 2 {
													g.P44VAL = append(g.P44VAL, 129, 0)
												} else {
													g.P44VAL = append(g.P44VAL, 129, 1)
												}
											}
										} else {
											if i == 3 {
												for _, char := range payload.TargetId {
													g.P44VAL = append(g.P44VAL, int(char))
												}
												g.P44VAL = append(g.P44VAL, 0)
											} else {
												if i < 2 {
													g.P44VAL = append(g.P44VAL, 129, 0)
												} else {
													g.P44VAL = append(g.P44VAL, 129, 2)
												}
											}
										}
									}
								}
							}
						}
					}
				}

				g.PayloadOrder = append(g.PayloadOrder, 44)
				for k, _ := range mapunion {
					g.KeyCount[k]++
					g.KeyHolderRaw[44] = append(g.KeyHolderRaw[44], k)
				}
			} else
			{
				//log.Println("NOTFOUND", ctag)
			}
		}
	}
	if strings.Contains(outstring, "[\"c\"]=") {
		g.Dheader = dheadermatch.FindStringSubmatch(outstring)[1]
	}

	//if ctag == "51584_2" {
	//	ioutil.WriteFile("./gentest/"+ctag+".js", []byte(starttag+outstring), 0644)
	//}

	g.IteratedFiles[ctag] = true
	for _, val := range g.Initializerholder[:] {
		g.FileLinker(val.CopyIndex, val.X, val.Y, val.Heap)
	}
}
func (g *GlobalHolder) GenerateHeaders() map[string]string {
	rand.Seed(time.Now().UnixNano())
	resString, err := g.getResString()
	if err != nil {
		return nil
	}

	globalKeySet := []int{}
	for i := 1; i < 5; i++ {
		appint := g.GlobalKeyArray[g.EncryptionKeyHolder[i]]
		if g.Negatives[g.EncryptionKeyHolder[i]] {
			appint = -appint
		}
		globalKeySet = append(globalKeySet, int(appint))
	}

	headerstruct := generation.VersionPayloadHolder{
		Hashseedbase: generation.StringHash(g.Seedstring),
		Hashseed1:    generation.StringHash(g.Seedstring),
		Hashseed2:    0,
		Currval:      0,
		Seedcount:    0,
		Keyholder:    g.KeyHolder,
		GenOrder:     g.PayloadOrder,
		BaseFile:     g.Base,
		SeedString:   g.Seedstring,
		Alphabet:     g.Alphabet,
		HashString:   g.Hashstring,
		Dheader:      g.Dheader,
		ResString:    resString,
		GlobalKeys:   globalKeySet,
		BaseKeys:     g.Basekeys,
		Calckey:      int(g.GlobalKeyArray[g.EncryptionKeyHolder[0]]),
		Encrounds:    g.Encrounds,
		P1VAL:        g.P1VAL,
		P2VAL:        g.P2VAL,
		P3VAL:        g.P3VAL,
		P4VAL:        g.P4VAL,
		P5VAL:        g.P5VAL,
		P11VAL:       g.P11VAL,
		P13VAL:       g.P13VAL,
		P14VAL:       g.P14VAL,
		P15VAL:       g.P15VAL,
		P16VAL:       g.P16VAL,
		P17VAL:       g.P17VAL,
		P18VAL:       g.P18VAL,
		P19VAL:       g.P19VAL,
		P20VAL:       g.P20VAL,
		P21VAL:       g.P21VAL,
		P23VAL:       g.P23VAL,
		P24VAL:       g.P24VAL,
		P25VAL:       g.P25VAL,
		P30VAL:       g.P30VAL,
		P31VAL:       g.P31VAL,
		P32VAL:       g.P32VAL,
		P35VAL:       g.P35VAL,
		P36VAL:       g.P36VAL,
		P40VAL:       g.P40VAL,
		P42VAL:       g.P42VAL,
		P44VAL:       g.P44VAL,
	}
	return headerstruct.GenerateHeaders()
}
