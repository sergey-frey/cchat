import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { App } from "./app";

import "./global.css";

async function enableMocking() {
  if (import.meta.env.MODE !== "mock") {
    return;
  }

  const { worker } = await import("@/shared/api/mocks/browser");
  return worker.start();
}

enableMocking()
  .then(() => {
    createRoot(document.getElementById("root")!).render(
      <StrictMode>
        <App />
      </StrictMode>,
    );
  })
  .catch(() => {
    console.log("Mock server failed to start");
  });
