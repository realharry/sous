package cli

import (
	"bytes"
	"flag"
	"strings"
	"testing"
)

var expectedHelpText = `
  -offset string
        source code relative repository offset
  -repo string
        source code repository location
  -revision string
        source code revision ID
  -tag string
        source code revision tag
`

func TestAddFlags(t *testing.T) {
	fs := flag.NewFlagSet("source", flag.ContinueOnError)

	actual := SourceFlags{}

	if err := AddFlags(fs, &actual, sourceFlagsHelp); err != nil {
		t.Fatal(err)
	}

	expected := SourceFlags{
		Repo:     "github.com/opentable/sous",
		Offset:   "",
		Tag:      "v1.0.0",
		Revision: "cabba9e",
	}

	args := []string{
		"-repo", expected.Repo,
		"-offset", expected.Offset,
		"-tag", expected.Tag,
		"-revision", expected.Revision,
	}
	if err := fs.Parse(args); err != nil {
		t.Fatal(err)
	}
	if actual.Repo != expected.Repo {
		t.Errorf("got %q; want %q", actual.Repo, expected.Repo)
	}
	if actual.Offset != expected.Offset {
		t.Errorf("got %q; want %q", actual.Offset, expected.Offset)
	}
	if actual.Tag != expected.Tag {
		t.Errorf("got %q; want %q", actual.Tag, expected.Tag)
	}
	if actual.Revision != expected.Revision {
		t.Errorf("got %q; want %q", actual.Revision, expected.Revision)
	}
	buf := &bytes.Buffer{}
	fs.SetOutput(buf)
	fs.PrintDefaults()
	actualHelp := strings.TrimSpace(buf.String())
	expectedHelp := strings.TrimSpace(expectedHelpText)
	actualFields := strings.Fields(actualHelp)
	expectedFields := strings.Fields(expectedHelp)
	// we're comparing the same words in the same order rather than being
	// concerned with whitespace differences.
	for i := range actualFields {
		if len(expectedFields)-1 < i || (actualFields[i] != expectedFields[i]) {
			t.Errorf("got help text:\n%s\nwant:\n%s", actualHelp, expectedHelp)
		}
	}
}

func TestParseUsage(t *testing.T) {
	in := `
		-someflag
			some usage text
	`
	out, err := parseUsage(in)
	if err != nil {
		t.Fatal(err)
	}
	actual, ok := out["someflag"]
	expected := "some usage text"
	if !ok {
		t.Fatalf("no usage text for -someflag; want %q", expected)
	}
	if actual != expected {
		t.Errorf("got %q; want %q", actual, expected)
	}
}

type AddFlagsInput struct {
	FlagSet *flag.FlagSet
	Target  interface{}
	Help    string
}

type BadFlagStruct struct {
	PtrField *string
}

func newFS() *flag.FlagSet { return flag.NewFlagSet("", flag.ContinueOnError) }

func TestAddFlags_badInputs(t *testing.T) {
	var s string
	stringPtr := &s
	var badAddFlagsInputs = map[AddFlagsInput]string{
		{nil, nil, ""}:                                     "cannot add flags to nil *flag.FlagSet",
		{newFS(), nil, ""}:                                 "target is <nil>; want pointer to struct",
		{newFS(), "", ""}:                                  "target is string; want pointer to struct",
		{newFS(), SourceFlags{}, ""}:                       "target is cli.SourceFlags; want pointer to struct",
		{newFS(), stringPtr, ""}:                           "target is *string; want pointer to struct",
		{newFS(), &BadFlagStruct{}, "\t-ptrfield\n\tblah"}: "target field cli.BadFlagStruct.PtrField is *string; want string, int",
		{newFS(), &SourceFlags{}, ""}:                      "no usage text for flag -repo",
	}
	for in, expected := range badAddFlagsInputs {
		actualErr := AddFlags(in.FlagSet, in.Target, in.Help)
		if actualErr == nil {
			t.Fatalf("got nil; want error %q", expected)
		}
		actual := actualErr.Error()
		if actual != expected {
			t.Errorf("got error %q; want error %q", actual, expected)
		}
	}
}