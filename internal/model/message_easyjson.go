// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package model

import (
	json "encoding/json"
	uuid "github.com/google/uuid"
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

func easyjson4086215fDecodeGithubComGoParkMailRu20251VelvetPullsInternalModel(in *jlexer.Lexer, out *SendMessage) {
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
		case "messageType":
			out.MessageType = MsgType(in.String())
		case "payload":
			if m, ok := out.Payload.(easyjson.Unmarshaler); ok {
				m.UnmarshalEasyJSON(in)
			} else if m, ok := out.Payload.(json.Unmarshaler); ok {
				_ = m.UnmarshalJSON(in.Raw())
			} else {
				out.Payload = in.Interface()
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
func easyjson4086215fEncodeGithubComGoParkMailRu20251VelvetPullsInternalModel(out *jwriter.Writer, in SendMessage) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"messageType\":"
		out.RawString(prefix[1:])
		out.String(string(in.MessageType))
	}
	{
		const prefix string = ",\"payload\":"
		out.RawString(prefix)
		if m, ok := in.Payload.(easyjson.Marshaler); ok {
			m.MarshalEasyJSON(out)
		} else if m, ok := in.Payload.(json.Marshaler); ok {
			out.Raw(m.MarshalJSON())
		} else {
			out.Raw(json.Marshal(in.Payload))
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v SendMessage) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson4086215fEncodeGithubComGoParkMailRu20251VelvetPullsInternalModel(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v SendMessage) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson4086215fEncodeGithubComGoParkMailRu20251VelvetPullsInternalModel(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *SendMessage) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson4086215fDecodeGithubComGoParkMailRu20251VelvetPullsInternalModel(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *SendMessage) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson4086215fDecodeGithubComGoParkMailRu20251VelvetPullsInternalModel(l, v)
}
func easyjson4086215fDecodeGithubComGoParkMailRu20251VelvetPullsInternalModel1(in *jlexer.Lexer, out *MessageList) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
		*out = nil
	} else {
		in.Delim('[')
		if *out == nil {
			if !in.IsDelim(']') {
				*out = make(MessageList, 0, 0)
			} else {
				*out = MessageList{}
			}
		} else {
			*out = (*out)[:0]
		}
		for !in.IsDelim(']') {
			var v1 Message
			(v1).UnmarshalEasyJSON(in)
			*out = append(*out, v1)
			in.WantComma()
		}
		in.Delim(']')
	}
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson4086215fEncodeGithubComGoParkMailRu20251VelvetPullsInternalModel1(out *jwriter.Writer, in MessageList) {
	if in == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
		out.RawString("null")
	} else {
		out.RawByte('[')
		for v2, v3 := range in {
			if v2 > 0 {
				out.RawByte(',')
			}
			(v3).MarshalEasyJSON(out)
		}
		out.RawByte(']')
	}
}

