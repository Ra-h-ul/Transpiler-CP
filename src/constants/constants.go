package constants

const CppHasStdFormat = false

const VersionString = "go transpiler"

const TupleType = "std::tuple"

const (
	HashMapSuffix = "_h__"
	KeysSuffix    = "_k__"
	SwitchPrefix  = "_s__"
	LabelPrefix   = "_l__"
	DeferPrefix   = "_d__"
)

var IncludeMap = map[string]string{
	"std::tuple":                       "",
	"std::endl":                        "",
	"std::cout":                        "",
	"std::string":                      "",
	"std::size":                        "",
	"std::unordered_map":               "",
	"std::hash":                        "",
	"std::size_t":                      "",
	"std::int8_t":                      "",
	"std::int16_t":                     "",
	"std::int32_t":                     "",
	"std::int64_t":                     "",
	"std::uint8_t":                     "",
	"std::uint16_t":                    "",
	"std::uint32_t":                    "",
	"std::uint64_t":                    "",
	"printf":                           "",
	"fprintf":                          "",
	"sprintf":                          "",
	"snprintf":                         "",
	"std::stringstream":                "",
	"std::is_pointer":                  "",
	"std::experimental::is_detected_v": "",
	"std::shared_ptr":                  "",
	"std::nullopt":                     "",
	"EXIT_SUCCESS":                     "",
	"EXIT_FAILURE":                     "",
	"std::vector":                      "",
	"std::unique_ptr":                  "",
	"std::runtime_error":               "",
	"std::regex_replace":               "",
	"std::regex_constants":             "",
	"std::to_string":                   "",
}

var Endings = []string{"{", ",", "}", ":"}

var (
	SwitchExpressionCounter = -1
	FirstCase               bool
	SwitchLabel             string
	LabelCounter            int
	IotaNumber              int // used for simple increases of iota constants
	DeferCounter            int
	UnfinishedDeferFunction bool
)
