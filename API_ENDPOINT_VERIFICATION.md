# API Endpoint Verification

## Endpoints dari main.go (47 endpoints total):

### Health Check (1 endpoint)
1. `GET /health` - Health check

### Public Endpoints (4 endpoints)
2. `GET /api/terms-conditions` - Get terms & conditions
3. `GET /api/privacy-policy` - Get privacy policy
4. `GET /api/onboardings` - Get all onboardings (public)
5. `GET /api/onboardings/:id` - Get onboarding by ID (public)

### Authentication Endpoints (3 endpoints)
6. `POST /api/login` - Banking login step 1 (send OTP)
7. `POST /api/login/verify` - Banking login step 2 (verify OTP)
8. `POST /api/refresh` - Refresh access token

### Admin Authentication (2 endpoints)
9. `POST /api/admin/login` - Admin login
10. `POST /api/admin/logout` - Admin logout

### Admin Management (5 endpoints)
11. `GET /api/admin/admins` - Get all admins
12. `GET /api/admin/admins/:id` - Get admin by ID
13. `POST /api/admin/admins` - Create new admin
14. `PUT /api/admin/admins/:id` - Update admin
15. `DELETE /api/admin/admins/:id` - Delete admin

### User Profile Management (3 endpoints)
16. `GET /api/profile` - Get user profile
17. `PUT /api/profile` - Update user profile
18. `PUT /api/change-pin` - Change PIN ATM

### Session Management (3 endpoints)
19. `GET /api/sessions` - Get active sessions
20. `POST /api/logout` - Logout
21. `POST /api/logout-others` - Logout other sessions

### Article Management (6 endpoints)
22. `GET /api/articles` - Get all articles
23. `GET /api/articles/:id` - Get article by ID
24. `PUT /api/articles/:id` - Update article
25. `DELETE /api/articles/:id` - Delete article
26. `GET /api/my-articles` - Get my articles
27. `POST /api/articles` - Create article

### Photo Management (5 endpoints)
28. `GET /api/photos` - Get all photos
29. `GET /api/photos/:id` - Get photo by ID
30. `PUT /api/photos/:id` - Update photo
31. `DELETE /api/photos/:id` - Delete photo
32. `POST /api/photos` - Create photo

### Bank Account Management (5 endpoints)
33. `GET /api/bank-accounts` - Get user's bank accounts
34. `POST /api/bank-accounts` - Create new bank account
35. `PUT /api/bank-accounts/:id` - Update bank account
36. `DELETE /api/bank-accounts/:id` - Delete bank account
37. `PUT /api/bank-accounts/:id/primary` - Set primary account

### Content Management (5 endpoints)
38. `POST /api/terms-conditions` - Set terms & conditions
39. `POST /api/privacy-policy` - Set privacy policy
40. `POST /api/onboardings` - Create onboarding
41. `PUT /api/onboardings/:id` - Update onboarding
42. `DELETE /api/onboardings/:id` - Delete onboarding

### User Management (5 endpoints)
43. `GET /api/users` - List all users
44. `GET /api/users/:id` - Get user by ID
45. `DELETE /api/users/:id` - Delete user by ID
46. `POST /api/users` - Create user
47. `PUT /api/users/:id` - Update user by ID

### Config Management (4 endpoints)
48. `POST /api/config` - Set config value
49. `GET /api/configs` - Get all configs
50. `DELETE /api/config/:key` - Delete config by key
51. `GET /api/config/:key` - Get config value by key

## PERHATIAN:
Saya menghitung 51 endpoint dalam analisis ini, tetapi analisis sebelumnya menyebutkan 47. 
Perlu verifikasi ulang untuk memastikan hitungan yang akurat.
