<script setup lang="ts">
import { watchDebounced } from '@vueuse/core';
import { Giveaway } from 'libs/grpc/dist/types/generated/api/api/giveaways';
import {
	NButton,
	NCard,
	NInput,
	NRadio,
	NRadioGroup,
	NSlider,
	NInputNumber,
} from 'naive-ui';
import { ref } from 'vue';
import { useI18n } from 'vue-i18n';

import { useCurrentGiveaway } from '@/api';

const { t } = useI18n();

const props = defineProps<{
	giveaway: Giveaway
}>();



const giveawayUpdate = giveawaysManager.update;

const giveawayType = ref(props.giveaway.type);
const isGiveawaysRunning = ref(props.giveaway.isRunning);
const keyword = ref(props.giveaway.keyword);
const subscribersLuck = ref(props.giveaway.subscribersLuck);
const followersLuck = ref(props.giveaway.followersLuck);
const minNumber = ref(props.giveaway.randomNumberFrom);
const maxNumber = ref(props.giveaway.randomNumberTo);
const description = ref(props.giveaway.description);
const winnersCount = ref(props.giveaway.winnerCount);



async function changeGiveawayStatus() {
	isGiveawaysRunning.value = !isGiveawaysRunning.value;

	await giveawayUpdate.mutateAsync({
		giveawayId: props.giveaway.id,
		isRunning: isGiveawaysRunning.value,
	});
}


const hasAccessToManageGiveaways = useUserAccessFlagChecker('MANAGE_GIVEAWAYS');
</script>
<template>
	<div class="flex-container">
		<n-card
			:title="t('giveaways.settings.title')" content-style="padding: 0;" header-style="padding: 10px;"
			style="min-width: 300px; height: 100%; overflow-y: auto" segmented
		>
			<div style="height: 95%; padding: 10px">
				<n-input :value="description" style="margin-bottom: 25px" type="textarea" placeholder="Giveaway description" />
				<n-radio-group :value="giveawayType" style="margin-bottom: 25px" name="giveawaysTypesGroup">
					<div>
						<div style="margin-bottom: 7px">
							{{ t('giveaways.settings.giveawayType') }}
						</div>
						<n-space>
							<n-radio value="BY_KEYWORD" label="Keyword" />
							<n-radio value="BY_RANDOM_NUMBER" label="Random number" />
						</n-space>
					</div>
				</n-radio-group>
				<n-space vertical>
					<n-input
						v-if="giveawayType === 'BY_KEYWORD'" :value="keyword" style="margin-bottom: 25px"
						:disabled="giveawayType !== 'BY_KEYWORD'" type="text" placeholder="Keyword Phrase"
					/>
					<div v-else style="margin-bottom: 25px">
						<n-input-number :value="minNumber" style="margin-bottom: 10px" placeholder="Minimum number" min="0" />
						<n-input-number :value="maxNumber" placeholder="Maximum number" min="0" />
					</div>
					<div style="margin-bottom: 25px">
						<div style="margin-bottom: 3px">
							{{ t('giveaways.settings.subscribersLuck') }}
						</div>
						<n-slider :value="subscribersLuck" :step="1" :max="10" :min="0" />
					</div>

					<div style="margin-bottom: 25px">
						<div style="margin-bottom: 3px">
							{{ t('giveaways.settings.followersLuck') }}
						</div>
						<n-slider :value="followersLuck" :step="1" :max="10" :min="0" />
					</div>

					<n-input-number :value="maxNumber" style="margin-bottom: 25px" placeholder="Minimum watch time" min="0" />

					<n-input-number :value="winnersCount" style="margin-bottom: 25px" placeholder="Winners count" min="1" />

					<n-button type="primary" style="width: 100%">
						{{
							isGiveawaysRunning && participantsCount && participantsCount > 0
								? t('giveaways.settings.roll')
								: t('giveaways.settings.cantRoll')
						}}
					</n-button>
				</n-space>
			</div>
			<template #header-extra>
				<n-button secondary type="primary" :disabled="!hasAccessToManageGiveaways" @click="changeGiveawayStatus">
					{{
						isGiveawaysRunning
							? t('giveaways.settings.running')
							: t('giveaways.settings.stopped')
					}}
				</n-button>
			</template>
		</n-card>
		</n-card>
	</div>
</template>

<style scoped></style>
