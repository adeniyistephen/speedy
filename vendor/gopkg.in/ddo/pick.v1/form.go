package pick

import (
	"io"
	"strings"
)

// input can be input/select/textarea
// TODO: now only get 1st value, need []string as result
// TODO: pick select option that has no value attr
// TODO: radio
func PickForm(reader io.Reader, attr *Attr) (input map[string][]string) {
	htmlArr := PickHtml(&Option{
		PageSource: reader,
		TagName:    "form",
		Attr:       attr,
	}, 1)
	if len(htmlArr) == 0 {
		return
	}

	html := htmlArr[0]
	input = make(map[string][]string)

	// pick input
	inputNameArr := PickAttr(&Option{
		PageSource: strings.NewReader(html),
		TagName:    "input",
		Attr:       nil,
	}, "name", 0)

	for i := 0; i < len(inputNameArr); i++ {
		values := PickAttr(&Option{
			PageSource: strings.NewReader(html),
			TagName:    "input",
			Attr: &Attr{
				Label: "name",
				Value: inputNameArr[i],
			},
		}, "value", 1)
		if len(values) == 0 {
			input[inputNameArr[i]] = []string{""}
			continue
		}

		input[inputNameArr[i]] = []string{values[0]}
	}

	// pick textarea
	textNameArr := PickAttr(&Option{
		PageSource: strings.NewReader(html),
		TagName:    "textarea",
		Attr:       nil,
	}, "name", 0)

	for i := 0; i < len(textNameArr); i++ {
		values := PickText(&Option{
			PageSource: strings.NewReader(html),
			TagName:    "textarea",
			Attr: &Attr{
				Label: "name",
				Value: textNameArr[i],
			},
		}, 1)
		if len(values) == 0 {
			input[textNameArr[i]] = []string{""}
			continue
		}

		input[textNameArr[i]] = []string{values[0]}
	}

	// pick select
	selectNameArr := PickAttr(&Option{
		PageSource: strings.NewReader(html),
		TagName:    "select",
		Attr:       nil,
	}, "name", 0)

	for i := 0; i < len(selectNameArr); i++ {
		selectHtmlArr := PickHtml(&Option{
			PageSource: strings.NewReader(html),
			TagName:    "select",
			Attr: &Attr{
				Label: "name",
				Value: selectNameArr[i],
			},
		}, 1)
		if len(selectHtmlArr) == 0 {
			continue
		}

		// get selected option
		values := PickAttr(&Option{
			PageSource: strings.NewReader(selectHtmlArr[0]),
			TagName:    "option",
			Attr: &Attr{
				Label: "selected",
				Value: "",
			},
		}, "value", 1)
		if len(values) != 0 {
			input[selectNameArr[i]] = []string{values[0]}
			continue
		}

		values = PickAttr(&Option{
			PageSource: strings.NewReader(selectHtmlArr[0]),
			TagName:    "option",
			Attr: &Attr{
				Label: "selected",
				Value: "selected",
			},
		}, "value", 1)
		if len(values) != 0 {
			input[selectNameArr[i]] = []string{values[0]}
			continue
		}

		// if no select get the 1st option value
		values = PickAttr(&Option{
			PageSource: strings.NewReader(selectHtmlArr[0]),
			TagName:    "option",
			Attr:       nil,
		}, "value", 1)
		if len(values) == 0 {
			input[selectNameArr[i]] = []string{""}
			continue
		}

		input[selectNameArr[i]] = []string{values[0]}
	}

	return
}
