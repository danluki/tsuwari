package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.57

import (
	"context"
	"fmt"

	model "github.com/satont/twir/libs/gomodels"
	data_loader "github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/dataloader"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/graph"
)

// TwitchProfile is the resolver for the twitchProfile field.
func (r *emoteStatisticTopUserResolver) TwitchProfile(ctx context.Context, obj *gqlmodel.EmoteStatisticTopUser) (*gqlmodel.TwirUserTwitchInfo, error) {
	return data_loader.GetHelixUserById(ctx, obj.UserID)
}

// TwitchProfile is the resolver for the twitchProfile field.
func (r *emoteStatisticUserUsageResolver) TwitchProfile(ctx context.Context, obj *gqlmodel.EmoteStatisticUserUsage) (*gqlmodel.TwirUserTwitchInfo, error) {
	return data_loader.GetHelixUserById(ctx, obj.UserID)
}

// EmotesStatistics is the resolver for the emotesStatistics field.
func (r *queryResolver) EmotesStatistics(ctx context.Context, opts gqlmodel.EmotesStatisticsOpts) (*gqlmodel.EmotesStatisticResponse, error) {
	dashboardId, err := r.deps.Sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return nil, err
	}

	var page int
	perPage := 10

	if opts.Page.IsSet() {
		page = *opts.Page.Value()
	}

	if opts.PerPage.IsSet() {
		perPage = *opts.PerPage.Value()
	}

	query := r.deps.Gorm.WithContext(ctx).
		Where(`"channelId" = ?`, dashboardId).
		Limit(perPage).
		Offset(page * perPage)

	if opts.Search.IsSet() && *opts.Search.Value() != "" {
		query = query.Where(`"emote" LIKE ?`, "%"+*opts.Search.Value()+"%")
	}

	var order gqlmodel.EmotesStatisticsOptsOrder
	if opts.Order.IsSet() {
		order = *opts.Order.Value()
	} else {
		order = gqlmodel.EmotesStatisticsOptsOrderDesc
	}

	var entities []emoteEntityModelWithCount
	if err :=
		query.
			Select(`"emote", COUNT(emote) as count`).
			Group("emote").
			Order(fmt.Sprintf("count %s", order.String())).
			Find(&entities).
			Error; err != nil {
		return nil, err
	}

	var totalCount int64
	if err := r.deps.Gorm.
		WithContext(ctx).
		Raw(
			`
				SELECT COUNT(DISTINCT emote)
				FROM channels_emotes_usages
				WHERE "channelId" = ?
				`,
			dashboardId,
		).
		Scan(&totalCount).Error; err != nil {
		return nil, fmt.Errorf("failed to get total count: %w", err)
	}

	models := make([]gqlmodel.EmotesStatistic, 0, len(entities))
	for _, entity := range entities {
		lastUsedEntity := &model.ChannelEmoteUsage{}
		if err := r.deps.Gorm.
			WithContext(ctx).
			Where(`"channelId" = ? AND "emote" = ?`, dashboardId, entity.Emote).
			Order(`"createdAt" DESC`).
			First(lastUsedEntity).Error; err != nil {
			return nil, err
		}

		var rangeType gqlmodel.EmoteStatisticRange
		if opts.GraphicRange.IsSet() {
			rangeType = *opts.GraphicRange.Value()
		} else {
			rangeType = gqlmodel.EmoteStatisticRangeLastDay
		}

		graphicUsages, err := r.getEmoteStatisticUsagesForRange(
			ctx,
			entity.Emote,
			rangeType,
		)
		if err != nil {
			return nil, err
		}

		models = append(
			models, gqlmodel.EmotesStatistic{
				EmoteName:         entity.Emote,
				TotalUsages:       entity.Count,
				LastUsedTimestamp: int(lastUsedEntity.CreatedAt.UTC().UnixMilli()),
				GraphicUsages:     graphicUsages,
			},
		)
	}

	return &gqlmodel.EmotesStatisticResponse{
		Emotes: models,
		Total:  int(totalCount),
	}, nil
}

