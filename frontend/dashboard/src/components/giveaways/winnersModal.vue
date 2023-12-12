<script setup lang="ts">
import { NCard, NSkeleton } from 'naive-ui';
import { computed } from 'vue';

import { useCurrentGiveaway, useCurrentGiveawaysWinners } from '@/api';


const { data: currentGiveaway } = useCurrentGiveaway();
const giveawayId = computed(() => currentGiveaway.value?.id);

const { data: winners, isLoading, isError } = useCurrentGiveawaysWinners(giveawayId);

</script>
<template>
	<n-card
		size="large"
		:bordered="false"
		aria-modal="true"
	>
		<n-skeleton v-if="isLoading" height="100%" width="100%" />
		<div v-else-if="isError">
			Error fetching giveaway winners
		</div>
		<ul style="list-style-type: none">
			<li v-for="winner of winners?.winners" :key="winner.userId">
				{{ winner.displayName }}
			</li>
		</ul>
	</n-card>
</template>
