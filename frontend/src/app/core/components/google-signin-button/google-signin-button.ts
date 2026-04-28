import {
  AfterViewInit,
  ChangeDetectionStrategy,
  ChangeDetectorRef,
  Component,
  ElementRef,
  EventEmitter,
  Input,
  NgZone,
  OnChanges,
  Output,
  SimpleChanges,
  ViewChild,
  inject,
} from '@angular/core';
import { CommonModule } from '@angular/common';
import { GoogleCredentialResponse, GoogleIdentityService } from '../../services/google-identity';

@Component({
  selector: 'app-google-signin-button',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './google-signin-button.html',
  styleUrl: './google-signin-button.css',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class GoogleSigninButtonComponent implements AfterViewInit, OnChanges {
  private readonly googleIdentityService = inject(GoogleIdentityService);
  private readonly cdr = inject(ChangeDetectorRef);
  private readonly ngZone = inject(NgZone);

  @Input() clientId = '';
  @Input() context: 'signin' | 'signup' | 'use' = 'signin';
  @Input() disabled = false;
  @Output() credentialReceived = new EventEmitter<string>();
  @Output() initializationFailed = new EventEmitter<string>();
  @ViewChild('buttonHost') private buttonHost?: ElementRef<HTMLDivElement>;

  protected isInitializing = false;
  protected errorMessage = '';
  private viewInitialized = false;

  ngAfterViewInit(): void {
    this.viewInitialized = true;
    void this.prepareButton();
  }

  ngOnChanges(changes: SimpleChanges): void {
    if (!this.viewInitialized) {
      return;
    }

    if (changes['clientId'] || changes['context']) {
      void this.prepareButton();
    }
  }

  private async prepareButton(): Promise<void> {
    const host = this.buttonHost?.nativeElement;
    if (!host) {
      return;
    }

    host.innerHTML = '';
    this.errorMessage = '';

    if (!this.clientId.trim()) {
      this.cdr.markForCheck();
      return;
    }

    this.isInitializing = true;
    this.cdr.markForCheck();

    try {
      await this.googleIdentityService.initialize(
        this.clientId,
        (response) => this.ngZone.run(() => this.handleCredentialResponse(response)),
        this.context
      );

      const width = Math.min(Math.max(host.clientWidth || 260, 260), 320);
      this.googleIdentityService.renderButton(
        host,
        width,
        this.context === 'signup' ? 'signup_with' : 'continue_with'
      );
    } catch (error) {
      this.errorMessage = error instanceof Error ? error.message : 'Google sign-in is currently unavailable.';
      this.initializationFailed.emit(this.errorMessage);
    } finally {
      this.isInitializing = false;
      this.cdr.markForCheck();
    }
  }

  private handleCredentialResponse(response: GoogleCredentialResponse): void {
    if (!response.credential) {
      this.initializationFailed.emit('Google sign-in did not return a credential.');
      return;
    }

    this.credentialReceived.emit(response.credential);
  }
}
