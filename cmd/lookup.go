package cmd

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/yogischogi/phylosnip/snp"
)

// Lookup checks if SNPs exist in the ISOGG database and
// adds the ISOGG information to the CSV files.
func Lookup(cmdLine []string) {
	flags := flag.NewFlagSet("", flag.ContinueOnError)
	var (
		in      = flags.String("in", "", "Input file in FTDNA CSV format.")
		out     = flags.String("out", "", "Output file for list of SNPs in CSV format.")
		isoggdb = flags.String("isoggdb", "", "Input file for ISOGG SNP data base in CSV format.")
	)
	flags.Parse(cmdLine)

	if *in == "" {
		fmt.Printf("Parameter in not specified.\n")
		os.Exit(1)
	}

	if *isoggdb == "" {
		fmt.Printf("Parameter isoggdb not specified.\n")
		os.Exit(1)
	}

	var snpDB *snp.DB
	if *isoggdb != "" {
		snpDB = snp.NewDB()
		err := snpDB.ReadISOGGcsv(*isoggdb)
		checkFatal(err, "Error reading SNP definitions from ISOGG CSV file")
	}

	inFiles, outFiles, err := inToOutFilenames(*in, ".csv", *out, ".csv")
	checkFatal(err, "Error converting filenames from parameter in to out")
	for i, _ := range inFiles {
		snps, err := snp.ReadCSV(inFiles[i])
		checkFatal(err, "Error reading input CSV file")

		// Convert SNPs to CSV records with enhanced information.
		var records snp.CSVRecords
		for s, _ := range snps {
			var rec snp.CSVRecord
			dbRec, exists := snpDB.EntryByKey(s)
			if exists {
				rec = snp.CSVRecord{
					Pos:     strconv.Itoa(dbRec.Key.Pos),
					Ref:     dbRec.Key.Ref,
					Alt:     dbRec.Key.Alt,
					Name:    dbRec.Name,
					Comment: dbRec.Comment,
				}
			} else {
				rec = snp.CSVRecord{Pos: strconv.Itoa(s.Pos), Ref: s.Ref, Alt: s.Alt}
			}
			records = append(records, rec)
		}

		if outFiles[i] != "" {
			// Write to file.
			err := records.WriteCSV(outFiles[i])
			checkFatal(err, "Error writing to CSV file")
		} else {
			// Write to stdout.
			for _, r := range records {
				os.Stdout.WriteString(r.String())
			}
		}
	}
}
