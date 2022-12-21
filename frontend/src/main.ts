import { createApp } from 'vue'
import App from './App.vue'
import vuetify from './plugins/vuetify'
import { loadFonts } from './plugins/webfontloader'
import { aliases, mdi } from 'vuetify/iconsets/mdi'

loadFonts()
createApp(App)
  .use(vuetify,{
      icons: {
          defaultSet: 'mdi',
          aliases,
          sets: {
              mdi,
          }
      }})
  .mount('#app')
