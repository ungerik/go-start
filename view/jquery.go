package view

import (
	"fmt"
	"io"
)

func JQuery(request *Request, session *Session, response *Response, writer io.Writer) (err error) {
	writer.Write([]byte(`<script src="//ajax.googleapis.com/ajax/libs/jquery/1.7.1/jquery.min.js"></script>`))
	writer.Write([]byte(`<script>window.jQuery || document.write('<script src="/js/libs/jquery-1.7.1.min.js"><\/script>')</script>`))
	writer.Write([]byte{'\n'})
	return nil
}

func JQueryUI(request *Request, session *Session, response *Response, writer io.Writer) (err error) {
	writer.Write([]byte(`<script src="//ajax.googleapis.com/ajax/libs/jqueryui/1.8.17/jquery-ui.min.js"></script>`))
	writer.Write([]byte(`<script>window.jQuery.ui || document.write('<script src="/js/libs/jquery-ui-1.8.17.custom.min.js"><\/script>')</script>`))
	writer.Write([]byte{'\n'})
	return nil
}

func JQueryUIAutocompleteFromURL(domSelector string, dataURL URL, minLength int) PageWriteFunc {
	return func(request *Request, session *Session, response *Response, writer io.Writer) (err error) {
		url := dataURL.URL(request, session, response)
		fmt.Fprintf(writer, "<script>$('%s').autocomplete({source:'%s',minLength:%d});</script>", domSelector, url, minLength)
		return nil
	}
}

func JQueryUIAutocomplete(domSelector string, options []string, minLength int) PageWriteFunc {
	return func(request *Request, session *Session, response *Response, writer io.Writer) (err error) {
		fmt.Fprintf(writer, "<script>$('%s').autocomplete({source:[", domSelector)
		for i := range options {
			if i > 0 {
				writer.Write([]byte{','})
			}
			writer.Write([]byte{'"'})
			writer.Write([]byte(options[i]))
			writer.Write([]byte{'"'})
		}
		fmt.Fprintf(writer, "],minLength:%d});</script>", minLength)
		return nil
	}
}
