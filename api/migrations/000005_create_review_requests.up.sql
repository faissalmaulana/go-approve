CREATE TABLE IF NOT EXISTS review_requests (
    id               VARCHAR(36) PRIMARY KEY NOT NULL,
    is_read             BOOL DEFAULT FALSE NOT NULL,
    status           ENUM('pending', 'accepted', 'rejected') DEFAULT 'pending',
    approval_room_id VARCHAR(36) NOT NULL,
    invitee_id       VARCHAR(36) NOT NULL,
    requester_id     VARCHAR(36) NOT NULL,
    created_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (approval_room_id) REFERENCES approval_rooms(id)
        ON DELETE CASCADE
        ON UPDATE NO ACTION,
    FOREIGN KEY (invitee_id) REFERENCES users(id)
        ON DELETE CASCADE
        ON UPDATE NO ACTION,
    FOREIGN KEY (requester_id) REFERENCES users(id)
        ON DELETE CASCADE
        ON UPDATE NO ACTION
);
