<script setup lang="ts">

import { useCurrentGiveaway } from '@/api/index.js';

</script>

<template>
	<div class="main-container">
		<div class="flex-container">
			<Users />
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
