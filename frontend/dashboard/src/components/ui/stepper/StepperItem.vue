<script lang="ts" setup>
import { StepperItem, useForwardProps } from 'reka-ui'
import { type HTMLAttributes, computed } from 'vue'

import type { StepperItemProps } from 'reka-ui'

import { cn } from '@/lib/utils'

const props = defineProps<StepperItemProps & { class?: HTMLAttributes['class'] }>()

const delegatedProps = computed(() => {
	const { class: _, ...delegated } = props

	return delegated
})

const forwarded = useForwardProps(delegatedProps)
</script>

<template>
	<StepperItem
		v-slot="slotProps"
		v-bind="forwarded"
		:class="cn('flex items-center gap-2 group data-[disabled]:pointer-events-none', props.class)"
	>
		<slot v-bind="slotProps" />
	</StepperItem>
</template>
