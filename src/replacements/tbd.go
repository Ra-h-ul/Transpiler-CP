package replacements

import "strings"

func AddFunctions(source string, useFormatOutput, haveStructs bool) (output string) {

	// TODO: Make the fmtSprintf implementation more watertight. Use variadic templates and parameter packs, while waiting for std::format to arrive in the C++20 implementations.

	output = source
	replacements := map[string]string{
		"strconv.ParseFloat": `using error = std::optional<std::string>;
auto strconvParseFloat(std::string s, int bitSize) -> std::tuple<double, error> {
	try {
		return std::tuple { std::stod(s), std::nullopt };
	} catch (const std::invalid_argument& ia) {
		return std::tuple { 0.0, std::optional { "invalid argument" } };
	}
}
`,
		"strconv.ParseInt": `using error = std::optional<std::string>;
auto strconvParseInt(std::string s, int base, int bitSize) -> std::tuple<int, error> {
	try {
		return std::tuple { std::stoi(s), std::nullopt };
	} catch (const std::invalid_argument& ia) {
		return std::tuple { 0.0, std::optional { "invalid argument" } };
	}
}
`,
		"strings.Contains":  `inline auto stringsContains(std::string const& haystack, std::string const& needle) -> bool { return haystack.find(needle) != std::string::npos; }`,
		"strings.HasPrefix": `inline auto stringsHasPrefix(std::string const& haystack, std::string const& prefix) -> auto { return 0 == haystack.find(prefix); }`,
		"_format_output": `template<typename T>
using _str_t = decltype( std::declval<T&>()._str() );

template<typename T>
using _p_str_t = decltype( std::declval<T&>()->_str() );

template <typename T> void _format_output(std::ostream& out, T x)
{
    if constexpr (std::is_same<T, bool>::value) {
        out << std::boolalpha << x << std::noboolalpha;
    } else if constexpr (std::is_integral<T>::value) {
        out << static_cast<int>(x);
    } else if constexpr (std::is_object<T>::value && !std::is_pointer<T>::value && std::experimental::is_detected_v<_str_t, T>) {
        out << x._str();
    } else if constexpr (std::is_object<T>::value && std::is_pointer<T>::value && std::experimental::is_detected_v<_p_str_t, T>) {
        out << "&" << x->_str();
    } else {
        out << x;
    }
}`,
		"strings.TrimSpace": `inline auto stringsTrimSpace(std::string const& s) -> std::string { std::string news {}; for (auto l : s) { if (l != ' ' && l != '\n' && l != '\t' && l != '\v' && l != '\f' && l != '\r') { news += l; } } return news; }`,
		"fmt.Sprintf": `
template <typename T>
std::string fmtSprintf(const std::string& fmt, T arg1)
{
    std::string tmp = fmt;
    tmp = std::regex_replace(tmp, std::regex("%(s|v|d|f)"), std::to_string(arg1), std::regex_constants::format_first_only);
    return tmp;
}
template <typename T>
std::string fmtSprintf(const std::string& fmt, T arg1, T arg2)
{
    std::string tmp = fmt;
    tmp = std::regex_replace(tmp, std::regex("%(s|v|d|f)"), std::to_string(arg1), std::regex_constants::format_first_only);
    tmp = std::regex_replace(tmp, std::regex("%(s|v|d|f)"), std::to_string(arg2), std::regex_constants::format_first_only);
    return tmp;
}
template <typename T>
std::string fmtSprintf(const std::string& fmt, T arg1, T arg2, T arg3)
{
    std::string tmp = fmt;
    tmp = std::regex_replace(tmp, std::regex("%(s|v|d|f)"), std::to_string(arg1), std::regex_constants::format_first_only);
    tmp = std::regex_replace(tmp, std::regex("%(s|v|d|f)"), std::to_string(arg2), std::regex_constants::format_first_only);
    tmp = std::regex_replace(tmp, std::regex("%(s|v|d|f)"), std::to_string(arg3), std::regex_constants::format_first_only);
    return tmp;
}
template <typename T>
std::string fmtSprintf(const std::string& fmt, T arg1, T arg2, T arg3, T arg4)
{
    std::string tmp = fmt;
    tmp = std::regex_replace(tmp, std::regex("%(s|v|d|f)"), std::to_string(arg1), std::regex_constants::format_first_only);
    tmp = std::regex_replace(tmp, std::regex("%(s|v|d|f)"), std::to_string(arg2), std::regex_constants::format_first_only);
    tmp = std::regex_replace(tmp, std::regex("%(s|v|d|f)"), std::to_string(arg3), std::regex_constants::format_first_only);
    tmp = std::regex_replace(tmp, std::regex("%(s|v|d|f)"), std::to_string(arg4), std::regex_constants::format_first_only);
    return tmp;
}
template <typename T>
std::string fmtSprintf(const std::string& fmt, T arg1, T arg2, T arg3, T arg4, T arg5)
{
    std::string tmp = fmt;
    tmp = std::regex_replace(tmp, std::regex("%(s|v|d|f)"), std::to_string(arg1), std::regex_constants::format_first_only);
    tmp = std::regex_replace(tmp, std::regex("%(s|v|d|f)"), std::to_string(arg2), std::regex_constants::format_first_only);
    tmp = std::regex_replace(tmp, std::regex("%(s|v|d|f)"), std::to_string(arg3), std::regex_constants::format_first_only);
    tmp = std::regex_replace(tmp, std::regex("%(s|v|d|f)"), std::to_string(arg4), std::regex_constants::format_first_only);
    tmp = std::regex_replace(tmp, std::regex("%(s|v|d|f)"), std::to_string(arg5), std::regex_constants::format_first_only);
    return tmp;
}
`,
		"len": `
template <typename T>
inline auto len(T x) -> int { return x.size(); }
`,
	}
	if useFormatOutput && !haveStructs {
		replacements["_format_output"] = `template <typename T> void _format_output(std::ostream& out, T x)
		
			{
			    if constexpr (std::is_same<T, bool>::value) {
			        out << std::boolalpha << x << std::noboolalpha;
			    } else if constexpr (std::is_integral<T>::value) {
			        out << static_cast<int>(x);
			    } else {
			        out << x;
			    }
			}`
		replacements["_format_output"] = ""
	}
	for k, v := range replacements {
		if strings.Contains(output, k) {
			output = strings.Replace(output, k, strings.Replace(k, ".", "", -1), -1)
			output = v + "\n" + output
		}
	}
	return output
}
