package pick

import (
	"io"

	"golang.org/x/net/html"
)

// Option
// PageSource: html reader
// TagName: element's tag. ex: div, p, form, input...
type Option struct {
	PageSource io.Reader
	TagName    string
	Attr       *Attr // optional
}

// Attr to define an element's attribute
// <input type="hidden" />
// -> Attr{Label: "type", Value: "hidden"}
type Attr struct {
	Label string
	Value string
}

// PickAttr gets attribute Option.Attr of elements that matches the Option.TagName
func PickAttr(option *Option, AttrLabel string, limit int) (res []string) {
	if option == nil || option.PageSource == nil {
		return
	}

	z := html.NewTokenizer(option.PageSource)

	for {
		tokenType := z.Next()

		switch tokenType {

		// ignore the error token
		// quit on eof
		case html.ErrorToken:
			if z.Err() == io.EOF {
				return
			}

		case html.StartTagToken, html.SelfClosingTagToken:
			tagName, attr := z.TagName()

			if string(tagName) != option.TagName {
				continue
			}

			var label, value []byte

			matched := false
			tmpRes := []string{}

			// get attr
			for attr {
				label, value, attr = z.TagAttr()

				labelStr := string(label)
				valueStr := string(value)

				// check the attr
				if option.Attr == nil || (option.Attr.Label == labelStr && option.Attr.Value == valueStr) {
					matched = true
				}

				// get the result - even the matched false or true
				if labelStr == AttrLabel {
					tmpRes = append(tmpRes, valueStr)
				}
			}

			// skip the non matched one
			if !matched {
				continue
			}

			// send the result for matched only
			res = append(res, tmpRes...)

			// return when limit
			if limit > 0 && len(res) >= limit {
				return
			}
		}
	}

	return
}

// PickText gets text of elements that matches the Option.TagName (including it's child elements)
func PickText(option *Option, limit int) (res []string) {
	if option == nil || option.PageSource == nil {
		return
	}

	z := html.NewTokenizer(option.PageSource)

	depth := 0

	for {
		tokenType := z.Next()

		switch tokenType {

		// ignore the error token
		// quit on eof
		case html.ErrorToken:
			if z.Err() == io.EOF {
				return
			}

		// get text
		case html.TextToken:
			if depth > 0 {
				// append to the last element
				res[len(res)-1] = res[len(res)-1] + string(z.Text())
			}

		case html.EndTagToken:
			if depth > 0 {
				depth--
			}

		case html.StartTagToken:
			tagName, attr := z.TagName()

			// inside the target
			if depth > 0 && !isSelfClosingTag(tagName) {
				depth++
				continue
			}

			// check limit
			if limit > 0 && len(res) >= limit {
				return
			}

			if string(tagName) != option.TagName {
				continue
			}

			var label, value []byte

			matched := false

			// empty attr element
			if !attr && option.Attr == nil {
				matched = true
			}

			// get attr
			for attr {
				label, value, attr = z.TagAttr()

				if option.Attr == nil || (option.Attr.Label == string(label) && option.Attr.Value == string(value)) {
					matched = true
				}
			}

			if !matched {
				continue
			}

			depth++

			// init an empty element
			res = append(res, "")
		}
	}

	return
}

// isSelfClosingTag returns if it's a self closing tag
// html.SelfClosingTagToken is not always correct
// <input/> -> html.SelfClosingTagToken
// <input> -> html.StartTagToken
func isSelfClosingTag(tag []byte) bool {
	switch string(tag) {
	case "area", "base", "br", "col", "embed", "hr", "img", "input", "keygen", "link", "meta", "param", "source", "track", "wbr":
		return true
	}

	return false
}

// PickHtml gets html of elements that matches the Option.TagName (including it's child elements)
func PickHtml(option *Option, limit int) (res []string) {
	if option == nil || option.PageSource == nil {
		return
	}

	z := html.NewTokenizer(option.PageSource)

	depth := 0

	for {
		tokenType := z.Next()

		switch tokenType {

		// ignore the error token
		// quit on eof
		case html.ErrorToken:
			if z.Err() == io.EOF {
				return
			}

		// get text
		case html.TextToken:
			if depth > 0 {
				// append to the last element
				res[len(res)-1] = res[len(res)-1] + string(z.Raw())
			}

		case html.EndTagToken:
			if depth > 0 {
				depth--
			}

			if depth > 0 {
				// append to the last element
				res[len(res)-1] = res[len(res)-1] + string(z.Raw())
			}

		case html.SelfClosingTagToken:
			if depth > 0 {
				// append to the last element
				res[len(res)-1] = res[len(res)-1] + string(z.Raw())
			}

		case html.StartTagToken:
			tagName, attr := z.TagName()

			// inside the target
			if depth > 0 {
				// append to the last element
				res[len(res)-1] = res[len(res)-1] + string(z.Raw())

				if !isSelfClosingTag(tagName) {
					depth++
				}
				continue
			}

			// check limit
			if limit > 0 && len(res) >= limit {
				return
			}

			if string(tagName) != option.TagName {
				continue
			}

			var label, value []byte

			matched := false

			// empty attr element
			if !attr && option.Attr == nil {
				matched = true
			}

			// get attr
			for attr {
				label, value, attr = z.TagAttr()

				if option.Attr == nil || (option.Attr.Label == string(label) && option.Attr.Value == string(value)) {

					matched = true
				}
			}

			if !matched {
				continue
			}

			depth++

			// init an empty element
			res = append(res, "")
		}
	}

	return
}
