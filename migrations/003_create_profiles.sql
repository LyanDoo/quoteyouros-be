-- Create profiles table for storing profile and resume information
CREATE TABLE IF NOT EXISTS profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    about TEXT,
    resume_file_name VARCHAR(255),
    resume_file_size INTEGER,
    resume_file_path VARCHAR(500),
    resume_mime_type VARCHAR(100),
    resume_uploaded_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create index for faster lookups
CREATE INDEX idx_profiles_created_at ON profiles(created_at DESC);
