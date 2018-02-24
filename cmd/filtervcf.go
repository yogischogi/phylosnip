package cmd

import (
	"flag"
	"fmt"
	"math"
	"os"

	"github.com/yogischogi/phylosnip/snp"
)

// FilterVCF filters VCF files for SNPs.
// cmdLine: command line parameters without the subcommand.
func FilterVCF(cmdLine []string) {
	flags := flag.NewFlagSet("", flag.ContinueOnError)
	var (
		in      = flags.String("in", "", "List of VCF files or directory.")
		out     = flags.String("out", "", "Output file for list of SNPs in CSV format.")
		quality = flags.Float64("quality", math.Inf(1), "Quality of SNP entry in VCF file.")
		reads   = flags.Int("reads", 3, "Minimum of total reads.")
		ratio   = flags.Int("ratio", 3, "Minimum ratio of ALT to REF reads.")
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

	inNames, outNames, err := inToOutFilenames(*in, ".vcf", *out, ".csv")
	checkFatal(err, "Error converting filenames from parameter in to out")
	for i, _ := range inNames {
		snps, err := snp.ReadVCF(inNames[i], *quality, *reads, *ratio)
		checkFatal(err, "Error reading VCF file")
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
