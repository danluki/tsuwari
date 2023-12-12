-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TYPE "channels_giveaways_type_enum" AS ENUM (
	'BY_KEYWORD',
	'BY_RANDOM_NUMBER'
);

CREATE TABLE "channels_giveaways" (
	"id" uuid default uuid_generate_v4() primary key,
	"description" text not null,
	"type" channels_giveaways_type_enum not null,
	"channel_id" text not null,
	"created_at" timestamp not null default now(),
	"start_at" timestamp not null,
	"end_at" timestamp not null,
	"closed_at" timestamp not null,
	"is_running" boolean not null default 'false',
	"is_finished" boolean not null default 'false',
	"required_min_watch_time" integer,
	"required_min_follow_time" integer,
	"require_min_messages" integer,
	"require_min_subscriber_tier" integer,
	"required_min_subscriber_time" integer,
	"eligible_user_groups" text NOT NULL,
	"keyword" varchar default 'BY_KEYWORD',
	"random_number_from" integer,
	"random_number_to" integer,
	"winner_random_number" integer,
	"winners_count" integer not null,
	"followers_luck" integer not null default '0',
	"subscribers_luck" integer not null default '0',
	"subscribers_tier1_luck" integer not null default '0',
	"subscribers_tier2_luck" integer not null default '0',
	"subscribers_tier3_luck" integer not null default '0',
	"messages_count_luck" integer not null default '0'
);

CREATE TABLE "channels_giveaways_participants" (
	"id" uuid default uuid_generate_v4() primary key,
	"giveaway_id" uuid not null,
	"is_winner" boolean not null default false,
	"user_id" text not null,
	"display_name" text not null,
	"is_subscriber" boolean not null default false,
	"is_follower" boolean not null default false,
	"is_moderator" boolean not null default false,
	"is_vip" boolean not null default false,
	"subscriber_tier" integer,
	"user_follow_since" timestamp,
	"user_stats_watched_time" bigint not null,
	"messages_count" integer not null default '0',
	"start_ticket" int not null default '0',
	"end_ticket" int not null default '0'
);

ALTER TABLE "channels_giveaways" ADD CONSTRAINT "channels_giveaways_channels_channel_fk" FOREIGN KEY ("channel_id") REFERENCES "channels"("id") ON DELETE CASCADE ON UPDATE NO ACTION;
ALTER TABLE "channels_giveaways_participants" ADD CONSTRAINT "channel_giveaways_channel_giveaways_participants_giveaway_fk" FOREIGN KEY ("giveaway_id") REFERENCES "channels_giveaways"("id") ON DELETE NO ACTION ON UPDATE NO ACTION;
ALTER TABLE "channels_giveaways_participants" ADD CONSTRAINT "channel_giveaways_channel_giveaways_participants_user_fk" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE NO ACTION ON UPDATE NO ACTION;
ALTER TABLE "channels_modules_settings" ADD CONSTRAINT "channel_giveaways_channel_modules_settings_user_fk" FOREIGN KEY ("userId") REFERENCES "users"("id") ON DELETE NO ACTION ON UPDATE NO ACTION;
ALTER TABLE "channels_giveaways_participants" ADD CONSTRAINT "channels_giveaways_participants_unique" UNIQUE("giveaway_id", "user_id");

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
ALTER TABLE channels_giveaways_participants
DROP CONSTRAINT channel_giveaways_channel_giveaways_participants_user_fk;

ALTER TABLE channels_giveaways_participants
DROP CONSTRAINT channel_giveaways_channel_giveaways_participants_giveaway_fk;

ALTER TABLE channels_giveaways
DROP CONSTRAINT channels_giveaways_channels_channel_fk;

-- Drop the tables
DROP TABLE IF EXISTS channels_giveaways_participants;
DROP TABLE IF EXISTS channels_giveaways;

-- Drop the enum type
DROP TYPE IF EXISTS channels_giveaways_type_enum;
-- +goose StatementEnd
