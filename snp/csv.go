package snp

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"os"
	"strconv"
)

// CSVRecord represents a single line in a FTDNA CSV file.
type CSVRecord struct {
	Pos,
	Ref,
	Alt,
	Name,
	Comment string
}

type CSVRecords []CSVRecord

// String returns a CSVRecord as a single line of text
// including the line ending CRLF.
func (c *CSVRecord) String() string {
	var b bytes.Buffer
	b.WriteString(c.Pos)
	b.WriteString(",")
	b.WriteString(c.Ref)
	b.WriteString(",")
	b.WriteString(c.Alt)
	b.WriteString(",")
	b.WriteString(c.Name)
	b.WriteString(",\"")
	b.WriteString(c.Comment)
	b.WriteString("\"\r\n")
	return b.String()
}

// ReadFTDNAcsv reads CSV records from a FTDNA encoded CSV file.
// If mutationsOnly == true, only true mutations are included in the result.
// If novelsOnly == true, only novel variants are reported.
func ReadFTDNAcsv(filename string, mutationsOnly bool, novelsOnly bool, db *DB) (CSVRecords, error) {
	infile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer infile.Close()

	// Read all CSV records from file.
	csvReader := csv.NewReader(infile)
	csvReader.FieldsPerRecord = -1
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	result := make([]CSVRecord, 0)
	for _, record := range records {
		rec, exists := ftdnaFieldsToCSVRecord(record, mutationsOnly, novelsOnly, db)
		if exists {
			result = append(result, rec)
		}
	}
	return result, nil
}

// ftdnaFieldsToSNP parses a line in FTDNA's CSV format and
// tries to extract an SNP mutation.
// fields are the fields of a single CSV line.
// If mutationsOnly == true, only true mutations are included in the result.
// If novelsOnly == true, only novel variants are reported.
func ftdnaFieldsToCSVRecord(fields []string, mutationsOnly bool, novelsOnly bool, db *DB) (rec CSVRecord, exists bool) {
	// Positions of the entries.
	const (
		variant = 0
		pos     = 1
		name    = 2
		ref     = 5
		alt     = 6
	)
	// Skip header.
	if fields[variant] == "Type" {
		return rec, false
	}
	// Skip uncertain values.
	if fields[alt] == "?" {
		return rec, false
	}
	if mutationsOnly && fields[ref] == fields[alt] {
		return rec, false
	}
	if novelsOnly && fields[variant] != "Novel Variant" {
		return rec, false
	}

	rec.Ref = fields[ref]
	rec.Alt = fields[alt]
	rec.Name = fields[name]
	_, err := strconv.Atoi(fields[pos])
	if err == nil {
		rec.Pos = fields[pos]
	} else {
		rec.Pos = "n/a"
	}

	if db != nil {
		snpEntry, found := db.EntryByName(fields[name])
		if found {
			rec.Comment = snpEntry.Comment
			if rec.Pos == "n/a" {
				rec.Pos = strconv.Itoa(snpEntry.Key.Pos)
			}
		}
	}
	return rec, true
}

func (c CSVRecords) WriteCSV(filename string) error {
	// Open file.
	outfile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer outfile.Close()

	w := bufio.NewWriter(outfile)
	for _, rec := range c {
		_, err := w.WriteString(rec.String())
		if err != nil {
			return err
		}
	}
	err = w.Flush()
	return err
}
