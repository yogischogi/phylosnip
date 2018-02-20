package cmd

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/yogischogi/phylosnip/snp"
)

// Union calculates the set union of the SNPs contained in the files
// specified be the parameter a.
// cmdLine: command line parameters without the subcommand.
func Union(cmdLine []string) {
	flags := flag.NewFlagSet("", flag.ContinueOnError)
	var (
		in  = flags.String("in", "", "Input list of CSV files separated by commas.")
		out = flags.String("out", "", "Output file in CSV format.")
	)
	flags.Parse(cmdLine)

	if *in == "" {
		fmt.Printf("Parameter in for input files not specified.\n")
		os.Exit(1)
	}

	// Parse filename parameter.
	filenames, err := parameterToFilenames(*in, ".csv")
	checkFatal(err, "Error parsing parameter in")

	// Calculate union.
	var snps snp.SNPs = make(map[snp.SNP]bool)
	err = unionWithFiles(snps, filenames)
	checkFatal(err, "Error calculating the union of files")

	// Output SNPs.
	if *out != "" {
		err := snps.WriteCSV(*out)
		checkFatal(err, "Error writing SNPs to CSV output file")
	} else {
		for snp, _ := range snps {
			fmt.Print(snp.String())
		}
	}
}

// unionWithFiles calculates the set union of snps with the SNPs
// contained in the files specified by filenames.
func unionWithFiles(snps snp.SNPs, filenames []string) error {
	for _, filename := range filenames {
		s, err := snp.ReadCSV(filename)
		if err != nil {
			return errors.New(fmt.Sprintf("reading CSV file, %v", err))
		}
		snps.Union(s)
	}
	return nil
}
