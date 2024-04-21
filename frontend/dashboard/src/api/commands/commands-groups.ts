import { useQuery } from '@urql/vue';
import { defineStore } from 'pinia';

import { invalidationKey as commandsInvalidationKey } from './commands.js';

import { useMutation } from '@/composables/use-mutation';
import { graphql } from '@/gql';


const invalidationKey = 'CommandsGroupsInvalidateKey';

export const useCommandsGroupsApi = defineStore('api/commands-groups', () => {
	const useQueryGroups = () => useQuery({
		query: graphql(`
			query GetAllCommandsGroups {
				commandsGroups {
					id
					name
					color
				}
			}
		`),
		variables: {},
		context: {
			additionalTypenames: [invalidationKey],
		},
	});

	const useMutationDeleteGroup = () => useMutation(
		graphql(`
			mutation DeleteCommandGroup($id: ID!) {
				commandsGroupsRemove(id: $id)
			}
		`),
		[invalidationKey, commandsInvalidationKey],
	);

	const useMutationCreateGroup = () => useMutation(
		graphql(`
			mutation CreateCommandGroup($opts: CommandsGroupsCreateOpts!) {
				commandsGroupsCreate(opts: $opts)
			}
		`),
		[invalidationKey],
	);

	const useMutationUpdateGroup = () => useMutation(
		graphql(`
			mutation UpdateCommandGroup($id: ID!, $opts: CommandsGroupsUpdateOpts!) {
				commandsGroupsUpdate(id: $id,opts: $opts)
			}
		`),
		[invalidationKey],
	);

	return {
		useQueryGroups,
		useMutationDeleteGroup,
		useMutationCreateGroup,
		useMutationUpdateGroup,
	};
});