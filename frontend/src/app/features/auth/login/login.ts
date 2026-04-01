import { Component, inject, ChangeDetectorRef, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormBuilder, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { Router, RouterLink } from '@angular/router';
import { AuthService } from '../../../core/services/auth';
import { ToastService } from '../../../core/services/toast';

@Component({
  selector: 'app-login',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, RouterLink],
  templateUrl: './login.html',
  styleUrl: './login.css',
})
export class Login implements OnInit {
  private fb = inject(FormBuilder);
  private authService = inject(AuthService);
  private router = inject(Router);
  private toastService = inject(ToastService);
  private cdr = inject(ChangeDetectorRef);

  loginForm: FormGroup = this.fb.nonNullable.group({
    email: ['', [Validators.required, Validators.email]],
    password: ['', [Validators.required, Validators.minLength(6)]]
  });

  isLoading = false;

  ngOnInit(): void {
    const state = history.state;
    if (state?.email) {
      this.loginForm.patchValue({ email: state.email });
    }
    if (state?.password) {
      this.loginForm.patchValue({ password: state.password });
    }
  }
  errorMessage = '';
  showPassword = false;

  togglePasswordVisibility() {
    this.showPassword = !this.showPassword;
  }

  onSubmit() {
    if (this.loginForm.invalid) {
      this.loginForm.markAllAsTouched();
      return;
    }

    this.isLoading = true;
    this.errorMessage = '';

    this.authService.login(this.loginForm.getRawValue()).subscribe({
      next: (response) => {
        // Handle success (e.g., store token, navigate)
        console.log('Login successful', response);
        this.router.navigate(['/dashboard']);
        this.isLoading = false;
        this.cdr.detectChanges();
      },
      error: (err) => {
        console.error('Login Error:', err);
        const errMsg = err.error?.message || 'Invalid email or password. Please try again.';
        this.errorMessage = errMsg;
        this.toastService.showError(errMsg);
        this.isLoading = false;
        this.cdr.detectChanges();
      }
    });
  }
}
