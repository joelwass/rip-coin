<template>
  <div class="container">
    <v-card flat>
      <v-card-text>
        <v-container fluid>
          {{msg}}
          <v-layout row>
            <v-flex xs4>
              <v-subheader>New Rip</v-subheader>
            </v-flex>
            <v-flex xs4>
              <v-text-field
                name="ripper"
                label="Ripper Address"
                id="ripper"
                v-model="ripAddress"
              ></v-text-field>
            </v-flex>
            <v-flex xs4>
              <v-text-field
                name="rip"
                label="Rip"
                id="rip"
                v-model="rip"
              ></v-text-field>
            </v-flex>
            <v-btn
              @click.native="addRip"
            >
              Add Rip
            </v-btn>
          </v-layout>
        </v-container>
      </v-card-text>
    </v-card>
  </div>
</template>

<script>
import Socket from '@/lib/socket';

export default {
  name: 'RipStream',
  mounted() {
    const self = this;
    this.sock = new Socket({
      opened() {
      },
      closed() {
      },
      onmessage(e) {
        self.msg = e.data;
      },
    });
  },
  data() {
    return {
      rip: '',
      ripAddress: '',
      msg: 'Welcome to Your Vue.js App',
    };
  },
  methods: {
    addRip() {
      const { rip, ripAddress } = this;
      this.sock.send({
        type: 'new_rip',
        data: {
          TotalAmount: 1,
          RipperPublicKey: ripAddress,
          Rip: {
            Rip: rip,
          },
        },
      });
    },
  },
};
</script>

<style>
  .container {
    padding: 10px;
  }
</style>
