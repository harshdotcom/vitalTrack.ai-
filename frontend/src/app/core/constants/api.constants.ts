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
  get DOCUMENTS_CALENDAR_URL() {
    return `${this.BASE_URL}/documents/calendar`;
  },
  get FILES_UPLOAD_URL() {
    return `${this.BASE_URL}/files/upload`;
  },
  get DOCUMENTS_URL() {
    return `${this.BASE_URL}/documents`;
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
  get USER_USAGE_URL() {
    return `${this.BASE_URL}/user-details/usage`;
  }
};
