package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
	"github.com/gordonklaus/portaudio"
)

func main() {
	portaudio.Initialize()
	defer portaudio.Terminate()

	stream, err := portaudio.OpenDefaultStream(&portaudio.StreamParameters{
		Input: portaudio.StreamDeviceParameters{Channels: 1},
	}, nil, 16000, 1024, processAudio)
	if err != nil {
		log.Fatal(err)
	}

	err = stream.Start()
	if err != nil {
		log.Fatal(err)
	}
	defer stream.Stop()

	fmt.Println("Listening... Press Ctrl+C to stop.")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}

func processAudio(in []int32) {
	// Process audio data here
	saveAudioDataToWAV(in)
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
