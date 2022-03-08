/*  Common codes to Smart JSON-YAML functions
    (C) 2021-2022 Péter Deák (hyper80@gmail.com)
    License: Apache 2.0
*/

// Do not use this package separated. It uses by smartjson and smartyaml packages
// It contains the common codes of that packages
package smartjsonyamlstub

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// SmartJsonYamlConfig is contains configurations for SmartJSON and SmartYAML
// This structure is available in that stuctures with name Config
type SmartJsonYamlConfig struct {
	// NotFoundOrInvalidNotation contains the string which passed by query
	// functions as type when the requested node is not found. The default value is "none"
	NotFoundOrInvalidNotation string
	// YamlGeneratorIndenter holds the indent string used by Yaml generator function.
	// The indentation is only apply to maps. The default is two spaces.
	YamlGeneratorIndenter string
	// If the YamlAlwaysUseQuotesForString value is true the Yaml generator
	// put every string in quotes. The default is false.
	YamlAlwaysUseQuotesForString bool
	// If the OutputMapKeyOrder string array is not empty
	// the Yaml and Json generator functions prefers this string order in map type nodes.
	// This is a workaround to get rid of side effects of randomized go maps, which is used by parsers.
	OutputMapKeyOrder []string
}

// SmartJsonYamlBase is the base structure of SmartJSON and SmartYAML
// The common used function connected to this structure
type SmartJsonYamlBase struct {
	// ParsedData holds the unmarsaled data structures
	ParsedData interface{}
	// ParsedFrom contains the name the original form of the parsed data. "yaml" or "json"
	ParsedFrom string
	// Config holds configurable items
	Config SmartJsonYamlConfig
}

// InitConfig sets the default values of Config
func (conf *SmartJsonYamlConfig) InitConfig() {
	conf.NotFoundOrInvalidNotation = "none"
	conf.YamlGeneratorIndenter = "  "
	conf.YamlAlwaysUseQuotesForString = false
	conf.OutputMapKeyOrder = []string{}
}

// Yaml generates a yaml output
func (sjyb SmartJsonYamlBase) Yaml() (out string) {
	return sjyb.yamlNodeToString(sjyb.ParsedData, "", "top") + "\n"
}

// JsonIndented generates an indented JSON
func (sjyb SmartJsonYamlBase) JsonIndented() (out string) {
	return sjyb.jsonNodeToString(sjyb.ParsedData, "", true) + "\n"
}

// JsonIndented generates an compacted JSON
func (sjyb SmartJsonYamlBase) JsonCompacted() (out string) {
	return sjyb.jsonNodeToString(sjyb.ParsedData, "", false)
}

func (sjyb SmartJsonYamlBase) jsonNodeToString(v interface{}, indent string, prettyOutput bool) (out string) {
	out = ""
	if m, isMap := v.(map[string]interface{}); isMap {
		if prettyOutput {
			out += "{\n" + indent + "  "
		} else {
			out += "{"
		}
		c := 0

		done := []string{}
		for _, orderedKey := range sjyb.Config.OutputMapKeyOrder {
			if _, ok := m[orderedKey]; ok {
				sep := ""
				if c > 0 {
					if prettyOutput {
						sep = ",\n  " + indent
					} else {
						sep = ","
					}
				}
				out += sep + "\"" + n + "\":" + sjyb.jsonNodeToString(v, indent+"  ", prettyOutput)
				done = append(done, orderedKey)
				c++
			}
		}

		for n, v := range m {
			if contains(done, n) {
				continue
			}
			sep := ""
			if c > 0 {
				if prettyOutput {
					sep = ",\n  " + indent
				} else {
					sep = ","
				}
			}
			out += sep + "\"" + n + "\":" + sjyb.jsonNodeToString(v, indent+"  ", prettyOutput)
			c++
		}
		if prettyOutput {
			out += "\n" + indent + "}"
		} else {
			out += "}"
		}
		return out
	}
	if arr, isArray := v.([]interface{}); isArray {
		if prettyOutput {
			out += "[\n" + indent + "  "
		} else {
			out += "["
		}
		l := len(arr)
		for i := 0; i < l; i++ {
			sep := ""
			if i > 0 {
				if prettyOutput {
					sep = ",\n  " + indent
				} else {
					sep = ","
				}
			}
			out += sep + sjyb.jsonNodeToString(arr[i], indent+"  ", prettyOutput)
		}
		if prettyOutput {
			out += "\n" + indent + "]"
		} else {
			out += "]"
		}
		return out
	}
	if str, isStr := v.(string); isStr {
		out += "\"" + sjyb.jsonStringToOutput(str) + "\""
		return out
	}
	if intval, isInt := v.(int); isInt {
		out += fmt.Sprintf("%d", intval)
		return out
	}
	if flt, isFlt := v.(float64); isFlt {
		out += strconv.FormatFloat(flt, 'g', 10, 64)
		return out
	}
	if timeval, isTime := v.(time.Time); isTime {
		if timeval.Hour() == 0 && timeval.Minute() == 0 && timeval.Second() == 0 && timeval.Nanosecond() == 0 {
			out += "\"" + fmt.Sprintf("%s", timeval.Format("2006-01-02")) + "\""
		} else {
			out += "\"" + fmt.Sprintf("%s", timeval.Format("2006-01-02 15:04:05")) + "\""
		}
		return out
	}
	if b, isBool := v.(bool); isBool {
		if b {
			out += "true"
		} else {
			out += "false"
		}
		return out
	}
	if v == nil {
		out += "null"
		return out
	}
	return ""
}

