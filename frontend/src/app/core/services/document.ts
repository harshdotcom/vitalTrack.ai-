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

  getInfiniteScrollDocuments(cursor = '', limit = 12): Observable<any> {
    const params: Record<string, string> = {
      limit: String(limit)
    };

    if (cursor) {
      params['cursor'] = cursor;
    }

    return this.http.get(API_CONSTANTS.DOCUMENTS_INFINITE_SCROLL_URL, { params });
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

  updateDocument(id: string, payload: any): Observable<any> {
    let formData = new FormData();
    if (payload instanceof FormData) {
      formData = payload;
    } else {
      Object.keys(payload).forEach(key => {
        if (payload[key] !== null && payload[key] !== undefined) {
          formData.append(key, payload[key]);
        }
      });
    }
    return this.http.patch(`${API_CONSTANTS.UPDATE_DOCUMENT_URL}/${id}`, formData);
  }

  getDocumentDetails(id: string): Observable<any> {
    return this.http.get(`${API_CONSTANTS.DOCUMENTS_URL}/${id}`);
  }

  getFileUrl(fileId: string): Observable<any> {
    return this.http.get(`${API_CONSTANTS.BASE_URL}/files/${fileId}`);
  }

  deleteDocument(id: string): Observable<any> {
    return this.http.delete(`${API_CONSTANTS.FILE_URL}/${id}`);
  }

  deleteHealthMetric(id: string): Observable<any> {
    return this.http.delete(`${API_CONSTANTS.HEALTH_METRIC_URL}/${id}`);
  }

  getAiAnalysis(fileId: string): Observable<any> {
    return this.http.get(`${API_CONSTANTS.FILES_AI_URL}/${fileId}`);
  }

  saveHealthMetric(payload: any): Observable<any> {
    return this.http.post(API_CONSTANTS.HEALTH_METRIC_SAVE_URL, payload);
  }
}
