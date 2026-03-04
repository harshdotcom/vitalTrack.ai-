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
}
```

Response:

```json
{
    "users": {
        "user_id": 39,
        "email": "john@gmail.com",
        "password": "$2a$12$wH6JBdlTpvmj67kzOu867O8H41r7Kc1CUoW5FyzphYZQn5GZPI3uq",
        "GoogleId": null,
        "name": "John Doe",
        "age": null,
        "gender": "",
        "profile_pic": null,
        "CreatedAt": "2026-03-05T02:16:08.562392+05:30",
        "UpdatedAt": "2026-03-05T02:16:08.562392+05:30"
    }
}
```

---

### POST /auth/login

Request:

```json
{
    "email": "mikasa@gmail.com",
    "password": "123456789"
}
```

Response:

```json
{
    "message": "Successfully login",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6Im1pa2FzYUBnbWFpbC5jb20iLCJleHAiOjE3NzI2NjE2MzgsInVzZXJJZCI6MzZ9.R1M7xjJGezETqKNvQPKBLn5wZ3ejFLCLFXbuzpieh1Y",
    "user": {
        "user_id": 36,
        "email": "mikasa@gmail.com",
        "password": "$2a$12$eTDWJf3e44Qa9cjcAk18Q.ca2oSEOu/3oArO0M9bzzWuYd3U1tSNK",
        "GoogleId": null,
        "name": "Mikasa Ackerman",
        "age": null,
        "gender": "",
        "profile_pic": null,
        "CreatedAt": "2026-03-02T15:50:12.765834+05:30",
        "UpdatedAt": "2026-03-02T15:50:12.765834+05:30"
    }
}
```

---

### POST /auth/google-login [Needs to be done]

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

### POST /auth/update-user [Needs to be done]

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

### POST /files/upload

Request (FormData):

```
files: report.pdf

```

Response:

```json
{
    "files": [
        {
            "file_id": "a000c1d5-5bc7-40e2-8d32-962a99127832",
            "original_name": "Screenshot_22.png"
        }
    ]
}
```

### GET /files/:id

Request (Params):

```
id: 8044891e-e5cf-44fc-bdc4-2b3e8f7efcfb

```

Response:

```json
{
    "url": "https://vitatrack-documents-dev.s3.ap-south-1.amazonaws.com/ba4dcca3-05d9-410e-8122-a9783668bce0.jpg?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Checksum-Mode=ENABLED&X-Amz-Credential=AKIAWLX2NNJXXXTUF5PS%2F20260304%2Fap-south-1%2Fs3%2Faws4_request&X-Amz-Date=20260304T205954Z&X-Amz-Expires=300&X-Amz-SignedHeaders=host&x-id=GetObject&X-Amz-Signature=71a4d3d723eb10526f18aafe626197da12e9ad7c06c99d0178d387b00033c43e"
}
```

---

### POST /file/get-document-list [Needs to done]

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

### GET /documents/:documentId

Request (Params):

```
id: 0e23fd6b-f6ad-42b4-a676-3682ab20b8cf

```

Response:

```json
{
    "id": "0e23fd6b-f6ad-42b4-a676-3682ab20b8cf",
    "user_id": 36,
    "file_id": "7d94725a-cba7-46f9-b04c-2fe7bf7f56bb",
    "category": "misc",
    "report_type": "nothing",
    "file_type": "pdf",
    "tags": "[\"no tag\"]",
    "status": "uploaded",
    "report_date": "2026-02-25T05:30:00+05:30"
}
```

---

### POST /documents

Request:

```json
   {
        "file_id":"7d94725a-cba7-46f9-b04c-2fe7bf7f56bb",
        "category": "misc",
        "report_type":"nothing",
        "file_type":"pdf",
        "tags": ["no tag"],
        "report_date": "2026-02-25"
    }
```

Response:

```json
{
    "document_id": "0e23fd6b-f6ad-42b4-a676-3682ab20b8cf",
    "status": "uploaded"
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
