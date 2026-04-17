import { DOCUMENT } from '@angular/common';
import { Injectable, computed, inject, signal } from '@angular/core';

export type AppTheme = 'light' | 'dark';

@Injectable({ providedIn: 'root' })
export class ThemeService {
  private readonly storageKey = 'vitatrack-theme';
  private readonly document = inject(DOCUMENT);
  private readonly themeState = signal<AppTheme>(this.resolveInitialTheme());

  readonly theme = this.themeState.asReadonly();
  readonly isDark = computed(() => this.themeState() === 'dark');
  readonly toggleLabel = computed(() =>
    this.isDark() ? 'Switch to light theme' : 'Switch to dark theme'
  );

  constructor() {
    this.applyTheme(this.themeState());
  }

  toggleTheme(): void {
    this.setTheme(this.isDark() ? 'light' : 'dark');
  }

  setTheme(theme: AppTheme): void {
    this.themeState.set(theme);
    this.applyTheme(theme);
    this.persistTheme(theme);
  }

  private resolveInitialTheme(): AppTheme {
    const storedTheme = this.readStoredTheme();
    if (storedTheme) {
      return storedTheme;
    }

    return this.prefersDarkTheme() ? 'dark' : 'light';
  }

  private applyTheme(theme: AppTheme): void {
    const root = this.document.documentElement;
    root.setAttribute('data-theme', theme);
    root.style.colorScheme = theme;
  }

  private persistTheme(theme: AppTheme): void {
    try {
      localStorage.setItem(this.storageKey, theme);
    } catch {
      // Ignore storage failures and keep theme in memory.
    }
  }

  private readStoredTheme(): AppTheme | null {
    try {
      const value = localStorage.getItem(this.storageKey);
      return value === 'light' || value === 'dark' ? value : null;
    } catch {
      return null;
    }
  }

  private prefersDarkTheme(): boolean {
    return typeof window !== 'undefined'
      && typeof window.matchMedia === 'function'
      && window.matchMedia('(prefers-color-scheme: dark)').matches;
  }
}
