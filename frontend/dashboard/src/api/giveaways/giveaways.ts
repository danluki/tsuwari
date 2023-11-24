import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query';
import { ChooseWinnersRequest, ClearParticipantsRequest, CreateRequest, UpdateOrCreateRequest } from '@twir/grpc/generated/api/api/giveaways';
import { Ref, ref } from 'vue';

import { protectedApiClient } from '@/api/twirp';



export const useCurrentGiveaway = () =>
	useQuery({
		queryKey: ['currentGiveaway'],
		queryFn: async () => {
			const req = await protectedApiClient.giveawaysGetCurrent({});
			return req.response;
		},
	});

export const useParticipants = (id: Ref<string | undefined>, query: Ref<string> = ref('')) => useQuery({
	queryKey: ['participants'],
	queryFn: async () => {
		const req = await protectedApiClient.giveawaysGetParticipants({
			giveawayId: id.value as string,
			query: query.value,
		});
		return req.response;
	},
	enabled: !!id?.value,
});

export const useCreateGiveaway = () => {
	const queryClient = useQueryClient();

	return useMutation({
		mutationKey: ['createGiveaway'],
		mutationFn: async (opts: CreateRequest) => {
			await protectedApiClient.giveawaysCreate({
				...opts,
			});
		},
		onSuccess: async () => {
			await queryClient.invalidateQueries({
				queryKey: ['currentGiveaway'],
			});
		},
	});
};

export const useUpdateOrCreateGiveaway = () => {
	const queryClient = useQueryClient();

	return useMutation({
		mutationKey: ['updateGiveaway'],
		mutationFn: async (opts: UpdateOrCreateRequest) => {
			await protectedApiClient.giveawaysUpdateOrCreate({
				...opts,
			});
		},
		onSuccess: async () => {
			await queryClient.invalidateQueries({
				queryKey: ['currentGiveaway'],
			});
		},
	});
};

export const useChooseGiveawayWinners = () => {
	const queryClient = useQueryClient();

	return useMutation({
		mutationKey: ['chooseWinners'],
		mutationFn: async (opts: ChooseWinnersRequest) => {
			await protectedApiClient.giveawaysChooseWinners({
				...opts,
			});
		},
		onSuccess: async () => {
			await queryClient.invalidateQueries(
				['giveawayWinners'],
			);
		},
	});
};

export const useClearGiveawayParticipants = () => {
	const queryClient = useQueryClient();

	return useMutation({
		mutationKey: ['chooseWinners'],
		mutationFn: async (opts: ClearParticipantsRequest) => {
			await protectedApiClient.giveawaysClearParticipants({
				...opts,
			});
		},
		onSuccess: async () => {
			await queryClient.invalidateQueries(
				['participants'],
			);
		},
	});
};

export const useCurrentGiveawaysWinners = (giveawayId: Ref<string | undefined>) => useQuery({
	queryKey: ['giveawayWinners'],
	queryFn: async () => {
		const req = await protectedApiClient.giveawaysGetWinners({
			giveawayId: giveawayId.value as string,
		});
		return req.response;
	},
});
