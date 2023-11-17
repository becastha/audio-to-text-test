package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
	"github.com/gordonklaus/portaudio"
)

var (
	recording       = false
	stream          *portaudio.Stream
	recordedSamples []int32
)

func main() {
	portaudio.Initialize()
	defer portaudio.Terminate()

	http.HandleFunc("/toggleRecording", toggleRecordingHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func toggleRecordingHandler(w http.ResponseWriter, r *http.Request) {
	if recording {
		stopRecording()
		fmt.Fprint(w, "Stopped recording.")
	} else {
		err := startRecording()
		if err != nil {
			log.Println("Error starting recording:", err)
			fmt.Fprint(w, "Error starting recording.")
			return
		}
		fmt.Fprint(w, "Started recording.")
	}
}

func startRecording() error {
	recordedSamples = nil // Reset recorded samples
	var err error
	stream, err = portaudio.OpenDefaultStream(&portaudio.StreamParameters{
		Input: portaudio.StreamDeviceParameters{Channels: 1},
	}, nil, 16000, 1024, processAudio)
	if err != nil {
		return err
	}
	err = stream.Start()
	if err != nil {
		return err
	}
	recording = true
	return nil
}

func stopRecording() {
	if stream != nil {
		stream.Stop()
		stream.Close()
	}
	saveAudioDataToWAV(recordedSamples)
	recording = false
}

func processAudio(in []int32) {
	// Append samples to recordedSamples slice
	recordedSamples = append(recordedSamples, in...)
}

func saveAudioDataToWAV(data []int32) {
	fileName := fmt.Sprintf("recording-%s.wav", time.Now().Format("20060102150405"))
	file, err := os.Create(fileName)
	if err != nil {
		log.Println("Error creating WAV file:", err)
		return
	}
	defer file.Close()

	enc := wav.NewEncoder(file, 16000, 16, 1, 1)

	// Convert int32 slice to int slice
	intData := make([]int, len(data))
	for i, d := range data {
		intData[i] = int(d)
	}

	buf := &audio.IntBuffer{
		Format: &audio.Format{
			NumChannels: 1,
			SampleRate:  16000,
		},
		Data: intData,
	}

	err = enc.Write(buf)
	if err != nil {
		log.Println("Error writing audio samples:", err)
		return
	}
}
