<script setup lang="ts">
import { UpdateRequest } from '@twir/grpc/generated/api/api/giveaways';
import { watchDebounced } from '@vueuse/core';
import {
  NButton,
  NCard,
  NInput,
  NRadio,
  NRadioGroup,
  NSlider,
  NInputNumber,
  NForm,
  NSpace,
} from 'naive-ui';
import { computed, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';

import {
  useCurrentGiveaway,
  useUserAccessFlagChecker,
  useParticipants,
  useUpdateOrCreateGiveaway,
  useChooseGiveawayWinners,
} from '@/api';

const props = defineProps<{
  showWinnersModal: () => void;
}>();

const { t } = useI18n();

const { data: currentGiveaway, isError, isFetched, isLoading } = useCurrentGiveaway();
const giveawayId = computed(() => currentGiveaway.value?.id);

const { data: participants } = useParticipants(giveawayId);
const participantsCount = computed(() => participants.value?.totalCount);

const formValue = ref<UpdateRequest>({
  description: '',
  type: 'BY_KEYWORD',
  keyword: '',
  randomNumberFrom: 0,
  randomNumberTo: 10,
  subscribersLuck: 1,
  followersLuck: 1,
  winnersCount: 1,
});
const updateOrCreateGiveaway = useUpdateOrCreateGiveaway();

watch(currentGiveaway, () => {
  formValue.value.description = currentGiveaway.value?.description;
  formValue.value.type = currentGiveaway.value?.type;
  formValue.value.keyword = currentGiveaway.value?.keyword;
  formValue.value.randomNumberFrom = currentGiveaway.value?.randomNumberFrom;
  formValue.value.randomNumberTo = currentGiveaway.value?.randomNumberTo;
  formValue.value.subscribersLuck = currentGiveaway.value?.subscribersLuck;
  formValue.value.followersLuck = currentGiveaway.value?.followersLuck;
  formValue.value.winnersCount = currentGiveaway.value?.winnerCount;
  formValue.value.isFinished = currentGiveaway.value?.isFinished;
  formValue.value.isRunning = currentGiveaway.value?.isRunning;
});

watchDebounced(
  formValue.value,
  async () => {
    await updateOrCreateGiveaway.mutateAsync({
      description: formValue.value.description,
      followersLuck: formValue.value.followersLuck,
      keyword: formValue.value.keyword,
      randomNumberFrom: formValue.value.randomNumberFrom,
      randomNumberTo: formValue.value.randomNumberTo,
      subscribersLuck: formValue.value.subscribersLuck,
      winnersCount: formValue.value.winnersCount,
      isRunning: formValue.value.isRunning,
      isFinished: formValue.value.isFinished,
      type: formValue.value.type,
    });
  },
  { debounce: 500, maxWait: 1000 },
);

async function changeGiveawayStatus() {
  await updateOrCreateGiveaway.mutateAsync({
    description: formValue.value.description,
    followersLuck: formValue.value.followersLuck,
    keyword: formValue.value.keyword,
    randomNumberFrom: formValue.value.randomNumberFrom,
    randomNumberTo: formValue.value.randomNumberTo,
    subscribersLuck: formValue.value.subscribersLuck,
    winnersCount: formValue.value.winnersCount,
    isRunning: currentGiveaway ? !currentGiveaway.value?.isRunning : true,
    isFinished: false,
    type: formValue.value.type,
  });
}

const chooseWinner = useChooseGiveawayWinners();

async function rollGiveaway() {
  await chooseWinner.mutateAsync({
    giveawayId: giveawayId.value,
  });
  props.showWinnersModal();

}


async function finishGiveaway() {
  await updateOrCreateGiveaway.mutateAsync({
    isFinished: true,
    type: currentGiveaway.value?.type,
  });
}

const hasAccessToManageGiveaways = useUserAccessFlagChecker('MANAGE_GIVEAWAYS');

const isAbleToRoll = computed(
  () =>
    currentGiveaway.value &&
    !currentGiveaway.value.isFinished &&
    participantsCount.value &&
    participantsCount.value > 0 &&
		currentGiveaway.value?.winnerCount < participantsCount.value,
);
</script>
<template>
	<n-form style="height: 100%; min-height: 100%">
		<n-card
			:title="t('giveaways.settings.title')"
			content-style="padding: 0;"
			header-style="padding: 10px;"
			style="min-width: 300px; height: 100%; overflow-y: auto"
			segmented
		>
			<div style="padding: 10px">
				<n-input
					v-model:value="formValue.description"
					style="margin-bottom: 25px"
					type="textarea"
					placeholder="Giveaway description"
				/>
				<n-radio-group
					v-model:value="formValue.type"
					style="margin-bottom: 25px"
					name="giveawaysTypesGroup"
				>
					<div>
						<div style="margin-bottom: 7px">
							{{ t("giveaways.settings.giveawayType") }}
						</div>
						<n-space>
							<n-radio value="BY_KEYWORD" label="Keyword" />
							<n-radio value="BY_RANDOM_NUMBER" label="Random number" />
						</n-space>
					</div>
				</n-radio-group>
				<n-space vertical>
					<n-input
						v-if="formValue.type === 'BY_KEYWORD'"
						v-model:value="formValue.keyword"
						style="margin-bottom: 25px"
						:disabled="formValue.type !== 'BY_KEYWORD'"
						type="text"
						placeholder="Keyword Phrase"
					/>
					<div v-else style="margin-bottom: 25px">
						<n-input-number
							v-model:value="formValue.randomNumberFrom"
							style="margin-bottom: 10px"
							placeholder="Minimum number"
							min="0"
						/>
						<n-input-number
							v-model:value="formValue.randomNumberTo"
							placeholder="Maximum number"
							min="0"
						/>
					</div>
					<div style="margin-bottom: 25px">
						<div style="margin-bottom: 3px">
							{{ t("giveaways.settings.subscribersLuck") }}
						</div>
						<n-slider
							v-model:value="formValue.subscribersLuck"
							:step="1"
							:max="10"
							:min="0"
						/>
					</div>

					<div style="margin-bottom: 25px">
						<div style="margin-bottom: 3px">
							{{ t("giveaways.settings.followersLuck") }}
						</div>
						<n-slider
							v-model:value="formValue.followersLuck"
							:step="1"
							:max="10"
							:min="0"
						/>
					</div>

					<n-input-number
						v-model:value="formValue.requiredMinWatchTime"
						style="margin-bottom: 25px"
						placeholder="Minimum watch time"
						min="0"
					/>

					<n-input-number
						v-model:value="formValue.winnersCount"
						style="margin-bottom: 25px"
						placeholder="Winners count"
						min="1"
					/>

					<n-button
						type="primary"
						style="width: 100%"
						:disabled="!isAbleToRoll"
						@click="rollGiveaway"
					>
						{{
							isAbleToRoll
								? t("giveaways.settings.roll")
								: t("giveaways.settings.cantRoll")
						}}
					</n-button>
					<n-space />
					<n-button
						type="primary"
						style="width: 100%"
						:disabled="!hasAccessToManageGiveaways"
						@click="finishGiveaway"
					>
						{{ t("giveaways.settings.finish") }}
					</n-button>
				</n-space>
			</div>
			<template #header-extra>
				<n-button
					secondary
					type="primary"
					:disabled="!hasAccessToManageGiveaways"
					@click="changeGiveawayStatus"
				>
					{{
						currentGiveaway?.isRunning
							? t("giveaways.settings.running")
							: t("giveaways.settings.stopped")
					}}
				</n-button>
			</template>
		</n-card>
	</n-form>
</template>

<style scoped></style>