// EmotesStatisticEmoteDetailedInformation is the resolver for the emotesStatisticEmoteDetailedInformation field.
func (r *queryResolver) EmotesStatisticEmoteDetailedInformation(ctx context.Context, opts gqlmodel.EmotesStatisticEmoteDetailedOpts) (*gqlmodel.EmotesStatisticEmoteDetailedResponse, error) {
	if opts.EmoteName == "" {
		return nil, nil
	}

	dashboardId, err := r.deps.Sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return nil, err
	}

	graphicUsages, err := r.getEmoteStatisticUsagesForRange(ctx, opts.EmoteName, opts.Range)
	if err != nil {
		return nil, err
	}

	lastUsedEntity := &model.ChannelEmoteUsage{}
	if err := r.deps.Gorm.
		WithContext(ctx).
		Where(`"channelId" = ? AND "emote" = ?`, dashboardId, opts.EmoteName).
		Order(`"createdAt" DESC`).
		First(lastUsedEntity).Error; err != nil {
		return nil, err
	}

	var usages int64
	if err := r.deps.Gorm.
		WithContext(ctx).
		Model(&model.ChannelEmoteUsage{}).
		Where(`"channelId" = ? AND "emote" = ?`, dashboardId, opts.EmoteName).
		Count(&usages).Error; err != nil {
		return nil, err
	}

	var usagedByUserPage int
	usagesByUserPerPage := 10
	if opts.UsagesByUsersPage.IsSet() {
		usagedByUserPage = *opts.UsagesByUsersPage.Value()
	}
	if opts.UsagesByUsersPerPage.IsSet() {
		usagesByUserPerPage = *opts.UsagesByUsersPerPage.Value()
	}

	var usagesHistoryEntities []model.ChannelEmoteUsage
	if err := r.deps.Gorm.
		WithContext(ctx).
		Where(`"channelId" = ? AND "emote" = ?`, dashboardId, opts.EmoteName).
		Order(`"createdAt" DESC`).
		Limit(usagesByUserPerPage).
		Offset(usagedByUserPage * usagesByUserPerPage).
		Find(&usagesHistoryEntities).Error; err != nil {
		return nil, err
	}
	var usagesByUsersTotalCount int64
	if err := r.deps.Gorm.
		WithContext(ctx).
		Model(&model.ChannelEmoteUsage{}).
		Where(`"channelId" = ? AND "emote" = ?`, dashboardId, opts.EmoteName).
		Count(&usagesByUsersTotalCount).Error; err != nil {
		return nil, err
	}

	var topUsersPage int
	topUsersPerPage := 10
	if opts.TopUsersPage.IsSet() {
		topUsersPage = *opts.TopUsersPage.Value()
	}
	if opts.TopUsersPerPage.IsSet() {
		topUsersPerPage = *opts.TopUsersPerPage.Value()
	}

	var topUsersEntities []emoteEntityModelWithCount
	if err := r.deps.Gorm.
		WithContext(ctx).
		Where(`"channelId" = ? AND "emote" = ?`, dashboardId, opts.EmoteName).
		Select(`"userId", COUNT("userId") as count`).
		Group("userId").
		Order("count DESC").
		Limit(topUsersPerPage).
		Offset(topUsersPage * topUsersPerPage).
		Find(&topUsersEntities).Error; err != nil {
		return nil, err
	}

	var topUsersTotalCount int64
	if err := r.deps.Gorm.
		WithContext(ctx).
		Model(&model.ChannelEmoteUsage{}).
		Where(`"channelId" = ? AND "emote" = ?`, dashboardId, opts.EmoteName).
		Group(`"userId"`).
		Count(&topUsersTotalCount).Error; err != nil {
		return nil, err
	}

	usagesHistory := make([]gqlmodel.EmoteStatisticUserUsage, 0, len(usagesHistoryEntities))
	for _, usage := range usagesHistoryEntities {
		usagesHistory = append(
			usagesHistory,
			gqlmodel.EmoteStatisticUserUsage{
				UserID: usage.UserID,
				Date:   usage.CreatedAt,
			},
		)
	}

	topUsers := make([]gqlmodel.EmoteStatisticTopUser, 0, len(topUsersEntities))
	for _, user := range topUsersEntities {
		topUsers = append(
			topUsers,
			gqlmodel.EmoteStatisticTopUser{
				UserID: user.UserID,
				Count:  user.Count,
			},
		)
	}

	return &gqlmodel.EmotesStatisticEmoteDetailedResponse{
		EmoteName:          opts.EmoteName,
		TotalUsages:        int(usages),
		LastUsedTimestamp:  int(lastUsedEntity.CreatedAt.UTC().UnixMilli()),
		GraphicUsages:      graphicUsages,
		UsagesHistory:      usagesHistory,
		UsagesByUsersTotal: int(usagesByUsersTotalCount),
		TopUsers:           topUsers,
		TopUsersTotal:      int(topUsersTotalCount),
	}, nil
}

// EmoteStatisticTopUser returns graph.EmoteStatisticTopUserResolver implementation.
func (r *Resolver) EmoteStatisticTopUser() graph.EmoteStatisticTopUserResolver {
	return &emoteStatisticTopUserResolver{r}
}

// EmoteStatisticUserUsage returns graph.EmoteStatisticUserUsageResolver implementation.
func (r *Resolver) EmoteStatisticUserUsage() graph.EmoteStatisticUserUsageResolver {
	return &emoteStatisticUserUsageResolver{r}
}

type emoteStatisticTopUserResolver struct{ *Resolver }
type emoteStatisticUserUsageResolver struct{ *Resolver }
