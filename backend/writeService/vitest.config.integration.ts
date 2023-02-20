// vitest.config.integration.ts
import { defineConfig } from 'vitest/config'

export default defineConfig({
  test: {
    include: ['test/**/*.test.ts'],
    threads: false,
    setupFiles: ['src/helpers/setup.util.ts'],
  },
//   resolve: {
//     alias: {
//       auth: '/src/auth',
//       quotes: '/src/quotes',
//       lib: '/src/lib'
//     }
//   }
})