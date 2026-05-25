# Gallery API Endpoints Documentation

This document describes the endpoints for the **NFT Gallery (Windows Picture and Fax Viewer)**. They allow public visitors to list gallery items and view photos, and allow administrators to manage items.

---

## 🖼️ Endpoints Summary

| Method | Endpoint | Description | Auth Required |
| :--- | :--- | :--- | :--- |
| **GET** | `/api/gallery` | Retrieve all gallery items | No |
| **GET** | `/api/gallery/images/:filename` | Serve the raw binary of a gallery photo | No |
| **POST** | `/api/gallery` | Add a new NFT photo to the gallery | **Yes (Admin)** |
| **PUT** | `/api/gallery/:id` | Update metadata and/or replace image of an item | **Yes (Admin)** |
| **DELETE**| `/api/gallery/:id` | Delete a gallery item and its image file | **Yes (Admin)** |

---

## 🔓 1. Get All Gallery Items

Fetch the list of all uploaded NFT photographs.

*   **URL:** `/api/gallery`
*   **Method:** `GET`
*   **Auth Required:** None

### Success Response
*   **Code:** `200 OK`
*   **Payload (JSON):**
    ```json
    {
      "success": true,
      "message": "Gallery items retrieved successfully",
      "data": [
        {
          "id": "270db65a-063a-4ef6-ac16-2f0857ef1960",
          "title": "Neon Sunset",
          "description": "Captured during sunset at the beach, featuring vibrant pink and blue colors.",
          "author": "John Doe",
          "image": "/api/gallery/images/1716670400_neon_sunset.png",
          "created_at": "2026-05-25 22:45:00",
          "updated_at": "2026-05-25 22:45:00"
        }
      ]
    }
    ```

> [!TIP]
> Use the relative path returned in `image` prepended with your backend host URL (e.g. `http://localhost:8000/api/gallery/images/1716670400_neon_sunset.png`) to render it in an `<img>` tag.

---

## 🔓 2. Get Gallery Image

Retrieve the raw binary image stream to render inside the browser (`PhotoViewer.jsx` or similar).

*   **URL:** `/api/gallery/images/:filename`
*   **Method:** `GET`
*   **Auth Required:** None
*   **URL Params:**
    *   `filename` (string, required) - e.g. `1716670400_neon_sunset.png`

### Success Response
*   **Code:** `200 OK`
*   **Content-Type:** Matches image format (`image/png`, `image/jpeg`, `image/gif`, `image/webp`)
*   **Payload:** Raw binary image file.

### Error Response
*   **Code:** `404 Not Found` if the file doesn't exist.

---

## 🔒 3. Create Gallery Item

Upload a new photograph to the NFT gallery.

*   **URL:** `/api/gallery`
*   **Method:** `POST`
*   **Auth Required:** Yes
*   **Headers:**
    *   `Authorization: Bearer <jwt_token>`
    *   `Content-Type: multipart/form-data`
*   **Request Body (form-data):**
    *   `title` (string, required): The name of the NFT.
    *   `description` (string, required): Background story/attributes.
    *   `author` (string, required): The name of the artist/creator.
    *   `image` (file, required): The binary image file to upload. Allowed formats: PNG, JPG/JPEG, GIF, WebP (max 20MB).

### Success Response
*   **Code:** `210 Created`
*   **Payload (JSON):**
    ```json
    {
      "success": true,
      "message": "Gallery item created successfully",
      "data": {
        "id": "270db65a-063a-4ef6-ac16-2f0857ef1960",
        "title": "Neon Sunset",
        "description": "Captured during sunset at the beach, featuring vibrant pink and blue colors.",
        "author": "John Doe",
        "image": "/api/gallery/images/1716670400_neon_sunset.png",
        "created_at": "2026-05-25 22:45:00",
        "updated_at": "2026-05-25 22:45:00"
      }
    }
    ```

---

## 🔒 4. Update Gallery Item

Update the details of an existing gallery item. You can optionally replace the image file as well.

*   **URL:** `/api/gallery/:id`
*   **Method:** `PUT`
*   **Auth Required:** Yes
*   **Headers:**
    *   `Authorization: Bearer <jwt_token>`
    *   `Content-Type: multipart/form-data`
*   **URL Params:**
    *   `id` (string, required) - The UUID of the gallery item.
*   **Request Body (form-data):**
    *   `title` (string, optional): Updated name of the NFT.
    *   `description` (string, optional): Updated background story.
    *   `author` (string, optional): Updated name of the artist/creator.
    *   `image` (file, optional): A new binary image file to replace the old one (max 20MB). If omitted, the old image is kept.

### Success Response
*   **Code:** `200 OK`
*   **Payload (JSON):**
    ```json
    {
      "success": true,
      "message": "Gallery item updated successfully",
      "data": {
        "id": "270db65a-063a-4ef6-ac16-2f0857ef1960",
        "title": "Updated Title",
        "description": "Updated description details.",
        "author": "Jane Doe",
        "image": "/api/gallery/images/1716679999_new_sunset.png",
        "created_at": "2026-05-25 22:45:00",
        "updated_at": "2026-05-25 22:50:00"
      }
    }
    ```

---

## 🔒 5. Delete Gallery Item

Remove an item from the gallery. This deletes the database entry and cleans up the associated image file from disk.

*   **URL:** `/api/gallery/:id`
*   **Method:** `DELETE`
*   **Auth Required:** Yes
*   **Headers:**
    *   `Authorization: Bearer <jwt_token>`
*   **URL Params:**
    *   `id` (string, required) - The UUID of the gallery item.

### Success Response
*   **Code:** `200 OK`
*   **Payload (JSON):**
    ```json
    {
      "success": true,
      "message": "Gallery item deleted successfully",
      "data": {}
    }
    ```
