package repository

import "github.com/skarakasoglu/aybush-bot-api/data/models"

type OrderEpisodeExperienceBy uint8

const (
	OrderEpisodeExperience_ExperiencePoints OrderEpisodeExperienceBy = iota
	OrderEpisodeExperience_MessageCount
	OrderEpisodeExperience_ActiveVoiceMinutes
)

var orderEpisodeExperiences = map[OrderEpisodeExperienceBy]string{
	OrderEpisodeExperience_ExperiencePoints: "\"experience_points\"",
	OrderEpisodeExperience_ActiveVoiceMinutes: "\"active_voice_minutes\"",
	OrderEpisodeExperience_MessageCount: "\"message_count\"",
}

func (o OrderEpisodeExperienceBy) String() string {
	str, ok := orderEpisodeExperiences[o]
	if !ok {
		return ""
	}

	return str
}

type DiscordRepository interface{
	GetTopNDiscordMemberLevels(n int) ([]models.DiscordMemberLevel, error)
	GetCurrentlyActiveEpisodes() ([]models.DiscordEpisode, error)
	GetEpisodeLeaderboard(n int, episodeId int, orderBy OrderEpisodeExperienceBy) ([]models.DiscordEpisodeExperience, error)
}
