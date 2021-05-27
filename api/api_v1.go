package api

import (
	"github.com/gin-gonic/gin"
	"github.com/skarakasoglu/aybush-bot-api/api/response"
	"github.com/skarakasoglu/aybush-bot-api/repository"
	"log"
	"net/http"
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
	}

	resp.Code = http.StatusOK

	for i, member := range memberLevels {
		resp.MemberLevels = append(resp.MemberLevels, response.MemberLevelStatus{
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