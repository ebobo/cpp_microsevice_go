<template>
  <v-container>
    <h3 class="ma-4 primary--text">Calculator</h3>
    <h3 class="ma-4 blue-grey--text">parameters :</h3>

    <v-row class="ma-2">
      <v-col cols="6">
        <v-text-field
          v-model="numberA"
          class="mt-0 pt-0"
          hide-details
          single-line
          type="number"
          style="height: 60px"
          min="0"
          max="100"
        ></v-text-field>
      </v-col>
      <v-col cols="6">
        <v-text-field
          v-model="numberB"
          class="mt-0 pt-0"
          hide-details
          single-line
          type="number"
          style="height: 60px"
          min="0"
          max="100"
        ></v-text-field>
      </v-col>
    </v-row>
    <v-row class="ma-2">
      <v-btn color="primary" @click="send">Send</v-btn>
    </v-row>

    <v-row class="ma-2">
      <h3 class="ma-2 mt-8 blue-grey--text">{{ 'result : ' + result }}</h3>
    </v-row>
  </v-container>
</template>

<script lang="ts">
import Vue from 'vue';
import { setParameters } from '../services/data';

export default Vue.extend({
  components: {},
  data(): {
    numberA: number;
    numberB: number;
    result: number;
  } {
    return {
      numberA: 0,
      numberB: 0,
      result: 0,
    };
  },

  methods: {
    send() {
      setParameters({ A: this.numberA, B: this.numberB })
        .then((response) => console.log(response))
        .catch((error) => {
          console.log(error);
        });
    },
  },
});
</script>
