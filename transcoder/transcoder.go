package transcoder

import (
	"fmt"
	"github.com/asticode/go-astisub"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
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

func transcodeVideo(inputFile, outputFolder, chunkDuration, videoCodec, videoScale, introFile string) error {
	log.Println("Début du transcodage en HLS...")
	log.Println("Transcodage de la vidéo...")

	// Create a temporary file to store the list of input files
	listFile, err := os.CreateTemp("", "ffmpeg_list_*.txt")
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer os.Remove(listFile.Name())

	inputFile = strings.ReplaceAll(inputFile, "'", "'\\''")

	// Write the list of input files to the temporary file
	_, err = listFile.WriteString("file '" + introFile + "'\n")
	if err != nil {
		return err
	}
	_, err = listFile.WriteString("file '" + inputFile + "'\n")
	if err != nil {
		return err
	}

	// Close the temporary file
	err = listFile.Close()
	if err != nil {
		return fmt.Errorf("failed to close temporary file: %w", err)
	}

	// Initialize common ffmpeg command arguments
	ffmpegArgs := []string{
		"-f", "concat",
		"-safe", "0",
		"-i", listFile.Name(),
		"-map", "0:v:0", // Sélectionnez seulement la première piste vidéo
	}

	if videoCodec == "h264" {
		// If the original video is h264, copy the codec for the main video
		log.Println("La vidéo est déjà encodée en h264, copie du codec pour la vidéo principale...")
		ffmpegArgs = append(ffmpegArgs, "-c:v", "copy")
	} else {
		// Otherwise, transcode the main video
		log.Println("La vidéo n'est pas encodée en h264, transcodage de la vidéo principale...")
		ffmpegArgs = append(ffmpegArgs,
			"-vf", fmt.Sprintf("scale=%s,format=yuv420p", videoScale), // rescaling to 720p
			"-c:v", "libx264",
			"-profile:v", "main", // Using the Main profile
			"-preset", "veryfast",
			"-crf", "23",
			"-pix_fmt", "yuv420p",
		)
	}
	ffmpegArgs = append(ffmpegArgs,
		"-hls_time", chunkDuration,
		"-hls_playlist_type", "vod",
		"-hls_segment_filename", filepath.Join(outputFolder, "segment_%03d.ts"),
		"-hls_flags", "delete_segments",
		"-f", "hls", filepath.Join(outputFolder, "index.m3u8"),
	)
	cmd := exec.Command("ffmpeg", ffmpegArgs...)
	//cmd.Stdout = os.Stdout
	//cmd.Stderr = os.Stderr
	log.Println("Commande ffmpeg :", cmd.String())
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to execute command: %w", err)
	}
	log.Println("Vidéo extraite :", "index.m3u8")
	return nil
}

func extractAudioStreams(inputFile, outputFolder, chunkDuration string, audioStreams []string, introFile string) error {
	log.Println("Transcodage des pistes audio...")

	for _, stream := range audioStreams {
		outputFile := filepath.Join(outputFolder, fmt.Sprintf("audio_%s.m3u8", stream))
		cmd := exec.Command("ffmpeg",
			"-i", introFile,
			"-i", inputFile,
			"-filter_complex", "[0:a:0][1:"+stream+"]concat=n=2:v=0:a=1[outa]",
			"-map", "[outa]",
			"-c:a", "aac",
			"-b:a", "160k",
			"-ac", "2",
			"-hls_time", chunkDuration,
			"-hls_playlist_type", "vod",
			"-hls_segment_filename", filepath.Join(outputFolder, fmt.Sprintf("audio_%s_%%03d.ts", stream)),
			outputFile,
		)
		//cmd.Stdout = os.Stdout
		//cmd.Stderr = os.Stderr
		log.Println("Commande ffmpeg :", cmd.String())

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to execute command: %w", err)
		}
		log.Println("Piste audio extraite :", outputFile)
	}
	return nil
}

func convertSubtitleFormat(inputFile, outputFile string) error {
	cmd := exec.Command("ffmpeg",
		"-i", inputFile,
		outputFile,
	)
	//cmd.Stdout = os.Stdout
	//cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to convert subtitle format: %w", err)
	}
	return nil
}

