type WSOptions = {
  path: string;
};

type WSEvent = { event: string };

export class WS<ClientToServerEventsMap extends WSEvent> {
  private _options: WSOptions;
  private _instance: WebSocket | null = null;

  constructor(options: WSOptions) {
    this._options = options;
  }

  public connect() {
    this._instance = new WebSocket(this._options.path);
  }

  public send(data: ClientToServerEventsMap) {
    this._instance?.send(JSON.stringify(data));
  }

  public onMessage(callback: (event: MessageEvent) => void) {
    this._instance?.addEventListener("message", callback);
  }
}
