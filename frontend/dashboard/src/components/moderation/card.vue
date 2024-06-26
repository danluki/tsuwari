<script setup lang="ts">
import { IconSettings, IconTrash } from '@tabler/icons-vue'
import { NButton, NPopconfirm, NSwitch, useNotification } from 'naive-ui'
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'

import { Icons } from './helpers.js'

import type { ItemWithId } from '@twir/api/messages/moderation/moderation'

import { useModerationManager, useUserAccessFlagChecker } from '@/api/index.js'
import Card from '@/components/card/card.vue'
import { ChannelRolePermissionEnum } from '@/gql/graphql'

const props = defineProps<{
	item: ItemWithId
}>()

defineEmits<{
	showSettings: []
}>()

const manager = useModerationManager()
const patcher = manager.patch!
const deleter = manager.deleteOne

const patchExecuting = ref(false)

const { t } = useI18n()

const userCanManageModeration = useUserAccessFlagChecker(ChannelRolePermissionEnum.ManageModeration)

const message = useNotification()

async function switchState(id: string, v: boolean) {
	patchExecuting.value = true

	try {
		await patcher.mutateAsync({ id, enabled: v })
	} catch (error) {
		console.error(error)
	} finally {
		patchExecuting.value = false
	}
}

async function removeItem() {
	await deleter.mutateAsync({ id: props.item.id })
	message.success({
		title: t('sharedTexts.deleted'),
		duration: 2000,
	})
}
</script>

<template>
	<Card
		:title="t(`moderation.types.${item.data!.type}.name`)"
		:icon="Icons[item.data!.type]"
		style="height:100%"
	>
		<template #headerExtra>
			<NSwitch
				:disabled="!userCanManageModeration"
				:value="item.data!.enabled"
				:loading="patchExecuting"
				@update:value="(v) => switchState(item.id, v)"
			/>
		</template>

		<template #content>
			{{ t(`moderation.types.${item.data!.type}.description`) }}
		</template>

		<template #footer>
			<div class="flex gap-2">
				<NButton
					:disabled="!userCanManageModeration"
					secondary
					size="large"
					@click="$emit('showSettings')"
				>
					<div class="flex gap-1">
						<span>{{ t('sharedButtons.settings') }}</span>
						<IconSettings />
					</div>
				</NButton>
				<NPopconfirm
					:positive-text="t('deleteConfirmation.confirm')"
					:negative-text="t('deleteConfirmation.cancel')"
					@positive-click="removeItem"
				>
					<template #trigger>
						<NButton
							:disabled="!userCanManageModeration"
							secondary
							size="large"
							type="error"
						>
							<div class="flex gap-1">
								<span>{{ t('sharedButtons.delete') }}</span>
								<IconTrash />
							</div>
						</NButton>
					</template>
					{{ t('deleteConfirmation.text') }}
				</NPopconfirm>
			</div>
		</template>
	</Card>
</template>
