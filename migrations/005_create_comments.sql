-- Create comments table for blog post comments
CREATE TABLE IF NOT EXISTS comments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    blog_post_id UUID NOT NULL REFERENCES blog_posts(id) ON DELETE CASCADE,
    reply_to_comment_id UUID REFERENCES comments(id) ON DELETE CASCADE,
    author_name VARCHAR(255) NOT NULL,
    author_email VARCHAR(255),
    content TEXT NOT NULL,
    rating INT CHECK (rating >= 1 AND rating <= 5),
    is_spam BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for faster lookups
CREATE INDEX idx_comments_blog_post_id ON comments(blog_post_id);
CREATE INDEX idx_comments_reply_to_id ON comments(reply_to_comment_id);
CREATE INDEX idx_comments_is_spam ON comments(is_spam);
CREATE INDEX idx_comments_created_at ON comments(created_at DESC);
