package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/yogischogi/phylosnip/snp"
)

// FilterYFull filters YFull novel SNP files for different SNP qualities.
// YFull provides CSV files that contain only novel SNPs. These SNPs come
// in three different qualities: best, acceptable and ambiguous.
// cmdLine: command line parameters without the subcommand.
func FilterYFull(cmdLine []string) {
	flags := flag.NewFlagSet("", flag.ContinueOnError)
	var (
		in      = flags.String("in", "", "Input file in FTDNA CSV format.")
		out     = flags.String("out", "", "Output file for list of SNPs in CSV format.")
		quality = flags.String("quality", "acceptable", "Minimum quality for the output.")
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

	if *quality != "best" && *quality != "acceptable" && *quality != "ambiguous" {
		fmt.Printf("Parameter quality must be best, acceptable or ambiguous.\n")
		os.Exit(1)
	}

	inNames, outNames, err := inToOutFilenames(*in, ".csv", *out, ".csv")
	checkFatal(err, "Error converting filenames from parameter in to out")
	for i, _ := range inNames {
		snps, err := snp.ReadYFull(inNames[i], *quality)
		checkFatal(err, "Error reading YFull CSV file")
		if outNames[i] != "" {
			// Write to file.
			err := snps.WriteCSV(outNames[i])
			checkFatal(err, "Error writing to CSV file")
		} else {
			// Write to stdout.
			for snp := range snps {
				os.Stdout.WriteString(snp.String())
			}
		}
	}
}
