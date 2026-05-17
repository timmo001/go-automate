import type { CliRenderer } from "@opentui/core";
import { Effect, Schema } from "effect";
import type { NotifyConfig } from "../types.js";
import type { Toast } from "../tui/Toast.js";

const log = (msg: string) =>
  console.error(`[go-automate-tui:CommandRunner] ${msg}`);

/** Typed error for command execution failures */
export class CommandRunnerError extends Schema.TaggedErrorClass<CommandRunnerError>()(
  "CommandRunnerError",
  {
    message: Schema.String,
  },
) {}

/** Service for executing shell commands with TUI suspend/resume lifecycle */
export interface CommandRunnerService {
  /** Suspend the TUI, run the command with inherited stdio, then resume.
   *  When wait is true, shows "Press any key to continue" before resuming. */
  readonly runSuspended: (cmd: string, wait: boolean) => Promise<void>;

  /** Run a command in the background without suspending the TUI.
   *  Returns immediately; fails with {@link CommandRunnerError} on non-zero exit. */
  readonly runSilent: (cmd: string) => Effect.Effect<void, CommandRunnerError>;

  /** Run a command silently with toast notifications for progress and result.
   *  Errors are recovered internally via toast — the returned Effect always succeeds. */
  readonly runNotify: (cmd: string, notify: NotifyConfig) => Effect.Effect<void>;

  /** Replace the TUI process with a command (no return). */
  readonly runReplace: (cmd: string) => Promise<never>;
}

/** Create a {@link CommandRunnerService} bound to the given renderer for suspend/resume */
export function createCommandRunner(
  renderer: CliRenderer,
  toast: Toast,
): CommandRunnerService {
  return {
    runSuspended: async (cmd, wait) => {
      log(`Suspending for: ${cmd}`);
      renderer.suspend();
      renderer.currentRenderBuffer.clear();

      try {
        const cols = process.stdout.columns || 80;
        const label = ` ${cmd} `;
        const pad = Math.max(0, cols - label.length);
        const left = Math.floor(pad / 2);
        const right = pad - left;
        const header =
          "\x1b[90m" +
          "─".repeat(left) +
          "\x1b[0m\x1b[1m" +
          label +
          "\x1b[0m\x1b[90m" +
          "─".repeat(right) +
          "\x1b[0m";
        process.stdout.write(`\n\n${header}\n\n`);

        const proc = Bun.spawn(["bash", "-c", cmd], {
          stdin: "inherit",
          stdout: "inherit",
          stderr: "inherit",
        });
        await proc.exited;

        if (wait) {
          process.stdout.write(
            "\n\x1b[90mPress any key to continue...\x1b[0m",
          );
          await new Promise<void>((resolve) => {
            const wasRaw = process.stdin.isRaw;
            if (process.stdin.isTTY) process.stdin.setRawMode(true);
            process.stdin.resume();
            process.stdin.once("data", () => {
              if (process.stdin.isTTY) process.stdin.setRawMode(wasRaw);
              process.stdin.pause();
              resolve();
            });
          });
        }
      } finally {
        renderer.currentRenderBuffer.clear();
        renderer.resume();
        renderer.requestRender();
        log("Resumed after command");
      }
    },

    runSilent: Effect.fn("CommandRunner.runSilent")(function* (
      cmd: string,
    ): Effect.fn.Return<void, CommandRunnerError> {
      log(`Running silently: ${cmd}`);

      const { exitCode, stderr } = yield* Effect.tryPromise({
        try: async () => {
          const proc = Bun.spawn(["bash", "-c", cmd], {
            stdout: "pipe",
            stderr: "pipe",
          });
          const exitCode = await proc.exited;
          const stderr = await new Response(proc.stderr).text();
          return { exitCode, stderr };
        },
        catch: (error) =>
          new CommandRunnerError({
            message: error instanceof Error ? error.message : String(error),
          }),
      });

      if (exitCode !== 0) {
        return yield* Effect.fail(
          new CommandRunnerError({
            message: stderr.trim().split("\n")[0] || `Command failed (exit ${exitCode})`,
          }),
        );
      }

      log(`Silent command completed: ${cmd}`);
    }),

    runNotify: Effect.fn("CommandRunner.runNotify")(function* (
      cmd: string,
      notify: NotifyConfig,
    ): Effect.fn.Return<void, never> {
      log(`Running with notification: ${cmd}`);
      toast.show(notify.id, notify.progress, "info");

      const result = yield* Effect.tryPromise({
        try: async () => {
          const proc = Bun.spawn(["bash", "-c", cmd], {
            stdout: "pipe",
            stderr: "pipe",
          });
          const exitCode = await proc.exited;
          const stderr = await new Response(proc.stderr).text();
          return { exitCode, stderr };
        },
        catch: (error) =>
          new CommandRunnerError({
            message: error instanceof Error ? error.message : String(error),
          }),
      }).pipe(
        Effect.catch((err: CommandRunnerError) =>
          Effect.succeed({ exitCode: 1, stderr: err.message }),
        ),
      );

      if (result.exitCode !== 0) {
        const errMsg =
          result.stderr.trim().split("\n")[0] || "Command failed";
        log(`Notify command failed (exit ${result.exitCode}): ${result.stderr}`);
        toast.show(notify.id, errMsg, "error");
      } else {
        log(`Notify command completed: ${cmd}`);
        toast.show(notify.id, notify.success, "success");
      }
    }),

    runReplace: async (cmd) => {
      log(`Replacing process with: ${cmd}`);
      renderer.suspend();

      const proc = Bun.spawn(["bash", "-c", cmd], {
        stdin: "inherit",
        stdout: "inherit",
        stderr: "inherit",
      });
      const exitCode = await proc.exited;
      process.exit(exitCode ?? 0);
    },
  };
}