func (sjyb SmartJsonYamlBase) jsonStringToOutput(str string) string {
	str = strings.Replace(str, "\"", "\\\"", -1)
	return str
}

func (sjyb SmartJsonYamlBase) yamlNodeToString(v interface{}, pindent string, parent string) (out string) {
	out = ""
	if parent == "top" {
		out += "---\n"
	}
	if m, isMap := v.(map[string]interface{}); isMap {
		if parent == "map" {
			out += "\n"
		}
		addindent := ""
		if parent == "map" {
			addindent = sjyb.Config.YamlGeneratorIndenter
		}
		c := 0

		done := []string{}
		for _, orderedKey := range sjyb.Config.OutputMapKeyOrder {
			if _, ok := m[orderedKey]; ok {
				if parent != "array" || c != 0 {
					out += pindent + addindent
				}
				out += orderedKey + ":" + sjyb.yamlNodeToString(m[orderedKey], pindent+addindent, "map")
				done = append(done, orderedKey)
				c++
			}
		}

		for n, v := range m {
			if contains(done, n) {
				continue
			}
			if parent != "array" || c != 0 {
				out += pindent + addindent
			}
			out += n + ":" + sjyb.yamlNodeToString(v, pindent+addindent, "map")
			c++
		}
		return out
	}
	if arr, isArray := v.([]interface{}); isArray {
		if parent == "map" {
			out += "\n"
		}
		l := len(arr)
		for i := 0; i < l; i++ {
			out += pindent
			out += "- " + sjyb.yamlNodeToString(arr[i], pindent+"  ", "array")
		}
		return out
	}

	if v == nil {
		out += "\n"
		return out
	}

	if parent != "array" {
		out += " "
	}

	if str, isStr := v.(string); isStr {
		out += sjyb.yamlStringToOutput(str) + "\n"
		return out
	}
	if intval, isInt := v.(int); isInt {
		out += fmt.Sprintf("%d", intval) + "\n"
		return out
	}
	if flt, isFlt := v.(float64); isFlt {
		if sjyb.ParsedFrom == "json" && flt == math.Floor(flt) {
			out += fmt.Sprintf("%d", int(flt)) + "\n"
			return out
		}
		out += strconv.FormatFloat(flt, 'g', 10, 64) + "\n"
		return out
	}
	if timeval, isTime := v.(time.Time); isTime {
		if timeval.Hour() == 0 && timeval.Minute() == 0 && timeval.Second() == 0 && timeval.Nanosecond() == 0 {
			out += "\"" + fmt.Sprintf("%s", timeval.Format("2006-01-02")) + "\"\n"
		} else {
			out += "\"" + fmt.Sprintf("%s", timeval.Format("2006-01-02 15:04:05")) + "\"\n"
		}
		return out
	}
	if b, isB := v.(bool); isB {
		if b {
			out += "true\n"
		} else {
			out += "false\n"
		}
		return out
	}

	return ""
}

