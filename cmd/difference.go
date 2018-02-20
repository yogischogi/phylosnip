package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/yogischogi/phylosnip/snp"
)

// Difference calculates the set difference a\b of the SNPs contained in the files
// specified be the parameters a and b.
// cmdLine: command line parameters without the subcommand.
func Difference(cmdLine []string) {
	flags := flag.NewFlagSet("", flag.ContinueOnError)
	var (
		ain = flags.String("ain", "", "Input list of CSV files separated by commas.")
		bin = flags.String("bin", "", "Input list of CSV files separated by commas.")
		out = flags.String("out", "", "Output file in CSV format.")
	)
	flags.Parse(cmdLine)

	if *ain == "" {
		fmt.Printf("Parameter ain for input files not specified.\n")
		os.Exit(1)
	}
	if *bin == "" {
		fmt.Printf("Parameter bin for input files not specified.\n")
		os.Exit(1)
	}

	// Parse filename parameter.
	aFilenames, err := parameterToFilenames(*ain, ".csv")
	checkFatal(err, "Error parsing parameter ain")

	bFilenames, err := parameterToFilenames(*bin, ".csv")
	checkFatal(err, "Error parsing parameter bin")

	// Calculate unions of a and b and afterwards the difference a\b.
	var a snp.SNPs = make(map[snp.SNP]bool)
	err = unionWithFiles(a, aFilenames)
	checkFatal(err, "Error calculating the union of files for parameter ain")

	var b snp.SNPs = make(map[snp.SNP]bool)
	err = unionWithFiles(b, bFilenames)
	checkFatal(err, "Error calculating the union of files for parameter bin")

	a.Difference(b)

	// Output SNPs.
	if *out != "" {
		err := a.WriteCSV(*out)
		checkFatal(err, "Error writing SNPs to CSV output file")
	} else {
		for snp, _ := range a {
			fmt.Print(snp.String())
		}
	}
}
