<template>
  <div class="container">
    <v-card flat>
      <v-card-text>
        <v-container fluid>
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

          <v-layout row>
            <v-flex xs12>
              <v-data-table
                :items="ripStream"
                class="elevation-1"
                hide-actions
                hide-headers
              >
              <template slot="items" slot-scope="props">
                <td>{{ props.item.id }}</td>
                <td>{{ props.item.rip.rip }}</td>
              </template>
            </v-data-table>
            </v-flex>
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
        self.ripStream.unshift(JSON.parse(e.data));
      },
    });
  },
  data() {
    return {
      ripStream: [],
      rip: '',
      ripAddress: '',
      msg: '',
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
