<script setup lang="ts">

import Chat from '@/components/giveaways/chat.vue';
import Settings from '@/components/giveaways/settings.vue';
import Users from '@/components/giveaways/users.vue';

</script>

<template>
	<div class="main-container">
		<div class="flex-container">
			<Users />
		</div>
		<div class="flex-container">
			<Settings />
		</div>
		<div class="flex-container">
			<Chat />
		</div>
	</div>
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
