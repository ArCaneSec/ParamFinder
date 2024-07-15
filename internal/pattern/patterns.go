package pattern

import "regexp"

var (
	jsObjKeysPattern  = regexp.MustCompile(`(?:[{,]\s*(?:['"])?)([\w]+?)(?:\s*)(?:['"]?:)`) // JS object keys
	htmlNamePattern   = regexp.MustCompile(`(?:<input.*?name=)(?:"|')?(.*?)(?:'|"|\s)`)     // html name keys
	htmlLabelNames    = regexp.MustCompile(`(?:<label.*?for=["']?)(.*?)(?:[>'" ])`)         // multipart data name keys
	jsVariablePattern = regexp.MustCompile(`(?:^|[{(\s])(?:let|const|var\s*)(\w+)`)         // JS variable names
	htmlIdPattern     = regexp.MustCompile(`(?:<input.*?id=)(?:'|")?(.*?)(?:'|"|\s)`)       // html id keys
	JsFiles           = regexp.MustCompile("(?:[\"'`])(.+\\.js(on)?)(?:[\"'`])")            // js files path in respones
	httpUrlParams     = regexp.MustCompile(`(?:[?&]|%3F|%26|&amp;)([\w]+?)(?:=|%3D)`)       // http url parameters
	Patterns = []*regexp.Regexp{htmlNamePattern, htmlIdPattern, jsVariablePattern, jsObjKeysPattern, httpUrlParams, htmlLabelNames}
)

