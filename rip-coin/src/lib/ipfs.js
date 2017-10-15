import IPFS from 'ipfs';

export default {
  new() {
    return new Promise((resolve, reject) => {
      // ipfs pubsub pub rip-coin-tx "yo dawg死"
      const node = new IPFS();

      node.on('ready', () => {
        const multihashStr = 'QmS4ustL54uo8FzR9455qaxZwuMiUhyvMcX9Ba8nUH4uVv';
        node.files.cat(multihashStr, (err, file) => {
          if (err) {
            reject(err);
          }

          let data = '';
          file.on('data', (d) => {
            data += d;
          });
          file.on('end', () => {
            return resolve(data);
          });
        });
      });
    });
  },
};
