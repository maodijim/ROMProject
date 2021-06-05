package config

type TeamConfig struct {
	TeamName         string `yaml:"team_name"`
	LeaderName       string `yaml:"leader_name"`
	FollowTeamMember string `yaml:"follow_team_member"`
	FollowTeamLeader bool   `yaml:"follow_team_leader"`
}

func (c *TeamConfig) GetLeaderName() string {
	return c.LeaderName
}

func (c *TeamConfig) GetTeamName() string {
	return c.TeamName
}
