package main

import "regexp"

var (
	JsObjKeysPattern  = regexp.MustCompile(`(?:[{,]\s*(?:['"])?)([\w]+?)(?:\s*)(?:['"]?:)`) // JS object keys
	htmlNamePattern   = regexp.MustCompile(`(?:<input.*?name=)(?:"|')?(.*?)(?:'|"|\s)`)     // html name keys
	htmlLabelNames    = regexp.MustCompile(`(?:<label.*?for=["']?)(.*?)(?:[>'" ])`)         // multipart data name keys
	JsVariablePattern = regexp.MustCompile(`(?:^|[{(\s])(?:let|const|var\s*)(\w+)`)         // JS variable names
	htmlIdPattern     = regexp.MustCompile(`(?:<input.*?id=)(?:'|")?(.*?)(?:'|"|\s)`)       // html id keys
	jsFiles           = regexp.MustCompile("(?:[\"'`])(.+\\.js(on)?)(?:[\"'`])")            // js files path in respones
	HttpUrlParams     = regexp.MustCompile(`(?:[?&]|%3F|%26)([\w]+?)(?:=|%3D)`)             // http url parameters
)

var patterns = []*regexp.Regexp{htmlNamePattern, htmlIdPattern, JsVariablePattern, JsObjKeysPattern, HttpUrlParams, htmlLabelNames}
