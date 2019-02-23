package mergemp3

import (
	"fmt"
	"io"
	"os"

	"github.com/dmulholland/mp3lib"
)

// Create a new file at the specified output path containing the merged
// contents of the list of input files.
func Merge(outpath string, inpaths []string, force, tag bool) {
	var totalFrames uint32
	var totalBytes uint32
	var totalFiles int
	var firstBitRate int
	var isVBR bool

	// Only overwrite an existing file if the --force flag has been used.
	if _, err := os.Stat(outpath); err == nil {
		if !force {
			fmt.Fprintf(
				os.Stderr,
				"Error: the file '%v' already exists.\n", outpath)
			os.Exit(1)
		}
	}

	// If the list of input files includes the output file we'll end up in an
	// infinite loop.
	for _, filepath := range inpaths {
		if filepath == outpath {
			fmt.Fprintln(
				os.Stderr,
				"Error: the list of input files includes the output file.")
			os.Exit(1)
		}
	}

	// Create the output file.
	outfile, err := os.Create(outpath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Loop over the input files and append their MP3 frames to the output
	// file.
	for _, inpath := range inpaths {

		fmt.Println("+", inpath)

		infile, err := os.Open(inpath)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		isFirstFrame := true

		for {

			// Read the next frame from the input file.
			frame := mp3lib.NextFrame(infile)
			if frame == nil {
				break
			}

			// Skip the first frame if it's a VBR header.
			if isFirstFrame {
				isFirstFrame = false
				if mp3lib.IsXingHeader(frame) || mp3lib.IsVbriHeader(frame) {
					continue
				}
			}

			// If we detect more than one bitrate we'll need to add a VBR
			// header to the output file.
			if firstBitRate == 0 {
				firstBitRate = frame.BitRate
			} else if frame.BitRate != firstBitRate {
				isVBR = true
			}

			// Write the frame to the output file.
			_, err := outfile.Write(frame.RawBytes)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			totalFrames += 1
			totalBytes += uint32(len(frame.RawBytes))
		}

		infile.Close()
		totalFiles += 1
	}

	outfile.Close()

	// If we detected multiple bitrates, prepend a VBR header to the file.
	if isVBR {
		fmt.Println("• Multiple bitrates detected. Adding VBR header.")
		addXingHeader(outpath, totalFrames, totalBytes)
	}

	// Copy the ID3v2 tag from the first input file if requested. Order of
	// operations is important here. The ID3 tag must be the first item in
	// the file - in particular, it must come *before* any VBR header.
	if tag {
		fmt.Println("• Adding ID3 tag.")
		addID3v2Tag(outpath, inpaths[0])
	}

	// Print a count of the number of files merged.
	fmt.Printf("• %v files merged.\n", totalFiles)
}

// Prepend an Xing VBR header to the specified MP3 file.
func addXingHeader(filepath string, totalFrames, totalBytes uint32) {

	outputFile, err := os.Create(filepath + ".mp3cat.tmp")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	inputFile, err := os.Open(filepath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	xingHeader := mp3lib.NewXingHeader(totalFrames, totalBytes)

	_, err = outputFile.Write(xingHeader.RawBytes)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	_, err = io.Copy(outputFile, inputFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	outputFile.Close()
	inputFile.Close()

	err = os.Remove(filepath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	err = os.Rename(filepath+".mp3cat.tmp", filepath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// Prepend an ID3v2 tag to the MP3 file at mp3Path, copying from tagPath.
func addID3v2Tag(mp3Path, tagPath string) {

	tagFile, err := os.Open(tagPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	id3tag := mp3lib.NextID3v2Tag(tagFile)
	tagFile.Close()

	if id3tag != nil {
		outputFile, err := os.Create(mp3Path + ".mp3cat.tmp")
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		inputFile, err := os.Open(mp3Path)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		_, err = outputFile.Write(id3tag.RawBytes)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		_, err = io.Copy(outputFile, inputFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		outputFile.Close()
		inputFile.Close()

		err = os.Remove(mp3Path)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		err = os.Rename(mp3Path+".mp3cat.tmp", mp3Path)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}
