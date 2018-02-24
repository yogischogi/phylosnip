# Phylosnip

Phylosnip provides basic set operations for Y-chromosome SNP mutations.


## Examples

### Extract SNPs from Family Tree DNA, YFull or VCF files

phylosnip filterftdna -in=ftdna.csv -out=ftdna-novels.csv -mutationsonly=true -isoggdb=snps_hg38.csv -novelsonly=true

phylosnip filteryfull -in=yfull.csv -out=yfull-amb.csv -quality=ambiguous

phylosnip filtervcf -in=000.vcf -out=000.csv


### Set operations

phylosnip union -in=01.csv,02.csv -out=result.csv

phylosnip interesction -in=01.csv,02.csv -out=result.csv

phylosnip difference -ain=01.csv -bin=02.csv -out=result.csv


## Lookup SNPs in ISOGG database

phylosnip lookup -in=00.csv -isoggdb=snps_hg38.csv

phylosnip lookup -in=indir -out=outdir -isoggdb=snps_hg38.csv


## Documentation

* [Source Code](http://godoc.org/github.com/yogischogi/phylosnip)

