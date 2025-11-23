package model

type Stats struct {
	TotalTeams       int `json:"total_teams"`
	TotalUsers       int `json:"total_users"`
	ActiveUsers      int `json:"active_users"`
	TotalPRs         int `json:"total_prs"`
	OpenPRs          int `json:"open_prs"`
	MergedPRs        int `json:"merged_prs"`
	TotalAssignments int `json:"total_assignments"`
}

type UsersStats struct {
	TotalUsers  int `json:"total_users"`
	ActiveUsers int `json:"active_users"`
}

type PRStats struct {
	TotalPRs  int `json:"total_prs"`
	OpenPRs   int `json:"open_prs"`
	MergedPRs int `json:"merged_prs"`
}
