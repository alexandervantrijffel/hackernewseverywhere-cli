package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"github.com/alexandervantrijffel/goutil/errorcheck"
	"github.com/alexandervantrijffel/goutil/logging"
	"github.com/alexandervantrijffel/hackernewseverywhere-cli/pkg/mergemp3"
	"github.com/alexandervantrijffel/hackernewseverywhere-cli/pkg/ssmltext"
	texttospeechpb "google.golang.org/genproto/googleapis/cloud/texttospeech/v1"
)

func main() {
	logging.InitWith("hackernewseverywhere-cli", false)
	ctx := context.Background()
	client, err := texttospeech.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	content, _ := ioutil.ReadAll(os.Stdin)
	if len(content) == 0 {
		log.Fatal("No content! Please pipe content to me")
		return
	}
	chunks, err := ssmltext.MakeChunks(string(content), 5000)
	errorcheck.CheckLogFatal(err, "No content to synthesize, please pipe text to me")
	if len(chunks) == 1 {
		SynthesizeSsmlToFile(client, ctx, chunks[0], "output.mp3")
		return
	}
	var sourceFiles []string
	for i, c := range chunks {
		src := strconv.Itoa(i) + ".mp3"
		SynthesizeSsmlToFile(client, ctx, c, src)
		sourceFiles = append(sourceFiles, src)
	}
	mergemp3.Merge("output.mp3", sourceFiles, true, false)
	for _, s := range sourceFiles {
		os.Remove(s)
	}
}

func SynthesizeSsmlToFile(client *texttospeech.Client, ctx context.Context, ssml, destinationFile string) {
	// Perform the text-to-speech request on the text input with the selected
	// voice parameters and audio file type.
	req := texttospeechpb.SynthesizeSpeechRequest{
		Input: &texttospeechpb.SynthesisInput{
			InputSource: &texttospeechpb.SynthesisInput_Ssml{Ssml: ssml},
		},

		// Set the text input to be synthesized.
		// Input: &texttospeechpb.SynthesisInput{
		// 	InputSource: &texttospeechpb.SynthesisInput_Text{Text: string(content)},

		// Build the voice request, select the language code ("en-US") and the SSML
		// voice gender ("neutral").
		Voice: &texttospeechpb.VoiceSelectionParams{
			LanguageCode: "en-US",
			// Wavenet male voice:         "en-US-Wavenet-D",
			// Wavenet female voice: en-US-Wavenet-C
			// Standard voice: en-US-Standard-B
			Name: "en-US-Wavenet-D",
		},
		AudioConfig: &texttospeechpb.AudioConfig{
			//AudioEncoding: texttospeechpb.AudioEncoding_MP3,
			Pitch: -6.00,
			SpeakingRate: 1.00,
			AudioEncoding:  texttospeechpb.AudioEncoding_LINEAR16,
		},
	}

	resp, err := client.SynthesizeSpeech(ctx, &req)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(destinationFile, resp.AudioContent, 0644)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Audio content written to file: %v\n", destinationFile)
}
