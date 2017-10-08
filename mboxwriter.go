package mboxparser

import (
	"bufio"
	"github.com/blabber/mbox"
	"net/mail"
	"os"
)

func WriteFile(mboxFile *Mbox, outputPath string) error {
	var encodedMessages []*mail.Message
	for _, mboxMessage := range mboxFile.Messages {
		encodedMessage, err := Encode(mboxMessage)
		if err != nil {
			return err
		}
		encodedMessages = append(encodedMessages, encodedMessage)
	}

	fileToWrite, err := os.Create(outputPath)
	if err != nil {
		return err
	}

	fileWriter := bufio.NewWriter(fileToWrite)
	defer fileWriter.Flush()
	mboxWriter := mbox.NewWriter(fileWriter)

	for _, mailMessage := range encodedMessages {
		// TODO: What shoudl I do with the bytes written var?
		_, err := mboxWriter.WriteMessage(mailMessage)
		if err != nil {
			return err
		}
	}
	return nil
}