// MarshalJSON supports json.Marshaler interface
func (v MessageList) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson4086215fEncodeGithubComGoParkMailRu20251VelvetPullsInternalModel1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v MessageList) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson4086215fEncodeGithubComGoParkMailRu20251VelvetPullsInternalModel1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *MessageList) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson4086215fDecodeGithubComGoParkMailRu20251VelvetPullsInternalModel1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *MessageList) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson4086215fDecodeGithubComGoParkMailRu20251VelvetPullsInternalModel1(l, v)
}
func easyjson4086215fDecodeGithubComGoParkMailRu20251VelvetPullsInternalModel2(in *jlexer.Lexer, out *MessageInput) {
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
		case "message":
			out.Message = string(in.String())
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
func easyjson4086215fEncodeGithubComGoParkMailRu20251VelvetPullsInternalModel2(out *jwriter.Writer, in MessageInput) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"message\":"
		out.RawString(prefix[1:])
		out.String(string(in.Message))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v MessageInput) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson4086215fEncodeGithubComGoParkMailRu20251VelvetPullsInternalModel2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v MessageInput) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson4086215fEncodeGithubComGoParkMailRu20251VelvetPullsInternalModel2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *MessageInput) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson4086215fDecodeGithubComGoParkMailRu20251VelvetPullsInternalModel2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *MessageInput) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson4086215fDecodeGithubComGoParkMailRu20251VelvetPullsInternalModel2(l, v)
}
func easyjson4086215fDecodeGithubComGoParkMailRu20251VelvetPullsInternalModel3(in *jlexer.Lexer, out *Message) {
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
			if data := in.UnsafeBytes(); in.Ok() {
				in.AddError((out.ID).UnmarshalText(data))
			}
		case "parent_message_id":
			if in.IsNull() {
				in.Skip()
				out.ParentMessageID = nil
			} else {
				if out.ParentMessageID == nil {
					out.ParentMessageID = new(uuid.UUID)
				}
				if data := in.UnsafeBytes(); in.Ok() {
					in.AddError((*out.ParentMessageID).UnmarshalText(data))
				}
			}
		case "chat_id":
			if data := in.UnsafeBytes(); in.Ok() {
				in.AddError((out.ChatID).UnmarshalText(data))
			}
		case "user_id":
			if data := in.UnsafeBytes(); in.Ok() {
				in.AddError((out.UserID).UnmarshalText(data))
			}
		case "body":
			out.Body = string(in.String())
		case "sent_at":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.SentAt).UnmarshalJSON(data))
			}
		case "is_redacted":
			out.IsRedacted = bool(in.Bool())
		case "avatar_path":
			if in.IsNull() {
				in.Skip()
				out.AvatarPath = nil
			} else {
				if out.AvatarPath == nil {
					out.AvatarPath = new(string)
				}
				*out.AvatarPath = string(in.String())
			}
		case "user":
			out.Username = string(in.String())
		case "message_type":
			out.MessageType = string(in.String())
		case "files":
			if in.IsNull() {
				in.Skip()
				out.FilesDTO = nil
			} else {
				in.Delim('[')
				if out.FilesDTO == nil {
					if !in.IsDelim(']') {
						out.FilesDTO = make([]Payload, 0, 1)
					} else {
						out.FilesDTO = []Payload{}
					}
				} else {
					out.FilesDTO = (out.FilesDTO)[:0]
				}
				for !in.IsDelim(']') {
					var v4 Payload
					(v4).UnmarshalEasyJSON(in)
					out.FilesDTO = append(out.FilesDTO, v4)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "photos":
			if in.IsNull() {
				in.Skip()
				out.PhotosDTO = nil
			} else {
				in.Delim('[')
				if out.PhotosDTO == nil {
					if !in.IsDelim(']') {
						out.PhotosDTO = make([]Payload, 0, 1)
					} else {
						out.PhotosDTO = []Payload{}
					}
				} else {
					out.PhotosDTO = (out.PhotosDTO)[:0]
				}
				for !in.IsDelim(']') {
					var v5 Payload
					(v5).UnmarshalEasyJSON(in)
					out.PhotosDTO = append(out.PhotosDTO, v5)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "sticker":
			out.Sticker = string(in.String())
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
func easyjson4086215fEncodeGithubComGoParkMailRu20251VelvetPullsInternalModel3(out *jwriter.Writer, in Message) {
	out.RawByte('{')
	first := true
	_ = first
	if true {
		const prefix string = ",\"id\":"
		first = false
		out.RawString(prefix[1:])
		out.RawText((in.ID).MarshalText())
	}
	if in.ParentMessageID != nil {
		const prefix string = ",\"parent_message_id\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.RawText((*in.ParentMessageID).MarshalText())
	}
	if true {
		const prefix string = ",\"chat_id\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.RawText((in.ChatID).MarshalText())
	}
	if true {
		const prefix string = ",\"user_id\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.RawText((in.UserID).MarshalText())
	}
	if in.Body != "" {
		const prefix string = ",\"body\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Body))
	}
	if true {
		const prefix string = ",\"sent_at\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Raw((in.SentAt).MarshalJSON())
	}
	if in.IsRedacted {
		const prefix string = ",\"is_redacted\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Bool(bool(in.IsRedacted))
	}
	if in.AvatarPath != nil {
		const prefix string = ",\"avatar_path\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(*in.AvatarPath))
	}
	if in.Username != "" {
		const prefix string = ",\"user\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Username))
	}
	{
		const prefix string = ",\"message_type\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.MessageType))
	}
	if len(in.FilesDTO) != 0 {
		const prefix string = ",\"files\":"
		out.RawString(prefix)
		{
			out.RawByte('[')
			for v6, v7 := range in.FilesDTO {
				if v6 > 0 {
					out.RawByte(',')
				}
				(v7).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	if len(in.PhotosDTO) != 0 {
		const prefix string = ",\"photos\":"
		out.RawString(prefix)
		{
			out.RawByte('[')
			for v8, v9 := range in.PhotosDTO {
				if v8 > 0 {
					out.RawByte(',')
				}
				(v9).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"sticker\":"
		out.RawString(prefix)
		out.String(string(in.Sticker))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Message) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson4086215fEncodeGithubComGoParkMailRu20251VelvetPullsInternalModel3(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Message) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson4086215fEncodeGithubComGoParkMailRu20251VelvetPullsInternalModel3(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Message) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson4086215fDecodeGithubComGoParkMailRu20251VelvetPullsInternalModel3(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Message) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson4086215fDecodeGithubComGoParkMailRu20251VelvetPullsInternalModel3(l, v)
}
func easyjson4086215fDecodeGithubComGoParkMailRu20251VelvetPullsInternalModel4(in *jlexer.Lexer, out *LastMessage) {
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
			if data := in.UnsafeBytes(); in.Ok() {
				in.AddError((out.ID).UnmarshalText(data))
			}
		case "user_id":
			if data := in.UnsafeBytes(); in.Ok() {
				in.AddError((out.UserID).UnmarshalText(data))
			}
		case "body":
			out.Body = string(in.String())
		case "sent_at":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.SentAt).UnmarshalJSON(data))
			}
		case "user":
			out.Username = string(in.String())
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
func easyjson4086215fEncodeGithubComGoParkMailRu20251VelvetPullsInternalModel4(out *jwriter.Writer, in LastMessage) {
	out.RawByte('{')
	first := true
	_ = first
	if true {
		const prefix string = ",\"id\":"
		first = false
		out.RawString(prefix[1:])
		out.RawText((in.ID).MarshalText())
	}
	if true {
		const prefix string = ",\"user_id\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.RawText((in.UserID).MarshalText())
	}
	if in.Body != "" {
		const prefix string = ",\"body\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Body))
	}
	if true {
		const prefix string = ",\"sent_at\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Raw((in.SentAt).MarshalJSON())
	}
	if in.Username != "" {
		const prefix string = ",\"user\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Username))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v LastMessage) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson4086215fEncodeGithubComGoParkMailRu20251VelvetPullsInternalModel4(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v LastMessage) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson4086215fEncodeGithubComGoParkMailRu20251VelvetPullsInternalModel4(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *LastMessage) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson4086215fDecodeGithubComGoParkMailRu20251VelvetPullsInternalModel4(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *LastMessage) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson4086215fDecodeGithubComGoParkMailRu20251VelvetPullsInternalModel4(l, v)
}
