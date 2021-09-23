<template>
  <v-container>
    <h3 class="ma-2 primary--text">Calculator</h3>
    <h3 class="ma-4 blue-grey--text">parameters :</h3>

    <v-row class="ma-2">
      <v-col cols="3">
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
      <v-col cols="3">
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
      <v-col cols="3">
        <v-btn
          color="primary"
          :disabled="isWebSocketConnectionOpen"
          @click="send"
          >Send</v-btn
        >
      </v-col>
    </v-row>
    <v-row class="ma-2">
      <v-col cols="8">
        <v-progress-linear v-model="progress"></v-progress-linear>
      </v-col>
    </v-row>

    <v-row class="ma-2">
      <h3 class="ma-2 mt-8 blue-grey--text">{{ 'result : ' + result }}</h3>
    </v-row>
  </v-container>
</template>

<script lang="ts">
import Vue from 'vue';
import { setParameters, CalcData } from '../services/data';

export default Vue.extend({
  components: {},
  data(): {
    numberA: string;
    numberB: string;
    progress: number;
    websocketConnection: null | WebSocket;
    isWebSocketConnectionOpen: boolean;
    result: number;
  } {
    return {
      numberA: '0',
      numberB: '0',
      progress: 0,
      websocketConnection: null,
      isWebSocketConnectionOpen: false,
      result: 0,
    };
  },

  methods: {
    //send the parameters
    send() {
      const data: CalcData = {
        A: parseInt(this.numberA),
        B: parseInt(this.numberB),
      };
      setParameters(data)
        .then((response) => this.parametersSended(response))
        .catch((error) => {
          console.log(error);
        });
    },

    //server got the parameter
    parametersSended(data: any) {
      this.progress = 0;

      if (!this.websocketConnection) {
        console.log('make a new connection');
        this.websocketConnection = new WebSocket(
          process.env.VUE_APP_WS_BASE_PATH
        );

        this.websocketConnection.onopen = () => {
          console.log('Connected to WebSocket server');
          this.isWebSocketConnectionOpen = true;
        };

        this.websocketConnection.onclose = () => {
          console.log('Disconnected with WebSocket server');
          this.isWebSocketConnectionOpen = false;
          this.progress = 100;
          this.websocketConnection = null;
        };

        this.websocketConnection.onmessage = (event) => {
          this.progress = event?.data;
        };
      }
    },
  },
});
</script>
