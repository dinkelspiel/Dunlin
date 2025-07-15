// src/providers/AuthUserProvider.ts
import { type InjectionKey, reactive, readonly, provide, inject } from 'vue'

export type AuthUser = {
  id: number
  username: string
  email: string
}

type AuthUserContext = {
  authUser: Readonly<{ value: AuthUser | null }>
  setAuthUser: (newUser: AuthUser | null) => void
  clearAuthUser: () => void
}

const AuthUserSymbol: InjectionKey<AuthUserContext> = Symbol('AuthUser')

export function createAuthUserProvider(initialUser: AuthUser | null = null) {
  const state = reactive<{ value: AuthUser | null }>({
    value: initialUser,
  })

  const setAuthUser = (newUser: AuthUser | null) => {
    state.value = newUser
  }

  const clearAuthUser = () => {
    state.value = null
  }

  provide(AuthUserSymbol, {
    authUser: readonly(state),
    setAuthUser,
    clearAuthUser,
  })
}

export function useAuthUser(): AuthUserContext {
  const context = inject(AuthUserSymbol)
  if (!context) {
    throw new Error('useAuthUser must be used within an AuthUserProvider.')
  }
  return context
}
