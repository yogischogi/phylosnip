package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/yogischogi/phylosnip/snp"
)

// FilterFTDNA filters FTDNA CSV files for SNPs.
// cmdLine: command line parameters without the subcommand.
func FilterFTDNA(cmdLine []string) {
	flags := flag.NewFlagSet("", flag.ContinueOnError)
	var (
		in            = flags.String("in", "", "Input file in FTDNA CSV format.")
		out           = flags.String("out", "", "Output file for list of SNPs in CSV format.")
		mutationsonly = flags.Bool("mutationsonly", true, "If mutationsonly=true only mutations are reported.")
		novelsonly    = flags.Bool("novelsonly", false, "If novelsonly=true only novel variants are reported.")
		isoggdb       = flags.String("isoggdb", "", "Input file for ISOGG SNP data base in CSV format.")
	)
	flags.Parse(cmdLine)

	if *in == "" {
		fmt.Printf("Parameter in not specified.\n")
		os.Exit(1)
	}

	if *in == *out {
		fmt.Printf("Parameter in and out may not be identical.\n")
		os.Exit(1)
	}

	var snpDB *snp.DB
	if *isoggdb != "" {
		snpDB = snp.NewDB()
		err := snpDB.ReadISOGGcsv(*isoggdb)
		checkFatal(err, "Error reading SNP definitions from ISOGG CSV file")
	}

	inNames, outNames, err := inToOutFilenames(*in, ".csv", *out, ".csv")
	checkFatal(err, "Error converting filenames from parameter in to out")
	for i, _ := range inNames {
		recs, err := snp.ReadFTDNAcsv(inNames[i], *mutationsonly, *novelsonly, snpDB)
		checkFatal(err, "Error reading FTDNA CSV file")
		if outNames[i] != "" {
			// Write to file.
			err := recs.WriteCSV(outNames[i])
			checkFatal(err, "Error writing to CSV file")
		} else {
			// Write to stdout.
			for _, r := range recs {
				os.Stdout.WriteString(r.String())
			}
		}
	}
}
