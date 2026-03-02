# Vita Track AI — Documents API (v1)

This README documents the Documents APIs currently implemented in the backend.

====================================================================

BASE URL
--------------------------------------------------------------------

http://localhost:8081/api/v1

All endpoints require JWT authentication.

====================================================================

AUTHENTICATION
--------------------------------------------------------------------

Required Headers:

Authorization: Bearer <JWT_TOKEN>
Content-Type: application/json

User ID is extracted from middleware:

userID := c.MustGet("user_id").(int64)

All document operations are scoped to the authenticated user.

====================================================================

DOCUMENT MODEL (REFERENCE)
--------------------------------------------------------------------

type Document struct {
    ID         string
    UserID     int64
    FileID     string
    Category   string
    ReportType string
    FileType   string
    Tags       string
    Status     string
    ReportDate time.Time
}

====================================================================

1) CREATE DOCUMENT
--------------------------------------------------------------------

Endpoint:
POST /documents

Description:
Creates a new document linked to an uploaded file.

Request Payload:

{
  "file_id": "a254c4e2-fa43-422d-8326-3882d1cd01ea",
  "category": "medical",
  "report_type": "blood_test",
  "file_type": "lab_report",
  "tags": ["cbc","vitamin_d"],
  "report_date": "2026-03-15"
}

Field Description:

file_id     (string)  required  -> File ID from file upload API
category    (string)  required  -> Document category
report_type (string)  required  -> Type of report
file_type   (string)  required  -> File classification
tags        (array)   optional  -> Tags for filtering/search
report_date (string)  required  -> Format: YYYY-MM-DD

Important:
report_date must match format:
2006-01-02 (Go time layout)

Successful Response (200):

{
  "document_id": "uuid",
  "status": "uploaded"
}

Validation Errors (400):

{
  "error": "invalid request"
}

or

{
  "error": "invalid report_date format (use YYYY-MM-DD)"
}

Server Error (500):

{
  "error": "failed to create document"
}

====================================================================

2) GET DOCUMENT BY ID
--------------------------------------------------------------------

Endpoint:
GET /documents/:id

Description:
Fetch a single document belonging to the authenticated user.

Example:
GET /documents/8a6a0e8c-xxxx-xxxx-xxxx-xxxxxxxxxxxx

Success Response (200):

{
  "id": "uuid",
  "user_id": 1,
  "file_id": "uuid",
  "category": "medical",
  "report_type": "blood_test",
  "file_type": "lab_report",
  "tags": "[\"cbc\",\"vitamin_d\"]",
  "status": "uploaded",
  "report_date": "2026-03-15T00:00:00Z"
}

Not Found (404):

{
  "error": "document not found"
}

====================================================================

3) DELETE DOCUMENT
--------------------------------------------------------------------

Endpoint:
DELETE /documents/:id

Description:
Deletes a document belonging to the authenticated user.

Example:
DELETE /documents/8a6a0e8c-xxxx-xxxx-xxxx-xxxxxxxxxxxx

Success Response (200):

{
  "message": "document deleted"
}

Not Found (404):

{
  "error": "document not found"
}

====================================================================

4) CALENDAR DOCUMENTS (MONTH VIEW)
--------------------------------------------------------------------

Endpoint:
POST /documents/calendar

Description:
Returns documents grouped by day for a specific month.
Used to build calendar UI.

Request Payload:

{
  "month": 3,
  "year": 2026,
  "category": "medical",
  "tags": ["cbc"]
}

Field Description:

month     (int)    required  -> Month number (1-12)
year      (int)    required  -> Year
category  (string) optional  -> Filter by category
tags      (array)  optional  -> Filter by tags

Backend Logic:

start := time.Date(year, month, 1, 0,0,0,0, time.UTC)
end   := start.AddDate(0,1,0)

Query:

WHERE user_id = ?
AND report_date >= start
AND report_date < end

Success Response (200):

{
  "month": 3,
  "year": 2026,
  "days": {
    "2026-03-15": {
      "count": 2,
      "documents": [
        {
          "id": "uuid",
          "category": "medical",
          "report_type": "blood_test",
          "report_date": "2026-03-15T00:00:00Z"
        }
      ]
    }
  }
}

If no data:

{
  "month": 3,
  "year": 2026,
  "days": {}
}

Error (400):

{
  "error": "invalid request"
}

Error (500):

{
  "error": "failed"
}

====================================================================

COMMON ISSUES / DEBUG NOTES
--------------------------------------------------------------------

1) Empty calendar response:

Cause:
report_date missing or invalid during upload.

Result:
report_date saved as:
0001-01-01 (Go zero time)

Fix:
Always send report_date in payload.

--------------------------------------------------------------------

2) Date parsing bug (IMPORTANT)

Wrong:

parsedDate, _ := time.Parse("2006-01-02", req.ReportDate)

Correct:

parsedDate, err := time.Parse("2006-01-02", req.ReportDate)
if err != nil {
    return error
}

--------------------------------------------------------------------

3) Required date format

YYYY-MM-DD

Example:
2026-03-15

--------------------------------------------------------------------

4) User scoped data

All queries include:

WHERE user_id = current_user

Documents created without user_id will never appear.

====================================================================

TYPICAL API FLOW
--------------------------------------------------------------------

Step 1:
Upload file (File API)
-> returns file_id

Step 2:
Create document
POST /documents
-> save metadata + report_date

Step 3:
Fetch calendar
POST /documents/calendar
-> grouped documents by day

Step 4:
View document
GET /documents/:id

Step 5:
Delete if needed
DELETE /documents/:id

====================================================================

BEST PRACTICES (RECOMMENDED)
--------------------------------------------------------------------

1) Make report_date required.
2) Never ignore time.Parse errors.
3) Use type:date in DB if only date is needed.
4) Always validate file_id exists before document creation.
5) Store tags as JSON string for flexible filtering.

====================================================================

END OF DOCUMENT
====================================================================