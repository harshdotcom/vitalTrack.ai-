import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { API_CONSTANTS } from '../constants/api.constants';

@Injectable({
  providedIn: 'root'
})
export class DocumentService {
  private http = inject(HttpClient);

  getMonthlyReports(month: number, year: number): Observable<any> {
    return this.http.post(API_CONSTANTS.DOCUMENTS_CALENDAR_URL, { month, year });
  }

  uploadFile(file: File, fileType: string): Observable<any> {
    const formData = new FormData();
    formData.append('files', file);
    formData.append('file_type', fileType);
    return this.http.post(API_CONSTANTS.FILES_UPLOAD_URL, formData);
  }

  submitDocument(details: any): Observable<any> {
    return this.http.post(API_CONSTANTS.DOCUMENTS_URL, details);
  }

  getDocumentDetails(id: string): Observable<any> {
    return this.http.get(`${API_CONSTANTS.DOCUMENTS_URL}/${id}`);
  }

  getFileUrl(fileId: string): Observable<any> {
    return this.http.get(`${API_CONSTANTS.BASE_URL}/files/${fileId}`);
  }

  deleteDocument(id: string): Observable<any> {
    return this.http.delete(`${API_CONSTANTS.DOCUMENTS_URL}/${id}`);
  }
}
