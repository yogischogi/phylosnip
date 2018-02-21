package snp

import (
	"encoding/csv"
	"os"
	"strconv"
)

// DB is an SNP data base.
type DB struct {
	snpRecords map[SNP]*DBRecord
	snpNames   map[string]*DBRecord
}

// DBRecord is an entry in the SNP data base.
// The entry is simplifies and does not contain all information
// that is listed at ISOGG.
type DBRecord struct {
	Key     SNP
	Name    string
	Comment string
}

func NewDB() *DB {
	return &DB{snpRecords: make(map[SNP]*DBRecord), snpNames: make(map[string]*DBRecord)}
}

func (db *DB) Add(entry DBRecord) {
	db.snpRecords[entry.Key] = &entry
	if entry.Name != "" {
		db.snpNames[entry.Name] = &entry
	}
}

// ReadISOGGcsv reads an ISOGG CSV file and adds the SNPs
// to the data base.
// The ISOGG file format can be found at http://ybrowse.org/gbrowse2/gff/.
func (db *DB) ReadISOGGcsv(filename string) error {
	infile, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer infile.Close()

	// Read all CSV records from file.
	csvReader := csv.NewReader(infile)
	csvReader.LazyQuotes = true
	//csvReader.FieldsPerRecord = -1
	records, err := csvReader.ReadAll()
	if err != nil {
		return err
	}

	for _, record := range records {
		entry, exists := newDBRecord(record)
		if exists {
			db.Add(entry)
		}
	}
	return nil
}

// newDBRecord creates a new DBRecord from the fields of a CSV file line.
// The CSV file must be in ISOGG format.
// http://ybrowse.org/gbrowse2/gff/
func newDBRecord(fields []string) (entry DBRecord, exists bool) {
	const (
		pos     = 3
		name    = 8
		ref     = 10
		alt     = 11
		comment = 18
	)
	if len(fields) < 19 {
		return entry, false
	}
	snpPos, err := strconv.ParseInt(fields[pos], 10, 0)
	if err != nil {
		return entry, false
	}
	entry.Key = SNP{Pos: int(snpPos), Ref: fields[ref], Alt: fields[alt]}
	entry.Name = fields[name]
	entry.Comment = fields[comment]
	exists = true
	return
}

func (db *DB) EntryByName(name string) (snp *DBRecord, exists bool) {
	snp, exists = db.snpNames[name]
	return
}

func (db *DB) EntryByKey(snp SNP) (entry *DBRecord, exists bool) {
	entry, exists = db.snpRecords[snp]
	return
}

/*
func (db *DB) EntryByPosition(position int) (snps []DBRecord, exists bool) {

}
*/
