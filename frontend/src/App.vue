<template>
  <v-app>
    <v-main>
      <calculator />
    </v-main>
  </v-app>
</template>

<script lang="ts">
import Vue from 'vue';
import Calculator from './components/Calculator.vue';

export default Vue.extend({
  name: 'App',

  components: {
    Calculator,
  },

  data: () => ({
    //
  }),
});
</script>
