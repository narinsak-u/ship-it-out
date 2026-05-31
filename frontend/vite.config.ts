import { defineConfig } from "vitest/config";
import vue from "@vitejs/plugin-vue";
import tailwindcss from "@tailwindcss/vite";
import path from "path";

export default defineConfig({
  plugins: [vue(), tailwindcss()],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
  },
  build: {
    rollupOptions: {
      output: {
        manualChunks: {
          vendor: [
            "vue",
            "vue-router",
            "pinia",
            "@tanstack/vue-query",
            "lucide-vue-next",
            "vue-sonner",
          ] as string[],
          ui: [
            "radix-vue",
            "reka-ui",
            "cmdk-vue",
            "vaul-vue",
            "class-variance-authority",
          ] as string[],
          maps: ["leaflet"] as string[],
        } as Record<string, string[]>,
      },
    },
  },
  test: {
    environment: "happy-dom",
    globals: true,
    include: ["src/**/*.spec.ts"],
    setupFiles: ["tests/setup.ts"],
    coverage: {
      provider: "v8",
      include: ["src/composables/**", "src/stores/**", "src/lib/**", "src/hooks/**"],
      exclude: ["src/components/ui/**", "src/main.ts", "src/vite-env.d.ts"],
      reporter: ["text", "html", "lcov"],
    },
  },
});
