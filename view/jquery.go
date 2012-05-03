package view

var (
	JQuery   HTML = `<script src="//ajax.googleapis.com/ajax/libs/jquery/1.7.1/jquery.min.js"></script><script>window.jQuery || document.write('<script src="/js/libs/jquery-1.7.1.min.js"><\/script>')</script>`
	JQueryUI HTML = `<script src="//ajax.googleapis.com/ajax/libs/jqueryui/1.8.17/jquery-ui.min.js"></script><script>window.jQuery.ui || document.write('<script src="/js/libs/jquery-ui-1.8.17.custom.min.js"><\/script>')</script>`
)

func JQueryUIAutocompleteFromURL(domSelector string, dataURL URL, minLength int) View {
	return RenderView(
		func(request *Request, session *Session, response *Response) (err error) {
			url := dataURL.URL(request, session, response)
			response.Printf("<script>$('%s').autocomplete({source:'%s',minLength:%d});</script>", domSelector, url, minLength)
			return nil
		},
	)
}

func JQueryUIAutocomplete(domSelector string, options []string, minLength int) View {
	return RenderView(
		func(request *Request, session *Session, response *Response) (err error) {
			response.Printf("<script>$('%s').autocomplete({source:[", domSelector)
			for i := range options {
				if i > 0 {
					response.WriteByte(',')
				}
				response.WriteByte('"')
				response.WriteString(options[i])
				response.WriteByte('"')
			}
			response.Printf("],minLength:%d});</script>", minLength)
			return nil
		},
	)
}
