# Dokedu Drive

### To Do's

- Files
  - Add "New folder" button
  - Improve renaming UX
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
- context menu (make right click work for full area, not only file list coz if no files, no context menu)

**Backend**
- Ensure image previews are cached by providing a cache-control header in the response
- Enforce total file size limit
  `SELECT SUM(size) FROM files WHERE organization_id = 1` must be less than 1TB
