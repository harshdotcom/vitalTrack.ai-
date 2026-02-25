# AI Health Journal – Final API Specification

BASE URL:
`/api/v1`

---

## AUTH APIs

### POST /auth/signup

Request:

```json
{
  "fullName": "Harsh Jha",
  "email": "harsh@email.com",
  "password": "StrongPassword123",
  "age": 28,
  "gender": "Male"
}
```

Response:

```json
{
  "success": true,
  "message": "User created successfully",
  "data": {
    "userId": "USER_ID"
  }
}
```

---

### POST /auth/login

Request:

```json
{
  "email": "harsh@email.com",
  "password": "StrongPassword123"
}
```

Response:

```json
{
  "success": true,
  "data": {
    "token": "JWT_TOKEN",
    "userId": "USER_ID",
    "loginType": "EMAIL"
  }
}
```

---

### POST /auth/google-login

Request:

```json
{
  "googleToken": "GOOGLE_ID_TOKEN"
}
```

Response:

```json
{
  "success": true,
  "data": {
    "token": "JWT_TOKEN",
    "userId": "USER_ID",
    "loginType": "GOOGLE"
  }
}
```

---

### POST /auth/update-user

Request:

```json
{
  "userId": "USER_ID",
  "fullName": "Harsh Jha",
  "age": 29,
  "weight": 72,
  "height": 175
}
```

Response:

```json
{
  "success": true,
  "message": "Profile updated"
}
```

---

## FILE / DOCUMENT APIs

### POST /file/upload

Request (FormData):

```
file: report.pdf
userId: USER_ID
documentType: Blood Test
date: 2026-02-25
```

Response:

```json
{
  "success": true,
  "data": {
    "documentId": "DOC_ID",
    "fileUrl": "storage/url/report.pdf"
  }
}
```

---

### POST /file/get-document-list

Request:

```json
{
  "userId": "USER_ID",
  "page": 1,
  "pageSize": 10
}
```

Response:

```json
{
  "success": true,
  "data": [
    {
      "documentId": "DOC_ID",
      "title": "Blood Report",
      "documentType": "Blood Test",
      "uploadDate": "2026-02-25"
    }
  ]
}
```

---

### POST /file/get-document-details

Request:

```json
{
  "userId": "USER_ID",
  "documentId": "DOC_ID"
}
```

Response:

```json
{
  "success": true,
  "data": {
    "documentId": "DOC_ID",
    "title": "Blood Report",
    "fileUrl": "storage/url/report.pdf",
    "documentType": "Blood Test",
    "uploadedAt": "2026-02-25"
  }
}
```

---

### POST /file/update-document

Request:

```json
{
  "documentId": "DOC_ID",
  "title": "Blood Report Feb",
  "tags": ["blood", "routine"]
}
```

Response:

```json
{
  "success": true,
  "message": "Document updated"
}
```

---

### POST /file/delete-document

Request:

```json
{
  "userId": "USER_ID",
  "documentId": "DOC_ID"
}
```

Response:

```json
{
  "success": true,
  "message": "Document deleted"
}
```

---

## AI APIs

### POST /ai/analyze-document

Request:

```json
{
  "userId": "USER_ID",
  "documentId": "DOC_ID"
}
```

Response:

```json
{
  "success": true,
  "data": {
    "analysisId": "ANALYSIS_ID",
    "summary": "Hemoglobin slightly low",
    "riskLevel": "Medium"
  }
}
```

---

### POST /ai/get-analysis

Request:

```json
{
  "userId": "USER_ID",
  "documentId": "DOC_ID"
}
```

Response:

```json
{
  "success": true,
  "data": {
    "analysisId": "ANALYSIS_ID",
    "summary": "Hemoglobin slightly low",
    "insights": [
      "Increase iron-rich food",
      "Stay hydrated"
    ]
  }
}
```

---

### POST /ai/generate-food-recommendation

Request:

```json
{
  "userId": "USER_ID",
  "analysisId": "ANALYSIS_ID"
}
```

Response:

```json
{
  "success": true,
  "message": "Food recommendation generated"
}
```

---

### POST /ai/get-food-recommendation

Request:

```json
{
  "userId": "USER_ID"
}
```

Response:

```json
{
  "success": true,
  "data": {
    "recommendedFoods": [
      "Spinach",
      "Lentils",
      "Pomegranate"
    ]
  }
}
```

---

## HEALTH JOURNAL API

### POST /user/get-health-timeline

Request:

```json
{
  "userId": "USER_ID",
  "startDate": "2026-02-01",
  "endDate": "2026-02-28"
}
```

Response:

```json
{
  "success": true,
  "data": [
    {
      "date": "2026-02-25",
      "documentId": "DOC_ID",
      "documentType": "Blood Test",
      "summary": "Hemoglobin slightly low"
    }
  ]
}
```

---

## ID FLOW (REFERENCE)

```
Login → USER_ID
Upload Document → DOC_ID
Analyze Document → ANALYSIS_ID
Generate Food Recommendation → USER_ID + ANALYSIS_ID
Timeline → USER_ID + Date Range
```