func extractSubtitleStreams(inputFile, outputFolder string, subtitleStreams []string, introFile string) error {
	log.Println("Transcodage des pistes de sous-titres...")

	// Obtenir la durée de la vidéo "intro"
	introDuration, err := getVideoDuration(introFile)
	if err != nil {
		return fmt.Errorf("failed to get intro video duration: %w", err)
	}

	log.Println("Durée de la vidéo d'introduction :", introDuration)

	for _, stream := range subtitleStreams {
		// Convert DVD subtitle to SRT format
		srtOutputFile := filepath.Join(outputFolder, fmt.Sprintf("subtitle_%s.srt", stream))
		if err := convertSubtitleFormat(inputFile, srtOutputFile); err != nil {
			return fmt.Errorf("failed to convert subtitle to SRT: %w", err)
		}

		// Convert SRT subtitle to WebVTT format
		vttOutputFile := filepath.Join(outputFolder, fmt.Sprintf("subtitle_%s.vtt", stream))
		if err := convertSubtitleFormat(srtOutputFile, vttOutputFile); err != nil {
			return fmt.Errorf("failed to convert SRT subtitle to WebVTT: %w", err)
		}

		if err := os.Remove(srtOutputFile); err != nil {
			log.Printf("Failed to remove intermediate SRT file: %s\n", srtOutputFile)
		}

		if err = shiftSubtitleTimecodes(vttOutputFile, introDuration); err != nil {
			return fmt.Errorf("failed to shift subtitle timestamps: %w", err)
		}

		log.Println("Piste de sous-titres extraite :", vttOutputFile)
	}
	return nil
}

func getVideoDuration(videoFile string) (time.Duration, error) {
	cmd := exec.Command("ffprobe",
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		videoFile,
	)
	output, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("failed to execute ffprobe command: %w", err)
	}

	durationSec, err := strconv.ParseFloat(strings.TrimSpace(string(output)), 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse duration value: %w", err)
	}

	duration := time.Duration(durationSec * float64(time.Second))
	return duration, nil
}

func shiftSubtitleTimecodes(subtitleFile string, duration time.Duration) error {
	// Ouvrir le fichier de sous-titres SRT
	subs, err := astisub.OpenFile(subtitleFile)
	if err != nil {
		return fmt.Errorf("failed to open subtitle file: %w", err)
	}

	// Décaler les timecodes de la durée de la vidéo d'introduction
	subs.Add(duration)

	// Enregistrer les modifications dans le fichier SRT
	err = subs.Write(subtitleFile)
	if err != nil {
		return fmt.Errorf("failed to write modified subtitle file: %w", err)
	}

	return nil
}

func ProcessFileTranscode(inputFilePath, introPath, mediaID, outputFolder, chunkDuration, videoScale string) (TranscodeResponse, error) {
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
	if err := transcodeVideo(inputFilePath, outputFileFolder, chunkDuration, videoCodec, videoScale, introPath); err != nil {
		return TranscodeResponse{}, err
	}
	log.Println("Temps de transcodage de la vidéo :", time.Since(beforeTranscode))

	beforeAudio := time.Now()
	if err := extractAudioStreams(inputFilePath, outputFileFolder, chunkDuration, audioStreams, introPath); err != nil {
		return TranscodeResponse{}, err
	}
	log.Println("Temps de transcodage des pistes audio :", time.Since(beforeAudio))

	beforeSubtitle := time.Now()
	if err := extractSubtitleStreams(inputFilePath, outputFileFolder, subtitleStreams, introPath); err != nil {
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

	// Set folder permissions to 777
	if err := os.Chmod(outputFileFolder, 0777); err != nil {
		log.Println("Failed to set folder permissions to 777 :", err)
	}

	return response, nil
}

/*func main() {
	const (
		introFile     = "/home/nospy/Téléchargements/intro.mkv"
		inputFile     = "/home/nospy/Téléchargements/video.mkv"
		inputFileID   = "123456"
		outputFolder  = "/home/nospy/Téléchargements/media/"
		chunkDuration = "15"       // durée des segments en secondes
		videoScale    = "1280:720" // dimension de la vidéo
	)
	response, err := ProcessFileTranscode(inputFile, introFile, inputFileID, outputFolder, chunkDuration, videoScale)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(response)
}*/
