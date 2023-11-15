<script setup lang="ts">
import { IconRotate } from '@tabler/icons-vue';
import { NCard, NInput, NButton } from 'naive-ui';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();

defineProps<{
	searchValue: string;
	participantsCount: number;
	users: giveaways[];
}>();

const emit = defineEmits<{
	(event: 'update:searchValue', payload: string): void;
}>();
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
					:value="searchValue"
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
					<IconRotate />
				</n-button>
			</div>
		</div>
		<div style="padding: 10px">
			<ul>
				<li v-for="user in users" :key="user">
					{{ user }}
				</li>
			</ul>
		</div>

		<template #header-extra>
			{{ participantsCount }} Users
		</template>
	</n-card>
</template>

<style scoped></style>
