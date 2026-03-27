"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { useAuthStore } from "@/stores/auth.store";
import Link from "next/link";

export default function HomePage() {
  const [todoText, setTodoText] = useState("");
  const { isAuthenticated, user, logout } = useAuthStore();

  const handleAdd = () => {
    if (!todoText.trim()) return;
    // TODO: wire up to API
    setTodoText("");
  };

  return (
    <div className="min-h-screen flex flex-col">
      {/* Top bar */}
      <header className="border-b">
        <div className="container mx-auto flex items-center justify-between h-16 px-4">
          <h1 className="text-2xl font-bold tracking-tight">
            Lets<span className="text-primary">Do</span>
          </h1>
          <div className="flex items-center gap-3">
            {isAuthenticated ? (
              <>
                <span className="text-sm text-muted-foreground">
                  {user?.email}
                </span>
                <Button variant="ghost" size="sm" onClick={logout}>
                  Sign out
                </Button>
              </>
            ) : (
              <>
                <Button variant="ghost" size="sm" asChild>
                  <Link href="/login">Log in</Link>
                </Button>
                <Button size="sm" asChild>
                  <Link href="/signup">Sign up</Link>
                </Button>
              </>
            )}
          </div>
        </div>
      </header>

      {/* Main content - centered input */}
      <main className="flex-1 flex items-center justify-center px-4">
        <div className="w-full max-w-md space-y-6">
          <div className="text-center space-y-2">
            <h2 className="text-4xl font-bold tracking-tight">
              Lets<span className="text-primary">Do</span>
            </h2>
            <p className="text-muted-foreground">
              What needs to be done?
            </p>
          </div>
          <div className="flex gap-2">
            <Input
              placeholder="Add a todo..."
              value={todoText}
              onChange={(e) => setTodoText(e.target.value)}
              onKeyDown={(e) => e.key === "Enter" && handleAdd()}
              className="h-12 text-base"
            />
            <Button onClick={handleAdd} className="h-12 px-6">
              Add
            </Button>
          </div>
        </div>
      </main>
    </div>
  );
}