func (sjyb SmartJsonYamlBase) yamlStringToOutput(str string) string {
	needquote := false
	if strings.Contains(str, "\"") ||
		strings.Contains(str, "\\") ||
		strings.Contains(str, ":") ||
		strings.Contains(str, "@") ||
		strings.Contains(str, ",") ||
		strings.Contains(str, "&") ||
		strings.Contains(str, "*") ||
		strings.Contains(str, "#") ||
		strings.Contains(str, "?") ||
		strings.Contains(str, "-") ||
		strings.Contains(str, "!") ||
		strings.Contains(str, "%") ||
		strings.Contains(str, "<") ||
		strings.Contains(str, ">") ||
		strings.Contains(str, "[:") ||
		strings.Contains(str, "]") ||
		strings.Contains(str, "{") ||
		strings.Contains(str, "}") {
		needquote = true
	}

	if str == "Yes" || str == "No" {
		needquote = true
	}

	foundNonNumeric := false
	for _, ch := range str {
		if (ch < '0' || ch > '9') && ch != '.' {
			foundNonNumeric = true
			break
		}
	}
	if !foundNonNumeric {
		needquote = true
	}

	str = strings.Replace(str, "\"", "\\\"", -1)

	if sjyb.Config.YamlAlwaysUseQuotesForString || needquote {
		return "\"" + str + "\""
	}
	return str
}

func (sjyb SmartJsonYamlBase) pathEvalNode(last interface{}) (interface{}, string) {
	if str, isStr := last.(string); isStr {
		return str, "string"
	}
	if flt, isFlt := last.(float64); isFlt {
		return flt, "float64"
	}
	if intval, isInt := last.(int); isInt {
		return intval, "int"
	}
	if bo, isBool := last.(bool); isBool {
		return bo, "bool"
	}
	if timeval, isTime := last.(time.Time); isTime {
		return timeval, "time"
	}
	if mp, isMap := last.(map[string]interface{}); isMap {
		return mp, "map"
	}
	if ar, isArr := last.([]interface{}); isArr {
		return ar, "array"
	}
	if last == nil {
		return nil, "null"
	}
	return nil, sjyb.Config.NotFoundOrInvalidNotation
}

func pathPreprocess(path string) string {
	p := path
	jp := false
	if len(p) > 9 && p[0:9] == "JsonPath:" {
		p = p[9:]
		jp = true
	}

	if len(p) > 2 && p[0:2] == "$." {
		p = p[2:]
		jp = true
	}

	if jp {
		p = strings.Replace(p, ".", "/", -1)
		p = strings.Replace(p, "[", "/[", -1)
	}

	p = strings.Replace(p, "//", "/", -1)

	for len(p) > 2 && p[0:1] == "/" {
		p = p[1:]
	}
	return p
}

// NodeExists return true or false depends on the json/yaml node specified by the path is exists or not
func (sjyb SmartJsonYamlBase) NodeExists(path string) bool {
	_, t := sjyb.GetNodeByPath(path)
	if t == sjyb.Config.NotFoundOrInvalidNotation {
		return false
	}
	return true
}

