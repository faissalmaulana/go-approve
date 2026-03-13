CREATE TABLE IF NOT EXISTS approval_rooms (
    id           VARCHAR(36) PRIMARY KEY NOT NULL,
    title        VARCHAR(255) NOT NULL,
    filepaths    MEDIUMTEXT NOT NULL,
    due_at       DATETIME NOT NULL,
    submitter_id VARCHAR(36),
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP NULL ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (submitter_id) REFERENCES users(id)
        ON DELETE CASCADE
        ON UPDATE NO ACTION
);
