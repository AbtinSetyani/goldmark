package goldmark_test

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	htmlToMdConverter "github.com/anytypeio/html-to-markdown"
	"log"

	"github.com/anytypeio/goldmark/blocksUtil"
	"io/ioutil"
	"os/exec"
	"testing"

	goldmark "github.com/anytypeio/goldmark"
)

var (
	pasteCmdArgs = "pbpaste"
	copyCmdArgs  = "pbcopy"
)

func getPasteCommand() *exec.Cmd {
	return exec.Command(pasteCmdArgs)
}

func readAll() (string, error) {
	pasteCmd := getPasteCommand()
	out, err := pasteCmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

type TestCase struct {
	HTML      string `json:"html"`
}

func TestConvertHTMLToBlocks(t *testing.T) {
	bs, err := ioutil.ReadFile("_test/spec.json")
	if err != nil {
		panic(err)
	}
	var testCases []TestCase
	if err := json.Unmarshal(bs, &testCases); err != nil {
		panic(err)
	}

	for _, c := range testCases {
		convertToBlocksAndPrint(c.HTML)
	}
/*	cases := []testutil.MarkdownTestCase{}
	for _, c := range testCases {
		cases = append(cases, testutil.MarkdownBlockTestCase{
			HTML: c.HTML,
		})
	}
	markdown := New(WithRendererOptions(
		html.WithXHTML(),
		html.WithUnsafe(),
	))
	testutil.DoTestCases(markdown, cases, t)	*/


	//sample1 := "<body><h1>My First Heading</h1><p>My first paragraph.</p></body>"

	//convertToBlocksAndPrint(sample1)
}

func TestConvertHTMLToBlocks2(t *testing.T) {
	bs, err := ioutil.ReadFile("_test/testData.json")
	if err != nil {
		panic(err)
	}
	var testCases []TestCase
	if err := json.Unmarshal(bs, &testCases); err != nil {
		panic(err)
	}

	//for i, c := range testCases {
		mdToBlocksConverter := goldmark.New()
		_, blocks := mdToBlocksConverter.HTMLToBlocks([]byte(testCases[3].HTML))
		fmt.Println(1, "html:", testCases[3].HTML)
		for i, b := range blocks {
			fmt.Println(i," block: ", b)
		}

	//}
}

func convertToBlocksAndPrint (html string) error {
	mdToBlocksConverter := goldmark.New()

	converter := htmlToMdConverter.NewConverter("", true, nil)
	md, err := converter.ConvertString(html)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println("md ->", md)

	var b bytes.Buffer
	writer := bufio.NewWriter(&b)
	BR := blocksUtil.NewRWriter(writer)

	err = mdToBlocksConverter.ConvertBlocks([]byte(md), BR)
	if err != nil {
		return err
	}

	fmt.Println("blocks:", BR.GetBlocks())
	return nil
}