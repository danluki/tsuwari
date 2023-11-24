<script setup lang="ts">

import {
  NModal,
} from 'naive-ui';
import { ref } from 'vue';

import Chat from '@/components/giveaways/chat.vue';
import Settings from '@/components/giveaways/settings.vue';
import Users from '@/components/giveaways/users.vue';
import Modal from '@/components/giveaways/winnersModal.vue';

const showWinnersModal = ref(false);

</script>

<template>
	<div class="main-container">
		<div class="flex-container">
			<Users :show-winners-modal="() => showWinnersModal = true" />
		</div>
		<div class="flex-container">
			<Settings :show-winners-modal="() => showWinnersModal = true" />
		</div>
		<div class="flex-container">
			<Chat />
		</div>
	</div>

	<n-modal
		v-model:show="showWinnersModal"
		:mask-closable="false"
		title="Giveaway winners"
		style="width: 700px; height: 500px"
		preset="card"
	>
		<modal />
	</n-modal>
</template>

<style scoped>
.flex-container {
	flex: 1;
	width: 100%;
	height: 85vh;
	margin: 10px;
}

.main-container {
	display: flex;
	align-items: center;
	justify-content: center;
	max-width: 100%;
	margin: 0 auto;
}

@media screen and (max-width: 768px) {
	.main-container {
		flex-direction: column;
		margin: 8px 0;
	}
}
</style>
