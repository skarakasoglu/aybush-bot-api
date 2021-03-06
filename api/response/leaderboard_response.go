package response

import "time"

type Level struct{
	Level int `json:"level"`
	RequiredExperiencePoints uint64 `json:"required_experience_points"`
}

type Member struct{
	MemberId string `json:"member_id"`
	AvatarUrl string `json:"avatar_url"`
	Username string `json:"username"`
	Discriminator string `json:"discriminator"`
	JoinedAt time.Time `json:"joined_at"`
}

type MemberLevelStatus struct{
	Member
	ExperiencePoints uint64 `json:"experience_points"`
	MessageCount uint64 `json:"message_count"`
	ActiveVoiceMinutes uint64 `json:"active_voice_minutes"`
	CurrentLevel Level `json:"current_level"`
	NextLevel Level `json:"next_level"`
	RoleName string `json:"role_name"`
	Position int `json:"position"`
}

type LeaderboardResponse struct {
	Response
	Leaderboard []MemberLevelStatus `json:"leaderboard"`
}
