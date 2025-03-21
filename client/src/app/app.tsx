/* dependencies */
import { Route, Router, Switch } from "wouter";

import { ProtectedRouter, UnauthorizedRoute } from "@/features/auth";
import { AuthPage } from "@/pages/auth-page";
import { ChatPage } from "@/pages/chat-page";
import { ChatsPage } from "@/pages/chats-page";
import { CreateChatPage } from "@/pages/create-chat-page";
import { ProfilePage } from "@/pages/profile-page";
import { NAVIGATION } from "@/shared/navigation";
import { Providers } from "./providers";
import { Confirm } from "@/shared/utils/confirm";

export const App = () => {
  return (
    <Providers>
      <main className="main max-w-[800px] mx-auto w-full">
        <Switch>
          <UnauthorizedRoute
            path="/auth"
            component={AuthPage}
            redirectPath={NAVIGATION.profile}
          />

          <ProtectedRouter base="/app">
            <Route path="/profile" component={ProfilePage} />

            <Router base="/chats">
              <Route path="/" component={ChatsPage} />

              <Switch>
                <Route path="/create" component={CreateChatPage} />
                <Route path="/:chatId" component={ChatPage} />
              </Switch>
            </Router>

            <Confirm />
          </ProtectedRouter>
        </Switch>
      </main>
    </Providers>
  );
};
