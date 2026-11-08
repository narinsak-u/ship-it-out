"use client";

import { usePathname } from "next/navigation";
import { type FormEvent, useEffect, useRef, useState } from "react";
import { useAuth } from "./AuthProvider";

export default function AuthControls() {
  const pathname = usePathname();
  const { user, isStaff, loading, login, signup, logout } = useAuth();
  const [open, setOpen] = useState(false);
  const [mode, setMode] = useState<"login" | "signup">("login");
  const [error, setError] = useState("");
  const [feedback, setFeedback] = useState("");
  const [pending, setPending] = useState(false);
  const dialog = useRef<HTMLDialogElement>(null);

  useEffect(() => {
    if (open && dialog.current && !dialog.current.open) dialog.current.showModal();
    if (!open && dialog.current?.open) dialog.current.close();
  }, [open]);

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
    <>
      <nav>
        <a
          href="/"
          className={pathname === "/" ? "active" : undefined}
          aria-current={pathname === "/" ? "page" : undefined}
        >
          Home
        </a>
        <a
          href="/orders"
          className={pathname.startsWith("/orders") ? "active" : undefined}
          aria-current={pathname.startsWith("/orders") ? "page" : undefined}
        >
          Orders
        </a>
        <a
          href="/carriers"
          className={pathname.startsWith("/carriers") ? "active" : undefined}
          aria-current={pathname.startsWith("/carriers") ? "page" : undefined}
        >
          Carriers
        </a>
        {loading ? (
          <span className="signin" role="status">
            Checking session…
          </span>
        ) : user ? (
          <span className="auth-status">
            <span className="signin" aria-label={`Signed in as ${user.role}`}>
              {user.role} ({user.name}){isStaff ? " · staff access" : " · read-only"}
            </span>
            <button className="auth-button" type="button" onClick={signOut} disabled={pending}>
              Sign out
            </button>
          </span>
        ) : (
          <button className="auth-button signin" type="button" onClick={() => setOpen(true)}>
            Sign in
          </button>
        )}
      </nav>
      {feedback && <span role="alert">{feedback}</span>}
      <dialog
        ref={dialog}
        className="auth-dialog"
        aria-labelledby="auth-title"
        onClose={() => setOpen(false)}
      >
        <div className="auth-dialog-header">
          <h2 id="auth-title">{mode === "login" ? "Sign in" : "Create account"}</h2>
          <button
            className="auth-close"
            type="button"
            aria-label="Close"
            onClick={() => setOpen(false)}
          >
            ×
          </button>
        </div>
        <div className="auth-tabs" aria-label="Authentication">
          <button
            type="button"
            aria-pressed={mode === "login"}
            onClick={() => {
              setMode("login");
              setError("");
            }}
          >
            Sign in
          </button>
          <button
            type="button"
            aria-pressed={mode === "signup"}
            onClick={() => {
              setMode("signup");
              setError("");
            }}
          >
            Sign up
          </button>
        </div>
        <form className="auth-form" key={mode} onSubmit={submit}>
          {mode === "signup" && (
            <label>
              Name
              <input name="name" type="text" autoComplete="name" required />
            </label>
          )}
          <label>
            Email
            <input
              name="email"
              type="email"
              defaultValue={mode === "login" ? "demo@shipitout.example" : ""}
              required
            />
          </label>
          <label>
            Password
            <input
              name="password"
              type="password"
              autoComplete={mode === "login" ? "current-password" : "new-password"}
              defaultValue={mode === "login" ? "demo-read-only-2026" : ""}
              required
            />
          </label>
          {error && (
            <p className="auth-error" role="alert">
              {error}
            </p>
          )}
          <button className="auth-submit" type="submit" disabled={pending}>
            {pending ? "Please wait…" : mode === "login" ? "Sign in" : "Create account"}
          </button>
          {mode === "login" && (
            <p className="demo-hint">
              Read-only demo account. Email: demo@shipitout.example · Password: demo-read-only-2026
            </p>
          )}
        </form>
      </dialog>
    </>
  );
}
