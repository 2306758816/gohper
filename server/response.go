package server

import (
	"net/http"

	"github.com/cosiner/golib/encoding"
)

type (
	Response interface {
		SetHeader(name, value string)
		AddHeader(name, value string)
		SetContentType(typ string)
		SetCookie(name, value string)
		SetCookieWithExpire(name, value string, lifetime int)
		DeleteClientCookie(name string)
		Redirect(url string)
		PermanentRedirect(url string)
		ReportError(statusCode int)
		Render(tmpl string) error
		WriteString(data string) (int, error)
		WriteJSON(val interface{}) error
		WriteXML(val interface{}) error
		Flush()
	}
	// response represent a response of request to user
	response struct {
		*context
		w      http.ResponseWriter
		header http.Header
	}

	// marshalFunc is the marshal function type
	marshalFunc func(interface{}) ([]byte, error)
)

// newResponse create a new response, and set default content type to HTML
func newResponse(ctx *context, w http.ResponseWriter) *response {
	resp := &response{
		context: ctx,
		w:       w,
		header:  w.Header(),
	}
	resp.SetContentType(CONTENTTYPE_HTML)
	return resp
}

// destroy destroy all reference that response keep
func (resp *response) destroy() {
	resp.context.destroy()
	resp.w = nil
	resp.header = nil
}

// SetHeader setup response header
func (resp *response) SetHeader(name, value string) {
	resp.header.Set(name, value)
}

// AddHeader add a value to response header
func (resp *response) AddHeader(name, value string) {
	resp.header.Add(name, value)
}

// SetContentType set content type of response
func (resp *response) SetContentType(typ string) {
	resp.SetHeader(HEADER_CONTENTTYPE, typ)
}

// contentType return current content type of response
func (resp *response) contentType() string {
	return resp.header.Get(HEADER_CONTENTTYPE)
}

// newCookie create a new Cookie and return it's displayed string
// parameter lifetime is time by second
func (*response) newCookie(name, value string, lifetime int) string {
	return (&http.Cookie{
		Name:   name,
		Value:  value,
		MaxAge: lifetime,
	}).String()
}

// SetCookie setup response cookie, default age is default browser opened time
func (resp *response) SetCookie(name, value string) {
	resp.SetCookieWithExpire(name, value, 0)
}

// SetCookieWithExpire setup response cookie with lifetime
func (resp *response) SetCookieWithExpire(name, value string, lifetime int) {
	resp.SetHeader(HEADER_SETCOOKIE, resp.newCookie(name, value, lifetime))
}

// DeleteClientCookie delete user briwser's cookie by name
func (resp *response) DeleteClientCookie(name string) {
	resp.SetCookieWithExpire(name, "", -1)
}

// setSessionCookie setup session cookie
func (resp *response) setSessionCookie(id string) {
	resp.SetCookie(_COOKIE_SESSION, id)
}

// Redirect redirect to new url
func (resp *response) Redirect(url string) {
	http.Redirect(resp.w, resp.request, url, http.StatusTemporaryRedirect)
}

// PermanentRedirect permanently redirect current request url to new url
func (resp *response) PermanentRedirect(url string) {
	http.Redirect(resp.w, resp.request, url, http.StatusMovedPermanently)
}

// Report Error report an http error with given status code
func (resp *response) ReportError(statusCode int) {
	resp.w.WriteHeader(statusCode)
}

// Render render template with context
func (resp *response) Render(tmpl string) error {
	return resp.Server().renderTemplate(resp, tmpl, resp.context)
}

func (resp *response) Write(data []byte) (int, error) {
	return resp.w.Write(data)
}

// WriteString write sting to client
func (resp *response) WriteString(data string) (int, error) {
	return encoding.WriteString(resp, data)
}

// WriteJSON write json data to client, and setup content type to json
func (resp *response) WriteJSON(val interface{}) error {
	return encoding.WriteJSON(resp, val)
}

// WriteXML write xml data to client, and setup content type to xml
func (resp *response) WriteXML(val interface{}) error {
	return encoding.WriteXML(resp, val)
}

// Flush flush response's output
func (resp *response) Flush() {
	if flusher, is := resp.w.(http.Flusher); is {
		flusher.Flush()
	}
}