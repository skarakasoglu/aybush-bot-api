package service

import (
	"database/sql"
	"github.com/skarakasoglu/aybush-bot-api/data/models"
	"log"
)

type DiscordService struct{
	db* sql.DB
}

func (d DiscordService) GetTopNDiscordMemberLevels(n int) ([]models.DiscordMemberLevel, error) {
	query := `
			SELECT dml.id, dm.avatar_url, dm.username, dm.discriminator, dml.message_count, dml.active_voice_minutes,
			dml.experience_points, dml.last_message_timestamp, dm.member_id, dm.guild_id,
			cdl.id "current_level",  cdl.required_experience_points "current_level_required", cdl.maximum_experience_points "current_level_maximum", 
			cdr.id "current_role_id", cdr.role_id "current_role_role_id", cdr.name "current_role_name",
			ndl.id "next_level", ndl.required_experience_points "next_level_required", ndl.maximum_experience_points "next_level_maximum",
			ndr.id "next_role_id", ndr.role_id "next_role_role_id", ndr.name "next_role_name"
			FROM "discord_member_levels" as dml 
			inner join "discord_members" as dm on dm.member_id = dml.member_id
			inner join "discord_levels" as cdl on dml.experience_points between cdl.required_experience_points and cdl.maximum_experience_points
			inner join "discord_levels" as ndl on cdl.maximum_experience_points = ndl.required_experience_points
			inner join "discord_roles" as cdr on cdr.role_id = cdl.role_id
			inner join "discord_roles" as ndr on ndr.role_id = ndl.role_id 
			WHERE dm.is_left = false
			ORDER BY dml.experience_points DESC LIMIT $1;
	`

	rows, err := d.db.Query(query, n)
	if err != nil {
		log.Printf("[DiscordService] Error on executing the query: %v", err)
		return nil, err
	}

	var memberLevels []models.DiscordMemberLevel

	for rows.Next() {
		var memberLevel models.DiscordMemberLevel
		var member models.DiscordMember
		var currentLevel models.DiscordLevel
		var nextLevel models.DiscordLevel

		err = rows.Scan(&memberLevel.Id, &member.AvatarUrl, &member.Username, &member.Discriminator,
			&memberLevel.MessageCount, &memberLevel.ActiveVoiceMinutes, &memberLevel.ExperiencePoints, &memberLevel.LastMessageTimestamp, &member.MemberId,
			&member.GuildId, &currentLevel.Id, &currentLevel.RequiredExperiencePoints, &currentLevel.MaximumExperiencePoints,
			&currentLevel.DiscordRole.Id, &currentLevel.DiscordRole.RoleId, &currentLevel.DiscordRole.Name,
			&nextLevel.Id, &nextLevel.RequiredExperiencePoints, &nextLevel.MaximumExperiencePoints,
			&nextLevel.DiscordRole.Id, &nextLevel.DiscordRole.RoleId, &nextLevel.DiscordRole.Name)

		if err != nil {
			log.Printf("[DiscordService] Error on scanning the row: %v", err)
			continue
		}

		memberLevel.DiscordMember = member
		memberLevel.CurrentLevel = currentLevel
		memberLevel.NextLevel = nextLevel

		memberLevels = append(memberLevels, memberLevel)
	}

	return memberLevels, nil
}

func NewDiscordService(db *sql.DB) *DiscordService{
	return &DiscordService{
		db: db,
	}
}