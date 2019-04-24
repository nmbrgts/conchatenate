package main

import (
	"testing"
)

func TestIdMessage(t *testing.T) {
	t.Run(
		"TestIdMessage should return an (integer, true) ID when `GetSenderId` invoked.",
		func(t *testing.T) {
			want := 1
			msg := IdMessage{want, ""}
			got, ok := msg.GetSenderId()
			if want != got || !ok {
				t.Errorf("Wanted return value of (%d, %t) but got (%d, %t)",
					want, true, got, ok)
			}
		},
	)
	t.Run(
		"TestIdMessage should return a (string, true) of content when `GetContent` invoked.",
		func(t *testing.T) {
			want := "hallo, this is dog"
			msg := IdMessage{0, want}
			got, ok := msg.GetContent()
			if want != got || !ok {
				t.Errorf("Wanted return value of (\"%s\", %t) but got (\"%s\", %t)",
					want, true, got, ok)
			}
		},
	)
}
func TestPlainMessage(t *testing.T) {
	t.Run(
		"PlainMessage should return a (0, false) ID when `GetSenderId` invoked.",
		func(t *testing.T) {
			want := 0
			msg := PlainMessage{""}
			got, ok := msg.GetSenderId()
			if want != got || ok {
				t.Errorf("Wanted return value of (%d, %t) but got (%d, %t)",
					want, false, got, ok)
			}
		},
	)
	t.Run(
		"PlainMessage should return a (string, true) of content when `GetContent` invoked.",
		func(t *testing.T) {
			want := "hallo, this is dog"
			msg := PlainMessage{want}
			got, ok := msg.GetContent()
			if want != got || !ok {
				t.Errorf("Wanted return value of (\"%s\", %t) but got (\"%s\", %t)",
					want, true, got, ok)
			}
		},
	)
}
func TestNilMessage(t *testing.T) {
	t.Run(
		"NilMessage should return a (0, false) ID when `GetSenderId` invoked.",
		func(t *testing.T) {
			want := 1
			msg := NilMessage{want}
			got, ok := msg.GetSenderId()
			if want != got || !ok {
				t.Errorf("Wanted return value of (%d, %t) but got (%d, %t)",
					want, true, got, ok)
			}
		},
	)
	t.Run(
		"NilMessage should return a (string, true) of content when `GetContent` invoked.",
		func(t *testing.T) {
			want := ""
			msg := NilMessage{1}
			got, ok := msg.GetContent()
			if want != got || ok {
				t.Errorf("Wanted return value of (\"\", %t) but got (\"%s\", %t)",
					false, got, ok)
			}
		},
	)
}
