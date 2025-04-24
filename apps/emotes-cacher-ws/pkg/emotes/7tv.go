package emotes

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/samber/lo"
)

type SevenTvEmote struct {
	Name string `json:"name"`
}

type SevenUserTvResponse struct {
	EmoteSet *struct {
		Emotes []SevenTvEmote `json:"emotes"`
	} `json:"emote_set"`
	User struct {
		EmoteSets []struct {
			ID string `json:"id"`
		} `json:"emote_sets"`
	} `json:"user"`
}

type SevenTvGlobalResponse struct {
	Emotes []SevenTvEmote `json:"emotes"`
}

func GetChannelSevenTvEmotesSets(channelID string) ([]string, error) {
	resp, err := http.Get("https://7tv.io/v3/users/twitch/" + channelID)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode > 299 {
		return nil, nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	reqData := SevenUserTvResponse{}
	err = json.Unmarshal(body, &reqData)
	if err != nil {
		return nil, fmt.Errorf("cannot fetch 7tv emotes: %w", err)
	}

	mappedEmotesSets := lo.Map(
		reqData.User.EmoteSets, func(item struct {
			ID string `json:"id"`
		}, _ int) string {
			return item.ID
		},
	)

	return mappedEmotesSets, nil
}

func GetChannelSevenTvEmotes(channelID string) ([]string, error) {
	resp, err := http.Get("https://7tv.io/v3/users/twitch/" + channelID)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode > 299 {
		return nil, nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	reqData := SevenUserTvResponse{}
	err = json.Unmarshal(body, &reqData)
	if err != nil {
		return nil, fmt.Errorf("cannot fetch 7tv emotes: %w", err)
	}

	if reqData.EmoteSet == nil {
		return []string{}, nil
	}

	mappedEmotes := lo.Map(
		reqData.EmoteSet.Emotes, func(item SevenTvEmote, _ int) string {
			return item.Name
		},
	)

	return mappedEmotes, nil
}

func GetGlobalSevenTvEmotes() ([]string, error) {
	resp, err := http.Get("https://7tv.io/v3/emote-sets/global")
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	reqData := SevenTvGlobalResponse{}
	err = json.Unmarshal(body, &reqData)
	if err != nil {
		return nil, errors.New("cannot fetch 7tv emotes")
	}

	mappedEmotes := lo.Map(
		reqData.Emotes, func(item SevenTvEmote, _ int) string {
			return item.Name
		},
	)

	return mappedEmotes, nil
}
