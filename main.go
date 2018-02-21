// Package main is a meta package for a bundle of SNP related programs.
package main

import (
	"fmt"
	"os"

	"github.com/yogischogi/phylosnip/cmd"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: phylosnip <subcommand> <options>\n" +
			"Current subcommands:\n" +
			"    filtervcf\n" +
			"        extracts SNPs from a VCF file.\n" +
			"    filterftdna\n" +
			"        extracts SNPs from FTDNA CSV file.\n" +
			"    filteryfull\n" +
			"        extracts SNPs from YFull novels SNP file.\n" +
			"    union\n" +
			"        unites SNPs from CSV files.\n" +
			"    intersection\n" +
			"        calculates the intersection of SNPs from CSV files.\n" +
			"    difference\n" +
			"        calculates the difference of SNPs from CSV files.\n" +
			"    lookup\n" +
			"        adds ISOGG data base information to SNP CSV files.\n")
		os.Exit(1)
	}

	// Relay control to subcommands.
	switch os.Args[1] {
	case "filtervcf":
		cmd.FilterVCF(os.Args[2:])
	case "filterftdna":
		cmd.FilterFTDNA(os.Args[2:])
	case "filteryfull":
		cmd.FilterYFull(os.Args[2:])
	case "union":
		cmd.Union(os.Args[2:])
	case "intersection":
		cmd.Intersection(os.Args[2:])
	case "difference":
		cmd.Difference(os.Args[2:])
	case "lookup":
		cmd.Lookup(os.Args[2:])
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
	}
}
