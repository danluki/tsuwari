<script setup lang="ts">
import { Winner } from '@twir/grpc/generated/api/api/giveaways';
import { NCard } from 'naive-ui';


const props = defineProps<{
	winners: Winner[],
}>();

</script>
<template>
	<n-card
		title="Giveaway winners"
		:bordered="false"
		size="huge"
		aria-modal="true"
	>
		<ul>
			<li v-for="(winner, index) of props.winners" :key="winner.userId">
				{{ winner.displayName }}
			</li>
		</ul>
	</n-card>
</template>
