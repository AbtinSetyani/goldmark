package goldmark_test

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	goldmark "github.com/anytypeio/goldmark"
	"github.com/anytypeio/goldmark/blocksUtil"
	htmlToMdConverter "github.com/anytypeio/html-to-markdown"
	"io/ioutil"
	"log"
	"os/exec"
	"testing"
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

	s := testCases[4].HTML
/*	asc := strings.Map(func(r rune) rune {
		if r > unicode.MaxASCII {
			return -1
		}
		return r
	}, )*/

/*	re := regexp.MustCompile("[[:^ascii:]]")
	reSpace := regexp.MustCompile(`[\s]+`)
	//re2 := regexp.MustCompile("[[:^ascii:]]")
	str := re.ReplaceAllLiteralString(s, "")
	str = str*/

/*	md := html2md.Convert(s)
	md = spaceReplace.WhitespaceNormalizeString(md)
	md = strings.ReplaceAll(md, "\n\n\n", "@#par-mark$")
	md = strings.ReplaceAll(md, "\n", "")
	md = strings.ReplaceAll(md, "@#par-mark$","\n\n")
	fmt.Println(md)*/

	mdToBlocksConverter := goldmark.New()
	_, blocks := mdToBlocksConverter.HTMLToBlocks([]byte(s))
	fmt.Println(1, "html:", testCases[4].HTML)
	for i, b := range blocks {
		fmt.Println(i," block: ", b)
	}

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