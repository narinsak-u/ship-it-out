"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { type FormEvent, useEffect, useState } from "react";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { cn } from "@/lib/utils";
import { useAuthStore } from "../store/auth-store";

export default function AuthControls() {
  const pathname = usePathname();
  const { user, isStaff, loading, login, signup, logout, initialize } = useAuthStore();
  const [open, setOpen] = useState(false);
  const [mode, setMode] = useState<"login" | "signup">("login");
  const [error, setError] = useState("");
  const [feedback, setFeedback] = useState("");
  const [pending, setPending] = useState(false);

  useEffect(() => {
    void initialize();
  }, [initialize]);

  async function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setPending(true);
    setError("");
    const form = new FormData(event.currentTarget);
    try {
      if (mode === "login") {
        await login(String(form.get("email")), String(form.get("password")));
      } else {
        await signup(
          String(form.get("name")),
          String(form.get("email")),
          String(form.get("password")),
        );
      }
      setOpen(false);
    } catch (cause) {
      setError(cause instanceof Error ? cause.message : "Authentication failed");
    } finally {
      setPending(false);
    }
  }

  async function signOut() {
    setPending(true);
    try {
      await logout();
      setFeedback("");
    } catch {
      setFeedback("Sign out failed. Please try again.");
    } finally {
      setPending(false);
    }
  }

  return (
    <Dialog
      open={open}
      onOpenChange={(nextOpen) => {
        setOpen(nextOpen);
        if (!nextOpen) setError("");
      }}
    >
      <nav aria-label="Main navigation" className="flex flex-wrap items-center justify-end gap-1">
        {[
          { label: "Home", href: "/", active: pathname === "/" },
          { label: "Orders", href: "/orders", active: pathname.startsWith("/orders") },
          { label: "Carriers", href: "/carriers", active: pathname.startsWith("/carriers") },
        ].map(({ label, href, active }) => (
          <Link
            key={href}
            href={href}
            aria-current={active ? "page" : undefined}
            className={cn(
              "inline-flex min-h-11 items-center rounded-md px-3 text-sm text-muted-foreground transition-colors hover:bg-muted hover:text-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring",
              active && "bg-secondary text-foreground",
            )}
          >
            {label}
          </Link>
        ))}
        {loading ? (
          <span className="px-3 py-2 text-sm text-muted-foreground" role="status">
            Checking session…
          </span>
        ) : user ? (
          <span className="flex flex-wrap items-center justify-end gap-2">
            <span
              className="px-2 text-xs text-muted-foreground"
              aria-label={`Signed in as ${user.role}`}
            >
              {user.role} ({user.name}){isStaff ? " · staff access" : " · read-only"}
            </span>
            <Button type="button" variant="outline" onClick={signOut} disabled={pending}>
              Sign out
            </Button>
          </span>
        ) : (
          <DialogTrigger asChild>
            <Button type="button" variant="outline">
              Sign in
            </Button>
          </DialogTrigger>
        )}
      </nav>
      {feedback && (
        <p className="m-0 px-4 pb-2 text-sm text-destructive" role="alert">
          {feedback}
        </p>
      )}
      <DialogContent className="max-h-[calc(100dvh-2rem)] overflow-y-auto sm:max-w-md">
        <DialogHeader>
          <DialogTitle>{mode === "login" ? "Sign in" : "Create account"}</DialogTitle>
          <DialogDescription>
            {mode === "login"
              ? "Access your Ship It Out account."
              : "Create a read-only customer account."}
          </DialogDescription>
        </DialogHeader>
        <div className="flex gap-2" aria-label="Authentication">
          <Button
            type="button"
            variant={mode === "login" ? "secondary" : "ghost"}
            aria-pressed={mode === "login"}
            onClick={() => {
              setMode("login");
              setError("");
            }}
          >
            Sign in
          </Button>
          <Button
            type="button"
            variant={mode === "signup" ? "secondary" : "ghost"}
            aria-pressed={mode === "signup"}
            onClick={() => {
              setMode("signup");
              setError("");
            }}
          >
            Sign up
          </Button>
        </div>
        <form key={mode} className="grid gap-4" onSubmit={submit}>
          {mode === "signup" && (
            <div className="grid gap-2">
              <Label htmlFor="auth-name">Name</Label>
              <Input id="auth-name" name="name" type="text" autoComplete="name" required />
            </div>
          )}
          <div className="grid gap-2">
            <Label htmlFor="auth-email">Email</Label>
            <Input
              id="auth-email"
              name="email"
              type="email"
              autoComplete="email"
              defaultValue={mode === "login" ? "demo@shipitout.example" : ""}
              required
            />
          </div>
          <div className="grid gap-2">
            <Label htmlFor="auth-password">Password</Label>
            <Input
              id="auth-password"
              name="password"
              type="password"
              autoComplete={mode === "login" ? "current-password" : "new-password"}
              defaultValue={mode === "login" ? "demo-read-only-2026" : ""}
              required
            />
          </div>
          {error && (
            <Alert variant="destructive">
              <AlertDescription>{error}</AlertDescription>
            </Alert>
          )}
          <Button type="submit" disabled={pending}>
            {pending ? "Please wait…" : mode === "login" ? "Sign in" : "Create account"}
          </Button>
          {mode === "login" && (
            <p className="m-0 rounded-md border border-border bg-secondary p-3 text-sm leading-relaxed text-muted-foreground">
              Read-only demo account. Email: demo@shipitout.example · Password: demo-read-only-2026
            </p>
          )}
        </form>
      </DialogContent>
    </Dialog>
  );
}
