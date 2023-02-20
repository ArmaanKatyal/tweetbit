// vitest.config.unit.ts
import { defineConfig } from 'vitest/config'

export default defineConfig({
  test: {

    include: [
      'test/**/*.test.ts',
      '!test'
    ]
  },
//   resolve: {
//     alias: {
//       auth: '/src/auth',
//       quotes: '/src/quotes',
//       lib: '/src/lib'
//     }
//   }
})