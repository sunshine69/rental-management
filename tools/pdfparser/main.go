package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"

	ag "github.com/sunshine69/automation-go/lib"
	u "github.com/sunshine69/golang-tools/utils"
	"github.com/sunshine69/rental-management/configs"
	"github.com/sunshine69/rental-management/model"
)

var Cfg configs.Config

func main() {
	filename := flag.String("f", "", "Path to pdf file")
	flag.Parse()
	fmt.Printf("INput file: '%s'\n", *filename)
	ParseQldRentalContract(*filename)
}

func ParseQldRentalContract(pdffile string) {
	tmpDir, err := os.MkdirTemp("", "pdfparser")
	u.CheckErr(err, "CreateTemp pdfparser ")
	defer os.RemoveAll(tmpDir)
	textfile := tmpDir + "/pdftotext.txt"
	o, err := u.RunSystemCommandV2("pdftotext '"+pdffile+"' "+textfile, true)
	u.CheckErr(err, "RunSystemCommandV2 pdftotext "+o)
	_, start, end, datalines := ag.ExtractTextBlockContains(textfile, []string{`Item 1.1 Lessor`}, []string{`Item 2.1 Tenant\/s`}, []string{`Name/trading name`})
	blocklines := datalines[start:end]

	ParseLessor(blocklines)
	// Parse tenant
	block, _, _, _ := ag.ExtractTextBlockContains(textfile, []string{`Item 2.1 Tenant\/s`}, []string{`2.2 Address for service`}, []string{`Tenant 1 Full name`})
	ParseTenant(block)

	_, start, end, blocklines = ag.ExtractTextBlock(textfile, []string{`Item 5.1 Address of the rental premises`}, []string{`5.2 Inclusions provided`})
	ParseProperty(blocklines[start:end])

	block, start, end, blocklines = ag.ExtractTextBlockContains(textfile, []string{`5.2 Inclusions provided`}, []string{`Part 2 Standard Terms`}, []string{`6.3 Ending on`})
	ParseContract(blocklines[start:end], block)
}

func ParseContract(blocklines []string, block string) {
	fixedTermPtn := regexp.MustCompile(`(?m)✔ fixed term agreement`)
	periodicTerm := regexp.MustCompile(`(?m)✔ periodic agreement`)
	startPtn := regexp.MustCompile(`(?m)6.2 Starting on\n([\d]+\/[\d]+\/[\d]+)`)
	endPtn := regexp.MustCompile(`(?m)6.3 Ending on[\n]+[^\d]+[\n]+([\d]+\/[\d]+\/[\d]+)`)
	if m := fixedTermPtn.FindStringSubmatch(block); m != nil {
		fmt.Println(m[0])
	}
	if m := periodicTerm.FindStringSubmatch(block); m != nil {
		fmt.Println(m[0])
	}
	if m := startPtn.FindStringSubmatch(block); m != nil {
		fmt.Println(m[1])
	}
	if m := endPtn.FindStringSubmatch(block); m != nil {
		println(m[1])
	}
	rentPtn := regexp.MustCompile(`(?m)Item Rent\n7\nper [✔]? week\n\$[\s]*([\d]+)`)
	if m := rentPtn.FindStringSubmatch(block); m != nil {
		println(m[1])
	}
}

func ParseProperty(blocklines []string) {
	postcodePtn := regexp.MustCompile(`Postcode ([\d]{4,4})`)
	addressPtn := regexp.MustCompile(`([a-zA-Z0-9]+ [a-zA-Z0-9]+ [a-zA-Z0-9]+[\s]*.*)`)
	suburbPtn := regexp.MustCompile(`([a-zA-Z]{3,})`)
	var address, postcode, suburb string
	for _, l := range blocklines[1:] {
		if m := postcodePtn.FindStringSubmatch(l); m != nil {
			postcode = m[1]
			continue
		}
		if m := addressPtn.FindStringSubmatch(l); m != nil {
			address = address + " " + m[0]
			continue
		}
		if m := suburbPtn.FindStringSubmatch(l); m != nil {
			suburb = m[1]
			continue
		}
	}
	property_code := strings.ReplaceAll(address, " ", "") + "_" + postcode
	pr := model.Property{Name: property_code, Address: address + ", " + suburb + " " + postcode}
	pr.Save()
}

func parseNames(name string) (firstName, lastName string) {
	lastName = name[strings.LastIndex(name, " ")+1:]
	firstName = strings.TrimSpace(strings.TrimSuffix(name, lastName))
	return
}

func ParseLessor(blocklines []string) {
	namePtn := regexp.MustCompile(`^Name/trading name (.*)$`)
	mobilePtn := regexp.MustCompile(`^([\d\s]+)$`)
	emailPtn := regexp.MustCompile(`([^\@]+\@[^\@]+)`)
	addressPtn := regexp.MustCompile(`([^,]+,[^,]+,[^,])`)
	postcodePtn := regexp.MustCompile(`([\d]{4,4})`)
	var fullname, mobile, email, address, postcode, firstName, lastName string
	for _, l := range blocklines {
		if m := namePtn.FindStringSubmatch(l); m != nil {
			fullname = m[1]
			firstName, lastName = parseNames(fullname)
			continue
		}
		if m := mobilePtn.FindStringSubmatch(l); m != nil {
			mobile = m[1]
			continue
		}
		if m := emailPtn.FindStringSubmatch(l); m != nil {
			email = m[1]
			continue
		}
		if m := addressPtn.FindStringSubmatch(l); m != nil {
			address = m[1]
			continue
		}
		if m := postcodePtn.FindStringSubmatch(l); m != nil {
			postcode = m[1]
			continue
		}

	}
	pm := model.Property_manager{Email: email, First_name: firstName, Last_name: lastName, Address: address + " " + postcode,
		Contact_number: mobile}
	pm.Save()
}

func ParseTenant(block string) {
	tenantBlocks := ag.SplitTextByPattern(block, `(?m)Tenant [\d]+ Full name`, true)
	println(u.JsonDump(tenantBlocks, ""))

	namePtn := regexp.MustCompile(`Full name\/s (.*)$`)
	mobilePtn := regexp.MustCompile(`^([\d\s]+)$`)
	emailPtn := regexp.MustCompile(`Email ([^\@]+\@[^\@]+)`)
	blockStart := regexp.MustCompile(`Tenant [\d]+`)
	var mobile, email, firstName, lastName string
	for _, b := range tenantBlocks {
		foundBlock := false
		for _, l := range strings.Split(b, "\n") {
			if blockStart.MatchString(l) {
				foundBlock = true
			}
			if m := namePtn.FindStringSubmatch(l); m != nil {
				firstName, lastName = parseNames(m[1])
				continue
			}
			if m := emailPtn.FindStringSubmatch(l); m != nil {
				email = m[1]
				continue
			}
			if m := mobilePtn.FindStringSubmatch(l); m != nil {
				mobile = m[1]
				continue
			}
		}
		if foundBlock {
			tn := model.Tenant{Email: email, Contact_number: mobile, First_name: firstName, Last_name: lastName}
			tn.Save()
			foundBlock = false
		}
	}
}
