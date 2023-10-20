// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package models

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjsonA0cbd601DecodeGitlabComBotsGoAvostBotPkgAnimevostModels(in *jlexer.Lexer, out *AnimeSpec) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "id":
			out.Id = int(in.Int())
		case "rating":
			out.Rating = int(in.Int())
		case "votes":
			out.Votes = int(in.Int())
		case "title":
			out.Title = string(in.String())
		case "description":
			out.Description = string(in.String())
		case "year":
			out.Year = string(in.String())
		case "urlImagePreview":
			out.UrlImagePreview = string(in.String())
		case "type":
			out.Type = string(in.String())
		case "director":
			out.Director = string(in.String())
		case "series":
			out.Series = string(in.String())
		case "screenImage":
			if in.IsNull() {
				in.Skip()
				out.ScreenImage = nil
			} else {
				in.Delim('[')
				if out.ScreenImage == nil {
					if !in.IsDelim(']') {
						out.ScreenImage = make([]string, 0, 4)
					} else {
						out.ScreenImage = []string{}
					}
				} else {
					out.ScreenImage = (out.ScreenImage)[:0]
				}
				for !in.IsDelim(']') {
					var v1 string
					v1 = string(in.String())
					out.ScreenImage = append(out.ScreenImage, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "timer":
			if m, ok := out.Timer.(easyjson.Unmarshaler); ok {
				m.UnmarshalEasyJSON(in)
			} else if m, ok := out.Timer.(json.Unmarshaler); ok {
				_ = m.UnmarshalJSON(in.Raw())
			} else {
				out.Timer = in.Interface()
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonA0cbd601EncodeGitlabComBotsGoAvostBotPkgAnimevostModels(out *jwriter.Writer, in AnimeSpec) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.Int(int(in.Id))
	}
	{
		const prefix string = ",\"rating\":"
		out.RawString(prefix)
		out.Int(int(in.Rating))
	}
	{
		const prefix string = ",\"votes\":"
		out.RawString(prefix)
		out.Int(int(in.Votes))
	}
	{
		const prefix string = ",\"title\":"
		out.RawString(prefix)
		out.String(string(in.Title))
	}
	{
		const prefix string = ",\"description\":"
		out.RawString(prefix)
		out.String(string(in.Description))
	}
	{
		const prefix string = ",\"year\":"
		out.RawString(prefix)
		out.String(string(in.Year))
	}
	{
		const prefix string = ",\"urlImagePreview\":"
		out.RawString(prefix)
		out.String(string(in.UrlImagePreview))
	}
	{
		const prefix string = ",\"type\":"
		out.RawString(prefix)
		out.String(string(in.Type))
	}
	{
		const prefix string = ",\"director\":"
		out.RawString(prefix)
		out.String(string(in.Director))
	}
	{
		const prefix string = ",\"series\":"
		out.RawString(prefix)
		out.String(string(in.Series))
	}
	{
		const prefix string = ",\"screenImage\":"
		out.RawString(prefix)
		if in.ScreenImage == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v2, v3 := range in.ScreenImage {
				if v2 > 0 {
					out.RawByte(',')
				}
				out.String(string(v3))
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"timer\":"
		out.RawString(prefix)
		if m, ok := in.Timer.(easyjson.Marshaler); ok {
			m.MarshalEasyJSON(out)
		} else if m, ok := in.Timer.(json.Marshaler); ok {
			out.Raw(m.MarshalJSON())
		} else {
			out.Raw(json.Marshal(in.Timer))
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (a AnimeSpec) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonA0cbd601EncodeGitlabComBotsGoAvostBotPkgAnimevostModels(&w, a)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (a AnimeSpec) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonA0cbd601EncodeGitlabComBotsGoAvostBotPkgAnimevostModels(w, a)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (a *AnimeSpec) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonA0cbd601DecodeGitlabComBotsGoAvostBotPkgAnimevostModels(&r, a)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (a *AnimeSpec) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonA0cbd601DecodeGitlabComBotsGoAvostBotPkgAnimevostModels(l, a)
}