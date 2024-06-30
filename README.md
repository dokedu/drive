# Dokedu Drive

### To Do's

- Files
  - Improve renaming UX - Aaron started this but seems to be broken rn
- Shared drives
  - Create
  - Edit sharing
  - Archive
- Trash
  - Restore file

- Settings
  - Account
    - Upload profile picture
    - Change first and last name
  - Admin (only for admins or owners)
    - Users
      - Invite user (via email)
      - Archive user
    - Billing
      - Write support@dokedu.org for more details


#### Technical

**Frontend**
- i18n (translations)
- favicon

**Backend**
- Enforce total file size limit
  `SELECT SUM(size) FROM files WHERE organization_id = 1` must be less than 1TB
