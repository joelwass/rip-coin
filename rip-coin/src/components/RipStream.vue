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
                  <td><v-btn color="primary" dark @click.stop="() => toggleVoteDialog(props.item)">Vote!</v-btn></td>
                </template>
              </v-data-table>
            </v-flex>
          </v-layout>

          <v-dialog v-model="voteDialog">
            <v-card>
              <v-card-title>
                Vote on rip
              </v-card-title>

              <v-card-text>
                <v-layout row>
                  <v-flex xs6>
                    <v-btn fab dark @click.native="() => vote(true)">
                      <v-icon dark>thumb_up</v-icon>
                    </v-btn>
                  </v-flex>

                  <v-flex xs6>
                    <v-btn fab dark @click.native="() => vote(false)">
                      <v-icon dark>thumb_down</v-icon>
                    </v-btn>
                  </v-flex>
                </v-layout>
              </v-card-text>

              <v-card-actions>
                <v-btn color="primary" flat @click.stop="voteDialog = false">Close</v-btn>
              </v-card-actions>
            </v-card>
          </v-dialog>
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
        const res = JSON.parse(e.data);
        const payload = JSON.parse(res.payload);
        switch (res.label) {
          case 'pub_key':
            self.key = payload.Pub;
            break;
          case 'new_tx':
            self.ripStream.unshift(payload);
            break;
          default:
            throw new Error(`Unhandled type ${res.label}`);
        }
      },
    });
  },
  data() {
    return {
      key: '',
      selectedTx: {},
      voteDialog: false,
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
    vote(approval) {
      const { selectedTx } = this;
      if (!selectedTx.rip.votes) selectedTx.rip.votes = [];

      // Add your vote
      selectedTx.rip.votes.push({
        address: this.key,
        approval,
      });

      this.sock.send({
        type: 'vote',
        data: selectedTx,
      });
    },
    toggleVoteDialog(tx) {
      this.selectedTx = tx;
      this.voteDialog = !this.voteDialog;
    },
  },
};
</script>

<style>
  .container {
    padding: 10px;
  }
</style>
