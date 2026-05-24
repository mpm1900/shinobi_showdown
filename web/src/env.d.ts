type AppEnv = {
  readonly API_URL
}

declare global {
  namespace NodeJS {
    interface ProcessEnv extends AppEnv {}
  }
}

export type { AppEnv }
