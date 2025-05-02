import { http, HttpHandler } from "msw";
import urlJoin from "url-join";

export class MockAPI {
  private _baseUrl = "";
  private _handlers: HttpHandler[] = [];

  constructor(baseUrl: string) {
    this._baseUrl = baseUrl;
  }

  public get(...[path, resolver]: Parameters<typeof http.get>) {
    this._handlers.push(
      http.get(urlJoin(this._baseUrl, path.toString()), resolver),
    );
  }

  public post(...[path, resolver]: Parameters<typeof http.post>) {
    this._handlers.push(
      http.post(urlJoin(this._baseUrl, path.toString()), resolver),
    );
  }

  public get handlers() {
    return this._handlers;
  }
}
