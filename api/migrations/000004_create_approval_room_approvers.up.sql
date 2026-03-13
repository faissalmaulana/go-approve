CREATE TABLE IF NOT EXISTS approval_room_approvers(
approval_id VARCHAR(36) NOT NULL,
approval_room_id VARCHAR(36) NOT NULL,
decision ENUM('pending','approved','rejected') DEFAULT 'pending',
FOREIGN KEY (approval_id) REFERENCES users(id)
    ON DELETE CASCADE
    ON UPDATE NO ACTION,
FOREIGN KEY (approval_room_id) REFERENCES approval_rooms(id)
    ON DELETE CASCADE
    ON UPDATE NO ACTION
);
