<script setup lang="ts">
import { IconRotate, IconTrophyFilled } from '@tabler/icons-vue';
import { NCard, NInput, NButton } from 'naive-ui';
import { computed, ref, onMounted } from 'vue';
import { useI18n } from 'vue-i18n';

import { useCurrentGiveaway, useParticipants } from '@/api';

const { t } = useI18n();

const emit = defineEmits<{
	(event: 'update:searchValue', payload: string): void;
}>();

const { data: currentGiveaway, isError, isFetched, isLoading } = useCurrentGiveaway();
const giveawayId = computed(() => currentGiveaway.value?.id);

const searchValue = ref('');

const {
  data: participants,
	isFetched: isParticipantsFetched,
	refetch: refreshParticipants,
} = useParticipants(giveawayId, searchValue);
const participantsCount = computed(() => participants.value?.totalCount);



const refreshUsers = () => {
	setInterval(() => {
		if (currentGiveaway.value?.isRunning) {
			refreshParticipants();
		}
	}, 5000);
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
					@update:value="$emit('update:searchValue', $event)"
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
				>
					<IconRotate />
				</n-button>
			</div>
		</div>
		<div style="padding: 10px">
			<ul>
				<li v-for="user in participants?.winners" :key="user.displayName">
					{{ user }}
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
