# QuoteYourOS Backend API Endpoints

Based on the frontend components we've built, here is a list of proposed RESTful API endpoints. They are organized by feature to help transition the site from hardcoded data to a dynamic, database-driven application.

## 📝 1. Blog (Internet Explorer)
Used to fetch blog posts for the `Blog.jsx` component.

*   `GET /api/blog`
    *   **Purpose:** Fetch a list of all blog posts (with pagination).
    *   **Response:** Array of objects `[{ id, title, date, excerpt }]`.
*   `GET /api/blog/:id`
    *   **Purpose:** Fetch the full content of a specific blog post.
    *   **Response:** Object `{ id, title, date, excerpt, content }`.

## 📁 2. Projects (Windows Explorer)
Used to fetch project data for the `Projects.jsx` component.

*   `GET /api/projects`
    *   **Purpose:** Fetch all portfolio projects.
    *   **Response:** Array of objects `[{ id, name, icon, desc, tech, url }]`.
*   *(Optional)* `GET /api/projects/:id`
    *   **Purpose:** Fetch detailed information about a specific project (if you want to expand the project view later).

## 📧 3. Contact (Outlook Express)
Used to handle form submissions from the `Contact.jsx` component.

*   `POST /api/contact`
    *   **Purpose:** Receive a message from a visitor.
    *   **Payload:** `{ from: "email@example.com", subject: "Hello", message: "..." }`
    *   **Action:** Stores the message in the database and/or triggers an email notification to you.

## 📄 4. Static Profile Data (Notepad & WordPad)
Currently, "About Me" and "Resume" are hardcoded. If you want to update them without redeploying the frontend, you can serve them from the backend.

*   `GET /api/profile/about`
    *   **Purpose:** Fetch the raw text for the "About Me" notepad.
*   `GET /api/profile/resume`
    *   **Purpose:** Fetch the structured JSON data for your work experience, skills, and education.
*   `GET /api/profile/resume/download`
    *   **Purpose:** Endpoint to download your actual PDF resume.

---

## 🔒 5. Backoffice / Admin (Protected Endpoints)
To actually *create* the dynamic content, you'll need an admin panel (perhaps hidden behind a specific desktop icon or a secret URL) and authenticated endpoints.

*   **Authentication**
    *   `POST /api/auth/login` (Returns a JWT or sets a session cookie)
    *   `POST /api/auth/logout`
    *   `GET /api/auth/me` (Verify current session)

*   **Content Management (Requires Authentication)**
    *   `POST /api/blog` (Create a new post)
    *   `PUT /api/blog/:id` (Update a post)
    *   `DELETE /api/blog/:id` (Delete a post)
    *   `POST /api/projects` (Add a new project)
    *   `PUT /api/projects/:id` (Edit a project)
    *   `DELETE /api/projects/:id` (Remove a project)
    *   `GET /api/messages` (Read your contact form submissions)
