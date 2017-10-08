package mboxparser

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestReadWriteFile(t *testing.T) {
	mbox, err := ReadFile("testdata/simple.mbox")
	if err != nil {
		t.Fatal(err)
	}

	if len(mbox.Messages) != 1 {
		t.Fatalf("invalid messages count: %d", len(mbox.Messages))
	}

	WriteFile(mbox, "testdata/simple-output.mbox")
	if _, err := os.Stat("testdata/simple-output.mbox"); os.IsNotExist(err) {
		t.Errorf("should write procesed file", err)
	}

	mbox2, err := ReadFile("testdata/simple-output.mbox")
	if err != nil {
		t.Fatal(err)
	}

	if len(mbox2.Messages) != 1 {
		t.Fatalf("invalid messages count: %d", len(mbox.Messages))
	}
}

func TestReadWriteNoBody(t *testing.T) {
	mbox, err := ReadFile("testdata/nobody.mbox")
	if err != nil {
		t.Fatal(err)
	}

	message := mbox.Messages[0]
	if len(message.Bodies) != 0 {
		t.Fatal("Invalid body found")
	}

	WriteFile(mbox, "testdata/nobody-output.mbox")
	if _, err := os.Stat("testdata/nobody-output.mbox"); os.IsNotExist(err) {
		t.Errorf("should write procesed file", err)
	}

	mbox2, err := ReadFile("testdata/nobody-output.mbox")
	if err != nil {
		t.Fatal(err)
	}

	if len(mbox2.Messages) != 1 {
		t.Fatalf("invalid messages count: %d", len(mbox.Messages))
	}
}

func TestReadWritePlain(t *testing.T) {
	mbox, err := ReadFile("testdata/plain.mbox")
	if err != nil {
		t.Fatal(err)
	}

	message := mbox.Messages[0]
	if len(message.Bodies) != 1 {
		t.Fatal("Invalid body found")
	}

	body := message.Bodies[0]
	buf := new(bytes.Buffer)
	buf.ReadFrom(body.Content)

	if strings.TrimSpace(buf.String()) != "This is a test" {
		t.Fatalf("Invalid body: %s", buf.String())
	}
	body.Content = buf

	WriteFile(mbox, "testdata/plain-output.mbox")
	if _, err := os.Stat("testdata/plain-output.mbox"); os.IsNotExist(err) {
		t.Errorf("should write procesed file", err)
	}

	mbox2, err := ReadFile("testdata/plain-output.mbox")
	if err != nil {
		t.Fatal(err)
	}

	message2 := mbox2.Messages[0]
	if len(message2.Bodies) != 1 {
		t.Fatal("Invalid body found")
	}

	body2 := message2.Bodies[0]
	buf2 := new(bytes.Buffer)
	buf2.ReadFrom(body2.Content)

	if strings.TrimSpace(buf2.String()) != "This is a test" {
		t.Fatalf("Invalid body: %s", buf2.String())
	}
}
