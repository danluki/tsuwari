import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query';
import { CreateOrGetRequest, UpdateRequest } from '@twir/grpc/generated/api/api/giveaways';

import { protectedApiClient } from '@/api/twirp';

export const useCurrentGiveaway = (params: CreateOrGetRequest) =>
	useQuery({
		queryKey: ['currentGiveaway'],
		queryFn: async () => {
			const req = await protectedApiClient.giveawaysCreateOrGet(params);
			return req.response;
		},
	});

export const userCurrentGiveawayUpdateSettings = () => {
	const queryClient = useQueryClient();

	return useMutation({
		mutationKey: ['currentGiveaway'],
		mutationFn: async (opts: UpdateRequest) => {
			const req = await protectedApiClient.giveawaysUpdate(opts);
		}
	}),
}
