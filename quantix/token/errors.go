package token

type Type int

const (
	Error Type = iota
	EOF        //$
	T_0        // !
	T_1        // (
	T_2        // )
	T_3        // .
	T_4        // :
	T_5        // ;
	T_6        // <
	T_7        // >
	T_8        // [
	T_9        // ]
	T_10       // any
	T_11       // char_lit
	T_12       // empty
	T_13       // letter
	T_14       // lowcase
	T_15       // not
	T_16       // nt
	T_17       // number
	T_18       // package
	T_19       // string_lit
	T_20       // tokid
	T_21       // upcase
	T_22       // {
	T_23       // |
	T_24       // }
)

var TypeToString = []string{
	"Error",
	"EOF",
	"T_0",
	"T_1",
	"T_2",
	"T_3",
	"T_4",
	"T_5",
	"T_6",
	"T_7",
	"T_8",
	"T_9",
	"T_10",
	"T_11",
	"T_12",
	"T_13",
	"T_14",
	"T_15",
	"T_16",
	"T_17",
	"T_18",
	"T_19",
	"T_20",
	"T_21",
	"T_22",
	"T_23",
	"T_24",
}

var StringToType = map[string]Type{
	"Error": Error,
	"EOF":   EOF,
	"T_0":   T_0,
	"T_1":   T_1,
	"T_2":   T_2,
	"T_3":   T_3,
	"T_4":   T_4,
	"T_5":   T_5,
	"T_6":   T_6,
	"T_7":   T_7,
	"T_8":   T_8,
	"T_9":   T_9,
	"T_10":  T_10,
	"T_11":  T_11,
	"T_12":  T_12,
	"T_13":  T_13,
	"T_14":  T_14,
	"T_15":  T_15,
	"T_16":  T_16,
	"T_17":  T_17,
	"T_18":  T_18,
	"T_19":  T_19,
	"T_20":  T_20,
	"T_21":  T_21,
	"T_22":  T_22,
	"T_23":  T_23,
	"T_24":  T_24,
}

var TypeToID = []string{
	"Error",
	"$",
	"!",
	"(",
	")",
	".",
	":",
	";",
	"<",
	">",
	"[",
	"]",
	"any",
	"char_lit",
	"empty",
	"letter",
	"lowcase",
	"not",
	"nt",
	"number",
	"package",
	"string_lit",
	"tokid",
	"upcase",
	"{",
	"|",
	"}",
}

var Suppress = []bool{
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
}

func (t Type) String() string {
	return TypeToString[t]
}

// ID returns the token type ID of token Type t
func (t Type) ID() string {
	return TypeToID[t]
}
