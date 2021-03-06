package json

import "github.com/gopherjs/gopherjs/js"

func getJSONTag(tag *js.Object) string {
	strLen := tag.Length()
	var val *js.Object
	for i := 0; i < strLen; i++ {
		start := i
		for i < strLen && tag.Call("charCodeAt", i).Int() == ' ' {
			i++
		}
		var char int
		for ; i < strLen; i++ {
			char = tag.Call("charCodeAt", i).Int()
			if char <= ' ' || char == ':' || char == '"' {
				break
			}
		}
		if i+1 >= strLen || char != ':' || tag.Call("charCodeAt", i+1).Int() != '"' {
			break
		}
		name := tag.Call("substring", start, i)
		i += 2
		start = i
		for ; i < strLen; i++ {
			char = tag.Call("charCodeAt", i).Int()
			if char == '\\' {
				i++
			} else if char == '"' {
				break
			}
		}
		if name.String() == "json" {
			val = tag.Call("substring", start, i)
			break
		}
	}
	if val == nil {
		return ""
	}
	return val.String()
}

type tagOptions struct {
	*js.Object
}

func parseTag(tag string) (string, tagOptions) {
	split := js.Global.Get("String").Get("prototype").Get("indexOf").Call("call", tag, ",").Int()
	if split < 0 {
		return tag, tagOptions{nil}
	}
	return js.Global.Get("String").Get("prototype").Get("substr").Call("call", tag, 0, split).String(), tagOptions{js.Global.Get("String").Get("prototype").Get("substr").Call("call", tag, split+1)}
}

func (o tagOptions) Contains(option string) bool {
	if o.Object == nil || o.Length() == 0 || o.Length() < len(option) {
		return false
	}
	if o.Length() == len(option) {
		return o.String() == option
	}
	return o.Call("substr", 0, len(option)+1).String() == option+"," || o.Call("substr", -1-len(option)).String() == ","+option || o.Call("includes", ","+option+",").Bool()
}
