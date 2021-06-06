package service

import (
	"database/sql"
	"fmt"
	"github.com/skarakasoglu/aybush-bot-api/data/models"
	"github.com/skarakasoglu/aybush-bot-api/repository"
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
			inner join "discord_levels" as cdl on dml.experience_points between cdl.required_experience_points and (cdl.maximum_experience_points - 1)
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

func (d DiscordService) GetCurrentlyActiveEpisodes() ([]models.DiscordEpisode, error) {
	var episodes []models.DiscordEpisode

	query := `SELECT id, name, start_timestamp, end_timestamp FROM "discord_episodes"
				WHERE NOW() BETWEEN start_timestamp and end_timestamp;`

	rows, err := d.db.Query(query)
	if err != nil {
		log.Printf("Error on executing the statement: %v", err)
		return episodes, err
	}
	defer rows.Close()

	for rows.Next() {
		var episode models.DiscordEpisode

		err = rows.Scan(&episode.Id, &episode.Name, &episode.StartTimestamp, &episode.EndTimestamp)
		if err != nil {
			log.Printf("Error on scanning the row: %v", err)
		}

		episodes = append(episodes, episode)
	}

	return episodes, nil
}

func (d DiscordService) GetEpisodeLeaderboard(n int, episodeId int, orderBy repository.OrderEpisodeExperienceBy) ([]models.DiscordEpisodeExperience, error) {
	var experiences []models.DiscordEpisodeExperience

	query := fmt.Sprintf(`
				SELECT dee.id, de.id "episode_id", dm.member_id, COALESCE(dm.avatar_url, ''), dm.username, dm.discriminator, 
				dee.experience_points "experience_points", dee.active_voice_minutes "active_voice_minutes", COALESCE(dmm.message_count, 0) "message_count",
				cdl.id "current_level", cdl.required_experience_points "current_level_required", cdl.maximum_experience_points "current_level_maximum",
				ndl.id "next_level", ndl.required_experience_points "next_level_required", ndl.maximum_experience_points "next_level_required"
				FROM "discord_episode_experiences" dee
				INNER JOIN "discord_episodes" de ON de.id = dee.episode_id
				LEFT JOIN LATERAL (SELECT member_id, COUNT(*) AS "message_count" FROM "discord_member_messages" dmm 
								   WHERE dmm.created_at BETWEEN de.start_timestamp and de.end_timestamp GROUP BY member_id) dmm 
				ON dmm.member_id = dee.member_id
				INNER JOIN "discord_members" dm ON dm.member_id = dee.member_id
				INNER JOIN "discord_levels" cdl ON dee.experience_points between cdl.required_experience_points AND (cdl.maximum_experience_points - 1)
				INNER JOIN "discord_levels" ndl ON cdl.maximum_experience_points = ndl.required_experience_points
				WHERE episode_id = $1 AND dm.is_left = false
				ORDER BY %s DESC LIMIT $2;
		`, orderBy.String())

	preparedStmt, err := d.db.Prepare(query)
	if err != nil {
		log.Printf("Error on preparing the statement: %v", err)
		return experiences, err
	}
	defer preparedStmt.Close()

	rows, err := preparedStmt.Query(episodeId, n)
	if err != nil {
		log.Printf("Error on querying the statement: %v", err)
		return experiences, err
	}
	defer rows.Close()

	for rows.Next() {
		var experience models.DiscordEpisodeExperience
		var episode models.DiscordEpisode
		var member models.DiscordMember
		var currentLevel models.DiscordLevel
		var nextLevel models.DiscordLevel

		var messageCount sql.NullInt64
		err = rows.Scan(&experience.Id, &episode.Id, &member.MemberId, &member.AvatarUrl, &member.Username, &member.Discriminator, &experience.ExperiencePoints, &experience.ActiveVoiceMinutes, &messageCount,
			&currentLevel.Id, &currentLevel.RequiredExperiencePoints, &currentLevel.MaximumExperiencePoints,
			&nextLevel.Id, &nextLevel.RequiredExperiencePoints, &nextLevel.MaximumExperiencePoints)
		if err != nil {
			log.Printf("Error on scanning the row: %v", err)
		}

		experience.DiscordEpisode = episode
		experience.DiscordMember = member
		experience.MessageCount = uint64(messageCount.Int64)
		experience.CurrentLevel = currentLevel
		experience.NextLevel = nextLevel

		experiences = append(experiences, experience)
	}

	return experiences, nil
}

func NewDiscordService(db *sql.DB) *DiscordService{
	return &DiscordService{
		db: db,
	}
}