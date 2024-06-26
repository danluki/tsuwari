import { QueryClient, VueQueryPlugin } from '@tanstack/vue-query'
import urql from '@urql/vue'
import { createApp } from 'vue'

import MainApp from './app.vue'
import { router } from './router'

import './assets/index.css'
import { urqlClient } from '@/urql-client'

const app = createApp(MainApp)
	.use(urql, urqlClient)
	.use(router)

VueQueryPlugin.install(app, {
	queryClient: new QueryClient({
		defaultOptions: {
			queries: {
				refetchOnWindowFocus: false,
				refetchOnMount: false,
				refetchOnReconnect: false,
				retry: false,
			},
		},
	}),
})

app.mount('#app')
