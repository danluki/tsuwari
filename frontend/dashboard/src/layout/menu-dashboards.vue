<script lang="ts" setup>
import { IconChevronRight } from '@tabler/icons-vue'
import { onClickOutside, onKeyStroke } from '@vueuse/core'
import {
	NAvatar,
	NInput,
	NPopover,
	NSpin,
	NText,
	NVirtualList,
	useThemeVars,
} from 'naive-ui'
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

import { useDashboard, useProfile } from '@/api/index.js'
import { resolveUserName } from '@/helpers/resolveUserName.js'
import { useSidebarCollapseStore } from '@/layout/use-sidebar-collapse'

const props = withDefaults(defineProps<{ isDrawer?: boolean }>(), {
	isDrawer: false,
})

const emits = defineEmits<{
	dashboardSelected: []
}>()

const { t } = useI18n()
const themeVars = useThemeVars()
const blockColor = computed(() => themeVars.value.buttonColor2)
const blockColor2 = computed(() => themeVars.value.buttonColor2Hover)

const { data: profile, isLoading: isProfileLoading } = useProfile()
const { setDashboard } = useDashboard()

const currentDashboard = computed(() => {
	const dashboard = profile.value?.availableDashboards.find(dashboard => dashboard.id === profile.value?.selectedDashboardId)
	if (!dashboard) return null

	return dashboard
})

const activeDashboard = ref('')
watch(currentDashboard, (v) => {
	if (!v) return
	activeDashboard.value = v.id
}, { immediate: true })

watch(activeDashboard, async (dashboardId) => {
	if (dashboardId === profile.value?.selectedDashboardId) return

	await setDashboard(dashboardId)
	emits('dashboardSelected')
})

const filterValue = ref('')

const menuOptions = computed(() => {
	return profile.value?.availableDashboards
		.filter(dashboard => {
			return dashboard.twitchProfile.displayName.includes(filterValue.value)
			  || dashboard.twitchProfile.login.includes(filterValue.value)
		})
		.map((u) => {
			return {
				key: u.id,
				label: resolveUserName(u.twitchProfile.login, u.twitchProfile.displayName),
				icon: u.twitchProfile.profileImageUrl,
			}
		}) ?? []
})

const isSelectDashboardPopoverOpened = ref(false)

function togglePopover(value?: boolean) {
	isSelectDashboardPopoverOpened.value = value ?? !isSelectDashboardPopoverOpened.value
}

function onSelectDashboard(key: string) {
	activeDashboard.value = key
	togglePopover(false)
}

onKeyStroke('k', (event) => {
	if (event.ctrlKey || event.metaKey) {
		event.preventDefault()
		togglePopover()
	}
})

const refPopoverList = ref<HTMLElement | null>()
const refPopover = ref<HTMLElement | null>()
onClickOutside(refPopover, (event) => {
	if (isSelectDashboardPopoverOpened.value) {
		event.stopPropagation()
		togglePopover(false)
	}
}, { ignore: [refPopoverList] })

const { isCollapsed } = useSidebarCollapseStore()

const displayNameLength = computed(() => {
	if (!currentDashboard.value) return 0
	return currentDashboard.value.twitchProfile.displayName.length
})

const isDrawerCollapsed = computed(() => {
	return props.isDrawer || !isCollapsed.value
})

const popoverPlacement = computed(() => {
	if (props.isDrawer) return 'bottom'
	return isCollapsed.value ? 'right-start' : 'bottom-start'
})
</script>

<template>
	<NPopover
		ref="refPopover"
		:placement="popoverPlacement"
		trigger="manual"
		class="w-[240px] !m-0"
		:show="isSelectDashboardPopoverOpened"
		:show-arrow="false"
	>
		<template #trigger>
			<div
				class="popover-trigger flex items-center gap-4 rounded-[10px] cursor-pointer"
				@click="isSelectDashboardPopoverOpened = true"
			>
				<div class="flex items-center justify-between w-full py-3 px-3.5 ">
					<div class="flex gap-3">
						<NAvatar
							round
							class="flex self-center"
							:src="currentDashboard?.twitchProfile.profileImageUrl"
						/>

						<div
							v-if="isDrawerCollapsed"
							class="flex flex-col whitespace-nowrap overflow-hidden overflow-ellipsis"
						>
							<NText :depth="3" class="whitespace-nowrap text-xs">
								{{ t(`dashboard.header.managingUser`) }}
							</NText>

							<NText :class="[displayNameLength > 16 ? 'text-xs' : 'text-sm']">
								{{ currentDashboard?.twitchProfile.displayName }}
							</NText>
						</div>
					</div>

					<IconChevronRight
						v-if="isDrawerCollapsed"
						:style="{
							transition: '0.2s transform ease',
							transform: `rotate(${!isSelectDashboardPopoverOpened ? 90 : -90}deg)`,
						}"
					/>
				</div>
			</div>
		</template>

		<NSpin v-if="isProfileLoading"></NSpin>

		<div v-else ref="refPopoverList" class="dashboards-container">
			<NText :depth="3" class="text-xs">
				{{ t(`dashboard.header.channelsAccess`) }}
			</NText>

			<NVirtualList
				class="max-h-[400px]"
				:item-size="42"
				trigger="none"
				:items="menuOptions"
				item-resizable
			>
				<template #default="{ item }">
					<div
						:key="item.key"
						class="item h-10"
						@click="onSelectDashboard(item.key)"
					>
						<NAvatar :src="item.icon" round size="small" />
						<span> {{ item.label }}</span>
					</div>
				</template>
			</NVirtualList>

			<template v-if="(profile?.availableDashboards.length ?? 0) > 10">
				<NInput v-model:value="filterValue" placeholder="Search" />
			</template>
		</div>
	</NPopover>
</template>

<style scoped>
.dashboards-container {
	@apply select-none;
}

.dashboards-container :deep(img) {
	-webkit-user-drag: none;
}

.item {
	@apply flex items-center gap-3 w-full p-1.5 rounded-md cursor-pointer;
	background-color: v-bind(blockColor);
}

.dashboards-menu > .item:hover {
	background-color: v-bind(blockColor2);
}

.popover-trigger {
	@apply flex w-full select-none;
}

.popover-trigger :deep(img) {
	-webkit-user-drag: none;
}

:deep(.v-vl) {
	@apply overflow-x-hidden;
}
</style>