// GetNodeByPath search the json/yaml node specified by the path and
// returns the value as interface{} and the type as string
func (sjyb SmartJsonYamlBase) GetNodeByPath(path string) (interface{}, string) {
	parts := strings.Split(pathPreprocess(path), "/")
	n := sjyb.ParsedData
	for i := 0; i < len(parts); i++ {
		if map_node, isMap_node := n.(map[string]interface{}); isMap_node {
			map_node_value, ok := map_node[parts[i]]
			if !ok {
				return nil, sjyb.Config.NotFoundOrInvalidNotation
			}
			if i == len(parts)-1 {
				return sjyb.pathEvalNode(map_node_value)
			}
			n = map_node_value
			continue
		}
		if arr_node, isArr_node := n.([]interface{}); isArr_node {
			if len(arr_node) == 0 {
				return nil, sjyb.Config.NotFoundOrInvalidNotation
			}
			var arr_node_item interface{}
			if parts[i] == "[]" {
				arr_node_item = arr_node[0]
			} else {
				r, rxerr := regexp.Compile(`^\[([0-9]+)\]$`)
				if rxerr != nil {
					return nil, sjyb.Config.NotFoundOrInvalidNotation
				}
				matches := r.FindStringSubmatch(parts[i])
				if len(matches) != 2 {
					return nil, sjyb.Config.NotFoundOrInvalidNotation
				}
				index, erratoi := strconv.Atoi(matches[1])
				if erratoi != nil {
					return nil, sjyb.Config.NotFoundOrInvalidNotation
				}
				if index >= len(arr_node) {
					return nil, sjyb.Config.NotFoundOrInvalidNotation
				}
				arr_node_item = arr_node[index]
			}
			if i == len(parts)-1 {
				return sjyb.pathEvalNode(arr_node_item)
			}
			n = arr_node_item
			continue
		}
		return nil, sjyb.Config.NotFoundOrInvalidNotation
	}
	return nil, sjyb.Config.NotFoundOrInvalidNotation
}

// GetSubtreeByPath returns a json/yaml subtree specified by the path
// Use GetSubjsonByPath or GetSubyamlByPath instead of this
func (sjyb SmartJsonYamlBase) GetSubtreeByPath(path string) (SmartJsonYamlBase, string) {
	p, str := sjyb.GetNodeByPath(path)
	s := SmartJsonYamlBase{}
	s.ParsedData = p
	s.Config.InitConfig()
	return s, str
}

// GetMapByPath search a map typed json/yaml node specified by the path and
// returns the value and the type as string
func (sjyb SmartJsonYamlBase) GetMapByPath(path string) (map[string]interface{}, string) {
	val, typ := sjyb.GetNodeByPath(path)
	if m, isMap := val.(map[string]interface{}); typ == "map" && isMap {
		return m, typ
	}
	return nil, sjyb.Config.NotFoundOrInvalidNotation
}

// GetArrayByPath search an array typed json/yaml node specified by the path and
// returns the value and the type as string
func (sjyb SmartJsonYamlBase) GetArrayByPath(path string) ([]interface{}, string) {
	val, typ := sjyb.GetNodeByPath(path)
	if a, isArray := val.([]interface{}); typ == "array" && isArray {
		return a, typ
	}
	return nil, sjyb.Config.NotFoundOrInvalidNotation
}

// GetStringByPath search a string typed json/yaml node specified by the path and
// returns the value and the type as string
func (sjyb SmartJsonYamlBase) GetStringByPath(path string) (string, string) {
	val, typ := sjyb.GetNodeByPath(path)
	if str, isStr := val.(string); typ == "string" && isStr {
		return str, typ
	}
	return "", sjyb.Config.NotFoundOrInvalidNotation
}

// GetFloat64ByPath search a float64 typed json/yaml node specified by the path and
// returns the value and the type as string
func (sjyb SmartJsonYamlBase) GetFloat64ByPath(path string) (float64, string) {
	val, typ := sjyb.GetNodeByPath(path)
	if f, isFlt := val.(float64); typ == "float64" && isFlt {
		return f, typ
	}
	return 0, sjyb.Config.NotFoundOrInvalidNotation
}

// GetIntegerByPath search a float64 typed yaml node specified by the path and
// returns the value and the type as string
func (sjyb SmartJsonYamlBase) GetIntegerByPath(path string) (int, string) {
	val, typ := sjyb.GetNodeByPath(path)
	if i, isInt := val.(int); typ == "int" && isInt {
		return i, typ
	}
	return 0, sjyb.Config.NotFoundOrInvalidNotation
}

// GetNumberByPath search an integer or float64 typed json/yaml node specified by the path and
// returns the value as float64 and the type as string
func (sjyb SmartJsonYamlBase) GetNumberByPath(path string) (float64, string) {
	val, typ := sjyb.GetNodeByPath(path)
	if i, isInt := val.(int); typ == "int" && isInt {
		return float64(i), typ
	}
	if f, isFlt := val.(float64); typ == "float64" && isFlt {
		return f, typ
	}
	return 0, sjyb.Config.NotFoundOrInvalidNotation
}

