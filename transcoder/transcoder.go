package transcoder

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type AudioTranscodeResponse struct {
	AudioIndex string `json:"audio_index"`
}

type SubtitleTranscodeResponse struct {
	SubtitleIndex string `json:"subtitle_index"`
}

type TranscodeResponse struct {
	VideoIndex string                      `json:"video_index"`
	Audios     []AudioTranscodeResponse    `json:"audios"`
	Subtitles  []SubtitleTranscodeResponse `json:"subtitles"`
}

func prepareOutputFolder(outputFolder string) error {
	if err := os.MkdirAll(outputFolder, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	files, err := os.ReadDir(outputFolder)
	if err != nil {
		return err
	}

	for _, f := range files {
		if err := os.Remove(filepath.Join(outputFolder, f.Name())); err != nil {
			return err
		}
	}

	return nil
}

func extractStreamsInfo(inputFile string) (audioStreams, subtitleStreams []string, videoCodec string, err error) {
	log.Println("Récupération des informations sur les pistes audio et sous-titres...")
	cmd := exec.Command("ffprobe",
		"-v", "error",
		"-show_entries", "stream=index,codec_name,codec_type",
		"-of", "csv=p=0",
		inputFile,
	)
	output, err := cmd.Output()
	if err != nil {
		return nil, nil, "", fmt.Errorf("failed to execute command: %w", err)
	}

	ffprobeOutput := strings.Split(strings.TrimSpace(string(output)), "\n")

	for _, line := range ffprobeOutput {
		fields := strings.Split(line, ",")
		streamIndex, codecName, codecType := fields[0], fields[1], fields[2]

		switch codecType {
		case "audio":
			audioStreams = append(audioStreams, streamIndex)
		case "subtitle":
			subtitleStreams = append(subtitleStreams, streamIndex)
		case "video":
			videoCodec = codecName
		}
	}

	log.Println("Pistes audio trouvées :", audioStreams)
	log.Println("Pistes de sous-titres trouvées :", subtitleStreams)
	log.Println("Codec vidéo :", videoCodec)

	return audioStreams, subtitleStreams, videoCodec, nil
}

func transcodeVideo(inputFile, outputFolder, chunkDuration, videoCodec, videoScale string) error {
	log.Println("Début du transcodage en HLS...")
	log.Println("Transcodage de la vidéo...")

	// Initialize common ffmpeg command arguments
	ffmpegArgs := []string{
		"-i", inputFile,
		"-map", "0:0", // Sélectionnez seulement la première piste vidéo
		"-hls_time", chunkDuration,
		"-hls_playlist_type", "vod",
		"-hls_segment_filename", filepath.Join(outputFolder, "segment_%03d.ts"),
		"-hls_flags", "delete_segments",
		"-f", "hls", filepath.Join(outputFolder, "index.m3u8"),
	}

	if videoCodec == "h264" {
		// If the original video is h264, copy the codec
		log.Println("La vidéo est déjà encodée en h264, copie du codec...")
		ffmpegArgs = append(ffmpegArgs, "-c:v", "copy")
	} else {
		// Otherwise, transcode the video
		log.Println("La vidéo n'est pas encodée en h264, transcodage...")
		ffmpegArgs = append(ffmpegArgs,
			"-vf", fmt.Sprintf("scale=%s,format=yuv420p", videoScale), // rescaling to 720p
			"-c:v", "libx264",
			"-profile:v", "main", // Using the Main profile
			"-preset", "veryfast",
			"-crf", "23",
			"-pix_fmt", "yuv420p",
		)
	}
	cmd := exec.Command("ffmpeg", ffmpegArgs...)
	log.Println("Commande ffmpeg :", cmd.String())
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to execute command: %w", err)
	}
	log.Println("Vidéo extraite :", "index.m3u8")
	return nil
}

