// src/providers/AuthUserProvider.ts
import { type InjectionKey, reactive, readonly, provide, inject } from 'vue'

export type AuthUser = {
  id: number
  username: string
  email: string
}

type AuthUserContext = {
  user: Readonly<{ value: AuthUser | null }>
  setUser: (newUser: AuthUser) => void
  clearUser: () => void
}

const AuthUserSymbol: InjectionKey<AuthUserContext> = Symbol('AuthUser')

export function createAuthUserProvider(initialUser: AuthUser | null = null) {
  const state = reactive<{ value: AuthUser | null }>({
    value: initialUser,
  })

  const setUser = (newUser: AuthUser) => {
    state.value = newUser
  }

  const clearUser = () => {
    state.value = null
  }

  provide(AuthUserSymbol, {
    user: readonly(state),
    setUser,
    clearUser,
  })
}

export function useAuthUser(): AuthUserContext {
  const context = inject(AuthUserSymbol)
  if (!context) {
    throw new Error('useAuthUser must be used within an AuthUserProvider.')
  }
  return context
}
