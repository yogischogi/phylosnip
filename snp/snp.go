// Package snp provides common operations on SNP mutations.
package snp

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// SNP is an SNP mutation.
type SNP struct {
	Pos int
	Ref string
	Alt string
}

// SNPs represents a set of SNPs.
type SNPs map[SNP]bool

// String returns a string representation in CSV format
// including CRLF.
func (s SNP) String() string {
	var b bytes.Buffer
	b.WriteString(strconv.Itoa(s.Pos))
	b.WriteString(",")
	b.WriteString(s.Ref)
	b.WriteString(",")
	b.WriteString(s.Alt)
	b.WriteString("\r\n")
	return b.String()
}

// Union calculates the set union of all SNPs in a with all SNPs in b.
// The result is stored in a.
func (a SNPs) Union(b SNPs) {
	for k, _ := range b {
		a[k] = true
	}
}

// Intersection calculates the set intersection of all SNPs in a with all SNPs in b.
// The result is stored in a.
func (a SNPs) Intersection(b SNPs) {
	for k, _ := range a {
		if b[k] == false {
			delete(a, k)
		}
	}
}

// Difference calculates the set difference a\b.
// The result is stored in a.
func (a SNPs) Difference(b SNPs) {
	for k, _ := range b {
		delete(a, k)
	}
}

// WriteCSV writes SNPs in a simplified format:
// Pos, Ref, Alt.
func (s SNPs) WriteCSV(filename string) error {
	// Open file.
	outfile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer outfile.Close()

	writer := bufio.NewWriter(outfile)
	for snp, _ := range s {
		writer.WriteString(snp.String())
	}
	err = writer.Flush()
	return err
}

// ReadCSV reads SNPs from a simple CSV file.
// Format: Pos, Ref, Alt.
func ReadCSV(filename string) (SNPs, error) {
	infile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer infile.Close()

	// Read all CSV records from file.
	csvReader := csv.NewReader(infile)
	csvReader.Comment = '#'
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	result := make(map[SNP]bool)
	for _, fields := range records {
		snpPos, err := strconv.Atoi(fields[0])
		if err != nil {
			return result, errors.New(fmt.Sprintf(" parsing SNP position %v\n", err))
		}
		snp := SNP{Pos: snpPos, Ref: fields[1], Alt: fields[2]}
		result[snp] = true
	}
	return result, nil
}

// ReadVCF reads SNPs from a VCF (Variant Call Format) file.
// The specification for VCF files is at https://github.com/samtools/hts-specs.
// The quality parameter specifies the minimum quality that an
// SNP must have to be included. SNPs that have passed the quality
// test are always included. Set quality to +Inf if you only want to
// include SNPs that have passed the quality test.
// If mutationsOnly == true only mutations are reported.
// reads is the required miminum number of reads.
// ratio is the required minimum ratio of ALT to REF reads.
func ReadVCF(filename string, quality float64, mutationsOnly bool, reads, ratio int) (SNPs, error) {
	infile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer infile.Close()

	// Read all CSV records from file.
	csvReader := csv.NewReader(infile)
	csvReader.Comma = '\t'
	csvReader.Comment = '#'
	csvReader.FieldsPerRecord = -1
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	result := make(map[SNP]bool)
	for _, record := range records {
		snp, exists := fieldsToSNP(record, quality, mutationsOnly, reads, ratio)
		if exists {
			result[snp] = true
		}
	}
	return result, nil
}

// fieldsToSNP tries to convert the entries of a VCF file line
// into an SNP.
// minReads and ratio are the same parameters as in ReadVCF.
func fieldsToSNP(fields []string, quality float64, mutationsOnly bool, minReads, ratio int) (snp SNP, exists bool) {
	// Positions of the entries.
	const (
		pos     = 1
		ref     = 3
		alt     = 4
		qual    = 5
		filter  = 6
		details = 9
	)
	// Exclude invalid lines and non-SNP mutations
	if len(fields) < details+1 || len(fields[ref]) != 1 {
		return snp, false
	}
	// Check for quality.
	snpQuality, err := strconv.ParseFloat(fields[qual], 64)
	if err != nil {
		return snp, false
	}
	if fields[filter] != "PASS" && snpQuality < quality {
		return snp, false
	}

	snpPos, err := strconv.Atoi(fields[pos])
	if err != nil {
		return snp, false
	}

	// Determine allele names and number of reads.
	// alleles are all different values that were read.
	// The first position is REF.
	alleles := []string{fields[ref]}
	altFields := strings.Split(fields[alt], ",")
	alleles = append(alleles, altFields...)

	reads := make([]int, len(alleles), len(alleles))
	strDetails := strings.Split(fields[details], ":")
	strReads := strings.Split(strDetails[1], ",")
	for i, _ := range strReads {
		r, err := strconv.Atoi(strReads[i])
		if err != nil {
			return snp, false
		}
		reads[i] = r
	}

	value, exists := snpValue(alleles, reads, minReads, ratio)
	if exists {
		if fields[ref] != value || mutationsOnly == false {
			return SNP{Pos: snpPos, Ref: fields[ref], Alt: value}, true
		}
	}
	return snp, false
}

// snpValues determines the valid allele value from a number of alleles
// and their number of reads.
func snpValue(alleles []string, reads []int, minReads, ratio int) (snp string, exists bool) {
	// Determine maximum number of reads.
	max := 0
	iMax := 0
	for i, r := range reads {
		if r > max {
			max = r
			iMax = i
		}
	}
	if max < minReads {
		return snp, false
	}
	// Check minimum ratio for reads.
	for i, r := range reads {
		if r == 0 || i == iMax {
			continue
		}
		if max/r < ratio {
			return snp, false
		}
	}

	if len(alleles[iMax]) == 1 {
		return alleles[iMax], true
	}
	return snp, false
}

// ReadYFull reads SNPs from a CSV encoded YFull file with
// novel SNPs.
// Quality: best, aceptable or ambiguous
func ReadYFull(filename string, quality string) (SNPs, error) {
	infile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer infile.Close()

	// Read all CSV records from file.
	csvReader := csv.NewReader(infile)
	csvReader.Comma = ';'
	csvReader.FieldsPerRecord = -1
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	result := make(map[SNP]bool)
	for _, record := range records {
		snp, exists := yfullFieldsToSNP(record, quality)
		if exists {
			result[snp] = true
		}
	}
	return result, nil
}

// yfullFieldsToSNP tries to convert the entries of a YFull line
// into an SNP.
// quality like in ReadYFull.
func yfullFieldsToSNP(fields []string, quality string) (snp SNP, exists bool) {
	// Positions of the entries.
	const (
		pos  = 2
		ref  = 3
		alt  = 4
		qual = 6
	)

	// Exclude invalid lines and non-SNP mutations
	if len(fields) < 7 || len(fields[ref]) != 1 || len(fields[alt]) != 1 {
		return snp, false
	}

	// Extract SNP.
	snpPos, err := strconv.Atoi(fields[pos])
	if err != nil {
		return snp, false
	}
	snp = SNP{Pos: snpPos, Ref: fields[ref], Alt: fields[alt]}

	// Check SNP quality.
	var qLevel = map[string]int{
		"best":            1,
		"acceptable":      2,
		"ambiguous":       3,
		"Best qual":       1,
		"Acceptable qual": 2,
		"Ambiguous qual":  3,
	}
	qSNP := qLevel[fields[qual]]
	qRequired := qLevel[quality]
	if qSNP <= qRequired {
		return snp, true
	}
	return snp, false
}
