package api

import (
	"github.com/gin-gonic/gin"
	"github.com/skarakasoglu/aybush-bot-api/api/response"
	"github.com/skarakasoglu/aybush-bot-api/repository"
	"log"
	"net/http"
	"strconv"
)

const (
	CurrentTopMemberCount = 50
)

type apiV1 struct{
	discordRepository repository.DiscordRepository
}

func NewApiV1(discordRepository repository.DiscordRepository) *apiV1 {
	return &apiV1{
		discordRepository: discordRepository,
	}
}

func (api *apiV1) getLeaderboard(ctx *gin.Context) {
	var resp response.LeaderboardResponse

	memberLevels, err := api.discordRepository.GetTopNDiscordMemberLevels(CurrentTopMemberCount)
	if err != nil {
		log.Printf("Error on fetching member levels: %v", err)
		resp.Error = err.Error()
		resp.Code = http.StatusInternalServerError

		ctx.JSON(http.StatusOK, resp)
		return
	}

	resp.Code = http.StatusOK

	for i, member := range memberLevels {
		resp.Leaderboard = append(resp.Leaderboard, response.MemberLevelStatus{
			Member: response.Member{
				MemberId:      member.MemberId,
				AvatarUrl: member.AvatarUrl,
				Username:      member.Username,
				Discriminator: member.Discriminator,
				JoinedAt:      member.JoinedAt,
			},
			ExperiencePoints: member.ExperiencePoints,
			MessageCount: member.MessageCount,
			ActiveVoiceMinutes: member.ActiveVoiceMinutes,
			CurrentLevel:     response.Level{
				Level:                    member.CurrentLevel.Id,
				RequiredExperiencePoints: member.CurrentLevel.RequiredExperiencePoints,
			},
			NextLevel:        response.Level{
				Level:                    member.NextLevel.Id,
				RequiredExperiencePoints: member.NextLevel.RequiredExperiencePoints,
			},
			RoleName:         member.CurrentLevel.DiscordRole.Name,
			Position: i + 1,
		})
	}

	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.JSON(http.StatusOK, resp)
}

func (api *apiV1) getLeaderboardByEpisode(ctx *gin.Context) {
	var resp response.LeaderboardResponse
	resp.Leaderboard = make([]response.MemberLevelStatus, 0)

	episodeIdParam := ctx.Query("episode")
	orderByParam := ctx.Query("order")

	episodeId, err := strconv.Atoi(episodeIdParam)
	if err != nil {
		log.Printf("Error on parsing episode id: %v", err)

		resp.Error = err.Error()
		resp.Code = http.StatusBadRequest

		ctx.JSON(resp.Code, resp)
		return
	}

	orderBy, err := strconv.Atoi(orderByParam)
	if err != nil {
		log.Printf("Error on parsing episode id: %v", err)

		resp.Error = err.Error()
		resp.Code = http.StatusBadRequest

		ctx.JSON(resp.Code, resp)
		return
	}

	episodes, err := api.discordRepository.GetEpisodeLeaderboard(CurrentTopMemberCount, episodeId, repository.OrderEpisodeExperienceBy(orderBy))
	resp.Code = http.StatusOK

	for i, member := range episodes {
		resp.Leaderboard = append(resp.Leaderboard, response.MemberLevelStatus{
			Member: response.Member{
				MemberId:      member.MemberId,
				AvatarUrl: member.AvatarUrl,
				Username:      member.Username,
				Discriminator: member.Discriminator,
				JoinedAt:      member.JoinedAt,
			},
			ExperiencePoints: member.ExperiencePoints,
			MessageCount: member.MessageCount,
			ActiveVoiceMinutes: member.ActiveVoiceMinutes,
			CurrentLevel:     response.Level{
				Level:                    member.CurrentLevel.Id,
				RequiredExperiencePoints: member.CurrentLevel.RequiredExperiencePoints,
			},
			NextLevel:        response.Level{
				Level:                    member.NextLevel.Id,
				RequiredExperiencePoints: member.NextLevel.RequiredExperiencePoints,
			},
			RoleName:         member.CurrentLevel.DiscordRole.Name,
			Position: i + 1,
		})
	}

	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.JSON(resp.Code, resp)
}

func (api *apiV1) getActiveEpisodes(ctx *gin.Context) {
	var resp response.EpisodesResponse

	episodes, err := api.discordRepository.GetCurrentlyActiveEpisodes()
	if err != nil {
		log.Printf("Error on fetching episodes: %v", err)
		return
	}

	resp.Code = http.StatusOK

	for _, episode := range episodes {
		resp.Episodes = append(resp.Episodes, response.Episode{
			Id:             episode.Id,
			Name:           episode.Name,
			StartTimestamp: episode.StartTimestamp,
			EndTimestamp:   episode.EndTimestamp,
		})
	}

	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.JSON(http.StatusOK, resp)
}