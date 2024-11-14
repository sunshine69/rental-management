package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
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

// This is subject to change depending on the form content changes. I run the pdftotext manually to get the text file and view them to find out the pattern
// to extract. The gov. changes the form pretty much frequently
func ParseQldRentalContract(pdffile string) {
	tmpDir, err := os.MkdirTemp("", "pdfparser")
	u.CheckErr(err, "CreateTemp pdfparser ")
	defer os.RemoveAll(tmpDir)
	textfile := tmpDir + "/pdftotext.txt"
	o, err := u.RunSystemCommandV2("pdftotext '"+pdffile+"' "+textfile, true)
	u.CheckErr(err, "RunSystemCommandV2 pdftotext "+o)
	_, start, end, datalines := ag.ExtractTextBlockContains(textfile, []string{`Item 1.1 Lessor`}, []string{`Item 2.1 Tenant\/s`}, []string{`Name/trading name`})
	blocklines := datalines[start:end]

	pm := ParseLessor(blocklines)
	if pm.Email == "" {
		panic("[ERROR] can not parse lessor\n")
	}
	// Parse tenant
	block, _, _, _ := ag.ExtractTextBlockContains(textfile, []string{`Item 2.1 Tenant\/s`}, []string{`2.2 Address for service`}, []string{`1. Full name/s`})
	tns := ParseTenant(block)
	if len(tns) == 0 {
		panic("[ERROR] can not parse tenant\n")
	}
	_, start, end, blocklines = ag.ExtractTextBlock(textfile, []string{`Item 5.1 Address of the rental premises`}, []string{`5.2 Inclusions provided`})
	property_code := ParseProperty(blocklines[start:end])

	block, start, end, blocklines = ag.ExtractTextBlockContains(textfile, []string{`5.2 Inclusions provided`}, []string{`Part 2 Standard Terms`}, []string{`6.3 Ending on`})
	ParseContract(blocklines[start:end], block, property_code, pm, tns)
}

func ParseContract(blocklines []string, block, property_code string, pm *model.Property_manager, tns []model.Tenant) model.Contract {

	var start_date, end_date, term string
	println(u.JsonDump(blocklines, ""))
	start_block_lines := ag.ExtractLineInLines(blocklines, `6.3 Ending on`, `([\d]+\/[\d]+\/[\d]+)`, `([\d]+\/[\d]+\/[\d]+)`)
	println(u.JsonDump(start_block_lines, ""))
	start_date = start_block_lines[0][1]
	end_block_lines := ag.ExtractLineInLines(blocklines, `([\d]+\/[\d]+\/[\d]+)`, `([\d]+\/[\d]+\/[\d]+)`, `Fixed term agreements only`)
	end_date = end_block_lines[0][1]

	fixedTermPtn := regexp.MustCompile(`✔ fixed term agreement`)
	if fixed := fixedTermPtn.MatchString(block); fixed {
		term = "fixed"
	} else {
		term = "periodic"
	}

	contract := model.Contract{Property: property_code, Start_date: start_date, End_date: end_date, Term: term, Property_manager: pm.Email, Tenant_main: tns[0].Email}
	if o := ag.ExtractLineInLines(blocklines, `Item Rent`, `\$ ([\d]+)`, `Item Rent must be paid on the`); o != nil {
		if rent, err := strconv.ParseInt(o[0][1], 10, 64); err == nil {
			contract.Rent = rent
		}
	}
	contract.Term = term
	if len(tns) > 1 {
		contract.Tenants = u.JsonDump(tns[1:], "")
	}
	contract.Save()
	return contract
}

func ParseProperty(blocklines []string) (property_code string) {
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
	property_code = strings.ReplaceAll(address, " ", "") + "_" + postcode
	pr := model.Property{Code: property_code, Address: address + ", " + suburb + " " + postcode}
	pr.Save()
	return
}

func parseNames(name string) (firstName, lastName string) {
	lastName = name[strings.LastIndex(name, " ")+1:]
	firstName = strings.TrimSpace(strings.TrimSuffix(name, lastName))
	return
}

func ParseLessor(blocklines []string) *model.Property_manager {
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
	return &pm
}

func ParseTenant(block string) (tenants []model.Tenant) {
	tenantBlocks := ag.SplitTextByPattern(block, `(?m)[\d]\. Full name\/s ([a-zA-Z0-9\s]+)`, true)
	// println(u.JsonDump(tenantBlocks, ""))

	for _, b := range tenantBlocks {
		var mobile, email, firstName, lastName string
		datalines := strings.Split(b, "\n")
		// println(u.JsonDump(datalines, ""))

		if o := ag.ExtractLineInLines(datalines, `Full name\/s (.*)$`, `Email ([^\@]+\@[^\@]+)`, `Emergency contact full name`); o != nil {
			// println(email)
			email = o[0][1]
		}
		if o := ag.ExtractLineInLines(datalines, `Full name\/s (.*)$`, `Full name\/s (.*)$`, `^([\d\s]+)$`); o != nil {
			firstName, lastName = parseNames(o[0][1])
		}
		if o := ag.ExtractLineInLines(datalines, `Full name\/s (.*)$`, `^([\d\s]+)$`, `Emergency contact full name`); o != nil {
			// println(mobile)
			mobile = o[0][1]
		}
		if email != "" {
			tn := model.Tenant{Email: email, Contact_number: mobile, First_name: firstName, Last_name: lastName}
			tn.Save()
			tenants = append(tenants, tn)
		} else {
			if len(tenants) > 0 {
				tenants[0].Note = fmt.Sprintf("%s\n%s %s %s", tenants[0].Note, firstName, lastName, mobile)
				tenants[0].Save()
			}
		}
	}
	return
}
