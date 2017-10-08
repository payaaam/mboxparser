package mboxparser

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/satori/go.uuid"
	"io"
	"mime/multipart"
	"mime/quotedprintable"
	"net/http"
	"net/mail"
	"net/textproto"
	"strings"
)

var QuotedPrintable string = "quoted-printable"

// Encode a mboxparser.Message to a writable mail.Message
func Encode(message *Message) (*mail.Message, error) {
	email := &mail.Message{
		Header: mail.Header(map[string][]string{}),
		Body:   nil,
	}

	if isMultiPartMessage(message) != true {
		return processPlainEmail(email, message)
	}

	messageBuffer := new(bytes.Buffer)
	mboxMessage := multipart.NewWriter(messageBuffer)

	// TODO: Get Boundary from Headers
	emailBoundary := uuid.NewV4().String()
	mboxMessage.SetBoundary(emailBoundary)

	// Add headers to *mail.Message from Message
	for key, header := range message.Header {
		if key == "Content-Type" {
			header[0] = fmt.Sprintf("%s; boundary=\"%s\"", header[0], emailBoundary)
		}
		email.Header[http.CanonicalHeaderKey(key)] = header
	}

	for _, part := range message.Bodies {
		partHeader := make(textproto.MIMEHeader)

		for key, header := range part.Header {
			partHeader.Set(http.CanonicalHeaderKey(key), strings.Join(header, "; "))
		}

		ContentTransferEncodingHeader := http.CanonicalHeaderKey("Content-Transfer-Encoding")
		if partHeader.Get(ContentTransferEncodingHeader) == "" {
			partHeader.Set(ContentTransferEncodingHeader, QuotedPrintable)
		}

		mboxPart, err := mboxMessage.CreatePart(partHeader)
		if err != nil {
			return nil, err
		}

		encodedMboxPart := quotedprintable.NewWriter(mboxPart)

		encodedMboxPart.Write(streamToBytes(part.Content))
		err = encodedMboxPart.Close()
		if err != nil {
			return nil, err
		}
	}

	err := mboxMessage.Close()
	if err != nil {
		return nil, err
	}
	email.Body = messageBuffer
	return email, nil
}

func isMultiPartMessage(message *Message) bool {
	for key, header := range message.Header {
		if key == "Content-Type" && len(header) > 0 {
			if strings.HasPrefix(header[0], "multipart/") {
				return true
			}
		}
	}
	return false
}

func processPlainEmail(email *mail.Message, message *Message) (*mail.Message, error) {
	for key, header := range message.Header {
		email.Header[http.CanonicalHeaderKey(key)] = header
	}
	if len(message.Bodies) != 1 {
		return nil, errors.New("Found multiple bodies in non multi-part email")
	}

	body := message.Bodies[0]
	contentType := body.Header.Get("Content-Type")
	email.Header["Content-Type"] = []string{contentType}
	email.Body = body.Content
	return email, nil
}

func streamToBytes(stream io.Reader) []byte {
	buf := new(bytes.Buffer)
	buf.ReadFrom(stream)
	return buf.Bytes()
}
