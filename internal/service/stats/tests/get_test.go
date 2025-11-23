package stats_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Dokhoyan/avito-pr-test/internal/model"
	repomocks "github.com/Dokhoyan/avito-pr-test/internal/repository/mocks"
	"github.com/Dokhoyan/avito-pr-test/internal/service/stats"
	"github.com/Dokhoyan/avito-pr-test/internal/service/testutil"
	"github.com/stretchr/testify/assert"
)

func TestGetStats_Success(t *testing.T) {
	ctx := context.Background()

	expectedStats := &model.Stats{
		TotalTeams:       2,
		TotalUsers:       5,
		ActiveUsers:      4,
		TotalPRs:         7,
		OpenPRs:          5,
		MergedPRs:        2,
		TotalAssignments: 9,
	}

	usersStats := &model.UsersStats{
		TotalUsers:  5,
		ActiveUsers: 4,
	}

	prStats := &model.PRStats{
		TotalPRs:  7,
		OpenPRs:   5,
		MergedPRs: 2,
	}

	mockStatsRepo := repomocks.NewMockStatsRepository(t)
	mockTrManager := testutil.NewTestTrManager(t)

	mockStatsRepo.EXPECT().GetTeamsCount(ctx).Return(2, nil)
	mockStatsRepo.EXPECT().GetUsersStats(ctx).Return(usersStats, nil)
	mockStatsRepo.EXPECT().GetPRStats(ctx).Return(prStats, nil)
	mockStatsRepo.EXPECT().GetAssignmentsCount(ctx).Return(9, nil)

	svc := stats.NewServiceWithTrManager(mockStatsRepo, mockTrManager)
	result, err := svc.GetStats(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedStats.TotalTeams, result.TotalTeams)
	assert.Equal(t, expectedStats.TotalUsers, result.TotalUsers)
	assert.Equal(t, expectedStats.ActiveUsers, result.ActiveUsers)
	assert.Equal(t, expectedStats.TotalPRs, result.TotalPRs)
	assert.Equal(t, expectedStats.OpenPRs, result.OpenPRs)
	assert.Equal(t, expectedStats.MergedPRs, result.MergedPRs)
	assert.Equal(t, expectedStats.TotalAssignments, result.TotalAssignments)
	mockStatsRepo.AssertExpectations(t)
}

func TestGetStats_GetTeamsCountError(t *testing.T) {
	ctx := context.Background()

	mockStatsRepo := repomocks.NewMockStatsRepository(t)
	mockTrManager := testutil.NewTestTrManager(t)

	mockStatsRepo.EXPECT().GetTeamsCount(ctx).Return(0, errors.New("database error"))

	svc := stats.NewServiceWithTrManager(mockStatsRepo, mockTrManager)
	result, err := svc.GetStats(ctx)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockStatsRepo.AssertExpectations(t)
}

func TestGetStats_GetUsersStatsError(t *testing.T) {
	ctx := context.Background()

	mockStatsRepo := repomocks.NewMockStatsRepository(t)
	mockTrManager := testutil.NewTestTrManager(t)

	mockStatsRepo.EXPECT().GetTeamsCount(ctx).Return(2, nil)
	mockStatsRepo.EXPECT().GetUsersStats(ctx).Return(nil, errors.New("database error"))

	svc := stats.NewServiceWithTrManager(mockStatsRepo, mockTrManager)
	result, err := svc.GetStats(ctx)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockStatsRepo.AssertExpectations(t)
}

func TestGetStats_GetPRStatsError(t *testing.T) {
	ctx := context.Background()

	usersStats := &model.UsersStats{
		TotalUsers:  5,
		ActiveUsers: 4,
	}

	mockStatsRepo := repomocks.NewMockStatsRepository(t)
	mockTrManager := testutil.NewTestTrManager(t)

	mockStatsRepo.EXPECT().GetTeamsCount(ctx).Return(2, nil)
	mockStatsRepo.EXPECT().GetUsersStats(ctx).Return(usersStats, nil)
	mockStatsRepo.EXPECT().GetPRStats(ctx).Return(nil, errors.New("database error"))

	svc := stats.NewServiceWithTrManager(mockStatsRepo, mockTrManager)
	result, err := svc.GetStats(ctx)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockStatsRepo.AssertExpectations(t)
}

func TestGetStats_GetAssignmentsCountError(t *testing.T) {
	ctx := context.Background()

	usersStats := &model.UsersStats{
		TotalUsers:  5,
		ActiveUsers: 4,
	}

	prStats := &model.PRStats{
		TotalPRs:  7,
		OpenPRs:   5,
		MergedPRs: 2,
	}

	mockStatsRepo := repomocks.NewMockStatsRepository(t)
	mockTrManager := testutil.NewTestTrManager(t)

	mockStatsRepo.EXPECT().GetTeamsCount(ctx).Return(2, nil)
	mockStatsRepo.EXPECT().GetUsersStats(ctx).Return(usersStats, nil)
	mockStatsRepo.EXPECT().GetPRStats(ctx).Return(prStats, nil)
	mockStatsRepo.EXPECT().GetAssignmentsCount(ctx).Return(0, errors.New("database error"))

	svc := stats.NewServiceWithTrManager(mockStatsRepo, mockTrManager)
	result, err := svc.GetStats(ctx)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockStatsRepo.AssertExpectations(t)
}
