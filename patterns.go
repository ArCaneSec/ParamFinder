package main

import "regexp"

var (
	JsObjKeysPattern   = regexp.MustCompile(`(?:[{,]\s*(?:['"])?)([\w_-]+?)(?:\s*)(?:['"]?:)`) // JS object keys
	htmlNamePattern    = regexp.MustCompile(`(?:<input.*?name)(?:="|')(.*?)(?:'|")`)           // html name keys
	htmlIdPattern      = regexp.MustCompile(`(?:<input.*?id)(?:="|')(.*?)(?:'|")`)             // html id keys
	JsVariablePattern  = regexp.MustCompile(`(?:(?:let|const|var)\s*)(\w+)`)                   // JS variable names
	HttpPostParams     = regexp.MustCompile(`(?:&|^)([a-zA-Z+\d%]*?)=`)                        // http post parameters
	HttpUrlParams      = regexp.MustCompile(`[?&](.*?)(?:=)`)                                  // http url parameters
	HttpMultipartNames = regexp.MustCompile(` name="(.*?)"`)                                   // multipart data name keys
	jsFiles            = regexp.MustCompile("(?:[\"'`])(.+\\.js(on)?)(?:[\"'`])")
)

var patterns = []*regexp.Regexp{htmlNamePattern, htmlIdPattern, JsVariablePattern, JsObjKeysPattern, HttpUrlParams, HttpPostParams, HttpMultipartNames}
