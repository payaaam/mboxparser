package mboxparser

import (
	"bytes"
	"github.com/satori/go.uuid"
	"io"
	"mime/multipart"
	"mime/quotedprintable"
	"net/http"
	"net/mail"
	"net/textproto"
	"strings"
)

// Encode a mboxparser.Message to a writable mail.Message
func Encode(message *Message) (*mail.Message, error) {
	email := &mail.Message{
		Header: nil,
		Body:   nil,
	}
	messageBuffer := new(bytes.Buffer)
	mboxMessage := multipart.NewWriter(messageBuffer)

	// Add headers to *mail.Message from Message
	for key, header := range message.Header {
		email.Header[http.CanonicalHeaderKey(key)] = header
	}

	// TODO: Get Boundary from Headers
	emailBoundary := uuid.NewV4().String()
	mboxMessage.SetBoundary(emailBoundary)

	for _, part := range message.Bodies {

		partHeader := make(textproto.MIMEHeader)
		for key, header := range part.Header {
			partHeader.Set(http.CanonicalHeaderKey(key), strings.Join(header, "; "))
		}

		mboxPart, err := mboxMessage.CreatePart(partHeader)
		if err != nil {
			return nil, err
		}

		encodedMboxPart := quotedprintable.NewWriter(mboxPart)
		defer encodedMboxPart.Close()

		encodedMboxPart.Write(streamToBytes(part.Content))
	}

	err := mboxMessage.Close()
	if err != nil {
		return nil, err
	}
	email.Body = messageBuffer
	return email, nil
}

func streamToBytes(stream io.Reader) []byte {
	buf := new(bytes.Buffer)
	buf.ReadFrom(stream)
	return buf.Bytes()
}
