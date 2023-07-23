<script setup lang='ts'>
import {
	IconTrash,
	IconChevronUp,
	IconChevronDown,
} from '@tabler/icons-vue';
import {
  type DataTableCreateSummary,
  NDataTable,
  NTag,
  NSpin,
  NSpace,
  NText,
  NCard,
  NButton,
	NTime,
} from 'naive-ui';
import type { TableColumn } from 'naive-ui/es/data-table/src/interface';
import { h, computed } from 'vue';
import { useI18n } from 'vue-i18n';

import { convertMillisToTime } from '@/components/songRequests/helpers.js';
import type { Video } from '@/components/songRequests/hook.js';

const props = defineProps<{
	queue: Video[]
}>();
const emits = defineEmits<{
	deleteVideo: [id: string]
	deleteAllVideos: []
	moveVideo: [id: string, newPosition: number]
}>();

const { t } = useI18n();

const columns = computed<TableColumn<Video>[]>(() => [
	{
		title: '#',
		key: 'position',
		width: 25,
		render(_, index) {
			return index+1;
		},
	},
	{
		title: t('sharedTexts.name'),
		key: 'title',
		render(row) {
			return h(NButton, {
				tag: 'a',
				type: 'primary',
				text: true,
				target: '_blank',
				href: row.songLink,
			}, {
				default: () => row.title,
			});
		},
	},
	{
		title: t('songRequests.table.columns.author'),
		key: 'author',
		render(row) {
			return h(NTag, { bordered: false, type: 'info' }, { default: () => row.orderedByDisplayName || row.orderedByName });
		},
	},
	{
		title: t('songRequests.table.columns.added'),
		key: 'createdAt',
		width: 150,
		render(row) {
			return h(NTime, { time: 0, to: Date.now() - new Date(row.createdAt).getTime(), type: 'relative' });
		},
	},
	{
		title: t('songRequests.table.columns.duration'),
		key: 'duration',
		width: 100,
		render(row) {
			return convertMillisToTime(row.duration * 1000);
		},
	},
	{
		title: '',
		key: 'actions',
		width: 150,
		render(row, index) {
			const deleteButton = h(
				NButton,
				{
					size: 'tiny',
					type: 'error',
					text: true,
					onClick: () => emits('deleteVideo', row.id),
				}, {
					default: () => h(IconTrash),
				},
			);

			const moveUpButton = h(NButton, {
				size: 'tiny',
				type: 'primary',
				text: true,
				disabled: index === 0,
				onClick: () => emits('moveVideo', row.id, index-1),
			}, {
				default: () => h(IconChevronUp),
			});

			const moveDownButton = h(NButton, {
				size: 'tiny',
				type: 'primary',
				text: true,
				disabled: index+1 === props.queue.length,
				onClick: () => emits('moveVideo', row.id, index+1),
			}, {
				default: () => h(IconChevronDown),
			});

			return h(NSpace, {
				justify: 'center',
				align: 'center',
			}, {
				default: () => [
					deleteButton,
					moveUpButton,
					moveDownButton,
				],
			});
		},
	},
]);

const createSummary: DataTableCreateSummary<Video> = (pageData) => {
	return{
		position: {
			value: h(
				'span',
				{ style: 'font-weight: bold;' },
				pageData.length,
			),
			colSpan: 4,
		},
		duration: {
			value: h(
				'span',
				{ style: 'font-weight: bold;' },
				convertMillisToTime(pageData.reduce((acc, cur) => acc + cur.duration * 1000, 0)),
			),
			colSpan: 2,
		},
	};
};
</script>

<template>
	<n-card
		:title="t('songRequests.table.title')"
		content-style="padding: 0;"
		header-style="padding: 10px;"
		segmented
	>
		<template #header-extra>
			<n-button tertiary size="small" @click="$emit('deleteAllVideos')">
				<IconTrash />
			</n-button>
		</template>
		<n-data-table
			:columns="columns"
			:data="queue"
			:loading="!queue.length"
			:bordered="false"
			:summary="createSummary"
		>
			<template #loading>
				<n-space vertical align="center" style="margin-top: 50px;">
					<n-spin :rotate="false" stroke="#959596">
						<template #description>
							<n-text>{{ t('songRequests.waiting') }}</n-text>
						</template>
					</n-spin>
				</n-space>
			</template>
		</n-data-table>
	</n-card>
</template>