package youtube

import (
	"context"
	"github.com/Genekoh/VikasBot/pkg/env"
	ytdl "github.com/kkdai/youtube/v2"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func QueryVideos(query string, maxResults int) ([]*youtube.SearchResult, error) {
	environ, err := env.GetEnviron("")
	apiKey := environ["YOUTUBE_API_KEY"]
	if apiKey == "" || err != nil {
		return nil, err
	}

	ctx := context.Background()
	service, err := youtube.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}

	call := service.Search.List([]string{"snippet"}).Q(query).MaxResults(int64(maxResults))
	response, err := call.Do()
	if err != nil {
		return nil, nil
	}

	return response.Items, nil
}

func DownloadVideo(vidId string) (string, error) {
	client := ytdl.Client{}

	video, err := client.GetVideo(vidId)
	if err != nil {
		return "", err
	}

	formats := video.Formats.WithAudioChannels()
	stream, err := client.GetStreamURL(video, &formats[0])
	if err != nil {
		return "", err
	}

	return stream, nil
}
