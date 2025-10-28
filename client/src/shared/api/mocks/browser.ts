import { setupWorker } from "msw/browser";
import { mockApi } from "./handlers";

export const worker = setupWorker(...mockApi.handlers);
