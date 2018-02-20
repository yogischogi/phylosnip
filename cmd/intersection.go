package cmd

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/yogischogi/phylosnip/snp"
)

// Intersection calculates the set intersection of the SNPs contained in the files
// specified be the parameters a and b.
// cmdLine: command line parameters without the subcommand.
func Intersection(cmdLine []string) {
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
	if len(filenames) == 0 {
		fmt.Printf("No files found for input parameter in.\n")
		os.Exit(1)
	}

	// Calculate intersection.
	snps, err := snp.ReadCSV(filenames[0])
	checkFatal(err, "Error reading CSV file")
	err = intersectionWithFiles(snps, filenames[1:])
	checkFatal(err, "Error calculating the intersection of files")

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

// intersectionWithFiles calculates the set intersection of snps with the SNPs
// contained in the files specified by filenames.
func intersectionWithFiles(snps snp.SNPs, filenames []string) error {
	for _, filename := range filenames {
		s, err := snp.ReadCSV(filename)
		if err != nil {
			return errors.New(fmt.Sprintf("reading CSV file, %v", err))
		}
		snps.Intersection(s)
	}
	return nil
}
