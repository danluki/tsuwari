<script setup lang="ts">
import { ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

import type { ChatWall } from '@/api/moderation-chat-wall.ts'

import { useModerationChatWall } from '@/api/moderation-chat-wall.ts'
import DialogOrSheet from '@/components/dialog-or-sheet.vue'
import { Button } from '@/components/ui/button'
import {
	Dialog,
	DialogHeader,
	DialogTitle,
	DialogTrigger,
} from '@/components/ui/dialog'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table/index.ts'

const { t } = useI18n()

const props = defineProps<{
	chatWall: ChatWall
}>()

const api = useModerationChatWall()
const { data, executeQuery } = api.useLogs(props.chatWall.id)

const dialogOpened = ref(false)

watch(dialogOpened, (v) => {
	if (v) {
		executeQuery()
	}
})
</script>

<template>
	<Dialog v-model:open="dialogOpened">
		<DialogTrigger as-child>
			<Button :disabled="!chatWall.affectedMessages" size="sm">
				{{ t('chatWall.table.affectedMessages') }} ({{ chatWall.affectedMessages }})
			</Button>
		</DialogTrigger>
		<DialogOrSheet>
			<DialogHeader>
				<DialogTitle>{{ t('chatWall.table.logs.title') }}</DialogTitle>
			</DialogHeader>

			<Table class="bg-sidebar rounded">
				<TableHeader>
					<TableRow>
						<TableHead class="w-[10%]">
							{{ t('chatWall.table.logs.user') }}
						</TableHead>
						<TableHead>
							{{ t('chatWall.table.logs.message') }}
						</TableHead>
					</TableRow>
				</TableHeader>
				<TableBody>
					<TableRow v-for="message of data?.chatWallLogs" :key="message.id">
						<TableCell class="w-[10%]">
							<a :href="`https://twitch.tv/${message.twitchProfile.login}`" class="flex items-center gap-2">
								<img :src="message.twitchProfile.profileImageUrl" class="size-6 rounded-full" />
								<span>
									{{ message.twitchProfile.displayName }}
								</span>
							</a>
						</TableCell>
						<TableCell>
							<span class="break-words">
								{{ message.text }}
							</span>
						</TableCell>
					</TableRow>
				</TableBody>
			</Table>
		</DialogOrSheet>
	</Dialog>
</template>
