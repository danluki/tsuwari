package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.57

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/samber/lo"
	data_loader "github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/dataloader"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/graph"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/mappers"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	"github.com/twirapp/twir/apps/api-gql/internal/services/roles"
	"github.com/twirapp/twir/apps/api-gql/internal/services/roles_with_roles_users"
)

// RolesCreate is the resolver for the rolesCreate field.
func (r *mutationResolver) RolesCreate(ctx context.Context, opts gqlmodel.RolesCreateOrUpdateOpts) (bool, error) {
	dashboardId, err := r.deps.Sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return false, err
	}

	user, err := r.deps.Sessions.GetAuthenticatedUser(ctx)
	if err != nil {
		return false, err
	}

	permissions := make([]string, len(opts.Permissions))
	for i, permission := range opts.Permissions {
		permissions[i] = permission.String()
	}

	users := make([]roles_with_roles_users.CreateInputUser, len(opts.Users))
	for idx, userId := range opts.Users {
		users[idx] = roles_with_roles_users.CreateInputUser{
			UserID: userId,
		}
	}

	err = r.deps.RolesWithUsersService.Create(
		ctx, roles_with_roles_users.CreateInput{
			Role: roles.CreateInput{
				ChannelID:                 dashboardId,
				ActorID:                   user.ID,
				Name:                      opts.Name,
				Type:                      entity.ChannelRoleTypeCustom,
				Permissions:               permissions,
				RequiredWatchTime:         int64(opts.Settings.RequiredWatchTime),
				RequiredMessages:          int32(opts.Settings.RequiredMessages),
				RequiredUsedChannelPoints: int64(opts.Settings.RequiredUserChannelPoints),
			},
			Users: users,
		},
	)
	if err != nil {
		return false, err
	}

	return true, nil
}

// RolesUpdate is the resolver for the rolesUpdate field.
func (r *mutationResolver) RolesUpdate(ctx context.Context, id uuid.UUID, opts gqlmodel.RolesCreateOrUpdateOpts) (bool, error) {
	dashboardId, err := r.deps.Sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return false, err
	}

	user, err := r.deps.Sessions.GetAuthenticatedUser(ctx)
	if err != nil {
		return false, err
	}

	permissions := make([]string, len(opts.Permissions))
	for i, permission := range opts.Permissions {
		permissions[i] = permission.String()
	}

	users := make([]roles_with_roles_users.CreateInputUser, len(opts.Users))
	for idx, userId := range opts.Users {
		users[idx] = roles_with_roles_users.CreateInputUser{
			UserID: userId,
		}
	}

	if err := r.deps.RolesWithUsersService.Update(
		ctx, roles_with_roles_users.UpdateInput{
			ID:        id,
			ChannelID: dashboardId,
			ActorID:   user.ID,
			Role: roles.UpdateInput{
				ChannelID:                 dashboardId,
				ActorID:                   user.ID,
				Name:                      &opts.Name,
				Permissions:               permissions,
				RequiredWatchTime:         lo.ToPtr(int64(opts.Settings.RequiredMessages)),
				RequiredMessages:          lo.ToPtr(int32(opts.Settings.RequiredMessages)),
				RequiredUsedChannelPoints: lo.ToPtr(int64(opts.Settings.RequiredUserChannelPoints)),
			},
			Users: users,
		},
	); err != nil {
		return false, err
	}

	return true, nil
}

// RolesRemove is the resolver for the rolesRemove field.
func (r *mutationResolver) RolesRemove(ctx context.Context, id uuid.UUID) (bool, error) {
	dashboardId, err := r.deps.Sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return false, err
	}

	user, err := r.deps.Sessions.GetAuthenticatedUser(ctx)
	if err != nil {
		return false, err
	}

	if err := r.deps.RolesService.Delete(
		ctx,
		roles.DeleteInput{
			ChannelID: dashboardId,
			ActorID:   user.ID,
			ID:        id,
		},
	); err != nil {
		return false, err
	}

	return true, nil
}

// Roles is the resolver for the roles field.
func (r *queryResolver) Roles(ctx context.Context) ([]gqlmodel.Role, error) {
	dashboardId, err := r.deps.Sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return nil, err
	}

	entities, err := r.deps.RolesService.GetManyByChannelID(ctx, dashboardId)
	if err != nil {
		return nil, err
	}

	result := make([]gqlmodel.Role, len(entities))
	for i, role := range entities {
		result[i] = mappers.RolesToGql(role)
	}

	return result, nil
}

// Users is the resolver for the users field.
func (r *roleResolver) Users(ctx context.Context, obj *gqlmodel.Role) ([]gqlmodel.TwirUserTwitchInfo, error) {
	if obj == nil {
		return nil, nil
	}

	users, err := r.deps.RolesUsersService.GetManyByRoleID(ctx, obj.ID)
	if err != nil {
		return nil, err
	}

	ids := make([]string, 0, len(users))
	for _, user := range users {
		ids = append(ids, user.UserID)
	}

	profiles, err := data_loader.GetHelixUsersByIds(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch twitch profiles: %w", err)
	}

	res := make([]gqlmodel.TwirUserTwitchInfo, 0, len(profiles))
	for _, profile := range profiles {
		res = append(res, *profile)
	}

	return res, nil
}

// Role returns graph.RoleResolver implementation.
func (r *Resolver) Role() graph.RoleResolver { return &roleResolver{r} }

type roleResolver struct{ *Resolver }
