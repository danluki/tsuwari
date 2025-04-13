package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.57

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/graph"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/mappers"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	"github.com/twirapp/twir/apps/api-gql/internal/services/badges"
	badges_users "github.com/twirapp/twir/apps/api-gql/internal/services/badges-users"
)

// Users is the resolver for the users field.
func (r *badgeResolver) Users(ctx context.Context, obj *gqlmodel.Badge) ([]string, error) {
	users, err := r.deps.BadgesUsersService.GetMany(
		ctx,
		badges_users.GetManyInput{
			BadgeID: obj.ID,
		},
	)
	if err != nil {
		r.deps.Logger.Error("cannot get badge users", slog.Any("err", err))
		return nil, err
	}

	userIds := make([]string, 0, len(users))
	for _, user := range users {
		userIds = append(userIds, user.UserID)
	}

	return userIds, nil
}

// BadgesDelete is the resolver for the badgesDelete field.
func (r *mutationResolver) BadgesDelete(ctx context.Context, id uuid.UUID) (bool, error) {
	if err := r.deps.BadgesService.Delete(ctx, id); err != nil {
		r.deps.Logger.Error("cannot delete badge", slog.Any("err", err))
		return false, err
	}

	return true, nil
}

// BadgesUpdate is the resolver for the badgesUpdate field.
func (r *mutationResolver) BadgesUpdate(ctx context.Context, id uuid.UUID, opts gqlmodel.TwirBadgeUpdateOpts) (*gqlmodel.Badge, error) {
	input := badges.UpdateInput{
		Name:    opts.Name.Value(),
		Enabled: opts.Enabled.Value(),
		FfzSlot: opts.FfzSlot.Value(),
	}

	if opts.File.IsSet() {
		f := opts.File.Value()

		input.File = &entity.Upload{
			File:        f.File,
			Filename:    f.Filename,
			Size:        f.Size,
			ContentType: f.ContentType,
		}
	}

	newBadge, err := r.deps.BadgesService.Update(ctx, id, input)
	if err != nil {
		r.deps.Logger.Error("cannot update badge", slog.Any("err", err))
		return nil, err
	}

	converted := mappers.BadgeEntityToGql(newBadge)
	return &converted, nil
}

// BadgesCreate is the resolver for the badgesCreate field.
func (r *mutationResolver) BadgesCreate(ctx context.Context, opts gqlmodel.TwirBadgeCreateOpts) (*gqlmodel.Badge, error) {
	input := badges.CreateInput{
		Name:    opts.Name,
		Enabled: true,
		FfzSlot: opts.FfzSlot,
		File: entity.Upload{
			File:        opts.File.File,
			Filename:    opts.File.Filename,
			Size:        opts.File.Size,
			ContentType: opts.File.ContentType,
		},
	}

	if opts.Enabled.IsSet() {
		input.Enabled = *opts.Enabled.Value()
	}

	newBadge, err := r.deps.BadgesService.Create(ctx, input)
	if err != nil {
		r.deps.Logger.Error("cannot create badge", slog.Any("err", err))
		return nil, err
	}

	converted := mappers.BadgeEntityToGql(newBadge)
	return &converted, nil
}

// BadgesAddUser is the resolver for the badgesAddUser field.
func (r *mutationResolver) BadgesAddUser(ctx context.Context, id uuid.UUID, userID string) (bool, error) {
	_, err := r.deps.BadgesUsersService.Create(
		ctx,
		badges_users.CreateInput{
			BadgeID: id,
			UserID:  userID,
		},
	)
	if err != nil {
		r.deps.Logger.Error("cannot add user to badge", slog.Any("err", err))
		return false, err
	}

	return true, nil
}

// BadgesRemoveUser is the resolver for the badgesRemoveUser field.
func (r *mutationResolver) BadgesRemoveUser(ctx context.Context, id uuid.UUID, userID string) (bool, error) {
	err := r.deps.BadgesUsersService.Delete(
		ctx,
		badges_users.DeleteInput{
			BadgeID: id,
			UserID:  userID,
		},
	)
	if err != nil {
		r.deps.Logger.Error("cannot remove user from badge", slog.Any("err", err))
		return false, err
	}

	return true, nil
}

// TwirBadges is the resolver for the twirBadges field.
func (r *queryResolver) TwirBadges(ctx context.Context) ([]gqlmodel.Badge, error) {
	entities, err := r.deps.BadgesService.GetMany(
		ctx,
		badges.GetManyInput{},
	)
	if err != nil {
		return nil, err
	}

	result := make([]gqlmodel.Badge, 0, len(entities))
	for _, b := range entities {
		result = append(result, mappers.BadgeEntityToGql(b))
	}

	return result, nil
}

// Badge returns graph.BadgeResolver implementation.
func (r *Resolver) Badge() graph.BadgeResolver { return &badgeResolver{r} }

type badgeResolver struct{ *Resolver }
