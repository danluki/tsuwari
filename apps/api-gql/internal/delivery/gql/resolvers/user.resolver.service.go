package resolvers

import (
	"context"

	"github.com/samber/lo"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
)

func (r *authenticatedUserResolver) getAvailableDashboards(
	ctx context.Context,
	obj *gqlmodel.AuthenticatedUser,
) ([]gqlmodel.Dashboard, error) {
	dashboardsEntities := make(map[string]gqlmodel.Dashboard)
	if obj.IsBotAdmin {
		var channels []model.Channels
		if err := r.deps.Gorm.WithContext(ctx).Find(&channels).Error; err != nil {
			return nil, err
		}

		for _, channel := range channels {
			dashboardsEntities[channel.ID] = gqlmodel.Dashboard{
				ID: channel.ID,
				Flags: []gqlmodel.ChannelRolePermissionEnum{
					gqlmodel.ChannelRolePermissionEnumCanAccessDashboard,
				},
			}
		}
	} else {
		dashboardsEntities[obj.ID] = gqlmodel.Dashboard{
			ID:    obj.ID,
			Flags: []gqlmodel.ChannelRolePermissionEnum{gqlmodel.ChannelRolePermissionEnumCanAccessDashboard},
		}

		var roles []model.ChannelRoleUser
		if err := r.deps.Gorm.
			WithContext(ctx).
			Where(
				`"userId" = ?`,
				obj.ID,
			).
			Preload("Role").
			Preload("Role.Channel").
			Find(&roles).
			Error; err != nil {
			return nil, err
		}

		for _, role := range roles {
			if role.Role == nil || role.Role.Channel == nil || len(role.Role.Permissions) == 0 {
				continue
			}

			var flags []gqlmodel.ChannelRolePermissionEnum
			for _, flag := range role.Role.Permissions {
				flags = append(flags, gqlmodel.ChannelRolePermissionEnum(flag))
			}

			dashboardsEntities[role.Role.Channel.ID] = gqlmodel.Dashboard{
				ID:    role.Role.Channel.ID,
				Flags: append(dashboardsEntities[role.Role.Channel.ID].Flags, flags...),
			}
		}
	}

	var usersStats []model.UsersStats
	if err := r.deps.Gorm.
		WithContext(ctx).
		Where(`"userId" = ?`, obj.ID).
		Find(&usersStats).Error; err != nil {
		return nil, err
	}

	for _, stat := range usersStats {
		var channelRoles []model.ChannelRole
		if err := r.deps.Gorm.WithContext(ctx).Where(
			`"channelId" = ?`,
			stat.ChannelID,
		).Find(&channelRoles).
			Error; err != nil {
			return nil, err
		}

		var role model.ChannelRole

		if stat.IsMod {
			role, _ = lo.Find(
				channelRoles,
				func(role model.ChannelRole) bool {
					return role.Type == model.ChannelRoleTypeModerator
				},
			)
		} else if stat.IsVip {
			role, _ = lo.Find(
				channelRoles,
				func(role model.ChannelRole) bool {
					return role.Type == model.ChannelRoleTypeVip
				},
			)
		} else if stat.IsSubscriber {
			role, _ = lo.Find(
				channelRoles,
				func(role model.ChannelRole) bool {
					return role.Type == model.ChannelRoleTypeSubscriber
				},
			)
		}

		var flags []gqlmodel.ChannelRolePermissionEnum
		for _, flag := range role.Permissions {
			flags = append(flags, gqlmodel.ChannelRolePermissionEnum(flag))
		}

		if role.ID != "" && len(flags) > 0 {
			dashboardsEntities[role.ChannelID] = gqlmodel.Dashboard{
				ID:    role.ChannelID,
				Flags: append(dashboardsEntities[role.ChannelID].Flags, flags...),
			}
		}
	}

	return lo.MapToSlice(
		dashboardsEntities,
		func(_ string, value gqlmodel.Dashboard) gqlmodel.Dashboard {
			return value
		},
	), nil
}
