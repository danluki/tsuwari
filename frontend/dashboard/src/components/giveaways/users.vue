<script setup lang="ts">
import { IconRotate, IconTrophyFilled } from '@tabler/icons-vue';
import { watchDebounced } from '@vueuse/core';
import { NCard, NInput, NButton } from 'naive-ui';
import { computed, ref, onMounted } from 'vue';
import { useI18n } from 'vue-i18n';

import { useClearGiveawayParticipants, useCurrentGiveaway, useCurrentGiveawaysWinners, useParticipants } from '@/api';

const props = defineProps<{
	showWinnersModal: () => void
}>();

const { t } = useI18n();

const { data: currentGiveaway, isError, isFetched, isLoading } = useCurrentGiveaway();
const giveawayId = computed(() => currentGiveaway.value?.id);
const searchValue = ref('');

const {
  data: participants,
	isFetched: isParticipantsFetched,
	refetch: refreshParticipants,
} = useParticipants(giveawayId, searchValue);
const participantsCount = computed(() => participants.value?.totalCount);

watchDebounced(
  searchValue,
  async () => {
    refreshParticipants();
  },
  { debounce: 200, maxWait: 500 },
);

const refreshUsers = () => {
	setInterval(async () => {
		if (!currentGiveaway.value?.isFinished) {
			await refreshParticipants();
		}
	}, 5000);
};

const resetParticipants = useClearGiveawayParticipants();

const { data: winners } = useCurrentGiveawaysWinners(giveawayId);

const resetUsers = async () => {
	await resetParticipants.mutateAsync({
		giveawayId: giveawayId.value,
	});
};

const showWinners = async () => {
	props.showWinnersModal();
};

onMounted(refreshUsers);

</script>
<template>
	<n-card
		:title="t('giveaways.users.title')"
		content-style="padding: 0;"
		header-style="padding: 10px;"
		style="min-width: 300px; height: 100%"
		segmented
	>
		<div style="padding: 10px">
			<div
				style="
					display: flex;
					flex-direction: row;
					align-items: flex-start;
					justify-content: flex-start;
				"
			>
				<n-input
					v-model:value="searchValue"
					type="text"
					placeholder="eg. TwirApp"
				/>
				<n-button
					type="tertiary"
					style="
						display: flex;
						align-items: center;
						justify-content: center;
						height: 34px;
						margin: 0 5px;
					"
					:disabled="winners && winners.length > 0"
					@click="showWinners"
				>
					<IconTrophyFilled />
				</n-button>
				<n-button
					type="tertiary"
					style="
						display: flex;
						align-items: center;
						justify-content: center;
						height: 34px;
						margin: 0 5px;
					"
					@click="resetUsers"
				>
					<IconRotate />
				</n-button>
			</div>
		</div>
		<div>
			<ul style="list-style-type: none">
				<li v-for="user in participants?.winners" :key="user.displayName">
					{{ user.displayName }}
				</li>
			</ul>
		</div>

		<template #header-extra>
			<div v-if="isParticipantsFetched">
				{{ participantsCount }} Users
			</div>
			<div v-else>
				0 Users
			</div>
		</template>
	</n-card>
</template>

<style scoped></style>
