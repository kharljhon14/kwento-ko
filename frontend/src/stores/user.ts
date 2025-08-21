import type { User } from '@/types/user';
import { create } from 'zustand';

interface State {
  user: undefined | User;
}

interface Action {
  setAuthState: (state: State['user']) => void;
}

export const useUserStore = create<State & Action>((set) => ({
  user: undefined,
  isAuthenticated: false,
  setAuthState: (user) => set({ user })
}));
