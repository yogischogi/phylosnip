package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/yogischogi/phylosnip/snp"
)

// Filter performs various filter operations on SNP CSV files.
// cmdLine: command line parameters without the subcommand.
func Filter(cmdLine []string) {
	flags := flag.NewFlagSet("", flag.ContinueOnError)
	var (
		in      = flags.String("in", "", "List of VCF files or directory.")
		out     = flags.String("out", "", "Output file for list of SNPs in CSV format.")
		bed     = flags.String("bed", "", "Input BED file.")
		exclude = flags.String("exclude", "", "Input file with list of SNPs that should be excluded.")
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

	var err error
	var bedRegs snp.BEDRegions
	if *bed != "" {
		bedRegs, err = snp.ReadBED(*bed)
		checkFatal(err, "Error reading BED file")
	}
	var ex snp.SNPs
	if *exclude != "" {
		ex, err = snp.ReadCSV(*exclude)
		checkFatal(err, "Error reading excludes file")
	}
	inFiles, outFiles, err := inToOutFilenames(*in, ".csv", *out, ".csv")
	checkFatal(err, "Error converting filenames from parameter in to out")

	for i, _ := range inFiles {
		snps, err := snp.ReadCSV(inFiles[i])
		checkFatal(err, "Error reading input CSV file")

		// Peform filter operations.
		if *exclude != "" {
			snps.Difference(ex)
		}
		if *bed != "" {
			for s, _ := range snps {
				if bedRegs.Includes(s.Pos) == false {
					delete(snps, s)
				}
			}
		}

		// Write to file or stdout.
		if outFiles[i] != "" {
			err := snps.WriteCSV(outFiles[i])
			checkFatal(err, "Error writing to CSV file")
		} else {
			for snp := range snps {
				os.Stdout.WriteString(snp.String())
			}
		}
	}
}