func extractAudioStreams(inputFile, outputFolder, chunkDuration string, audioStreams []string) error {
	log.Println("Transcodage des pistes audio...")
	for _, stream := range audioStreams {
		outputFile := filepath.Join(outputFolder, fmt.Sprintf("audio_%s.m3u8", stream))
		cmd := exec.Command("ffmpeg",
			"-i", inputFile,
			"-map", "0:"+stream,
			"-c:a", "aac",
			"-b:a", "160k",
			"-ac", "2",
			"-hls_time", chunkDuration,
			"-hls_playlist_type", "vod",
			"-hls_segment_filename", filepath.Join(outputFolder, fmt.Sprintf("audio_%s_%%03d.ts", stream)),
			outputFile,
		)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to execute command: %w", err)
		}
		log.Println("Piste audio extraite :", outputFile)
	}
	return nil
}

func extractSubtitleStreams(inputFile, outputFolder string, subtitleStreams []string) error {
	log.Println("Transcodage des pistes de sous-titres...")
	for _, stream := range subtitleStreams {
		outputFile := filepath.Join(outputFolder, fmt.Sprintf("subtitle_%s.vtt", stream))
		cmd := exec.Command("ffmpeg",
			"-i", inputFile,
			"-map", "0:"+stream,
			outputFile,
		)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to execute command: %w", err)
		}
		log.Println("Piste de sous-titres extraite :", outputFile)
	}
	return nil
}

func ProcessFileTranscode(inputFilePath, mediaID, outputFolder, chunkDuration, videoScale string) (TranscodeResponse, error) {
	start := time.Now()
	log.Println("Début du transcodage du fichier :", inputFilePath)

	outputFileFolder := filepath.Join(outputFolder, mediaID)
	if err := prepareOutputFolder(outputFileFolder); err != nil {
		return TranscodeResponse{}, err
	}

	audioStreams, subtitleStreams, videoCodec, err := extractStreamsInfo(inputFilePath)
	if err != nil {
		return TranscodeResponse{}, err
	}

	beforeTranscode := time.Now()
	if err := transcodeVideo(inputFilePath, outputFileFolder, chunkDuration, videoCodec, videoScale); err != nil {
		return TranscodeResponse{}, err
	}
	log.Println("Temps de transcodage de la vidéo :", time.Since(beforeTranscode))

	beforeAudio := time.Now()
	if err := extractAudioStreams(inputFilePath, outputFileFolder, chunkDuration, audioStreams); err != nil {
		return TranscodeResponse{}, err
	}
	log.Println("Temps de transcodage des pistes audio :", time.Since(beforeAudio))

	beforeSubtitle := time.Now()
	if err := extractSubtitleStreams(inputFilePath, outputFileFolder, subtitleStreams); err != nil {
		return TranscodeResponse{}, err
	}
	log.Println("Temps de transcodage des pistes de sous-titres :", time.Since(beforeSubtitle))

	log.Println("Transcodage terminé. Fichiers HLS générés dans :", outputFileFolder)
	response := TranscodeResponse{
		VideoIndex: "index.m3u8",
	}
	for _, stream := range audioStreams {
		response.Audios = append(response.Audios, AudioTranscodeResponse{
			AudioIndex: fmt.Sprintf("audio_%s.m3u8", stream),
		})
	}
	for _, stream := range subtitleStreams {
		response.Subtitles = append(response.Subtitles, SubtitleTranscodeResponse{
			SubtitleIndex: fmt.Sprintf("subtitle_%s.vtt", stream),
		})
	}
	log.Println("Temps de transcodage :", time.Since(start))
	return response, nil
}

/*func main() {
	const (
		inputFile     = "/home/nospy/Téléchargements/Mashle.S01E07.VOSTFR.1080p.WEBRiP.x265-KAF.mkv"
		inputFileID   = "123456"
		outputFolder  = "/home/nospy/Téléchargements/media/"
		chunkDuration = "15"       // durée des segments en secondes
		videoScale    = "1280:720" // dimension de la vidéo
	)
	response, err := ProcessFileTranscode(inputFile, inputFileID, outputFolder, chunkDuration, videoScale)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(response)
}*/
