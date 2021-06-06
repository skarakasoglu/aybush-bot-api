package models

type DiscordEpisodeExperience struct {
	Id int
	DiscordEpisode
	DiscordMember
	CurrentLevel DiscordLevel
	NextLevel DiscordLevel
	ExperiencePoints uint64
	ActiveVoiceMinutes uint64
	MessageCount uint64
}
