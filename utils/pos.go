package utils

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"golang.org/x/net/html/charset"
)

type PDFpage struct {
	Top      int    `xml:"top,attr"`
	Left     int    `xml:"left,attr"`
	Width    int    `xml:"width,attr"`
	Height   int    `xml:"height,attr"`
	Number   int    `xml:"number,attr"`
	Position string `xml:"position,attr"`

	Texts []PDFtext `xml:"text"`
}

type PDFtext struct {
	Top    int    `xml:"top,attr"`
	Left   int    `xml:"left,attr"`
	Font   int    `xml:"font,attr"`
	Width  int    `xml:"width,attr"`
	Height int    `xml:"height,attr"`
	Value  string `xml:",chardata"`
}

type PDF struct {
	Pages []PDFpage `xml:"page"`
}

func (this *PDF) ocrConvert(filePath string) (ocrFilepath string) {

	var filePdfs, fileImages []string
	fileDir := filepath.Base(filePath)
	fileDir = strings.TrimSuffix(fileDir, filepath.Ext(fileDir))

	if fileDir == "" {
		return
	}

	if err := CreateDir(fileDir); err != nil {
		log.Printf(err.Error())
		return
	}
	defer RemoveFileDir(fileDir)

	cmdPdfimages := `pdfimages`
	cmdPdfimagesArgs := []string{"-j", filePath, fileDir + "/pdf"}
	if _, err := ExecCommand(cmdPdfimages, cmdPdfimagesArgs); err != nil {
		return
	}

	cmdJpegoptimArgs := []string{"-c", fmt.Sprintf("jpegoptim -f -m90 -q %s/*jpg", fileDir)}
	ExecCommand("/bin/sh", cmdJpegoptimArgs)

	fileImages, err := ListDir("/" + fileDir)
	if err != nil || len(fileImages) == 0 {
		log.Printf(err.Error())
		return
	}

	ocrChan := make(chan bool)
	for _, fileImage := range fileImages {
		if strings.HasSuffix(fileImage, ".jpg") {
			go func(fileImagePath string) {
				filePdf := strings.TrimSuffix(fileImagePath, filepath.Ext(fileImagePath)) + ".pdf"
				ExecCommand(`./curl_recognize.sh`, []string{fileImagePath, filePdf, "-f", "pdfTextAndImages"})
				ocrChan <- true
			}(fileDir + "/" + fileImage)
		}
	}
	for _, fileImage := range fileImages {
		if strings.HasSuffix(fileImage, ".jpg") {
			<-ocrChan
		}
	}

	if filePdfs, err = ListDir("/" + fileDir); err != nil || len(filePdfs) == 0 {
		log.Printf(err.Error())
		return
	}

	var cmdPdfuniteArgs []string
	for _, fileName := range filePdfs {
		if strings.HasSuffix(fileName, ".pdf") {
			cmdPdfuniteArgs = append(cmdPdfuniteArgs, fileDir+"/"+fileName)
		}
	}

	cmdPdfunite := `pdfunite`
	if len(cmdPdfuniteArgs) > 0 {
		ocrFilepathNew := strings.TrimSuffix(filePath, filepath.Ext(filePath))
		cmdPdfuniteArgs = append(cmdPdfuniteArgs, fmt.Sprintf("%s-OCR.pdf", ocrFilepathNew))
		ocrFilepath = cmdPdfuniteArgs[len(cmdPdfuniteArgs)-1]
		ExecCommand(cmdPdfunite, cmdPdfuniteArgs)
	}

	return
}

//postscipt convert
func (this *PDF) Ps2pdfConvert(filePath string) (pdfFilepath string) {
	pdfFilepath = filePath + ".pdf"
	cmdPs2pdf := `ps2pdfwr`
	cmdPs2pdfArgs := []string{filePath, pdfFilepath}
	if _, err := ExecCommand(cmdPs2pdf, cmdPs2pdfArgs); err != nil {
		log.Println(err.Error())
		pdfFilepath = ""
		return
	}
	return
}

func (this *PDF) Scrape(source, password, filePath string) error {

	// <-time.NewTimer(time.Second * 3).C

	if strings.Compare(source, "Scanned PDF") == 0 {
		ocrFilepath := fmt.Sprintf("%s-OCR.pdf", strings.TrimSuffix(filePath, filepath.Ext(filePath)))
		if _, err := os.Stat(ocrFilepath); err == nil {
			filePath = ocrFilepath
		} else {
			if ocrFilepath = this.ocrConvert(filePath); strings.HasSuffix(ocrFilepath, "-OCR.pdf") {
				filePath = ocrFilepath
			}
		}
	}

	cmdExec := `pdftohtml`
	cmdArgs := []string{filePath, "-i", "-q", "-stdout", "-xml"}

	if len(password) > 0 {
		cmdArgs = []string{filePath, "-i", "-q", "-stdout", "-xml", "-opw", strings.TrimSpace(password)}
	}

	cmdOutput, err := exec.Command(cmdExec, cmdArgs...).CombinedOutput()
	if err != nil {
		fmt.Printf("%v %v", err, cmdArgs)
		return err
	}

	replacer := strings.NewReplacer(`<A`, `<a`, `</span><span class="ft4">`, ``, `<b>`, ``, `</b>`, ``)
	cmdOutput = []byte(replacer.Replace(string(cmdOutput)))

	if strings.Contains(string(cmdOutput), "Error:") {
		return errors.New(string(cmdOutput))
	}

	decoder := xml.NewDecoder(bytes.NewReader(cmdOutput))
	decoder.CharsetReader = charset.NewReaderLabel
	err = decoder.Decode(&this)
	if err != nil {
		return err
	}

	for _, curPage := range this.Pages {
		sort.Sort(ByTop(curPage.Texts))

		var curTop int
		var textRow []PDFtext
		var newPage PDFpage
		for index, text := range curPage.Texts {
			if index > 0 && curTop != text.Top {
				sort.Sort(ByLeft(textRow))
				newPage.Texts = append(newPage.Texts, textRow...)
				textRow = textRow[:0]
			}
			textRow = append(textRow, text)
			curTop = text.Top
		}
		if len(textRow) > 0 {
			sort.Sort(ByLeft(textRow))
			newPage.Texts = append(newPage.Texts, textRow...)
		}
		curPage.Texts = curPage.Texts[:0]
		curPage.Texts = newPage.Texts
	}

	return nil
}

// ByTop implements sort.Interface for []PDFtext based on the Top field
type ByTop []PDFtext

func (a ByTop) Len() int           { return len(a) }
func (a ByTop) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTop) Less(i, j int) bool { return a[i].Top < a[j].Top }

// ByLeft implements sort.Interface for []PDFtext based on the Left field
type ByLeft []PDFtext

func (a ByLeft) Len() int           { return len(a) }
func (a ByLeft) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByLeft) Less(i, j int) bool { return a[i].Left < a[j].Left }
