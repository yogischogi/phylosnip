package snp

import (
	"encoding/csv"
	"os"
	"strconv"
	"strings"
)

// BEDRegion is a region on the Y-chromosome.
type BEDRegion struct {
	Start int
	End   int
}

// Includes test if the given position is included in the region.
func (b *BEDRegion) Includes(pos int) bool {
	return pos >= b.Start && pos < b.End
}

type BEDRegions []BEDRegion

// Includes test if the given position is included in the regions.
func (b *BEDRegions) Includes(pos int) bool {
	for _, region := range *b {
		if region.Includes(pos) {
			return true
		}
	}
	return false
}

// ReadBED reads Y-chromosome regions from a BED file
// as described in http://genome.ucsc.edu/FAQ/FAQformat#format1
func ReadBED(filename string) (BEDRegions, error) {
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

	var result BEDRegions
	for _, r := range records {
		exists, region := bedFromRecord(r)
		if exists {
			result = append(result, region)
		}
	}
	return result, nil
}

// bedFromRecord tries to convert a line from a BED file into a
// BEDRegion. Note that not all lines in a BED file describe regions.
func bedFromRecord(fields []string) (exists bool, region BEDRegion) {
	if len(fields) < 3 {
		return false, region
	}
	if strings.Index(fields[0], "chrY") != 0 {
		return false, region
	}
	var err1, err2 error
	region.Start, err1 = strconv.Atoi(fields[1])
	region.End, err2 = strconv.Atoi(fields[2])
	if err1 != nil || err2 != nil {
		return false, region
	}
	return true, region
}
