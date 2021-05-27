package repository

import "github.com/skarakasoglu/aybush-bot-api/data/models"

type DiscordRepository interface{
	GetTopNDiscordMemberLevels(n int) ([]models.DiscordMemberLevel, error)
}
