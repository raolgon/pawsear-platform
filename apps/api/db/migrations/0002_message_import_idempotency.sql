CREATE UNIQUE INDEX idx_messages_conversation_external_unique
ON messages(conversation_id, external_message_id)
WHERE external_message_id IS NOT NULL;
