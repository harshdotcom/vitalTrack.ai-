import { Injectable } from '@angular/core';

export interface GoogleCredentialResponse {
  credential: string;
  select_by: string;
}

interface GoogleAccountsIdConfiguration {
  callback: (response: GoogleCredentialResponse) => void;
  cancel_on_tap_outside?: boolean;
  client_id: string;
  context?: 'signin' | 'signup' | 'use';
  ux_mode?: 'popup' | 'redirect';
}

interface GoogleRenderButtonOptions {
  logo_alignment?: 'left' | 'center';
  shape?: 'pill' | 'rectangular';
  size?: 'large' | 'medium' | 'small';
  text?: 'continue_with' | 'signin_with' | 'signup_with';
  theme?: 'outline' | 'filled_blue' | 'filled_black';
  type?: 'standard' | 'icon';
  width?: number;
}

interface GoogleAccountsIdAPI {
  initialize(config: GoogleAccountsIdConfiguration): void;
  renderButton(parent: HTMLElement, options: GoogleRenderButtonOptions): void;
}

declare global {
  interface Window {
    google?: {
      accounts: {
        id: GoogleAccountsIdAPI;
      };
    };
  }
}

@Injectable({
  providedIn: 'root'
})
export class GoogleIdentityService {
  private readonly scriptId = 'google-identity-services';
  private scriptLoadingPromise: Promise<void> | null = null;

  loadClient(): Promise<void> {
    if (window.google?.accounts?.id) {
      return Promise.resolve();
    }

    if (this.scriptLoadingPromise) {
      return this.scriptLoadingPromise;
    }

    this.scriptLoadingPromise = new Promise<void>((resolve, reject) => {
      const existingScript = document.getElementById(this.scriptId) as HTMLScriptElement | null;
      if (existingScript) {
        existingScript.addEventListener('load', () => resolve(), { once: true });
        existingScript.addEventListener('error', () => reject(new Error('Failed to load Google sign-in.')), { once: true });
        return;
      }

      const script = document.createElement('script');
      script.id = this.scriptId;
      script.src = 'https://accounts.google.com/gsi/client';
      script.async = true;
      script.defer = true;
      script.onload = () => resolve();
      script.onerror = () => reject(new Error('Failed to load Google sign-in.'));

      document.head.appendChild(script);
    });

    return this.scriptLoadingPromise;
  }

  async initialize(
    clientId: string,
    callback: (response: GoogleCredentialResponse) => void,
    context: 'signin' | 'signup' | 'use'
  ): Promise<void> {
    if (!clientId.trim()) {
      throw new Error('Google sign-in is unavailable.');
    }

    await this.loadClient();
    window.google?.accounts.id.initialize({
      client_id: clientId,
      callback,
      context,
      ux_mode: 'popup',
      cancel_on_tap_outside: true,
    });
  }

  renderButton(parent: HTMLElement, width: number, text: GoogleRenderButtonOptions['text']): void {
    window.google?.accounts.id.renderButton(parent, {
      type: 'standard',
      theme: 'outline',
      size: 'medium',
      shape: 'pill',
      text,
      logo_alignment: 'left',
      width,
    });
  }
}
