export default class Socket {
  constructor(config) {
    const ws = new WebSocket(config.url || 'ws://localhost:6969/ws');

    ws.onopen = config.opened;
    ws.onmessage = config.onmessage;
    ws.onclose = config.closed;
    this.ws = ws;
  }

  send(msg) {
    let payload = msg;
    if (typeof msg === 'object') {
      payload = JSON.stringify(msg);
    }

    this.ws.send(payload);
  }
}
