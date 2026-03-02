import { Injectable, signal } from '@angular/core';

export type ToastType = 'success' | 'error' | 'info';

export interface ToastMessage {
  message: string;
  type: ToastType;
  id: number;
}

@Injectable({
  providedIn: 'root'
})
export class ToastService {
  private counter = 0;
  
  // Expose a read-only signal of the active toasts
  toasts = signal<ToastMessage[]>([]);

  show(message: string, type: ToastType = 'info', durationMs: number = 4000) {
    const id = this.counter++;
    const newToast: ToastMessage = { message, type, id };
    
    // Add to the list
    this.toasts.update(current => [...current, newToast]);

    // Auto remove after duration
    if (durationMs > 0) {
      setTimeout(() => this.remove(id), durationMs);
    }
  }

  showSuccess(message: string, durationMs?: number) {
    this.show(message, 'success', durationMs);
  }

  showError(message: string, durationMs: number = 5000) {
    this.show(message, 'error', durationMs);
  }

  showInfo(message: string, durationMs?: number) {
    this.show(message, 'info', durationMs);
  }

  remove(id: number) {
    this.toasts.update(current => current.filter(t => t.id !== id));
  }
}
