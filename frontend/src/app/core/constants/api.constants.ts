import { environment } from '../../../environments/environment';

export const API_CONSTANTS = {
  BASE_URL: environment.apiUrl,
  // when you create the build please paste below in the above line
  // http://43.205.212.220:8081/api/v1
  get LOGIN_URL() {
    return `${this.BASE_URL}/users/login`;
  },
  get SIGNUP_URL() {
    return `${this.BASE_URL}/users/signup`;
  },
  get GOOGLE_LOGIN_URL() {
    return `${this.BASE_URL}/users/google`;
  },
  get DOCUMENTS_CALENDAR_URL() {
    return `${this.BASE_URL}/documents/calendar`;
  },
  get DOCUMENTS_INFINITE_SCROLL_URL() {
    return `${this.BASE_URL}/documents/infiniteScroll`;
  },
  get FILES_UPLOAD_URL() {
    return `${this.BASE_URL}/files/upload`;
  },
  get DOCUMENTS_URL() {
    return `${this.BASE_URL}/documents`;
  },
  get UPDATE_DOCUMENT_URL() {
    return `${this.BASE_URL}/documents/update`;
  },
  get FILE_URL() {
    return `${this.BASE_URL}/files`;
  },
  get FILES_AI_URL() {
    return `${this.BASE_URL}/files/ai`;
  },
  get VERIFY_OTP_URL() {
    return `${this.BASE_URL}/users/verify-otp`;
  },
  get RESEND_OTP_URL() {
    return `${this.BASE_URL}/users/resend-otp`;
  },
  get FORGOT_PASSWORD_URL() {
    return `${this.BASE_URL}/users/forgot-password`;
  },
  get RESET_PASSWORD_URL() {
    return `${this.BASE_URL}/users/reset-password`;
  },
  get USER_USAGE_URL() {
    return `${this.BASE_URL}/user-details/usage`;
  },
  get AI_CREDITS_URL() {
    return `${this.BASE_URL}/user-details/ai-credits`;
  },
  get UPDATE_PROFILE_URL() {
    return `${this.BASE_URL}/user-details/update`;
  },
  get HEALTH_METRIC_SAVE_URL() {
    return `${this.BASE_URL}/health-metric/save`;
  },
  get HEALTH_METRIC_URL() {
    return `${this.BASE_URL}/health-metric`;
  }
};
