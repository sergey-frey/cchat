/* dependencies */
import { AnimatePresence } from "framer-motion";
import { useEffect, useRef } from "react";
import { Redirect, Route, Router, Switch } from "wouter";

import { ProtectedRouter, UnauthorizedRoute } from "@/features/auth";
import { BottomNavigation, SlidePageAnimation } from "@/features/navigation";
import { AuthPage } from "@/pages/auth-page";
import { ChatPage } from "@/pages/chat-page";
import { ChatsPage } from "@/pages/chats-page";
import { CreateChatPage } from "@/pages/create-chat-page";
import { EditProfilePage } from "@/pages/edit-profile-page";
import { ProfilePage } from "@/pages/profile-page";
import { NAVIGATION } from "@/shared/navigation";
import {
  setContainerRefSelector,
  useAppContainer,
} from "@/shared/utils/app-container";
import { Confirm } from "@/shared/utils/confirm";
import { Providers } from "./providers";

export const App = () => {
  const containerRef = useRef<HTMLElement>(null);
  const setContainerRef = useAppContainer(setContainerRefSelector);

  useEffect(() => {
    setContainerRef(containerRef);
  }, [setContainerRef]);

  return (
    <Providers>
      <main
        className="relative main max-w-[800px] mx-auto w-full overflow-x-hidden"
        ref={containerRef}
      >
        <AnimatePresence mode="wait">
          <Switch>
            <Route
              path="/"
              component={() => <Redirect to={NAVIGATION.profile} />}
            />

            <UnauthorizedRoute
              path="/auth"
              component={AuthPage}
              redirectPath={NAVIGATION.profile}
            />

            <ProtectedRouter base="/app">
              <SlidePageAnimation>
                <Router base="/profile">
                  <Route path="/edit">
                    <EditProfilePage />
                  </Route>

                  <Route path="/">
                    <ProfilePage />
                  </Route>
                </Router>
              </SlidePageAnimation>

              <Router base="/chats">
                <Route path="/" component={ChatsPage} />

                <Switch>
                  <Route path="/create" component={CreateChatPage} />
                  <Route path="/:chatId" component={ChatPage} />
                </Switch>
              </Router>

              <Confirm />

              <BottomNavigation />
            </ProtectedRouter>
          </Switch>
        </AnimatePresence>
      </main>
    </Providers>
  );
};
