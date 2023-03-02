<template>
  <v-app>
  <v-main>
    <router-view/>
    <messages style="height:100vh;width: 100%;position: absolute"></messages>
  </v-main>
  </v-app>
</template>

<script>
import messages from '@/components/Message'

export default {
  components: {
    'messages':messages,
  },
  data:()=>({
  }),
}
</script>
