package config

type TeamConfig struct {
	TeamName            string   `yaml:"team_name"`
	LeaderName          string   `yaml:"leader_name"`
	LeaderId            *uint64  `yaml:"leader_id"`
	FollowTeamMember    string   `yaml:"follow_team_member"`
	FollowTeamLeader    bool     `yaml:"follow_team_leader"`
	AutoAcceptTeamApply bool     `yaml:"auto_accept"`
	AllowedTeamMembers  []string `yaml:"allowed_team_members"`
}

func (c *TeamConfig) GetLeaderName() string {
	return c.LeaderName
}

func (c *TeamConfig) GetLeaderId() *uint64 {
	if c.LeaderId == nil {
		return new(uint64)
	}
	return c.LeaderId
}

func (c *TeamConfig) GetTeamName() string {
	return c.TeamName
}

func (c *TeamConfig) AutoAccept() bool {
	return c.AutoAcceptTeamApply
}
