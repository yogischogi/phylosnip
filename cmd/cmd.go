// Package cmd provides a set of subcommands for the main phylosnip package.
package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// checkFatal checks for an error. In case the
// error has occurred, the function prints out a message
// and exits the program.
func checkFatal(err error, msg string) {
	if err != nil {
		fmt.Printf("%s, %v.\n", msg, err)
		os.Exit(1)
	}
}

// parameterToFilenames parses a command line parameter for filenames.
// The parameter containes a list of filenames separated by commas.
// If a filename is a directory parameterToFilenames returns all files
// within that directory that satisfy the given extension ext.
func parameterToFilenames(filesParameter string, ext string) (filenames []string, err error) {
	files := strings.Split(filesParameter, ",")
	for _, file := range files {
		fileInfo, err := os.Stat(file)
		switch {
		case err != nil:
			return filenames, errors.New(fmt.Sprintf("unknown file, %v\n", err))
		case fileInfo.IsDir():
			dirFiles, err := namesWithExt(file, ext)
			if err != nil {
				return filenames, errors.New(fmt.Sprintf("reading directory %s, %v\n", file, err))
			}
			for _, name := range dirFiles {
				name = filepath.Join(file, name)
				filenames = append(filenames, name)
			}
		case strings.HasSuffix(strings.ToLower(file), ext):
			filenames = append(filenames, file)
		}
	}
	return filenames, nil
}

// namesWithExt returns the names of all files in a directory
// ending with the extension ext.
// If there are no matching files in the directory
// an empty slice is returned.
func namesWithExt(dirName string, ext string) (filenames []string, err error) {
	filenames = make([]string, 0, 100)
	dir, err := os.Open(dirName)
	if err != nil {
		return filenames, errors.New(fmt.Sprintf("could not open directory, %s\n", err))
	}
	defer dir.Close()

	files, err := dir.Readdirnames(0)
	if err != nil {
		return files, errors.New(fmt.Sprintf("could not read files from directory, %s\n", err))
	}
	for _, filename := range files {
		if filepath.Ext(filename) == ext {
			filenames = append(filenames, filename)
		}
	}
	return filenames, err
}

// inToOutFilenames takes two files or directories as input and
// maps all input names to corresponding output names.
// If out is an empty string, all input files are mapped to empty strings.
func inToOutFilenames(in, inExt, out, outExt string) (inNames, outNames []string, err error) {
	inInfo, err := os.Stat(in)
	if err != nil {
		return inNames, outNames, errors.New(fmt.Sprintf("unknown file, %v\n", err))
	}

	// Simplest case: in file is regular.
	if inInfo.Mode().IsRegular() {
		inNames = append(inNames, in)
		outNames = append(outNames, out)
		return inNames, outNames, nil
	}

	// in file is a directory.
	if inInfo.IsDir() {
		// Read input files.
		names, err := namesWithExt(in, inExt)
		if err != nil {
			return inNames, outNames, errors.New(fmt.Sprintf("could not read from directory in, %v\n", err))
		}
		// Create output filenames.
		for _, name := range names {
			inName := filepath.Join(in, name)
			inNames = append(inNames, inName)
			base := filepath.Base(name)
			base = strings.TrimSuffix(base, inExt)
			outName := base + outExt
			outName = filepath.Join(out, outName)
			outNames = append(outNames, outName)
		}
		return inNames, outNames, nil
	}

	// If out file is an empty string, all input files are mapped to empty strings.
	if out == "" {
		if inInfo.Mode().IsRegular() {
			inNames = append(inNames, in)
			outNames = append(outNames, "")
			return inNames, outNames, nil
		}
		if inInfo.IsDir() {
			inNames, err = namesWithExt(in, inExt)
			if err != nil {
				return inNames, outNames, errors.New(fmt.Sprintf("could not read from directory in, %v\n", err))
			}
			outNames = make([]string, len(inNames), len(inNames))
			return inNames, outNames, nil
		}
	}

	// Everything else is forbidden.
	return inNames, outNames, errors.New("in and out must be both files or directories")
}
