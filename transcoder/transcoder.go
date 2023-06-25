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

func extractStreamsInfo(inputFile string) (audioStreams, subtitleStreams []string, videoCodec string, aspectRatio string, err error) {
	log.Println("Récupération des informations sur les pistes audio et sous-titres...")
	cmd := exec.Command("ffprobe",
		"-v", "error",
		"-show_entries", "stream=index,codec_name,codec_type,display_aspect_ratio",
		"-of", "csv=p=0",
		inputFile,
	)
	output, err := cmd.Output()
	if err != nil {
		return nil, nil, "", "", fmt.Errorf("failed to execute command: %w", err)
	}

	ffprobeOutput := strings.Split(strings.TrimSpace(string(output)), "\n")

	for _, line := range ffprobeOutput {
		fields := strings.Split(line, ",")
		streamIndex, codecName, codecType := fields[0], fields[1], fields[2]

		switch codecType {
		case "audio":
			audioStreams = append(audioStreams, streamIndex)
		case "subtitle":
			log.Println("Piste de sous-titres trouvée :", streamIndex, codecName)
			if codecName != "dvd_subtitle" && codecName != "hdmv_pgs_subtitle" {
				subtitleStreams = append(subtitleStreams, streamIndex)
			}
		case "video":
			if videoCodec != "" {
				continue
			}
			videoCodec = codecName
			aspectRatio = fields[3]
		}
	}

	log.Println("Pistes audio trouvées :", audioStreams)
	log.Println("Pistes de sous-titres trouvées :", subtitleStreams)
	log.Println("Codec vidéo :", videoCodec)

	return audioStreams, subtitleStreams, videoCodec, aspectRatio, nil
}

func transcodeVideo(inputFile, outputFolder, chunkDuration, videoScale, introFile string) error {
	log.Println("Début du transcodage en HLS...")
	log.Println("Transcodage de la vidéo...")

	// Initialize common ffmpeg command arguments
	ffmpegArgs := []string{
		"-i", introFile,
		"-i", inputFile,
		"-filter_complex", fmt.Sprintf("[0:v:0]scale=%s,format=yuv420p,setsar=sar=1/1[v0]; [1:v:0]scale=%s,format=yuv420p,setsar=sar=1/1[v1]; [v0][v1]concat=n=2:v=1[outv]", videoScale, videoScale),
		"-map", "[outv]",
		"-c:v", "libx264",
		"-profile:v", "high", // Using the Main profile
		"-preset", "veryfast",
		"-crf", "25",
		"-pix_fmt", "yuv420p",
		"-hls_time", chunkDuration,
		"-hls_playlist_type", "vod",
		"-hls_segment_filename", filepath.Join(outputFolder, "segment_%03d.ts"),
		"-hls_flags", "delete_segments",
		"-f", "hls", filepath.Join(outputFolder, "index.m3u8"),
	}

	cmd := exec.Command("ffmpeg", ffmpegArgs...)
	//cmd.Stdout = os.Stdout
	//cmd.Stderr = os.Stderr
	log.Println("Commande ffmpeg :", cmd.String())
	err := cmd.Run()
	if err != nil {
		cmd = exec.Command("ffmpeg", ffmpegArgs...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
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
			if err != nil {
				cmd = exec.Command("ffmpeg",
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
				cmd.Stderr = os.Stderr
				cmd.Stdout = os.Stdout
				err = cmd.Run()
				return fmt.Errorf("failed to execute command: %w", err)
			}
		}
		log.Println("Piste audio extraite :", outputFile)
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
		outputFile := filepath.Join(outputFolder, fmt.Sprintf("subtitle_%s.vtt", stream))
		cmd := exec.Command("ffmpeg",
			"-i", inputFile,
			"-map", "0:"+stream,
			outputFile,
		)
		//cmd.Stdout = os.Stdout
		//cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			cmd := exec.Command("ffmpeg",
				"-i", inputFile,
				"-map", "0:"+stream,
				outputFile,
			)
			cmd.Stderr = os.Stderr
			cmd.Stdout = os.Stdout
			err = cmd.Run()
			return fmt.Errorf("failed to execute command: %w", err)
		}
		if err = shiftSubtitleTimecodes(outputFile, introDuration); err != nil {
			log.Printf("failed to shift subtitle timestamps: %v", err)
		}

		log.Println("Piste de sous-titres extraite :", outputFile)
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

func ProcessFileTranscode(inputFilePath, introPath, intro219Path, mediaID, outputFolder, chunkDuration, videoScale, videoScale219 string) (TranscodeResponse, error) {
	start := time.Now()
	log.Println("Début du transcodage du fichier :", inputFilePath)

	outputFileFolder := filepath.Join(outputFolder, mediaID)
	if err := prepareOutputFolder(outputFileFolder); err != nil {
		return TranscodeResponse{}, err
	}

	audioStreams, subtitleStreams, _, aspectRatio, err := extractStreamsInfo(inputFilePath)
	if err != nil {
		return TranscodeResponse{}, err
	}

	beforeTranscode := time.Now()
	ratioX, err := strconv.ParseFloat(strings.Split(aspectRatio, ":")[0], 64)
	if err != nil {
		log.Println("Erreur lors de la récupération du ratio de la vidéo :", err)
		ratioX = 16
	}
	ratioY, err := strconv.ParseFloat(strings.Split(aspectRatio, ":")[1], 64)
	if err != nil {
		log.Println("Erreur lors de la récupération du ratio de la vidéo :", err)
		ratioY = 9
	}

	if ratioX/ratioY > 1.8 {
		log.Println("La vidéo est au format 21:9")
		if err := transcodeVideo(inputFilePath, outputFileFolder, chunkDuration, videoScale219, intro219Path); err != nil {
			os.RemoveAll(outputFileFolder)
			return TranscodeResponse{}, err
		}
	} else {
		log.Println("La vidéo est au format 16:9")
		if err := transcodeVideo(inputFilePath, outputFileFolder, chunkDuration, videoScale, introPath); err != nil {
			os.RemoveAll(outputFileFolder)
			return TranscodeResponse{}, err
		}
	}
	log.Println("Temps de transcodage de la vidéo :", time.Since(beforeTranscode))

	beforeAudio := time.Now()
	if err := extractAudioStreams(inputFilePath, outputFileFolder, chunkDuration, audioStreams, introPath); err != nil {
		os.RemoveAll(outputFileFolder)
		return TranscodeResponse{}, err
	}
	log.Println("Temps de transcodage des pistes audio :", time.Since(beforeAudio))

	beforeSubtitle := time.Now()
	if err := extractSubtitleStreams(inputFilePath, outputFileFolder, subtitleStreams, introPath); err != nil {
		os.RemoveAll(outputFileFolder)
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
		introFile     = "/home/nospy/Projets/bingemate/media-indexer/assets/intro.mkv"
		introFile219  = "/home/nospy/Projets/bingemate/media-indexer/assets/intro_21-9.mp4"
		inputFile     = "/media/nospy/Data/Encodage/Encoded/Star Wars - Episode IV - A New Hope - 1977.mkv"
		inputFileID   = "123456"
		outputFolder  = "/home/nospy/Téléchargements/media/"
		chunkDuration = "15"       // durée des segments en secondes
		videoScale    = "1280:720" // dimension de la vidéo
		videoScale219 = "1920:816" // dimension de la vidéo
	)
	response, err := ProcessFileTranscode(inputFile, introFile, introFile219, inputFileID, outputFolder, chunkDuration, videoScale, videoScale219)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(response)
}*/