// GetTimeByPath search a date or time typed yaml node specified by the path and
// returns the value and the type as string
func (sjyb SmartJsonYamlBase) GetTimeByPath(path string) (time.Time, string) {
	val, typ := sjyb.GetNodeByPath(path)
	if tv, isTime := val.(time.Time); typ == "time.Time" && isTime {
		return tv, typ
	}
	return time.Time{}, sjyb.Config.NotFoundOrInvalidNotation
}

// GetBoolByPath search a boolean typed json/yaml node specified by the path and
// returns the value and the type as string
func (sjyb SmartJsonYamlBase) GetBoolByPath(path string) (bool, string) {
	val, typ := sjyb.GetNodeByPath(path)
	if b, isBool := val.(bool); typ == "bool" && isBool {
		return b, typ
	}
	return false, sjyb.Config.NotFoundOrInvalidNotation
}

// GetStringByPathWithDefault search a string typed json/yaml node specified by the path
// and returns the value if found, otherwise the value of def is returned without error
func (sjyb SmartJsonYamlBase) GetStringByPathWithDefault(path string, def string) string {
	val, typ := sjyb.GetNodeByPath(path)
	if str, isStr := val.(string); typ == "string" && isStr {
		return str
	}
	return def
}

// GetFloat64ByPathWithDefault search a float64 typed json/yaml node specified by the path
// and returns the value if found, otherwise the value of def is returned without error
func (sjyb SmartJsonYamlBase) GetFloat64ByPathWithDefault(path string, def float64) float64 {
	val, typ := sjyb.GetNodeByPath(path)
	if f, isFlt := val.(float64); typ == "float64" && isFlt {
		return f
	}
	return def
}

// GetIntegerByPathWithDefault search a int typed json/yaml node specified by the path
// and returns the value if found, otherwise the value of def is returned without error
func (sjyb SmartJsonYamlBase) GetIntegerByPathWithDefault(path string, def int) int {
	val, typ := sjyb.GetNodeByPath(path)
	if i, isInt := val.(int); typ == "int" && isInt {
		return i
	}
	return def
}

// GetNumberByPathWithDefault search a float64 or int typed json/yaml node specified by the path
// and returns the value if found, otherwise the value of def is returned without error
func (sjyb SmartJsonYamlBase) GetNumberByPathWithDefault(path string, def float64) float64 {
	val, typ := sjyb.GetNodeByPath(path)
	if i, isInt := val.(int); typ == "int" && isInt {
		return float64(i)
	}
	if f, isFlt := val.(float64); typ == "float64" && isFlt {
		return f
	}
	return def
}

// GetTimeByPathWithDefault search a date or time typed yaml node specified by the path
// and returns the value if found, otherwise the value of def is returned without error
func (sjyb SmartJsonYamlBase) GetTimeByPathWithDefault(path string, def time.Time) time.Time {
	val, typ := sjyb.GetNodeByPath(path)
	if tv, isTime := val.(time.Time); typ == "time" && isTime {
		return tv
	}
	return def
}

// GetBoolByPathWithDefault search a boolean typed json/yaml node specified by the path
// and returns the value if found, otherwise the value of def is returned without error
func (sjyb SmartJsonYamlBase) GetBoolByPathWithDefault(path string, def bool) bool {
	val, typ := sjyb.GetNodeByPath(path)
	if b, isBool := val.(bool); typ == "bool" && isBool {
		return b
	}
	return def
}

// GetCountDescendantsByPath search a json/yaml node specified by the path
// and returns the number of descendant nodes. For a non existing node it returns zero.
func (sjyb SmartJsonYamlBase) GetCountDescendantsByPath(path string) int {
	val, typ := sjyb.GetNodeByPath(path)
	if a, isArray := val.([]interface{}); typ == "array" && isArray {
		return len(a)
	}
	if m, isMap := val.(map[string]interface{}); typ == "map" && isMap {
		return len(m)
	}
	return 0
}

func contains(elems []string, v string) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}
